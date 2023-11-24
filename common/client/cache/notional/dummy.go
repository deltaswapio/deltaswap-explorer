package notional

import "github.com/deltaswapio/deltaswap-explorer/common/domain"

// DummyNotionalCache is a dummy notional cache.
type DummyNotionalCache struct {
}

// NewDummyNotionalCache init a new dummy notional cache.
func NewDummyNotionalCache() *DummyNotionalCache {
	return &DummyNotionalCache{}
}

// Get get notional cache value.
func (c *DummyNotionalCache) Get(symbol domain.Symbol) (PriceData, error) {
	return PriceData{}, nil
}

// Close the dummy cache.
func (c *DummyNotionalCache) Close() error {
	return nil
}
