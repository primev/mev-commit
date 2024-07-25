package inmemstorage

import (
	"github.com/armon/go-radix"
	"github.com/primev/mev-commit/p2p/pkg/storage"
)

type inmemStorage struct {
	Tree *radix.Tree
}

func New() *inmemStorage {
	return &inmemStorage{
		Tree: radix.New(),
	}
}

func (s *inmemStorage) Close() error {
	return nil
}

func (s *inmemStorage) Get(key string) ([]byte, error) {
	v, found := s.Tree.Get(key)
	if !found {
		return nil, storage.ErrKeyNotFound
	}
	return v.([]byte), nil
}

func (s *inmemStorage) Put(key string, value []byte) error {
	_, _ = s.Tree.Insert(key, value)
	return nil
}

func (s *inmemStorage) Delete(key string) error {
	_, _ = s.Tree.Delete(key)
	return nil
}

func (s *inmemStorage) DeletePrefix(prefix string) error {
	_ = s.Tree.DeletePrefix(prefix)
	return nil
}

func (s *inmemStorage) WalkPrefix(prefix string, fn func(key string, val []byte) bool) error {
	s.Tree.WalkPrefix(prefix, func(k string, v interface{}) bool {
		return fn(k, v.([]byte))
	})
	return nil
}
