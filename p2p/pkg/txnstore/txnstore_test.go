package txnstore_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	"github.com/primev/mev-commit/p2p/pkg/txnstore"
)

func TestStore_Save(t *testing.T) {
	st := inmem.New()
	store := txnstore.New(st)

	txHash := common.HexToHash("0x1234567890abcdef")
	nonce := uint64(123)
	err := store.Save(context.Background(), txHash, nonce)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that the transaction details are saved in the store
	var txnDetails txnstore.TxnDetails
	buf, err := st.Get(txnstore.TxKey(txHash))
	if err != nil {
		t.Errorf("failed to get transaction details from store: %v", err)
	}
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&txnDetails); err != nil {
		t.Errorf("failed to decode transaction details: %v", err)
	}
	if txnDetails.Hash != txHash {
		t.Errorf("unexpected transaction hash: got %s, want %s", txnDetails.Hash.Hex(), txHash.Hex())
	}
	if txnDetails.Nonce != nonce {
		t.Errorf("unexpected nonce: got %d, want %d", txnDetails.Nonce, nonce)
	}
}

func TestStore_Update(t *testing.T) {
	st := inmem.New()
	store := txnstore.New(st)

	txHash := common.HexToHash("0x1234567890abcdef")
	nonce := uint64(123)
	err := store.Save(context.Background(), txHash, nonce)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = store.Update(context.Background(), txHash, "completed")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that the transaction details are removed from the store
	_, err = st.Get(txnstore.TxKey(txHash))
	if err != storage.ErrKeyNotFound {
		t.Errorf("expected transaction details to be removed from store, got: %v", err)
	}
}

func TestStore_PendingTxns(t *testing.T) {
	st := inmem.New()
	store := txnstore.New(st)

	// Save some pending transactions in the store
	txHash1 := common.HexToHash("0x1234567890abcdef")
	nonce1 := uint64(123)
	err := store.Save(context.Background(), txHash1, nonce1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	txHash2 := common.HexToHash("0xabcdef1234567890")
	nonce2 := uint64(456)
	err = store.Save(context.Background(), txHash2, nonce2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Retrieve the pending transactions from the store
	pendingTxns, err := store.PendingTxns()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify the number of pending transactions
	if len(pendingTxns) != 2 {
		t.Errorf("unexpected number of pending transactions: got %d, want %d", len(pendingTxns), 2)
	}

	// Verify the details of the pending transactions
	if pendingTxns[0].Hash != txHash1 {
		t.Errorf("unexpected transaction hash: got %s, want %s", pendingTxns[0].Hash.Hex(), txHash1.Hex())
	}
	if pendingTxns[0].Nonce != nonce1 {
		t.Errorf("unexpected nonce: got %d, want %d", pendingTxns[0].Nonce, nonce1)
	}
	if pendingTxns[1].Hash != txHash2 {
		t.Errorf("unexpected transaction hash: got %s, want %s", pendingTxns[1].Hash.Hex(), txHash2.Hex())
	}
	if pendingTxns[1].Nonce != nonce2 {
		t.Errorf("unexpected nonce: got %d, want %d", pendingTxns[1].Nonce, nonce2)
	}
}

func TestStore_PendingTxns_SortedByCreated(t *testing.T) {
	st := inmem.New()
	store := txnstore.New(st)

	// Save some pending transactions in the store with different creation times
	txHash1 := common.HexToHash("0x1234567890abcdef")
	nonce1 := uint64(123)
	err := store.Save(context.Background(), txHash1, nonce1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	txHash2 := common.HexToHash("0xabcdef1234567890")
	nonce2 := uint64(456)
	err = store.Save(context.Background(), txHash2, nonce2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Retrieve the pending transactions from the store
	pendingTxns, err := store.PendingTxns()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that the pending transactions are sorted by creation time
	if pendingTxns[0].Created > pendingTxns[1].Created {
		t.Errorf("pending transactions are not sorted by creation time")
	}
}
