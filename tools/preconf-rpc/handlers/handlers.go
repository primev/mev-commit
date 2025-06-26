package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"resenje.org/multex"
)

const (
	blockTime = 12 // seconds, typical Ethereum block time
)

var (
	preconfBlockHashPrefix = hex.EncodeToString([]byte("mev-commit"))
)

type Bidder interface {
	Estimate() (int64, error)
}

type Pricer interface {
	EstimatePrice(
		ctx context.Context,
		txn *types.Transaction,
	) (*pricer.BlockPrice, error)
}

type Store interface {
	GetPreconfirmedTransaction(
		ctx context.Context,
		txnHash common.Hash,
	) (*types.Transaction, []*bidderapiv1.Commitment, error)
	GetPreconfirmedTransactionsForBlock(
		ctx context.Context,
		blockNumber int64,
	) ([]*types.Transaction, error)
	GetBalance(ctx context.Context, account common.Address) (*big.Int, error)
}

type BlockTracker interface {
	LatestBlockNumber() uint64
}

type Sender interface {
	Enqueue(txn *sender.Transaction) error
}

type accountNonce struct {
	Account string `json:"account"`
	Nonce   uint64 `json:"nonce"`
	Block   int64  `json:"block"`
}

type bidResult struct {
	noOfProviders int
	blockNumber   uint64
	optedInSlot   bool
	bidAmount     *big.Int
	commitments   []*bidderapiv1.Commitment
}

type rpcMethodHandler struct {
	logger         *slog.Logger
	bidder         Bidder
	pricer         Pricer
	store          Store
	blockTracker   BlockTracker
	sndr           Sender
	depositAddress common.Address
	bridgeAddress  common.Address
	chainID        *big.Int
	nonceLock      *multex.Multex[string]
	nonceMap       map[string]accountNonce
	nonceMapLock   sync.RWMutex
}

func NewRPCMethodHandler(
	logger *slog.Logger,
	bidder Bidder,
	store Store,
	blockTracker BlockTracker,
	sndr Sender,
	depositAddress common.Address,
	bridgeAddress common.Address,
	chainId *big.Int,
) *rpcMethodHandler {
	return &rpcMethodHandler{
		logger:         logger,
		bidder:         bidder,
		store:          store,
		blockTracker:   blockTracker,
		sndr:           sndr,
		depositAddress: depositAddress,
		bridgeAddress:  bridgeAddress,
		chainID:        chainId,
		nonceLock:      multex.New[string](),
		nonceMap:       make(map[string]accountNonce),
	}
}

func (h *rpcMethodHandler) RegisterMethods(server *rpcserver.JSONRPCServer) {
	// Ethereum JSON-RPC methods overridden
	server.RegisterHandler("eth_getBlockNumber", h.handleGetBlockNumber)
	server.RegisterHandler("eth_chainId", h.handleChainID)
	server.RegisterHandler("eth_sendRawTransaction", h.handleSendRawTx)
	server.RegisterHandler("eth_getTransactionReceipt", h.handleGetTxReceipt)
	server.RegisterHandler("eth_getTransactionCount", h.handleGetTxCount)
	server.RegisterHandler("eth_getBlockByHash", h.handleGetBlockByHash)
	// Custom methods for MEV Commit
	server.RegisterHandler("mevcommit_getTransactionCommitments", h.handleGetTxCommitments)
	server.RegisterHandler("mevcommit_getBalance", h.handleMevCommitGetBalance)
	server.RegisterHandler("mevcommit_optInBlock", h.handleMevCommitOptInBlock)
}

func (h *rpcMethodHandler) handleGetBlockNumber(
	ctx context.Context,
	params ...any,
) (json.RawMessage, bool, error) {
	if len(params) != 0 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getBlockNumber does not require any parameters",
		)
	}

	blockNumber := h.blockTracker.LatestBlockNumber()
	h.logger.Info("Retrieved latest block number", "blockNumber", blockNumber)

	blockNumberJSON, err := json.Marshal(hexutil.Uint64(blockNumber))
	if err != nil {
		h.logger.Error("Failed to marshal block number to JSON", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal block number",
		)
	}

	return blockNumberJSON, false, nil
}

