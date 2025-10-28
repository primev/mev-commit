package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"

	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
)

type ExecInfo struct {
	BlockNumber int64
	Slot        int64
	ProposerIdx *int64
	Timestamp   *time.Time
	RelayTag    *string
	BuilderHex  *string
	FeeRecHex   *string
	RewardEth   *float64
}

// appendAPIKey appends the API key to the URL if provided
func appendAPIKey(url, apiKey string) string {
	if apiKey == "" {
		return url
	}
	if strings.Contains(url, "?") {
		return url + "&apikey=" + apiKey
	}
	return url + "?apikey=" + apiKey
}

func FetchBeaconExecutionBlock(ctx context.Context, httpc *retryablehttp.Client, limiter *rate.Limiter, beaconBase string, apiKey string, blockNum int64) (*ExecInfo, error) {
	t0 := time.Now()

	// Rate limit beacon API calls
	if limiter != nil {
		if err := limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}
	}
	rateLimitWait := time.Since(t0)

	url := appendAPIKey(fmt.Sprintf("%s/execution/block/%d", beaconBase, blockNum), apiKey)

	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	t1 := time.Now()
	var wrap struct {
		Data []map[string]any `json:"data"`
	}
	if err := httputil.FetchJSON(ctx, httpc, url, &wrap); err != nil || len(wrap.Data) == 0 {
		httpDuration := time.Since(t1)
		return nil, fmt.Errorf("no exec block %d (http_ms=%d, rate_wait_ms=%d): %w",
			blockNum, httpDuration.Milliseconds(), rateLimitWait.Milliseconds(), err)
	}
	httpDuration := time.Since(t1)
	j := wrap.Data[0]
	out := &ExecInfo{BlockNumber: blockNum}

	// posConsensus.slot & proposerIndex
	if pc, ok := j["posConsensus"].(map[string]any); ok {
		if v, ok := pc["slot"].(float64); ok {
			out.Slot = int64(v)
		} else if s, ok := pc["slot"].(string); ok {
			if n, err := strconv.ParseInt(s, 10, 64); err == nil {
				out.Slot = n
			}
		}
		if v, ok := pc["proposerIndex"].(float64); ok {
			x := int64(v)
			out.ProposerIdx = &x
		} else if s, ok := pc["proposerIndex"].(string); ok {
			if n, err := strconv.ParseInt(s, 10, 64); err == nil {
				out.ProposerIdx = &n
			}
		}
	}

	// timestamp
	if v, ok := j["timestamp"]; ok {
		switch t := v.(type) {
		case float64:
			u := time.Unix(int64(t), 0).UTC()
			out.Timestamp = &u
		case string:
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				u := time.Unix(n, 0).UTC()
				out.Timestamp = &u
			}
		}
	}

	// relay
	if rel, ok := j["relay"].(map[string]any); ok {
		if s, ok := rel["tag"].(string); ok {
			out.RelayTag = &s
		}
		if s, ok := rel["builderPubkey"].(string); ok {
			out.BuilderHex = &s
		}
		if s, ok := rel["producerFeeRecipient"].(string); ok {
			out.FeeRecHex = &s
		}
	}

	// reward eth from blockMevReward or producerReward
	if v, ok := j["blockMevReward"]; ok {
		switch t := v.(type) {
		case float64:
			f := t
			if f > 1e10 {
				f = f / 1e18 // wei -> ETH
			}
			out.RewardEth = &f
		case string:
			if strings.HasPrefix(t, "0x") {
				if bi, ok := new(big.Int).SetString(t[2:], 16); ok {
					f, _ := new(big.Rat).SetFrac(bi, big.NewInt(1e18)).Float64()
					out.RewardEth = &f
				}
			} else if f, err := strconv.ParseFloat(t, 64); err == nil {
				out.RewardEth = &f
			}
		}
	} else if v, ok := j["producerReward"]; ok {
		if f, ok := v.(float64); ok {
			out.RewardEth = &f
		}
	}

	// sanity
	if out.Slot == 0 {
		totalDuration := time.Since(t0)
		return nil, fmt.Errorf("exec block missing posConsensus.slot for %d (total_ms=%d)",
			blockNum, totalDuration.Milliseconds())
	}

	_ = httpDuration // Used for error messages
	return out, nil
}

