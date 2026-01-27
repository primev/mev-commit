package fastswap

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
)

// ============ Types ============

// Intent mirrors the Solidity Intent struct from FastSettlementV3.
type Intent struct {
	User        common.Address `json:"user"`
	InputToken  common.Address `json:"inputToken"`
	OutputToken common.Address `json:"outputToken"`
	InputAmt    *big.Int       `json:"inputAmt"`
	UserAmtOut  *big.Int       `json:"userAmtOut"`
	Recipient   common.Address `json:"recipient"`
	Deadline    *big.Int       `json:"deadline"`
	Nonce       *big.Int       `json:"nonce"`
}

// SwapCall mirrors the Solidity SwapCall struct.
type SwapCall struct {
	To    common.Address `json:"to"`
	Value *big.Int       `json:"value"`
	Data  []byte         `json:"data"`
}

// SwapRequest is the HTTP request body for the /fastswap endpoint.
type SwapRequest struct {
	Intent    Intent `json:"intent"`
	Signature []byte `json:"signature"` // EIP-712 Permit2 signature
}

// SwapResult is the HTTP response for /fastswap.
type SwapResult struct {
	TxHash       string `json:"txHash,omitempty"` // Transaction hash (when submitted)
	OutputAmount string `json:"outputAmount"`
	GasLimit     uint64 `json:"gasLimit"`
	Status       string `json:"status"` // "success", "error"
	Error        string `json:"error,omitempty"`
}

// ETHSwapRequest is the HTTP request body for /fastswap/eth (user-submitted ETH swaps).
// User swaps native ETH for an ERC20 token and submits the transaction themselves.
type ETHSwapRequest struct {
	OutputToken common.Address `json:"outputToken"`
	InputAmt    *big.Int       `json:"inputAmt"`   // ETH amount in wei
	UserAmtOut  *big.Int       `json:"userAmtOut"` // minAmountOut from dapp quote
	Sender      common.Address `json:"sender"`     // User address (also recipient)
	Deadline    *big.Int       `json:"deadline"`   // Unix timestamp
}

// ETHSwapResponse is the response for /fastswap/eth containing unsigned tx data.
type ETHSwapResponse struct {
	To       string `json:"to"`       // FastSettlementV3 contract address
	Data     string `json:"data"`     // Hex-encoded calldata (0x-prefixed)
	Value    string `json:"value"`    // ETH value to send (same as inputAmt)
	ChainID  uint64 `json:"chainId"`  // Chain ID for the transaction
	GasLimit uint64 `json:"gasLimit"` // Estimated gas limit
	Status   string `json:"status"`   // "success", "error"
	Error    string `json:"error,omitempty"`
}

// BarterResponse represents the parsed response from Barter API.
type BarterResponse struct {
	To       common.Address `json:"to"`
	GasLimit string         `json:"gasLimit"`
	Value    string         `json:"value"`
	Data     string         `json:"data"`
	Route    struct {
		OutputAmount  string `json:"outputAmount"`
		GasEstimation uint64 `json:"gasEstimation"`
		BlockNumber   uint64 `json:"blockNumber"`
	} `json:"route"`
}

// barterRequest is the request body for the Barter API.
type barterRequest struct {
	Source     string     `json:"source"`
	Target     string     `json:"target"`
	SellAmount string     `json:"sellAmount"`
	Recipient  string     `json:"recipient"`
	Origin     string     `json:"origin"`
	MinReturn  string     `json:"minReturn"`
	Deadline   string     `json:"deadline"`
	SourceFee  *sourceFee `json:"sourceFee,omitempty"`
}

type sourceFee struct {
	Amount    string `json:"amount"`
	Recipient string `json:"recipient"`
}

// IntentTuple is a struct that matches the ABI tuple for Intent.
type IntentTuple struct {
	User        common.Address
	InputToken  common.Address
	OutputToken common.Address
	InputAmt    *big.Int
	UserAmtOut  *big.Int
	Recipient   common.Address
	Deadline    *big.Int
	Nonce       *big.Int
}

// SwapCallTuple is a struct that matches the ABI tuple for SwapCall.
type SwapCallTuple struct {
	To    common.Address
	Value *big.Int
	Data  []byte
}

// ============ Service ============

// Mainnet WETH address
var mainnetWETH = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

// Signer interface for signing transactions
type Signer interface {
	SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	GetAddress() common.Address
}

// TxEnqueuer interface for enqueuing transactions to be sent
type TxEnqueuer interface {
	Enqueue(ctx context.Context, tx *sender.Transaction) error
}

