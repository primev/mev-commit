package store

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	balanceNS = "bbs/"
)

var (
	balanceKey = func(bidder common.Address, provider common.Address) string {
		return fmt.Sprintf("%s%s/%s", balanceNS, bidder, provider)
	}
	balancePrefix = func(bidder common.Address) string {
		return fmt.Sprintf("%s%s", balanceNS, bidder)
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

func (s *Store) SetBalance(bidder common.Address, provider common.Address, depositedAmount *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(balanceKey(bidder, provider), depositedAmount.Bytes())
}

func (s *Store) GetBalance(bidder common.Address, provider common.Address) (*big.Int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, err := s.st.Get(balanceKey(bidder, provider))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return new(big.Int).SetBytes(val), nil
}

func (s *Store) DeleteBalance(bidder common.Address, provider common.Address) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.st.Delete(balanceKey(bidder, provider))
}

func (s *Store) RefundBalanceIfExists(
	bidder common.Address,
	provider common.Address,
	amount *big.Int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.st.Get(balanceKey(bidder, provider))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return status.Errorf(codes.FailedPrecondition, "balance not found, no increase needed")
	case err != nil:
		return err
	}

	newAmount := new(big.Int).Add(new(big.Int).SetBytes(val), amount)
	return s.st.Put(balanceKey(bidder, provider), newAmount.Bytes())
}

func (s *Store) DeductBalanceIfExists(
	bidder common.Address,
	provider common.Address,
	amount *big.Int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.st.Get(balanceKey(bidder, provider))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return status.Errorf(codes.FailedPrecondition, "balance not found, no decrease needed")
	case err != nil:
		return err
	}

	newBalance := new(big.Int).Sub(new(big.Int).SetBytes(val), amount)
	if newBalance.Cmp(big.NewInt(0)) < 0 {
		return status.Errorf(codes.Internal, "balance cannot be decreased below 0")
	}
	return s.st.Put(balanceKey(bidder, provider), newBalance.Bytes())
}

func (s *Store) BalanceEntries(bidder common.Address) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := 0
	prefix := balancePrefix(bidder)
	err := s.st.WalkPrefix(prefix, func(key string, val []byte) bool {
		entries++
		return false
	})
	if err != nil {
		return 0, err
	}

	return entries, nil
}
