package ccipdataprovider

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
)

//go:generate mockery --quiet --name PriceRegistry --filename price_registry_mock.go --case=underscore
type PriceRegistry interface {
	NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error)
}

type EvmPriceRegistry struct {
	lp          logpoller.LogPoller
	ec          client.Client
	lggr        logger.Logger
	pluginLabel string
}

func NewEvmPriceRegistry(lp logpoller.LogPoller, ec client.Client, lggr logger.Logger, pluginLabel string) *EvmPriceRegistry {
	return &EvmPriceRegistry{
		lp:          lp,
		ec:          ec,
		lggr:        lggr,
		pluginLabel: pluginLabel,
	}
}

func (p *EvmPriceRegistry) NewPriceRegistryReader(_ context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error) {
	destPriceRegistryReader, err := factory.NewPriceRegistryReader(p.lggr, factory.NewEvmVersionFinder(), addr, p.lp, p.ec)
	if err != nil {
		return nil, err
	}
	return observability.NewPriceRegistryReader(destPriceRegistryReader, p.ec.ConfiguredChainID().Int64(), p.pluginLabel), nil
}
