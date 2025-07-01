package store

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"google.golang.org/protobuf/proto"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrNotFound            = errors.New("not found")
)

var transactionsTable = `
CREATE TABLE IF NOT EXISTS mcTransactions (
	hash TEXT PRIMARY KEY,
	nonce BIGINT,
	raw_transaction TEXT,
	block_number BIGINT,
	sender TEXT,
	tx_type INTEGER,
	status TEXT,
	details TEXT
);`

var commitmentsTable = `
CREATE TABLE IF NOT EXISTS commitments (
	commitment_digest TEXT PRIMARY KEY,
	transaction_hash TEXT,
	provider_address TEXT,
	commitment_data BYTEA,
	FOREIGN KEY (transaction_hash) REFERENCES mcTransactions (hash) ON DELETE CASCADE
);`

var balancesTable = `
CREATE TABLE IF NOT EXISTS balances (
	account TEXT PRIMARY KEY,
	balance NUMERIC(24, 0)
);`

type rpcstore struct {
	db *sql.DB
}

func New(db *sql.DB) (*rpcstore, error) {
	for _, table := range []string{
		transactionsTable,
		commitmentsTable,
		balancesTable,
	} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &rpcstore{
		db: db,
	}, nil
}

func (s *rpcstore) Close() error {
	return s.db.Close()
}

func (s *rpcstore) AddQueuedTransaction(ctx context.Context, tx *sender.Transaction) error {
	insertQuery := `
	INSERT INTO mcTransactions (hash, nonce, raw_transaction, sender, tx_type, status)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := s.db.ExecContext(
		ctx,
		insertQuery,
		tx.Hash().Hex(),
		tx.Nonce(),
		tx.Raw,
		tx.Sender.Hex(),
		int(tx.Type),
		string(sender.TxStatusPending),
	)
	if err != nil {
		return fmt.Errorf("failed to add queued transaction: %w", err)
	}

	return nil
}

func parseTransactionsFromRows(rows *sql.Rows) ([]*sender.Transaction, error) {
	var transactions []*sender.Transaction
	for rows.Next() {
		var (
			rawTransaction string
			senderAddress  string
			txType         int
			blockNum       sql.NullInt64
			status         string
			details        sql.NullString
		)
		err := rows.Scan(&rawTransaction, &blockNum, &senderAddress, &txType, &status, &details)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		txStr, err := hex.DecodeString(strings.TrimPrefix(rawTransaction, "0x"))
		if err != nil {
			return nil, fmt.Errorf("failed to decode raw transaction: %w", err)
		}
		parsedTxn := new(types.Transaction)
		if err := parsedTxn.UnmarshalBinary(txStr); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
		}
		txn := &sender.Transaction{
			Transaction: parsedTxn,
			Raw:         rawTransaction,
			BlockNumber: blockNum.Int64,
			Sender:      common.HexToAddress(senderAddress),
			Type:        sender.TxType(txType),
			Status:      sender.TxStatus(status),
			Details:     details.String,
		}
		transactions = append(transactions, txn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

// GetQueuedTransactions retrieves the next pending transaction for each sender.
func (s *rpcstore) GetQueuedTransactions(ctx context.Context) ([]*sender.Transaction, error) {
	query := `
	SELECT t1.raw_transaction, t1.block_number, t1.sender, t1.tx_type, t1.status, t1.details
	FROM mcTransactions t1
	INNER JOIN (
		SELECT sender, MIN(nonce) AS min_nonce
		FROM mcTransactions
		WHERE status = 'pending'
		GROUP BY sender
	) t2
	ON t1.sender = t2.sender AND t1.nonce = t2.min_nonce
	WHERE t1.status = 'pending';
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*sender.Transaction{}, nil // No pending transactions found
		}
		return nil, fmt.Errorf("failed to get queued transactions: %w", err)
	}

	transactions, err := parseTransactionsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transactions from rows: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

