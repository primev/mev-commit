package sim_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/sim"
)

// Mock eth_simulateV1 response for a successful simple transfer
var simulateV1ResponseSimple = `[{
	"number": "0x1",
	"gasUsed": "0x5208",
	"calls": [{
		"status": "0x1",
		"gasUsed": "0x5208",
		"returnData": "0x",
		"logs": []
	}]
}]`

// Mock eth_simulateV1 response for a swap transaction with SushiSwap/Uniswap V2 Swap event
var simulateV1ResponseSwap = `[{
	"number": "0x1",
	"gasUsed": "0x20000",
	"calls": [{
		"status": "0x1",
		"gasUsed": "0x20000",
		"returnData": "0x",
		"logs": [{
			"address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"topics": [
				"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822",
				"0x0000000000000000000000007a250d5630b4cf539739df2c5dacb4c659f2488d",
				"0x000000000000000000000000ae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e"
			],
			"data": "0x0000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000000"
		}]
	}]
}]`

// Mock eth_simulateV1 response for a reverted transaction
var simulateV1ResponseRevert = `[{
	"number": "0x1",
	"gasUsed": "0x10000",
	"calls": [{
		"status": "0x0",
		"gasUsed": "0x10000",
		"returnData": "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001a496e73756666696369656e742062616c616e636500000000000000000000000000",
		"logs": [],
		"error": {
			"code": 3,
			"message": "execution reverted"
		}
	}]
}]`

// Mock debug_traceCall response for a successful simple transfer
var traceCallResponseSimple = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x1234567890123456789012345678901234567890",
	"value": "0xde0b6b3a7640000",
	"gas": "0x5208",
	"gasUsed": "0x5208",
	"input": "0x",
	"output": "0x",
	"logs": []
}`

// Mock debug_traceCall response for a swap transaction with SushiSwap/Uniswap V2 Swap event
var traceCallResponseSwap = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
	"value": "0x0",
	"gas": "0x30000",
	"gasUsed": "0x20000",
	"input": "0x38ed1739",
	"output": "0x",
	"logs": [
		{
			"address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"topics": [
				"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822",
				"0x0000000000000000000000007a250d5630b4cf539739df2c5dacb4c659f2488d",
				"0x000000000000000000000000ae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e"
			],
			"data": "0x0000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000000"
		}
	],
	"calls": []
}`

// Mock debug_traceCall response for a reverted transaction
var traceCallResponseRevert = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x1234567890123456789012345678901234567890",
	"value": "0x0",
	"gas": "0x30000",
	"gasUsed": "0x10000",
	"input": "0x",
	"output": "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001a496e73756666696369656e742062616c616e636500000000000000000000000000",
	"error": "execution reverted",
	"logs": []
}`

// Mock debug_traceCall response with nested calls containing Uniswap V3 swap
var traceCallResponseNestedSwap = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
	"value": "0x0",
	"gas": "0x50000",
	"gasUsed": "0x40000",
	"input": "0x",
	"output": "0x",
	"logs": [],
	"calls": [
		{
			"type": "CALL",
			"from": "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
			"to": "0xe592427a0aece92de3edee1f18e0157c05861564",
			"value": "0x0",
			"gas": "0x40000",
			"gasUsed": "0x30000",
			"input": "0x",
			"output": "0x",
			"logs": [
				{
					"address": "0xe592427a0aece92de3edee1f18e0157c05861564",
					"topics": [
						"0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67",
						"0x000000000000000000000000ae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
						"0x000000000000000000000000ae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e"
					],
					"data": "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"
				}
			]
		}
	]
}`

