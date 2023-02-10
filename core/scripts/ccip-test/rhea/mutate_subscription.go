package rhea

import "testing"

func UpgradeLane(t *testing.T, sourceClient *EvmDeploymentConfig, destClient *EvmDeploymentConfig) {
	if !sourceClient.DeploySettings.DeployRamp || !destClient.DeploySettings.DeployRamp {
		sourceClient.Logger.Errorf("Please set \"DeployRamp\" to true for the given EvmChainConfigs and make sure "+
			"the right ones are set. Source: %d, Dest %d", sourceClient.ChainConfig.ChainId, destClient.ChainConfig.ChainId)
		return
	}

	upgradeOnRamp(t, sourceClient, destClient)
	upgradeOffRamp(t, sourceClient, destClient)
}

func upgradeOnRamp(t *testing.T, sourceClient *EvmDeploymentConfig, destClient *EvmDeploymentConfig) {
	sourceClient.Logger.Infof("Upgrading onRamp")
	deployOnRamp(t, sourceClient, destClient.ChainConfig.ChainId, destClient.ChainConfig.SupportedTokens)
	setOnRampOnTokenPools(t, sourceClient)

	sourceClient.Logger.Info("Please deploy new commit jobs")
}

func upgradeOffRamp(t *testing.T, sourceClient *EvmDeploymentConfig, destClient *EvmDeploymentConfig) {
	destClient.Logger.Infof("Upgrading offRamp")
	deployOffRamp(t, destClient, sourceClient.ChainConfig.ChainId, sourceClient.ChainConfig.SupportedTokens, sourceClient.LaneConfig.OnRamp)
	setRouterOnOffRamp(t, destClient)
	setOffRampOnRouter(t, destClient)
	setOffRampOnTokenPools(t, destClient)

	destClient.Logger.Info("Please deploy new execution jobs")
}

/*
func removeOffRamp(t *testing.T, destClient *EvmDeploymentConfig, offRampAddress common.Address) {
	// Pause contract
	revokeOffRampOnRouter(t, destClient, offRampAddress)
	revokeOffRampOnTokenPools(t, destClient, offRampAddress)
}
*/
