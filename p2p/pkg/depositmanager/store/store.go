package store

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/storage"
)

const (
	balanceNS = "bbs/"
)

var (
	balanceKey = func(window *big.Int, bidder common.Address) string {
		return fmt.Sprintf("%s%s/%s", balanceNS, window, bidder)
	}
	blockBalanceKey = func(window *big.Int, bidder common.Address, blockNumber int64) string {
		return fmt.Sprintf("%s%s/%s/%d", balanceNS, window, bidder, blockNumber)
	}
	balancePrefix = func(window *big.Int) string {
		return fmt.Sprintf("%s%s", balanceNS, window)
	}
)

type Store struct {
	mu sync.RWMutex
	st storage.Storage
}

func New(st storage.Storage) *Store {
	return &Store{
		st: st,
	}
}

func (s *Store) SetBalance(bidder common.Address, windowNumber, depositedAmount *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(balanceKey(windowNumber, bidder), depositedAmount.Bytes())
}

func (s *Store) GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, err := s.st.Get(balanceKey(windowNumber, bidder))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return new(big.Int).SetBytes(val), nil
}

func (s *Store) ClearBalances(windowNumber *big.Int) ([]*big.Int, error) {
	if windowNumber == nil || windowNumber.Cmp(big.NewInt(0)) == -1 {
		return nil, nil
	}

	s.mu.RLock()
	windows := make([]*big.Int, 0)
	err := s.st.WalkPrefix(balanceNS, func(key string, _ []byte) bool {
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
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	for _, w := range windows {
		err := s.st.DeletePrefix(balancePrefix(w))
		if err != nil {
			s.mu.Unlock()
			return nil, err
		}
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

	val, err := s.st.Get(blockBalanceKey(window, bidder, blockNumber))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return new(big.Int).SetBytes(val), nil
}

func (s *Store) SetBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(blockBalanceKey(window, bidder, blockNumber), amount.Bytes())
}

func (s *Store) RefundBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.st.Get(blockBalanceKey(window, bidder, blockNumber))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return s.st.Put(blockBalanceKey(window, bidder, blockNumber), amount.Bytes())
	case err != nil:
		return err
	}

	newAmount := new(big.Int).Add(new(big.Int).SetBytes(val), amount)
	return s.st.Put(blockBalanceKey(window, bidder, blockNumber), newAmount.Bytes())
}

func (s *Store) BalanceEntries(windowNumber *big.Int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := 0
	prefix := balancePrefix(windowNumber)
	err := s.st.WalkPrefix(prefix, func(key string, val []byte) bool {
		entries++
		return false
	})
	if err != nil {
		return 0, err
	}

	return entries, nil
}
