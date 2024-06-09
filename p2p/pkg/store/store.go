package store

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-radix"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
)

const (
	commitmentNS = "cm/"
	balanceNS    = "bbs/"
	aesKeysNS    = "aes/"

	// provider related keys
	eciesPrivateKeyNS = "ecies/"
	nikePrivateKeyNS  = "nike/"

	// txns related keys
	txNS = "tx/"
)

var (
	commitmentKey = func(blockNum int64, index []byte) string {
		return fmt.Sprintf("%s%d/%s", commitmentNS, blockNum, string(index))
	}
	blockCommitmentPrefix = func(blockNum int64) string {
		return fmt.Sprintf("%s%d", commitmentNS, blockNum)
	}

	balanceKey = func(window *big.Int, bidder common.Address) string {
		return fmt.Sprintf("%s%s/%s", balanceNS, window, bidder)
	}
	blockBalanceKey = func(window *big.Int, bidder common.Address, blockNumber int64) string {
		return fmt.Sprintf("%s%s/%s/%d", balanceNS, window, bidder, blockNumber)
	}
	balancePrefix = func(window *big.Int) string {
		return fmt.Sprintf("%s%s", balanceNS, window)
	}

	bidderAesKey = func(bidder common.Address) string {
		return fmt.Sprintf("%s%s", aesKeysNS, bidder)
	}

	txKey = func(txHash common.Hash) string {
		return fmt.Sprintf("%s%s", txNS, txHash.Hex())
	}
)

type Store struct {
	mu sync.RWMutex
	*radix.Tree
}

type EncryptedPreConfirmationWithDecrypted struct {
	*preconfpb.EncryptedPreConfirmation
	*preconfpb.PreConfirmation
	TxnHash common.Hash
}

func NewStore() *Store {
	return &Store{
		Tree: radix.New(),
	}
}

func (s *Store) LastBlock() (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.Tree.Get("last_block")
	if !ok {
		return 0, nil
	}
	return val.(uint64), nil
}

func (s *Store) SetLastBlock(blockNum uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.Tree.Insert("last_block", blockNum)
	return nil
}

func (s *Store) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := commitmentKey(commitment.Bid.BlockNumber, commitment.EncryptedPreConfirmation.Commitment)
	_, _ = s.Tree.Insert(key, commitment)
}

func (s *Store) GetCommitmentsByBlockNumber(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockCommitmentsKey := blockCommitmentPrefix(blockNum)
	commitments := make([]*EncryptedPreConfirmationWithDecrypted, 0)

	s.Tree.WalkPrefix(blockCommitmentsKey, func(key string, value interface{}) bool {
		commitments = append(commitments, value.(*EncryptedPreConfirmationWithDecrypted))
		return false
	})
	return commitments, nil
}

func (s *Store) DeleteCommitmentByBlockNumber(blockNum int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	blockCommitmentsKey := blockCommitmentPrefix(blockNum)
	_ = s.Tree.DeletePrefix(blockCommitmentsKey)
	return nil
}

func (s *Store) DeleteCommitmentByDigest(blockNum int64, digest [32]byte) (*EncryptedPreConfirmationWithDecrypted, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := commitmentKey(blockNum, digest[:])
	val, deleted := s.Tree.Delete(key)
	return val.(*EncryptedPreConfirmationWithDecrypted), deleted
}

func (s *Store) SetCommitmentIndexByCommitmentDigest(cDigest, cIndex [32]byte) (*EncryptedPreConfirmationWithDecrypted, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var commitmentToReturn *EncryptedPreConfirmationWithDecrypted

	s.Tree.WalkPrefix(commitmentNS, func(key string, value interface{}) bool {
		commitment := value.(*EncryptedPreConfirmationWithDecrypted)
		if bytes.Equal(commitment.EncryptedPreConfirmation.Commitment, cDigest[:]) {
			commitment.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
			commitmentToReturn = commitment
			return true
		}
		return false
	})

	return commitmentToReturn, commitmentToReturn != nil
}

func (s *Store) SetAESKey(bidder common.Address, key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.Tree.Insert(bidderAesKey(bidder), key)
	return nil
}

func (s *Store) GetAESKey(bidder common.Address) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := bidderAesKey(bidder)
	val, ok := s.Tree.Get(key)
	if !ok {
		return nil, nil
	}
	return val.([]byte), nil
}

func (s *Store) SetECIESPrivateKey(key *ecies.PrivateKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.Tree.Insert(eciesPrivateKeyNS, key)
	return nil
}

