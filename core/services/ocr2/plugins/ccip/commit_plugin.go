package ccip

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/contractutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const (
	COMMIT_PRICE_UPDATES = "Commit price updates"
	COMMIT_CCIP_SENDS    = "Commit ccip sends"
)

func NewCommitServices(lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer, new bool, pr pipeline.Runner, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}
	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return nil, errors.New("chainID must be provided in relay config")
	}
	destChainID := int64(chainIDInterface.(float64))
	destChain, err := chainSet.Get(strconv.FormatInt(destChainID, 10))
	if err != nil {
		return nil, errors.Wrap(err, "get chainset")
	}
	commitStore, err := contractutil.LoadCommitStore(common.HexToAddress(spec.ContractID), CommitPluginLabel, destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading commitStore")
	}
	staticConfig, err := commitStore.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "failed getting the static config from the commitStore")
	}
	chainId, err := chainselectors.ChainIdFromSelector(staticConfig.SourceChainSelector)
	if err != nil {
		return nil, err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open source chain")
	}
	offRamp, err := contractutil.LoadOffRamp(common.HexToAddress(pluginConfig.OffRamp), CommitPluginLabel, destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading offRamp")
	}
	commitLggr := lggr.Named("CCIPCommit").With(
		"sourceChain", ChainName(int64(chainId)),
		"destChain", ChainName(destChainID))
	onRampReader, err := ccipdata.NewOnRampReader(commitLggr, staticConfig.SourceChainSelector, staticConfig.ChainSelector, staticConfig.OnRamp, sourceChain.LogPoller(), sourceChain.Client(), sourceChain.Config().EVM().FinalityTagEnabled())
	if err != nil {
		return nil, err
	}
	priceGetterObject, err := pricegetter.NewPipelineGetter(pluginConfig.TokenPricesUSDPipeline, pr, jb.ID, jb.ExternalJobID, jb.Name.ValueOrZero(), lggr)
	if err != nil {
		return nil, err
	}
	sourceRouter, err := router.NewRouter(onRampReader.Router(), sourceChain.Client())
	if err != nil {
		return nil, err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(nil)
	if err != nil {
		return nil, err
	}

	wrappedPluginFactory := NewCommitReportingPluginFactory(
		CommitPluginConfig{
			lggr:                commitLggr,
			sourceLP:            sourceChain.LogPoller(),
			destLP:              destChain.LogPoller(),
			onRampReader:        onRampReader,
			destReader:          ccipdata.NewLogPollerReader(destChain.LogPoller(), commitLggr, destChain.Client()),
			offRamp:             offRamp,
			onRampAddress:       staticConfig.OnRamp,
			priceGetter:         priceGetterObject,
			sourceNative:        sourceNative,
			sourceFeeEstimator:  sourceChain.GasEstimator(),
			sourceChainSelector: staticConfig.SourceChainSelector,
			destClient:          destChain.Client(),
			sourceClient:        sourceChain.Client(),
			commitStore:         commitStore,
		})

	err = wrappedPluginFactory.UpdateLogPollerFilters(utils.ZeroAddress, qopts...)
	if err != nil {
		return nil, err
	}

	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPCommit", string(spec.Relay), destChain.ID())
	argsNoPlugin.Logger = relaylogger.NewOCRWrapper(commitLggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	commitLggr.Infow("NewCommitServices",
		"pluginConfig", pluginConfig,
		"staticConfig", staticConfig,
		// TODO bring back
		//"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceNative,
		"sourceRouter", sourceRouter.Address())
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{oraclelib.NewBackfilledOracle(
			commitLggr,
			sourceChain.LogPoller(),
			destChain.LogPoller(),
			pluginConfig.SourceStartBlock,
			pluginConfig.DestStartBlock,
			job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

// CommitReportToEthTxMeta generates a txmgr.EthTxMeta from the given commit report.
// sequence numbers of the committed messages will be added to tx metadata
func CommitReportToEthTxMeta(report []byte) (*txmgr.TxMeta, error) {
	commitReport, err := abihelpers.DecodeCommitReport(report)
	if err != nil {
		return nil, err
	}
	n := int(commitReport.Interval.Max-commitReport.Interval.Min) + 1
	seqRange := make([]uint64, n)
	for i := 0; i < n; i++ {
		seqRange[i] = uint64(i) + commitReport.Interval.Min
	}
	return &txmgr.TxMeta{
		SeqNumbers: seqRange,
	}, nil
}

func getCommitPluginSourceLpFilters(onRamp common.Address) []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(COMMIT_CCIP_SENDS, onRamp.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.SendRequested},
			Addresses: []common.Address{onRamp},
		},
	}
}

func getCommitPluginDestLpFilters(priceRegistry common.Address, offRamp common.Address) []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(COMMIT_PRICE_UPDATES, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.UsdPerUnitGasUpdated, abihelpers.EventSignatures.UsdPerTokenUpdated},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistry),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenAdded},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistry),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenRemoved},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_ADDED, offRamp),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolAdded},
			Addresses: []common.Address{offRamp},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_REMOVED, offRamp),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolRemoved},
			Addresses: []common.Address{offRamp},
		},
	}
}

