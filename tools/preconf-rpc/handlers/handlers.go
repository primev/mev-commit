package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
)

const (
	blockTime = 12 // seconds, typical Ethereum block time
)

type Bidder interface {
	Estimate() (int64, error)
	Bid(ctx context.Context, bidAmount *big.Int, slashAmount *big.Int, rawTx string) (<-chan optinbidder.BidStatus, error)
}

type rpcMethodHandler struct {
	logger *slog.Logger
	bidder Bidder
}

func (h *rpcMethodHandler) handleSendRawTx(ctx context.Context, params ...any) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeInvalidRequest, "sendRawTx requires exactly one parameter")
	}
	if params[0] == nil {
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeParseError, "sendRawTx parameter cannot be null")
	}

	rawTxHex := params[0].(string)
	if len(rawTxHex) < 2 || rawTxHex[:2] != "0x" {
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeParseError, "sendRawTx parameter must be a hex string starting with '0x'")
	}

	decodedTxn, err := hex.DecodeString(rawTxHex[2:])
	if err != nil {
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeParseError, "sendRawTx parameter must be a valid hex string")
	}

	txn := new(types.Transaction)
	if err := txn.UnmarshalBinary(decodedTxn); err != nil {
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeParseError, "sendRawTx parameter must be a valid transaction")
	}

	timeToOptIn, err := h.bidder.Estimate()
	if err != nil {
		h.logger.Error("Failed to estimate time to opt-in", "error", err)
		return nil, false, rpcserver.NewJSONErr(rpcserver.CodeCustomError, "failed to estimate time to opt-in")
	}

	if timeToOptIn <= blockTime {
		bidC, err := h.bidder.Bid(ctx, big.NewInt(0), big.NewInt(0), rawTxHex)
		if err != nil {
			h.logger.Error("Failed to place bid", "error", err)
			return nil, false, rpcserver.NewJSONErr(rpcserver.CodeCustomError, "failed to place bid")
		}
		noOfProviders, noOfCommitments, blockNumber := 0, 0, 0
		commitments := make([]*bidderapiv1.Commitment, 0)
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
					return nil, false, rpcserver.NewJSONErr(rpcserver.CodeCustomError, "bid cancelled")
				case optinbidder.BidStatusFailed:
					return nil, false, rpcserver.NewJSONErr(rpcserver.CodeCustomError, "bid failed")
				}
			}
		}
	}

	return nil, false, nil
}

func (h *rpcMethodHandler) handleGetTxReceipt(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	return nil, false, nil
}

func (h *rpcMethodHandler) handleGetBalance(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	return nil, false, nil
}

func (h *rpcMethodHandler) handleGetTxCount(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	return nil, false, nil
}
