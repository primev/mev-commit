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
	fastsettlementv3 "github.com/primev/mev-commit/contracts-abi/clients/FastSettlementV3"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
)

// ============ Types ============

// Intent is an alias to the generated contract binding type for HTTP JSON parsing.
// Uses the exact same struct layout as the Solidity definition.
type Intent = fastsettlementv3.IFastSettlementV3Intent

// SwapCall is an alias to the generated contract binding type.
type SwapCall = fastsettlementv3.IFastSettlementV3SwapCall

// SwapRequest is the HTTP request body for the /fastswap endpoint.
type SwapRequest struct {
	User        common.Address `json:"user"`
	InputToken  common.Address `json:"inputToken"`
	OutputToken common.Address `json:"outputToken"`
	InputAmt    *big.Int       `json:"inputAmt"`
	UserAmtOut  *big.Int       `json:"userAmtOut"`
	Recipient   common.Address `json:"recipient"`
	Deadline    *big.Int       `json:"deadline"`
	Nonce       *big.Int       `json:"nonce"`
	Signature   []byte         `json:"signature"`          // EIP-712 Permit2 signature
	Slippage    string         `json:"slippage,omitempty"` // User slippage percentage (e.g. "1.0" for 1%)
}

// ToIntent converts SwapRequest to the generated Intent type for ABI encoding.
func (r *SwapRequest) ToIntent() Intent {
	return Intent{
		User:        r.User,
		InputToken:  r.InputToken,
		OutputToken: r.OutputToken,
		InputAmt:    r.InputAmt,
		UserAmtOut:  r.UserAmtOut,
		Recipient:   r.Recipient,
		Deadline:    r.Deadline,
		Nonce:       r.Nonce,
	}
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
	InputAmt    *big.Int       `json:"inputAmt"`           // ETH amount in wei
	UserAmtOut  *big.Int       `json:"userAmtOut"`         // minAmountOut from dapp quote
	Sender      common.Address `json:"sender"`             // User address (also recipient)
	Deadline    *big.Int       `json:"deadline"`           // Unix timestamp
	Slippage    string         `json:"slippage,omitempty"` // User slippage percentage (e.g. "1.0" for 1%)
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
	To        common.Address `json:"to"`
	GasLimit  string         `json:"gasLimit"`
	Value     string         `json:"value"`
	Data      string         `json:"data"`
	MinReturn string         `json:"minReturn"` // Guaranteed minimum amount from Barter
	Route     struct {
		OutputAmount  string `json:"outputAmount"`
		GasEstimation uint64 `json:"gasEstimation"`
		BlockNumber   uint64 `json:"blockNumber"`
	} `json:"route"`
}

// barterRequest is the request body for the Barter API.
type barterRequest struct {
	Source            string     `json:"source"`
	Target            string     `json:"target"`
	SellAmount        string     `json:"sellAmount"`
	Recipient         string     `json:"recipient"`
	Origin            string     `json:"origin"`
	MinReturnFraction float64    `json:"minReturnFraction"` // e.g. 0.99 for 1% slippage
	Deadline          string     `json:"deadline"`
	SourceFee         *sourceFee `json:"sourceFee,omitempty"`
}

type sourceFee struct {
	Amount    string `json:"amount"`
	Recipient string `json:"recipient"`
}

// ============ Service ============

// Mainnet WETH address
var mainnetWETH = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

// parsedABI holds the pre-parsed contract ABI from generated bindings
var parsedABI *abi.ABI

func init() {
	// Use pre-parsed ABI from generated bindings - no panic possible
	parsed, err := fastsettlementv3.Fastsettlementv3MetaData.GetAbi()
	if err != nil {
		// This should never happen with generated bindings
		parsed = nil
	}
	parsedABI = parsed
}

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

// NonceStore interface for getting the current nonce from internal store
type NonceStore interface {
	GetCurrentNonce(ctx context.Context, sender common.Address) (uint64, bool)
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
	nonceStore   NonceStore
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
func (s *Service) SetExecutorDeps(signer Signer, txEnqueuer TxEnqueuer, blockTracker BlockTracker, nonceStore NonceStore) {
	s.signer = signer
	s.txEnqueuer = txEnqueuer
	s.blockTracker = blockTracker
	s.nonceStore = nonceStore
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
		s.barterBaseURL+"/swap",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create barter request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.barterAPIKey))

	s.logger.Info("calling Barter API",
		"type", logDescription,
		"inputToken", reqBody.Source,
		"outputToken", reqBody.Target,
		"inputAmount", reqBody.SellAmount,
		"minReturnFraction", reqBody.MinReturnFraction,
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
		return nil, fmt.Errorf("decode barter response: %w", err)
	}

	s.logger.Info("Barter API response",
		"type", logDescription,
		"to", barterResp.To.Hex(),
		"outputAmount", barterResp.Route.OutputAmount,
		"gasLimit", barterResp.GasLimit,
	)

	return &barterResp, nil
}

