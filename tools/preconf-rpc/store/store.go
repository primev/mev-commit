package store

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/common"
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

func (s *rpcstore) DeductBalance(
	ctx context.Context,
	account common.Address,
	amount *big.Int,
) error {
	if account == (common.Address{}) || amount == nil || amount.Sign() <= 0 {
		fmt.Println("invalid account or amount: %s, %s", account.Hex(), amount.String())
		return errors.New("invalid account or amount")
	}

	balanceKey := []byte(fmt.Sprintf("balance:%s", account.Hex()))
	currentBalance, closer, err := s.db.Get(balanceKey)
	if err != nil {
		return err
	}
	defer func() {
		_ = closer.Close()
	}()

	currentBalanceBig := new(big.Int).SetBytes(currentBalance)
	if currentBalanceBig.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance for account %s: %w", account, ErrInsufficientBalance)
	}
	newBalance := new(big.Int).Sub(currentBalanceBig, amount)
	if err := s.db.Set(balanceKey, newBalance.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to update balance for account %s: %w", account, err)
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
			currentBalance = []byte("0") // Default balance for a new account
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