func (h *rpcMethodHandler) handleChainID(
	ctx context.Context,
	params ...any,
) (json.RawMessage, bool, error) {
	if len(params) != 0 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"chainID does not require any parameters",
		)
	}

	chainIDJSON, err := json.Marshal(hexutil.Uint64(h.chainID.Uint64()))
	if err != nil {
		h.logger.Error("Failed to marshal chain ID to JSON", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal chain ID",
		)
	}

	return chainIDJSON, false, nil
}

func (h *rpcMethodHandler) handleGetBlockByHash(
	ctx context.Context,
	params ...any,
) (json.RawMessage, bool, error) {
	if len(params) == 0 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getBlockByHash requires one or two parameter",
		)
	}

	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getBlock parameter cannot be null",
		)
	}

	blockHashStr := params[0].(string)
	if !strings.HasPrefix(blockHashStr, preconfBlockHashPrefix) {
		return nil, true, nil // Not a preconf block hash, proxy
	}

	details := false
	if len(params) > 1 && params[1] != nil {
		details, _ = params[1].(bool)
	}

	blockNumberWithPadding := strings.TrimPrefix(blockHashStr, preconfBlockHashPrefix)
	blockNumber, err := strconv.ParseUint(blockNumberWithPadding[:8], 10, 64)
	if err != nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getBlock parameter must be a valid preconf block hash",
		)
	}

	txns, err := h.store.GetPreconfirmedTransactionsForBlock(ctx, int64(blockNumber))
	if err != nil {
		h.logger.Error("Failed to get preconfirmed transactions for block", "error", err, "blockNumber", blockNumber)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get preconfirmed transactions for block",
		)
	}

	block := map[string]interface{}{
		"number":           hexutil.Uint64(blockNumber),
		"hash":             blockHashStr,
		"parentHash":       (common.Hash{}).Hex(),
		"nonce":            "0x0000000000000000",
		"sha3Uncles":       (common.Hash{}).Hex(),
		"logsBloom":        hexutil.Bytes(types.Bloom{}.Bytes()),
		"transactionsRoot": (common.Hash{}).Hex(),
		"stateRoot":        (common.Hash{}).Hex(),
		"miner":            (common.Address{}).Hex(),
		"difficulty":       hexutil.Uint64(0),
		"totalDifficulty":  hexutil.Uint64(0),
		"size":             hexutil.Uint64(0),
		"extraData":        "0x",
		"gasLimit":         hexutil.Uint64(0),
		"gasUsed":          hexutil.Uint64(0),
		"timestamp":        hexutil.Uint64(0),
		"baseFeePerGas":    hexutil.EncodeBig(big.NewInt(0)),
		"withdrawals":      nil,
	}

	var txnsToReturn any
	for i, txn := range txns {
		if !details {
			if txnsToReturn == nil {
				txnsToReturn = make([]string, 0, len(txns))
			}
			txnsToReturn = append(
				txnsToReturn.([]string),
				txn.Hash().Hex(),
			)
			continue
		}
		if txnsToReturn == nil {
			txnsToReturn = make([]map[string]interface{}, len(txns))
		}
		r, s, v := txn.RawSignatureValues()
		sender, err := types.Sender(types.LatestSignerForChainID(txn.ChainId()), txn)
		if err != nil {
			h.logger.Error("Failed to get transaction sender", "error", err, "txnHash", txn.Hash().Hex())
			continue
		}
		txnsToReturn = append(
			txnsToReturn.([]map[string]interface{}),
			map[string]interface{}{
				"hash":                 txn.Hash().Hex(),
				"blockHash":            blockHashStr,
				"blockNumber":          hexutil.Uint64(blockNumber),
				"transactionIndex":     hexutil.Uint64(i),
				"type":                 hexutil.Uint(txn.Type()),
				"accessList":           nil, // Access lists are not used in preconf blocks
				"maxFeePerGas":         hexutil.EncodeBig(txn.GasFeeCap()),
				"maxPriorityFeePerGas": hexutil.EncodeBig(txn.GasTipCap()),
				"to":                   txn.To().Hex(),
				"value":                hexutil.EncodeBig(txn.Value()),
				"input":                hexutil.Encode(txn.Data()),
				"from":                 sender.Hex(),
				"nonce":                hexutil.Uint64(txn.Nonce()),
				"gas":                  hexutil.Uint64(txn.Gas()),
				"gasPrice":             hexutil.EncodeBig(txn.GasPrice()),
				"r":                    hexutil.EncodeBig(r),
				"s":                    hexutil.EncodeBig(s),
				"v":                    hexutil.EncodeBig(v),
			},
		)
	}
	block["transactions"] = txnsToReturn
	blockJSON, err := json.Marshal(block)
	if err != nil {
		h.logger.Error("Failed to marshal block to JSON", "error", err, "blockNumber", blockNumber)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal block",
		)
	}

	h.logger.Info("Retrieved preconf block", "blockNumber", blockNumber, "txCount", len(txns))
	return blockJSON, false, nil
}

