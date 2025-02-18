package metrics

// NoopMetrics is a no-op implementation of the Metrics interface.
type NoopMetrics struct {
}

// NewNoopMetrics returns a new instance of NoopMetrics.
func NewNoopMetrics() *NoopMetrics {
	return &NoopMetrics{}
}

func (p *NoopMetrics) IncFailedMeasurement(measurement string) {
}

func (p *NoopMetrics) IncSuccessfulMeasurement(measurement string) {
}

func (p *NoopMetrics) IncMissingNotional(symbol string) {
}

func (p *NoopMetrics) IncFoundNotional(symbol string) {
}

func (p *NoopMetrics) IncMissingToken(chain, token string) {
}

func (p *NoopMetrics) IncFoundToken(chain, token string) {
}
