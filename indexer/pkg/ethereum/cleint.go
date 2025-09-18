package ethereum

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/primev/mev-commit/indexer/pkg/config"
	httputil "github.com/primev/mev-commit/indexer/pkg/http"
)

const optInABIJSON = `[
  {
    "inputs":[{"internalType":"bytes[]","name":"valBLSPubKeys","type":"bytes[]"}],
    "name":"areValidatorsOptedIn",
    "outputs":[{"components":[
        {"internalType":"bool","name":"isVanillaOptedIn","type":"bool"},
        {"internalType":"bool","name":"isAvsOptedIn","type":"bool"},
        {"internalType":"bool","name":"isMiddlewareOptedIn","type":"bool"}
      ],"internalType":"struct OptInStatus[]","name":"","type":"tuple[]"}],
    "stateMutability":"view","type":"function"
  }
]`

func BuildAreOptedInCallData(pubkey []byte) ([]byte, error) {
	ab, err := abi.JSON(strings.NewReader(optInABIJSON))
	if err != nil {
		return nil, err
	}
	// Solidity expects bytes[]; pass []byte{pubkey} length 1
	return ab.Pack("areValidatorsOptedIn", [][]byte{pubkey})
}

// JSON-RPC helper (Infura / Alchemy / any node)
func EthCallJSONRPC(httpc *http.Client, rpcURL string, to string, data []byte, blockNum int64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tag := "0x" + strconv.FormatInt(blockNum, 16)
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_call",
		"params": []any{
			map[string]any{
				"to":   to,
				"data": "0x" + hex.EncodeToString(data),
			},
			tag,
		},
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("eth_call marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("build eth_call request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("eth_call http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Error != nil {
		return nil, fmt.Errorf("eth_call rpc error %d: %s", out.Error.Code, out.Error.Message)
	}
	if out.Result == "" || out.Result == "0x" {
		return nil, fmt.Errorf("eth_call empty result (err=%v)", out.Error)
	}
	return hex.DecodeString(strings.TrimPrefix(out.Result, "0x"))
}

func CallAreOptedInAtBlock(httpc *http.Client, cfg *config.Config, blockNum int64, pubkey []byte) (bool, error) {
	if len(pubkey) == 0 {
		return false, fmt.Errorf("empty pubkey")
	}
	data, err := BuildAreOptedInCallData(pubkey)
	if err != nil {
		return false, err
	}

	var ret []byte
	if cfg.InfuraRPC != "" {
		// Preferred: direct JSON-RPC via Infura
		ret, err = EthCallJSONRPC(httpc, cfg.InfuraRPC, cfg.OptInContract, data, blockNum)
		if err != nil {
			return false, err
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Fallback: Etherscan proxy
		tag := "0x" + strconv.FormatInt(blockNum, 16)
		url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_call&to=%s&data=0x%s&tag=%s",
			cfg.OptInContract, hex.EncodeToString(data), tag)
		if cfg.EtherscanKey != "" {
			url += "&apikey=" + cfg.EtherscanKey
		}
		var resp struct {
			Result string `json:"result"`
		}
		if err := httputil.FetchJSONWithRetry(ctx, httpc, url, &resp, cfg.MaxRetries, cfg.BaseRetryDelay); err != nil {
			return false, err
		}
		if resp.Result == "" || resp.Result == "0x" {
			return false, fmt.Errorf("empty result")
		}
		ret, err = hex.DecodeString(strings.TrimPrefix(resp.Result, "0x"))
		if err != nil {
			return false, err
		}
	}

	// Decode OptInStatus[] where OptInStatus=(bool,bool,bool)
	ab, err := abi.JSON(strings.NewReader(optInABIJSON))
	if err != nil {
		return false, err
	}
	var out []struct {
		IsVanillaOptedIn    bool
		IsAvsOptedIn        bool
		IsMiddlewareOptedIn bool
	}
	if err := ab.UnpackIntoInterface(&out, "areValidatorsOptedIn", ret); err != nil {
		return false, err
	}
	if len(out) == 0 {
		return false, nil
	}
	o := out[0]
	return o.IsVanillaOptedIn || o.IsAvsOptedIn || o.IsMiddlewareOptedIn, nil
}

// GetLatestBlockNumber gets the latest block number from Ethereum RPC
func GetLatestBlockNumber(httpc *http.Client, rpcURL string) (int64, error) {
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_blockNumber",
		"params":  []any{},
	}

	buf, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", rpcURL, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Convert hex to int64
	blockNum, err := strconv.ParseInt(result.Result[2:], 16, 64)
	return blockNum, err
}
