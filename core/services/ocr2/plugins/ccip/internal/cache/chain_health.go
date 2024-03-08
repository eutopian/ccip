package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

// ChainHealthcheck checks the health of the both source and destination chain.
// Based on the values returned, CCIP can make a decision to stop or continue processing messages.
// There are four things verified here:
// 1. Source chain is healthy (this is verified by checking if source LogPoller saw finality violation)
// 2. Dest chain is healthy (this is verified by checking if destination LogPoller saw finality violation)
// 3. CommitStore is down (this is verified by checking if CommitStore is down and destination RMN is not cursed)
// 4. Source chain is cursed (this is verified by checking if source RMN is not cursed)
//
// Whenever any of the above checks fail, the chain is considered unhealthy and the CCIP should stop
// processing messages. Additionally, when the chain is unhealthy, this information is considered "sticky"
// and is cached for a certain period of time based on defaultGlobalStatusDuration.
// This may lead to some false-positives, but in this case we want to be extra cautious and avoid executing any reorged messages.
//
// Additionally, to reduce the number of calls to the RPC, we refresh RMN state in the background based on defaultRMNStateRefreshInterval
//
//go:generate mockery --quiet --name ChainHealthcheck --filename chain_health_mock.go --case=underscore
type ChainHealthcheck interface {
	// IsHealthy checks if the chain is healthy and returns true if it is, false otherwise
	// If forceRefresh is set to true, it will refresh the RMN curse state. Should be used in the Observation and ShouldTransmit phases of OCR2.
	// Otherwise, it will use the cached value of the RMN curse state.
	IsHealthy(ctx context.Context) (bool, error)
}

const (
	// RMN curse state is refreshed every 20 seconds or when ForceIsHealthy is called
	defaultRMNStateRefreshInterval = 10 * time.Second
	defaultGlobalStatusDuration    = 30 * time.Minute

	globalStatusKey = "globalStatus"
	rmnStatusKey    = "rmnCurseCheck"
)

type chainHealthcheck struct {
	cache                    *cache.Cache
	globalStatusKey          string
	rmnStatusKey             string
	globalStatusExpiration   time.Duration
	rmnStatusRefreshInterval time.Duration

	lggr        logger.Logger
	onRamp      ccipdata.OnRampReader
	commitStore ccipdata.CommitStoreReader
}

func NewChainHealthcheck(ctx context.Context, lggr logger.Logger, onRamp ccipdata.OnRampReader, commitStore ccipdata.CommitStoreReader) *chainHealthcheck {
	return newChainHealthcheckWithCustomEviction(
		ctx,
		lggr,
		onRamp,
		commitStore,
		defaultGlobalStatusDuration,
		defaultRMNStateRefreshInterval,
	)
}

func newChainHealthcheckWithCustomEviction(ctx context.Context, lggr logger.Logger, onRamp ccipdata.OnRampReader, commitStore ccipdata.CommitStoreReader, globalStatusDuration time.Duration, rmnStatusRefreshInterval time.Duration) *chainHealthcheck {
	ch := &chainHealthcheck{
		cache:                    cache.New(rmnStatusRefreshInterval, 0),
		rmnStatusKey:             rmnStatusKey,
		globalStatusKey:          globalStatusKey,
		globalStatusExpiration:   globalStatusDuration,
		rmnStatusRefreshInterval: rmnStatusRefreshInterval,

		lggr:        lggr,
		onRamp:      onRamp,
		commitStore: commitStore,
	}
	ch.spawnBackgroundRefresher(ctx)
	return ch
}

type rmnResponse struct {
	healthy bool
	err     error
}

