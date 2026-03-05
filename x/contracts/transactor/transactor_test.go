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

// autoAllowWatcher is a watcher that always allows and records sent txs.
type autoAllowWatcher struct {
	txnChan chan *types.Transaction
}

func (w *autoAllowWatcher) Allow(_ context.Context, _ uint64) bool {
	return true
}

func (w *autoAllowWatcher) Sent(_ context.Context, tx *types.Transaction) {
	if w.txnChan != nil {
		w.txnChan <- tx
	}
}

func (w *autoAllowWatcher) NonceOverride() <-chan uint64 {
	return make(chan uint64)
}

// nonceOverrideWatcher is a watcher that allows all txs and supports nonce override.
type nonceOverrideWatcher struct {
	txnChan      chan *types.Transaction
	overrideChan chan uint64
}

func (w *nonceOverrideWatcher) Allow(_ context.Context, _ uint64) bool {
	return true
}

func (w *nonceOverrideWatcher) Sent(_ context.Context, tx *types.Transaction) {
	if w.txnChan != nil {
		w.txnChan <- tx
	}
}

func (w *nonceOverrideWatcher) NonceOverride() <-chan uint64 {
	return w.overrideChan
}

func TestNonceOverrideFromMonitor(t *testing.T) {
	t.Parallel()

	backend := &testBackend{
		nonce: 100,
	}
	overrideChan := make(chan uint64, 1)
	watcher := &nonceOverrideWatcher{
		txnChan:      make(chan *types.Transaction, 16),
		overrideChan: overrideChan,
	}
	txnSender := transactor.NewTransactor(backend, watcher)

	// First call: local nonce is 0, chain says 100 → returns 100
	nonce, err := txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	if nonce != 100 {
		t.Fatalf("expected nonce 100, got %d", nonce)
	}

	// Send nonces 100-109 to advance local nonce to 110
	for i := uint64(100); i <= 109; i++ {
		backend.nonce = i + 1
		if i > 100 {
			nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
			if err != nil {
				t.Fatal(err)
			}
		}
		err = txnSender.SendTransaction(context.Background(), types.NewTransaction(nonce, common.Address{}, nil, 0, nil, nil))
		if err != nil {
			t.Fatal(err)
		}
		<-watcher.txnChan
		nonce = i + 1
	}

	// Local nonce is now 110. Simulate monitor detecting stuck state and
	// sending a nonce override to reset to the confirmed nonce (100).
	backend.nonce = 100
	overrideChan <- 100

	nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	if nonce != 100 {
		t.Fatalf("expected nonce to reset to 100 via monitor override, got %d", nonce)
	}
}

func TestNonceOverrideNotTriggeredWithoutSignal(t *testing.T) {
	t.Parallel()

	backend := &testBackend{
		nonce: 100,
	}
	overrideChan := make(chan uint64, 1)
	watcher := &nonceOverrideWatcher{
		txnChan:      make(chan *types.Transaction, 16),
		overrideChan: overrideChan,
	}
	txnSender := transactor.NewTransactor(backend, watcher)

	// Get initial nonce from chain (100)
	nonce, err := txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	// Send txs 100-104 to advance local nonce to 105
	for i := uint64(100); i <= 104; i++ {
		backend.nonce = i + 1
		if i > 100 {
			nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
			if err != nil {
				t.Fatal(err)
			}
		}
		err = txnSender.SendTransaction(context.Background(), types.NewTransaction(nonce, common.Address{}, nil, 0, nil, nil))
		if err != nil {
			t.Fatal(err)
		}
		<-watcher.txnChan
		nonce = i + 1
	}

	// No override signal sent — local nonce should be used even if
	// chain nonce is lower, because the monitor hasn't detected a stuck state.
	// After sending 100-104, local nonce is 105. The last PendingNonceAt in the
	// loop saw chain=105, so transactor internal nonce is 105+1=106 after the
	// last SendTransaction.
	backend.nonce = 100
	nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	if nonce != 106 {
		t.Fatalf("expected local nonce 106 (no override), got %d", nonce)
	}
}

func TestSendTransactionRetriesExhausted(t *testing.T) {
	t.Parallel()

	backend := &testBackend{
		nonce:     10,
		sendTxErr: context.DeadlineExceeded, // all sends timeout
	}
	watcher := &autoAllowWatcher{txnChan: make(chan *types.Transaction, 1)}
	txnSender := transactor.NewTransactor(backend, watcher)

	// Get initial nonce
	nonce, err := txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	if nonce != 10 {
		t.Fatalf("expected nonce 10, got %d", nonce)
	}

	// SendTransaction should fail after exhausting all retries
	err = txnSender.SendTransaction(context.Background(), types.NewTransaction(nonce, common.Address{}, nil, 0, nil, nil))
	if err == nil {
		t.Fatal("expected error when all retries exhausted, got nil")
	}

	// The nonce should NOT have been incremented — it should be reusable.
	// Since the defer puts the nonce back, the next PendingNonceAt should return 10.
	backend.sendTxErr = nil // clear the error so future sends work
	nonce, err = txnSender.PendingNonceAt(context.Background(), common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	if nonce != 10 {
		t.Fatalf("expected nonce to remain 10 after failed retries, got %d", nonce)
	}

	// Verify no transaction was reported as sent
	select {
	case <-watcher.txnChan:
		t.Fatal("watcher.Sent should not have been called for failed transaction")
	default:
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

func (w *testWatcher) NonceOverride() <-chan uint64 {
	return make(chan uint64)
}

type testBackend struct {
	bind.ContractBackend
	nonce           uint64
	errNonce        uint64
	pendingNonceErr error
	sendTxErr       error // if set, SendTransaction always returns this error
}

func (b *testBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if b.pendingNonceErr != nil {
		return 0, b.pendingNonceErr
	}
	return b.nonce, nil
}

func (b *testBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if b.sendTxErr != nil {
		return b.sendTxErr
	}
	if b.errNonce == tx.Nonce() {
		return errors.New("nonce error")
	}
	return nil
}
