package ccip

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ccip/hasher"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ccip/merklemulti"
)

const MaxCommitReportLength = 1000

var (
	_ types.ReportingPluginFactory = &CommitReportingPluginFactory{}
	_ types.ReportingPlugin        = &CommitReportingPlugin{}
)

// EncodeCommitReport abi encodes an offramp.InternalCommitReport.
func EncodeCommitReport(commitReport *commit_store.ICommitStoreCommitReport) (types.Report, error) {
	report, err := makeCommitReportArgs().PackValues([]interface{}{
		commitReport,
	})
	if err != nil {
		return nil, err
	}
	return report, nil
}

// DecodeCommitReport abi decodes a types.Report to an ICommitStoreCommitReport
func DecodeCommitReport(report types.Report) (*commit_store.ICommitStoreCommitReport, error) {
	unpacked, err := makeCommitReportArgs().Unpack(report)
	if err != nil {
		return nil, err
	}
	if len(unpacked) != 1 {
		return nil, errors.New("expected single struct value")
	}

	commitReport, ok := unpacked[0].(struct {
		PriceUpdates struct {
			FeeTokenPriceUpdates []struct {
				SourceFeeToken common.Address `json:"sourceFeeToken"`
				UsdPerFeeToken *big.Int       `json:"usdPerFeeToken"`
			} `json:"feeTokenPriceUpdates"`
			DestChainId   uint64   `json:"destChainId"`
			UsdPerUnitGas *big.Int `json:"usdPerUnitGas"`
		} `json:"priceUpdates"`
		Interval struct {
			Min uint64 `json:"min"`
			Max uint64 `json:"max"`
		} `json:"interval"`
		MerkleRoot [32]byte `json:"merkleRoot"`
	})
	if !ok {
		return nil, errors.Errorf("invalid commit report got %T", unpacked[0])
	}

	var feeTokenUpdates []commit_store.InternalFeeTokenPriceUpdate
	for _, u := range commitReport.PriceUpdates.FeeTokenPriceUpdates {
		feeTokenUpdates = append(feeTokenUpdates, commit_store.InternalFeeTokenPriceUpdate{
			SourceFeeToken: u.SourceFeeToken,
			UsdPerFeeToken: u.UsdPerFeeToken,
		})
	}

	return &commit_store.ICommitStoreCommitReport{
		PriceUpdates: commit_store.InternalPriceUpdates{
			DestChainId:          commitReport.PriceUpdates.DestChainId,
			UsdPerUnitGas:        commitReport.PriceUpdates.UsdPerUnitGas,
			FeeTokenPriceUpdates: feeTokenUpdates,
		},
		Interval: commit_store.ICommitStoreInterval{
			Min: commitReport.Interval.Min,
			Max: commitReport.Interval.Max,
		},
		MerkleRoot: commitReport.MerkleRoot,
	}, nil
}

func isCommitStoreDownNow(lggr logger.Logger, commitStore *commit_store.CommitStore) bool {
	paused, err := commitStore.Paused(nil)
	if err != nil {
		// Air on side of caution by halting if we cannot read the state?
		lggr.Errorw("Unable to read offramp paused", "err", err)
		return true
	}
	healthy, err := commitStore.IsAFNHealthy(nil)
	if err != nil {
		lggr.Errorw("Unable to read offramp afn", "err", err)
		return true
	}
	return paused || !healthy
}

type InflightReport struct {
	report    *commit_store.ICommitStoreCommitReport
	createdAt time.Time
}

type InflightFeeUpdate struct {
	priceUpdates commit_store.InternalPriceUpdates
	createdAt    time.Time
}

