package evmclient

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	maxSentTxs uint64 = 1024
	batchSize  int    = 64
)

var (
	ErrTxnCancelled  = errors.New("transaction was cancelled")
	ErrMonitorClosed = errors.New("monitor was closed")
)

type waitCheck struct {
	nonce uint64
	block uint64
}

type txmonitor struct {
	baseCtx            context.Context
	baseCancel         context.CancelFunc
	owner              common.Address
	mtx                sync.Mutex
	waitMap            map[uint64]map[common.Hash][]chan Result
	client             EVM
	newTxAdded         chan struct{}
	waitDone           chan struct{}
	checkerDone        chan struct{}
	blockUpdate        chan waitCheck
	logger             *slog.Logger
	metrics            *metrics
	lastConfirmedNonce atomic.Uint64
}

func newTxMonitor(
	owner common.Address,
	client EVM,
	logger *slog.Logger,
	m *metrics,
) *txmonitor {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	if m == nil {
		m = newMetrics()
	}
	tm := &txmonitor{
		baseCtx:     baseCtx,
		baseCancel:  baseCancel,
		owner:       owner,
		client:      client,
		logger:      logger,
		metrics:     m,
		waitMap:     make(map[uint64]map[common.Hash][]chan Result),
		newTxAdded:  make(chan struct{}),
		waitDone:    make(chan struct{}),
		checkerDone: make(chan struct{}),
		blockUpdate: make(chan waitCheck),
	}
	go tm.watchLoop()
	go tm.checkLoop()

	return tm
}

type Result struct {
	Receipt *types.Receipt
	Err     error
}

func (t *txmonitor) watchLoop() {
	defer close(t.waitDone)

	queryTicker := time.NewTicker(500 * time.Millisecond)
	defer queryTicker.Stop()

	defer func() {
		t.mtx.Lock()
		defer t.mtx.Unlock()

		for _, v := range t.waitMap {
			for _, c := range v {
				for _, c := range c {
					c <- Result{nil, ErrMonitorClosed}
					close(c)
				}
			}
		}
	}()

	lastBlock := uint64(0)
	for {
		newTx := false
		select {
		case <-t.baseCtx.Done():
			return
		case <-t.newTxAdded:
			newTx = true
		case <-queryTicker.C:
		}

		currentBlock, err := t.client.BlockNumber(t.baseCtx)
		if err != nil {
			t.logger.Error("failed to get block number", "err", err)
			continue
		}

		if currentBlock <= lastBlock && !newTx {
			continue
		}

		t.metrics.CurrentBlockNumber.Set(float64(currentBlock))

		lastNonce, err := t.client.NonceAt(
			t.baseCtx,
			t.owner,
			new(big.Int).SetUint64(currentBlock),
		)
		if err != nil {
			t.logger.Error("failed to get nonce", "err", err)
			continue
		}

		t.lastConfirmedNonce.Store(lastNonce)
		t.metrics.LastConfirmedNonce.Set(float64(lastNonce))

		select {
		case t.blockUpdate <- waitCheck{lastNonce, currentBlock}:
		default:
		}
		lastBlock = currentBlock
	}
}

func (t *txmonitor) checkLoop() {
	defer close(t.checkerDone)

	for {
		select {
		case <-t.baseCtx.Done():
			return
		case check := <-t.blockUpdate:
			t.check(check.block, check.nonce)
		}
	}
}

func (t *txmonitor) Close() error {
	t.baseCancel()
	done := make(chan struct{})
	go func() {
		<-t.checkerDone
		<-t.waitDone

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("failed to close txmonitor")
	}
}

func (t *txmonitor) getOlderTxns(nonce uint64) map[uint64][]common.Hash {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	txnMap := make(map[uint64][]common.Hash)
	for k, v := range t.waitMap {
		if k >= nonce {
			continue
		}

		for h := range v {
			txnMap[k] = append(txnMap[k], h)
		}
	}

	return txnMap
}

func (t *txmonitor) notify(
	nonce uint64,
	txn common.Hash,
	res Result,
) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	waiters := 0
	for _, c := range t.waitMap[nonce][txn] {
		c <- res
		waiters++
		close(c)
	}
	delete(t.waitMap[nonce], txn)
	if len(t.waitMap[nonce]) == 0 {
		delete(t.waitMap, nonce)
	}
}

func (t *txmonitor) check(newBlock uint64, lastNonce uint64) {
	checkTxns := t.getOlderTxns(lastNonce)
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

		// Prepare the batch
		batch := make([]rpc.BatchElem, end-start)
		for i, hash := range txHashes[start:end] {
			batch[i] = rpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{hash},
				Result: new(types.Receipt),
			}
		}

		opStart := time.Now()
		// Execute the batch request
		err := t.client.Batcher().BatchCallContext(t.baseCtx, batch)
		if err != nil {
			t.logger.Error("failed to execute batch call", "err", err)
			return
		}
		t.metrics.GetReceiptBatchOperationTimeMs.Set(float64(time.Since(opStart).Milliseconds()))

		// Process the responses
		for i, result := range batch {
			tHash := txHashes[start+i]
			nonce := nonceMap[tHash]
			if result.Error != nil {
				if errors.Is(result.Error, ethereum.NotFound) {
					t.notify(nonce, tHash, Result{nil, ErrTxnCancelled})
					continue
				}
				var tt *TransactionTrace
				if dbg, ok := t.client.(Debugger); ok {
					if tt, err = dbg.TraceTransaction(t.baseCtx, tHash); err != nil {
						t.logger.Error("retrieving transaction trace failed", "error", err)
					}
				}
				t.logger.Error("failed to get receipt", "error", result.Error, "transaction_trace", tt)
				continue
			}
			if result.Result == nil {
				continue
			}
			t.notify(nonce, tHash, Result{result.Result.(*types.Receipt), nil})
		}
	}
}

func (t *txmonitor) allowNonce(nonce uint64) bool {
	return nonce <= t.lastConfirmedNonce.Load()+maxSentTxs
}

func (t *txmonitor) watchTx(txHash common.Hash, nonce uint64) (<-chan Result, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if t.waitMap[nonce] == nil {
		t.waitMap[nonce] = make(map[common.Hash][]chan Result)
	}

	c := make(chan Result, 1)
	t.waitMap[nonce][txHash] = append(t.waitMap[nonce][txHash], c)

	select {
	case t.newTxAdded <- struct{}{}:
	default:
	}
	return c, nil
}
