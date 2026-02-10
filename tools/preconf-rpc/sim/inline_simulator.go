package sim

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

// TraceLog represents a log entry from simulation.
type TraceLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}

// SimulateV1CallResult represents a call result from eth_simulateV1.
type SimulateV1CallResult struct {
	Status     hexutil.Uint64 `json:"status"`
	GasUsed    hexutil.Uint64 `json:"gasUsed"`
	ReturnData hexutil.Bytes  `json:"returnData"`
	Logs       []TraceLog     `json:"logs"`
	Error      *SimulateError `json:"error,omitempty"`
}

// SimulateError represents an error returned by eth_simulateV1.
type SimulateError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// SimulateV1Block represents a block result from eth_simulateV1.
type SimulateV1Block struct {
	Number  hexutil.Uint64         `json:"number"`
	GasUsed hexutil.Uint64         `json:"gasUsed"`
	Calls   []SimulateV1CallResult `json:"calls"`
}

// TraceCallResult represents the result of debug_traceCall with callTracer.
type TraceCallResult struct {
	Type    string            `json:"type"`
	From    string            `json:"from"`
	To      string            `json:"to"`
	Value   string            `json:"value,omitempty"`
	Gas     string            `json:"gas"`
	GasUsed string            `json:"gasUsed"`
	Input   string            `json:"input"`
	Output  string            `json:"output"`
	Error   string            `json:"error,omitempty"`
	Calls   []TraceCallResult `json:"calls,omitempty"`
	Logs    []TraceLog        `json:"logs,omitempty"`
}

type rpcEndpoint struct {
	client *rpc.Client
}

// InlineSimulator simulates transactions using eth_simulateV1 with debug_traceCall as fallback.
// It prefers eth_simulateV1 for better performance, falling back to debug_traceCall when
// the RPC doesn't support eth_simulateV1. Multiple endpoints can be configured for redundancy.
type InlineSimulator struct {
	endpoints []rpcEndpoint
	metrics   *metrics
	logger    *slog.Logger
}

// NewInlineSimulator creates a simulator with the given RPC endpoints.
// The first URL is primary; others are used as fallbacks on network errors.
func NewInlineSimulator(rpcURLs []string, logger *slog.Logger) (*InlineSimulator, error) {
	if len(rpcURLs) == 0 {
		return nil, errors.New("at least one RPC URL is required")
	}

	endpoints := make([]rpcEndpoint, 0, len(rpcURLs))
	for i, url := range rpcURLs {
		client, err := rpc.Dial(url)
		if err != nil {
			if logger != nil {
				logger.Warn("failed to connect to RPC endpoint", "endpointIndex", i, "error", err)
			}
			continue
		}
		endpoints = append(endpoints, rpcEndpoint{client: client})
	}

	if len(endpoints) == 0 {
		return nil, errors.New("failed to connect to any RPC endpoint")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &InlineSimulator{
		endpoints: endpoints,
		metrics:   newMetrics(),
		logger:    logger,
	}, nil
}

// Metrics returns prometheus collectors for monitoring.
func (s *InlineSimulator) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.attempts,
		s.metrics.success,
		s.metrics.fail,
		s.metrics.latency,
	}
}

// Simulate runs a transaction simulation and returns logs, swap detection result, and any error.
// State can be "latest" or "pending".
func (s *InlineSimulator) Simulate(ctx context.Context, txRaw string, state SimState) ([]*types.Log, bool, error) {
	start := time.Now()
	defer func() {
		s.metrics.latency.Observe(float64(time.Since(start).Milliseconds()))
	}()

	s.metrics.attempts.Inc()

	rawBytes, err := hex.DecodeString(strings.TrimPrefix(txRaw, "0x"))
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("invalid hex: %w", err)
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(rawBytes); err != nil {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("invalid transaction: %w", err)
	}

	sender, err := recoverSender(tx)
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("failed to recover sender: %w", err)
	}

	// Build call object. We use "input" here; debug_traceCall expects "data" so we convert later.
	callObj := map[string]interface{}{
		"from":  sender.Hex(),
		"gas":   hexutil.Uint64(tx.Gas()),
		"value": hexutil.EncodeBig(tx.Value()),
		"input": hexutil.Encode(tx.Data()),
	}
	if tx.To() != nil {
		callObj["to"] = tx.To().Hex()
	}

	switch tx.Type() {
	case types.DynamicFeeTxType, types.BlobTxType:
		callObj["maxFeePerGas"] = hexutil.EncodeBig(tx.GasFeeCap())
		callObj["maxPriorityFeePerGas"] = hexutil.EncodeBig(tx.GasTipCap())
	default:
		callObj["gasPrice"] = hexutil.EncodeBig(tx.GasPrice())
	}

	logs, isSwap, err := s.simulateWithFallback(ctx, callObj, state)
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, err
	}

	s.metrics.success.Inc()
	return logs, isSwap, nil
}

