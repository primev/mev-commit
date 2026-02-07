package fastswap_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/fastswap"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"github.com/primev/mev-commit/x/util"
	"github.com/stretchr/testify/require"
)

// ============ Mocks ============

type mockSigner struct {
	address common.Address
}

func (m *mockSigner) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Return the tx unchanged (not actually signed, but tests don't verify signatures)
	return tx, nil
}

func (m *mockSigner) GetAddress() common.Address {
	return m.address
}

type mockTxEnqueuer struct {
	enqueuedTxs []*sender.Transaction
}

func (m *mockTxEnqueuer) Enqueue(ctx context.Context, tx *sender.Transaction) error {
	m.enqueuedTxs = append(m.enqueuedTxs, tx)
	return nil
}

type mockBlockTracker struct {
	nonce       uint64
	nextBaseFee *big.Int
}

func (m *mockBlockTracker) AccountNonce(ctx context.Context, account common.Address) (uint64, error) {
	return m.nonce, nil
}

func (m *mockBlockTracker) NextBaseFee() *big.Int {
	return m.nextBaseFee
}

// mockNonceStore implements fastswap.NonceStore interface
type mockNonceStore struct {
	nonce  uint64
	hasTxs bool
}

func (m *mockNonceStore) GetCurrentNonce(_ context.Context, _ common.Address) (uint64, bool) {
	return m.nonce, m.hasTxs
}

// ============ Test Helpers ============

func newTestBarterResponse() fastswap.BarterResponse {
	// Mock response matching real Barter /swap API format
	// Real example: USDT -> USDC swap returns:
	//   to: swap executor contract
	//   gasLimit: estimated gas as string
	//   value: ETH value (usually "0" for ERC20 swaps)
	//   data: encoded swap calldata
	//   route.outputAmount: expected output tokens
	//   route.gasEstimation: gas estimate as uint64
	//   route.blockNumber: current block number
	return fastswap.BarterResponse{
		To:        common.HexToAddress("0x179dc3fb0f2230094894317f307241a52cdb38aa"), // Barter swap executor
		GasLimit:  "1227112",
		Value:     "0",
		Data:      "0xf0d7bb940000000000000000000000002c0552e5dcb79b064fd23e358a86810bc5994244", // truncated for test
		MinReturn: "250000000",                                                                  // slightly less than output amount
		Route: struct {
			OutputAmount  string `json:"outputAmount"`
			GasEstimation uint64 `json:"gasEstimation"`
			BlockNumber   uint64 `json:"blockNumber"`
		}{
			OutputAmount:  "250212361", // e.g., 250.2 USDC (6 decimals)
			GasEstimation: 217861,
			BlockNumber:   24322525,
		},
	}
}

func setupTestServer(t *testing.T, barterResp fastswap.BarterResponse) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/swap" {
				http.NotFound(w, r)
				return
			}
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Verify authorization header
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing authorization", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(barterResp); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}),
	)
}

// ============ Tests ============

func TestNewService(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		"https://api.barter.com",
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	require.NotNil(t, svc)
}

func TestCallBarterAPI(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	intent := fastswap.Intent{
		User:        common.HexToAddress("0xUserAddress"),
		InputToken:  common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // USDC
		OutputToken: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
		InputAmt:    big.NewInt(1000000000),                                            // 1000 USDC
		UserAmtOut:  big.NewInt(100),                                                   // Low amount to pass validation
		Recipient:   common.HexToAddress("0xRecipientAddress"),
		Deadline:    big.NewInt(1700000000),
		Nonce:       big.NewInt(1),
	}

	ctx := context.Background()
	resp, err := svc.CallBarterAPI(ctx, intent, "")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, barterResp.To, resp.To)
	require.Equal(t, barterResp.Route.OutputAmount, resp.Route.OutputAmount)
}

func TestCallBarterAPIForETH(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	req := fastswap.ETHSwapRequest{
		OutputToken: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // USDC
		InputAmt:    big.NewInt(1000000000000000000),                                   // 1 ETH
		UserAmtOut:  big.NewInt(100),                                                   // Low amount to pass validation
		Sender:      common.HexToAddress("0xSenderAddress"),
		Deadline:    big.NewInt(1700000000),
	}

	ctx := context.Background()
	resp, err := svc.CallBarterAPIForETH(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, barterResp.To, resp.To)
}

