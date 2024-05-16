package transactor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/contracts/transactor"
)

func TestTrasactor(t *testing.T) {
	t.Parallel()

	backend := &testBackend{
		nonce:    5,
		errNonce: 6,
	}
	watcher := &testWatcher{
		allowChan: make(chan uint64),
		txnChan:   make(chan *types.Transaction, 1),
	}
	txnSender := transactor.NewTransactor(backend, watcher)

	nonce, err := txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	if nonce != 5 {
		t.Errorf("expected nonce to be 5, got %d", nonce)
	}

	// If the transaction was not sent, the PendingNonceAt should block until the
	// context is canceled.
	ctx, cancel := context.WithCancel(context.Background())
	errC := make(chan error)
	go func() {
		_, err := txnSender.PendingNonceAt(ctx, common.Address{})
		errC <- err
	}()
	cancel()

	err = <-errC
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	go func() {
		nonce := <-watcher.allowChan
		if nonce != 5 {
			t.Errorf("expected nonce to be 5, got %d", nonce)
		}
	}()

	err = txnSender.SendTransaction(context.Background(), types.NewTransaction(nonce, common.Address{}, nil, 0, nil, nil))
	if err != nil {
		t.Fatal(err)
	}

	select {
	case txn := <-watcher.txnChan:
		if txn.Nonce() != 5 {
			t.Errorf("expected nonce to be 5, got %d", txn.Nonce())
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for transaction")
	}

	nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	if nonce != 6 {
		t.Errorf("expected nonce to be 6, got %d", nonce)
	}

	type nonceResult struct {
		nonce uint64
		err   error
	}
	nonceChan := make(chan nonceResult, 1)
	go func() {
		nonce, err := txnSender.PendingNonceAt(context.Background(), common.Address{})
		nonceChan <- nonceResult{nonce, err}
	}()

	go func() {
		nonce := <-watcher.allowChan
		if nonce != 6 {
			t.Errorf("expected nonce to be 6, got %d", nonce)
		}
	}()
	err = txnSender.SendTransaction(context.Background(), types.NewTransaction(nonce, common.Address{}, nil, 0, nil, nil))
	if err == nil {
		t.Error("expected error, got nil")
	}

	result := <-nonceChan
	if result.err != nil {
		t.Fatal(result.err)
	}

	if result.nonce != 6 {
		t.Errorf("expected nonce to be 6, got %d", result.nonce)
	}

	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	backend.errNonce = 7
	err = txnSender.SendTransaction(ctx, types.NewTransaction(6, common.Address{}, nil, 0, nil, nil))
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	backend.pendingNonceErr = errors.New("nonce error")
	_, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err == nil {
		t.Error("expected error, got nil")
	}

	backend.pendingNonceErr = nil
	nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	if nonce != 6 {
		t.Errorf("expected nonce to be 6, got %d", nonce)
	}
}

type testWatcher struct {
	allowChan chan uint64
	txnChan   chan *types.Transaction
}

func (w *testWatcher) Allow(ctx context.Context, nonce uint64) bool {
	select {
	case <-ctx.Done():
		return false
	case w.allowChan <- nonce:
	}
	return true
}

func (w *testWatcher) Sent(ctx context.Context, tx *types.Transaction) {
	w.txnChan <- tx
}

type testBackend struct {
	bind.ContractTransactor
	nonce           uint64
	errNonce        uint64
	pendingNonceErr error
}

func (b *testBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if b.pendingNonceErr != nil {
		return 0, b.pendingNonceErr
	}
	return b.nonce, nil
}

func (b *testBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if b.errNonce == tx.Nonce() {
		return errors.New("nonce error")
	}
	return nil
}