func (h *rpcMethodHandler) handleSendRawTx(
	ctx context.Context,
	params ...any,
) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"sendRawTx requires exactly one parameter",
		)
	}
	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"sendRawTx parameter cannot be null",
		)
	}

	rawTxHex := params[0].(string)
	if len(rawTxHex) < 2 || rawTxHex[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"sendRawTx parameter must be a hex string starting with '0x'",
		)
	}

	decodedTxn, err := hex.DecodeString(rawTxHex[2:])
	if err != nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"sendRawTx parameter must be a valid hex string",
		)
	}

	txn := new(types.Transaction)
	if err := txn.UnmarshalBinary(decodedTxn); err != nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"sendRawTx parameter must be a valid transaction",
		)
	}

	txSender, err := types.Sender(types.LatestSignerForChainID(txn.ChainId()), txn)
	if err != nil {
		h.logger.Error("Failed to get transaction sender", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get transaction sender",
		)
	}

	txType := sender.TxTypeRegular
	switch {
	case txn.To().Cmp(h.depositAddress) == 0:
		txType = sender.TxTypeDeposit
	case txn.To().Cmp(h.bridgeAddress) == 0:
		txType = sender.TxTypeInstantBridge
	}

	err = h.sndr.Enqueue(&sender.Transaction{
		Transaction: txn,
		Raw:         rawTxHex,
		Sender:      txSender,
		Type:        txType,
	})
	if err != nil {
		h.logger.Error("Failed to enqueue transaction for sending", "error", err, "sender", txSender.Hex())
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to enqueue transaction for sending",
		)
	}

	txHashJSON, err := json.Marshal(txn.Hash().Hex())
	if err != nil {
		h.logger.Error("Failed to marshal transaction hash to JSON", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal transaction hash",
		)
	}

	return txHashJSON, false, nil
}

func (h *rpcMethodHandler) handleGetTxReceipt(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getTxReceipt requires exactly one parameter",
		)
	}
	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxReceipt parameter cannot be null",
		)
	}

	txHashStr := params[0].(string)
	if len(txHashStr) < 2 || txHashStr[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxReceipt parameter must be a hex string starting with '0x'",
		)
	}

	txHash := common.HexToHash(txHashStr)

	h.logger.Info("Retrieving transaction receipt", "txHash", txHash)
	txn, commitments, err := h.store.GetPreconfirmedTransaction(ctx, txHash)
	if err != nil {
		return nil, true, nil
	}

	if h.blockTracker.LatestBlockNumber() > uint64(commitments[0].BlockNumber) {
		return nil, true, nil
	}

	sender, err := types.Sender(types.LatestSignerForChainID(txn.ChainId()), txn)
	if err != nil {
		h.logger.Error("Failed to get transaction sender", "error", err, "txHash", txHash)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get transaction sender",
		)
	}

	blockHash := fmt.Sprintf("%s%08d", preconfBlockHashPrefix, commitments[0].BlockNumber)
	padding := strings.Repeat("0", 66-len(blockHash))
	blockHash = blockHash + padding

	result := map[string]interface{}{
		"type":              hexutil.Uint(txn.Type()),
		"transactionHash":   txn.Hash().Hex(),
		"transactionIndex":  hexutil.Uint(0),
		"blockHash":         blockHash,
		"blockNumber":       hexutil.EncodeBig(big.NewInt(commitments[0].BlockNumber)),
		"from":              sender.Hex(),
		"to":                nil,
		"contractAddress":   (common.Address{}).Hex(),
		"gasUsed":           hexutil.Uint64(0),
		"cumulativeGasUsed": hexutil.Uint64(1),
		"logs":              []*types.Log{}, // should be [] not null
		"logsBloom":         hexutil.Bytes(types.Bloom{}.Bytes()),
		"status":            hexutil.Uint64(types.ReceiptStatusSuccessful),
		"effectiveGasPrice": hexutil.EncodeBig(big.NewInt(0)),
	}

	receiptJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal receipt to JSON", "error", err, "txHash", txHash)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal receipt",
		)
	}

	return receiptJSON, false, nil
}

