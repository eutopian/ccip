package rhea

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/ping_pong_demo"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func setOffRampOnTokenPools(t *testing.T, destClient *EvmDeploymentConfig) {
	for _, tokenConfig := range destClient.ChainConfig.SupportedTokens {
		pool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenConfig.Pool, destClient.Client)
		shared.RequireNoError(t, err)

		rampUpdate := lock_release_token_pool.IPoolRampUpdate{
			Ramp:    destClient.LaneConfig.OffRamp,
			Allowed: true,
		}

		// Configure offramp address on pool
		tx, err := pool.ApplyRampUpdates(destClient.Owner, []lock_release_token_pool.IPoolRampUpdate{}, []lock_release_token_pool.IPoolRampUpdate{rampUpdate})
		shared.RequireNoError(t, err)
		shared.WaitForMined(t, destClient.Logger, destClient.Client, tx.Hash(), true)
		destClient.Logger.Infof("Offramp pool configured with offramp address: %s", helpers.ExplorerLink(int64(destClient.ChainConfig.ChainId), tx.Hash()))
	}
}

func SetPriceRegistryPrices(t *testing.T, client *EvmDeploymentConfig, destChainId uint64) {
	priceRegistry, err := price_registry.NewPriceRegistry(client.ChainConfig.PriceRegistry, client.Client)
	shared.RequireNoError(t, err)

	priceUpdates := price_registry.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry.InternalTokenPriceUpdate{},
		DestChainId:       destChainId,
		// Set 1e18 units of gas to $2k, being fairly reasonable for eth
		// These values will get auto updated by the DON
		UsdPerUnitGas: big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
	}

	for _, feeToken := range client.ChainConfig.FeeTokens {
		priceUpdates.TokenPriceUpdates = append(priceUpdates.TokenPriceUpdates, price_registry.InternalTokenPriceUpdate{
			SourceToken: client.ChainConfig.SupportedTokens[feeToken].Token,
			// USD per wei.
			UsdPerToken: client.ChainConfig.SupportedTokens[feeToken].Price,
		})
	}

	tx, err := priceRegistry.UpdatePrices(client.Owner, priceUpdates)
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
}

func setOnRampPrices(t *testing.T, client *EvmDeploymentConfig) {
	var tokens []common.Address
	var prices []*big.Int
	for _, tokenConfig := range client.ChainConfig.SupportedTokens {
		tokens = append(tokens, tokenConfig.Token)
		prices = append(prices, tokenConfig.Price)
	}

	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(client.LaneConfig.OnRamp, client.Client)
	shared.RequireNoError(t, err)

	// Prices are used by the rate limiter and dictate what tokens are supported
	tx, err := onRamp.SetPrices(client.Owner, tokens, prices)
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
	client.Logger.Infof(fmt.Sprintf("OnRamp prices set on %s in tx %s", client.LaneConfig.OnRamp.String(), helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash())))
}

func setOffRampPrices(t *testing.T, client *EvmDeploymentConfig) {
	var tokens []common.Address
	var prices []*big.Int
	for _, tokenConfig := range client.ChainConfig.SupportedTokens {
		tokens = append(tokens, tokenConfig.Token)
		prices = append(prices, tokenConfig.Price)
	}

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(client.LaneConfig.OffRamp, client.Client)
	shared.RequireNoError(t, err)

	// Prices are used by the rate limiter and dictate what tokens are supported
	tx, err := offRamp.SetPrices(client.Owner, tokens, prices)
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
	client.Logger.Infof(fmt.Sprintf("OffRamp prices set on %s in tx %s", client.LaneConfig.OnRamp.String(), helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash())))
}

func setOnRampOnRouter(t *testing.T, sourceClient *EvmDeploymentConfig, destChainId uint64) {
	sourceClient.Logger.Infof("Setting the onRamp on the Router")
	routerContract, err := router.NewRouter(sourceClient.ChainConfig.Router, sourceClient.Client)
	shared.RequireNoError(t, err)
	sourceClient.Logger.Infof("Registering new onRamp")
	tx, err := routerContract.ApplyRampUpdates(sourceClient.Owner, []router.RouterOnRampUpdate{{DestChainId: destChainId, OnRamp: sourceClient.LaneConfig.OnRamp}}, nil)
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, sourceClient.Logger, sourceClient.Client, tx.Hash(), true)
}

