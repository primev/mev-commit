package txmonitor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	batchSize int = 64
)

var (
	ErrTxnCancelled  = errors.New("transaction was cancelled")
	ErrTxnFailed     = errors.New("transaction failed")
	ErrMonitorClosed = errors.New("monitor was closed")
)

type TxnDetails struct {
	Hash    common.Hash
	Nonce   uint64
	Created int64
}

type Saver interface {
	Save(ctx context.Context, txHash common.Hash, nonce uint64) error
	Update(ctx context.Context, txHash common.Hash, status string) error
	PendingTxns() ([]*TxnDetails, error)
}

type EVMHelper interface {
	BatchReceiptGetter
	Debugger
}

type EVM interface {
	BlockNumber(ctx context.Context) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}

type waitCheck struct {
	nonce uint64
	block uint64
}

// Monitor is a transaction monitor that keeps track of the transactions sent by the owner
// and waits for the receipt to be available to query them in batches. The monitor also
// ensures rate-limiting by waiting for the nonce to be confirmed by the client before
// allowing the transaction to be sent. The alternative is that each client keeps
// polling the backend for receipts which is inefficient.
type Monitor struct {
	owner              common.Address
	mtx                sync.Mutex
	waitMap            map[uint64]map[common.Hash][]chan Result
	client             EVM
	helper             EVMHelper
	saver              Saver
	newTxAdded         chan struct{}
	nonceUpdate        chan struct{}
	blockUpdate        chan waitCheck
	logger             *slog.Logger
	lastConfirmedNonce atomic.Uint64
	maxPendingTxs      uint64
	metrics            *metrics
}

func New(
	owner common.Address,
	client EVM,
	helper EVMHelper,
	saver Saver,
	logger *slog.Logger,
	maxPendingTxs uint64,
) *Monitor {
	if saver == nil {
		saver = noopSaver{}
	}
	m := &Monitor{
		owner:         owner,
		client:        client,
		logger:        logger,
		helper:        helper,
		saver:         saver,
		maxPendingTxs: maxPendingTxs,
		metrics:       newMetrics(),
		waitMap:       make(map[uint64]map[common.Hash][]chan Result),
		newTxAdded:    make(chan struct{}),
		nonceUpdate:   make(chan struct{}),
		blockUpdate:   make(chan waitCheck),
	}

	pending, err := saver.PendingTxns()
	if err != nil {
		logger.Error("failed to get pending transactions", "err", err)
	}

	for _, txn := range pending {
		m.WatchTx(txn.Hash, txn.Nonce)
	}

	return m
}

func (m *Monitor) Metrics() []prometheus.Collector {
	return m.metrics.Metrics()
}

func (m *Monitor) Start(ctx context.Context) <-chan struct{} {
	wg := sync.WaitGroup{}
	done := make(chan struct{})

	wg.Add(2)
	go func() {
		defer wg.Done()

		queryTicker := time.NewTicker(500 * time.Millisecond)
		defer queryTicker.Stop()

		defer func() {
			m.mtx.Lock()
			defer m.mtx.Unlock()

			for _, v := range m.waitMap {
				for _, c := range v {
					for _, c := range c {
						c <- Result{nil, ErrMonitorClosed}
						close(c)
					}
				}
			}
		}()

		m.logger.Info("monitor started")
		lastBlock := uint64(0)
		for {
			newTx := false
			select {
			case <-ctx.Done():
				return
			case <-m.newTxAdded:
				newTx = true
			case <-queryTicker.C:
			}

			currentBlock, err := m.client.BlockNumber(ctx)
			if err != nil {
				m.logger.Error("failed to get block number", "err", err)
				continue
			}

			if currentBlock <= lastBlock && !newTx {
				continue
			}

			lastNonce, err := m.client.NonceAt(
				ctx,
				m.owner,
				new(big.Int).SetUint64(currentBlock),
			)
			if err != nil {
				m.logger.Error("failed to get nonce", "err", err)
				continue
			}

			m.lastConfirmedNonce.Store(lastNonce)
			m.metrics.lastConfirmedNonce.Set(float64(lastNonce))
			m.metrics.lastBlockNumber.Set(float64(currentBlock))
			m.triggerNonceUpdate()

			select {
			case m.blockUpdate <- waitCheck{lastNonce, currentBlock}:
			default:
			}
			m.logger.Debug("checking for receipts", "block", currentBlock, "lastNonce", lastNonce)
			lastBlock = currentBlock
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case check := <-m.blockUpdate:
				m.check(ctx, check.block, check.nonce)
			}
		}
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

// Allow waits until a sufficiently high nonce is confirmed by the client. If
// the context is cancelled, the function returns false.
func (m *Monitor) Allow(ctx context.Context, nonce uint64) bool {
	for {
		if nonce <= m.lastConfirmedNonce.Load()+m.maxPendingTxs {
			return true
		}
		select {
		case <-ctx.Done():
			return false
		case <-m.nonceUpdate:
		}
	}
}

// Sent saves the transaction and starts monitoring it for the receipt. The Saver
// is used to save the transaction details in the storage backend of choice.
func (m *Monitor) Sent(ctx context.Context, tx *types.Transaction) {
	if err := m.saver.Save(ctx, tx.Hash(), tx.Nonce()); err != nil {
		m.logger.Error("failed to save transaction", "err", err)
	}

	m.metrics.lastUsedNonce.Set(float64(tx.Nonce()))
	m.metrics.lastUsedGas.Set(float64(tx.Gas()))
	m.metrics.lastUsedGasPrice.Set(float64(tx.GasPrice().Int64()))
	m.metrics.lastUsedGasTip.Set(float64(tx.GasTipCap().Int64()))

	res := m.WatchTx(tx.Hash(), tx.Nonce())
	go func() {
		r := <-res
		status := "success"
		if r.Err != nil {
			status = fmt.Sprintf("failed: %v", r.Err)
		}
		if err := m.saver.Update(context.Background(), tx.Hash(), status); err != nil {
			m.logger.Error("failed to update transaction", "err", err)
		}
		m.logger.Info("transaction status",
			"txHash", tx.Hash(),
			"txStatus", status,
		)
	}()
}

func (m *Monitor) WatchTx(txHash common.Hash, nonce uint64) <-chan Result {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.waitMap[nonce] == nil {
		m.waitMap[nonce] = make(map[common.Hash][]chan Result)
	}

	c := make(chan Result, 1)
	m.waitMap[nonce][txHash] = append(m.waitMap[nonce][txHash], c)

	m.triggerNewTx()
	return c
}

func (m *Monitor) WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	res := m.WatchTx(tx.Hash(), tx.Nonce())
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-res:
		return r.Receipt, r.Err
	}
}

