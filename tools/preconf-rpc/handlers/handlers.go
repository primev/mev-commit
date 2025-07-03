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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
)

var (
	preconfBlockHashPrefix = hex.EncodeToString([]byte("mev-commit"))
)

type Bidder interface {
	Estimate() (int64, error)
}

type Pricer interface {
	EstimatePrice(ctx context.Context, txn *types.Transaction) (*pricer.BlockPrice, error)
}

type Store interface {
	GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*sender.Transaction, error)
	GetTransactionsForBlock(ctx context.Context, blockNumber int64) ([]*sender.Transaction, error)
	GetTransactionCommitments(ctx context.Context, txnHash common.Hash) ([]*bidderapiv1.Commitment, error)
	GetBalance(ctx context.Context, account common.Address) (*big.Int, error)
	GetCurrentNonce(ctx context.Context, account common.Address) uint64
}

type BlockTracker interface {
	LatestBlockNumber() uint64
}

type Sender interface {
	Enqueue(ctx context.Context, txn *sender.Transaction) error
}

type rpcMethodHandler struct {
	logger         *slog.Logger
	pricer         Pricer
	bidder         Bidder
	store          Store
	blockTracker   BlockTracker
	sndr           Sender
	depositAddress common.Address
	bridgeAddress  common.Address
	chainID        *big.Int
}

func NewRPCMethodHandler(
	logger *slog.Logger,
	pricer Pricer,
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
		pricer:         pricer,
		bidder:         bidder,
		store:          store,
		blockTracker:   blockTracker,
		sndr:           sndr,
		depositAddress: depositAddress,
		bridgeAddress:  bridgeAddress,
		chainID:        chainId,
	}
}

