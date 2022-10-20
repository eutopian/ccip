package dione

import (
	"fmt"

	null2 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/integration-tests/client"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/rhea"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

func RelaySpecToString(spec job.Job) string {
	onRamp := spec.OCR2OracleSpec.PluginConfig["onRampIDs"].([]string)[0]

	const relayTemplate = `
# CCIPRelaySpec
type               = "offchainreporting2"
name               = "%s"
pluginType         = "ccip-relay"
relay              = "evm"
schemaVersion      = 1
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"

[pluginConfig]
sourceChainID      = %s
destChainID        = %s
onRampIDs          = ["%s"]
pollPeriod         = "%s"
SourceStartBlock   = %d
DestStartBlock     = %d

[relayConfig]
chainID            = %s
`

	return fmt.Sprintf(relayTemplate+"\n",
		spec.Name.String,
		spec.OCR2OracleSpec.ContractID,
		spec.OCR2OracleSpec.OCRKeyBundleID.String,
		spec.OCR2OracleSpec.TransmitterID.String,
		spec.OCR2OracleSpec.PluginConfig["sourceChainID"],
		spec.OCR2OracleSpec.PluginConfig["destChainID"],
		onRamp,
		PollPeriod,
		spec.OCR2OracleSpec.PluginConfig["SourceStartBlock"],
		spec.OCR2OracleSpec.PluginConfig["DestStartBlock"],
		spec.OCR2OracleSpec.PluginConfig["destChainID"],
	)
}

func ExecSpecToString(spec job.Job) string {
	const execTemplate = `
# CCIPExecutionSpec
type               = "offchainreporting2"
name               = "%s"
pluginType         = "ccip-execution"
relay              = "evm"
schemaVersion      = 1
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"

[pluginConfig]
sourceChainID      = %s
destChainID        = %s
onRampID           = "%s"
blobVerifierID     = "%s"
SourceStartBlock   = %d
DestStartBlock     = %d
tokensPerFeeCoinPipeline = %s

[relayConfig]
chainID            = %s
`

	return fmt.Sprintf(execTemplate+"\n",
		spec.Name.String,
		spec.OCR2OracleSpec.ContractID,
		spec.OCR2OracleSpec.OCRKeyBundleID.String,
		spec.OCR2OracleSpec.TransmitterID.String,
		spec.OCR2OracleSpec.PluginConfig["sourceChainID"],
		spec.OCR2OracleSpec.PluginConfig["destChainID"],
		spec.OCR2OracleSpec.PluginConfig["onRampID"],
		spec.OCR2OracleSpec.PluginConfig["blobVerifierID"],
		spec.OCR2OracleSpec.PluginConfig["SourceStartBlock"],
		spec.OCR2OracleSpec.PluginConfig["DestStartBlock"],
		spec.OCR2OracleSpec.PluginConfig["tokensPerFeeCoinPipeline"],
		spec.OCR2OracleSpec.PluginConfig["destChainID"],
	)
}

func GetOCRkeysForChainType(OCRKeys client.OCR2Keys, chainType string) client.OCR2KeyData {
	for _, key := range OCRKeys.Data {
		if key.Attributes.ChainType == chainType {
			return key
		}
	}

	panic("Keys not found for chain")
}

func generateRelayJobSpecs(sourceClient *rhea.EvmChainConfig, destClient *rhea.EvmChainConfig) job.Job {
	return job.Job{
		Name: null2.StringFrom(fmt.Sprintf("ccip-relay-%s-%s", helpers.ChainName(sourceClient.ChainId.Int64()), helpers.ChainName(destClient.ChainId.Int64()))),
		Type: "offchainreporting2",
		OCR2OracleSpec: &job.OCR2OracleSpec{
			PluginType:                  job.CCIPRelay,
			ContractID:                  destClient.BlobVerifier.Hex(),
			Relay:                       "evm",
			RelayConfig:                 map[string]interface{}{"chainID": destClient.ChainId.String()},
			P2PV2Bootstrappers:          []string{},     // Set in env vars
			OCRKeyBundleID:              null2.String{}, // Set per node
			TransmitterID:               null2.String{}, // Set per node
			ContractConfigConfirmations: 2,
			PluginConfig: map[string]interface{}{
				"sourceChainID":    sourceClient.ChainId.String(),
				"destChainID":      destClient.ChainId.String(),
				"onRampIDs":        []string{sourceClient.OnRamp.String()},
				"pollPeriod":       PollPeriod,
				"SourceStartBlock": sourceClient.DeploySettings.DeployedAt,
				"DestStartBlock":   destClient.DeploySettings.DeployedAt,
			},
		},
	}
}

func generateExecutionJobSpecs(sourceClient *rhea.EvmChainConfig, destClient *rhea.EvmChainConfig) job.Job {
	return job.Job{
		Name: null2.StringFrom(fmt.Sprintf("ccip-exec-%s-%s", helpers.ChainName(sourceClient.ChainId.Int64()), helpers.ChainName(destClient.ChainId.Int64()))),
		Type: "offchainreporting2",
		OCR2OracleSpec: &job.OCR2OracleSpec{
			PluginType:                  job.CCIPExecution,
			ContractID:                  destClient.OffRamp.Hex(),
			Relay:                       "evm",
			RelayConfig:                 map[string]interface{}{"chainID": destClient.ChainId.String()},
			P2PV2Bootstrappers:          []string{},     // Set in env vars
			OCRKeyBundleID:              null2.String{}, // Set per node
			TransmitterID:               null2.String{}, // Set per node
			ContractConfigConfirmations: 2,
			PluginConfig: map[string]interface{}{
				"sourceChainID":            sourceClient.ChainId.String(),
				"destChainID":              destClient.ChainId.String(),
				"onRampID":                 sourceClient.OnRamp.String(),
				"pollPeriod":               PollPeriod,
				"blobVerifierID":           destClient.BlobVerifier.Hex(),
				"SourceStartBlock":         sourceClient.DeploySettings.DeployedAt,
				"DestStartBlock":           destClient.DeploySettings.DeployedAt,
				"tokensPerFeeCoinPipeline": fmt.Sprintf(`"""merge [type=merge left="{}" right="{\\\"%s\\\":\\\"1000000000000000000\\\"}"];"""`, destClient.LinkToken.Hex()),
			},
		},
	}
}