// Mock debug_traceCall response for multi-hop aggregator swap
var traceCallResponseMultiHop = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x1111111254EEB25477B68fb85Ed929f73A960582",
	"value": "0x0",
	"gas": "0x80000",
	"gasUsed": "0x60000",
	"input": "0x",
	"output": "0x",
	"logs": [],
	"calls": [
		{
			"type": "CALL",
			"from": "0x1111111254EEB25477B68fb85Ed929f73A960582",
			"to": "0xsomepool1",
			"value": "0x0",
			"gas": "0x30000",
			"gasUsed": "0x20000",
			"input": "0x",
			"output": "0x",
			"logs": [
				{
					"address": "0xsomepool1",
					"topics": [
						"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"
					],
					"data": "0x"
				}
			]
		},
		{
			"type": "CALL",
			"from": "0x1111111254EEB25477B68fb85Ed929f73A960582",
			"to": "0xsomepool2",
			"value": "0x0",
			"gas": "0x30000",
			"gasUsed": "0x20000",
			"input": "0x",
			"output": "0x",
			"logs": [
				{
					"address": "0xsomepool2",
					"topics": [
						"0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"
					],
					"data": "0x"
				}
			]
		}
	]
}`

// Mock debug_traceCall response for Curve StableSwap NG swap
var traceCallResponseCurve = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x99a58482BD75cbab83b27EC03CA68fF489b5788f",
	"value": "0x0",
	"gas": "0x50000",
	"gasUsed": "0x40000",
	"input": "0x",
	"output": "0x",
	"logs": [
		{
			"address": "0xbebc44782c7db0a1a60cb6fe97d0b483032ff1c7",
			"topics": [
				"0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140"
			],
			"data": "0x"
		}
	]
}`

// Mock debug_traceCall response for Balancer swap
var traceCallResponseBalancer = `{
	"type": "CALL",
	"from": "0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e",
	"to": "0x9008D19f58AAbD9eD0D60971565AA8510560ab41",
	"value": "0x0",
	"gas": "0x100000",
	"gasUsed": "0x80000",
	"input": "0x",
	"output": "0x",
	"logs": [],
	"calls": [
		{
			"type": "CALL",
			"from": "0x9008D19f58AAbD9eD0D60971565AA8510560ab41",
			"to": "0xsomepool",
			"value": "0x0",
			"gas": "0x50000",
			"gasUsed": "0x30000",
			"input": "0x",
			"output": "0x",
			"logs": [
				{
					"address": "0xsomepool",
					"topics": [
						"0x2170c741c41531aec20e7c107c24eecfdd15e69c9bb0a8dd37b1840b9e0b207b"
					],
					"data": "0x"
				}
			]
		}
	]
}`

// Real mainnet transaction test vectors for e2e validation of parsing and sender recovery.
var realTxVectors = []struct {
	name           string
	rawHex         string
	expectedType   uint8
	expectedSender string
	hasChainID     bool
}{
	{
		// Pre-EIP-155 legacy transaction (no chain ID replay protection)
		// Block 46147 - early mainnet transaction
		name:           "PreEIP155_Legacy",
		rawHex:         "f86780862d79883d2000825208945df9b87991262f6ba471f09758cde1c0fc1de734827a69801ca088ff6cf0fefd94db46111149ae4bfc179e9b94721fffd821d38d16464b3f71d0a045e0aff800961cfce805daef7016f9ae479c0a24afba38dd33c2ecdbb01dcacf",
		expectedType:   types.LegacyTxType,
		expectedSender: "0xD3678D173368032b34E00AE057C31b083FBAb830",
		hasChainID:     false,
	},
	{
		// EIP-1559 dynamic fee transaction (type 2)
		name:           "EIP1559_DynamicFee",
		rawHex:         "02f8730101843b9aca00850c92a69c0082520894d8da6bf26964af9d7eed9e03e53415d37aa9604588016345785d8a000080c001a0a9f0aabbfa2b831dd37d0f8d48d941f35f4fd40a1f2e2fa74a7df3e60aa534c8a0488e799fae157d086b8e0b624ab63627f14509482fe037e88f516a3725070896",
		expectedType:   types.DynamicFeeTxType,
		expectedSender: "0xcEC000D467698070C6D8D73D8ff1F60FD7DCb531",
		hasChainID:     true,
	},
}

