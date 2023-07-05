package ccip

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cache"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/hasher"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// exec Report should make sure to cap returned payload to this limit
const MaxExecutionReportLength = 250_000

var (
	_ types.ReportingPluginFactory = &ExecutionReportingPluginFactory{}
	_ types.ReportingPlugin        = &ExecutionReportingPlugin{}
)

type ExecutionPluginConfig struct {
	lggr                  logger.Logger
	sourceLP, destLP      logpoller.LogPoller
	onRamp                evm_2_evm_onramp.EVM2EVMOnRampInterface
	offRamp               evm_2_evm_offramp.EVM2EVMOffRampInterface
	commitStore           commit_store.CommitStoreInterface
	srcPriceRegistry      price_registry.PriceRegistryInterface
	srcWrappedNativeToken common.Address
	destClient            evmclient.Client
	destGasEstimator      gas.EvmFeeEstimator
	leafHasher            hasher.LeafHasherInterface[[32]byte]
}

type ExecutionReportingPlugin struct {
	config             ExecutionPluginConfig
	F                  int
	lggr               logger.Logger
	inflightReports    *inflightExecReportsContainer
	snoozedRoots       map[[32]byte]time.Time
	destPriceRegistry  price_registry.PriceRegistryInterface
	destWrappedNative  common.Address
	onchainConfig      ccipconfig.ExecOnchainConfig
	offchainConfig     ccipconfig.ExecOffchainConfig
	cachedSrcFeeTokens *cache.CachedChain[[]common.Address]
	cachedDstTokens    *cache.CachedChain[cache.CachedTokens]
}

type ExecutionReportingPluginFactory struct {
	config ExecutionPluginConfig

	// We keep track of the registered filters
	srcChainFilters []logpoller.Filter
	dstChainFilters []logpoller.Filter
	filtersMu       *sync.Mutex
}

func NewExecutionReportingPluginFactory(config ExecutionPluginConfig) *ExecutionReportingPluginFactory {
	return &ExecutionReportingPluginFactory{
		config:    config,
		filtersMu: &sync.Mutex{},
	}
}

