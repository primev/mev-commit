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
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
)

type ExecInfo struct {
	BlockNumber       int64
	Slot              int64
	ProposerIdx       *int64
	Timestamp         *time.Time
	RelayTag          *string
	BuilderPublicKey  *string
	ProposerFeeRecHex *string
	MevRewardEth      *float64
	ProposerRewardEth *float64
	FeeRecipient      *string
}

func FetchBeaconExecutionBlock(ctx context.Context, httpc *retryablehttp.Client, beaconBase string, blockNum int64) (*ExecInfo, error) {
	url := fmt.Sprintf("%s/execution/block/%d", beaconBase, blockNum)

	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	var wrap struct {
		Data []map[string]any `json:"data"`
	}
	if err := httputil.FetchJSON(ctx, httpc, url, &wrap); err != nil || len(wrap.Data) == 0 {
		return nil, fmt.Errorf("no exec block %d", blockNum)
	}
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
			out.BuilderPublicKey = &s
		}
		if s, ok := rel["producerFeeRecipient"].(string); ok {
			out.ProposerFeeRecHex = &s
		}
	}

	// reward eth from blockMevReward
	if v, ok := j["blockMevReward"]; ok {
		switch t := v.(type) {
		case float64:
			f := t
			if f > 1e10 { // likely wei
				f = f / 1e18
			}
			out.MevRewardEth = &f
		case string:
			if strings.HasPrefix(t, "0x") {
				if bi, ok := new(big.Int).SetString(t[2:], 16); ok {
					f, _ := new(big.Rat).SetFrac(bi, big.NewInt(1e18)).Float64()
					out.MevRewardEth = &f
				}
			} else if f, err := strconv.ParseFloat(t, 64); err == nil {
				out.MevRewardEth = &f
			}
		}
	}

	// producerReward → out.ProposerRewardEth (ETH units)
	if v, ok := j["producerReward"]; ok {
		switch t := v.(type) {
		case float64:
			f := t
			if f > 1e10 {
				f = f / 1e18
			}
			out.ProposerRewardEth = &f
		case string:
			if strings.HasPrefix(t, "0x") {
				if bi, ok := new(big.Int).SetString(t[2:], 16); ok {
					f, _ := new(big.Rat).SetFrac(bi, big.NewInt(1e18)).Float64()
					out.ProposerRewardEth = &f
				}
			} else if f, err := strconv.ParseFloat(t, 64); err == nil {
				out.ProposerRewardEth = &f
			}
		}
	}

	if fr, ok := j["feeRecipient"].(string); ok && strings.TrimSpace(fr) != "" {
		out.FeeRecipient = &fr
	}
	// sanity
	if out.Slot == 0 {
		return nil, fmt.Errorf("exec block missing posConsensus.slot for %d", blockNum)
	}
	return out, nil
}

// validator pubkey from proposer index
func FetchValidatorPubkey(ctx context.Context, httpc *retryablehttp.Client, beaconBase string, proposerIndex int64) ([]byte, error) {
	url := fmt.Sprintf("%s/validator/%d", beaconBase, proposerIndex)

	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	var resp struct {
		Data struct {
			Pubkey string `json:"pubkey"`
		} `json:"data"`
	}
	if err := httputil.FetchJSON(ctx, httpc, url, &resp); err != nil {
		return nil, err
	}
	if strings.TrimSpace(resp.Data.Pubkey) == "" {
		return nil, fmt.Errorf("validator %d pubkey empty", proposerIndex)
	}
	return common.FromHex(resp.Data.Pubkey), nil
}

// to fetch blocks from Alchemy RPC
func fetchBlockFromRPC(httpc *retryablehttp.Client, rpcURL string, blockNumber int64) (*ExecInfo, error) {
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBlockByNumber",
		"params":  []any{fmt.Sprintf("0x%x", blockNumber), true},
	}
	buf, _ := json.Marshal(payload)

	req, _ := retryablehttp.NewRequest("POST", rpcURL, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rpc HTTP %d", resp.StatusCode)
	}

	var result struct {
		Result struct {
			Number    string `json:"number"`
			Timestamp string `json:"timestamp"`
			Miner     string `json:"miner"`
			Author    string `json:"author"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", result.Error.Code, result.Error.Message)
	}
	if result.Result.Number == "" {
		return nil, fmt.Errorf("block not found")
	}

	// timestamp hex -> time
	tsHex := strings.TrimPrefix(result.Result.Timestamp, "0x")
	secs, err := strconv.ParseInt(tsHex, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("bad timestamp: %w", err)
	}
	t := time.Unix(secs, 0).UTC()

	fr := result.Result.Miner
	if fr == "" {
		fr = result.Result.Author
	}

	out := &ExecInfo{
		BlockNumber: blockNumber,
		Timestamp:   &t,
	}
	if strings.TrimSpace(fr) != "" {
		out.FeeRecipient = &fr
	}
	return out, nil
}
func FetchCombinedBlockData(ctx context.Context, httpc *retryablehttp.Client, rpcURL, beaconBase string, blockNumber int64) (*ExecInfo, error) {
	execBlock, err := fetchBlockFromRPC(httpc, rpcURL, blockNumber)
	if err != nil {
		return nil, err
	}

	if beaconData, _ := FetchBeaconExecutionBlock(ctx, httpc, beaconBase, blockNumber); beaconData != nil {
		execBlock.Slot = beaconData.Slot
		execBlock.ProposerIdx = beaconData.ProposerIdx
		execBlock.RelayTag = beaconData.RelayTag
		execBlock.BuilderPublicKey = beaconData.BuilderPublicKey
		execBlock.ProposerFeeRecHex = beaconData.ProposerFeeRecHex
		execBlock.MevRewardEth = beaconData.MevRewardEth
		execBlock.ProposerRewardEth = beaconData.ProposerRewardEth

		// If RPC didn't provide fee recipient for any reason, fall back to beacon (if present)
		if execBlock.FeeRecipient == nil && beaconData.FeeRecipient != nil {
			execBlock.FeeRecipient = beaconData.FeeRecipient
		}
	}
	return execBlock, nil
}