type CommitPluginConfig struct {
	lggr                                 logger.Logger
	source, dest                         logpoller.LogPoller
	seqParsers                           func(log logpoller.Log) (uint64, error)
	reqEventSig                          EventSignatures
	onRamp                               common.Address
	offRamp                              *evm_2_evm_offramp.EVM2EVMOffRamp
	priceRegistry                        *price_registry.PriceRegistry
	priceGetter                          PriceGetter
	sourceNative                         common.Address
	sourceGasEstimator, destGasEstimator gas.Estimator
	sourceChainID                        uint64
	commitStore                          *commit_store.CommitStore
	hasher                               LeafHasherInterface[[32]byte]
	inflightCacheExpiry                  time.Duration
}

type CommitReportingPluginFactory struct {
	config CommitPluginConfig
}

// NewCommitReportingPluginFactory return a new CommitReportingPluginFactory.
func NewCommitReportingPluginFactory(config CommitPluginConfig) types.ReportingPluginFactory {
	return &CommitReportingPluginFactory{config: config}
}

// NewReportingPlugin returns the ccip CommitReportingPlugin and satisfies the ReportingPluginFactory interface.
func (rf *CommitReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	offchainConfig, err := Decode(config.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	execTokens, err := rf.config.offRamp.GetDestinationTokens(nil)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	// TODO: Hack assume link is the first token  https://smartcontract-it.atlassian.net/browse/CCIP-304
	// Only set link token as fee token for now
	linkToken := execTokens[0]

	return &CommitReportingPlugin{
			config:         rf.config,
			F:              config.F,
			feeTokens:      []common.Address{linkToken},
			inFlight:       make(map[[32]byte]InflightReport),
			offchainConfig: offchainConfig,
		},
		types.ReportingPluginInfo{
			Name:          "CCIPCommit",
			UniqueReports: true,
			Limits: types.ReportingPluginLimits{
				MaxQueryLength:       MaxQueryLength,
				MaxObservationLength: MaxObservationLength,
				MaxReportLength:      MaxCommitReportLength,
			},
		}, nil
}

type CommitReportingPlugin struct {
	config    CommitPluginConfig
	F         int
	feeTokens []common.Address
	// We need to synchronize access to the inflight structure
	// as reporting plugin methods may be called from separate goroutines,
	// e.g. reporting vs transmission protocol.
	inFlightMu         sync.RWMutex
	inFlight           map[[32]byte]InflightReport
	inFlightFeeUpdates []InflightFeeUpdate
	offchainConfig     OffchainConfig
}

func (r *CommitReportingPlugin) nextMinSeqNumForInFlight() uint64 {
	r.inFlightMu.RLock()
	defer r.inFlightMu.RUnlock()
	max := uint64(0)
	for _, report := range r.inFlight {
		if report.report.Interval.Max > max {
			max = report.report.Interval.Max
		}
	}
	return max + 1
}

func (r *CommitReportingPlugin) nextMinSeqNum() (uint64, error) {
	nextMin, err := r.config.commitStore.GetExpectedNextSequenceNumber(nil)
	if err != nil {
		return 0, err
	}
	nextMinInFlight := r.nextMinSeqNumForInFlight()
	if nextMinInFlight > nextMin {
		nextMin = nextMinInFlight
	}
	return nextMin, nil
}

func (r *CommitReportingPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

func calculateUsdPerUnitGas(sourceGasPrice *big.Int, usdPerFeeCoin *big.Int) *big.Int {
	// (wei / gas) * (usd / eth) * (1 eth / 1e18 wei)  = usd/gas
	tmp := big.NewInt(0).Mul(sourceGasPrice, usdPerFeeCoin)
	return tmp.Div(tmp, big.NewInt(1e18))
}

// deviation_parts_per_billion = ((x2 - x1) / x1) * 1e9
func (r *CommitReportingPlugin) deviates(x1, x2 *big.Int) bool {
	gasPriceDeviation := big.NewInt(0).Sub(x1, x2)
	gasPriceDeviation.Mul(gasPriceDeviation, big.NewInt(1e9))
	gasPriceDeviation.Div(gasPriceDeviation, x1)
	return gasPriceDeviation.CmpAbs(big.NewInt(int64(r.offchainConfig.FeeUpdateDeviationPPB))) > 0
}

// All prices are USD denominated.
func (r *CommitReportingPlugin) canSkipFeeUpdate(gasPrice *big.Int, tokenPrices map[common.Address]*big.Int) (bool, error) {
	var latestGasUpdateTimestamp time.Time
	var latestGasPrice *big.Int
	gasUpdatesWithinHeartBeat, err := r.config.dest.IndexedLogsCreatedAfter(GasFeeUpdated, r.config.priceRegistry.Address(), 1, []common.Hash{EvmWord(r.config.sourceChainID)}, time.Now().Add(-r.offchainConfig.FeeUpdateHeartBeat.Duration()))
	if err != nil {
		return false, err
	}
	if len(gasUpdatesWithinHeartBeat) > 0 {
		// Ordered by ascending timestamps
		priceUpdate, err := r.config.priceRegistry.ParseUsdPerUnitGasUpdated(gasUpdatesWithinHeartBeat[len(gasUpdatesWithinHeartBeat)-1].GetGethLog())
		if err != nil {
			return false, err
		}
		latestGasUpdateTimestamp = time.Unix(priceUpdate.Timestamp.Int64(), 0)
		latestGasPrice = priceUpdate.Value
	}
	r.inFlightMu.RLock()
	for _, inflight := range r.inFlightFeeUpdates {
		if !inflight.createdAt.Before(latestGasUpdateTimestamp) {
			latestGasUpdateTimestamp = inflight.createdAt
			latestGasPrice = inflight.priceUpdates.UsdPerUnitGas
		}
	}
	r.inFlightMu.RUnlock()

	if latestGasPrice == nil || time.Since(latestGasUpdateTimestamp) > r.offchainConfig.FeeUpdateHeartBeat.Duration() || r.deviates(gasPrice, latestGasPrice) {
		return false, nil
	}

	for _, feeToken := range r.feeTokens {
		var latestFeeTokenPriceTimestamp time.Time
		var latestFeeTokenPrice *big.Int
		feeTokenUpdatesWithinHeartBeat, err := r.config.dest.IndexedLogsCreatedAfter(GasFeeUpdated, r.config.priceRegistry.Address(), 1, []common.Hash{feeToken.Hash()}, time.Now().Add(-r.offchainConfig.FeeUpdateHeartBeat.Duration()))
		if err != nil {
			return false, err
		}
		if len(feeTokenUpdatesWithinHeartBeat) > 0 {
			// Ordered by ascending timestamps
			parsed, err := r.config.priceRegistry.ParseUsdPerFeeTokenUpdated(feeTokenUpdatesWithinHeartBeat[len(gasUpdatesWithinHeartBeat)-1].GetGethLog())
			if err != nil {
				return false, err
			}
			latestFeeTokenPriceTimestamp = time.Unix(parsed.Timestamp.Int64(), 0)
			latestFeeTokenPrice = parsed.Value
		}
		r.inFlightMu.RLock()
		for _, inflight := range r.inFlightFeeUpdates {
			for _, update := range inflight.priceUpdates.FeeTokenPriceUpdates {
				if update.SourceFeeToken == feeToken && !inflight.createdAt.Before(latestGasUpdateTimestamp) {
					latestGasUpdateTimestamp = inflight.createdAt
					latestFeeTokenPrice = update.UsdPerFeeToken
				}
			}
		}
		r.inFlightMu.RUnlock()

		if latestFeeTokenPrice == nil || time.Since(latestFeeTokenPriceTimestamp) > r.offchainConfig.FeeUpdateHeartBeat.Duration() || r.deviates(tokenPrices[feeToken], latestFeeTokenPrice) {
			return false, nil
		}
	}
	// If we make it all the way here, nothing needs updating, and we can skip.
	return true, nil
}

func (r *CommitReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	lggr := r.config.lggr.Named("CommitObservation")
	if isCommitStoreDownNow(lggr, r.config.commitStore) {
		return nil, ErrCommitStoreIsDown
	}
	r.expireInflight(lggr)

	// Will return 0,0 if no messages are found. This is a valid case as the report could
	// still contain fee updates.
	min, max, err := r.calculateMinMaxSequenceNumbers(lggr)
	if err != nil {
		return nil, err
	}
	// Include wrapped native in our token query as way to identify the source native USD price.
	tokenPricesUSD, err := r.config.priceGetter.TokenPricesUSD(context.Background(), append(r.feeTokens, r.config.sourceNative))
	if err != nil {
		return nil, err
	}
	// Observe a source chain price for pricing.
	// TODO: 1559 support https://smartcontract-it.atlassian.net/browse/CCIP-316
	sourceGasPriceWei, _, err := r.config.sourceGasEstimator.GetLegacyGas(ctx, nil, BatchGasLimit, assets.NewWei(big.NewInt(MaxGasPrice)))
	if err != nil {
		return nil, err
	}

	sourceGasPriceUSD := calculateUsdPerUnitGas(sourceGasPriceWei.ToInt(), tokenPricesUSD[r.config.sourceNative])
	if canSkip, err := r.canSkipFeeUpdate(sourceGasPriceUSD, tokenPricesUSD); err != nil {
		return nil, err
	} else if canSkip {
		sourceGasPriceUSD = nil // vote skip
	}

	return CommitObservation{
		Interval: commit_store.ICommitStoreInterval{
			Min: min,
			Max: max,
		},
		TokenPricesUSD:    tokenPricesUSD,
		SourceGasPriceUSD: sourceGasPriceUSD,
	}.Marshal()
}

func (r *CommitReportingPlugin) calculateMinMaxSequenceNumbers(lggr logger.Logger) (uint64, uint64, error) {
	nextMin, err := r.nextMinSeqNum()
	if err != nil {
		return 0, 0, err
	}
	// All available messages that have not been committed yet and have sufficient confirmations.
	lggr.Infof("Looking for requests with sig %s and nextMin %d on onRampAddr %s", r.config.reqEventSig.SendRequested.Hex(), nextMin, r.config.onRamp.Hex())
	reqs, err := r.config.source.LogsDataWordGreaterThan(r.config.reqEventSig.SendRequested, r.config.onRamp, r.config.reqEventSig.SendRequestedSequenceNumberIndex, EvmWord(nextMin), int(r.offchainConfig.SourceIncomingConfirmations))
	if err != nil {
		return 0, 0, err
	}
	lggr.Infof("%d requests found for onRampAddr %s", len(reqs), r.config.onRamp.Hex())
	if len(reqs) == 0 {
		return 0, 0, nil
	}
	var seqNrs []uint64
	for _, req := range reqs {
		seqNr, err2 := r.config.seqParsers(req)
		if err2 != nil {
			lggr.Errorw("error parsing seq num", "err", err2)
			continue
		}
		seqNrs = append(seqNrs, seqNr)
	}
	min := seqNrs[0]
	max := seqNrs[len(seqNrs)-1]
	if min != nextMin {
		// Still report the observation as even partial reports have value e.g. all nodes are
		// missing a single, different log each, they would still be able to produce a valid report.
		lggr.Warnf("Missing sequence number range [%d-%d] for onRamp %s", nextMin, min, r.config.onRamp.Hex())
	}
	if !contiguousReqs(lggr, min, max, seqNrs) {
		return 0, 0, errors.New("unexpected gap in seq nums")
	}
	lggr.Infof("OnRamp %v: min %v max %v", r.config.onRamp, min, max)
	return min, max, nil
}

// buildReport assumes there is at least one message in reqs.
func (r *CommitReportingPlugin) buildReport(interval commit_store.ICommitStoreInterval, priceUpdates commit_store.InternalPriceUpdates) (*commit_store.ICommitStoreCommitReport, error) {
	lggr := r.config.lggr.Named("BuildReport")

	// If no messages are needed only include fee updates
	if interval.Min == 0 {
		return &commit_store.ICommitStoreCommitReport{
			PriceUpdates: priceUpdates,
			MerkleRoot:   [32]byte{},
			Interval:     interval,
		}, nil
	}

	leaves, err := leavesFromIntervals(lggr, r.config.onRamp, r.config.reqEventSig, r.config.seqParsers, interval, r.config.source, r.config.hasher, int(r.offchainConfig.SourceIncomingConfirmations))
	if err != nil {
		return nil, err
	}

	if len(leaves) == 0 {
		return nil, fmt.Errorf("tried building a tree without leaves for onRampAddr %s. %+v", r.config.onRamp.Hex(), leaves)
	}
	tree, err := merklemulti.NewTree(hasher.NewKeccakCtx(), leaves)
	if err != nil {
		return nil, err
	}

	return &commit_store.ICommitStoreCommitReport{
		PriceUpdates: priceUpdates,
		MerkleRoot:   tree.Root(),
		Interval:     interval,
	}, nil
}

func (r *CommitReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	lggr := r.config.lggr.Named("Report")
	if isCommitStoreDownNow(lggr, r.config.commitStore) {
		return false, nil, ErrCommitStoreIsDown
	}
	nonEmptyObservations := getNonEmptyObservations[CommitObservation](lggr, observations)
	var intervals []commit_store.ICommitStoreInterval
	for _, obs := range nonEmptyObservations {
		intervals = append(intervals, obs.Interval)
	}
	if len(intervals) <= r.F {
		lggr.Debugf("Observations for OnRamp %s 1 < #obs <= F, need at least F+1 to continue", r.config.onRamp.Hex())
		return false, nil, nil
	}

	agreedInterval, err := calculateIntervalConsensus(intervals, r.F, r.nextMinSeqNum)
	if err != nil {
		return false, nil, err
	}

	feeUpdates := r.calculateFeeUpdates(nonEmptyObservations)
	// If there are no fee updates and the interval is zero there is no report to produce.
	if len(feeUpdates.FeeTokenPriceUpdates) == 0 && feeUpdates.DestChainId == 0 && agreedInterval.Min == 0 {
		return false, nil, nil
	}

	report, err := r.buildReport(agreedInterval, feeUpdates)
	if err != nil {
		return false, nil, err
	}
	encodedReport, err := EncodeCommitReport(report)
	if err != nil {
		return false, nil, err
	}
	lggr.Infow("Built report", "interval", agreedInterval)
	return true, encodedReport, nil
}

func (r *CommitReportingPlugin) calculateFeeUpdates(observations []CommitObservation) commit_store.InternalPriceUpdates {
	priceObservations := make(map[common.Address][]*big.Int)
	var sourceGasObservations []*big.Int
	var sourceGasPriceNilCount int

	for _, obs := range observations {
		hasAllPrices := true
		for _, token := range r.feeTokens {
			if _, ok := obs.TokenPricesUSD[token]; !ok {
				hasAllPrices = false
				break
			}
		}
		// Disallow partial updates
		if !hasAllPrices {
			continue
		}
		if obs.SourceGasPriceUSD == nil {
			sourceGasPriceNilCount++
		} else {
			// Add only non-nil source gas price
			sourceGasObservations = append(sourceGasObservations, obs.SourceGasPriceUSD)
		}
		// If it has all the prices, add each price to observations
		for token, price := range obs.TokenPricesUSD {
			priceObservations[token] = append(priceObservations[token], price)
		}
	}
	var feeUpdates []commit_store.InternalFeeTokenPriceUpdate
	for _, feeToken := range r.feeTokens {
		medianPrice := median(priceObservations[feeToken])
		feeUpdates = append(feeUpdates, commit_store.InternalFeeTokenPriceUpdate{
			SourceFeeToken: feeToken,
			UsdPerFeeToken: medianPrice,
		})
	}

	// If majority report a gas price, include it in the update
	if sourceGasPriceNilCount < len(sourceGasObservations) {
		return commit_store.InternalPriceUpdates{
			FeeTokenPriceUpdates: feeUpdates,
			DestChainId:          r.config.sourceChainID,
			UsdPerUnitGas:        median(sourceGasObservations),
		}
	}
	return commit_store.InternalPriceUpdates{
		FeeTokenPriceUpdates: feeUpdates,
		DestChainId:          0,
		UsdPerUnitGas:        big.NewInt(0),
	}
}

// Assumed at least f+1 valid observations
func calculateIntervalConsensus(intervals []commit_store.ICommitStoreInterval, f int, nextMinSeqNumForOffRamp func() (uint64, error)) (commit_store.ICommitStoreInterval, error) {
	if len(intervals) <= f {
		return commit_store.ICommitStoreInterval{}, errors.Errorf("Not enough intervals to form consensus intervals %d, f %d", len(intervals), f)
	}
	// Extract the min and max
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Min < intervals[j].Min
	})
	minSeqNum := intervals[f].Min

	// The only way a report could have a minSeqNum of 0 is when there are no messages to report
	// and the report is potentially still valid for gas fee updates.
	if minSeqNum == 0 {
		return commit_store.ICommitStoreInterval{Min: 0, Max: 0}, nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Max < intervals[j].Max
	})
	// We use a conservative maximum. If we pick a value that some honest oracles might not
	// have seen they’ll end up not agreeing on a msg, stalling the protocol.
	maxSeqNum := intervals[f].Max
	// TODO: Do we for sure want to fail everything here?
	if maxSeqNum < minSeqNum {
		return commit_store.ICommitStoreInterval{}, errors.New("max seq num smaller than min")
	}
	nextMin, err := nextMinSeqNumForOffRamp()
	if err != nil {
		return commit_store.ICommitStoreInterval{}, err
	}
	// Contract would revert
	if nextMin > minSeqNum {
		return commit_store.ICommitStoreInterval{}, errors.Errorf("invalid min seq number got %v want %v", minSeqNum, nextMin)
	}

	return commit_store.ICommitStoreInterval{
		Min: minSeqNum,
		Max: maxSeqNum,
	}, nil
}