// UnregisterCommitPluginLpFilters unregisters all the registered filters for both source and dest chains.
func UnregisterCommitPluginLpFilters(ctx context.Context, lggr logger.Logger, spec *job.OCR2OracleSpec, chainSet evm.LegacyChainContainer, qopts ...pg.QOpt) error {
	if spec == nil {
		return errors.New("spec is nil")
	}
	if !common.IsHexAddress(spec.ContractID) {
		return fmt.Errorf("invalid contract id address: %s", spec.ContractID)
	}

	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return err
	}

	destChainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return errors.New("chainID must be provided in relay config")
	}
	destChainIDf64, is := destChainIDInterface.(float64)
	if !is {
		return fmt.Errorf("chain id '%v' is not float64", destChainIDInterface)
	}
	destChainID := int64(destChainIDf64)
	destChain, err := chainSet.Get(strconv.FormatInt(destChainID, 10))
	if err != nil {
		return err
	}
	commitStore, err := contractutil.LoadCommitStore(common.HexToAddress(spec.ContractID), CommitPluginLabel, destChain.Client())
	if err != nil {
		return err
	}
	staticConfig, err := commitStore.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return err
	}
	chainId, err := chainselectors.ChainIdFromSelector(staticConfig.SourceChainSelector)
	if err != nil {
		return err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return err
	}
	onRampReader, _ := ccipdata.NewOnRampReader(lggr, staticConfig.SourceChainSelector, staticConfig.ChainSelector, staticConfig.OnRamp, sourceChain.LogPoller(), sourceChain.Client(), sourceChain.Config().EVM().FinalityTagEnabled())
	if err := onRampReader.Close(); err != nil {
		return err
	}

	// TODO: once offramp/commit are abstracted, we can call Close on the offramp/commit readers to unregister filters.
	return unregisterCommitPluginFilters(ctx, sourceChain.LogPoller(), destChain.LogPoller(), commitStore, common.HexToAddress(pluginConfig.OffRamp), qopts...)
}

func unregisterCommitPluginFilters(ctx context.Context, sourceLP, destLP logpoller.LogPoller, destCommitStore commit_store.CommitStoreInterface, offRamp common.Address, qopts ...pg.QOpt) error {
	//staticCfg, err := destCommitStore.GetStaticConfig(&bind.CallOpts{Context: ctx})
	//if err != nil {
	//	return err
	//}

	dynamicCfg, err := destCommitStore.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}

	//if err := logpollerutil.UnregisterLpFilters(
	//	sourceLP,
	//	getCommitPluginSourceLpFilters(staticCfg.OnRamp),
	//	qopts...,
	//); err != nil {
	//	return err
	//}

	return logpollerutil.UnregisterLpFilters(
		destLP,
		getCommitPluginDestLpFilters(dynamicCfg.PriceRegistry, offRamp),
		qopts...,
	)
}