// simulateWithFallback tries endpoints in order, using eth_simulateV1 first then debug_traceCall.
func (s *InlineSimulator) simulateWithFallback(ctx context.Context, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	var lastErr error

	for i, endpoint := range s.endpoints {
		logs, isSwap, err := s.executeSimulateV1(ctx, endpoint.client, callObj, state)
		if err == nil {
			if i > 0 {
				s.logger.Info("simulation succeeded on fallback endpoint", "endpointIndex", i, "method", "eth_simulateV1")
			}
			return logs, isSwap, nil
		}

		// If eth_simulateV1 isn't supported, try debug_traceCall on the same endpoint
		if isMethodNotSupported(err) {
			s.logger.Debug("eth_simulateV1 not supported, trying debug_traceCall", "endpointIndex", i)
			logs, isSwap, err = s.executeDebugTraceCall(ctx, endpoint.client, callObj, state)
			if err == nil {
				if i > 0 {
					s.logger.Info("simulation succeeded on fallback endpoint", "endpointIndex", i, "method", "debug_traceCall")
				}
				return logs, isSwap, nil
			}
		}

		lastErr = err

		// Don't retry on application errors (reverts, bad requests).
		// Only retry on transient errors (network issues, 5xx, rate limits).
		if !shouldFallback(err) {
			return nil, false, err
		}

		s.logger.Warn("endpoint failed, trying next",
			"endpointIndex", i,
			"error", err,
			"remainingEndpoints", len(s.endpoints)-i-1,
		)
	}

	return nil, false, fmt.Errorf("all endpoints failed: %w", lastErr)
}

// executeSimulateV1 runs simulation using eth_simulateV1.
// See: https://ethereum.github.io/execution-apis/ethsimulatev1-notes/
func (s *InlineSimulator) executeSimulateV1(ctx context.Context, client *rpc.Client, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	simRequest := map[string]interface{}{
		"blockStateCalls": []map[string]interface{}{
			{"calls": []map[string]interface{}{callObj}},
		},
		"validation": true,
	}

	var result []SimulateV1Block
	if err := client.CallContext(ctx, &result, "eth_simulateV1", simRequest, string(state)); err != nil {
		return nil, false, err
	}

	if len(result) == 0 {
		return nil, false, &NonRetryableError{Err: errors.New("empty response from eth_simulateV1")}
	}
	block := result[0]
	if len(block.Calls) == 0 {
		return nil, false, &NonRetryableError{Err: errors.New("no calls in eth_simulateV1 response")}
	}

	call := block.Calls[0]

	// Extract call target for error messages
	toAddr := "contract creation"
	if to, ok := callObj["to"].(string); ok && to != "" {
		toAddr = to
	}

	// status 0 means reverted
	if call.Status == 0 {
		reason := "execution reverted"
		if call.Error != nil && call.Error.Message != "" {
			reason = call.Error.Message
		} else if len(call.ReturnData) > 0 {
			reason = decodeRevert(hexutil.Encode(call.ReturnData), reason)
		}
		return nil, false, &NonRetryableError{Err: fmt.Errorf("reverted: %s (to=%s)", reason, toAddr)}
	}

	if call.GasUsed == 0 {
		return nil, false, &NonRetryableError{Err: errors.New("invalid response: zero gas used")}
	}

	isSwap, _ := DetectSwapsFromLogs(call.Logs)
	logs := convertTraceLogs(call.Logs)

	return logs, isSwap, nil
}

// executeDebugTraceCall runs simulation using debug_traceCall with callTracer.
func (s *InlineSimulator) executeDebugTraceCall(ctx context.Context, client *rpc.Client, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	// debug_traceCall expects "data" instead of "input"
	traceCallObj := make(map[string]interface{})
	for k, v := range callObj {
		if k == "input" {
			traceCallObj["data"] = v
		} else {
			traceCallObj[k] = v
		}
	}

	var result TraceCallResult
	err := client.CallContext(ctx, &result, "debug_traceCall",
		traceCallObj,
		string(state),
		map[string]interface{}{
			"tracer": "callTracer",
			"tracerConfig": map[string]interface{}{
				"withLog":          true,
				"enableReturnData": true,
			},
		},
	)
	if err != nil {
		return nil, false, fmt.Errorf("debug_traceCall failed (state=%s): %w", state, err)
	}

	if result.Error != "" {
		reason := decodeRevertFromTrace(result.Output, result.Error)
		toAddr := result.To
		if toAddr == "" {
			toAddr = "contract creation"
		}
		return nil, false, &NonRetryableError{Err: fmt.Errorf("reverted: %s (to=%s)", reason, toAddr)}
	}

	// Check nested calls for reverts (e.g., inner contract call failed)
	if innerErr := findInnerCallError(&result); innerErr != "" {
		return nil, false, &NonRetryableError{Err: fmt.Errorf("inner call reverted: %s", innerErr)}
	}

	gasUsed, err := hexutil.DecodeUint64(result.GasUsed)
	if err != nil || gasUsed == 0 {
		return nil, false, &NonRetryableError{Err: errors.New("invalid trace: zero gas used")}
	}

	var traceLogs []TraceLog
	collectTraceLogs(&result, &traceLogs)

	isSwap, _ := DetectSwapsFromLogs(traceLogs)
	logs := convertTraceLogs(traceLogs)

	return logs, isSwap, nil
}

