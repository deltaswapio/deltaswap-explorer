package metrics

import (
	sdk "github.com/deltaswapio/deltaswap/sdk/vaa"
)

const serviceName = "deltaswapscan-contract-watcher"

type Metrics interface {
	SetLastBlock(chain sdk.ChainID, block uint64)
	SetCurrentBlock(chain sdk.ChainID, block uint64)
	IncDestinationTrxSaved(chain sdk.ChainID)
	IncRpcRequest(client string, method string, statusCode int)
}