func TestTransactionParsingAndSenderRecovery(t *testing.T) {
	for _, tc := range realTxVectors {
		t.Run(tc.name, func(t *testing.T) {
			rawBytes, err := hex.DecodeString(tc.rawHex)
			if err != nil {
				t.Fatalf("failed to decode hex: %v", err)
			}

			tx := new(types.Transaction)
			if err := tx.UnmarshalBinary(rawBytes); err != nil {
				t.Fatalf("failed to parse tx: %v", err)
			}

			if tx.Type() != tc.expectedType {
				t.Errorf("expected tx type %d, got %d", tc.expectedType, tx.Type())
			}

			if tc.hasChainID {
				if tx.ChainId().Sign() == 0 {
					t.Error("expected non-zero chainId")
				}
			} else {
				if tx.ChainId().Sign() != 0 {
					t.Errorf("expected chainId 0, got %s", tx.ChainId().String())
				}
			}

			sender, err := recoverSenderForTest(tx)
			if err != nil {
				t.Fatalf("failed to recover sender: %v", err)
			}

			if sender == (common.Address{}) {
				t.Error("recovered zero address")
			}

			if tc.expectedSender != "" {
				expected := common.HexToAddress(tc.expectedSender)
				if sender != expected {
					t.Errorf("sender mismatch: got %s, want %s", sender.Hex(), expected.Hex())
				}
			}

			t.Logf("tx=%s sender=%s chainId=%s", tc.name, sender.Hex(), tx.ChainId().String())
		})
	}
}