// BlockTracker interface for getting nonce and gas pricing
type BlockTracker interface {
	AccountNonce(ctx context.Context, account common.Address) (uint64, error)
	NextBaseFee() *big.Int
}

// Service handles FastSwap operations.
type Service struct {
	barterBaseURL  string
	barterAPIKey   string
	settlementAddr common.Address
	chainID        uint64
	logger         *slog.Logger
	httpClient     *http.Client
	// Path 1 executor submission dependencies (set via SetExecutorDeps)
	signer       Signer
	txEnqueuer   TxEnqueuer
	blockTracker BlockTracker
}

// NewService creates a new FastSwap service.
// For Path 1 executor submission, call SetExecutorDeps after creation.
func NewService(
	barterBaseURL string,
	barterAPIKey string,
	settlementAddr common.Address,
	chainID uint64,
	logger *slog.Logger,
) *Service {
	return &Service{
		barterBaseURL:  barterBaseURL,
		barterAPIKey:   barterAPIKey,
		settlementAddr: settlementAddr,
		chainID:        chainID,
		logger:         logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetExecutorDeps sets the dependencies needed for Path 1 executor transaction submission.
// This is called after TxSender is created since there's a circular dependency.
func (s *Service) SetExecutorDeps(signer Signer, txEnqueuer TxEnqueuer, blockTracker BlockTracker) {
	s.signer = signer
	s.txEnqueuer = txEnqueuer
	s.blockTracker = blockTracker
}

// ============ ABI ============

// FastSettlementV3 ABI for executeWithPermit (executor path)
const executeWithPermitABI = `[{
	"inputs": [
		{
			"components": [
				{"internalType": "address", "name": "user", "type": "address"},
				{"internalType": "address", "name": "inputToken", "type": "address"},
				{"internalType": "address", "name": "outputToken", "type": "address"},
				{"internalType": "uint256", "name": "inputAmt", "type": "uint256"},
				{"internalType": "uint256", "name": "userAmtOut", "type": "uint256"},
				{"internalType": "address", "name": "recipient", "type": "address"},
				{"internalType": "uint256", "name": "deadline", "type": "uint256"},
				{"internalType": "uint256", "name": "nonce", "type": "uint256"}
			],
			"internalType": "struct IFastSettlementV3.Intent",
			"name": "intent",
			"type": "tuple"
		},
		{"internalType": "bytes", "name": "signature", "type": "bytes"},
		{
			"components": [
				{"internalType": "address", "name": "to", "type": "address"},
				{"internalType": "uint256", "name": "value", "type": "uint256"},
				{"internalType": "bytes", "name": "data", "type": "bytes"}
			],
			"internalType": "struct IFastSettlementV3.SwapCall",
			"name": "swapData",
			"type": "tuple"
		}
	],
	"name": "executeWithPermit",
	"outputs": [
		{"internalType": "uint256", "name": "received", "type": "uint256"},
		{"internalType": "uint256", "name": "surplus", "type": "uint256"}
	],
	"stateMutability": "nonpayable",
	"type": "function"
}]`

// FastSettlementV3 ABI for executeWithETH (user ETH swap path)
const executeWithETHABI = `[{
	"inputs": [
		{
			"components": [
				{"internalType": "address", "name": "user", "type": "address"},
				{"internalType": "address", "name": "inputToken", "type": "address"},
				{"internalType": "address", "name": "outputToken", "type": "address"},
				{"internalType": "uint256", "name": "inputAmt", "type": "uint256"},
				{"internalType": "uint256", "name": "userAmtOut", "type": "uint256"},
				{"internalType": "address", "name": "recipient", "type": "address"},
				{"internalType": "uint256", "name": "deadline", "type": "uint256"},
				{"internalType": "uint256", "name": "nonce", "type": "uint256"}
			],
			"internalType": "struct IFastSettlementV3.Intent",
			"name": "intent",
			"type": "tuple"
		},
		{
			"components": [
				{"internalType": "address", "name": "to", "type": "address"},
				{"internalType": "uint256", "name": "value", "type": "uint256"},
				{"internalType": "bytes", "name": "data", "type": "bytes"}
			],
			"internalType": "struct IFastSettlementV3.SwapCall",
			"name": "swapData",
			"type": "tuple"
		}
	],
	"name": "executeWithETH",
	"outputs": [
		{"internalType": "uint256", "name": "received", "type": "uint256"},
		{"internalType": "uint256", "name": "surplus", "type": "uint256"}
	],
	"stateMutability": "payable",
	"type": "function"
}]`

var (
	executeWithPermitABIParsed abi.ABI
	executeWithETHABIParsed    abi.ABI
)

func init() {
	var err error
	executeWithPermitABIParsed, err = abi.JSON(strings.NewReader(executeWithPermitABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse executeWithPermit ABI: %v", err))
	}
	executeWithETHABIParsed, err = abi.JSON(strings.NewReader(executeWithETHABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse executeWithETH ABI: %v", err))
	}
}

// ============ Barter API ============

// callBarter is the shared HTTP call logic for calling the Barter swap API.
func (s *Service) callBarter(ctx context.Context, reqBody barterRequest, logDescription string) (*BarterResponse, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal barter request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/swap", s.barterBaseURL),
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create barter request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.barterAPIKey))

	s.logger.Debug("calling Barter API",
		"type", logDescription,
		"url", req.URL.String(),
		"source", reqBody.Source,
		"target", reqBody.Target,
		"sellAmount", reqBody.SellAmount,
	)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("barter API request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read barter response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("barter API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var barterResp BarterResponse
	if err := json.Unmarshal(respBody, &barterResp); err != nil {
		return nil, fmt.Errorf("unmarshal barter response: %w", err)
	}

	s.logger.Info("Barter API response",
		"type", logDescription,
		"to", barterResp.To.Hex(),
		"outputAmount", barterResp.Route.OutputAmount,
		"gasLimit", barterResp.GasLimit,
	)

	return &barterResp, nil
}

// CallBarterAPI calls the Barter swap API with the given intent (Path 1).
func (s *Service) CallBarterAPI(ctx context.Context, intent Intent) (*BarterResponse, error) {
	// IMPORTANT: Recipient must be the settlement contract, not the user's recipient.
	// The contract receives the swap output, then distributes:
	//   - userAmtOut → intent.Recipient
	//   - surplus → treasury
	reqBody := barterRequest{
		Source:     intent.InputToken.Hex(),
		Target:     intent.OutputToken.Hex(),
		SellAmount: intent.InputAmt.String(),
		Recipient:  s.settlementAddr.Hex(),
		Origin:     intent.User.Hex(),
		MinReturn:  intent.UserAmtOut.String(),
		Deadline:   intent.Deadline.String(),
	}
	return s.callBarter(ctx, reqBody, "executor-swap")
}

// ============ Transaction Builder ============

// BuildExecuteTx constructs the calldata for FastSettlementV3.execute.
func (s *Service) BuildExecuteTx(intent Intent, signature []byte, barter *BarterResponse) ([]byte, error) {
	// Decode swap data from Barter response
	swapData, err := hex.DecodeString(strings.TrimPrefix(barter.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("decode barter data: %w", err)
	}

	// Parse value from Barter response
	value := new(big.Int)
	if barter.Value != "" && barter.Value != "0" {
		value.SetString(barter.Value, 10)
	}

	intentTuple := IntentTuple(intent)

	swapCallTuple := SwapCallTuple{
		To:    barter.To,
		Value: value,
		Data:  swapData,
	}

	calldata, err := executeWithPermitABIParsed.Pack("executeWithPermit", intentTuple, signature, swapCallTuple)
	if err != nil {
		return nil, fmt.Errorf("pack executeWithPermit calldata: %w", err)
	}

	s.logger.Debug("built executeWithPermit calldata",
		"calldataLen", len(calldata),
		"swapTarget", barter.To.Hex(),
	)

	return calldata, nil
}

// ============ Handler ============

// HandleSwap is the main orchestrator for processing a swap request.
// For Path 1, it builds, signs, and enqueues the transaction for the executor.
func (s *Service) HandleSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Validate executor dependencies are set
	if s.signer == nil || s.txEnqueuer == nil || s.blockTracker == nil {
		return &SwapResult{
			Status: "error",
			Error:  "executor dependencies not configured",
		}, nil
	}

	// 1. Call Barter API
	barterResp, err := s.CallBarterAPI(ctx, req.Intent)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("barter API error: %v", err),
		}, nil
	}

	// 2. Build execute transaction calldata
	calldata, err := s.BuildExecuteTx(req.Intent, req.Signature, barterResp)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("build tx error: %v", err),
		}, nil
	}

	// 3. Parse gas limit and add buffer
	gasLimit, _ := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	gasLimit += 100000 // Buffer for settlement contract overhead

	// 4. Get nonce for executor wallet
	executorAddr := s.signer.GetAddress()
	nonce, err := s.blockTracker.AccountNonce(ctx, executorAddr)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to get nonce: %v", err),
		}, nil
	}

	// 5. Calculate gas pricing: GasFeeCap = NextBaseFee only (no tip needed, mev-commit bid handles inclusion)
	nextBaseFee := s.blockTracker.NextBaseFee()
	if nextBaseFee == nil || nextBaseFee.Sign() == 0 {
		nextBaseFee = big.NewInt(30_000_000_000) // 30 gwei fallback
	}
	gasTipCap := big.NewInt(0) // No priority fee - mev-commit bid handles inclusion
	gasFeeCap := nextBaseFee   // GasFeeCap = BaseFee (tip is 0)

	// 6. Build the transaction
	chainID := big.NewInt(int64(s.chainID))
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &s.settlementAddr,
		Value:     big.NewInt(0), // executeWithPermit is non-payable
		Data:      calldata,
	})

	// 7. Sign the transaction
	signedTx, err := s.signer.SignTx(tx, chainID)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to sign tx: %v", err),
		}, nil
	}

	// 8. Encode to raw hex
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to encode tx: %v", err),
		}, nil
	}
	rawTxHex := "0x" + hex.EncodeToString(rawTxBytes)

	// 9. Enqueue the transaction
	senderTx := &sender.Transaction{
		Transaction: signedTx,
		Sender:      executorAddr,
		Raw:         rawTxHex,
		Type:        sender.TxTypeRegular,
	}

	if err := s.txEnqueuer.Enqueue(ctx, senderTx); err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to enqueue tx: %v", err),
		}, nil
	}

	s.logger.Info("swap transaction submitted",
		"txHash", signedTx.Hash().Hex(),
		"user", req.Intent.User.Hex(),
		"inputToken", req.Intent.InputToken.Hex(),
		"outputToken", req.Intent.OutputToken.Hex(),
		"inputAmt", req.Intent.InputAmt.String(),
		"outputAmount", barterResp.Route.OutputAmount,
		"gasLimit", gasLimit,
		"gasFeeCap", gasFeeCap.String(),
		"gasTipCap", gasTipCap.String(),
		"nonce", nonce,
	)

	return &SwapResult{
		TxHash:       signedTx.Hash().Hex(),
		OutputAmount: barterResp.Route.OutputAmount,
		GasLimit:     gasLimit,
		Status:       "success",
	}, nil
}

