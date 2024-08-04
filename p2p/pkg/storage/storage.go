package storage

import (
	"errors"
	"io"
)

var (
	// ErrKeyNotFound is returned when the key is not found.
	ErrKeyNotFound = errors.New("key not found")
)

type Reader interface {
	// Get returns the value for the given key.
	Get(key string) ([]byte, error)
	// WalkPrefix walks the values for the given prefix.
	WalkPrefix(prefix string, fn func(key string, val []byte) bool) error
}

type Writer interface {
	// Put puts the value for the given key.
	Put(key string, value []byte) error
	// Delete deletes the value for the given key.
	Delete(key string) error
	// DeletePrefix deletes all the values for the given prefix.
	DeletePrefix(prefix string) error
}

type Storage interface {
	Reader
	Writer

	io.Closer
}

type Batcher interface {
	// Batch returns a new batch.
	Batch() Batch
}

type Batch interface {
	Writer

	// Write writes the batch to the storage.
	Write() error
	// Reset resets the batch.
	Reset()
}
