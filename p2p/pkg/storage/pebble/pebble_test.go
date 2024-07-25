package pebblestorage_test

import (
	"io/ioutil"
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/storage"
	pebblestorage "github.com/primev/mev-commit/p2p/pkg/storage/pebble"
	storagetest "github.com/primev/mev-commit/p2p/pkg/storage/testing"
)

func TestPebbleStore(t *testing.T) {
	storagetest.RunStoreTests(t, func() storage.Storage {
		path, err := ioutil.TempDir("", "pebble")
		if err != nil {
			t.Fatal(err)
		}
		st, err := pebblestorage.New(path)
		if err != nil {
			t.Fatal(err)
		}
		return st
	})
}
