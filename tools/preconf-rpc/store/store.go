package store

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrNotFound            = errors.New("not found")
)

var transactionsTable = `
CREATE TABLE IF NOT EXISTS transactions (
	hash TEXT PRIMARY KEY,
	nonce BIGINT,
	raw_transaction TEXT,
	block_number BIGINT,
	sender TEXT,
	tx_type INTEGER,
	status TEXT,
	details TEXT,
);`

var balancesTable = `
CREATE TABLE IF NOT EXISTS balances (
	account TEXT PRIMARY KEY,
	balance NUMERIC(24, 0)
);`

type rpcstore struct {
	db *sql.DB
}

func New(path string) (*rpcstore, error) {
	return nil, nil
}

func (s *rpcstore) Close() error {
	return s.db.Close()
}

func (s *rpcstore) AddQueuedTransaction(tx *sender.Transaction) error {
	if tx.Status != sender.TxStatusPending {
		return fmt.Errorf("transaction must be in pending status, got %s", tx.Status)
	}
	insertQuery := `
	INSERT INTO transactions (hash, nonce, raw_transaction, sender, tx_type, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(insertQuery, tx.Hash().Hex(), tx.Nonce(), tx.Raw, tx.Sender.Hex(), int(tx.Type), string(tx.Status))
	if err != nil {
		return fmt.Errorf("failed to add queued transaction: %w", err)
	}

	return nil
}

func (s *rpcstore) GetQueuedTransactions() ([]*sender.Transaction, error) {
	query := `
	SELECT t1.raw_transaction, t1.sender, t1.tx_type
	FROM transactions t1
	INNER JOIN (
		SELECT sender, MIN(nonce) AS min_nonce
		FROM transactions
		WHERE status = 'pending'
		GROUP BY sender
	) t2
	ON t1.sender = t2.sender AND t1.nonce = t2.min_nonce
	WHERE t1.status = 'pending';
	`

	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*sender.Transaction{}, nil // No pending transactions found
		}
		return nil, fmt.Errorf("failed to get queued transactions: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var transactions []*sender.Transaction
	for rows.Next() {
		var (
			rawTransaction string
			senderAddress  string
			txType         int
		)
		err := rows.Scan(&rawTransaction, &senderAddress, &txType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		txStr, err := hex.DecodeString(rawTransaction)
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
			Sender:      common.HexToAddress(senderAddress),
			Type:        sender.TxType(txType),
			Status:      sender.TxStatusPending,
		}
		transactions = append(transactions, txn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

func (s *rpcstore) GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*sender.Transaction, error) {
	query := `
	SELECT raw_transaction, sender, tx_type, status, details
	FROM transactions
	WHERE hash = $1;
	`
	row := s.db.QueryRowContext(ctx, query, txnHash.Hex())
	var (
		rawTransaction string
		senderAddress  string
		txType         int
		status         string
		details        sql.NullString
	)
	err := row.Scan(&rawTransaction, &senderAddress, &txType, &status, &details)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction %s not found: %w", txnHash.Hex(), ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}
	txStr, err := hex.DecodeString(rawTransaction)
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
		Sender:      common.HexToAddress(senderAddress),
		Type:        sender.TxType(txType),
		Status:      sender.TxStatus(status),
		Details:     details.String,
	}

	return txn, nil
}

func (s *rpcstore) GetTransactionsForBlock(ctx context.Context, blockNumber int64) ([]*sender.Transaction, error) {
	query := `
	SELECT raw_transaction, sender, tx_type, status, details
	FROM transactions
	WHERE block_number = $1 AND status = 'pre-confirmed';
	`
	rows, err := s.db.QueryContext(ctx, query, blockNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*sender.Transaction{}, nil // No transactions found for this block
		}
		return nil, fmt.Errorf("failed to get transactions for block %d: %w", blockNumber, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()
	var transactions []*sender.Transaction
	for rows.Next() {
		var (
			rawTransaction string
			senderAddress  string
			txType         int
			status         string
			details        sql.NullString
		)
		err := rows.Scan(&rawTransaction, &senderAddress, &txType, &status, &details)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		txStr, err := hex.DecodeString(rawTransaction)
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

func (s *rpcstore) GetTransactionCommitments(ctx context.Context, txnHash common.Hash) ([]*bidderapiv1.Commitment, error) {
	return nil, nil
}

func (s *rpcstore) GetCurrentNonce(sender common.Address) uint64 {
	query := `
	SELECT COALESCE(MAX(nonce), 0) + 1
	FROM transactions
	WHERE sender = $1 AND status = 'pending';
	`
	row := s.db.QueryRow(query, sender.Hex())
	var nextNonce uint64
	err := row.Scan(&nextNonce)
	if err != nil {
		return 0 // If no pending transactions found, return 0 as the next nonce
	}
	return nextNonce
}

func (s *rpcstore) StorePreconfirmedTransaction(
	ctx context.Context,
	blockNumber int64,
	txn *types.Transaction,
	commitments []*bidderapiv1.Commitment,
) error {
	if blockNumber <= 0 || txn == nil || commitments == nil {
		return errors.New("invalid input parameters")
	}

	// Serialize the transaction and commitments
	txnData, err := txn.MarshalBinary()
	if err != nil {
		return err
	}

	txnDataLenBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(txnDataLenBuf, uint64(len(txnData)))
	txnDataWithLen := append(txnDataLenBuf, txnData...)

	commitmentsData, err := json.Marshal(commitments)
	if err != nil {
		return err
	}

	// Create a composite key for the block number and transaction hash
	key := []byte(fmt.Sprintf("%d:%s", blockNumber, txn.Hash().Hex()))
	// Store the transaction and commitments in the database
	if err := s.db.Set(key, append(txnDataWithLen, commitmentsData...), nil); err != nil {
		return err
	}

	blockNumBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(blockNumBuf, uint64(blockNumber))

	txnKey := []byte(fmt.Sprintf("txn:%s", txn.Hash().Hex()))
	if err := s.db.Set(txnKey, blockNumBuf, nil); err != nil {
		return err
	}

	return nil
}

func (s *rpcstore) GetPreconfirmedTransaction(
	ctx context.Context,
	txnHash common.Hash,
) (*types.Transaction, []*bidderapiv1.Commitment, error) {
	if txnHash == (common.Hash{}) {
		return nil, nil, errors.New("transaction hash cannot be empty")
	}

	txnKey := []byte(fmt.Sprintf("txn:%s", txnHash.Hex()))
	blkNumBuf, closer, err := s.db.Get(txnKey)
	if err != nil {
		return nil, nil, err
	}

	blockNumber := binary.BigEndian.Uint64(blkNumBuf)
	if blockNumber == 0 {
		return nil, nil, fmt.Errorf("transaction %s not found", txnHash)
	}

	_ = closer.Close() // Close the closer from Get

	key := []byte(fmt.Sprintf("%d:%s", blockNumber, txnHash.Hex()))
	txnData, closer, err := s.db.Get(key)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = closer.Close()
	}()

	// The first 8 bytes are the length of the transaction data
	txnDataLen := binary.BigEndian.Uint64(txnData[:8])

	txn := new(types.Transaction)
	if err := txn.UnmarshalBinary(txnData[8 : 8+txnDataLen]); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	var commitments []*bidderapiv1.Commitment
	if err := json.Unmarshal(txnData[8+txnDataLen:], &commitments); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal commitments: %w", err)
	}

	return txn, commitments, nil
}

func (s *rpcstore) GetPreconfirmedTransactionsForBlock(
	ctx context.Context,
	blockNumber int64,
) ([]*types.Transaction, error) {
	if blockNumber <= 0 {
		return nil, errors.New("invalid block number")
	}

	keyPrefix := []byte(fmt.Sprintf("%d:", blockNumber))
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: keyPrefix,
		UpperBound: append(keyPrefix, 0xFF),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create iterator for block %d: %w", blockNumber, err)
	}
	defer func() {
		_ = iter.Close()
	}()

	var transactions []*types.Transaction
	for iter.First(); iter.Valid(); iter.Next() {
		if !bytes.Equal(iter.Key()[:len(keyPrefix)], keyPrefix) {
			continue
		}
		txnData := iter.Value()
		if len(txnData) < 8 {
			return nil, fmt.Errorf("invalid transaction data length for block %d", blockNumber)
		}
		txnDataLen := binary.BigEndian.Uint64(txnData[:8])
		if len(txnData) < int(8+txnDataLen) {
			return nil, fmt.Errorf("invalid transaction data length for block %d", blockNumber)
		}

		txn := new(types.Transaction)
		if err := txn.UnmarshalBinary(txnData[8 : 8+txnDataLen]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
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
	err := s.db.ExecContext(ctx, query, amount.Bytes(), account.Hex())
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
		return errors.New("invalid account or amount")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account.Hex()))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			// If the account does not exist, we create a new one with the initial balance
			bal := new(big.Int)
			currentBalance = bal.Bytes() // Default balance for a new account
		} else {
			return fmt.Errorf("failed to get balance for account %s: %w", account, err)
		}
	}
	defer func() {
		if closer != nil {
			_ = closer.Close()
		}
	}()

	currentBalanceBig := new(big.Int).SetBytes(currentBalance)

	newBalance := new(big.Int).Add(currentBalanceBig, amount)
	if err := s.db.Set(balanceKey, newBalance.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to update balance for account %s: %w", account, err)
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

	balanceKey := []byte(fmt.Sprintf("balance:%s", account.Hex()))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return false
	}
	defer func() {
		_ = closer.Close()
	}()

	currentBalanceBig := new(big.Int).SetBytes(currentBalance)

	return currentBalanceBig.Cmp(amount) >= 0
}

func (s *rpcstore) GetBalance(
	ctx context.Context,
	account common.Address,
) (*big.Int, error) {
	if account == (common.Address{}) {
		return nil, errors.New("account cannot be empty")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account.Hex()))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = closer.Close()
	}()

	return new(big.Int).SetBytes(currentBalance), nil
}
