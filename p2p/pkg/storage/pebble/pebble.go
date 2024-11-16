package pebblestorage

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/primev/mev-commit/p2p/pkg/storage"
)

type pebbleStorage struct {
	db *pebble.DB
}

func New(path string) (*pebbleStorage, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &pebbleStorage{
		db: db,
	}, nil
}

func (s *pebbleStorage) Close() error {
	return errors.Join(s.db.Flush(), s.db.Close())
}

func (s *pebbleStorage) Get(key string) ([]byte, error) {
	buf, closer, err := s.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, errors.Join(storage.ErrKeyNotFound, err)
		}
		return nil, err
	}

	val := make([]byte, len(buf))
	copy(val, buf)

	_ = closer.Close()

	return val, nil
}

func (s *pebbleStorage) WalkPrefix(prefix string, fn func(key string, val []byte) bool) error {
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: []byte(prefix),
		UpperBound: upperBound([]byte(prefix)),
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if fn(string(iter.Key()), iter.Value()) {
			break
		}
	}
	return nil
}

func (s *pebbleStorage) Put(key string, value []byte) error {
	return s.db.Set([]byte(key), value, pebble.NoSync)
}

func (s *pebbleStorage) Delete(key string) error {
	return s.db.Delete([]byte(key), pebble.NoSync)
}

func (s *pebbleStorage) DeletePrefix(prefix string) error {
	batch := s.db.NewBatch()
	defer batch.Close()

	err := batch.DeleteRange([]byte(prefix), upperBound([]byte(prefix)), nil)
	if err != nil {
		return err
	}

	return batch.Commit(pebble.NoSync)
}

func (s *pebbleStorage) Batch() storage.Batch {
	return &pebbleBatch{
		batch: s.db.NewBatch(),
	}
}

type pebbleBatch struct {
	batch *pebble.Batch
}

func (b *pebbleBatch) Put(key string, value []byte) error {
	return b.batch.Set([]byte(key), value, nil)
}

func (b *pebbleBatch) Delete(key string) error {
	return b.batch.Delete([]byte(key), nil)
}

func (b *pebbleBatch) DeletePrefix(prefix string) error {
	return b.batch.DeleteRange([]byte(prefix), upperBound([]byte(prefix)), pebble.NoSync)
}

func (b *pebbleBatch) Write() error {
	return b.batch.Commit(pebble.NoSync)
}

func (b *pebbleBatch) Reset() {
	b.batch.Reset()
}
func upperBound(prefix []byte) []byte {
	// if the prefix is 0x01..., we want 0x02 as an upper bound.
	// if the prefix is 0x0000ff..., we want 0x0001 as an upper bound.
	// if the prefix is 0x0000ff01..., we want 0x0000ff02 as an upper bound.
	// if the prefix is 0xffffff..., we don't want an upper bound.
	// if the prefix is 0xff..., we don't want an upper bound.
	// if the prefix is empty, we don't want an upper bound.
	// basically, we want to find the last byte that can be lexicographically incremented.
	var upper []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		b := prefix[i]
		if b == 0xff {
			continue
		}
		upper = make([]byte, i+1)
		copy(upper, prefix)
		upper[i] = b + 1
		break
	}
	return upper
}
