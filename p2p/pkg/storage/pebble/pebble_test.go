package pebblestorage_test

import (
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/storage"
	pebblestorage "github.com/primev/mev-commit/p2p/pkg/storage/pebble"
	storagetest "github.com/primev/mev-commit/p2p/pkg/storage/testing"
)

func TestPebbleStore(t *testing.T) {
	storagetest.RunStoreTests(t, func() storage.Storage {
		st, err := pebblestorage.New(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		return st
	})
}