// CallBarterAPI calls the Barter API for swap routing (Path 1 - executor submitted).
func (s *Service) CallBarterAPI(ctx context.Context, intent Intent, slippageStr string) (*BarterResponse, error) {
	// Default slippage 0.5% if not provided
	fraction := 0.995
	if slippageStr != "" {
		if val, err := strconv.ParseFloat(slippageStr, 64); err == nil && val >= 0 && val <= 100 {
			fraction = 1.0 - (val / 100.0)
		}
	}

	reqBody := barterRequest{
		Source:            intent.InputToken.Hex(),
		Target:            intent.OutputToken.Hex(),
		SellAmount:        intent.InputAmt.String(),
		Recipient:         s.settlementAddr.Hex(),
		Origin:            intent.User.Hex(),
		MinReturnFraction: fraction,
		Deadline:          intent.Deadline.String(),
	}
	resp, err := s.callBarter(ctx, reqBody, "executor-swap")
	if err != nil {
		return nil, err
	}

	// VALIDATION: Ensure Barter's output meets User's requirement
	// We check minReturn (worst case) against user's requirement for safety
	outAmt, ok := new(big.Int).SetString(resp.MinReturn, 10)
	if !ok {
		return nil, fmt.Errorf("invalid minReturn from barter: %s", resp.MinReturn)
	}
	if outAmt.Cmp(intent.UserAmtOut) < 0 {
		// Barter's worst case < User's worst case.
		// Abort to prevent failed transaction.
		return nil, fmt.Errorf("barter minReturn (%s) < user required (%s)", outAmt.String(), intent.UserAmtOut.String())
	}
	return resp, nil
}

// ============ Transaction Building ============

