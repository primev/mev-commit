package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
	"github.com/primev/mev-commit/oracle/pkg/updater"
)

var settlementType = `
DO $$ BEGIN
    CREATE TYPE settlement_type AS ENUM ('reward', 'slash');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;`

var settlementsTable = `
CREATE TABLE IF NOT EXISTS settlements (
	commitment_index BYTEA PRIMARY KEY,
	transaction TEXT,
	block_number BIGINT,
	builder_address BYTEA,
	type settlement_type,
	amount NUMERIC(24, 0),
	bid_id BYTEA,
	chainhash BYTEA,
	nonce BIGINT,
	settled BOOLEAN,
	decay_percentage BIGINT,
	settlement_window BIGINT
);`

var encryptedCommitmentsTable = `
CREATE TABLE IF NOT EXISTS encrypted_commitments (
	commitment_index BYTEA PRIMARY KEY,
	committer BYTEA,
	commitment_hash BYTEA,
	commitment_signature BYTEA,
	dispatch_timestamp BIGINT
);`

var winnersTable = `
CREATE TABLE IF NOT EXISTS winners (
	block_number BIGINT PRIMARY KEY,
	builder_address BYTEA,
	settlement_window BIGINT
);`

var transactionsTable = `
CREATE TABLE IF NOT EXISTS sent_transactions (
	hash BYTEA PRIMARY KEY,
	nonce BIGINT,
	settled BOOLEAN,
	status TEXT
);`

var integerTable = `
CREATE TABLE IF NOT EXISTS integers (
	key TEXT PRIMARY KEY,
	value BIGINT
);`

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	for _, table := range []string{
		settlementType,
		settlementsTable,
		encryptedCommitmentsTable,
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
	window int64,
) error {
	insertStr := "INSERT INTO winners (block_number, builder_address, settlement_window) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, insertStr, blockNum, winner, window)
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
	err := s.db.QueryRowContext(
		ctx,
		"SELECT builder_address, settlement_window FROM winners WHERE block_number = $1",
		blockNum,
	).Scan(&winner.Winner, &winner.Window)
	if err != nil {
		return winner, err
	}
	return winner, nil
}

func (s *Store) AddEncryptedCommitment(
	ctx context.Context,
	commitmentIdx []byte,
	committer []byte,
	commitmentHash []byte,
	commitmentSignature []byte,
	dispatchTimestamp uint64,
) error {
	columns := []string{
		"commitment_index",
		"committer",
		"commitment_hash",
		"commitment_signature",
		"dispatch_timestamp",
	}
	values := []interface{}{
		commitmentIdx,
		committer,
		commitmentHash,
		commitmentSignature,
		dispatchTimestamp,
	}
	placeholder := make([]string, len(values))
	for i := range columns {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStr := fmt.Sprintf(
		"INSERT INTO encrypted_commitments (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholder, ", "),
	)

	_, err := s.db.ExecContext(ctx, insertStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount uint64,
	builder []byte,
	bidID []byte,
	settlementType updater.SettlementType,
	decayPercentage int64,
	window int64,
	postingTxnHash []byte,
	postingTxnNonce uint64,
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
		"settlement_window",
	}
	values := []interface{}{
		commitmentIdx,
		txHash,
		blockNum,
		builder,
		settlementType,
		amount,
		bidID,
		false,
		postingTxnHash,
		postingTxnNonce,
		decayPercentage,
		window,
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
	err := s.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM settlements WHERE commitment_index = $1)",
		commitmentIdx,
	).Scan(&settled)
	if err != nil {
		return false, err
	}

	return settled, nil
}

func (s *Store) Settlement(
	ctx context.Context,
	commitmentIdx []byte,
) (updater.Settlement, error) {
	var st updater.Settlement
	err := s.db.QueryRowContext(
		ctx,
		`
		SELECT
			transaction, block_number, builder_address, amount, bid_id, type,
			decay_percentage
		FROM settlements
		WHERE commitment_index = $1`,
		commitmentIdx,
	).Scan(
		&st.TxHash,
		&st.BlockNum,
		&st.Builder,
		&st.Amount,
		&st.BidID,
		&st.Type,
		&st.DecayPercentage,
	)
	if err != nil {
		return st, err
	}
	return st, nil
}

func (s *Store) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO sent_transactions (hash, nonce, settled) VALUES ($1, $2, false)",
		txHash.Bytes(),
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

	_, err = tx.ExecContext(
		ctx,
		"UPDATE sent_transactions SET status = $1, settled = true WHERE hash = $2",
		status,
		txHash.Bytes(),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE settlements SET settled = true WHERE chainhash = $1",
		txHash.Bytes(),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
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
