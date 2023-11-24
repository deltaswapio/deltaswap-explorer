package metrics

const serviceName = "deltaswapscan-analytics"

type Metrics interface {
	IncFailedMeasurement(measurement string)
	IncSuccessfulMeasurement(measurement string)
	IncMissingNotional(symbol string)
	IncFoundNotional(symbol string)
	IncMissingToken(chain, token string)
	IncFoundToken(chain, token string)
}
