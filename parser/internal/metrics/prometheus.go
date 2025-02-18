package metrics

import (
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics is a Prometheus implementation of Metric interface.
type PrometheusMetrics struct {
	vaaParseCount                 *prometheus.CounterVec
	vaaPayloadParserRequest       *prometheus.CounterVec
	vaaPayloadParserResponseCount *prometheus.CounterVec
}

// NewPrometheusMetrics returns a new instance of PrometheusMetrics.
func NewPrometheusMetrics(environment string) *PrometheusMetrics {
	vaaParseCount := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "parse_vaa_count_by_chain",
			Help: "Total number of vaa parser by chain",
			ConstLabels: map[string]string{
				"environment": environment,
				"service":     serviceName,
			},
		}, []string{"chain", "type"})
	vaaPayloadParserRequestCount := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "parse_vaa_payload_request_count_by_chain",
			Help: "Total number of request to payload parser component by chain",
			ConstLabels: map[string]string{
				"environment": environment,
				"service":     serviceName,
			},
		}, []string{"chain"})
	vaaPayloadParserResponseCount := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "parse_vaa_payload_response_count_by_chain",
			Help: "Total number of response from payload parser component by chain",
			ConstLabels: map[string]string{
				"environment": environment,
				"service":     serviceName,
			},
		}, []string{"chain", "status"})
	return &PrometheusMetrics{
		vaaParseCount:                 vaaParseCount,
		vaaPayloadParserRequest:       vaaPayloadParserRequestCount,
		vaaPayloadParserResponseCount: vaaPayloadParserResponseCount,
	}
}

// IncVaaConsumedQueue increments the number of consumed VAA.
func (m *PrometheusMetrics) IncVaaConsumedQueue(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaParseCount.WithLabelValues(chain, "consumed").Inc()
}

// IncVaaUnfiltered increments the number of unfiltered VAA.
func (m *PrometheusMetrics) IncVaaUnfiltered(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaParseCount.WithLabelValues(chain, "unfiltered").Inc()
}

// IncVaaUnexpired increments the number of unexpired VAA.
func (m *PrometheusMetrics) IncVaaUnexpired(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaParseCount.WithLabelValues(chain, "unexpired").Inc()
}

// IncVaaParsed increments the number of parsed VAA.
func (m *PrometheusMetrics) IncVaaParsed(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaParseCount.WithLabelValues(chain, "parsed").Inc()
}

// IncVaaParsedInserted increments the number of parsed VAA inserted into database.
func (m *PrometheusMetrics) IncVaaParsedInserted(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaParseCount.WithLabelValues(chain, "inserted").Inc()
}

// IncVaaPayloadParserRequestCount increments the number of vaa payload parser request.
func (m *PrometheusMetrics) IncVaaPayloadParserRequestCount(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaPayloadParserRequest.WithLabelValues(chain).Inc()
}

// IncVaaPayloadParserErrorCount increments the number of vaa payload parser error.
func (m *PrometheusMetrics) IncVaaPayloadParserErrorCount(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaPayloadParserResponseCount.WithLabelValues(chain, "failed").Inc()
}

// IncVaaPayloadParserSuccessCount increments the number of vaa payload parser success.
func (m *PrometheusMetrics) IncVaaPayloadParserSuccessCount(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaPayloadParserResponseCount.WithLabelValues(chain, "success").Inc()
}

// IncVaaPayloadParserSuccessCount increments the number of vaa payload parser not found.
func (m *PrometheusMetrics) IncVaaPayloadParserNotFoundCount(chainID uint16) {
	chain := vaa.ChainID(chainID).String()
	m.vaaPayloadParserResponseCount.WithLabelValues(chain, "not_found").Inc()
}