func (s *Store) GetECIESPrivateKey() (*ecies.PrivateKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.Tree.Get(eciesPrivateKeyNS)
	if !ok {
		return nil, nil
	}
	return val.(*ecies.PrivateKey), nil
}

func (s *Store) SetNikePrivateKey(key *ecdh.PrivateKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.Tree.Insert(nikePrivateKeyNS, key)
	return nil
}

func (s *Store) GetNikePrivateKey() (*ecdh.PrivateKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.Tree.Get(nikePrivateKeyNS)
	if !ok {
		return nil, nil
	}
	return val.(*ecdh.PrivateKey), nil
}

func (s *Store) SetBalance(bidder common.Address, windowNumber, depositedAmount *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := balanceKey(windowNumber, bidder)
	_, _ = s.Tree.Insert(key, depositedAmount)
	return nil
}

func (s *Store) GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := balanceKey(windowNumber, bidder)
	val, ok := s.Tree.Get(key)
	if !ok {
		return nil, nil
	}
	return val.(*big.Int), nil
}

func (s *Store) ClearBalances(windowNumber *big.Int) ([]*big.Int, error) {
	if windowNumber == nil || windowNumber.Cmp(big.NewInt(0)) == -1 {
		return nil, nil
	}

	s.mu.RLock()
	windows := make([]*big.Int, 0)
	s.Tree.WalkPrefix(balanceNS, func(key string, value interface{}) bool {
		parts := strings.Split(key, "/")
		if len(parts) != 3 {
			return false
		}
		w, ok := new(big.Int).SetString(parts[1], 10)
		if !ok {
			return false
		}
		switch w.Cmp(windowNumber) {
		case -1:
			windows = append(windows, w)
		case 0:
			windows = append(windows, w)
			return true
		}
		return false
	})
	s.mu.RUnlock()

	s.mu.Lock()
	for _, w := range windows {
		key := balancePrefix(w)
		_ = s.Tree.DeletePrefix(key)
	}
	s.mu.Unlock()

	return windows, nil
}

func (s *Store) GetBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	blockNumber int64,
) (*big.Int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := blockBalanceKey(window, bidder, blockNumber)
	val, ok := s.Tree.Get(key)
	if !ok {
		return nil, nil
	}
	return val.(*big.Int), nil
}

func (s *Store) SetBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := blockBalanceKey(window, bidder, blockNumber)
	_, _ = s.Tree.Insert(key, amount)
	return nil
}

func (s *Store) RefundBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := blockBalanceKey(window, bidder, blockNumber)
	val, ok := s.Tree.Get(key)
	if !ok {
		_, _ = s.Tree.Insert(key, amount)
		return nil
	}
	amount.Add(amount, val.(*big.Int))
	_, _ = s.Tree.Insert(key, amount)
	return nil
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Tree.Len()
}

// Following are the methods to save and update the transaction details in the store.
// The store is used to keep track of the transactions that are sent to the blockchain and
// have not yet received the transaction receipt. These are used by the debug service
// to show the pending transactions and cancel them if needed. The store hooks up
// to the txmonitor package which allows a component to get notified when the transaction
// is sent to the blockchain and when the transaction receipt is received. As of now,
// the store is in-memory and doesn't persist the transaction details, so the update
// method is used to remove the transaction from the store. This will no longer be seen
// in the pending transactions list.
type TxnDetails struct {
	Hash    common.Hash
	Nonce   uint64
	Created int64
}

// Save implements the txmonitor.Saver interface. It saves the transaction hash and nonce in the store once
// the transaction is sent to the blockchain.
func (s *Store) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.Tree.Insert(txKey(txHash), &TxnDetails{Hash: txHash, Nonce: nonce, Created: time.Now().Unix()})
	return nil
}

// Update implements the txmonitor.Saver interface. It is called to update the status of the
// transaction once the monitor receives the transaction receipt. For the in-memory store,
// we don't need to update the but rather remove the transaction from the store as we dont
// need to keep track of it anymore. Once we implement a persistent store, we will need to
// update the status of the transaction and keep it in the store for future reference.
func (s *Store) Update(ctx context.Context, txHash common.Hash, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, _ = s.Tree.Delete(txKey(txHash))
	return nil
}

func (s *Store) PendingTxns() ([]*TxnDetails, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	txns := make([]*TxnDetails, 0)
	s.Tree.WalkPrefix(txNS, func(key string, value interface{}) bool {
		txns = append(txns, value.(*TxnDetails))
		return false
	})

	slices.SortFunc(txns, func(a, b *TxnDetails) int {
		return int(a.Created - b.Created)
	})
	return txns, nil
}
