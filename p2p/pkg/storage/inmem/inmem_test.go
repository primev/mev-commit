package inmemstorage_test

import (
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/storage"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	storagetest "github.com/primev/mev-commit/p2p/pkg/storage/testing"
)

func TestInmemStorage(t *testing.T) {
	storagetest.RunStoreTests(t, func() storage.Storage {
		return inmemstorage.New()
	})
}
