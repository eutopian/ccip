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

type CCIPRelay struct {
	lggr              logger.Logger
	spec              *job.OCR2OracleSpec
	sourceChainPoller logpoller.LogPoller
	blobVerifier      *blob_verifier.BlobVerifier
	onRamp            *evm_2_evm_toll_onramp.EVM2EVMTollOnRamp
}

var _ plugins.OraclePlugin = &CCIPRelay{}

func NewCCIPRelay(lggr logger.Logger, spec *job.OCR2OracleSpec, chainSet evm.ChainSet) (*CCIPRelay, error) {
	var pluginConfig ccipconfig.RelayPluginConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return &CCIPRelay{}, err
	}
	err = pluginConfig.ValidateRelayPluginConfig()
	if err != nil {
		return &CCIPRelay{}, err
	}
	lggr.Infof("CCIP relay plugin initialized with offchainConfig: %+v", pluginConfig)

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
	blobVerifier, err := blob_verifier.NewBlobVerifier(common.HexToAddress(spec.ContractID), destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a new offramp")
	}
	if !common.IsHexAddress(string(pluginConfig.OnRampID)) {
		return nil, errors.Wrap(err, "OnRampID is not a valid hex address")
	}
	onRamp, err := evm_2_evm_toll_onramp.NewEVM2EVMTollOnRamp(common.HexToAddress(string(pluginConfig.OnRampID)), sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed creating a new onramp")
	}
	// Subscribe to all relevant relay logs.
	sourceChain.LogPoller().MergeFilter([]common.Hash{CCIPSendRequested}, onRamp.Address())
	return &CCIPRelay{
		lggr:              lggr,
		blobVerifier:      blobVerifier,
		onRamp:            onRamp,
		sourceChainPoller: sourceChain.LogPoller(),
	}, nil
}

func (c *CCIPRelay) GetPluginFactory() (plugin ocrtypes.ReportingPluginFactory, err error) {
	return NewRelayReportingPluginFactory(c.lggr, c.sourceChainPoller, c.blobVerifier, c.onRamp), nil
}

func (c *CCIPRelay) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