func (c *chainHealthcheck) IsHealthy(ctx context.Context) (bool, error) {
	// Verify if flag is raised to indicate that the chain is not healthy
	// If set to false then immediately return false without checking the chain
	if cachedValue, found := c.cache.Get(c.globalStatusKey); found {
		healthy, ok := cachedValue.(bool)
		// If cached value is properly casted to bool and not healthy it means the sticky flag is raised
		// and should be returned immediately
		if !ok {
			c.lggr.Criticalw("Failed to cast cached value to sticky healthcheck", "value", cachedValue)
		} else if ok && !healthy {
			return false, nil
		}
	}

	// These checks are cheap and don't require any communication with the database or RPC
	if healthy, err := c.checkIfReadersAreHealthy(ctx); err != nil {
		return false, err
	} else if !healthy {
		c.cache.Set(c.globalStatusKey, false, c.globalStatusExpiration)
		return healthy, nil
	}

	// First call might initialize cache if it's not initialized yet. Otherwise, it will use the cached value
	if healthy, err := c.checkIfRMNsAreHealthy(ctx); err != nil {
		return false, err
	} else if !healthy {
		c.cache.Set(c.globalStatusKey, false, c.globalStatusExpiration)
		return healthy, nil
	}
	return true, nil
}

func (c *chainHealthcheck) spawnBackgroundRefresher(ctx context.Context) {
	ticker := time.NewTicker(c.rmnStatusRefreshInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Refresh the RMN state
				c.refresh(ctx)
			}
		}
	}()
}

func (c *chainHealthcheck) refresh(ctx context.Context) {
	healthy, err := c.fetchRMNCurseState(ctx)
	c.cache.Set(c.rmnStatusKey, rmnResponse{healthy, err}, -1)
}

// checkIfReadersAreHealthy checks if the source and destination chains are healthy by calling underlying LogPoller
// These calls are cheap because they don't require any communication with the database or RPC, so we don't have
// to cache the result of these calls.
func (c *chainHealthcheck) checkIfReadersAreHealthy(ctx context.Context) (bool, error) {
	sourceChainHealthy, err := c.onRamp.IsSourceChainHealthy(ctx)
	if err != nil {
		return false, errors.Wrap(err, "onRamp IsSourceChainHealthy errored")
	}

	destChainHealthy, err := c.commitStore.IsDestChainHealthy(ctx)
	if err != nil {
		return false, errors.Wrap(err, "commitStore IsDestChainHealthy errored")
	}

	if !sourceChainHealthy || !destChainHealthy {
		c.lggr.Criticalw(
			"Lane processing is stopped because source or destination chain is reported unhealthy",
			"sourceChainHealthy", sourceChainHealthy,
			"destChainHealthy", destChainHealthy,
		)
	}
	return sourceChainHealthy && destChainHealthy, nil
}

func (c *chainHealthcheck) checkIfRMNsAreHealthy(ctx context.Context) (bool, error) {
	if cachedValue, found := c.cache.Get(c.rmnStatusKey); found {
		rmn, ok := cachedValue.(rmnResponse)
		if ok {
			return rmn.healthy, rmn.err
		}
		c.lggr.Criticalw("Failed to cast cached value to RMN response", "response", rmn)
	}

	// If the value is not found in the cache, fetch the RMN curse state in a sync manner for the first time
	healthy, err := c.fetchRMNCurseState(ctx)
	c.cache.Set(c.rmnStatusKey, rmnResponse{healthy, err}, -1)
	return healthy, err
}

func (c *chainHealthcheck) fetchRMNCurseState(ctx context.Context) (bool, error) {
	var (
		eg                = new(errgroup.Group)
		isCommitStoreDown bool
		isSourceCursed    bool
	)

	eg.Go(func() error {
		var err error
		isCommitStoreDown, err = c.commitStore.IsDown(ctx)
		if err != nil {
			return errors.Wrap(err, "commitStore isDown check errored")
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		isSourceCursed, err = c.onRamp.IsSourceCursed(ctx)
		if err != nil {
			return errors.Wrap(err, "onRamp isSourceCursed errored")
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return false, err
	}

	if isCommitStoreDown || isSourceCursed {
		c.lggr.Criticalw(
			"Lane processing is stopped because source chain is cursed or CommitStore is down",
			"isCommitStoreDown", isCommitStoreDown,
			"isSourceCursed", isSourceCursed,
		)
		return false, nil
	}
	return true, nil
}
