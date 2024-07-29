package storagetest

import (
	"errors"
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/storage"
)

func RunStoreTests(t *testing.T, factory func() storage.Storage) {
	t.Helper()

	t.Run("PutAndGet", func(t *testing.T) {
		TestPutAndGet(t, factory())
	})
	t.Run("Delete", func(t *testing.T) {
		TestDelete(t, factory())
	})
	t.Run("DeletePrefix", func(t *testing.T) {
		TestDeletePrefix(t, factory())
	})
	t.Run("WalkPrefix", func(t *testing.T) {
		TestWalkPrefix(t, factory())
	})
	t.Run("BatchOperations", func(t *testing.T) {
		TestBatchOperations(t, factory())
	})
}

func TestPutAndGet(t *testing.T, s storage.Storage) {
	key := "testKey"
	value := []byte("testValue")

	err := s.Put(key, value)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	retrievedValue, err := s.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrievedValue) != string(value) {
		t.Errorf("Expected value %s, got %s", value, retrievedValue)
	}
}

func TestDelete(t *testing.T, s storage.Storage) {
	key := "testKeyToDelete"
	value := []byte("value")

	err := s.Put(key, value)
	if err != nil {
		t.Fatalf("Setup failed, Put returned error: %v", err)
	}

	err = s.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = s.Get(key)
	if !errors.Is(err, storage.ErrKeyNotFound) {
		t.Fatalf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestDeletePrefix(t *testing.T, s storage.Storage) {
	prefix := "prefix/"
	keys := []string{"prefix/key1", "prefix/key2"}
	value := []byte("value")

	// Put keys with the prefix.
	for _, key := range keys {
		err := s.Put(key, value)
		if err != nil {
			t.Fatalf("Setup failed, Put returned error: %v", err)
		}
	}

	// Delete the prefix.
	err := s.DeletePrefix(prefix)
	if err != nil {
		t.Fatalf("DeletePrefix failed: %v", err)
	}

	// Verify keys are gone.
	for _, key := range keys {
		_, err := s.Get(key)
		if !errors.Is(err, storage.ErrKeyNotFound) {
			t.Fatalf("Expected ErrKeyNotFound for key %s, got %v", key, err)
		}
	}
}

func TestWalkPrefix(t *testing.T, s storage.Storage) {
	prefix := "walkPrefix/"
	keys := map[string][]byte{
		"walkPrefix/key1": []byte("value1"),
		"walkPrefix/key2": []byte("value2"),
	}

	// Put keys with the prefix.
	for key, value := range keys {
		err := s.Put(key, value)
		if err != nil {
			t.Fatalf("Setup failed, Put returned error: %v", err)
		}
	}

	// Walk the prefix.
	foundKeys := make(map[string]bool)
	err := s.WalkPrefix(prefix, func(key string, val []byte) bool {
		if expectedVal, ok := keys[key]; ok && string(val) == string(expectedVal) {
			foundKeys[key] = true
		}
		return false
	})
	if err != nil {
		t.Fatalf("WalkPrefix failed: %v", err)
	}

	// Verify all keys were found.
	if len(foundKeys) != len(keys) {
		t.Fatalf("WalkPrefix did not find all keys, found: %v", foundKeys)
	}
}

func TestBatchOperations(t *testing.T, s storage.Storage) {
	if batcher, ok := s.(storage.Batcher); ok {
		batch := batcher.Batch()
		key := "batchKey"
		value := []byte("batchValue")

		// Put in batch.
		err := batch.Put(key, value)
		if err != nil {
			t.Fatalf("Batch Put failed: %v", err)
		}

		// Write batch.
		err = batch.Write()
		if err != nil {
			t.Fatalf("Batch Write failed: %v", err)
		}

		// Verify key exists.
		retrievedValue, err := s.Get(key)
		if err != nil {
			t.Fatalf("Get after batch write failed: %v", err)
		}
		if string(retrievedValue) != string(value) {
			t.Errorf("Expected value %s, got %s", value, retrievedValue)
		}

		// Reset and verify batch is empty or reset functionality works as expected.
		// This might depend on the implementation details of the batch.
		batch.Reset()
	} else {
		t.Skip("Storage does not implement Batcher interface")
	}
}
