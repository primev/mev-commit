package util

import (
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func PadKeyTo32Bytes(key *big.Int) []byte {
	keyBytes := key.Bytes()
	if len(keyBytes) < 32 {
		padding := make([]byte, 32-len(keyBytes))
		keyBytes = append(padding, keyBytes...)
	}
	return keyBytes
}

func NewTestLogger(w io.Writer) *slog.Logger {
	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

// NewLogger initializes a *slog.Logger with specified level, format, and sink.
//   - lvl: string representation of slog.Level
//   - logFmt: format of the log output: "text", "json", "none" defaults to "json"
//   - tags: comma-separated list of <name:value> pairs that will be inserted into each log line
//   - sink: destination for log output (e.g., os.Stdout, file)
//
// Returns a configured *slog.Logger on success or nil on failure.
func NewLogger(lvl, logFmt, tags string, sink io.Writer) (*slog.Logger, error) {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(lvl)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
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
	case "json", "none":
		handler = slog.NewJSONHandler(sink, options)
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	logger := slog.New(handler)

	if tags == "" {
		return logger, nil
	}

	var args []any
	for i, p := range strings.Split(tags, ",") {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag at index %d", i)
		}
		args = append(args, strings.ToValidUTF8(kv[0], "�"), strings.ToValidUTF8(kv[1], "�"))
	}

	return logger.With(args...), nil
}

func ResolveFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}

type Kind string

const (
	KindETHSend          Kind = "ETH Send"
	KindSelfTransfer     Kind = "Self ETH Transfer"
	KindContractCreation Kind = "Contract Creation"
	KindContractCall     Kind = "Contract Call"
	KindERC20Transfer    Kind = "ERC20 Transfer"
	KindERC20Approve     Kind = "ERC20/721 Approval"
	KindERC20Permit      Kind = "ERC20 Permit (EIP-2612)"
	KindERC721Transfer   Kind = "ERC721 Transfer"
	KindERC721ApproveAll Kind = "ERC721/1155 SetApprovalForAll"
	KindERC1155Transfer  Kind = "ERC1155 TransferSingle"
	KindERC1155Batch     Kind = "ERC1155 TransferBatch"
	KindUnknown          Kind = "Unknown"
)

type Classification struct {
	Kind    Kind
	Details string // small hint (e.g. function name)
}

// Common selectors (first 4 bytes of calldata)
var (
	selERC20Transfer     = [4]byte{0xa9, 0x05, 0x9c, 0xbb} // transfer(address,uint256)
	selERC20Approve      = [4]byte{0x09, 0x5e, 0xa7, 0xb3} // approve(address,uint256)  (same as ERC721 approve(address,uint256))
	selERC20TransferFrom = [4]byte{0x23, 0xb8, 0x72, 0xdd} // transferFrom(address,address,uint256)
	selPermit2612        = [4]byte{0xd5, 0x05, 0xac, 0xcf} // permit(address,address,uint256,uint256,uint8,bytes32,bytes32)
	sel721SafeTx3        = [4]byte{0x42, 0x84, 0x2e, 0x0e} // safeTransferFrom(address,address,uint256)
	sel721SafeTx4        = [4]byte{0xb8, 0x8d, 0x4f, 0xde} // safeTransferFrom(address,address,uint256,bytes)
	selSetApprovalForAll = [4]byte{0xa2, 0x2c, 0xb4, 0x65} // setApprovalForAll(address,bool)
	sel1155SafeTx        = [4]byte{0xf2, 0x42, 0x43, 0x2a} // safeTransferFrom(address,address,uint256,uint256,bytes)
	sel1155SafeBatchTx   = [4]byte{0x2e, 0xb2, 0xc2, 0xd6} // safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)
)

// ClassifyTxOnly uses only tx envelope fields (no receipt/logs).
// `from` is required because go-ethereum’s Transaction doesn’t carry it.
func ClassifyTxOnly(tx *types.Transaction, from common.Address) Classification {
	to := tx.To()
	data := tx.Data()
	val := tx.Value()

	// Creation vs call vs plain ETH send
	if to == nil {
		return Classification{Kind: KindContractCreation, Details: "contract deploy"}
	}
	if len(data) == 0 {
		if val != nil && val.Sign() > 0 {
			if from == *to {
				return Classification{Kind: KindSelfTransfer, Details: fmt.Sprintf("%s wei", val.String())}
			}
			return Classification{Kind: KindETHSend, Details: fmt.Sprintf("to %s, %s wei", to.Hex(), val.String())}
		}
		// empty data + zero value ⇒ likely “noop” call (rare) or account-base tx
		return Classification{Kind: KindContractCall, Details: "empty calldata"}
	}

	// Function selector heuristics
	if len(data) >= 4 {
		var sel [4]byte
		copy(sel[:], data[:4])
		switch sel {
		case selERC20Transfer:
			return Classification{Kind: KindERC20Transfer, Details: "transfer(address,uint256)"}
		case selERC20Approve:
			return Classification{Kind: KindERC20Approve, Details: "approve(...)"} // could be ERC721 approve; indistinguishable without logs/ABI
		case selERC20TransferFrom:
			// Could be ERC20 transferFrom; ERC721 transferFrom has same selector (0x23b872dd).
			return Classification{Kind: KindERC721Transfer, Details: "transferFrom(...)"} // label generically as transferFrom
		case selPermit2612:
			return Classification{Kind: KindERC20Permit, Details: "permit(...) (EIP-2612)"}
		case sel721SafeTx3, sel721SafeTx4:
			return Classification{Kind: KindERC721Transfer, Details: "safeTransferFrom"}
		case selSetApprovalForAll:
			return Classification{Kind: KindERC721ApproveAll, Details: "setApprovalForAll"}
		case sel1155SafeTx:
			return Classification{Kind: KindERC1155Transfer, Details: "ERC1155 safeTransferFrom"}
		case sel1155SafeBatchTx:
			return Classification{Kind: KindERC1155Batch, Details: "ERC1155 safeBatchTransferFrom"}
		}
	}

	// Fallback
	return Classification{Kind: KindContractCall, Details: "call with calldata"}
}