func recoverSenderForTest(tx *types.Transaction) (common.Address, error) {
	var signer types.Signer
	switch tx.Type() {
	case types.LegacyTxType:
		if tx.ChainId().Sign() == 0 {
			signer = types.HomesteadSigner{}
		} else {
			signer = types.NewEIP155Signer(tx.ChainId())
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

// TestSimulateWithRealTransactions tests the full simulation pipeline with real transaction hex.
// This verifies parsing, sender recovery, call object construction, and RPC interaction.
func TestSimulateWithRealTransactions(t *testing.T) {
	// Create mock server that accepts any valid simulation request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
			ID     int               `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		w.Header().Set("Content-Type", "application/json")

		// Return successful simulation response for eth_simulateV1
		if req.Method == "eth_simulateV1" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result":  json.RawMessage(simulateV1ResponseSimple),
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// Fallback to debug_traceCall
		if req.Method == "debug_traceCall" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result":  json.RawMessage(traceCallResponseSimple),
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// Unknown method
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	simulator, err := sim.NewInlineSimulator([]string{srv.URL}, nil)
	if err != nil {
		t.Fatalf("failed to create simulator: %v", err)
	}
	defer func() { _ = simulator.Close() }()

	for _, tc := range realTxVectors {
		t.Run(tc.name, func(t *testing.T) {
			// Test full simulation pipeline with real tx hex
			logs, isSwap, err := simulator.Simulate(context.Background(), tc.rawHex, sim.Latest)
			if err != nil {
				t.Fatalf("simulation failed: %v", err)
			}

			// Verify we got a response (mock returns simple transfer with no logs)
			if logs == nil {
				t.Error("expected non-nil logs slice")
			}

			t.Logf("tx=%s simulated successfully, logs=%d, isSwap=%v", tc.name, len(logs), isSwap)
		})
	}

	// Test with pending state
	t.Run("PendingState", func(t *testing.T) {
		logs, isSwap, err := simulator.Simulate(context.Background(), realTxVectors[0].rawHex, sim.Pending)
		if err != nil {
			t.Fatalf("simulation with pending state failed: %v", err)
		}
		t.Logf("pending state simulation: logs=%d, isSwap=%v", len(logs), isSwap)
	})
}

func TestInlineSimulator(t *testing.T) {
	// eth_simulateV1 responses
	simV1Responses := map[string]string{
		"simple": simulateV1ResponseSimple,
		"swap":   simulateV1ResponseSwap,
		"revert": simulateV1ResponseRevert,
	}

	// debug_traceCall responses (used as fallback)
	traceResponses := map[string]string{
		"simple":     traceCallResponseSimple,
		"swap":       traceCallResponseSwap,
		"revert":     traceCallResponseRevert,
		"nestedSwap": traceCallResponseNestedSwap,
		"multiHop":   traceCallResponseMultiHop,
		"curve":      traceCallResponseCurve,
		"balancer":   traceCallResponseBalancer,
	}

	// Helper to create test server with configurable eth_simulateV1 support
	createTestServer := func(supportSimulateV1 bool) *httptest.Server {
		return httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					Method string            `json:"method"`
					Params []json.RawMessage `json:"params"`
					ID     int               `json:"id"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				defer func() { _ = r.Body.Close() }()

				w.Header().Set("Content-Type", "application/json")

				if req.Method == "eth_simulateV1" {
					if !supportSimulateV1 {
						// Return JSON-RPC error for method not found
						response := map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      req.ID,
							"error": map[string]interface{}{
								"code":    -32601,
								"message": "Method not found",
							},
						}
						_ = json.NewEncoder(w).Encode(response)
						return
					}

					// Parse the simulateV1 request to get the call
					var simReq map[string]interface{}
					if err := json.Unmarshal(req.Params[0], &simReq); err != nil {
						http.Error(w, "bad params", http.StatusBadRequest)
						return
					}

					// Get the call object
					blockStateCalls, _ := simReq["blockStateCalls"].([]interface{})
					if len(blockStateCalls) == 0 {
						http.Error(w, "no block state calls", http.StatusBadRequest)
						return
					}
					blockState, _ := blockStateCalls[0].(map[string]interface{})
					calls, _ := blockState["calls"].([]interface{})
					if len(calls) == 0 {
						http.Error(w, "no calls", http.StatusBadRequest)
						return
					}
					callObj, _ := calls[0].(map[string]interface{})
					to, _ := callObj["to"].(string)

					var responseKey string
					switch strings.ToLower(to) {
					case "0x1234567890123456789012345678901234567890":
						if input, ok := callObj["input"].(string); ok && input == "0xrevert" {
							responseKey = "revert"
						} else {
							responseKey = "simple"
						}
					case "0x7a250d5630b4cf539739df2c5dacb4c659f2488d":
						responseKey = "swap"
					default:
						responseKey = "simple"
					}

					response := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result":  json.RawMessage(simV1Responses[responseKey]),
					}
					_ = json.NewEncoder(w).Encode(response)
					return
				}

				if req.Method == "debug_traceCall" {
					// Parse the call object to get the "to" address for routing
					var callObj map[string]interface{}
					if err := json.Unmarshal(req.Params[0], &callObj); err != nil {
						http.Error(w, "bad params", http.StatusBadRequest)
						return
					}

					// Route based on the "to" address
					to, _ := callObj["to"].(string)
					var responseKey string
					switch strings.ToLower(to) {
					case "0x1234567890123456789012345678901234567890":
						if data, ok := callObj["data"].(string); ok && data == "0xrevert" {
							responseKey = "revert"
						} else {
							responseKey = "simple"
						}
					case "0x7a250d5630b4cf539739df2c5dacb4c659f2488d":
						responseKey = "swap"
					case "0x68b3465833fb72a70ecdf485e0e4c7bd8665fc45":
						responseKey = "nestedSwap"
					case "0x1111111254eeb25477b68fb85ed929f73a960582":
						responseKey = "multiHop"
					case "0x99a58482bd75cbab83b27ec03ca68ff489b5788f":
						responseKey = "curve"
					case "0x9008d19f58aabd9ed0d60971565aa8510560ab41":
						responseKey = "balancer"
					default:
						responseKey = "simple"
					}

					response := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result":  json.RawMessage(traceResponses[responseKey]),
					}
					_ = json.NewEncoder(w).Encode(response)
					return
				}

				// Unknown method
				response := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req.ID,
					"error": map[string]interface{}{
						"code":    -32601,
						"message": "Method not found",
					},
				}
				_ = json.NewEncoder(w).Encode(response)
			}),
		)
	}

	// Test with eth_simulateV1 support
	t.Run("WithSimulateV1Support", func(t *testing.T) {
		srv := createTestServer(true)
		defer srv.Close()

		simulator, err := sim.NewInlineSimulator([]string{srv.URL}, nil)
		if err != nil {
			t.Fatalf("failed to create inline simulator: %v", err)
		}
		defer func() { _ = simulator.Close() }()

		t.Run("InvalidTransaction", func(t *testing.T) {
			_, _, err := simulator.Simulate(context.Background(), "invalid", sim.Latest)
			if err == nil {
				t.Error("expected error for invalid transaction")
			}
		})

		t.Run("InvalidHex", func(t *testing.T) {
			_, _, err := simulator.Simulate(context.Background(), "0xZZZZ", sim.Latest)
			if err == nil {
				t.Error("expected error for invalid hex")
			}
		})
	})

	// Test fallback to debug_traceCall when eth_simulateV1 is not supported
	t.Run("FallbackToDebugTraceCall", func(t *testing.T) {
		srv := createTestServer(false)
		defer srv.Close()

		simulator, err := sim.NewInlineSimulator([]string{srv.URL}, nil)
		if err != nil {
			t.Fatalf("failed to create inline simulator: %v", err)
		}
		defer func() { _ = simulator.Close() }()

		t.Run("InvalidTransaction", func(t *testing.T) {
			_, _, err := simulator.Simulate(context.Background(), "invalid", sim.Latest)
			if err == nil {
				t.Error("expected error for invalid transaction")
			}
		})

		t.Run("InvalidHex", func(t *testing.T) {
			_, _, err := simulator.Simulate(context.Background(), "0xZZZZ", sim.Latest)
			if err == nil {
				t.Error("expected error for invalid hex")
			}
		})
	})
}

