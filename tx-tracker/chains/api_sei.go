package chains

import (
	"context"
	"time"

	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

type seiTx struct {
	TxHash string
	Sender string
}

func seiTxSearchExtractor(tx *cosmosTxSearchResponse, logs []cosmosLogWrapperResponse) (*seiTx, error) {
	var sender string
	for _, l := range logs {
		for _, e := range l.Events {
			if e.Type == "message" {
				for _, attr := range e.Attributes {
					if attr.Key == "sender" {
						sender = attr.Value
					}
				}
				break
			}
		}
	}
	return &seiTx{TxHash: tx.Result.Txs[0].Hash, Sender: sender}, nil
}

type apiSei struct {
	deltachainUrl         string
	deltachainRateLimiter *time.Ticker
	p2pNetwork            string
}

func fetchSeiDetail(ctx context.Context, baseUrl string, rateLimiter *time.Ticker, sequence, timestamp, srcChannel, dstChannel string) (*seiTx, error) {
	params := &cosmosTxSearchParams{Sequence: sequence, Timestamp: timestamp, SrcChannel: srcChannel, DstChannel: dstChannel}
	return fetchTxSearch[seiTx](ctx, baseUrl, rateLimiter, params, seiTxSearchExtractor)
}

func (a *apiSei) fetchSeiTx(
	ctx context.Context,
	rateLimiter *time.Ticker,
	baseUrl string,
	txHash string,
) (*TxDetail, error) {
	txHash = txHashLowerCaseWith0x(txHash)
	deltachainTx, err := fetchDeltachainDetail(ctx, a.deltachainUrl, a.deltachainRateLimiter, txHash)
	if err != nil {
		return nil, err
	}
	seiTx, err := fetchSeiDetail(ctx, baseUrl, rateLimiter, deltachainTx.sequence, deltachainTx.timestamp, deltachainTx.srcChannel, deltachainTx.dstChannel)
	if err != nil {
		return nil, err
	}
	return &TxDetail{
		NativeTxHash: txHash,
		From:         deltachainTx.receiver,
		Attribute: &AttributeTxDetail{
			Type: "deltachain-gateway",
			Value: &WorchainAttributeTxDetail{
				OriginChainID: vaa.ChainIDSei,
				OriginTxHash:  seiTx.TxHash,
				OriginAddress: seiTx.Sender,
			},
		},
	}, nil
}
