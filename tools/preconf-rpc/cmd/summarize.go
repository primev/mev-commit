package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"google.golang.org/protobuf/proto"
)

// TransactionRow mirrors the columns selected in our query.
type TransactionRow struct {
	Hash          string
	Nonce         int64
	RawTx         string
	BlockNumber   int64
	Sender        string
	TxType        int64
	Status        string
	Details       sql.NullString
	Options       []byte
	CommitDigest  sql.NullString
	CommitData    []byte
	SimulationLog sql.NullString
}

// SimulationLog represents the JSON structure stored in simulationLogs.logs.
// Adjust these fields if your JSON is different.
type SimulationLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

// CommitmentDBRow holds commitment data for a single row.
type CommitmentDBRow struct {
	Digest string
	Data   []byte
}

// TxAgg aggregates a transaction with all its commitments.
type TxAgg struct {
	Tx          TransactionRow
	Commitments []CommitmentDBRow
}

func main() {
	// DB URL: e.g.
	// export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	log.Println("Connected to Postgres")

	rows, err := queryTransactionsWithCommitmentsAndLogs(ctx, db)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	// First: aggregate rows into unique txs with multiple commitments
	txMap := make(map[string]*TxAgg)

	for rows.Next() {
		var r TransactionRow
		err := rows.Scan(
			&r.Hash,
			&r.Nonce,
			&r.RawTx,
			&r.BlockNumber,
			&r.Sender,
			&r.TxType,
			&r.Status,
			&r.Details,
			&r.Options,
			&r.CommitDigest,
			&r.CommitData,
			&r.SimulationLog,
		)
		if err != nil {
			log.Fatalf("scan failed: %v", err)
		}

		agg, ok := txMap[r.Hash]
		if !ok {
			// First time we see this transaction
			agg = &TxAgg{
				Tx:          r,
				Commitments: nil,
			}
			txMap[r.Hash] = agg
		}

		// Attach this commitment (if there is one) to the aggregated tx
		if r.CommitDigest.Valid && len(r.CommitData) > 0 {
			agg.Commitments = append(agg.Commitments, CommitmentDBRow{
				Digest: r.CommitDigest.String,
				Data:   r.CommitData,
			})
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("row iteration error: %v", err)
	}

	// Now: classify per unique tx and unmarshal all its commitments
	var (
		swapCount        int
		ethTransferCount int
		otherCount       int
		total            int
		totalCommitments int
		latencyTotal     int
	)

	for _, agg := range txMap {
		total++

		// classify using the tx & logs
		switch classifyTransaction(agg.Tx) {
		case "swap":
			swapCount++
		case "eth_transfer":
			ethTransferCount++
		default:
			otherCount++
		}

		// unmarshal *all* commitments for this tx
		for _, c := range agg.Commitments {
			commitMsg, err := unmarshalCommitment(c.Data)
			if err != nil {
				log.Printf("failed to unmarshal commitment for tx %s (digest %s): %v",
					agg.Tx.Hash, c.Digest, err)
				continue
			}
			latency := checkPreconfLatency(commitMsg)
			latencyTotal += int(latency)
			totalCommitments++
		}
	}

	fmt.Println("=== Transaction Type Counts (status: confirmed / pre-confirmed) ===")
	fmt.Printf("Total unique txs: %d\n", total)
	fmt.Printf("Swap:             %d\n", swapCount)
	fmt.Printf("ETH transfer:     %d\n", ethTransferCount)
	fmt.Printf("Others:           %d\n", otherCount)
	fmt.Printf("Total commitments processed: %d\n", totalCommitments)
	if totalCommitments > 0 {
		avgLatency := float64(latencyTotal) / float64(totalCommitments)
		fmt.Printf("Average pre-confirmation latency: %.2f ms\n", avgLatency)
	} else {
		fmt.Println("No commitments processed; cannot compute average latency.")
	}
}

// queryTransactionsWithCommitmentsAndLogs fetches all confirmed / pre-confirmed txs,
// with optional commitments and simulation logs.
func queryTransactionsWithCommitmentsAndLogs(ctx context.Context, db *sql.DB) (*sql.Rows, error) {
	const q = `
SELECT
    t.hash,
    t.nonce,
    t.raw_transaction,
    t.block_number,
    t.sender,
    t.tx_type,
    t.status,
    t.details,
    t.options,
    c.commitment_digest,
    c.commitment_data,
    s.logs
FROM mcTransactions t
LEFT JOIN commitments c
    ON c.transaction_hash = t.hash
LEFT JOIN simulationLogs s
    ON s.transaction_hash = t.hash
WHERE t.status IN ('confirmed', 'pre-confirmed');
`
	return db.QueryContext(ctx, q)
}

// classifyTransaction decides whether a transaction is a swap, ETH transfer, or other.
//
// Heuristics:
// 1. Decode raw_transaction into an Ethereum types.Transaction.
// 2. If (data == empty) and (value > 0) → "eth_transfer".
// 3. Parse simulation logs and look for known DEX Swap events → "swap".
// 4. Otherwise → "other".
func classifyTransaction(r TransactionRow) string {
	// 1. Try to decode the raw transaction
	tx, err := decodeRawTx(r.RawTx)
	if err != nil {
		log.Printf("failed to decode raw tx %s: %v", r.Hash, err)
	} else {
		// ETH transfer: to != nil, non-zero value, no calldata
		if isSimpleEthTransfer(tx) {
			return "eth_transfer"
		}
	}

	// 2. If logs exist, try to detect swap based on event signatures
	logs := parseSimulationLogs(r.SimulationLog)
	if containsSwapEvent(logs) {
		return "swap"
	}

	// 3. Fallback
	return "other"
}

// decodeRawTx decodes a hex-encoded signed Ethereum transaction into *types.Transaction.
func decodeRawTx(raw string) (*types.Transaction, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "0x")

	if raw == "" {
		return nil, fmt.Errorf("empty raw transaction")
	}

	data, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}

	var tx types.Transaction
	if err := tx.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("unmarshal tx: %w", err)
	}
	return &tx, nil
}

