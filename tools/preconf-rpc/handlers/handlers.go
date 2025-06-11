package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"

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
	) (<-chan optinbidder.BidStatus, error)
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
		txnHash string,
	) (*types.Transaction, []*bidderapiv1.Commitment, error)
	DeductBalance(
		ctx context.Context,
		account string,
		amount *big.Int,
	) error
	HasBalance(
		ctx context.Context,
		account string,
		amount *big.Int,
	) bool
}

type accountNonce struct {
	Account string `json:"account"`
	Nonce   uint64 `json:"nonce"`
	Block   int64  `json:"block"`
}

type rpcMethodHandler struct {
	logger       *slog.Logger
	bidder       Bidder
	store        Store
	pricer       Pricer
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
	unlock := true
	defer func() {
		if unlock {
			h.nonceLock.Unlock(sender.Hex())
		}
	}()

	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate time to opt-in", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to estimate time to opt-in",
		)
	}

	optedInSlot := false
	if timeToOptIn <= blockTime {
		optedInSlot = true
	}

	price, err := h.pricer.EstimatePrice(ctx, txn)
	if err != nil {
		h.logger.Error("Failed to estimate transaction price", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to estimate transaction price",
		)
	}

	if !h.store.HasBalance(ctx, sender.Hex(), price.BidAmount) {
		h.logger.Error("Insufficient balance for sender", "sender", sender.Hex())
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"insufficient balance for sender",
		)
	}

	bidC, err := h.bidder.Bid(
		ctx,
		price.BidAmount,
		big.NewInt(0),
		rawTxHex,
		&optinbidder.BidOpts{
			WaitForOptIn: optedInSlot,
			BlockNumber:  uint64(price.BlockNumber),
		},
	)
	if err != nil {
		h.logger.Error("Failed to place bid", "error", err)
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeCustomError, "failed to place bid")
	}
	noOfProviders, noOfCommitments, blockNumber := 0, 0, 0
	cancelled, failed := false, false
	commitments := make([]*bidderapiv1.Commitment, 0)
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Context cancelled while waiting for bid status")
			return nil, false, ctx.Err()
		case bidStatus := <-bidC:
			switch bidStatus.Type {
			case optinbidder.BidStatusNoOfProviders:
				noOfProviders = bidStatus.Arg.(int)
			case optinbidder.BidStatusAttempted:
				blockNumber = bidStatus.Arg.(int)
			case optinbidder.BidStatusCommitment:
				noOfCommitments++
				commitments = append(commitments, bidStatus.Arg.(*bidderapiv1.Commitment))
			case optinbidder.BidStatusCancelled:
				cancelled = true
				break BID_LOOP
			case optinbidder.BidStatusFailed:
				failed = true
				break BID_LOOP
			}
		}
	}
	switch {
	case noOfProviders == 0:
		h.logger.Info("No providers available for the bid", "noOfProviders", noOfProviders)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"no providers available for the bid",
		)
	case cancelled || failed:
		h.logger.Info(
			"Bid cancelled or failed",
			"cancelled", cancelled,
			"failed", failed,
			"noOfCommitments", noOfCommitments,
		)
		if noOfCommitments == 0 {
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"bid cancelled with no commitments",
			)
		}
	case noOfCommitments == 0:
		h.logger.Info("Bid completed with no commitments", "noOfCommitments", noOfCommitments)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"bid completed with no commitments",
		)
	}
	h.logger.Info(
		"Bid successful with commitments",
		"noOfProviders", noOfProviders,
		"noOfCommitments", noOfCommitments,
		"blockNumber", blockNumber,
	)
	// If we reach here, we have a successful bid with commitments
	txHashJSON, err := json.Marshal(txn.Hash().Hex())
	if err != nil {
		h.logger.Error("Failed to marshal transaction hash to JSON", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal transaction hash",
		)
	}

	if err := h.store.DeductBalance(ctx, sender.Hex(), price.BidAmount); err != nil {
		h.logger.Error("Failed to deduct balance for sender", "sender", sender.Hex(), "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to deduct balance for sender",
		)
	}

	if err := h.store.StorePreconfirmedTransaction(ctx, int64(blockNumber), txn, commitments); err != nil {
		h.logger.Error("Failed to store preconfirmed transaction", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to store preconfirmed transaction",
		)
	}

	if noOfProviders == noOfCommitments && optedInSlot {
		h.logger.Info("All providers committed, updating nonce", "account", sender.Hex(), "nonce", txn.Nonce()+1)
		h.nonceMapLock.Lock()
		h.nonceMap[sender.Hex()] = accountNonce{
			Account: sender.Hex(),
			Nonce:   txn.Nonce() + 1,
			Block:   int64(blockNumber),
		}
		h.nonceMapLock.Unlock()
	} else {
		unlock = false
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

	txHash := params[0].(string)
	if len(txHash) < 2 || txHash[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeParseError,
			"getTxReceipt parameter must be a hex string starting with '0x'",
		)
	}

	h.logger.Info("Retrieving transaction receipt", "txHash", txHash)
	txn, commitments, err := h.store.GetPreconfirmedTransaction(ctx, txHash[2:])
	if err != nil {
		return nil, true, nil
	}

	receipt := &types.Receipt{
		TxHash:      txn.Hash(),
		Type:        txn.Type(),
		Status:      1, // Assuming success, as this is a preconfirmed transaction
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
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeInvalidRequest,
			"getTxCount requires exactly one parameter",
		)
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
