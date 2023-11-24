package metrics

const serviceName = "deltaswapscan-tx-tracker"

type Metrics interface {
	IncVaaConsumedQueue(chainID uint16)
	IncVaaUnfiltered(chainID uint16)
	IncOriginTxInserted(chainID uint16)
}
