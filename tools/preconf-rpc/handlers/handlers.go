package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"resenje.org/multex"
)

const (
	blockTime = 12 // seconds, typical Ethereum block time
)

var (
	preconfBlockHashPrefix    = hex.EncodeToString([]byte("mev-commit"))
	preconfBlockHashPrefixLen = len(preconfBlockHashPrefix)
)

type Bidder interface {
	Estimate() (int64, error)
	Bid(
		ctx context.Context,
		bidAmount *big.Int,
		slashAmount *big.Int,
		rawTx string,
		opts *optinbidder.BidOpts,
	) (chan optinbidder.BidStatus, error)
}

type Pricer interface {
	EstimatePrice(
		ctx context.Context,
		txn *types.Transaction,
	) (*pricer.BlockPrice, error)
}

type Store interface {
	StorePreconfirmedTransaction(
		ctx context.Context,
		blockNumber int64,
		txn *types.Transaction,
		commitments []*bidderapiv1.Commitment,
	) error
	GetPreconfirmedTransaction(
		ctx context.Context,
		txnHash common.Hash,
	) (*types.Transaction, []*bidderapiv1.Commitment, error)
	GetPreconfirmedTransactionsForBlock(
		ctx context.Context,
		blockNumber int64,
	) ([]*types.Transaction, error)
	DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error
	HasBalance(ctx context.Context, account common.Address, amount *big.Int) bool
	GetBalance(ctx context.Context, account common.Address) (*big.Int, error)
	AddBalance(ctx context.Context, account common.Address, amount *big.Int) error
}

type BlockTracker interface {
	CheckTxnInclusion(ctx context.Context, txnHash common.Hash, blockNumber uint64) (bool, error)
	LatestBlockNumber() uint64
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
	logger       *slog.Logger
	bidder       Bidder
	store        Store
	pricer       Pricer
	blockTracker BlockTracker
	owner        common.Address
	chainID      *big.Int
	nonceLock    *multex.Multex[string]
	nonceMap     map[string]accountNonce
	nonceMapLock sync.RWMutex
}

func NewRPCMethodHandler(
	logger *slog.Logger,
	bidder Bidder,
	store Store,
	pricer Pricer,
	blockTracker BlockTracker,
	owner common.Address,
) *rpcMethodHandler {
	return &rpcMethodHandler{
		logger:       logger,
		bidder:       bidder,
		store:        store,
		pricer:       pricer,
		blockTracker: blockTracker,
		owner:        owner,
		nonceLock:    multex.New[string](),
		nonceMap:     make(map[string]accountNonce),
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
	server.RegisterHandler("mevcommit_estimateFastBid", h.handleMevCommitEstimateFastBid)
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
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getBlock requires exactly one parameter",
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
		"transactions":     make([]string, len(txns)),
		"transactionsRoot": (common.Hash{}).Hex(),
		"stateRoot":        (common.Hash{}).Hex(),
		"miner":            h.owner.Hex(),
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

	for i, txn := range txns {
		block["transactions"].([]string)[i] = txn.Hash().Hex()
	}
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

	sender, err := types.Sender(types.LatestSignerForChainID(txn.ChainId()), txn)
	if err != nil {
		h.logger.Error("Failed to get transaction sender", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to get transaction sender",
		)
	}

	// Once we are ready to send the bid, we need to ensure that the nonce for the
	// sender is not locked by another transaction.
	h.nonceLock.Lock(sender.Hex())
	defer h.nonceLock.Unlock(sender.Hex())

	// This is a txn to add balance to the bidder's account, so we will pay this
	// out of the owner's account. We will add the balance to the bidder's
	// account and then proceed with the bid process.
	depositTxn := txn.To().Cmp(h.owner) == 0 && txn.Value().Cmp(big.NewInt(0)) > 0

BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"context cancelled while processing transaction",
			)
		default:
		}

		result, err := h.sendBid(ctx, txn, sender, rawTxHex, depositTxn)
		switch {
		case err != nil:
			h.logger.Error("Failed to send bid", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to send bid",
			)
		case result.optedInSlot:
			if result.noOfProviders == len(result.commitments) {
				// This means that all builders have committed to the bid and it
				// is a primev opted in slot. We can safely proceed to inform the
				// user that the txn was successfully sent and will be processed
				if err := h.storePreconfAndDeductBalance(
					ctx,
					txn,
					result.commitments,
					sender,
					int64(result.blockNumber),
					result.bidAmount,
					depositTxn,
				); err != nil {
					return nil, false, rpcserver.NewJSONErr(
						rpcserver.CodeCustomError,
						"failed to update preconfirmed transaction and deduct balance",
					)
				}
				// Update the nonce locally if user wants to send more transactions
				h.nonceMapLock.Lock()
				h.nonceMap[sender.Hex()] = accountNonce{
					Account: sender.Hex(),
					Nonce:   txn.Nonce() + 1,
					Block:   int64(result.blockNumber),
				}
				h.nonceMapLock.Unlock()
				break BID_LOOP
			}
		default:
		}

		// Wait for block number to be updated to confirm transaction. If failed
		// we will retry the bid process till user cancels the operation
		included, err := h.blockTracker.CheckTxnInclusion(ctx, txn.Hash(), result.blockNumber)
		if err != nil {
			h.logger.Error("Failed to check transaction inclusion", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to check transaction inclusion",
			)
		}
		if included {
			if err := h.storePreconfAndDeductBalance(
				ctx,
				txn,
				result.commitments,
				sender,
				int64(result.blockNumber),
				result.bidAmount,
				depositTxn,
			); err != nil {
				h.logger.Error("Failed to update preconfirmed transaction and deduct balance", "error", err)
				return nil, false, rpcserver.NewJSONErr(
					rpcserver.CodeCustomError,
					"failed to update preconfirmed transaction and deduct balance",
				)
			}
			break BID_LOOP
		}
	}

	// If we reach here, we have a successful bid with commitments
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