func setOnRampOnTokenPools(t *testing.T, sourceClient *EvmDeploymentConfig) {
	for _, tokenConfig := range sourceClient.ChainConfig.SupportedTokens {
		pool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenConfig.Pool, sourceClient.Client)
		shared.RequireNoError(t, err)

		rampUpdate := lock_release_token_pool.IPoolRampUpdate{
			Ramp:    sourceClient.LaneConfig.OnRamp,
			Allowed: true,
		}

		// Configure offramp address on pool
		tx, err := pool.ApplyRampUpdates(sourceClient.Owner, []lock_release_token_pool.IPoolRampUpdate{rampUpdate}, []lock_release_token_pool.IPoolRampUpdate{})
		shared.RequireNoError(t, err)
		shared.WaitForMined(t, sourceClient.Logger, sourceClient.Client, tx.Hash(), true)
		sourceClient.Logger.Infof("Onramp pool configured with offramp address: %s", helpers.ExplorerLink(int64(sourceClient.ChainConfig.ChainId), tx.Hash()))
	}
}

func setOffRampOnRouter(t *testing.T, sourceChainId uint64, client *EvmDeploymentConfig) {
	client.Logger.Infof("Setting the offRamp on the Router")
	routerContract, err := router.NewRouter(client.ChainConfig.Router, client.Client)
	shared.RequireNoError(t, err)

	offRamps, err := routerContract.GetOffRamps(&bind.CallOpts{}, sourceChainId)
	shared.RequireNoError(t, err)
	for _, offRamp := range offRamps {
		if offRamp == client.LaneConfig.OffRamp {
			client.Logger.Infof("OffRamp already configured on router. Skipping")
			return
		}
	}

	tx, err := routerContract.ApplyRampUpdates(client.Owner, nil, []router.RouterOffRampUpdate{
		{SourceChainId: sourceChainId, OffRamps: []common.Address{client.LaneConfig.OffRamp}}})
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
}

func setPriceRegistryUpdater(t *testing.T, client *EvmDeploymentConfig) {
	priceRegistry, err := price_registry.NewPriceRegistry(client.ChainConfig.PriceRegistry, client.Client)
	shared.RequireNoError(t, err)

	tx, err := priceRegistry.ApplyPriceUpdatersUpdates(client.Owner, []common.Address{client.LaneConfig.CommitStore}, []common.Address{})
	shared.RequireNoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
}

func fillPoolWithTokens(t *testing.T, client *EvmDeploymentConfig, pool *lock_release_token_pool.LockReleaseTokenPool, tokenAddress common.Address) {
	token, err := link_token_interface.NewLinkToken(tokenAddress, client.Client)
	shared.RequireNoError(t, err)

	// fill offramp token pool with 0.5 token
	amount := big.NewInt(5e17)
	tx, err := token.Approve(client.Owner, pool.Address(), amount)
	shared.RequireNoError(t, err)
	client.Logger.Infof("Approving token to the token pool: %s", helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash()))
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)

	tx, err = pool.AddLiquidity(client.Owner, amount)
	shared.RequireNoError(t, err)
	client.Logger.Infof("Adding liquidity to the token pool: %s", helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash()))
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)

	client.Logger.Infof("Pool filled with tokens: %s", helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash()))
}

func FundPingPong(t *testing.T, client *EvmDeploymentConfig, fundingAmount *big.Int, tokenAddress common.Address) {
	pingDapp, err := ping_pong_demo.NewPingPongDemo(client.LaneConfig.PingPongDapp, client.Client)
	require.NoError(t, err)

	linkToken, err := link_token_interface.NewLinkToken(tokenAddress, client.Client)
	require.NoError(t, err)

	tx, err := linkToken.Approve(client.Owner, client.LaneConfig.PingPongDapp, fundingAmount)
	require.NoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)

	tx, err = pingDapp.Fund(client.Owner, fundingAmount)
	require.NoError(t, err)
	shared.WaitForMined(t, client.Logger, client.Client, tx.Hash(), true)
	client.Logger.Infof("Ping pong funded with %s in tx: %s", fundingAmount.String(), helpers.ExplorerLink(int64(client.ChainConfig.ChainId), tx.Hash()))
}

func UpdateDeployedAt(t *testing.T, source *EvmDeploymentConfig, dest *EvmDeploymentConfig) {
	sourceBlock, err := source.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	source.ChainConfig.DeploySettings.DeployedAtBlock = sourceBlock
	source.LaneConfig.DeploySettings.DeployedAtBlock = sourceBlock

	destBlock, err := dest.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	dest.ChainConfig.DeploySettings.DeployedAtBlock = destBlock
	dest.LaneConfig.DeploySettings.DeployedAtBlock = destBlock
}
