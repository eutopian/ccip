package observability

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/price_registry"
)

type ObservedPriceRegistryReader struct {
	price_registry.PriceRegistryReader
	metric metricDetails
}

func NewPriceRegistryReader(origin price_registry.PriceRegistryReader, chainID int64, pluginName string) *ObservedPriceRegistryReader {
	return &ObservedPriceRegistryReader{
		PriceRegistryReader: origin,
		metric: metricDetails{
			interactionDuration: readerHistogram,
			resultSetSize:       readerDatasetSize,
			pluginName:          pluginName,
			readerName:          "PriceRegistryReader",
			chainId:             chainID,
		},
	}
}

func (o *ObservedPriceRegistryReader) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]ccipdata.Event[price_registry.TokenPriceUpdate], error) {
	return withObservedInteractionAndResults(o.metric, "GetTokenPriceUpdatesCreatedAfter", func() ([]ccipdata.Event[price_registry.TokenPriceUpdate], error) {
		return o.PriceRegistryReader.GetTokenPriceUpdatesCreatedAfter(ctx, ts, confs)
	})
}

func (o *ObservedPriceRegistryReader) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]ccipdata.Event[price_registry.GasPriceUpdate], error) {
	return withObservedInteractionAndResults(o.metric, "GetGasPriceUpdatesCreatedAfter", func() ([]ccipdata.Event[price_registry.GasPriceUpdate], error) {
		return o.PriceRegistryReader.GetGasPriceUpdatesCreatedAfter(ctx, chainSelector, ts, confs)
	})
}

func (o *ObservedPriceRegistryReader) GetFeeTokens(ctx context.Context) ([]common.Address, error) {
	return withObservedInteraction(o.metric, "GetFeeTokens", func() ([]common.Address, error) {
		return o.PriceRegistryReader.GetFeeTokens(ctx)
	})
}

func (o *ObservedPriceRegistryReader) GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]price_registry.TokenPriceUpdate, error) {
	return withObservedInteractionAndResults(o.metric, "GetTokenPrices", func() ([]price_registry.TokenPriceUpdate, error) {
		return o.PriceRegistryReader.GetTokenPrices(ctx, wantedTokens)
	})
}
