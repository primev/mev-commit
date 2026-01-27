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

// TraceLog represents a log entry from simulation (used by both eth_simulateV1 and debug_traceCall)
type TraceLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}

// SimulateV1CallResult represents a single call result from eth_simulateV1
type SimulateV1CallResult struct {
	Status     hexutil.Uint64 `json:"status"`
	GasUsed    hexutil.Uint64 `json:"gasUsed"`
	ReturnData hexutil.Bytes  `json:"returnData"`
	Logs       []TraceLog     `json:"logs"`
	Error      *SimulateError `json:"error,omitempty"`
}

// SimulateError represents an error from eth_simulateV1
type SimulateError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// SimulateV1Block represents a block result from eth_simulateV1
type SimulateV1Block struct {
	Number  hexutil.Uint64         `json:"number"`
	GasUsed hexutil.Uint64         `json:"gasUsed"`
	Calls   []SimulateV1CallResult `json:"calls"`
}

// TraceCallResult represents the result of debug_traceCall with callTracer
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

// rpcEndpoint holds the RPC client
type rpcEndpoint struct {
	client *rpc.Client
}

// InlineSimulator simulates transactions using eth_simulateV1 (primary) with debug_traceCall fallback.
// eth_simulateV1 is lighter and preferred for performance, debug_traceCall is used when
// eth_simulateV1 is not supported or for edge cases requiring deeper tracing.
// Supports multiple RPC endpoints with fallback on connection errors.
type InlineSimulator struct {
	endpoints []rpcEndpoint
	metrics   *metrics
	logger    *slog.Logger
}

// NewInlineSimulator creates a new inline simulator with fallback support
// The first URL is the primary endpoint, subsequent URLs are fallbacks
func NewInlineSimulator(rpcURLs []string, logger *slog.Logger) (*InlineSimulator, error) {
	if len(rpcURLs) == 0 {
		return nil, errors.New("at least one RPC URL is required")
	}

	endpoints := make([]rpcEndpoint, 0, len(rpcURLs))
	for i, url := range rpcURLs {
		client, err := rpc.Dial(url)
		if err != nil {
			// Log warning but continue - we'll fail later if all endpoints are down
			if logger != nil {
				logger.Warn("failed to connect to RPC endpoint", "endpointIndex", i, "error", err)
			}
			continue
		}
		endpoints = append(endpoints, rpcEndpoint{client: client})
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("failed to connect to any RPC endpoint")
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

// Metrics returns prometheus collectors for the simulator
func (s *InlineSimulator) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.attempts,
		s.metrics.success,
		s.metrics.fail,
		s.metrics.latency,
	}
}

// Simulate executes a transaction simulation using eth_simulateV1 (primary) or debug_traceCall (fallback).
// eth_simulateV1 is lighter and preferred for performance.
// debug_traceCall is used when eth_simulateV1 is not supported by the RPC.
// Supported states: "latest" and "pending"
// If the primary endpoint fails with a connection error, fallback endpoints are tried.
func (s *InlineSimulator) Simulate(ctx context.Context, txRaw string, state SimState) ([]*types.Log, bool, error) {
	start := time.Now()
	defer func() {
		s.metrics.latency.Observe(float64(time.Since(start).Milliseconds()))
	}()

	s.metrics.attempts.Inc()

	// Decode the raw transaction
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

	// Get sender
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("failed to get sender: %w", err)
	}

	// Build call object (used by both eth_simulateV1 and debug_traceCall)
	callObj := map[string]interface{}{
		"from":  sender.Hex(),
		"gas":   hexutil.Uint64(tx.Gas()),
		"value": hexutil.EncodeBig(tx.Value()),
		"input": hexutil.Encode(tx.Data()), // eth_simulateV1 uses "input", debug_traceCall uses "data"
	}
	if tx.To() != nil {
		callObj["to"] = tx.To().Hex()
	}

	// Set gas price fields based on transaction type
	switch tx.Type() {
	case types.DynamicFeeTxType, types.BlobTxType:
		callObj["maxFeePerGas"] = hexutil.EncodeBig(tx.GasFeeCap())
		callObj["maxPriorityFeePerGas"] = hexutil.EncodeBig(tx.GasTipCap())
	default:
		callObj["gasPrice"] = hexutil.EncodeBig(tx.GasPrice())
	}

	// Try eth_simulateV1 first (lighter, better performance)
	logs, isSwap, err := s.simulateWithFallback(ctx, callObj, state)
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, err
	}

	s.metrics.success.Inc()
	return logs, isSwap, nil
}

