package tvl

import (
	"context"
	"errors"
	"time"

	errs "github.com/deltaswapio/deltaswap-explorer/api/internal/errors"
	deltaswapscanCache "github.com/deltaswapio/deltaswap-explorer/common/client/cache"
	"go.uber.org/zap"
)

// Tvl is the tvl client.
type Tvl struct {
	api        *TvlAPI
	cache      deltaswapscanCache.Cache
	tvlKey     string
	expiration time.Duration
	logger     *zap.Logger
}

// NewTVL init a new tvl client.
func NewTVL(p2pNetwork string, cache deltaswapscanCache.Cache, tvlKey string, expiration int, logger *zap.Logger) *Tvl {
	return &Tvl{
		api:        NewTvlAPI(p2pNetwork),
		cache:      cache,
		tvlKey:     tvlKey,
		expiration: time.Duration(expiration) * time.Second,
		logger:     logger}
}

// Get get tvl value from cache if exists or call wormhole api to get tvl value and set the in cache for t.expiration time.
func (t *Tvl) Get(ctx context.Context) (string, error) {

	// Get tvl from cache
	tvl, err := t.cache.Get(ctx, t.tvlKey)
	if err == nil {
		return tvl, nil
	}
	if errors.Is(err, deltaswapscanCache.ErrInternal) {
		t.logger.Error("error getting tvl from cache",
			zap.Error(err),
			zap.String("key", t.tvlKey))
	}

	// Get tvl from wormhole api
	tvlUSD, err := t.api.GetNotionalUSD(ctx, []string{"all"})
	if err != nil {
		t.logger.Error("error getting tvl from wormhole api",
			zap.Error(err))
	}
	if tvlUSD == nil {
		return "", errs.ErrNotFound
	}

	// Set tvl in cache with t.expiration time
	err = t.cache.Set(ctx, t.tvlKey, *tvlUSD, t.expiration)
	if err != nil {
		t.logger.Error("error setting tvl in cache",
			zap.Error(err),
			zap.String("key", t.tvlKey))
	}
	return *tvlUSD, nil
}
