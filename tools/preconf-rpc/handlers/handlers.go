package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"resenje.org/multex"
)

const (
	blockTime = 12 // seconds, typical Ethereum block time
)

type Bidder interface {
	Estimate() (int64, error)
	Bid(ctx context.Context, bidAmount *big.Int, slashAmount *big.Int, rawTx string) (<-chan optinbidder.BidStatus, error)
}

type Store interface {
	StorePreconfirmedTransaction(
		ctx context.Context,
		blockNumber int64,
		txn *types.Transaction,
		commitments []*bidderapiv1.Commitment,
	) error
	UpdateAccountNonce(
		ctx context.Context,
		account string,
		nonce uint64,
	) error
	GetAccountNonce(
		ctx context.Context,
		account string,
	) (uint64, error)
}

type rpcMethodHandler struct {
	logger    *slog.Logger
	bidder    Bidder
	store     Store
	nonceLock *multex.Multex[string]
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

	bidC, err := h.bidder.Bid(ctx, big.NewInt(0), big.NewInt(0), rawTxHex)
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

	if err := h.store.StorePreconfirmedTransaction(ctx, int64(blockNumber), txn, commitments); err != nil {
		h.logger.Error("Failed to store preconfirmed transaction", "error", err)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to store preconfirmed transaction",
		)
	}

	if noOfProviders == noOfCommitments && optedInSlot {
		if err := h.store.UpdateAccountNonce(ctx, sender.Hex(), txn.Nonce()+1); err != nil {
			h.logger.Error("Failed to update account nonce", "error", err)
			return nil, false, rpcserver.NewJSONErr(
				rpcserver.CodeCustomError,
				"failed to update account nonce",
			)
		}
	} else {
		h.nonceLock.Lock(sender.Hex())
	}

	return txHashJSON, false, nil
}

func (h *rpcMethodHandler) handleGetTxReceipt(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	return nil, false, nil
}

func (h *rpcMethodHandler) handleGetBalance(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	return nil, false, nil
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
	nonce, err := h.store.GetAccountNonce(ctx, account)
	if err != nil {
		h.logger.Error("Failed to get account nonce", "error", err, "account", account)
		return nil, true, nil
	}
	nonceJSON, err := json.Marshal(nonce)
	if err != nil {
		h.logger.Error("Failed to marshal nonce to JSON", "error", err, "account", account)
		return nil, false, rpcserver.NewJSONErr(
			rpcserver.CodeCustomError,
			"failed to marshal nonce",
		)
	}
	h.logger.Info("Retrieved account nonce", "account", account, "nonce", nonce)
	return nonceJSON, false, nil
}
