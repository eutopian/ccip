package ccip

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blob_verifier"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/evm_2_evm_toll_onramp"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	ccipconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ccip/config"
)

type CCIPExecution struct {
	lggr                               logger.Logger
	spec                               *job.OCR2OracleSpec
	sourceChainPoller, destChainPoller logpoller.LogPoller
	destChain                          evm.Chain
	blobVerifier                       *blob_verifier.BlobVerifier
	onRamp                             *evm_2_evm_toll_onramp.EVM2EVMTollOnRamp
}

var _ plugins.OraclePlugin = &CCIPExecution{}

func NewCCIPExecution(lggr logger.Logger, spec *job.OCR2OracleSpec, chainSet evm.ChainSet) (*CCIPExecution, error) {
	var pluginConfig ccipconfig.ExecutionPluginConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return &CCIPExecution{}, err
	}
	err = pluginConfig.ValidateExecutionPluginConfig()
	if err != nil {
		return &CCIPExecution{}, err
	}
	lggr.Infof("CCIP execution plugin initialized with offchainConfig: %+v", pluginConfig)

	sourceChain, err := chainSet.Get(big.NewInt(pluginConfig.SourceChainID))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open source chain")
	}
	destChain, err := chainSet.Get(big.NewInt(pluginConfig.DestChainID))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open destination chain")
	}

	if !common.IsHexAddress(spec.ContractID) {
		return nil, errors.Wrap(err, "spec.OffRampID is not a valid hex address")
	}
	offRamp, err := blob_verifier.NewBlobVerifier(common.HexToAddress(string(pluginConfig.OffRampID)), destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a new onramp")
	}
	if !common.IsHexAddress(string(pluginConfig.OnRampID)) {
		return nil, errors.Wrap(err, "spec.OffRampID is not a valid hex address")
	}
	onRamp, err := evm_2_evm_toll_onramp.NewEVM2EVMTollOnRamp(common.HexToAddress(string(pluginConfig.OnRampID)), sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a new onramp")
	}

	// Subscribe to all relevant execution logs.
	sourceChain.LogPoller().MergeFilter([]common.Hash{CCIPSendRequested}, onRamp.Address())
	destChain.LogPoller().MergeFilter([]common.Hash{ReportAccepted}, offRamp.Address())
	destChain.LogPoller().MergeFilter([]common.Hash{CrossChainMessageExecuted}, offRamp.Address())
	return &CCIPExecution{
		lggr:              lggr,
		spec:              spec,
		blobVerifier:      offRamp,
		onRamp:            onRamp,
		sourceChainPoller: sourceChain.LogPoller(),
		destChainPoller:   destChain.LogPoller(),
	}, nil
}

func (c *CCIPExecution) GetPluginFactory() (plugin ocrtypes.ReportingPluginFactory, err error) {
	return NewExecutionReportingPluginFactory(
		c.lggr,
		c.onRamp,
		c.blobVerifier,
		c.sourceChainPoller,
		c.destChainPoller,
		common.HexToAddress(c.spec.ContractID),
	), nil
}

func (c *CCIPExecution) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
