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

	httputil "github.com/primev/mev-commit/indexer/pkg/http"
	"github.com/primev/mev-commit/indexer/pkg/utils"
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

func FetchBeaconExecutionBlock(httpc *http.Client, beaconBase string, blockNum int64) (*ExecInfo, error) {
	url := fmt.Sprintf("%s/execution/block/%d", beaconBase, blockNum)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wrap struct {
		Data []map[string]any `json:"data"`
	}
	if err := httputil.FetchJSONWithRetry(ctx, httpc, url, &wrap, 3, 300*time.Millisecond); err != nil || len(wrap.Data) == 0 {
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
		return nil, fmt.Errorf("exec block missing posConsensus.slot for %d", blockNum)
	}
	return out, nil
}

// validator pubkey from proposer index
func FetchValidatorPubkey(httpc *http.Client, beaconBase string, proposerIndex int64) ([]byte, error) {
	url := fmt.Sprintf("%s/validator/%d", beaconBase, proposerIndex)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var resp struct {
		Data struct {
			Pubkey string `json:"pubkey"`
		} `json:"data"`
	}
	if err := httputil.FetchJSONWithRetry(ctx, httpc, url, &resp, 3, 300*time.Millisecond); err != nil {
		return nil, err
	}
	if strings.TrimSpace(resp.Data.Pubkey) == "" {
		return nil, fmt.Errorf("validator %d pubkey empty", proposerIndex)
	}
	return utils.HexToBytes(resp.Data.Pubkey), nil
}

// Add this new function to fetch blocks from Alchemy RPC
func FetchBlockFromRPC(httpc *http.Client, rpcURL string, blockNumber int64) (*ExecInfo, error) {
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

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
func FetchCombinedBlockData(httpc *http.Client, rpcURL string, beaconBase string, blockNumber int64) (*ExecInfo, error) {
	// Get execution block from Alchemy (always available)
	execBlock, err := FetchBlockFromRPC(httpc, rpcURL, blockNumber)
	if err != nil {
		return nil, err
	}

	// Convert block number to slot for beacon chain query
	slotNumber := utils.BlockNumberToSlot(blockNumber)

	// Try to get beacon chain data using slot number (may not exist for recent blocks)
	beaconData, _ := FetchBeaconExecutionBlock(httpc, beaconBase, slotNumber)

	// Merge data - use Alchemy as primary, beacon as supplement
	if beaconData != nil {
		execBlock.Slot = beaconData.Slot
		execBlock.ProposerIdx = beaconData.ProposerIdx
		execBlock.RelayTag = beaconData.RelayTag
		execBlock.RewardEth = beaconData.RewardEth
	} else {
		// Set the calculated slot if beacon data not available
		execBlock.Slot = slotNumber
	}

	return execBlock, nil
}