// Handler returns an HTTP handler for the /fastswap endpoint.
func (s *Service) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Limit body size
		r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024) // 1MB
		defer func() { _ = r.Body.Close() }()

		var req SwapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Intent.User == (common.Address{}) {
			http.Error(w, "missing intent.user", http.StatusBadRequest)
			return
		}
		if req.Intent.InputToken == (common.Address{}) {
			http.Error(w, "missing intent.inputToken", http.StatusBadRequest)
			return
		}
		if req.Intent.OutputToken == (common.Address{}) {
			http.Error(w, "missing intent.outputToken", http.StatusBadRequest)
			return
		}
		if req.Intent.InputAmt == nil || req.Intent.InputAmt.Sign() <= 0 {
			http.Error(w, "invalid intent.inputAmt", http.StatusBadRequest)
			return
		}
		if len(req.Signature) == 0 {
			http.Error(w, "missing signature", http.StatusBadRequest)
			return
		}

		result, err := s.HandleSwap(r.Context(), req)
		if err != nil {
			s.logger.Error("HandleSwap error", "error", err)
			http.Error(w, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if result.Status == "error" {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			s.logger.Error("failed to encode response", "error", err)
		}
	}
}

// ============ Path 2: User ETH Swap ============

// CallBarterAPIForETH calls the Barter swap API for an ETH swap (Path 2).
// Uses WETH as the source token since Barter works with ERC20s.
func (s *Service) CallBarterAPIForETH(ctx context.Context, req ETHSwapRequest) (*BarterResponse, error) {
	reqBody := barterRequest{
		Source:     mainnetWETH.Hex(), // WETH (Barter uses ERC20s)
		Target:     req.OutputToken.Hex(),
		SellAmount: req.InputAmt.String(),
		Recipient:  s.settlementAddr.Hex(),
		Origin:     req.Sender.Hex(),
		MinReturn:  req.UserAmtOut.String(),
		Deadline:   req.Deadline.String(),
	}
	return s.callBarter(ctx, reqBody, "eth-swap")
}

