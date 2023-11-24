package metrics

import sdk "github.com/deltaswapio/deltaswap/sdk/vaa"

const serviceName = "wormscan-fly"

type Metrics interface {
	// vaa metrics
	IncVaaFromGossipNetwork(chain sdk.ChainID)
	IncVaaUnfiltered(chain sdk.ChainID)
	IncVaaConsumedFromQueue(chain sdk.ChainID)
	IncVaaInserted(chain sdk.ChainID)
	IncVaaSendNotification(chain sdk.ChainID)
	IncVaaTotal()

	// observation metrics
	IncObservationFromGossipNetwork(chain sdk.ChainID)
	IncObservationUnfiltered(chain sdk.ChainID)
	IncObservationInserted(chain sdk.ChainID)
	IncObservationWithoutTxHash(chain sdk.ChainID)
	IncObservationTotal()

	// heartbeat metrics
	IncHeartbeatFromGossipNetwork(phylaxName string)
	IncHeartbeatInserted(phylaxName string)

	// governor config metrics
	IncGovernorConfigFromGossipNetwork(phylaxName string)
	IncGovernorConfigInserted(phylaxName string)

	// governor status metrics
	IncGovernorStatusFromGossipNetwork(phylaxName string)
	IncGovernorStatusInserted(phylaxName string)

	// max sequence cache metrics
	IncMaxSequenceCacheError(chain sdk.ChainID)
}