func (h *rpcMethodHandler) handleGetTxCount(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
	if len(params) == 2 {
		state := params[1].(string)
		if state != "latest" && state != "pending" {
			return nil, true, nil
		}
	}

	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxCount parameter cannot be null",
		)
	}

	account := params[0].(string)
	if len(account) < 2 || account[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxCount parameter must be a hex string starting with '0x'",
		)
	}

	h.nonceLock.Lock(account)
	defer h.nonceLock.Unlock(account)

	h.nonceMapLock.RLock()
	accNonce, found := h.nonceMap[account]
	h.nonceMapLock.RUnlock()

	if !found {
		return nil, true, nil
	}

	if h.blockTracker.LatestBlockNumber() > uint64(accNonce.Block) {
		h.nonceMapLock.Lock()
		delete(h.nonceMap, account)
		h.nonceMapLock.Unlock()
		return nil, true, nil
	}

	nonceJSON, err := json.Marshal(accNonce.Nonce)
	if err != nil {
		h.logger.Error("Failed to marshal nonce to JSON", "error", err, "account", account)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal nonce",
		)
	}

	h.logger.Info("Retrieved account nonce from cache", "account", account, "nonce", accNonce.Nonce)
	return nonceJSON, false, nil
}

func (h *rpcMethodHandler) handleGetTxCommitments(
	ctx context.Context,
	params ...any,
) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getTxCommitments requires exactly one parameter",
		)
	}

	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxCommitments parameter cannot be null",
		)
	}

	txHashStr := params[0].(string)
	if len(txHashStr) < 2 || txHashStr[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxCommitments parameter must be a hex string starting with '0x'",
		)
	}

	txHash := common.HexToHash(txHashStr)

	_, commitments, err := h.store.GetPreconfirmedTransaction(ctx, txHash)
	if err != nil {
		return nil, true, nil
	}

	if len(commitments) == 0 {
		h.logger.Info("No commitments found for transaction", "txHash", txHash)
		return json.RawMessage("[]"), false, nil
	}

	commitmentsJSON, err := json.Marshal(commitments)
	if err != nil {
		h.logger.Error("Failed to marshal commitments to JSON", "error", err, "txHash", txHash)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal commitments",
		)
	}

	return commitmentsJSON, false, nil
}

func (h *rpcMethodHandler) handleMevCommitGetBalance(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"mevcommit_getBalance requires exactly one parameter",
		)
	}

	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"mevcommit_getBalance parameters cannot be null",
		)
	}

	account := params[0].(string)
	if len(account) < 2 || account[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"mevcommit_getBalance account must be a hex string starting with '0x'",
		)
	}

	balance, err := h.store.GetBalance(ctx, common.HexToAddress(account))
	if err != nil {
		h.logger.Error("Failed to get balance for account", "error", err, "account", account)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get balance for account",
		)
	}

	return json.RawMessage(fmt.Sprintf(`{"balance": "%s"}`, balance)), false, nil
}

func (h *rpcMethodHandler) handleMevCommitOptInBlock(
	ctx context.Context,
	_ ...any,
) (json.RawMessage, bool, error) {
	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate time to opt in", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to estimate opt in time",
		)
	}

	return json.RawMessage(fmt.Sprintf(`{"timeInSecs": "%d"}`, timeToOptIn)), false, nil
}