func (s *rpcstore) GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*sender.Transaction, error) {
	query := `
	SELECT raw_transaction, block_number, sender, tx_type, status, details
	FROM mcTransactions
	WHERE hash = $1;
	`
	row := s.db.QueryRowContext(ctx, query, txnHash.Hex())
	var (
		rawTransaction string
		senderAddress  string
		txType         int
		status         string
		blockNum       sql.NullInt64
		details        sql.NullString
	)
	err := row.Scan(&rawTransaction, &blockNum, &senderAddress, &txType, &status, &details)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction %s not found: %w", txnHash.Hex(), ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}
	txStr, err := hex.DecodeString(strings.TrimPrefix(rawTransaction, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw transaction: %w", err)
	}
	parsedTxn := new(types.Transaction)
	if err := parsedTxn.UnmarshalBinary(txStr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	txn := &sender.Transaction{
		Transaction: parsedTxn,
		Raw:         rawTransaction,
		BlockNumber: blockNum.Int64,
		Sender:      common.HexToAddress(senderAddress),
		Type:        sender.TxType(txType),
		Status:      sender.TxStatus(status),
		Details:     details.String,
	}

	return txn, nil
}

func (s *rpcstore) GetTransactionsForBlock(ctx context.Context, blockNumber int64) ([]*sender.Transaction, error) {
	query := `
	SELECT raw_transaction, block_number, sender, tx_type, status, details
	FROM mcTransactions
	WHERE block_number = $1 AND status = 'pre-confirmed';
	`
	rows, err := s.db.QueryContext(ctx, query, blockNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*sender.Transaction{}, nil // No transactions found for this block
		}
		return nil, fmt.Errorf("failed to get transactions for block %d: %w", blockNumber, err)
	}
	transactions, err := parseTransactionsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transactions from rows: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows for block %d: %w", blockNumber, err)
	}

	// If no transactions found, return an empty slice
	if len(transactions) == 0 {
		return []*sender.Transaction{}, nil
	}

	return transactions, nil
}

func (s *rpcstore) StoreTransaction(
	ctx context.Context,
	txn *sender.Transaction,
	commitments []*bidderapiv1.Commitment,
) error {
	if txn.Status == sender.TxStatusPending {
		return fmt.Errorf("transaction must not be in pending status, got %s", txn.Status)
	}

	if txn.BlockNumber == 0 && txn.Status != sender.TxStatusFailed {
		return fmt.Errorf("block number must be set for successful transactions, got %d", txn.BlockNumber)
	}

	dbTxn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	updateTxns := `
	UPDATE mcTransactions
	SET block_number = $1, status = $2, details = $3
	WHERE hash = $4;
	`

	_, err = dbTxn.ExecContext(ctx, updateTxns, txn.BlockNumber, string(txn.Status), txn.Details, txn.Hash().Hex())
	if err != nil {
		_ = dbTxn.Rollback()
		return fmt.Errorf("failed to update transaction %s: %w", txn.Hash().Hex(), err)
	}

	if txn.Status != sender.TxStatusFailed {
		for _, commitment := range commitments {
			insertCommitment := `
			INSERT INTO commitments (commitment_digest, transaction_hash, provider_address, commitment_data)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (commitment_digest) DO NOTHING;
			`
			commitmentData, err := proto.Marshal(commitment)
			if err != nil {
				_ = dbTxn.Rollback()
				return fmt.Errorf("failed to marshal commitment: %w", err)
			}

			_, err = dbTxn.ExecContext(
				ctx,
				insertCommitment,
				commitment.CommitmentDigest,
				txn.Hash().Hex(),
				commitment.ProviderAddress,
				commitmentData,
			)
			if err != nil {
				_ = dbTxn.Rollback()
				return fmt.Errorf("failed to insert commitment for transaction %s: %w", txn.Hash().Hex(), err)
			}
		}
	}

	if err := dbTxn.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *rpcstore) GetTransactionCommitments(ctx context.Context, txnHash common.Hash) ([]*bidderapiv1.Commitment, error) {
	query := `
	SELECT commitment_data
	FROM commitments
	WHERE transaction_hash = $1;
	`
	rows, err := s.db.QueryContext(ctx, query, txnHash.Hex())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no commitments found for transaction %s: %w", txnHash.Hex(), ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get commitments for transaction %s: %w", txnHash.Hex(), err)
	}

	var commitments []*bidderapiv1.Commitment
	for rows.Next() {
		var commitmentData []byte
		err := rows.Scan(&commitmentData)
		if err != nil {
			return nil, fmt.Errorf("failed to scan commitment data for transaction %s: %w", txnHash.Hex(), err)
		}
		commitment := &bidderapiv1.Commitment{}
		if err := proto.Unmarshal(commitmentData, commitment); err != nil {
			return nil, fmt.Errorf("failed to unmarshal commitment data for transaction %s: %w", txnHash.Hex(), err)
		}
		commitments = append(commitments, commitment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows for transaction %s: %w", txnHash.Hex(), err)
	}
	if len(commitments) == 0 {
		return nil, fmt.Errorf("no commitments found for transaction %s: %w", txnHash.Hex(), ErrNotFound)
	}
	return commitments, nil
}

