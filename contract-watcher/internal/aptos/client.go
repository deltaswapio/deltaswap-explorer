package aptos

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/deltaswapio/deltaswap-explorer/contract-watcher/internal/metrics"
	"github.com/go-resty/resty/v2"
	"go.uber.org/ratelimit"
)

var ErrTooManyRequests = fmt.Errorf("too many requests")

const clientName = "aptos"

// AptosSDK is a client for the Aptos API.
type AptosSDK struct {
	client  *resty.Client
	rl      ratelimit.Limiter
	metrics metrics.Metrics
}

type GetLatestBlock struct {
	BlockHeight string `json:"block_height"`
}

type Payload struct {
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []any    `json:"arguments"`
	Type          string   `json:"type"`
}

type Transaction struct {
	Version string  `json:"version"`
	Hash    string  `json:"hash"`
	Payload Payload `json:"payload,omitempty"`
}

type GetBlockResult struct {
	BlockHeight    string        `json:"block_height"`
	BlockHash      string        `json:"block_hash"`
	BlockTimestamp string        `json:"block_timestamp"`
	Transactions   []Transaction `json:"transactions"`
}

type GetTransactionResult struct {
	Version             string `json:"version"`
	Hash                string `json:"hash"`
	StateChangeHash     string `json:"state_change_hash"`
	EventRootHash       string `json:"event_root_hash"`
	StateCheckpointHash any    `json:"state_checkpoint_hash"`
	GasUsed             string `json:"gas_used"`
	Success             bool   `json:"success"`
	VMStatus            string `json:"vm_status"`
}

func (r *GetBlockResult) GetBlockTime() (*time.Time, error) {
	t, err := strconv.ParseUint(r.BlockTimestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	tm := time.UnixMicro(int64(t))
	return &tm, nil
}

// NewAptosSDK creates a new AptosSDK.
func NewAptosSDK(url string, rl ratelimit.Limiter, metrics metrics.Metrics) *AptosSDK {
	return &AptosSDK{
		rl:      rl,
		client:  resty.New().SetBaseURL(url),
		metrics: metrics,
	}
}

func (s *AptosSDK) GetLatestBlock(ctx context.Context) (uint64, error) {
	s.rl.Take()
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&GetLatestBlock{}).
		Get("v1")

	if err != nil {
		return 0, err
	}

	s.metrics.IncRpcRequest(clientName, "get-latest-block", resp.StatusCode())

	if resp.IsError() {
		return 0, fmt.Errorf("status code: %s. %s", resp.Status(), string(resp.Body()))
	}

	result := resp.Result().(*GetLatestBlock)
	if result == nil {
		return 0, fmt.Errorf("empty response")
	}
	if result.BlockHeight == "" {
		return 0, fmt.Errorf("empty block height")
	}
	return strconv.ParseUint(result.BlockHeight, 10, 64)
}

func (s *AptosSDK) GetBlock(ctx context.Context, block uint64) (*GetBlockResult, error) {
	s.rl.Take()
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&GetBlockResult{}).
		SetQueryParam("with_transactions", "true").
		Get(fmt.Sprintf("v1/blocks/by_height/%d", block))

	if err != nil {
		return nil, err
	}

	s.metrics.IncRpcRequest(clientName, "get-block", resp.StatusCode())

	if resp.IsError() {
		if resp.StatusCode() == http.StatusTooManyRequests {
			return nil, ErrTooManyRequests
		}
		return nil, fmt.Errorf("status code: %s. %s", resp.Status(), string(resp.Body()))
	}

	return resp.Result().(*GetBlockResult), nil
}

func (s *AptosSDK) GetTransaction(ctx context.Context, version string) (*GetTransactionResult, error) {
	s.rl.Take()
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&GetTransactionResult{}).
		SetQueryParam("with_transactions", "true").
		Get(fmt.Sprintf("v1/transactions/by_version/%s", version))

	if err != nil {
		return nil, err
	}
	s.metrics.IncRpcRequest(clientName, "get-transaction", resp.StatusCode())

	if resp.IsError() {
		if resp.StatusCode() == http.StatusTooManyRequests {
			return nil, ErrTooManyRequests
		}
		return nil, fmt.Errorf("status code: %s. %s", resp.Status(), string(resp.Body()))
	}

	return resp.Result().(*GetTransactionResult), nil
}
