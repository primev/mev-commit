package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
)

var settlementType = `
DO $$ BEGIN
    CREATE TYPE settlement_type AS ENUM ('reward', 'slash');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;`

var settlementsTable = `
CREATE TABLE IF NOT EXISTS settlements (
	commitment_index TEXT PRIMARY KEY,
	transaction TEXT,
	block_number BIGINT,
	builder_address TEXT,
	type settlement_type,
	amount NUMERIC(24, 0),
	bid_id TEXT,
	chainhash TEXT,
	nonce BIGINT,
	settled BOOLEAN,
	decay_percentage BIGINT,
	options TEXT
);`

var winnersTable = `
CREATE TABLE IF NOT EXISTS winners (
	block_number BIGINT PRIMARY KEY,
	builder_address TEXT
);`

var transactionsTable = `
CREATE TABLE IF NOT EXISTS sent_transactions (
	hash TEXT PRIMARY KEY,
	nonce BIGINT,
	settled BOOLEAN,
	status TEXT
);`

var integerTable = `
CREATE TABLE IF NOT EXISTS integers (
	key TEXT PRIMARY KEY,
	value BIGINT
);`

var ErrNotFound = fmt.Errorf("not found")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	for _, table := range []string{
		settlementType,
		settlementsTable,
		winnersTable,
		transactionsTable,
		integerTable,
	} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) RegisterWinner(
	ctx context.Context,
	blockNum int64,
	winner []byte,
) error {
	insertStr := "INSERT INTO winners (block_number, builder_address) VALUES ($1, $2)"

	// Convert winner to base64 string for storage
	winnerBase64 := base64.StdEncoding.EncodeToString(winner)

	_, err := s.db.ExecContext(ctx, insertStr, blockNum, winnerBase64)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetWinner(
	ctx context.Context,
	blockNum int64,
) (updater.Winner, error) {
	winner := updater.Winner{}
	var winnerBase64 string
	err := s.db.QueryRowContext(
		ctx,
		"SELECT builder_address FROM winners WHERE block_number = $1",
		blockNum,
	).Scan(&winnerBase64)
	if err != nil {
		return winner, err
	}

	// Convert winner from base64 string to raw bytes
	winner.Winner, err = base64.StdEncoding.DecodeString(winnerBase64)
	if err != nil {
		return winner, err
	}

	return winner, nil
}

func (s *Store) LastWinnerBlock() (int64, error) {
	var lastBlock sql.NullInt64
	err := s.db.QueryRow("SELECT block_number FROM winners ORDER BY block_number DESC LIMIT 1").Scan(&lastBlock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if !lastBlock.Valid {
		return 0, nil
	}
	return lastBlock.Int64, nil
}

func (s *Store) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount *big.Int,
	builder []byte,
	bidID []byte,
	settlementType updater.SettlementType,
	decayPercentage int64,
	postingTxnHash []byte,
	postingTxnNonce uint64,
	options []byte,
) error {
	columns := []string{
		"commitment_index",
		"transaction",
		"block_number",
		"builder_address",
		"type",
		"amount",
		"bid_id",
		"settled",
		"chainhash",
		"nonce",
		"decay_percentage",
		"options",
	}

	// Convert byte slices to base64 strings for storage
	commitmentIdxBase64 := base64.StdEncoding.EncodeToString(commitmentIdx)
	builderBase64 := base64.StdEncoding.EncodeToString(builder)
	bidIDBase64 := base64.StdEncoding.EncodeToString(bidID)
	postingTxnHashBase64 := base64.StdEncoding.EncodeToString(postingTxnHash)
	optionsBase64 := base64.StdEncoding.EncodeToString(options)

	values := []interface{}{
		commitmentIdxBase64,
		txHash,
		blockNum,
		builderBase64,
		settlementType,
		amount.String(),
		bidIDBase64,
		false,
		postingTxnHashBase64,
		postingTxnNonce,
		decayPercentage,
		optionsBase64,
	}
	placeholder := make([]string, len(values))
	for i := range columns {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStr := fmt.Sprintf(
		"INSERT INTO settlements (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholder, ", "),
	)

	_, err := s.db.ExecContext(ctx, insertStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) IsSettled(
	ctx context.Context,
	commitmentIdx []byte,
) (bool, error) {
	var settled bool
	commitmentIdxBase64 := base64.StdEncoding.EncodeToString(commitmentIdx)
	err := s.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM settlements WHERE commitment_index = $1)",
		commitmentIdxBase64,
	).Scan(&settled)
	if err != nil {
		return false, err
	}

	return settled, nil
}

func (s *Store) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	txHashBase64 := base64.StdEncoding.EncodeToString(txHash.Bytes())
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO sent_transactions (hash, nonce, settled) VALUES ($1, $2, false)",
		txHashBase64,
		nonce,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Update(ctx context.Context, txHash common.Hash, status string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txHashBase64 := base64.StdEncoding.EncodeToString(txHash.Bytes())

	_, err = tx.ExecContext(
		ctx,
		"UPDATE sent_transactions SET status = $1, settled = true WHERE hash = $2",
		status,
		txHashBase64,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if status == "success" {
		_, err = tx.ExecContext(
			ctx,
			"UPDATE settlements SET settled = true WHERE chainhash = $1",
			txHashBase64,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) PendingTxnCount() (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM sent_transactions WHERE hash IS NOT NULL AND settled = false",
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	rows, err := s.db.Query("SELECT hash, nonce FROM sent_transactions WHERE settled = false")
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	var txns []*txmonitor.TxnDetails
	for rows.Next() {
		var hashBase64 string
		var nonce uint64
		if err := rows.Scan(&hashBase64, &nonce); err != nil {
			return nil, err
		}

		hash, err := base64.StdEncoding.DecodeString(hashBase64)
		if err != nil {
			return nil, err
		}

		txns = append(txns, &txmonitor.TxnDetails{
			Hash:  common.BytesToHash(hash),
			Nonce: nonce,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return txns, nil
}

func (s *Store) LastBlock() (uint64, error) {
	var lastBlock sql.NullInt64
	err := s.db.QueryRow("SELECT value FROM integers WHERE key = 'last_block'").Scan(&lastBlock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	if !lastBlock.Valid {
		return 0, nil
	}
	return uint64(lastBlock.Int64), nil
}

func (s *Store) SetLastBlock(blockNum uint64) error {
	_, err := s.db.Exec(
		"INSERT INTO integers (key, value) VALUES ('last_block', $1) ON CONFLICT (key) DO UPDATE SET value = $1",
		blockNum,
	)
	if err != nil {
		return err
	}
	return nil
}
