package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// timeMilliZ is like time.StampMilli but with time zone suffix "Z".
const timeMilliZ = "2006-01-02T15:04:05.000Z"

// newLogger initializes a *slog.Logger with specified level, format, and sink.
//   - sink: destination for log output (e.g., os.Stdout, file)
//   - logLvl: string representation of slog.Level; it defaults to slog.LevelInfo if empty
//   - logFmt: format of the log output: "text", "json", defaults to "json"; it defaults to "text" if empty
//   - logTags: comma-separated list of <name:value> pairs that will be inserted into each log line
//
// Returns a configured *slog.Logger on success or nil on failure.
func newLogger(sink io.Writer, logLvl, logFmt, logTags string) (*slog.Logger, error) {
	var level slog.Leveler = slog.LevelInfo
	if logLvl != "" {
		l := new(slog.LevelVar)
		if err := l.UnmarshalText([]byte(logLvl)); err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		level = l
	}

	var (
		handler slog.Handler
		options = &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
	)
	switch logFmt {
	case "text":
		handler = slog.NewTextHandler(sink, options)
	case "json":
		handler = slog.NewJSONHandler(sink, options)
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	logger := slog.New(handler)

	if logTags == "" {
		return logger, nil
	}

	var args []any
	for i, p := range strings.Split(logTags, ",") {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag at index %d", i)
		}
		args = append(args, strings.ToValidUTF8(kv[0], "�"), strings.ToValidUTF8(kv[1], "�"))
	}

	return logger.With(args...), nil
}

// parseBlockTransactions extracts transaction data from a given block.
// Returns a slice of TransactionItem and a slice of transaction hashes.
func parseBlockTransactions(block *types.Block) ([]*TransactionItem, []common.Hash, error) {
	var (
		txs = make([]*TransactionItem, len(block.Transactions()))
		txh = make([]common.Hash, len(block.Transactions()))
	)

	for i, tx := range block.Transactions() {
		sender, err := types.Sender(types.NewCancunSigner(tx.ChainId()), tx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to derive sender: %w", err)
		}

		v, r, s := tx.RawSignatureValues()
		transaction := &TransactionItem{
			Hash:        tx.Hash().Hex(),
			From:        sender.Hex(),
			Gas:         tx.Gas(),
			Nonce:       tx.Nonce(),
			BlockHash:   block.Hash().Hex(),
			BlockNumber: block.NumberU64(),
			ChainId:     tx.ChainId().String(),
			V:           v.String(),
			R:           r.String(),
			S:           s.String(),
			Input:       hex.EncodeToString(tx.Data()),
			Timestamp:   tx.Time().UTC().Format(timeMilliZ),
		}

		if tx.To() != nil {
			transaction.To = tx.To().Hex()
		}
		if tx.GasPrice() != nil {
			transaction.GasPrice = tx.GasPrice().Uint64()
		}
		if tx.GasTipCap() != nil {
			transaction.GasTipCap = tx.GasTipCap().Uint64()
		}
		if tx.GasFeeCap() != nil {
			transaction.GasFeeCap = tx.GasFeeCap().Uint64()
		}
		if tx.Value() != nil {
			transaction.Value = tx.Value().String()
		}

		txs[i] = transaction
		txh[i] = tx.Hash()
	}

	return txs, txh, nil
}