// BuildExecuteWithETHTx constructs the calldata for FastSettlementV3.executeWithETH.
func (s *Service) BuildExecuteWithETHTx(req ETHSwapRequest, barter *BarterResponse) ([]byte, error) {
	// Decode swap data from Barter response
	swapData, err := hex.DecodeString(strings.TrimPrefix(barter.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("decode barter data: %w", err)
	}

	// Parse value from Barter response
	value := new(big.Int)
	if barter.Value != "" && barter.Value != "0" {
		value.SetString(barter.Value, 10)
	}

	// Build Intent tuple - inputToken is address(0) for ETH
	intentTuple := IntentTuple{
		User:        req.Sender,
		InputToken:  common.Address{}, // address(0) indicates native ETH
		OutputToken: req.OutputToken,
		InputAmt:    req.InputAmt,
		UserAmtOut:  req.UserAmtOut,
		Recipient:   req.Sender, // User is also recipient
		Deadline:    req.Deadline,
		Nonce:       big.NewInt(0), // Unused for Path 2
	}

	swapCallTuple := SwapCallTuple{
		To:    barter.To,
		Value: value,
		Data:  swapData,
	}

	calldata, err := executeWithETHABIParsed.Pack("executeWithETH", intentTuple, swapCallTuple)
	if err != nil {
		return nil, fmt.Errorf("pack executeWithETH calldata: %w", err)
	}

	s.logger.Debug("built executeWithETH calldata",
		"calldataLen", len(calldata),
		"swapTarget", barter.To.Hex(),
	)

	return calldata, nil
}

// HandleETHSwap is the main orchestrator for processing an ETH swap request (Path 2).
func (s *Service) HandleETHSwap(ctx context.Context, req ETHSwapRequest) (*ETHSwapResponse, error) {
	// 1. Call Barter API
	barterResp, err := s.CallBarterAPIForETH(ctx, req)
	if err != nil {
		return &ETHSwapResponse{
			Status: "error",
			Error:  fmt.Sprintf("barter API error: %v", err),
		}, nil
	}

	// 2. Build executeWithETH transaction calldata
	txData, err := s.BuildExecuteWithETHTx(req, barterResp)
	if err != nil {
		return &ETHSwapResponse{
			Status: "error",
			Error:  fmt.Sprintf("build tx error: %v", err),
		}, nil
	}

	// 3. Parse gas limit
	gasLimit, _ := strconv.ParseUint(barterResp.GasLimit, 10, 64)

	// Add buffer for settlement contract overhead (WETH wrap, approve, transfer, etc.)
	// Rough estimate: 150k gas for settlement logic with WETH wrapping
	gasLimit += 150000

	s.logger.Info("ETH swap request processed",
		"sender", req.Sender.Hex(),
		"outputToken", req.OutputToken.Hex(),
		"inputAmt", req.InputAmt.String(),
		"outputAmount", barterResp.Route.OutputAmount,
		"gasEstimation", barterResp.Route.GasEstimation,
		"gasLimit", gasLimit,
	)

	return &ETHSwapResponse{
		To:       s.settlementAddr.Hex(),
		Data:     "0x" + hex.EncodeToString(txData),
		Value:    req.InputAmt.String(),
		ChainID:  s.chainID,
		GasLimit: gasLimit,
		Status:   "success",
	}, nil
}

// ETHHandler returns an HTTP handler for the /fastswap/eth endpoint.
func (s *Service) ETHHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Limit body size
		r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024) // 1MB
		defer func() { _ = r.Body.Close() }()

		var req ETHSwapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Sender == (common.Address{}) {
			http.Error(w, "missing sender", http.StatusBadRequest)
			return
		}
		if req.OutputToken == (common.Address{}) {
			http.Error(w, "missing outputToken", http.StatusBadRequest)
			return
		}
		if req.InputAmt == nil || req.InputAmt.Sign() <= 0 {
			http.Error(w, "invalid inputAmt", http.StatusBadRequest)
			return
		}
		if req.UserAmtOut == nil || req.UserAmtOut.Sign() <= 0 {
			http.Error(w, "invalid userAmtOut", http.StatusBadRequest)
			return
		}
		if req.Deadline == nil || req.Deadline.Sign() <= 0 {
			http.Error(w, "invalid deadline", http.StatusBadRequest)
			return
		}

		result, err := s.HandleETHSwap(r.Context(), req)
		if err != nil {
			s.logger.Error("HandleETHSwap error", "error", err)
			http.Error(w, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if result.Status == "error" {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			s.logger.Error("failed to encode response", "error", err)
		}
	}
}