func TestBuildExecuteTx(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		"https://api.barter.com",
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	intent := fastswap.Intent{
		User:        common.HexToAddress("0xUserAddress"),
		InputToken:  common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		OutputToken: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		InputAmt:    big.NewInt(1000000000),
		UserAmtOut:  big.NewInt(100),
		Recipient:   common.HexToAddress("0xRecipientAddress"),
		Deadline:    big.NewInt(1700000000),
		Nonce:       big.NewInt(1),
	}

	barterResp := newTestBarterResponse()
	signature := []byte{0x01, 0x02, 0x03, 0x04}

	calldata, err := svc.BuildExecuteTx(intent, signature, &barterResp)

	require.NoError(t, err)
	require.NotEmpty(t, calldata)
	// Check it starts with the function selector (first 4 bytes)
	require.True(t, len(calldata) > 4)
}

func TestBuildExecuteWithETHTx(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		"https://api.barter.com",
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	req := fastswap.ETHSwapRequest{
		OutputToken: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		InputAmt:    big.NewInt(1000000000000000000),
		UserAmtOut:  big.NewInt(100),
		Sender:      common.HexToAddress("0xSenderAddress"),
		Deadline:    big.NewInt(1700000000),
	}

	barterResp := newTestBarterResponse()

	calldata, err := svc.BuildExecuteWithETHTx(req, &barterResp)

	require.NoError(t, err)
	require.NotEmpty(t, calldata)
	require.True(t, len(calldata) > 4)
}

func TestHandleSwap(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	settlementAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		settlementAddr,
		1,
		logger,
	)

	// Set up executor dependencies
	mockSigner := &mockSigner{address: common.HexToAddress("0xExecutorAddress")}
	mockEnqueuer := &mockTxEnqueuer{}
	mockTracker := &mockBlockTracker{
		nonce:       5,
		nextBaseFee: big.NewInt(30000000000), // 30 gwei
	}
	mockStore := &mockNonceStore{nonce: 4, hasTxs: true} // store nonce + 1 should match tracker nonce
	svc.SetExecutorDeps(mockSigner, mockEnqueuer, mockTracker, mockStore)

	req := fastswap.SwapRequest{
		User:        common.HexToAddress("0xUserAddress"),
		InputToken:  common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		OutputToken: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		InputAmt:    big.NewInt(1000000000),
		UserAmtOut:  big.NewInt(100),
		Recipient:   common.HexToAddress("0xRecipientAddress"),
		Deadline:    big.NewInt(1700000000),
		Nonce:       big.NewInt(1),
		Signature:   []byte{0x01, 0x02, 0x03, 0x04},
	}

	ctx := context.Background()
	result, err := svc.HandleSwap(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "success", result.Status)
	require.NotEmpty(t, result.TxHash)
	require.Equal(t, barterResp.Route.OutputAmount, result.OutputAmount)

	// Verify tx was enqueued
	require.Len(t, mockEnqueuer.enqueuedTxs, 1)
	enqueuedTx := mockEnqueuer.enqueuedTxs[0]
	require.Equal(t, mockSigner.address, enqueuedTx.Sender)
}

func TestHandleSwap_NoExecutorDeps(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		"https://api.barter.com",
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)
	// No executor deps set

	req := fastswap.SwapRequest{
		User:       common.HexToAddress("0xUserAddress"),
		InputToken: common.HexToAddress("0xInputToken"),
		InputAmt:   big.NewInt(1000),
		Signature:  []byte{0x01, 0x02},
	}

	ctx := context.Background()
	result, err := svc.HandleSwap(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "error", result.Status)
	require.Contains(t, result.Error, "executor dependencies not configured")
}

