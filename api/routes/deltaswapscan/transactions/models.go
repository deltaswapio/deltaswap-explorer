package transactions

import (
	"time"

	"github.com/deltaswapio/deltaswap-explorer/api/handlers/transactions"
	sdk "github.com/deltaswapio/deltaswap/sdk/vaa"
)

// TransactionDetail is a brief description of a transaction (e.g. ID, txHash, payload, etc.)
type TransactionDetail struct {
	ID           string      `json:"id"`
	Timestamp    time.Time   `json:"timestamp"`
	TxHash       string      `json:"txHash,omitempty"`
	EmitterChain sdk.ChainID `json:"emitterChain"`
	// EmitterAddress contains the VAA's emitter address, encoded in hex.
	EmitterAddress string `json:"emitterAddress"`
	// EmitterNativeAddress contains the VAA's emitter address, encoded in the emitter chain's native format.
	EmitterNativeAddress   string                             `json:"emitterNativeAddress,omitempty"`
	TokenAmount            string                             `json:"tokenAmount,omitempty"`
	UsdAmount              string                             `json:"usdAmount,omitempty"`
	Symbol                 string                             `json:"symbol,omitempty"`
	Payload                map[string]interface{}             `json:"payload,omitempty"`
	StandardizedProperties map[string]interface{}             `json:"standardizedProperties,omitempty"`
	GlobalTx               *transactions.GlobalTransactionDoc `json:"globalTx,omitempty"`
}

// ListTransactionsResponse is the "200 OK" response model for `GET /api/v1/transactions`.
type ListTransactionsResponse struct {
	Transactions []*TransactionDetail `json:"transactions"`
}
