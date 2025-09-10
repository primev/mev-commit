package follower

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"

	"github.com/primev/mev-commit/p2p/pkg/storage"
)

type Store struct {
	logger *slog.Logger
	kv     storage.Storage
}

const (
	progressKey = "follower/last_block"
)

func NewStore(logger *slog.Logger, kv storage.Storage) *Store {
	return &Store{
		logger: logger.With("component", "FollowerStore"),
		kv:     kv,
	}
}

func (s *Store) GetLastProcessed() (uint64, error) {
	if s.kv == nil {
		return 0, errors.New("kv is nil")
	}
	buf, err := s.kv.Get(progressKey)
	switch {
	case err == nil:
		if len(buf) != 8 {
			return 0, fmt.Errorf("invalid %q length: got %d, want 8", progressKey, len(buf))
		}
		return binary.BigEndian.Uint64(buf), nil
	case errors.Is(err, storage.ErrKeyNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (s *Store) SetLastProcessed(height uint64) error {
	if s.kv == nil {
		return errors.New("kv is nil")
	}
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], height)
	return s.kv.Put(progressKey, b[:])
}

func (s *Store) Close() error {
	if s.kv != nil {
		return s.kv.Close()
	}
	return nil
}