func (m *Monitor) triggerNewTx() {
	select {
	case m.newTxAdded <- struct{}{}:
	default:
	}
}

func (m *Monitor) triggerNonceUpdate() {
	select {
	case m.nonceUpdate <- struct{}{}:
	default:
	}
}

// waitMap holds the transactions that are waiting for the receipt with the nonce as the key.
// The passed nonce is the last nonce that was confirmed by the client, so any transactions
// with a nonce less than this value are supposed to be confirmed and waiting for the receipt.
func (m *Monitor) getOlderTxns(nonce uint64) map[uint64][]common.Hash {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	txnMap := make(map[uint64][]common.Hash)
	for k, v := range m.waitMap {
		if k >= nonce {
			continue
		}

		for h := range v {
			txnMap[k] = append(txnMap[k], h)
		}
	}

	return txnMap
}

func (m *Monitor) notify(
	nonce uint64,
	txn common.Hash,
	res Result,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	waiters := 0
	for _, c := range m.waitMap[nonce][txn] {
		c <- res
		waiters++
		close(c)
	}
	delete(m.waitMap[nonce], txn)
	if len(m.waitMap[nonce]) == 0 {
		delete(m.waitMap, nonce)
	}
}

// check retrieves the receipts for the transactions with nonce less than the lastNonce
// and notifies the waiting clients.
func (m *Monitor) check(ctx context.Context, newBlock uint64, lastNonce uint64) {
	checkTxns := m.getOlderTxns(lastNonce)
	nonceMap := make(map[common.Hash]uint64)

	if len(checkTxns) == 0 {
		return
	}

	txHashes := make([]common.Hash, 0, len(checkTxns))
	for n, txns := range checkTxns {
		for _, txn := range txns {
			txHashes = append(txHashes, txn)
			nonceMap[txn] = n
		}
	}

	for start := 0; start < len(txHashes); start += batchSize {
		end := start + batchSize
		if end > len(txHashes) {
			end = len(txHashes)
		}

		receipts, err := m.helper.BatchReceipts(ctx, txHashes[start:end])
		if err != nil {
			m.logger.Error("failed to get receipts", "err", err)
			return
		}

		for i, r := range receipts {
			nonce := nonceMap[txHashes[start+i]]
			if r.Err != nil {
				if errors.Is(r.Err, ethereum.NotFound) {
					m.notify(nonce, txHashes[start+i], Result{nil, ErrTxnCancelled})
					continue
				}
				m.logger.Error("failed to get receipt", "error", r.Err, "txHash", txHashes[start+i])
				continue
			}
			if r.Receipt.Status != types.ReceiptStatusSuccessful {
				// reason, err := m.helper.RevertReason(ctx, r.Receipt, m.owner)
				// if err != nil {
				// 	m.logger.Error(
				// 		"retrieving transaction revert reason failed",
				// 		"error", err,
				// 		"txHash", txHashes[start+i],
				// 	)
				// 	reason = "unknown"
				// }
				reason := "unknown"
				m.logger.Error(
					"failed transaction",
					"txHash", txHashes[start+i],
					"status", r.Receipt.Status,
					"reason", reason,
				)
				m.notify(nonce, txHashes[start+i], Result{r.Receipt, fmt.Errorf("%w: %v", ErrTxnFailed, reason)})
				continue
			}

			m.notify(nonce, txHashes[start+i], Result{r.Receipt, nil})
		}
	}
}

type noopSaver struct{}

func (noopSaver) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	return nil
}

func (noopSaver) Update(ctx context.Context, txHash common.Hash, status string) error {
	return nil
}

func (noopSaver) PendingTxns() ([]*TxnDetails, error) {
	return nil, nil
}