func (rf *ExecutionReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	ctx := context.TODO()

	onchainConfig, err := abihelpers.DecodeAbiStruct[ccipconfig.ExecOnchainConfig](config.OnchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	offchainConfig, err := ccipconfig.DecodeOffchainConfig[ccipconfig.ExecOffchainConfig](config.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	priceRegistry, err := observability.NewObservedPriceRegistry(onchainConfig.PriceRegistry, ExecPluginLabel, rf.config.destClient)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	destRouter, err := router.NewRouter(onchainConfig.Router, rf.config.destClient)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	destWrappedNative, err := destRouter.GetWrappedNative(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	if err = rf.UpdateLogPollerFilters(onchainConfig.PriceRegistry); err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	cachedSrcFeeTokens := cache.NewCachedFeeTokens(rf.config.sourceLP, rf.config.srcPriceRegistry, int64(offchainConfig.SourceFinalityDepth))
	cachedDstTokens := cache.NewCachedSupportedTokens(rf.config.destLP, rf.config.offRamp, priceRegistry, int64(offchainConfig.DestOptimisticConfirmations))
	rf.config.lggr.Infow("Starting exec plugin",
		"offchainConfig", offchainConfig,
		"onchainConfig", onchainConfig)

	return &ExecutionReportingPlugin{
			config:             rf.config,
			F:                  config.F,
			lggr:               rf.config.lggr.Named("ExecutionReportingPlugin"),
			snoozedRoots:       make(map[[32]byte]time.Time),
			inflightReports:    newInflightExecReportsContainer(offchainConfig.InflightCacheExpiry.Duration()),
			destPriceRegistry:  priceRegistry,
			destWrappedNative:  destWrappedNative,
			onchainConfig:      onchainConfig,
			offchainConfig:     offchainConfig,
			cachedDstTokens:    cachedDstTokens,
			cachedSrcFeeTokens: cachedSrcFeeTokens,
		}, types.ReportingPluginInfo{
			Name: "CCIPExecution",
			// Setting this to false saves on calldata since OffRamp doesn't require agreement between NOPs
			// (OffRamp is only able to execute committed messages).
			UniqueReports: false,
			Limits: types.ReportingPluginLimits{
				MaxObservationLength: MaxObservationLength,
				MaxReportLength:      MaxExecutionReportLength,
			},
		}, nil
}

func (r *ExecutionReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ExecutionReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	lggr := r.lggr.Named("ExecutionObservation")
	if isCommitStoreDownNow(ctx, lggr, r.config.commitStore) {
		return nil, ErrCommitStoreIsDown
	}
	// Expire any inflight reports.
	r.inflightReports.expire(lggr)
	inFlight := r.inflightReports.getAll()

	observationBuildStart := time.Now()
	// IMPORTANT: We build executable set based on the leaders token prices, ensuring consistency across followers.
	executableObservations, err := r.getExecutableObservations(ctx, lggr, timestamp, inFlight)
	measureObservationBuildDuration(timestamp, time.Since(observationBuildStart))
	if err != nil {
		return nil, err
	}
	// cap observations which fits MaxObservationLength (after serialized)
	capped := sort.Search(len(executableObservations), func(i int) bool {
		var encoded []byte
		encoded, err = NewExecutionObservation(executableObservations[:i+1]).Marshal()
		if err != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > MaxObservationLength
	})
	if err != nil {
		return nil, err
	}
	executableObservations = executableObservations[:capped]
	lggr.Infow("Observation", "executableMessages", executableObservations)
	// Note can be empty
	return NewExecutionObservation(executableObservations).Marshal()
}

// UpdateLogPollerFilters updates the log poller filters for the source and destination chains.
// pass zeroAddress if dstPriceRegistry is unknown, filters with zero address are omitted.
func (rf *ExecutionReportingPluginFactory) UpdateLogPollerFilters(dstPriceRegistry common.Address) error {
	rf.filtersMu.Lock()
	defer rf.filtersMu.Unlock()

	// source chain filters
	srcFiltersBefore, srcFiltersNow := rf.srcChainFilters, getExecutionPluginSourceLpChainFilters(rf.config.onRamp.Address(), rf.config.srcPriceRegistry.Address())
	created, deleted := filtersDiff(srcFiltersBefore, srcFiltersNow)
	if err := unregisterLpFilters(nilQueryer, rf.config.sourceLP, deleted); err != nil {
		return err
	}
	if err := registerLpFilters(nilQueryer, rf.config.sourceLP, created); err != nil {
		return err
	}
	rf.srcChainFilters = srcFiltersNow

	// destination chain filters
	dstFiltersBefore, dstFiltersNow := rf.dstChainFilters, getExecutionPluginDestLpChainFilters(rf.config.commitStore.Address(), rf.config.offRamp.Address(), dstPriceRegistry)
	created, deleted = filtersDiff(dstFiltersBefore, dstFiltersNow)
	if err := unregisterLpFilters(nilQueryer, rf.config.destLP, deleted); err != nil {
		return err
	}
	if err := registerLpFilters(nilQueryer, rf.config.destLP, created); err != nil {
		return err
	}
	rf.dstChainFilters = dstFiltersNow

	return nil
}

// helper struct to hold the send request and some metadata
type evm2EVMOnRampCCIPSendRequestedWithMeta struct {
	evm_2_evm_offramp.InternalEVM2EVMMessage
	blockTimestamp time.Time
	executed       bool
	finalized      bool
}

func (r *ExecutionReportingPlugin) getExecutableObservations(ctx context.Context, lggr logger.Logger, timestamp types.ReportTimestamp, inflight []InflightInternalExecutionReport) ([]ObservedMessage, error) {
	unexpiredReports, err := getUnexpiredCommitReports(
		ctx,
		r.config.destLP,
		r.config.commitStore,
		r.onchainConfig.PermissionLessExecutionThresholdDuration(),
	)
	if err != nil {
		return nil, err
	}
	lggr.Infow("Unexpired roots", "n", len(unexpiredReports))
	if len(unexpiredReports) == 0 {
		return []ObservedMessage{}, nil
	}

	// This could result in slightly different values on each call as
	// the function returns the allowed amount at the time of the last block.
	// Since this will only increase over time, the highest observed value will
	// always be the lower bound of what would be available on chain
	// since we already account for inflight txs.
	allowedTokenAmount := LazyFetch(func() (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
		return r.config.offRamp.CurrentRateLimiterState(&bind.CallOpts{Context: ctx})
	})
	srcToDstTokens, supportedDestTokens, err := r.sourceDestinationTokens(ctx)
	if err != nil {
		return nil, err
	}
	srcTokensPrices := LazyFetch(func() (map[common.Address]*big.Int, error) {
		srcFeeTokens, err1 := r.cachedSrcFeeTokens.Get(ctx)
		if err1 != nil {
			return nil, err1
		}
		return getTokensPrices(ctx, srcFeeTokens, r.config.srcPriceRegistry, []common.Address{r.config.srcWrappedNativeToken})
	})
	destTokensPrices := LazyFetch(func() (map[common.Address]*big.Int, error) {
		dstTokens, err1 := r.cachedDstTokens.Get(ctx)
		if err1 != nil {
			return nil, err1
		}
		return getTokensPrices(ctx, dstTokens.FeeTokens, r.destPriceRegistry, append(supportedDestTokens, r.destWrappedNative))
	})
	destGasPrice := LazyFetch(func() (*big.Int, error) {
		return r.estimateDestinationGasPrice(ctx)
	})

	lggr.Infow("Processing unexpired reports", "n", len(unexpiredReports))
	measureNumberOfReportsProcessed(timestamp, len(unexpiredReports))
	reportIterationStart := time.Now()
	defer func() {
		measureReportsIterationDuration(timestamp, time.Since(reportIterationStart))
	}()

	unexpiredReportsWithSendReqs, err := r.getReportsWithSendRequests(ctx, unexpiredReports)
	if err != nil {
		return nil, err
	}

	for _, rep := range unexpiredReportsWithSendReqs {
		if ctx.Err() != nil {
			lggr.Warn("Processing of roots killed by context")
			break
		}

		merkleRoot := rep.commitReport.MerkleRoot

		rootLggr := lggr.With("root", hexutil.Encode(merkleRoot[:]),
			"minSeqNr", rep.commitReport.Interval.Min,
			"maxSeqNr", rep.commitReport.Interval.Max,
		)

		snoozeUntil, haveSnoozed := r.snoozedRoots[merkleRoot]
		if haveSnoozed && time.Now().Before(snoozeUntil) {
			rootLggr.Debug("Skipping snoozed root")
			continue
		}

		if err := rep.validate(); err != nil {
			rootLggr.Errorw("Skipping invalid report", "err", err)
			continue
		}

		// If all messages are already executed and finalized, snooze the root for
		// config.PermissionLessExecutionThresholdSeconds so it will never be considered again.
		if allMsgsExecutedAndFinalized := rep.allRequestsAreExecutedAndFinalized(); allMsgsExecutedAndFinalized {
			rootLggr.Infof("Snoozing root %s forever since there are no executable txs anymore", hex.EncodeToString(merkleRoot[:]))
			r.snoozedRoots[merkleRoot] = time.Now().Add(r.onchainConfig.PermissionLessExecutionThresholdDuration())
			incSkippedRequests(reasonAllExecuted)
			continue
		}

		blessed, err := r.config.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, merkleRoot)
		if err != nil {
			return nil, err
		}
		if !blessed {
			rootLggr.Infow("Report is accepted but not blessed")
			incSkippedRequests(reasonNotBlessed)
			continue
		}

		allowedTokenAmountValue, err := allowedTokenAmount()
		if err != nil {
			return nil, err
		}
		srcTokensPricesValue, err := srcTokensPrices()
		if err != nil {
			return nil, err
		}

		destTokensPricesValue, err := destTokensPrices()
		if err != nil {
			return nil, err
		}

		buildBatchDuration := time.Now()
		batch := r.buildBatch(rootLggr, rep, inflight, allowedTokenAmountValue.Tokens,
			srcTokensPricesValue, destTokensPricesValue, destGasPrice, srcToDstTokens)
		measureBatchBuildDuration(timestamp, time.Since(buildBatchDuration))
		if len(batch) != 0 {
			return batch, nil
		}
		r.snoozedRoots[merkleRoot] = time.Now().Add(r.offchainConfig.RootSnoozeTime.Duration())
	}
	return []ObservedMessage{}, nil
}

func (r *ExecutionReportingPlugin) estimateDestinationGasPrice(ctx context.Context) (*big.Int, error) {
	destGasPriceWei, _, err := r.config.destGasEstimator.GetFee(ctx, nil, 0, assets.NewWei(big.NewInt(int64(r.offchainConfig.MaxGasPrice))))
	if err != nil {
		return nil, errors.Wrap(err, "could not estimate destination gas price")
	}
	destGasPrice := destGasPriceWei.Legacy.ToInt()
	if destGasPriceWei.DynamicFeeCap != nil {
		destGasPrice = destGasPriceWei.DynamicFeeCap.ToInt()
	}
	return destGasPrice, nil
}

func (r *ExecutionReportingPlugin) sourceDestinationTokens(ctx context.Context) (map[common.Address]common.Address, []common.Address, error) {
	dstTokens, err := r.cachedDstTokens.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	srcToDstTokens := dstTokens.SupportedTokens
	supportedDestTokens := make([]common.Address, 0, len(srcToDstTokens))
	for _, destToken := range srcToDstTokens {
		supportedDestTokens = append(supportedDestTokens, destToken)
	}
	return srcToDstTokens, supportedDestTokens, nil
}

// Calculates a map that indicated whether a sequence number has already been executed
// before. It doesn't matter if the executed succeeded, since we don't retry previous
// attempts even if they failed. Value in the map indicates whether the log is finalized or not.
func (r *ExecutionReportingPlugin) getExecutedSeqNrsInRange(ctx context.Context, min, max uint64, latestBlock int64) (map[uint64]bool, error) {
	executedLogs, err := r.config.destLP.IndexedLogsTopicRange(
		abihelpers.EventSignatures.ExecutionStateChanged,
		r.config.offRamp.Address(),
		abihelpers.EventSignatures.ExecutionStateChangedSequenceNumberIndex,
		logpoller.EvmWord(min),
		logpoller.EvmWord(max),
		int(r.offchainConfig.DestOptimisticConfirmations),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}
	executedMp := make(map[uint64]bool)
	for _, executedLog := range executedLogs {
		exec, err := r.config.offRamp.ParseExecutionStateChanged(executedLog.GetGethLog())
		if err != nil {
			return nil, err
		}
		finalized := (latestBlock - executedLog.BlockNumber) >= int64(r.offchainConfig.DestFinalityDepth)
		executedMp[exec.SequenceNumber] = finalized
	}
	return executedMp, nil
}

// Builds a batch of transactions that can be executed, takes into account
// the available gas, rate limiting, execution state, nonce state, and
// profitability of execution.
func (r *ExecutionReportingPlugin) buildBatch(
	lggr logger.Logger,
	report commitReportWithSendRequests,
	inflight []InflightInternalExecutionReport,
	aggregateTokenLimit *big.Int,
	srcTokenPricesUSD map[common.Address]*big.Int,
	destTokenPricesUSD map[common.Address]*big.Int,
	execGasPriceEstimate LazyFunction[*big.Int],
	srcToDestToken map[common.Address]common.Address,
) (executableMessages []ObservedMessage) {
	inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, err := inflightAggregates(inflight, destTokenPricesUSD, srcToDestToken)
	if err != nil {
		lggr.Errorw("Unexpected error computing inflight values", "err", err)
		return []ObservedMessage{}
	}
	availableGas := uint64(r.offchainConfig.BatchGasLimit)
	aggregateTokenLimit.Sub(aggregateTokenLimit, inflightAggregateValue)
	expectedNonces := make(map[common.Address]uint64)
	for _, msg := range report.sendRequestsWithMeta {
		msgLggr := lggr.With("messageID", hexutil.Encode(msg.MessageId[:]))
		if msg.executed && msg.finalized {
			msgLggr.Infow("Skipping message already executed", "seqNr", msg.SequenceNumber)
			continue
		}
		if _, isInflight := inflightSeqNrs[msg.SequenceNumber]; isInflight {
			msgLggr.Infow("Skipping message already inflight", "seqNr", msg.SequenceNumber)
			continue
		}
		if _, ok := expectedNonces[msg.Sender]; !ok {
			// First message in batch, need to populate expected nonce
			if maxInflight, ok := maxInflightSenderNonces[msg.Sender]; ok {
				// Sender already has inflight nonce, populate from there
				expectedNonces[msg.Sender] = maxInflight + 1
			} else {
				// Nothing inflight take from chain.
				// Chain holds existing nonce.
				nonce, err := r.config.offRamp.GetSenderNonce(nil, msg.Sender)
				if err != nil {
					lggr.Errorw("unable to get sender nonce", "err", err)
					continue
				}
				expectedNonces[msg.Sender] = nonce + 1
			}
		}
		// Check expected nonce is valid
		if msg.Nonce != expectedNonces[msg.Sender] {
			msgLggr.Warnw("Skipping message invalid nonce", "have", msg.Nonce, "want", expectedNonces[msg.Sender])
			continue
		}
		msgValue, err := aggregateTokenValue(destTokenPricesUSD, srcToDestToken, msg.TokenAmounts)
		if err != nil {
			msgLggr.Errorw("Skipping message unable to compute aggregate value", "err", err)
			continue
		}
		// if token limit is smaller than message value skip message
		if aggregateTokenLimit.Cmp(msgValue) == -1 {
			msgLggr.Warnw("token limit is smaller than message value", "aggregateTokenLimit", aggregateTokenLimit.String(), "msgValue", msgValue.String())
			continue
		}
		// Fee boosting
		execGasPriceEstimateValue, err := execGasPriceEstimate()
		if err != nil {
			msgLggr.Errorw("Unexpected error fetching gas price estimate", "err", err)
			return []ObservedMessage{}
		}

		dstWrappedNativePrice, exists := destTokenPricesUSD[r.destWrappedNative]
		if !exists {
			msgLggr.Errorw("token not in dst token prices", "token", r.destWrappedNative)
			continue
		}

		execCostUsd := computeExecCost(msg.GasLimit, execGasPriceEstimateValue, dstWrappedNativePrice)
		// calculating the source chain fee, dividing by 1e18 for denomination.
		// For example:
		// FeeToken=link; FeeTokenAmount=1e17 i.e. 0.1 link, price is 6e18 USD/link (1 USD = 1e18),
		// availableFee is 1e17*6e18/1e18 = 6e17 = 0.6 USD

		srcFeeTokenPrice, exists := srcTokenPricesUSD[msg.FeeToken]
		if !exists {
			msgLggr.Errorw("token not in src token prices", "token", msg.FeeToken)
			continue
		}

		availableFee := big.NewInt(0).Mul(msg.FeeTokenAmount, srcFeeTokenPrice)
		availableFee = availableFee.Div(availableFee, big.NewInt(1e18))
		availableFeeUsd := waitBoostedFee(time.Since(msg.blockTimestamp), availableFee, r.offchainConfig.RelativeBoostPerWaitHour)
		if availableFeeUsd.Cmp(execCostUsd) < 0 {
			msgLggr.Infow("Insufficient remaining fee", "availableFeeUsd", availableFeeUsd, "execCostUsd", execCostUsd,
				"srcBlockTimestamp", msg.blockTimestamp, "waitTime", time.Since(msg.blockTimestamp), "boost", r.offchainConfig.RelativeBoostPerWaitHour)
			continue
		}

		messageMaxGas, err := calculateMessageMaxGas(
			msg.GasLimit,
			len(report.sendRequestsWithMeta),
			len(msg.Data),
			len(msg.TokenAmounts),
		)
		if err != nil {
			msgLggr.Errorw("calculate message max gas", "err", err)
			continue
		}

		// Check sufficient gas in batch
		if availableGas < messageMaxGas {
			msgLggr.Infow("Insufficient remaining gas in batch limit", "availableGas", availableGas, "messageMaxGas", messageMaxGas)
			continue
		}
		availableGas -= messageMaxGas
		aggregateTokenLimit.Sub(aggregateTokenLimit, msgValue)

		var tokenData [][]byte

		// TODO add attestation data for USDC here
		for range msg.TokenAmounts {
			tokenData = append(tokenData, []byte{})
		}

		msgLggr.Infow("Adding msg to batch", "seqNum", msg.SequenceNumber, "nonce", msg.Nonce,
			"value", msgValue, "aggregateTokenLimit", aggregateTokenLimit)
		executableMessages = append(executableMessages, NewObservedMessage(msg.SequenceNumber, tokenData))
		expectedNonces[msg.Sender] = msg.Nonce + 1
	}
	return executableMessages
}

func calculateMessageMaxGas(gasLimit *big.Int, numRequests, dataLen, numTokens int) (uint64, error) {
	if !gasLimit.IsUint64() {
		return 0, fmt.Errorf("gas limit %s cannot be casted to uint64", gasLimit)
	}

	gasLimitU64 := gasLimit.Uint64()
	gasOverHeadGas := maxGasOverHeadGas(numRequests, dataLen, numTokens)
	messageMaxGas := gasLimitU64 + gasOverHeadGas

	if messageMaxGas < gasLimitU64 || messageMaxGas < gasOverHeadGas {
		return 0, fmt.Errorf("message max gas overflow, gasLimit=%d gasOverHeadGas=%d", gasLimitU64, gasOverHeadGas)
	}

	return messageMaxGas, nil
}

// helper struct to hold the commitReport and the related send requests
type commitReportWithSendRequests struct {
	commitReport         commit_store.CommitStoreCommitReport
	sendRequestsWithMeta []evm2EVMOnRampCCIPSendRequestedWithMeta
}

func (r *commitReportWithSendRequests) validate() error {
	// make sure that number of messages is the expected
	if exp := int(r.commitReport.Interval.Max - r.commitReport.Interval.Min + 1); len(r.sendRequestsWithMeta) != exp {
		return errors.Errorf(
			"unexpected missing sendRequestsWithMeta in committed root %x have %d want %d", r.commitReport.MerkleRoot, len(r.sendRequestsWithMeta), exp)
	}

	return nil
}

func (r *commitReportWithSendRequests) allRequestsAreExecutedAndFinalized() bool {
	for _, req := range r.sendRequestsWithMeta {
		if !req.executed || !req.finalized {
			return false
		}
	}
	return true
}

// checks if the send request fits the commit report interval
func (r *commitReportWithSendRequests) sendReqFits(sendReq evm2EVMOnRampCCIPSendRequestedWithMeta) bool {
	return sendReq.SequenceNumber >= r.commitReport.Interval.Min &&
		sendReq.SequenceNumber <= r.commitReport.Interval.Max
}

// getReportsWithSendRequests returns the target reports with populated send requests.
func (r *ExecutionReportingPlugin) getReportsWithSendRequests(
	ctx context.Context,
	reports []commit_store.CommitStoreCommitReport,
) ([]commitReportWithSendRequests, error) {
	if len(reports) == 0 {
		return nil, nil
	}

	// find interval from all the reports
	intervalMin := reports[0].Interval.Min
	intervalMax := reports[0].Interval.Max
	for _, report := range reports[1:] {
		if report.Interval.Max > intervalMax {
			intervalMax = report.Interval.Max
		}
		if report.Interval.Min < intervalMin {
			intervalMin = report.Interval.Min
		}
	}

	// use errgroup to fetch send request logs and executed sequence numbers in parallel
	eg := &errgroup.Group{}

	var sendRequestLogs []logpoller.Log
	eg.Go(func() error {
		// get logs from all the reports
		rawLogs, err := r.config.sourceLP.LogsDataWordRange(
			abihelpers.EventSignatures.SendRequested,
			r.config.onRamp.Address(),
			abihelpers.EventSignatures.SendRequestedSequenceNumberWord,
			logpoller.EvmWord(intervalMin),
			logpoller.EvmWord(intervalMax),
			int(r.offchainConfig.SourceFinalityDepth),
			pg.WithParentCtx(ctx),
		)
		if err != nil {
			return err
		}
		sendRequestLogs = rawLogs
		return nil
	})

	var executedSeqNums map[uint64]bool
	eg.Go(func() error {
		latestBlock, err := r.config.destLP.LatestBlock()
		if err != nil {
			return err
		}
		// get executable sequence numbers
		executedMp, err := r.getExecutedSeqNrsInRange(ctx, intervalMin, intervalMax, latestBlock)
		if err != nil {
			return err
		}
		executedSeqNums = executedMp
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	reportsWithSendReqs := make([]commitReportWithSendRequests, len(reports))
	for i, report := range reports {
		reportsWithSendReqs[i] = commitReportWithSendRequests{
			commitReport:         report,
			sendRequestsWithMeta: make([]evm2EVMOnRampCCIPSendRequestedWithMeta, 0, len(sendRequestLogs)),
		}
	}

	for _, rawLog := range sendRequestLogs {
		ccipSendRequested, err := r.config.onRamp.ParseCCIPSendRequested(gethtypes.Log{
			Topics: rawLog.GetTopics(),
			Data:   rawLog.Data,
		})
		if err != nil {
			r.lggr.Errorw("unable to parse message", "err", err, "tx", rawLog.TxHash, "logIdx", rawLog.LogIndex)
			continue
		}
		msg := abihelpers.OnRampMessageToOffRampMessage(ccipSendRequested.Message)

		// if value exists in the map then it's executed
		// if value exists, and it's true then it's considered finalized
		finalized, executed := executedSeqNums[msg.SequenceNumber]

		sendReq := evm2EVMOnRampCCIPSendRequestedWithMeta{
			InternalEVM2EVMMessage: msg,
			blockTimestamp:         rawLog.BlockTimestamp,
			executed:               executed,
			finalized:              finalized,
		}

		// attach the msg to the appropriate reports
		for i := range reportsWithSendReqs {
			if reportsWithSendReqs[i].sendReqFits(sendReq) {
				reportsWithSendReqs[i].sendRequestsWithMeta = append(reportsWithSendReqs[i].sendRequestsWithMeta, sendReq)
			}
		}
	}

	return reportsWithSendReqs, nil
}

func aggregateTokenValue(destTokenPricesUSD map[common.Address]*big.Int, srcToDst map[common.Address]common.Address, tokensAndAmount []evm_2_evm_offramp.ClientEVMTokenAmount) (*big.Int, error) {
	sum := big.NewInt(0)
	for i := 0; i < len(tokensAndAmount); i++ {
		price, ok := destTokenPricesUSD[srcToDst[tokensAndAmount[i].Token]]
		if !ok {
			return nil, errors.Errorf("do not have price for src token %v", tokensAndAmount[i].Token)
		}
		sum.Add(sum, new(big.Int).Quo(new(big.Int).Mul(price, tokensAndAmount[i].Amount), big.NewInt(1e18)))
	}
	return sum, nil
}

func (r *ExecutionReportingPlugin) parseSeqNr(log logpoller.Log) (uint64, error) {
	s, err := r.config.onRamp.ParseCCIPSendRequested(log.ToGethLog())
	if err != nil {
		return 0, err
	}
	return s.Message.SequenceNumber, nil
}

// Assumes non-empty report. Messages to execute can span more than one report, but are assumed to be in order of increasing
// sequence number.
func (r *ExecutionReportingPlugin) buildReport(ctx context.Context, lggr logger.Logger, observedMessages []ObservedMessage) ([]byte, error) {
	if err := validateSeqNumbers(ctx, r.config.commitStore, observedMessages); err != nil {
		return nil, err
	}
	commitReport, err := getCommitReportForSeqNum(ctx, r.config.destLP, r.config.commitStore, observedMessages[0].SeqNr)
	if err != nil {
		return nil, err
	}
	lggr.Infow("Building execution report", "observations", observedMessages, "merkleRoot", hexutil.Encode(commitReport.MerkleRoot[:]), "report", commitReport)

	msgsInRoot, leaves, tree, err := getProofData(ctx, lggr, r.config.leafHasher, r.parseSeqNr, r.config.onRamp.Address(), r.config.sourceLP, commitReport.Interval)
	if err != nil {
		return nil, err
	}

	messages := make([]*evm_2_evm_offramp.InternalEVM2EVMMessage, len(msgsInRoot))
	for i, msg := range msgsInRoot {
		decodedMessage, err2 := abihelpers.DecodeOffRampMessage(msg.Data)
		if err2 != nil {
			return nil, err
		}
		messages[i] = decodedMessage
	}

	// cap messages which fits MaxExecutionReportLength (after serialized)
	capped := sort.Search(len(observedMessages), func(i int) bool {
		report, _ := buildExecutionReportForMessages(messages, leaves, tree, commitReport.Interval, observedMessages[:i+1])
		var encoded []byte
		encoded, err = abihelpers.EncodeExecutionReport(report)
		if err != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > MaxObservationLength
	})
	if err != nil {
		return nil, err
	}

	execReport, hashes := buildExecutionReportForMessages(messages, leaves, tree, commitReport.Interval, observedMessages[:capped])
	encodedReport, err := abihelpers.EncodeExecutionReport(execReport)
	if err != nil {
		return nil, err
	}

	if capped < len(observedMessages) {
		lggr.Warnf(
			"Capping report to fit MaxExecutionReportLength: msgsCount %d -> %d, bytes %d, bytesLimit %d",
			len(observedMessages), capped, len(encodedReport), MaxExecutionReportLength,
		)
	}

	// Double check this verifies before sending.
	res, err := r.config.commitStore.Verify(&bind.CallOpts{Context: ctx}, hashes, execReport.Proofs, execReport.ProofFlagBits)
	if err != nil {
		lggr.Errorw("Unable to call verify", "observations", observedMessages[:capped], "root", commitReport.MerkleRoot[:], "seqRange", commitReport.Interval, "err", err)
		return nil, err
	}
	// No timestamp, means failed to verify root.
	if res.Cmp(big.NewInt(0)) == 0 {
		root := tree.Root()
		lggr.Errorf("Root does not verify for messages: %v, our inner root 0x%x", observedMessages[:capped], root)
		return nil, errors.New("root does not verify")
	}
	return encodedReport, nil
}

func (r *ExecutionReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	lggr := r.lggr.Named("ExecutionReport")
	parsableObservations := getParsableObservations[ExecutionObservation](lggr, observations)
	// Need at least F+1 observations
	if len(parsableObservations) <= r.F {
		lggr.Warn("Non-empty observations <= F, need at least F+1 to continue")
		return false, nil, nil
	}

	observedMessages, err := calculateObservedMessagesConsensus(parsableObservations, r.F)
	if err != nil {
		return false, nil, err
	}
	if len(observedMessages) == 0 {
		return false, nil, nil
	}

	report, err := r.buildReport(ctx, lggr, observedMessages)
	if err != nil {
		return false, nil, err
	}
	lggr.Infow("Report", "executableObservations", observedMessages)
	return true, report, nil
}

type tallyKey struct {
	seqNr         uint64
	tokenDataHash [32]byte
}

type tallyVal struct {
	tally     int
	tokenData [][]byte
}

func calculateObservedMessagesConsensus(observations []ExecutionObservation, f int) ([]ObservedMessage, error) {
	tally := make(map[tallyKey]tallyVal)
	for _, obs := range observations {
		for seqNr, msgData := range obs.Messages {
			tokenDataHash, err := bytesOfBytesKeccak(msgData.TokenData)
			if err != nil {
				return nil, fmt.Errorf("bytes of bytes keccak: %w", err)
			}

			key := tallyKey{seqNr: seqNr, tokenDataHash: tokenDataHash}
			if val, ok := tally[key]; ok {
				tally[key] = tallyVal{tally: val.tally + 1, tokenData: msgData.TokenData}
			} else {
				tally[key] = tallyVal{tally: 1, tokenData: msgData.TokenData}
			}
		}
	}

	// We might have different token data for the same sequence number.
	// For that purpose we want to keep the token data with the most occurrences.
	seqNumTally := make(map[uint64]tallyVal)
	for key, tallyInfo := range tally {
		existingTally, exists := seqNumTally[key.seqNr]
		if tallyInfo.tally > f && (!exists || tallyInfo.tally > existingTally.tally) {
			seqNumTally[key.seqNr] = tallyInfo
		}
	}

	finalSequenceNumbers := make([]ObservedMessage, 0, len(seqNumTally))
	for seqNr, tallyInfo := range seqNumTally {
		finalSequenceNumbers = append(finalSequenceNumbers, NewObservedMessage(seqNr, tallyInfo.tokenData))
	}
	// buildReport expects sorted sequence numbers (tally map is non-deterministic).
	sort.Slice(finalSequenceNumbers, func(i, j int) bool {
		return finalSequenceNumbers[i].SeqNr < finalSequenceNumbers[j].SeqNr
	})
	return finalSequenceNumbers, nil
}

func (r *ExecutionReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldAcceptFinalizedReport")
	messages, err := abihelpers.MessagesFromExecutionReport(report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, err
	}
	// If the first message is executed already, this execution report is stale, and we do not accept it.
	stale, err := r.isStaleReport(messages)
	if err != nil {
		return false, err
	}
	if stale {
		lggr.Infow("Execution report is stale", "messages", messages)
		return false, nil
	}
	// Else just assume in flight
	if err = r.inflightReports.add(lggr, messages); err != nil {
		return false, err
	}
	lggr.Info("Accepting finalized report")
	return true, nil
}

func (r *ExecutionReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	messages, err := abihelpers.MessagesFromExecutionReport(report)
	if err != nil {
		return false, nil
	}
	// If report is not stale we transmit.
	// When the executeTransmitter enqueues the tx for tx manager,
	// we mark it as execution_sent, removing it from the set of inflight messages.
	stale, err := r.isStaleReport(messages)
	return !stale, err
}

func (r *ExecutionReportingPlugin) isStaleReport(messages []evm_2_evm_offramp.InternalEVM2EVMMessage) (bool, error) {
	// If the first message is executed already, this execution report is stale.
	// Note the default execution state, including for arbitrary seq number not yet committed
	// is ExecutionStateUntouched.
	msgState, err := r.config.offRamp.GetExecutionState(nil, messages[0].SequenceNumber)
	if err != nil {
		return true, err
	}
	if state := abihelpers.MessageExecutionState(msgState); state == abihelpers.ExecutionStateFailure || state == abihelpers.ExecutionStateSuccess {
		return true, nil
	}

	return false, nil
}

func (r *ExecutionReportingPlugin) Close() error {
	return nil
}

func inflightAggregates(
	inflight []InflightInternalExecutionReport,
	destTokenPrices map[common.Address]*big.Int,
	srcToDst map[common.Address]common.Address,
) (map[uint64]struct{}, *big.Int, map[common.Address]uint64, error) {
	inflightSeqNrs := make(map[uint64]struct{})
	inflightAggregateValue := big.NewInt(0)
	maxInflightSenderNonces := make(map[common.Address]uint64)
	for _, rep := range inflight {
		for _, message := range rep.messages {
			inflightSeqNrs[message.SequenceNumber] = struct{}{}
			msgValue, err := aggregateTokenValue(destTokenPrices, srcToDst, message.TokenAmounts)
			if err != nil {
				return nil, nil, nil, err
			}
			inflightAggregateValue.Add(inflightAggregateValue, msgValue)
			maxInflightSenderNonce, ok := maxInflightSenderNonces[message.Sender]
			if !ok || message.Nonce > maxInflightSenderNonce {
				maxInflightSenderNonces[message.Sender] = message.Nonce
			}
		}
	}
	return inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, nil
}

// getTokensPrices returns token prices of the given price registry,
// results include feeTokens and passed-in tokens
// price values are USD per 1e18 of smallest token denomination, in base units 1e18 (e.g. 5$ = 5e18 USD per 1e18 units).
// this function is used for price registry of both source and destination chains.
func getTokensPrices(ctx context.Context, feeTokens []common.Address, priceRegistry price_registry.PriceRegistryInterface, tokens []common.Address) (map[common.Address]*big.Int, error) {
	prices := make(map[common.Address]*big.Int)

	wantedTokens := append(feeTokens, tokens...)
	wantedPrices, err := priceRegistry.GetTokenPrices(&bind.CallOpts{Context: ctx}, wantedTokens)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get token prices of %v", wantedTokens)
	}
	for i, token := range wantedTokens {
		prices[token] = wantedPrices[i].Value
	}

	return prices, nil
}

func getUnexpiredCommitReports(
	ctx context.Context,
	dstLogPoller logpoller.LogPoller,
	commitStore commit_store.CommitStoreInterface,
	permissionExecutionThreshold time.Duration,
) ([]commit_store.CommitStoreCommitReport, error) {
	logs, err := dstLogPoller.LogsCreatedAfter(
		abihelpers.EventSignatures.ReportAccepted,
		commitStore.Address(),
		time.Now().Add(-permissionExecutionThreshold),
		0,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}
	var reports []commit_store.CommitStoreCommitReport
	for _, log := range logs {
		reportAccepted, err := commitStore.ParseReportAccepted(gethtypes.Log{
			Topics: log.GetTopics(),
			Data:   log.Data,
		})
		if err != nil {
			return nil, err
		}
		reports = append(reports, reportAccepted.Report)
	}
	return reports, nil
}
