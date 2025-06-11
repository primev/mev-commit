package store

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type rpcstore struct {
	db *pebble.DB
}

func New(path string) (*rpcstore, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &rpcstore{
		db: db,
	}, nil
}

func (s *rpcstore) Close() error {
	return errors.Join(s.db.Flush(), s.db.Close())
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

	commitmentsData, err := json.Marshal(commitments)
	if err != nil {
		return err
	}

	// Create a composite key for the block number and transaction hash
	key := []byte(fmt.Sprintf("%d:%s", blockNumber, txn.Hash().Hex()))
	// Store the transaction and commitments in the database
	if err := s.db.Set(key, append(txnData, commitmentsData...), nil); err != nil {
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
	txnHash string,
) (*types.Transaction, []*bidderapiv1.Commitment, error) {
	if txnHash == "" {
		return nil, nil, errors.New("transaction hash cannot be empty")
	}

	txnKey := []byte(fmt.Sprintf("txn:%s", txnHash))
	blkNumBuf, closer, err := s.db.Get(txnKey)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = closer.Close()
	}()

	blockNumber := binary.BigEndian.Uint64(blkNumBuf)
	if blockNumber == 0 {
		return nil, nil, fmt.Errorf("transaction %s not found", txnHash)
	}

	key := []byte(fmt.Sprintf("%d:%s", blockNumber, txnHash))
	txnData, closer, err := s.db.Get(key)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = closer.Close()
	}()

	var commitments []*bidderapiv1.Commitment
	txnLen := len(txnData)
	if txnLen < 32 {
		return nil, nil, fmt.Errorf("invalid transaction data length: %d", txnLen)
	}

	txn := new(types.Transaction)
	if err := txn.UnmarshalBinary(txnData[:txnLen-32]); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	if err := json.Unmarshal(txnData[txnLen-32:], &commitments); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal commitments: %w", err)
	}

	return txn, commitments, nil
}

func (s *rpcstore) DeductBalance(
	ctx context.Context,
	account string,
	amount *big.Int,
) error {
	if account == "" || amount == nil || amount.Sign() <= 0 {
		return errors.New("invalid account or amount")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return err
	}
	defer func() {
		_ = closer.Close()
	}()

	currentBalanceBig := new(big.Int)
	if err := currentBalanceBig.UnmarshalText(currentBalance); err != nil {
		return fmt.Errorf("failed to unmarshal current balance: %w", err)
	}
	if currentBalanceBig.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance for account %s: %w", account, ErrInsufficientBalance)
	}
	newBalance := new(big.Int).Sub(currentBalanceBig, amount)
	if err := s.db.Set(balanceKey, newBalance.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to update balance for account %s: %w", account, err)
	}

	return nil
}

func (s *rpcstore) RefundBalance(
	ctx context.Context,
	account string,
	amount *big.Int,
) error {
	if account == "" || amount == nil || amount.Sign() <= 0 {
		return errors.New("invalid account or amount")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return err
	}
	defer func() {
		_ = closer.Close()
	}()

	currentBalanceBig := new(big.Int)
	if err := currentBalanceBig.UnmarshalText(currentBalance); err != nil {
		return fmt.Errorf("failed to unmarshal current balance: %w", err)
	}

	newBalance := new(big.Int).Add(currentBalanceBig, amount)
	if err := s.db.Set(balanceKey, newBalance.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to update balance for account %s: %w", account, err)
	}

	return nil
}

func (s *rpcstore) HasBalance(
	ctx context.Context,
	account string,
	amount *big.Int,
) error {
	if account == "" || amount == nil || amount.Sign() <= 0 {
		return errors.New("invalid account or amount")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return err
	}
	defer func() {
		_ = closer.Close()
	}()

	currentBalanceBig := new(big.Int)
	if err := currentBalanceBig.UnmarshalText(currentBalance); err != nil {
		return fmt.Errorf("failed to unmarshal current balance: %w", err)
	}

	if currentBalanceBig.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance for account %s: %w", account, ErrInsufficientBalance)
	}
	return nil
}