func (h *rpcMethodHandler) sendBid(
	ctx context.Context,
	txn *types.Transaction,
	sender common.Address,
	rawTxHex string,
	depositTxn bool,
) (bidResult, error) {
	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate time to opt-in", "error", err)
		if !errors.Is(err, optinbidder.ErrNoSlotInCurrentEpoch) && !errors.Is(err, optinbidder.ErrNoEpochInfo) {
			return bidResult{}, err
		}
		// If we cannot estimate the time to opt-in, we assume a default value and
		// proceed with the bid process. The default value should be higher than
		// the typical block time to ensure we consider the next slot as a non-opt-in slot.
		timeToOptIn = blockTime * 32
	}

	optedInSlot := timeToOptIn <= blockTime

	price, err := h.pricer.EstimatePrice(ctx, txn)
	if err != nil {
		h.logger.Error("Failed to estimate transaction price", "error", err)
		return bidResult{}, fmt.Errorf("failed to estimate transaction price: %w", err)
	}

	if !depositTxn && !h.store.HasBalance(ctx, sender, price.BidAmount) {
		h.logger.Error("Insufficient balance for sender", "sender", sender.Hex())
		return bidResult{}, fmt.Errorf("insufficient balance for sender: %s", sender.Hex())
	}

	bidC, err := h.bidder.Bid(
		ctx,
		price.BidAmount,
		big.NewInt(0),
		rawTxHex[2:],
		&optinbidder.BidOpts{
			WaitForOptIn: optedInSlot,
			// BlockNumber:  uint64(price.BlockNumber),
		},
	)
	if err != nil {
		h.logger.Error("Failed to place bid", "error", err)
		return bidResult{}, fmt.Errorf("failed to place bid: %w", err)
	}

	result := bidResult{
		commitments: make([]*bidderapiv1.Commitment, 0),
		bidAmount:   price.BidAmount,
	}
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Context cancelled while waiting for bid status")
			return bidResult{}, ctx.Err()
		case bidStatus, more := <-bidC:
			if !more {
				h.logger.Info("Bid channel closed, no more bid statuses")
				break BID_LOOP
			}
			switch bidStatus.Type {
			case optinbidder.BidStatusNoOfProviders:
				result.noOfProviders = bidStatus.Arg.(int)
			case optinbidder.BidStatusAttempted:
				result.blockNumber = bidStatus.Arg.(uint64)
			case optinbidder.BidStatusCommitment:
				result.commitments = append(result.commitments, bidStatus.Arg.(*bidderapiv1.Commitment))
			case optinbidder.BidStatusCancelled:
				h.logger.Warn("Bid context cancelled by the bidder")
				break BID_LOOP
			case optinbidder.BidStatusFailed:
				h.logger.Error("Bid failed", "error", bidStatus.Arg)
				break BID_LOOP
			}
		}
	}
	if len(result.commitments) == 0 {
		h.logger.Error("Bid completed with no commitments")
		return bidResult{}, fmt.Errorf("bid completed with no commitments")
	}
	h.logger.Info(
		"Bid successful with commitments",
		"noOfProviders", result.noOfProviders,
		"noOfCommitments", len(result.commitments),
		"blockNumber", result.blockNumber,
		"optedInSlot", optedInSlot,
	)

	result.optedInSlot = optedInSlot
	return result, nil
}

func (h *rpcMethodHandler) storePreconfAndDeductBalance(
	ctx context.Context,
	txn *types.Transaction,
	commitments []*bidderapiv1.Commitment,
	sender common.Address,
	blockNumber int64,
	amount *big.Int,
	depositTxn bool,
) error {
	if err := h.store.StorePreconfirmedTransaction(ctx, blockNumber, txn, commitments); err != nil {
		h.logger.Error("Failed to store preconfirmed transaction", "error", err)
		return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
	}

	if !depositTxn {
		if err := h.store.DeductBalance(ctx, sender, amount); err != nil {
			h.logger.Error("Failed to deduct balance for sender", "sender", sender.Hex(), "error", err)
			return fmt.Errorf("failed to deduct balance for sender: %w", err)
		}
	} else {
		if err := h.store.AddBalance(ctx, sender, txn.Value()); err != nil {
			h.logger.Error("Failed to add balance for sender", "sender", sender.Hex(), "error", err)
			return fmt.Errorf("failed to add balance for sender: %w", err)
		}
	}

	return nil
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

func (h *rpcMethodHandler) handleMevCommitEstimateFastBid(
	ctx context.Context,
	_ ...any,
) (json.RawMessage, bool, error) {
	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate fast bid", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to estimate fast bid",
		)
	}

	return json.RawMessage(fmt.Sprintf(`{"timeInSecs": "%d"}`, timeToOptIn)), false, nil
}
