package transactor

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	txnRetriesLimit = 3

	// defaultMaxNonceDrift is the maximum allowed difference between the
	// local nonce counter and the chain's pending nonce before the transactor
	// resets to the chain nonce. Since sends are serialized via a size-1
	// channel, under normal operation the drift should be 0. A drift larger
	// than this threshold indicates that previously sent transactions were
	// lost (e.g. mempool cleared after a node restart).
	defaultMaxNonceDrift uint64 = 5
)

// Watcher is an interface that is used to manage the lifecycle of a transaction.
// The Allow method is used to determine if a transaction should be sent. The context
// is passed to the method so that the watcher can determine this based on the context.
// The Sent method is is used to notify the watcher that the transaction has been sent.
type Watcher interface {
	Allow(ctx context.Context, nonce uint64) bool
	Sent(ctx context.Context, tx *types.Transaction)
}

// Transactor is a wrapper around a bind.ContractBackend that ensures that
// transactions are sent in nonce order and that the nonce is updated correctly.
// It also uses rate-limiting to ensure that the transactions are sent at a
// reasonable rate. The Watcher is used to manage the tx lifecycle. It is used to
// determine if a transaction should be sent and to notify the watcher when a
// transaction is sent.
// The purpose of this type is to use the abi generated code to interact with the
// contract and to manage the nonce and rate-limiting. To understand the synchronization
// better, the abi generated code calls the PendingNonceAt method to get the nonce
// and then calls SendTransaction to send the transaction. The PendingNonceAt method
// and the SendTransaction method are both called in the same goroutine. So this ensures
// that the nonce is updated correctly and that the transactions are sent in order. In case
// of an error, the nonce is put back into the channel so that it can be reused.
type Transactor struct {
	bind.ContractBackend
	nonceChan     chan uint64
	watcher       Watcher
	maxNonceDrift uint64
	logger        *slog.Logger
}

// Option is a functional option for configuring the Transactor.
type Option func(*Transactor)

// WithMaxNonceDrift sets the maximum allowed drift between the local nonce
// counter and the chain's pending nonce before the transactor self-heals by
// resetting to the chain nonce.
func WithMaxNonceDrift(n uint64) Option {
	return func(t *Transactor) {
		t.maxNonceDrift = n
	}
}

// WithLogger sets the logger for the transactor.
func WithLogger(l *slog.Logger) Option {
	return func(t *Transactor) {
		t.logger = l
	}
}

func NewTransactor(
	backend bind.ContractBackend,
	watcher Watcher,
	opts ...Option,
) *Transactor {
	nonceChan := make(chan uint64, 1)
	// We need to send a value to the channel so that the first transaction
	// can be sent. The value is not important as the first transaction will
	// get the nonce from the blockchain.
	nonceChan <- 0
	t := &Transactor{
		ContractBackend: backend,
		watcher:         watcher,
		nonceChan:       nonceChan,
		maxNonceDrift:   defaultMaxNonceDrift,
		logger:          slog.Default(),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Transactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case nonce := <-t.nonceChan:
		pendingNonce, err := t.ContractBackend.PendingNonceAt(ctx, account)
		if err != nil {
			// this naked write is safe as only the SendTransaction writes to
			// the channel. The goroutine which is trying to send the transaction
			// won't be calling SendTransaction if there is an error here.
			t.nonceChan <- nonce
			return 0, err
		}
		if pendingNonce > nonce {
			return pendingNonce, nil
		}
		// Self-healing: if local nonce drifts too far ahead of the chain's
		// pending nonce, transactions were likely lost (e.g. mempool cleared
		// after a node restart). Reset to chain nonce to close the gap.
		// This does not add any extra RPC calls — we already query the
		// chain's pending nonce above.
		if nonce > pendingNonce+t.maxNonceDrift {
			t.logger.Warn("nonce drift detected, resetting to chain pending nonce",
				"localNonce", nonce,
				"chainPendingNonce", pendingNonce,
				"drift", nonce-pendingNonce,
			)
			return pendingNonce, nil
		}
		return nonce, nil
	}
}

func (t *Transactor) SendTransaction(ctx context.Context, tx *types.Transaction) (retErr error) {
	defer func() {
		if retErr != nil {
			// If the transaction fails, we need to put the nonce back into the channel
			// so that it can be reused.
			t.nonceChan <- tx.Nonce()
		}
	}()

	if !t.watcher.Allow(ctx, tx.Nonce()) {
		return ctx.Err()
	}

	delay := 1 * time.Second
	for tries := 0; tries <= txnRetriesLimit; tries++ {
		cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := t.ContractBackend.SendTransaction(cctx, tx)
		cancel()

		if err != nil {
			if err == context.DeadlineExceeded {
				delay *= 2
				retryTimer := time.NewTimer(delay)
				select {
				case <-ctx.Done():
					retryTimer.Stop()
					return ctx.Err()
				case <-retryTimer.C:
				}
				continue
			}
			return err
		}
		break
	}

	// If the transaction is successful, we need to update the nonce and notify the
	// watcher.
	t.watcher.Sent(ctx, tx)
	t.nonceChan <- tx.Nonce() + 1

	return nil
}