// TestSwapDetection tests the swap detector with realistic trace responses
func TestSwapDetection(t *testing.T) {
	t.Run("NestedTraceLogCollection", func(t *testing.T) {
		logs := []sim.TraceLog{
			// First hop - SushiSwap (uses same signature as Uniswap V2 Swap)
			{
				Topics: []common.Hash{
					common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
				},
			},
			// Second hop - Uniswap V3
			{
				Topics: []common.Hash{
					common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"),
				},
			},
		}

		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if !isSwap {
			t.Error("expected swap detection for multi-hop aggregator trade")
		}
		if len(kinds) != 2 {
			t.Errorf("expected 2 swap kinds for multi-hop, got %v", kinds)
		}
	})

	// Test that we can detect swaps even with Transfer events mixed in
	t.Run("SwapWithTransferEvents", func(t *testing.T) {
		logs := []sim.TraceLog{
			// Transfer event (should be ignored)
			{
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				},
			},
			// Approval event (should be ignored)
			{
				Topics: []common.Hash{
					common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"),
				},
			},
			// Actual swap event (SushiSwap/Uniswap V2 Swap)
			{
				Topics: []common.Hash{
					common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
				},
			},
		}

		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if !isSwap {
			t.Error("expected swap detection even with Transfer/Approval events")
		}
		if len(kinds) != 1 || kinds[0] != "sushiswap_swap" {
			t.Errorf("expected sushiswap_swap, got %v", kinds)
		}
	})
}

