package txnstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
)

const (
	txNS = "tx/"
)

var (
	txKey = func(txHash common.Hash) string {
		return fmt.Sprintf("%s%s", txNS, txHash.Hex())
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

// Following are the methods to save and update the transaction details in the store.
// The store is used to keep track of the transactions that are sent to the blockchain and
// have not yet received the transaction receipt. These are used by the debug service
// to show the pending transactions and cancel them if needed. The store hooks up
// to the txmonitor package which allows a component to get notified when the transaction
// is sent to the blockchain and when the transaction receipt is received. As of now,
// we do not persist the transaction details, so the update method is used to remove the
// transaction from the store. This will no longer be seen in the pending transactions list.
// Save implements the txmonitor.Saver interface. It saves the transaction hash and nonce in the store once
// the transaction is sent to the blockchain.
func (s *Store) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(&txmonitor.TxnDetails{Hash: txHash, Nonce: nonce, Created: time.Now().Unix()}); err != nil {
		return err
	}

	return s.st.Put(txKey(txHash), b.Bytes())
}

// Update implements the txmonitor.Saver interface. It is called to update the status of the
// transaction once the monitor receives the transaction receipt. For the in-memory store,
// we don't need to update the but rather remove the transaction from the store as we dont
// need to keep track of it anymore. Once we implement a persistent store, we will need to
// update the status of the transaction and keep it in the store for future reference.
func (s *Store) Update(ctx context.Context, txHash common.Hash, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Delete(txKey(txHash))
}

func (s *Store) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	txns := make([]*txmonitor.TxnDetails, 0)
	err := s.st.WalkPrefix(txNS, func(key string, value []byte) bool {
		txn := new(txmonitor.TxnDetails)
		if err := gob.NewDecoder(bytes.NewReader(value)).Decode(txn); err != nil {
			return false
		}
		txns = append(txns, txn)
		return false
	})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(txns, func(a, b *txmonitor.TxnDetails) int {
		return int(a.Created - b.Created)
	})
	return txns, nil
}