// simulateWithFallback tries eth_simulateV1 first, falls back to debug_traceCall if not supported
func (s *InlineSimulator) simulateWithFallback(ctx context.Context, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	var lastErr error

	for i, endpoint := range s.endpoints {
		// Try eth_simulateV1 first
		logs, isSwap, err := s.executeSimulateV1(ctx, endpoint.client, callObj, state)
		if err == nil {
			if i > 0 {
				s.logger.Info("simulation succeeded on fallback endpoint", "endpointIndex", i, "method", "eth_simulateV1")
			}
			return logs, isSwap, nil
		}

		// Check if eth_simulateV1 is not supported - fall back to debug_traceCall
		if isMethodNotSupported(err) {
			s.logger.Debug("eth_simulateV1 not supported, falling back to debug_traceCall", "endpointIndex", i)
			logs, isSwap, err = s.executeDebugTraceCall(ctx, endpoint.client, callObj, state)
			if err == nil {
				if i > 0 {
					s.logger.Info("simulation succeeded on fallback endpoint", "endpointIndex", i, "method", "debug_traceCall")
				}
				return logs, isSwap, nil
			}
		}

		lastErr = err

		// Only fallback to next endpoint if it's not an application error
		if !shouldFallback(err) {
			return nil, false, err
		}

		s.logger.Warn("endpoint failed, trying fallback",
			"endpointIndex", i,
			"error", err,
			"remainingEndpoints", len(s.endpoints)-i-1,
		)
	}

	return nil, false, fmt.Errorf("all endpoints failed: %w", lastErr)
}

// executeSimulateV1 calls eth_simulateV1 (lighter, preferred method)
func (s *InlineSimulator) executeSimulateV1(ctx context.Context, client *rpc.Client, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	// Build eth_simulateV1 request
	// Format: https://ethereum.github.io/execution-apis/ethsimulatev1-notes/
	simRequest := map[string]interface{}{
		"blockStateCalls": []map[string]interface{}{
			{
				"calls": []map[string]interface{}{callObj},
			},
		},
		"validation": true,
	}

	var result []SimulateV1Block
	err := client.CallContext(ctx, &result, "eth_simulateV1", simRequest, string(state))
	if err != nil {
		return nil, false, err
	}

	// Validate response
	if len(result) == 0 {
		return nil, false, errors.New("empty response from eth_simulateV1")
	}
	block := result[0]
	if len(block.Calls) == 0 {
		return nil, false, errors.New("no calls in eth_simulateV1 response")
	}

	call := block.Calls[0]

	// Check for revert (status 0x0)
	if call.Status == 0 {
		reason := "execution reverted"
		if call.Error != nil && call.Error.Message != "" {
			reason = call.Error.Message
		} else if len(call.ReturnData) > 0 {
			reason = decodeRevert(hexutil.Encode(call.ReturnData), reason)
		}
		return nil, false, fmt.Errorf("reverted: %s", reason)
	}

	// Validate gas used
	if call.GasUsed == 0 {
		return nil, false, errors.New("empty response: missing or zero gas used")
	}

	// Detect swaps from logs
	isSwap, _ := DetectSwapsFromLogs(call.Logs)

	// Convert logs to types.Log
	logs := convertTraceLogs(call.Logs)

	return logs, isSwap, nil
}

