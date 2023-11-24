package health

import (
	"context"
	"time"
)

// PhylaxCheck definition.
type PhylaxCheck struct {
	maxHealthTimeDuration time.Duration
	lastPing              time.Time
}

// NewPhylaxCheck instanciate a new PhylaxCheck
func NewPhylaxCheck(maxHealthTimeSeconds int64) *PhylaxCheck {
	return &PhylaxCheck{maxHealthTimeDuration: time.Duration(maxHealthTimeSeconds * int64(time.Second)), lastPing: time.Now()}
}

// Change last ping.
func (g *PhylaxCheck) Ping(ctx context.Context) {
	g.lastPing = time.Now()
}

// IsAlive check if the phylaxs are alive.
func (g *PhylaxCheck) IsAlive() bool {
	healthTime := time.Now().Add(-1 * g.maxHealthTimeDuration)
	return !g.lastPing.Before(healthTime)
}