// BuildExecuteTx constructs the calldata for FastSettlementV3.executeWithPermit.
func (s *Service) BuildExecuteTx(intent Intent, signature []byte, barter *BarterResponse) ([]byte, error) {
	if parsedABI == nil {
		return nil, fmt.Errorf("contract ABI not initialized")
	}

	// Decode swap data from Barter response
	swapData, err := hex.DecodeString(strings.TrimPrefix(barter.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("decode barter data: %w", err)
	}

	// Parse value from Barter response with proper error handling
	value, ok := new(big.Int).SetString(barter.Value, 10)
	if !ok {
		value = big.NewInt(0)
	}

	swapCall := SwapCall{
		To:    barter.To,
		Value: value,
		Data:  swapData,
	}

	calldata, err := parsedABI.Pack("executeWithPermit", intent, signature, swapCall)
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

	// Convert request to Intent
	intent := req.ToIntent()

	// 1. Call Barter API using user's slippage if provided, or default
	barterResp, err := s.CallBarterAPI(ctx, intent, req.Slippage)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("barter API error: %v", err),
		}, nil
	}

	// 2. Build execute transaction calldata
	calldata, err := s.BuildExecuteTx(intent, req.Signature, barterResp)
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
	// Use same logic as sender's hasCorrectNonce
	executorAddr := s.signer.GetAddress()
	maxNonce, hasTxs := s.nonceStore.GetCurrentNonce(ctx, executorAddr)
	chainNonce, err := s.blockTracker.AccountNonce(ctx, executorAddr)
	if err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to get chain nonce: %v", err),
		}, nil
	}

	var nonce uint64
	if hasTxs {
		// Has transactions in store, next nonce is max + 1
		nonce = maxNonce + 1
	} else {
		// No transactions in store, use chain nonce
		nonce = chainNonce
	}
	// If chain has advanced beyond our tracking, use chain nonce
	if chainNonce > nonce {
		nonce = chainNonce
	}

	// 5. Calculate gas pricing: GasFeeCap = NextBaseFee only (no tip needed, mev-commit bid handles inclusion)
	nextBaseFee := s.blockTracker.NextBaseFee()
	if nextBaseFee == nil || nextBaseFee.Sign() == 0 {
		nextBaseFee = big.NewInt(30_000_000_000) // 30 gwei fallback
	}
	gasTipCap := big.NewInt(0) // No priority fee - mev-commit bid handles inclusion
	// add buffer to fee cap to account for changes
	gasFeeCap := new(big.Int).Mul(nextBaseFee, big.NewInt(2))

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

	// 9. Enqueue the transaction (uses TxTypeFastSwap to skip balance check)
	senderTx := &sender.Transaction{
		Transaction: signedTx,
		Sender:      executorAddr,
		Raw:         rawTxHex,
		Type:        sender.TxTypeFastSwap,
	}

	if err := s.txEnqueuer.Enqueue(ctx, senderTx); err != nil {
		return &SwapResult{
			Status: "error",
			Error:  fmt.Sprintf("failed to enqueue tx: %v", err),
		}, nil
	}

	s.logger.Info("swap transaction submitted",
		"txHash", signedTx.Hash().Hex(),
		"user", intent.User.Hex(),
		"inputToken", intent.InputToken.Hex(),
		"outputToken", intent.OutputToken.Hex(),
		"inputAmt", intent.InputAmt.String(),
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

// ============ HTTP Handlers ============

// Handler returns an HTTP handler for the /fastswap endpoint.
func (s *Service) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var rawReq struct {
			User        string `json:"user"`
			InputToken  string `json:"inputToken"`
			OutputToken string `json:"outputToken"`
			InputAmt    string `json:"inputAmt"`
			UserAmtOut  string `json:"userAmtOut"`
			Recipient   string `json:"recipient"`
			Deadline    string `json:"deadline"`
			Nonce       string `json:"nonce"`
			Signature   string `json:"signature"`
			Slippage    string `json:"slippage"` // Optional
		}

		if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
			http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if rawReq.User == "" || !common.IsHexAddress(rawReq.User) {
			http.Error(w, "missing or invalid user address", http.StatusBadRequest)
			return
		}
		if rawReq.InputToken == "" || !common.IsHexAddress(rawReq.InputToken) {
			http.Error(w, "missing or invalid inputToken", http.StatusBadRequest)
			return
		}
		if rawReq.OutputToken == "" || !common.IsHexAddress(rawReq.OutputToken) {
			http.Error(w, "missing or invalid outputToken", http.StatusBadRequest)
			return
		}
		if rawReq.Recipient == "" || !common.IsHexAddress(rawReq.Recipient) {
			http.Error(w, "missing or invalid recipient", http.StatusBadRequest)
			return
		}
		if rawReq.Signature == "" {
			http.Error(w, "missing signature", http.StatusBadRequest)
			return
		}

		// Parse big.Int fields
		inputAmt, ok := new(big.Int).SetString(rawReq.InputAmt, 10)
		if !ok || inputAmt.Sign() <= 0 {
			http.Error(w, "invalid inputAmt", http.StatusBadRequest)
			return
		}
		userAmtOut, ok := new(big.Int).SetString(rawReq.UserAmtOut, 10)
		if !ok {
			http.Error(w, "invalid userAmtOut", http.StatusBadRequest)
			return
		}
		deadline, ok := new(big.Int).SetString(rawReq.Deadline, 10)
		if !ok || deadline.Sign() <= 0 {
			http.Error(w, "invalid deadline", http.StatusBadRequest)
			return
		}
		nonce, ok := new(big.Int).SetString(rawReq.Nonce, 10)
		if !ok {
			http.Error(w, "invalid nonce", http.StatusBadRequest)
			return
		}

		// Decode signature from hex
		signature, err := hex.DecodeString(strings.TrimPrefix(rawReq.Signature, "0x"))
		if err != nil {
			http.Error(w, "invalid signature hex", http.StatusBadRequest)
			return
		}

		req := SwapRequest{
			User:        common.HexToAddress(rawReq.User),
			InputToken:  common.HexToAddress(rawReq.InputToken),
			OutputToken: common.HexToAddress(rawReq.OutputToken),
			InputAmt:    inputAmt,
			UserAmtOut:  userAmtOut,
			Recipient:   common.HexToAddress(rawReq.Recipient),
			Deadline:    deadline,
			Nonce:       nonce,
			Signature:   signature,
			Slippage:    rawReq.Slippage,
		}

		result, err := s.HandleSwap(r.Context(), req)
		if err != nil {
			http.Error(w, fmt.Sprintf("swap failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}
}

// ============ Path 2: User-Submitted ETH Swaps ============

// CallBarterAPIForETH calls the Barter API for ETH->Token swap routing (Path 2).
// Uses WETH as the source token since Barter works with ERC20s.
func (s *Service) CallBarterAPIForETH(ctx context.Context, req ETHSwapRequest) (*BarterResponse, error) {
	// Default slippage 0.5% if not provided
	fraction := 0.995
	if req.Slippage != "" {
		if val, err := strconv.ParseFloat(req.Slippage, 64); err == nil && val >= 0 && val <= 100 {
			fraction = 1.0 - (val / 100.0)
		}
	}

	reqBody := barterRequest{
		Source:            mainnetWETH.Hex(),
		Target:            req.OutputToken.Hex(),
		SellAmount:        req.InputAmt.String(),
		Recipient:         s.settlementAddr.Hex(),
		Origin:            req.Sender.Hex(),
		MinReturnFraction: fraction,
		Deadline:          req.Deadline.String(),
	}
	resp, err := s.callBarter(ctx, reqBody, "eth-swap")
	if err != nil {
		return nil, err
	}

	// VALIDATION
	outAmt, ok := new(big.Int).SetString(resp.MinReturn, 10)
	if !ok {
		return nil, fmt.Errorf("invalid minReturn from barter: %s", resp.MinReturn)
	}
	if outAmt.Cmp(req.UserAmtOut) < 0 {
		return nil, fmt.Errorf("barter minReturn (%s) < user required (%s)", outAmt.String(), req.UserAmtOut.String())
	}
	return resp, nil
}

// BuildExecuteWithETHTx constructs the calldata for FastSettlementV3.executeWithETH.
func (s *Service) BuildExecuteWithETHTx(req ETHSwapRequest, barter *BarterResponse) ([]byte, error) {
	if parsedABI == nil {
		return nil, fmt.Errorf("contract ABI not initialized")
	}

	// Decode swap data from Barter response
	swapData, err := hex.DecodeString(strings.TrimPrefix(barter.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("decode barter data: %w", err)
	}

	// Parse value from Barter response with proper error handling
	value, ok := new(big.Int).SetString(barter.Value, 10)
	if !ok {
		value = big.NewInt(0)
	}

	// Build Intent - inputToken is address(0) for ETH
	intent := Intent{
		User:        req.Sender,
		InputToken:  common.Address{}, // address(0) indicates native ETH
		OutputToken: req.OutputToken,
		InputAmt:    req.InputAmt,
		UserAmtOut:  req.UserAmtOut,
		Recipient:   req.Sender, // User is also recipient
		Deadline:    req.Deadline,
		Nonce:       big.NewInt(0), // Unused for Path 2
	}

	swapCall := SwapCall{
		To:    barter.To,
		Value: value,
		Data:  swapData,
	}

	calldata, err := parsedABI.Pack("executeWithETH", intent, swapCall)
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

	// 2. Build executeWithETH calldata
	calldata, err := s.BuildExecuteWithETHTx(req, barterResp)
	if err != nil {
		return &ETHSwapResponse{
			Status: "error",
			Error:  fmt.Sprintf("build tx error: %v", err),
		}, nil
	}

	// 3. Calculate gas limit with buffer
	gasLimit, _ := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	gasLimit += 150000 // Buffer for ETH wrap + settlement contract overhead

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
		Data:     "0x" + hex.EncodeToString(calldata),
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

		var rawReq struct {
			OutputToken string `json:"outputToken"`
			InputAmt    string `json:"inputAmt"`
			UserAmtOut  string `json:"userAmtOut"`
			Sender      string `json:"sender"`
			Deadline    string `json:"deadline"`
			Slippage    string `json:"slippage"` // Optional
		}

		if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
			http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if rawReq.Sender == "" || !common.IsHexAddress(rawReq.Sender) {
			http.Error(w, "missing or invalid sender address", http.StatusBadRequest)
			return
		}
		if rawReq.OutputToken == "" || !common.IsHexAddress(rawReq.OutputToken) {
			http.Error(w, "missing or invalid outputToken", http.StatusBadRequest)
			return
		}

		// Parse big.Int fields
		inputAmt, ok := new(big.Int).SetString(rawReq.InputAmt, 10)
		if !ok || inputAmt.Sign() <= 0 {
			http.Error(w, "invalid inputAmt", http.StatusBadRequest)
			return
		}
		userAmtOut, ok := new(big.Int).SetString(rawReq.UserAmtOut, 10)
		if !ok {
			http.Error(w, "invalid userAmtOut", http.StatusBadRequest)
			return
		}
		deadline, ok := new(big.Int).SetString(rawReq.Deadline, 10)
		if !ok || deadline.Sign() <= 0 {
			http.Error(w, "invalid deadline", http.StatusBadRequest)
			return
		}

		req := ETHSwapRequest{
			OutputToken: common.HexToAddress(rawReq.OutputToken),
			InputAmt:    inputAmt,
			UserAmtOut:  userAmtOut,
			Sender:      common.HexToAddress(rawReq.Sender),
			Deadline:    deadline,
			Slippage:    rawReq.Slippage,
		}

		result, err := s.HandleETHSwap(r.Context(), req)
		if err != nil {
			http.Error(w, fmt.Sprintf("ETH swap failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}
}
