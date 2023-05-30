package deployment_io

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/dione"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/rhea"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
)

func PrettyPrintLanes(env dione.Environment, source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	WriteChainConfigToFile(env, source)
	WriteChainConfigToFile(env, destination)

	writeLaneConfigToFile(env, source, source.LaneConfig, "lane")
	writeLaneConfigToFile(env, destination, destination.LaneConfig, "lane")

	writeLaneConfigToFile(env, source, source.UpgradeLaneConfig, "upgrade_lane")
	writeLaneConfigToFile(env, destination, destination.UpgradeLaneConfig, "upgrade_lane")
}

func WriteChainConfigToFile(env dione.Environment, chain *rhea.EvmDeploymentConfig) {
	sourceChainConfig := prettyPrint(chain.ChainConfig)
	chain.Logger.Info(string(sourceChainConfig))
	chainName := ccip.ChainName(int64(chain.ChainConfig.EvmChainId))
	attemptWriteToFile(sourceChainConfig, chainName, "chain", string(env))
}

func writeLaneConfigToFile(env dione.Environment, chain *rhea.EvmDeploymentConfig, config rhea.EVMLaneConfig, configType string) {
	destLaneConfig := prettyPrint(config)
	chain.Logger.Info(string(destLaneConfig))
	chainName := ccip.ChainName(int64(chain.ChainConfig.EvmChainId))
	attemptWriteToFile(destLaneConfig, chainName, configType, string(env))
}

func attemptWriteToFile(file []byte, chainName string, configType string, env string) {
	path := fmt.Sprintf("json/deployments/%s/%s/%s", env, configType, chainName)
	// ignore errors
	_ = os.MkdirAll(path, os.ModePerm)
	fileName := fmt.Sprintf("/%s", time.Now().Format("2006-01-02 15:04:05"))
	// ignore errors
	_ = os.WriteFile(path+fileName, file, 0600)
}

func prettyPrint(i interface{}) []byte {
	s, _ := json.MarshalIndent(i, "", "\t")
	return s
}