func TestHandleETHSwap(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	settlementAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		settlementAddr,
		1,
		logger,
	)

	req := fastswap.ETHSwapRequest{
		OutputToken: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		InputAmt:    big.NewInt(1000000000000000000),
		UserAmtOut:  big.NewInt(100),
		Sender:      common.HexToAddress("0xSenderAddress"),
		Deadline:    big.NewInt(1700000000),
	}

	ctx := context.Background()
	result, err := svc.HandleETHSwap(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "success", result.Status)
	require.Equal(t, settlementAddr.Hex(), result.To)
	require.True(t, strings.HasPrefix(result.Data, "0x"))
	require.Equal(t, req.InputAmt.String(), result.Value)
	require.Equal(t, uint64(1), result.ChainID)
	require.Greater(t, result.GasLimit, uint64(0))
}

func TestHandler(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	settlementAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		settlementAddr,
		1,
		logger,
	)

	mockSignerInst := &mockSigner{address: common.HexToAddress("0xExecutorAddress")}
	mockEnqueuer := &mockTxEnqueuer{}
	mockTracker := &mockBlockTracker{
		nonce:       0,
		nextBaseFee: big.NewInt(30000000000),
	}
	mockStore := &mockNonceStore{nonce: 0, hasTxs: false}
	svc.SetExecutorDeps(mockSignerInst, mockEnqueuer, mockTracker, mockStore)

	handler := svc.Handler()

	// Use raw JSON with string values for numeric fields (new handler format)
	reqJSON := `{
		"user": "0x0000000000000000000000000000000000000001",
		"inputToken": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		"outputToken": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"inputAmt": "1000000000",
		"userAmtOut": "100",
		"recipient": "0x0000000000000000000000000000000000000002",
		"deadline": "1700000000",
		"nonce": "1",
		"signature": "0x01020304"
	}`

	req := httptest.NewRequest(http.MethodPost, "/fastswap", strings.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result fastswap.SwapResult
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService("", "", common.Address{}, 1, logger)
	handler := svc.Handler()

	req := httptest.NewRequest(http.MethodGet, "/fastswap", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestHandler_MissingFields(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService("", "", common.Address{}, 1, logger)
	handler := svc.Handler()

	tests := []struct {
		name     string
		reqBody  map[string]interface{}
		expected string
	}{
		{
			name:     "missing user",
			reqBody:  map[string]interface{}{},
			expected: "missing or invalid user address",
		},
		{
			name: "missing inputToken",
			reqBody: map[string]interface{}{
				"user": "0x1234567890123456789012345678901234567890",
			},
			expected: "missing or invalid inputToken",
		},
		{
			name: "missing outputToken",
			reqBody: map[string]interface{}{
				"user":       "0x1234567890123456789012345678901234567890",
				"inputToken": "0x1234567890123456789012345678901234567890",
			},
			expected: "missing or invalid outputToken",
		},
		{
			name: "missing recipient",
			reqBody: map[string]interface{}{
				"user":        "0x1234567890123456789012345678901234567890",
				"inputToken":  "0x1234567890123456789012345678901234567890",
				"outputToken": "0x1234567890123456789012345678901234567890",
			},
			expected: "missing or invalid recipient",
		},
		{
			name: "missing signature",
			reqBody: map[string]interface{}{
				"user":        "0x1234567890123456789012345678901234567890",
				"inputToken":  "0x1234567890123456789012345678901234567890",
				"outputToken": "0x1234567890123456789012345678901234567890",
				"recipient":   "0x1234567890123456789012345678901234567890",
				"inputAmt":    "1000",
				"userAmtOut":  "900",
				"deadline":    "1700000000",
				"nonce":       "0",
			},
			expected: "missing signature",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/fastswap", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), tc.expected)
		})
	}
}

func TestETHHandler(t *testing.T) {
	barterResp := newTestBarterResponse()
	srv := setupTestServer(t, barterResp)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	settlementAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		settlementAddr,
		1,
		logger,
	)

	handler := svc.ETHHandler()

	// ETH handler expects string values for numeric fields
	reqJSON := `{
		"outputToken": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		"inputAmt": "1000000000000000000",
		"userAmtOut": "100",
		"sender": "0x0000000000000000000000000000000000000001",
		"deadline": "1700000000"
	}`

	req := httptest.NewRequest(http.MethodPost, "/fastswap/eth", strings.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result fastswap.ETHSwapResponse
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
	require.Equal(t, settlementAddr.Hex(), result.To)
}

func TestETHHandler_MethodNotAllowed(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService("", "", common.Address{}, 1, logger)
	handler := svc.ETHHandler()

	req := httptest.NewRequest(http.MethodGet, "/fastswap/eth", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestETHHandler_MissingFields(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService("", "", common.Address{}, 1, logger)
	handler := svc.ETHHandler()

	tests := []struct {
		name     string
		reqBody  map[string]interface{}
		expected string
	}{
		{
			name:     "missing sender",
			reqBody:  map[string]interface{}{},
			expected: "missing or invalid sender",
		},
		{
			name: "missing outputToken",
			reqBody: map[string]interface{}{
				"sender": "0x1234567890123456789012345678901234567890",
			},
			expected: "missing or invalid outputToken",
		},
		{
			name: "invalid inputAmt",
			reqBody: map[string]interface{}{
				"sender":      "0x1234567890123456789012345678901234567890",
				"outputToken": "0x1234567890123456789012345678901234567890",
				"inputAmt":    "invalid",
			},
			expected: "invalid inputAmt",
		},
		{
			name: "invalid userAmtOut",
			reqBody: map[string]interface{}{
				"sender":      "0x1234567890123456789012345678901234567890",
				"outputToken": "0x1234567890123456789012345678901234567890",
				"inputAmt":    "1000",
				"userAmtOut":  "invalid",
			},
			expected: "invalid userAmtOut",
		},
		{
			name: "invalid deadline",
			reqBody: map[string]interface{}{
				"sender":      "0x1234567890123456789012345678901234567890",
				"outputToken": "0x1234567890123456789012345678901234567890",
				"inputAmt":    "1000",
				"userAmtOut":  "900",
				"deadline":    "invalid",
			},
			expected: "invalid deadline",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/fastswap/eth", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), tc.expected)
		})
	}
}

func TestBarterAPIError(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		srv.URL,
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	intent := fastswap.Intent{
		User:        common.HexToAddress("0xUserAddress"),
		InputToken:  common.HexToAddress("0xInputToken"),
		OutputToken: common.HexToAddress("0xOutputToken"),
		InputAmt:    big.NewInt(1000),
		UserAmtOut:  big.NewInt(100),
		Deadline:    big.NewInt(1700000000),
		Nonce:       big.NewInt(1),
	}

	ctx := context.Background()
	resp, err := svc.CallBarterAPI(ctx, intent, "")

	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "barter API error")
}

func TestIntentTupleEncoding(t *testing.T) {
	// Verify that Intent can be used with the ABI encoding
	logger := util.NewTestLogger(os.Stdout)
	svc := fastswap.NewService(
		"https://api.barter.com",
		"test-api-key",
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		1,
		logger,
	)

	intent := fastswap.Intent{
		User:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		InputToken:  common.HexToAddress("0x2222222222222222222222222222222222222222"),
		OutputToken: common.HexToAddress("0x3333333333333333333333333333333333333333"),
		InputAmt:    big.NewInt(1000000000000000000),
		UserAmtOut:  big.NewInt(100),
		Recipient:   common.HexToAddress("0x4444444444444444444444444444444444444444"),
		Deadline:    big.NewInt(1700000000),
		Nonce:       big.NewInt(42),
	}

	barterResp := &fastswap.BarterResponse{
		To:       common.HexToAddress("0x5555555555555555555555555555555555555555"),
		GasLimit: "200000",
		Value:    "0",
		Data:     "0x" + hex.EncodeToString([]byte("test swap data")),
	}

	signature := make([]byte, 65)
	for i := range signature {
		signature[i] = byte(i)
	}

	calldata, err := svc.BuildExecuteTx(intent, signature, barterResp)
	require.NoError(t, err)
	require.NotEmpty(t, calldata)

	// The function selector for executeWithPermit should be at the beginning
	// We can verify the calldata is properly formed by checking its length
	// Function selector (4 bytes) + encoded params
	require.True(t, len(calldata) > 4)
}