// GetCurrentNonce retrieves the next nonce for a given sender address by looking at the
// pending transactions in the database. If there are no pending transactions, it returns 0.
// The RPC would proxy this call to the underlying Ethereum node to get the current nonce in
// case if 0 is returned.
func (s *rpcstore) GetCurrentNonce(ctx context.Context, sender common.Address) uint64 {
	query := `
	SELECT COALESCE(MAX(nonce), 0)
	FROM mcTransactions
	WHERE sender = $1 AND status = 'pending';
	`
	row := s.db.QueryRowContext(ctx, query, sender.Hex())
	var nextNonce uint64
	err := row.Scan(&nextNonce)
	if err != nil {
		return 0 // If no pending transactions found, return 0 as the next nonce
	}
	return nextNonce
}

func (s *rpcstore) DeductBalance(
	ctx context.Context,
	account common.Address,
	amount *big.Int,
) error {
	query := `
	UPDATE balances
	SET balance = balance - $1
	WHERE account = $2 AND balance >= $1;
	`
	_, err := s.db.ExecContext(ctx, query, amount.String(), account.Hex())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("account %s not found or insufficient balance: %w", account.Hex(), ErrInsufficientBalance)
		}
		return fmt.Errorf("failed to deduct balance for account %s: %w", account.Hex(), err)
	}

	return nil
}

func (s *rpcstore) AddBalance(
	ctx context.Context,
	account common.Address,
	amount *big.Int,
) error {
	if account == (common.Address{}) || amount == nil || amount.Sign() <= 0 {
		return fmt.Errorf("invalid account or amount: account=%s, amount=%s", account.Hex(), amount.String())
	}

	query := `
	INSERT INTO balances (account, balance)
	VALUES ($1, $2)
	ON CONFLICT (account) DO UPDATE SET balance = balances.balance + $2
	WHERE balances.balance + $2 >= 0;
	`

	_, err := s.db.ExecContext(ctx, query, account.Hex(), amount.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("account %s not found or insufficient balance: %w", account.Hex(), ErrInsufficientBalance)
		}
		return fmt.Errorf("failed to add balance for account %s: %w", account.Hex(), err)
	}

	return nil
}

func (s *rpcstore) HasBalance(
	ctx context.Context,
	account common.Address,
	amount *big.Int,
) bool {
	if account == (common.Address{}) || amount == nil || amount.Sign() <= 0 {
		return false
	}

	query := `
	SELECT balance
	FROM balances
	WHERE account = $1;
	`

	row := s.db.QueryRowContext(ctx, query, account.Hex())
	var currentBalanceString string
	err := row.Scan(&currentBalanceString)
	if err != nil {
		return false
	}
	currentBalance, ok := new(big.Int).SetString(currentBalanceString, 10)
	if !ok {
		return false
	}

	return currentBalance.Cmp(amount) >= 0
}

func (s *rpcstore) GetBalance(
	ctx context.Context,
	account common.Address,
) (*big.Int, error) {
	if account == (common.Address{}) {
		return nil, errors.New("account cannot be empty")
	}

	query := `
	SELECT balance
	FROM balances
	WHERE account = $1;
	`

	row := s.db.QueryRowContext(ctx, query, account.Hex())
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, fmt.Errorf("account %s not found: %w", account.Hex(), ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get balance for account %s: %w", account.Hex(), row.Err())
	}

	var balance string
	err := row.Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account %s not found: %w", account.Hex(), ErrNotFound)
		}
		return nil, fmt.Errorf("failed to scan balance for account %s: %w", account.Hex(), err)
	}

	// Convert the balance string to a big.Int
	balanceInt, ok := big.NewInt(0).SetString(balance, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert balance string to big.Int for account %s", account.Hex())
	}

	return balanceInt, nil
}
