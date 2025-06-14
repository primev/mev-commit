package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
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
	DeductBalance(
		ctx context.Context,
		account common.Address,
		amount *big.Int,
	) error
	HasBalance(
		ctx context.Context,
		account common.Address,
		amount *big.Int,
	) bool
}

type BlockTracker interface {
	CheckTxnInclusion(
		ctx context.Context,
		txnHash common.Hash,
		blockNumber uint64,
	) (bool, error)
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
	nonceLock    *multex.Multex[string]
	latestBlock  atomic.Pointer[types.Block]
	nonceMap     map[string]accountNonce
	nonceMapLock sync.RWMutex
}

func NewRPCMethodHandler(
	logger *slog.Logger,
	bidder Bidder,
	store Store,
	pricer Pricer,
) *rpcMethodHandler {
	return &rpcMethodHandler{
		logger:      logger,
		bidder:      bidder,
		store:       store,
		pricer:      pricer,
		nonceLock:   multex.New[string](),
		nonceMap:    make(map[string]accountNonce),
		latestBlock: atomic.Pointer[types.Block]{},
	}
}

func (h *rpcMethodHandler) RegisterMethods(server *rpcserver.JSONRPCServer) {
	server.RegisterHandler("eth_sendRawTransaction", h.handleSendRawTx)
	server.RegisterHandler("eth_getTransactionReceipt", h.handleGetTxReceipt)
	server.RegisterHandler("eth_getTransactionCount", h.handleGetTxCount)
	server.RegisterHandler("mevcommit_getTransactionCommitments", h.handleGetTxCommitments)
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

		result, err := h.sendBid(ctx, txn, sender, rawTxHex)
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
) (bidResult, error) {
	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate time to opt-in", "error", err)
		if !errors.Is(err, optinbidder.ErrNoSlotInCurrentEpoch) {
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

	if !h.store.HasBalance(ctx, sender, price.BidAmount) {
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
) error {
	if err := h.store.StorePreconfirmedTransaction(ctx, blockNumber, txn, commitments); err != nil {
		h.logger.Error("Failed to store preconfirmed transaction", "error", err)
		return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
	}

	if err := h.store.DeductBalance(ctx, sender, amount); err != nil {
		h.logger.Error("Failed to deduct balance for sender", "sender", sender.Hex(), "error", err)
		return fmt.Errorf("failed to deduct balance for sender: %w", err)
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

	receipt := &types.Receipt{
		TxHash:      txn.Hash(),
		Type:        txn.Type(),
		Status:      types.ReceiptStatusSuccessful, // Assuming success, as this is a preconfirmed transaction
		BlockNumber: big.NewInt(commitments[0].BlockNumber),
	}

	receiptJSON, err := json.Marshal(receipt)
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

	if h.latestBlock.Load().Number().Uint64() > uint64(accNonce.Block) {
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