// executeDebugTraceCall calls debug_traceCall (fallback for deeper tracing or unsupported eth_simulateV1)
func (s *InlineSimulator) executeDebugTraceCall(ctx context.Context, client *rpc.Client, callObj map[string]interface{}, state SimState) ([]*types.Log, bool, error) {
	// debug_traceCall uses "data" instead of "input"
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

	// Check for revert at top level
	if result.Error != "" {
		reason := decodeRevertFromTrace(result.Output, result.Error)
		return nil, false, fmt.Errorf("reverted: %s", reason)
	}

	// Check for inner call errors (recursive)
	if innerErr := findInnerCallError(&result); innerErr != "" {
		return nil, false, fmt.Errorf("inner call reverted: %s", innerErr)
	}

	// Validate trace response - a valid trace always has non-zero GasUsed
	gasUsed, err := hexutil.DecodeUint64(result.GasUsed)
	if err != nil || gasUsed == 0 {
		return nil, false, errors.New("empty trace response: missing or zero gas used")
	}

	// Collect all logs from trace (depth-first, execution order)
	var traceLogs []TraceLog
	collectTraceLogs(&result, &traceLogs)

	// Detect swaps from logs
	isSwap, _ := DetectSwapsFromLogs(traceLogs)

	// Convert logs to types.Log
	logs := convertTraceLogs(traceLogs)

	return logs, isSwap, nil
}

// isMethodNotSupported checks if the error indicates the RPC method is not supported
func isMethodNotSupported(err error) bool {
	if err == nil {
		return false
	}
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		// -32601 is the standard JSON-RPC error code for "Method not found"
		// -32600 is "Invalid Request"
		code := rpcErr.ErrorCode()
		if code == -32601 || code == -32600 {
			return true
		}
		// Also check error message for method not found
		msg := strings.ToLower(err.Error())
		return strings.Contains(msg, "method not found") ||
			strings.Contains(msg, "not supported") ||
			strings.Contains(msg, "unknown method")
	}
	return false
}

// shouldFallback returns true if the error should trigger a fallback to the next endpoint.
// JSON-RPC errors (invalid method, invalid params) should NOT trigger fallback.
// HTTP 4xx errors (except 429 rate limit) should NOT trigger fallback.
// Everything else (network errors, 5xx, 429) should fallback.
func shouldFallback(err error) bool {
	if err == nil {
		return false
	}

	// JSON-RPC errors are application-level - don't fallback
	// These include "method not found", "invalid params", etc.
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		return false
	}

	// HTTP errors: 4xx (except 429) are client errors - don't fallback
	// 5xx and 429 (rate limit) should fallback
	var httpErr rpc.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 && httpErr.StatusCode != 429 {
			return false
		}
		return true
	}

	// Everything else (network errors, timeouts, etc.) - fallback
	return true
}

// findInnerCallError recursively checks for errors in nested calls
// Returns the error reason with the failing call path (to, type)
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

// collectTraceLogs recursively collects logs from the trace result in depth-first order
// This matches the execution order of events
func collectTraceLogs(call *TraceCallResult, logs *[]TraceLog) {
	*logs = append(*logs, call.Logs...)
	for i := range call.Calls {
		collectTraceLogs(&call.Calls[i], logs)
	}
}

// convertTraceLogs converts TraceLog to types.Log
func convertTraceLogs(traceLogs []TraceLog) []*types.Log {
	logs := make([]*types.Log, 0, len(traceLogs))
	for i, tl := range traceLogs {
		log := &types.Log{
			Address: tl.Address,
			Topics:  tl.Topics,
			Data:    tl.Data,
			Index:   uint(i),
		}
		logs = append(logs, log)
	}
	return logs
}

// decodeRevertFromTrace attempts to decode a revert reason from trace output
func decodeRevertFromTrace(output string, fallback string) string {
	if output == "" || output == "0x" {
		return fallback
	}
	// Try to decode using the existing decodeRevert function
	if reason := decodeRevert(output, ""); reason != "" {
		return reason
	}
	return fallback
}

// Close closes all RPC clients and implements io.Closer
func (s *InlineSimulator) Close() error {
	for _, endpoint := range s.endpoints {
		if endpoint.client != nil {
			endpoint.client.Close()
		}
	}
	return nil
}