func TestSwapSignatures(t *testing.T) {
	swapTests := []struct {
		name         string
		topicHash    string
		expectedKind string
	}{
		// Uniswap V2 Sync event (emitted on every swap)
		{"UniswapV2Sync", "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1", "uniswap_v2_swap"},
		// Uniswap V3 Swap
		{"UniswapV3Swap", "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67", "uniswap_v3_swap"},
		// Uniswap V4 Swap
		{"UniswapV4Swap", "0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f", "uniswap_v4_swap"},
		// MetaMask Swap Router
		{"MetaMaskSwapRouter", "0xbeee1e6e7fe307ddcf84b0a16137a4430ad5e2480fc4f4a8e250ab56ccd7630d", "metamask_swap_router"},
		// Fluid DEX
		{"FluidSwap", "0xfbce846c23a724e6e61161894819ec46c90a8d3dd96e90e7342c6ef49ffb539c", "fluid_swap"},
		// Curve TokenExchange
		{"CurveFinanceSwap", "0x56d0661e240dfb199ef196e16e6f42473990366314f0226ac978f7be3cd9ee83", "curve_finance_swap"},
		// Curve tricrypto
		{"CurveTricryptoSwap", "0x143f1f8e861fbdeddd5b46e844b7d3ac7b86a122f36e8c463859ee6811b1f29c", "curve_tricrypto_swap"},
		// Curve StableSwap NG
		{"CurveStableswapNGSwap", "0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140", "curve_stableswap_ng_swap"},
		// Balancer V2 Swap
		{"BalancerSwap", "0x2170c741c41531aec20e7c107c24eecfdd15e69c9bb0a8dd37b1840b9e0b207b", "balancer_swap"},
		// Balancer LOG_SWAP
		{"BalancerLogSwap", "0x908fb5ee8f16c6bc9bc3690973819f32a4d4b10188134543c88706e0e1d43378", "balancer_log_swap"},
		// 1inch Aggregation Router V6
		{"OneInchAggregationRouterV6", "0xfec331350fce78ba658e082a71da20ac9f8d798a99b3c79681c8440cbfe77e07", "oneinch_aggregation_router_v6"},
		// SushiSwap Swap (same signature as Uniswap V2 Swap event)
		{"SushiSwapSwap", "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822", "sushiswap_swap"},
		// KyberSwap
		{"KyberSwapSwap", "0xd6d4f5681c246c9f42c203e287975af1601f8df8035a9251f79aab5c8f09e2f8", "kyberswap_swap"},
		// PancakeSwap
		{"PancakeSwapSwap", "0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83", "pancakeswap_swap"},
		// DODO
		{"DODOSwap", "0xc2c0245e056d5fb095f04cd6373bc770802ebd1e6c918eb78fdef843cdb37b0f", "dodoswap_swap"},
	}

	for _, tt := range swapTests {
		t.Run("Detect_"+tt.name, func(t *testing.T) {
			logs := []sim.TraceLog{
				{
					Topics: []common.Hash{
						common.HexToHash(tt.topicHash),
					},
				},
			}
			isSwap, kinds := sim.DetectSwapsFromLogs(logs)
			if !isSwap {
				t.Errorf("expected swap detection for %s event", tt.name)
			}
			if len(kinds) != 1 || kinds[0] != tt.expectedKind {
				t.Errorf("expected %s swap kind, got %v", tt.expectedKind, kinds)
			}
		})
	}

	t.Run("DetectMultipleSwaps", func(t *testing.T) {
		logs := []sim.TraceLog{
			{
				Topics: []common.Hash{
					common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"), // Uniswap V2 Sync
				},
			},
			{
				Topics: []common.Hash{
					common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"), // Uniswap V3
				},
			},
			{
				Topics: []common.Hash{
					common.HexToHash("0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140"), // Curve StableSwap NG
				},
			},
		}
		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if !isSwap {
			t.Error("expected swap detection for multiple swap events")
		}
		if len(kinds) != 3 {
			t.Errorf("expected 3 swap kinds, got %v", kinds)
		}
	})

	t.Run("DeduplicateSameSwapType", func(t *testing.T) {
		logs := []sim.TraceLog{
			{
				Topics: []common.Hash{
					common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"), // SushiSwap
				},
			},
			{
				Topics: []common.Hash{
					common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"), // SushiSwap again
				},
			},
		}
		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if !isSwap {
			t.Error("expected swap detection")
		}
		if len(kinds) != 1 || kinds[0] != "sushiswap_swap" {
			t.Errorf("expected single sushiswap_swap, got %v", kinds)
		}
	})

	t.Run("NoSwapDetected", func(t *testing.T) {
		logs := []sim.TraceLog{
			{
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), // Transfer event
				},
			},
		}
		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if isSwap {
			t.Error("expected no swap detection for Transfer event")
		}
		if len(kinds) != 0 {
			t.Errorf("expected no swap kinds, got %v", kinds)
		}
	})

	t.Run("EmptyLogs", func(t *testing.T) {
		isSwap, kinds := sim.DetectSwapsFromLogs([]sim.TraceLog{})
		if isSwap {
			t.Error("expected no swap detection for empty logs")
		}
		if len(kinds) != 0 {
			t.Errorf("expected no swap kinds, got %v", kinds)
		}
	})

	t.Run("LogWithNoTopics", func(t *testing.T) {
		logs := []sim.TraceLog{
			{
				Topics: []common.Hash{},
			},
		}
		isSwap, kinds := sim.DetectSwapsFromLogs(logs)
		if isSwap {
			t.Error("expected no swap detection for log with no topics")
		}
		if len(kinds) != 0 {
			t.Errorf("expected no swap kinds, got %v", kinds)
		}
	})
}
