package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
)

type rpcMethodHandler struct {
	logger *slog.Logger
}

func (h *rpcMethodHandler) handleSendRawTx(ctx context.Context, params ...json.RawMessage) (json.RawMessage, bool, error) {
	if len(params) != 1 {
		return nil, false, &jsonRPCError{
			Code:    CodeInvalidRequest,
			Message: "sendRawTx requires exactly one parameter",
		}
	}
	if params[0] == nil {
		return nil, false, &jsonRPCError{
			Code:    CodeParseError,
			Message: "sendRawTx parameter cannot be null",
		}
	}

	rawTxHex := params[0].(string)
	if len(rawTxHex) < 2 || rawTxHex[:2] != "0x" {
		return nil, false, &jsonRPCError{
			Code:    CodeParseError,
			Message: "sendRawTx parameter must be a hex string starting with '0x'",
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