// isSimpleEthTransfer returns true for "plain" ETH transfers:
//   - not a contract creation (To() != nil)
//   - no input data
//   - non-zero value
func isSimpleEthTransfer(tx *types.Transaction) bool {
	if tx == nil {
		return false
	}
	if tx.To() == nil {
		// contract creation
		return false
	}
	if len(tx.Data()) != 0 {
		// has calldata → not a simple ETH transfer
		return false
	}
	if tx.Value() == nil || tx.Value().Sign() <= 0 {
		// zero or nil value
		return false
	}
	return true
}

// parseSimulationLogs parses the JSON logs stored in simulationLogs.logs.
// If the field is NULL or invalid, it returns an empty slice.
func parseSimulationLogs(logStr sql.NullString) []SimulationLog {
	if !logStr.Valid || strings.TrimSpace(logStr.String) == "" {
		return nil
	}

	var logs []SimulationLog
	if err := json.Unmarshal([]byte(logStr.String), &logs); err != nil {
		log.Printf("failed to unmarshal simulation logs: %v", err)
		return nil
	}
	return logs
}

const (
	// keccak256("Swap(address,uint256,uint256,uint256,uint256,address)")
	uniswapV2SwapEventSig = "d78ad95fa46c994b6551d0da85fc275fe6131b74a83c8bbd3a27d6e5f8c3d7e1"

	// keccak256("Swap(address,address,int256,int256,uint160,uint128,int24)")
	uniswapV3SwapEventSig = "c42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"
)

// containsSwapEvent checks if any log has a known DEX Swap event signature.
func containsSwapEvent(logs []SimulationLog) bool {
	for _, lg := range logs {
		if len(lg.Topics) == 0 {
			continue
		}
		t0 := strings.TrimPrefix(strings.ToLower(lg.Topics[0]), "0x")
		switch t0 {
		case uniswapV2SwapEventSig, uniswapV3SwapEventSig:
			return true
		}
	}
	return false
}

// unmarshalCommitment converts the BYTEA `commitment_data` into a protobuf struct.
func unmarshalCommitment(data []byte) (*bidderapiv1.Commitment, error) {
	var msg bidderapiv1.Commitment
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func checkPreconfLatency(c *bidderapiv1.Commitment) int {
	actualStart := c.DecayStartTimestamp - 200
	if c.DispatchTimestamp < actualStart {
		return 0
	}
	return int(c.DispatchTimestamp - actualStart)
}