// validator pubkey from proposer index
func FetchValidatorPubkey(ctx context.Context, httpc *retryablehttp.Client, limiter *rate.Limiter, beaconBase string, apiKey string, proposerIndex int64) ([]byte, error) {
	t0 := time.Now()

	// Rate limit beacon API calls
	if limiter != nil {
		if err := limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}
	}
	rateLimitWait := time.Since(t0)

	url := appendAPIKey(fmt.Sprintf("%s/validator/%d", beaconBase, proposerIndex), apiKey)

	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	t1 := time.Now()
	var resp struct {
		Data struct {
			Pubkey string `json:"pubkey"`
		} `json:"data"`
	}
	if err := httputil.FetchJSON(ctx, httpc, url, &resp); err != nil {
		httpDuration := time.Since(t1)
		return nil, fmt.Errorf("http error (http_ms=%d, rate_wait_ms=%d): %w",
			httpDuration.Milliseconds(), rateLimitWait.Milliseconds(), err)
	}
	httpDuration := time.Since(t1)

	if strings.TrimSpace(resp.Data.Pubkey) == "" {
		totalDuration := time.Since(t0)
		return nil, fmt.Errorf("validator %d pubkey empty (total_ms=%d)", proposerIndex, totalDuration.Milliseconds())
	}

	_ = httpDuration // Used for error messages
	_ = rateLimitWait // Used for error messages
	return common.FromHex(resp.Data.Pubkey), nil
}

// to fetch blocks from Alchemy RPC
func fetchBlockFromRPC(httpc *retryablehttp.Client, rpcURL string, blockNumber int64) (*ExecInfo, error) {
	underlyingClient := httpc.HTTPClient
	// Get block data from Alchemy
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBlockByNumber",
		"params":  []any{fmt.Sprintf("0x%x", blockNumber), true}, // true for full transaction objects
	}

	buf, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", rpcURL, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := underlyingClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Result struct {
			Number    string `json:"number"`
			Timestamp string `json:"timestamp"`
			Miner     string `json:"miner"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Result.Number == "" {
		return nil, fmt.Errorf("block not found")
	}

	// Convert hex timestamp to time
	timestampHex := result.Result.Timestamp[2:] // Remove 0x
	timestamp, _ := strconv.ParseInt(timestampHex, 16, 64)
	blockTime := time.Unix(timestamp, 0)

	return &ExecInfo{
		BlockNumber: blockNumber,
		Timestamp:   &blockTime,
	}, nil
}
func FetchCombinedBlockData(ctx context.Context, httpc *retryablehttp.Client, limiter *rate.Limiter, rpcURL string, beaconBase string, apiKey string, blockNumber int64) (*ExecInfo, error) {
	// Get execution block from Alchemy (always available)
	execBlock, err := fetchBlockFromRPC(httpc, rpcURL, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("fetch from RPC: %w", err)
	}

	beaconData, beaconErr := FetchBeaconExecutionBlock(ctx, httpc, limiter, beaconBase, apiKey, blockNumber)

	// Merge data - use Alchemy as primary, beacon as supplement
	if beaconData != nil {
		execBlock.Slot = beaconData.Slot
		execBlock.ProposerIdx = beaconData.ProposerIdx
		execBlock.RelayTag = beaconData.RelayTag
		execBlock.RewardEth = beaconData.RewardEth
		execBlock.BuilderHex = beaconData.BuilderHex
		execBlock.FeeRecHex = beaconData.FeeRecHex
	} else if beaconErr != nil {
		// Return error with context about beacon fetch failure
		return nil, fmt.Errorf("beaconcha.in API failed for block %d: %w", blockNumber, beaconErr)
	}
	return execBlock, nil
}
