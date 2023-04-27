package testhelpers

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/integration-tests/client"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type JobType string

const (
	Commit    JobType = "commit"
	Execution JobType = "exec"
	Boostrap  JobType = "bootstrap"
)

func JobName(jobType JobType, source string, destination string) string {
	return fmt.Sprintf("ccip-%s-%s-%s", jobType, source, destination)
}

type CCIPJobSpecParams struct {
	Name                   string
	OffRamp                common.Address
	CommitStore            common.Address
	SourceChainName        string
	DestChainName          string
	DestChainId            uint64
	TokenPricesUSDPipeline string
	SourceStartBlock       uint64
	DestStartBlock         uint64
	P2PV2Bootstrappers     pq.StringArray
}

func (params CCIPJobSpecParams) Validate() error {
	if params.CommitStore == common.HexToAddress("0x0") {
		return fmt.Errorf("must set commit store address")
	}
	return nil
}

func (params CCIPJobSpecParams) ValidateCommitJobSpec() error {
	commonErr := params.Validate()
	if commonErr != nil {
		return commonErr
	}
	if params.OffRamp == common.HexToAddress("0x0") {
		return fmt.Errorf("OffRamp cannot be empty for execution job")
	}
	return nil
}

func (params CCIPJobSpecParams) ValidateExecJobSpec() error {
	commonErr := params.Validate()
	if commonErr != nil {
		return commonErr
	}
	if params.OffRamp == common.HexToAddress("0x0") {
		return fmt.Errorf("OffRamp cannot be empty for execution job")
	}
	return nil
}

// CommitJobSpec generates template for CCIP-relay job spec.
// OCRKeyBundleID,TransmitterID need to be set from the calling function
func (params CCIPJobSpecParams) CommitJobSpec() (*client.OCR2TaskJobSpec, error) {
	err := params.ValidateCommitJobSpec()
	if err != nil {
		return nil, err
	}
	ocrSpec := job.OCR2OracleSpec{
		Relay:                             relay.EVM,
		PluginType:                        job.CCIPCommit,
		ContractID:                        params.CommitStore.Hex(),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),
		P2PV2Bootstrappers:                params.P2PV2Bootstrappers,
		PluginConfig: map[string]interface{}{
			"tokenPricesUSDPipeline": fmt.Sprintf(`"""
%s
"""`, params.TokenPricesUSDPipeline),
		},
		RelayConfig: map[string]interface{}{
			"chainID": params.DestChainId,
		},
	}
	if params.DestStartBlock > 0 {
		ocrSpec.PluginConfig["destStartBlock"] = params.DestStartBlock
	}
	if params.SourceStartBlock > 0 {
		ocrSpec.PluginConfig["sourceStartBlock"] = params.SourceStartBlock
	}
	return &client.OCR2TaskJobSpec{
		OCR2OracleSpec: ocrSpec,
		JobType:        "offchainreporting2",
		Name:           JobName(Commit, params.SourceChainName, params.DestChainName),
	}, nil
}

// ExecutionJobSpec generates template for CCIP-execution job spec.
// OCRKeyBundleID,TransmitterID need to be set from the calling function
func (params CCIPJobSpecParams) ExecutionJobSpec() (*client.OCR2TaskJobSpec, error) {
	err := params.ValidateExecJobSpec()
	if err != nil {
		return nil, err
	}
	ocrSpec := job.OCR2OracleSpec{
		Relay:                             relay.EVM,
		PluginType:                        job.CCIPExecution,
		ContractID:                        params.OffRamp.Hex(),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),
		P2PV2Bootstrappers:                params.P2PV2Bootstrappers,
		PluginConfig:                      map[string]interface{}{},
		RelayConfig: map[string]interface{}{
			"chainID": params.DestChainId,
		},
	}
	if params.DestStartBlock > 0 {
		ocrSpec.PluginConfig["destStartBlock"] = params.DestStartBlock
	}
	if params.SourceStartBlock > 0 {
		ocrSpec.PluginConfig["sourceStartBlock"] = params.SourceStartBlock
	}
	return &client.OCR2TaskJobSpec{
		OCR2OracleSpec: ocrSpec,
		JobType:        "offchainreporting2",
		Name:           JobName(Execution, params.SourceChainName, params.DestChainName),
	}, err
}

func (params CCIPJobSpecParams) BootstrapJob(contractID string) *client.OCR2TaskJobSpec {
	bootstrapSpec := job.OCR2OracleSpec{
		ContractID:                        contractID,
		Relay:                             relay.EVM,
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),
		RelayConfig: map[string]interface{}{
			"chainID": params.DestChainId,
		},
	}
	if params.DestStartBlock > 0 {
		bootstrapSpec.RelayConfig["fromBlock"] = params.DestStartBlock
	}
	return &client.OCR2TaskJobSpec{
		Name:           fmt.Sprintf("%s-%s", Boostrap, params.DestChainName),
		JobType:        "bootstrap",
		OCR2OracleSpec: bootstrapSpec,
	}
}
