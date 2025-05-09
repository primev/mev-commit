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
	"github.com/primev/mev-commit/x/contracts/txmonitor"
)

var transfers = `
CREATE TABLE IF NOT EXISTS %s_transfers (
	transfer_idx BIGINT PRIMARY KEY,
	amount NUMERIC(24, 0),
	recipient TEXT,
	nonce BIGINT,
	chainhash TEXT,
	settled BOOLEAN
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
	db        *sql.DB
	component string
}

func NewStore(db *sql.DB, component string) (*Store, error) {
	for _, table := range []string{
		fmt.Sprintf(transfers, strings.ToLower(component)),
		transactionsTable,
		integerTable,
	} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		db:        db,
		component: strings.ToLower(component),
	}, nil
}

func (s *Store) StoreTransfer(
	ctx context.Context,
	transferIdx *big.Int,
	amount *big.Int,
	recipient common.Address,
	nonce uint64,
	chainHash common.Hash,
) error {
	recipientBase64 := base64.StdEncoding.EncodeToString(recipient.Bytes())
	chainHashBase64 := base64.StdEncoding.EncodeToString(chainHash.Bytes())
	insertQuery := fmt.Sprintf(
		"INSERT INTO %s_transfers (transfer_idx, amount, recipient, nonce, chainhash, settled) VALUES ($1, $2, $3, $4, $5, false)",
		s.component,
	)
	_, err := s.db.ExecContext(
		ctx,
		insertQuery,
		transferIdx.Uint64(),
		amount.String(),
		recipientBase64,
		nonce,
		chainHashBase64,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) MarkTransferSettled(
	ctx context.Context,
	transferIdx *big.Int,
) error {
	updateQuery := fmt.Sprintf(
		"UPDATE %s_transfers SET settled = true WHERE transfer_idx = $1",
		s.component,
	)
	_, err := s.db.ExecContext(ctx, updateQuery, transferIdx.Uint64())
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) IsSettled(
	ctx context.Context,
	transferIdx *big.Int,
) (bool, error) {
	var settled bool
	query := fmt.Sprintf(
		"SELECT settled FROM %s_transfers WHERE transfer_idx = $1",
		s.component,
	)
	err := s.db.QueryRowContext(ctx, query, transferIdx.Uint64()).Scan(&settled)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
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
	txHashBase64 := base64.StdEncoding.EncodeToString(txHash.Bytes())

	_, err := s.db.ExecContext(
		ctx,
		"UPDATE sent_transactions SET status = $1, settled = true WHERE hash = $2",
		status,
		txHashBase64,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	var query = `
    SELECT hash, nonce
    FROM sent_transactions
    WHERE settled = false AND EXISTS (
        SELECT 1 FROM %s_transfers WHERE %s_transfers.chainhash = sent_transactions.hash
    )
	`
	rows, err := s.db.Query(fmt.Sprintf(query, s.component, s.component))
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
	query := fmt.Sprintf("SELECT value FROM integers WHERE key = 'last_block_%s'", s.component)
	err := s.db.QueryRow(query).Scan(&lastBlock)
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
	query := fmt.Sprintf(
		"INSERT INTO integers (key, value) VALUES ('last_block_%s', $1) ON CONFLICT (key) DO UPDATE SET value = $1",
		s.component,
	)
	_, err := s.db.Exec(query, blockNum)
	if err != nil {
		return err
	}
	return nil
}