func (r *CommitReportingPlugin) expireInflight(lggr logger.Logger) {
	r.inFlightMu.Lock()
	defer r.inFlightMu.Unlock()
	// Reap any expired entries from inflight.
	for root, inFlightReport := range r.inFlight {
		if time.Since(inFlightReport.createdAt) > r.config.inflightCacheExpiry {
			// Happy path: inflight report was successfully transmitted onchain, we remove it from inflight and onchain state reflects inflight.
			// Sad path: inflight report reverts onchain, we remove it from inflight, onchain state does not reflect the chains so we retry.
			lggr.Infow("Inflight report expired", "rootOfRoots", hexutil.Encode(inFlightReport.report.MerkleRoot[:]))
			delete(r.inFlight, root)
		}
	}
	var stillInflight []InflightFeeUpdate
	for _, inFlightFeeUpdate := range r.inFlightFeeUpdates {
		if time.Since(inFlightFeeUpdate.createdAt) > r.config.inflightCacheExpiry {
			// Happy path: inflight report was successfully transmitted onchain, we remove it from inflight and onchain state reflects inflight.
			// Sad path: inflight report reverts onchain, we remove it from inflight, onchain state does not reflect the chains so we retry.
			lggr.Infow("Inflight price update expired", "updates", inFlightFeeUpdate.priceUpdates)
			stillInflight = append(stillInflight, inFlightFeeUpdate)
		}
	}
	r.inFlightFeeUpdates = stillInflight
}

