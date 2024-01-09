package liquiditymanager

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// Factory initializes a new liquidity manager instance.
//
//go:generate mockery --quiet --name Factory --output ../rebalancermocks --filename lm_factory_mock.go --case=underscore
type Factory interface {
	NewLiquidityManager(networkID models.NetworkID, address models.Address) (LiquidityManager, error)
}

type evmDep struct {
	lp        logpoller.LogPoller
	ethClient client.Client
}

type BaseLiquidityManagerFactory struct {
	evmDeps map[models.NetworkID]evmDep
}

func NewBaseLiquidityManagerFactory(opts ...func(f *BaseLiquidityManagerFactory)) *BaseLiquidityManagerFactory {
	f := &BaseLiquidityManagerFactory{
		evmDeps: make(map[models.NetworkID]evmDep),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithEvmDep(networkID models.NetworkID, lp logpoller.LogPoller, ethClient client.Client) func(f *BaseLiquidityManagerFactory) {
	return func(f *BaseLiquidityManagerFactory) {
		f.evmDeps[networkID] = evmDep{
			lp:        lp,
			ethClient: ethClient,
		}
	}
}

func (b *BaseLiquidityManagerFactory) NewLiquidityManager(networkID models.NetworkID, address models.Address) (LiquidityManager, error) {
	switch typ := networkID.Type(); typ {
	case models.NetworkTypeEvm:
		evmDeps, exists := b.evmDeps[networkID]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found")
		}
		return NewEvmLiquidityManager(address, evmDeps.ethClient, evmDeps.lp)
	default:
		return nil, fmt.Errorf("liquidity manager of type %v is not supported", typ)
	}
}
