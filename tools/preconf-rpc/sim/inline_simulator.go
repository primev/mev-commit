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

// TraceLog represents a log entry from the trace
type TraceLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
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

// InlineSimulator simulates transactions using debug_traceCall
// Supports multiple RPC endpoints with fallback
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

// Simulate executes a transaction simulation using debug_traceCall
// Supported states: "latest" and "pending"
// Requires an RPC provider that supports debug_traceCall with pending state
// (e.g., Alchemy, QuickNode, Erigon). This matches rethsim behavior.
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

	// Build call object for debug_traceCall
	// Format: https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracecall
	callObj := map[string]interface{}{
		"from":  sender.Hex(),
		"gas":   hexutil.Uint64(tx.Gas()),
		"value": hexutil.EncodeBig(tx.Value()),
		"data":  hexutil.Encode(tx.Data()),
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

	// Execute trace with fallback support
	result, err := s.executeTraceWithFallback(ctx, callObj, state)
	if err != nil {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("debug_traceCall failed (state=%s): %w", state, err)
	}

	// Check for revert at top level
	if result.Error != "" {
		s.metrics.fail.Inc()
		reason := decodeRevertFromTrace(result.Output, result.Error)
		return nil, false, fmt.Errorf("reverted: %s", reason)
	}

	// Check for inner call errors (recursive)
	if innerErr := findInnerCallError(result); innerErr != "" {
		s.metrics.fail.Inc()
		return nil, false, fmt.Errorf("inner call reverted: %s", innerErr)
	}

	// Collect all logs from trace (depth-first, execution order)
	var traceLogs []TraceLog
	collectTraceLogs(result, &traceLogs)

	// Validate trace response - a valid trace always has non-zero GasUsed
	// (at minimum, intrinsic gas of 21000 is consumed)
	gasUsed, err := hexutil.DecodeUint64(result.GasUsed)
	if err != nil || gasUsed == 0 {
		s.metrics.fail.Inc()
		return nil, false, errors.New("empty trace response: missing or zero gas used")
	}

	// Detect swaps from logs (same approach as rethsim - topic scanning only)
	isSwap, _ := DetectSwapsFromLogs(traceLogs)

	// Convert trace logs to types.Log
	logs := convertTraceLogs(traceLogs)

	s.metrics.success.Inc()
	return logs, isSwap, nil
}

// executeTraceWithFallback tries the primary endpoint first, then fallbacks on connection errors
func (s *InlineSimulator) executeTraceWithFallback(ctx context.Context, callObj map[string]interface{}, state SimState) (*TraceCallResult, error) {
	var lastErr error

	for i, endpoint := range s.endpoints {
		result, err := s.executeTrace(ctx, endpoint.client, callObj, state)
		if err == nil {
			if i > 0 {
				s.logger.Info("simulation succeeded on fallback endpoint", "endpointIndex", i)
			}
			return result, nil
		}

		lastErr = err

		// Only fallback if it's not an application error
		if !shouldFallback(err) {
			return nil, err
		}

		s.logger.Warn("endpoint failed, trying fallback",
			"endpointIndex", i,
			"error", err,
			"remainingEndpoints", len(s.endpoints)-i-1,
		)
	}

	return nil, fmt.Errorf("all endpoints failed: %w", lastErr)
}

// executeTrace calls debug_traceCall with the given parameters
func (s *InlineSimulator) executeTrace(ctx context.Context, client *rpc.Client, callObj map[string]interface{}, state SimState) (*TraceCallResult, error) {
	var result TraceCallResult
	err := client.CallContext(ctx, &result, "debug_traceCall",
		callObj,
		string(state), // "latest" or "pending"
		map[string]interface{}{
			"tracer": "callTracer",
			"tracerConfig": map[string]interface{}{
				"withLog":          true,
				"enableReturnData": true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
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