// isMethodNotSupported checks if the error indicates the RPC method doesn't exist.
func isMethodNotSupported(err error) bool {
	if err == nil {
		return false
	}
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		code := rpcErr.ErrorCode()
		// -32601: Method not found, -32600: Invalid Request
		if code == -32601 || code == -32600 {
			return true
		}
		msg := strings.ToLower(err.Error())
		return strings.Contains(msg, "method not found") ||
			strings.Contains(msg, "not supported") ||
			strings.Contains(msg, "unknown method")
	}
	return false
}

// shouldFallback determines if we should try the next endpoint.
// Returns false for application errors (reverts, bad requests) since retrying won't help.
// Returns true for transient errors (network issues, 5xx, rate limits).
func shouldFallback(err error) bool {
	if err == nil {
		return false
	}

	var nonRetryable *NonRetryableError
	if errors.As(err, &nonRetryable) {
		return false
	}

	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		return false
	}

	var httpErr rpc.HTTPError
	if errors.As(err, &httpErr) {
		// 4xx (except 429) are client errors, don't retry
		if httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 && httpErr.StatusCode != 429 {
			return false
		}
		return true
	}

	return true
}

// findInnerCallError recursively searches for errors in nested calls.
func findInnerCallError(call *TraceCallResult) string {
	for i := range call.Calls {
		if call.Calls[i].Error != "" {
			reason := decodeRevertFromTrace(call.Calls[i].Output, call.Calls[i].Error)
			return fmt.Sprintf("%s (to=%s, type=%s)", reason, call.Calls[i].To, call.Calls[i].Type)
		}
		if innerErr := findInnerCallError(&call.Calls[i]); innerErr != "" {
			return innerErr
		}
	}
	return ""
}

// collectTraceLogs gathers logs from the trace in execution order (depth-first).
func collectTraceLogs(call *TraceCallResult, logs *[]TraceLog) {
	*logs = append(*logs, call.Logs...)
	for i := range call.Calls {
		collectTraceLogs(&call.Calls[i], logs)
	}
}

func convertTraceLogs(traceLogs []TraceLog) []*types.Log {
	logs := make([]*types.Log, 0, len(traceLogs))
	for i, tl := range traceLogs {
		logs = append(logs, &types.Log{
			Address: tl.Address,
			Topics:  tl.Topics,
			Data:    tl.Data,
			Index:   uint(i),
		})
	}
	return logs
}

func decodeRevertFromTrace(output string, fallback string) string {
	if output == "" || output == "0x" {
		return fallback
	}
	if reason := decodeRevert(output, ""); reason != "" {
		return reason
	}
	return fallback
}

// Close releases all RPC connections.
func (s *InlineSimulator) Close() error {
	for _, endpoint := range s.endpoints {
		if endpoint.client != nil {
			endpoint.client.Close()
		}
	}
	return nil
}

// recoverSender extracts the sender address from a signed transaction.
// Uses the appropriate signer based on transaction type to handle edge cases
// like pre-EIP-155 transactions that lack chain ID replay protection.
func recoverSender(tx *types.Transaction) (common.Address, error) {
	var signer types.Signer

	switch tx.Type() {
	case types.LegacyTxType:
		chainID := tx.ChainId()
		if chainID.Sign() == 0 {
			signer = types.HomesteadSigner{}
		} else {
			signer = types.NewEIP155Signer(chainID)
		}
	case types.AccessListTxType:
		signer = types.NewEIP2930Signer(tx.ChainId())
	case types.DynamicFeeTxType:
		signer = types.NewLondonSigner(tx.ChainId())
	case types.BlobTxType:
		signer = types.NewCancunSigner(tx.ChainId())
	default:
		signer = types.LatestSignerForChainID(tx.ChainId())
	}

	return types.Sender(signer, tx)
}
