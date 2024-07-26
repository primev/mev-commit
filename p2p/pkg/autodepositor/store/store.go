package store

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/primev/mev-commit/p2p/pkg/storage"
)

const (
	// local deposit entries
	depositNS = "dep/"
)

var (
	depositKey = func(window *big.Int) string {
		return fmt.Sprintf("%s%s", depositNS, window)
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

func (s *Store) StoreDeposits(ctx context.Context, window []*big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range window {
		err := s.st.Put(depositKey(w), []byte{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListDeposits(ctx context.Context, till *big.Int) ([]*big.Int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deposits := make([]*big.Int, 0)
	err := s.st.WalkPrefix(depositNS, func(key string, _ []byte) bool {
		parts := strings.Split(key, "/")
		if len(parts) != 2 {
			return false
		}
		w, ok := new(big.Int).SetString(parts[1], 10)
		if !ok {
			return false
		}
		if till == nil || w.Cmp(till) != 1 {
			deposits = append(deposits, w)
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	return deposits, nil
}

func (s *Store) ClearDeposits(ctx context.Context, windows []*big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range windows {
		err := s.st.Delete(depositKey(w))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) IsDepositMade(ctx context.Context, window *big.Int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := s.st.Get(depositKey(window))
	return err == nil
}