func (h *rpcMethodHandler) RegisterMethods(server *rpcserver.JSONRPCServer) {
	// Ethereum JSON-RPC methods overridden
	server.RegisterHandler("eth_getBlockNumber", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		blockNumber := h.blockTracker.LatestBlockNumber()

		blockNumberJSON, err := json.Marshal(hexutil.Uint64(blockNumber))
		if err != nil {
			h.logger.Error("Failed to marshal block number to JSON", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to marshal block number",
			)
		}

		return blockNumberJSON, false, nil
	})
	server.RegisterHandler("eth_chainId", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		chainIDJSON, err := json.Marshal(hexutil.Uint64(h.chainID.Uint64()))
		if err != nil {
			h.logger.Error("Failed to marshal chain ID to JSON", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to marshal chain ID",
			)
		}
		return chainIDJSON, false, nil
	})
	server.RegisterHandler("eth_maxPriorityFeePerGas", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		// Return zero value for maxPriorityFeePerGas
		maxPriorityFee := big.NewInt(0)
		maxPriorityFeeJSON, err := json.Marshal(hexutil.EncodeBig(maxPriorityFee))
		if err != nil {
			h.logger.Error("Failed to marshal maxPriorityFeePerGas to JSON", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to marshal maxPriorityFeePerGas",
			)
		}
		return maxPriorityFeeJSON, false, nil
	})
	server.RegisterHandler("eth_sendRawTransaction", h.handleSendRawTx)
	server.RegisterHandler("eth_getTransactionReceipt", h.handleGetTxReceipt)
	server.RegisterHandler("eth_getTransactionCount", h.handleGetTxCount)
	server.RegisterHandler("eth_getBlockByHash", h.handleGetBlockByHash)
	// Custom methods for MEV Commit
	server.RegisterHandler("mevcommit_getTransactionCommitments", h.handleGetTxCommitments)
	server.RegisterHandler("mevcommit_getBalance", h.handleMevCommitGetBalance)
	server.RegisterHandler("mevcommit_optInBlock", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		timeToOptIn, err := h.bidder.Estimate()
		if err != nil {
			h.logger.Error("Failed to estimate time to opt in", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to estimate opt in time",
			)
		}
		return json.RawMessage(fmt.Sprintf(`{"timeInSecs": "%d"}`, timeToOptIn)), false, nil
	})
	server.RegisterHandler("mevcommit_estimateDeposit", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		blockPrice, err := h.pricer.EstimatePrice(
			ctx,
			types.NewTransaction(0, h.depositAddress, big.NewInt(0), 21000, big.NewInt(0), nil),
		)
		if err != nil {
			h.logger.Error("Failed to estimate deposit price", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to estimate deposit price",
			)
		}
		if blockPrice == nil {
			h.logger.Warn("No block price estimated for deposit")
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"no block price available for deposit",
			)
		}
		result := map[string]interface{}{
			"bidAmount":      hexutil.EncodeBig(blockPrice.BidAmount),
			"depositAddress": h.depositAddress.Hex(),
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			h.logger.Error("Failed to marshal deposit estimate to JSON", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to marshal deposit estimate",
			)
		}
		h.logger.Debug("Estimated deposit price", "bidAmount", blockPrice.BidAmount, "depositAddress", h.depositAddress.Hex())
		return resultJSON, false, nil
	})
	server.RegisterHandler("mevcommit_estimateBridge", func(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
		blockPrice, err := h.pricer.EstimatePrice(
			ctx,
			types.NewTransaction(0, h.bridgeAddress, big.NewInt(0), 21000, big.NewInt(0), nil),
		)
		if err != nil {
			h.logger.Error("Failed to estimate bridge price", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to estimate bridge price",
			)
		}
		if blockPrice == nil {
			h.logger.Warn("No block price estimated for bridge")
			return nil, true, nil // No price available, proxy
		}
		bridgeCost := new(big.Int).Mul(blockPrice.BidAmount, big.NewInt(2))
		result := map[string]interface{}{
			"bidAmount":     hexutil.EncodeBig(bridgeCost),
			"bridgeAddress": h.bridgeAddress.Hex(),
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			h.logger.Error("Failed to marshal bridge estimate to JSON", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to marshal bridge estimate",
			)
		}
		h.logger.Debug("Estimated bridge price", "bidAmount", blockPrice.BidAmount, "bridgeAddress", h.bridgeAddress.Hex())
		return resultJSON, false, nil
	})
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
	// Remove "0x" prefix if present
	blockHashStr = strings.TrimPrefix(blockHashStr, "0x")
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

	txns, err := h.store.GetTransactionsForBlock(ctx, int64(blockNumber))
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
		if txn.Status != sender.TxStatusPreConfirmed {
			h.logger.Warn("Skipping transaction not in preconfirmed state", "txnHash", txn.Hash().Hex(), "status", txn.Status)
			continue
		}
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
		txnsToReturn = append(
			txnsToReturn.([]map[string]interface{}),
			map[string]interface{}{
				"hash":                 txn.Hash().Hex(),
				"blockHash":            blockHashStr,
				"blockNumber":          hexutil.Uint64(blockNumber),
				"transactionIndex":     hexutil.Uint64(i),
				"type":                 hexutil.Uint(txn.Transaction.Type()),
				"accessList":           nil, // Access lists are not used in preconf blocks
				"maxFeePerGas":         hexutil.EncodeBig(txn.GasFeeCap()),
				"maxPriorityFeePerGas": hexutil.EncodeBig(txn.GasTipCap()),
				"to":                   txn.To().Hex(),
				"value":                hexutil.EncodeBig(txn.Value()),
				"input":                hexutil.Encode(txn.Data()),
				"from":                 txn.Sender.Hex(),
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

	err = h.sndr.Enqueue(ctx, &sender.Transaction{
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

	txn, err := h.store.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return nil, true, nil
	}

	if txn.Status != sender.TxStatusFailed &&
		(txn.Status != sender.TxStatusPreConfirmed || h.blockTracker.LatestBlockNumber() > uint64(txn.BlockNumber)) {
		return nil, true, nil
	}

	result := map[string]interface{}{
		"type":              hexutil.Uint(txn.Transaction.Type()),
		"transactionHash":   txn.Hash().Hex(),
		"transactionIndex":  hexutil.Uint(0),
		"from":              txn.Sender.Hex(),
		"to":                nil,
		"contractAddress":   (common.Address{}).Hex(),
		"gasUsed":           hexutil.Uint64(0),
		"cumulativeGasUsed": hexutil.Uint64(1),
		"logs":              []*types.Log{}, // should be [] not null
		"logsBloom":         hexutil.Bytes(types.Bloom{}.Bytes()),
		"effectiveGasPrice": hexutil.EncodeBig(big.NewInt(0)),
	}

	if txn.Status == sender.TxStatusFailed {
		result["status"] = hexutil.Uint64(types.ReceiptStatusFailed)
	} else {
		result["status"] = hexutil.Uint64(types.ReceiptStatusSuccessful)
		blockHash := fmt.Sprintf("0x%s%08d", preconfBlockHashPrefix, txn.BlockNumber)
		blockHash += strings.Repeat("0", 66-len(blockHash))
		result["blockHash"] = blockHash
		result["blockNumber"] = hexutil.EncodeBig(big.NewInt(txn.BlockNumber))
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

	accNonce := h.store.GetCurrentNonce(ctx, common.HexToAddress(account))
	if accNonce == 0 {
		return nil, true, nil
	}

	accNonce += 1

	nonceJSON, err := json.Marshal(accNonce)
	if err != nil {
		h.logger.Error("Failed to marshal nonce to JSON", "error", err, "account", account)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal nonce",
		)
	}

	h.logger.Info("Retrieved account nonce from cache", "account", account, "nonce", accNonce)
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

	commitments, err := h.store.GetTransactionCommitments(ctx, txHash)
	if err != nil {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get transaction commitments",
		)
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
