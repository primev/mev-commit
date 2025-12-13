package sim

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type SimCall struct {
	Status     string       `json:"status"`
	GasUsed    string       `json:"gasUsed"`
	ReturnData string       `json:"returnData"`
	Logs       []*types.Log `json:"logs"`
	Calls      []SimCall    `json:"calls,omitempty"`
}

type SimBlock struct {
	Number        string    `json:"number"`
	Hash          string    `json:"hash"`
	BaseFeePerGas string    `json:"baseFeePerGas"`
	LogsBloom     string    `json:"logsBloom"`
	Transactions  []string  `json:"transactions"`
	Calls         []SimCall `json:"calls"`
	TraceErrors   []string  `json:"traceErrors"`
	IsSwap        bool      `json:"swapDetected"`
}

type simResp struct {
	Result      []SimBlock `json:"result"`
	TraceErrors []string   `json:"traceErrors"`
}

type Simulator struct {
	apiURL string
	client *http.Client
}

func NewSimulator(apiURL string) *Simulator {
	return &Simulator{
		apiURL: apiURL,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				MaxIdleConns:        256,
				MaxIdleConnsPerHost: 256,
				IdleConnTimeout:     90 * time.Second,
				ForceAttemptHTTP2:   true,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
			Timeout: 15 * time.Second,
		},
	}
}

type reqBody struct {
	TxRaw      string `json:"raw"`
	Block      string `json:"block,omitempty"`
	Validation bool   `json:"validation,omitempty"`
	TraceCalls bool   `json:"traceCalls,omitempty"`
}

func (s *Simulator) Simulate(ctx context.Context, txRaw string) ([]*types.Log, bool, error) {
	body := reqBody{
		TxRaw:      txRaw,
		Block:      "latest",
		Validation: true,
		TraceCalls: true,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, false, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/rethsim/simulate/raw", s.apiURL),
		strings.NewReader(string(bodyJSON)),
	)
	if err != nil {
		return nil, false, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("bad status %d: %s", resp.StatusCode, string(respBody))
	}

	return parseResponse(respBody)
}

func parseResponse(body []byte) ([]*types.Log, bool, error) {
	trim := strings.TrimSpace(string(body))
	if len(trim) == 0 {
		return nil, false, errors.New("empty response")
	}

	var blk SimBlock
	var traceErrors []string

	if strings.HasPrefix(trim, "[") {
		var arr []SimBlock
		if err := json.Unmarshal(body, &arr); err != nil {
			return nil, false, fmt.Errorf("decode array: %w", err)
		}
		if len(arr) == 0 {
			return nil, false, errors.New("no blocks in response")
		}
		blk = arr[0]
	} else {
		var w simResp
		if err := json.Unmarshal(body, &w); err == nil && len(w.Result) > 0 {
			blk = w.Result[0]
			traceErrors = w.TraceErrors
		} else {
			if err := json.Unmarshal(body, &blk); err != nil {
				return nil, false, fmt.Errorf("decode object: %w", err)
			}
		}
	}

	if len(blk.Calls) == 0 {
		return nil, false, errors.New("no calls in response")
	}
	root := blk.Calls[0]

	// Failure → build extended error
	if strings.EqualFold(root.Status, "0x0") {
		reason := decodeRevert(root.ReturnData, "execution reverted")
		return nil, false, fmt.Errorf("reverted: %s", reason)
	}

	// Check trace errors for internal reverts
	if len(traceErrors) == 0 {
		traceErrors = blk.TraceErrors
	}
	for _, te := range traceErrors {
		if strings.Contains(strings.ToLower(te), "execution reverted") {
			return nil, false, errors.New(te)
		}
	}

	// Success → collect all logs (depth-first, execution order)
	var out []*types.Log
	collectLogs(&root, &out)
	return out, blk.IsSwap, nil
}

func collectLogs(n *SimCall, acc *[]*types.Log) {
	*acc = append(*acc, n.Logs...)
	for i := range n.Calls {
		collectLogs(&n.Calls[i], acc)
	}
}

func decodeRevert(dataHex string, fallback string) string {
	h := extract0x(dataHex)
	if h == "" || h == "0x" {
		return fallback
	}
	if len(h) < 10 {
		return fallback
	}
	sel := strings.ToLower(h[:10]) // 0x + 8 hex
	args := "0x" + h[10:]

	switch sel {
	case "0x08c379a0": // Error(string)
		if s, ok := abiDecodeString(args); ok {
			return s
		}
		return "execution reverted (Error)"
	case "0x4e487b71": // Panic(uint256)
		if n, ok := abiDecodeUint256(args); ok {
			return fmt.Sprintf("panic(0x%x)", n)
		}
		return "execution reverted (Panic)"
	default:
		return fmt.Sprintf("reverted (selector %s)", sel)
	}
}

func extract0x(x string) string {
	x = strings.TrimSpace(x)
	if strings.HasPrefix(x, "0x") {
		return x
	}
	// sometimes servers embed the hex in a quoted JSON string
	var s string
	if err := json.Unmarshal([]byte(x), &s); err == nil && strings.HasPrefix(s, "0x") {
		return s
	}
	return ""
}

func abiDecodeString(args string) (string, bool) {
	b, err := hex.DecodeString(strings.TrimPrefix(args, "0x"))
	if err != nil || len(b) < 64 {
		return "", false
	}
	// length at bytes 32..64
	l := new(big.Int).SetBytes(b[32:64]).Int64()
	if l < 0 || int(64+l) > len(b) {
		return "", false
	}
	return string(b[64 : 64+l]), true
}

func abiDecodeUint256(args string) (*big.Int, bool) {
	b, err := hex.DecodeString(strings.TrimPrefix(args, "0x"))
	if err != nil || len(b) < 32 {
		return nil, false
	}
	return new(big.Int).SetBytes(b[len(b)-32:]), true
}