func (r *CommitReportingPlugin) addToInflight(lggr logger.Logger, report *commit_store.ICommitStoreCommitReport) {
	r.inFlightMu.Lock()
	defer r.inFlightMu.Unlock()

	if report.MerkleRoot != [32]byte{} {
		// Set new inflight ones as pending
		lggr.Infow("Adding to inflight report", "rootOfRoots", hexutil.Encode(report.MerkleRoot[:]))
		r.inFlight[report.MerkleRoot] = InflightReport{
			report:    report,
			createdAt: time.Now(),
		}
	}

	if report.PriceUpdates.DestChainId != 0 || len(report.PriceUpdates.FeeTokenPriceUpdates) != 0 {
		lggr.Infow("Adding to inflight fee updates", "priceUpdates", report.PriceUpdates)
		r.inFlightFeeUpdates = append(r.inFlightFeeUpdates, InflightFeeUpdate{
			priceUpdates: report.PriceUpdates,
			createdAt:    time.Now(),
		})
	}
}

func (r *CommitReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.config.lggr.Named("ShouldAcceptFinalizedReport")
	parsedReport, err := DecodeCommitReport(report)
	if err != nil {
		return false, err
	}
	if parsedReport.MerkleRoot == [32]byte{} && parsedReport.PriceUpdates.DestChainId == 0 && len(parsedReport.PriceUpdates.FeeTokenPriceUpdates) == 0 {
		// Empty report, should not be put on chain
		return false, nil
	}

	if parsedReport.MerkleRoot != [32]byte{} {
		// Note it's ok to leave the unstarted requests behind, since the
		// 'Observe' is always based on the last reports onchain min seq num.
		if r.isStaleReport(parsedReport) {
			return false, nil
		}

		nextInflightMin, err := r.nextMinSeqNum()
		if err != nil {
			return false, err
		}
		if nextInflightMin != parsedReport.Interval.Min {
			// There are sequence numbers missing between the commitStore/inflight txs and the proposed report.
			// The report will fail onchain unless the inflight cache is in an incorrect state. A state like this
			// could happen for various reasons, e.g. a reboot of the node emptying the caches, and should be self-healing.
			// We do not submit a tx and wait for the protocol to self-heal by updating the caches or invalidating
			// inflight caches over time.
			lggr.Errorw("Next inflight min is not equal to the proposed min of the report", "nextInflightMin", nextInflightMin, "proposed min", parsedReport.Interval.Min)
			return false, errors.New("Next inflight min is not equal to the proposed min of the report")
		}
	}

	r.addToInflight(lggr, parsedReport)
	lggr.Infow("Accepting finalized report", "merkleRoot", hexutil.Encode(parsedReport.MerkleRoot[:]))
	return true, nil
}

func (r *CommitReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	parsedReport, err := DecodeCommitReport(report)
	if err != nil {
		return false, err
	}
	// If report is not stale we transmit.
	// When the commitTransmitter enqueues the tx for tx manager,
	// we mark it as fulfilled, effectively removing it from the set of inflight messages.
	return !r.isStaleReport(parsedReport), nil
}

func (r *CommitReportingPlugin) isStaleReport(report *commit_store.ICommitStoreCommitReport) bool {
	if isCommitStoreDownNow(r.config.lggr, r.config.commitStore) {
		return true
	}
	nextMin, err := r.config.commitStore.GetExpectedNextSequenceNumber(nil)
	if err != nil {
		// Assume it's a transient issue getting the last report
		// Will try again on the next round
		return true
	}
	// If the next min is already greater than this reports min,
	// this report is stale.
	if nextMin > report.Interval.Min {
		r.config.lggr.Warnw("report is stale", "onchain min", nextMin, "report min", report.Interval.Min)
		return true
	}
	return false
}

func (r *CommitReportingPlugin) Close() error {
	return nil
}
