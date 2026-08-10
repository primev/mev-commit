package sender

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ProcessQueuedForTest runs one queue processing pass with the given
// context. Tests use it to drive the dispatch path directly, without
// the ticker loop that Start runs.
func (t *TxSender) ProcessQueuedForTest(ctx context.Context) {
	if t.eg == nil {
		t.eg, t.egCtx = errgroup.WithContext(context.Background())
	}
	t.processMu.Lock()
	defer t.processMu.Unlock()
	t.processQueuedTransactions(ctx)
}

// WaitWorkersForTest waits for every worker goroutine that a
// ProcessQueuedForTest call started.
func (t *TxSender) WaitWorkersForTest() error {
	return t.eg.Wait()
}

// FillWorkerPoolForTest fills the worker pool until only free slots
// remain. Tests use it to make a batch slot acquisition block partway
// through. It returns a drain function that releases the filled slots;
// the drain function is safe to call more than once. Call it only
// while no worker goroutine runs.
func (t *TxSender) FillWorkerPoolForTest(free int) func() {
	filled := 0
	for cap(t.workerPool)-len(t.workerPool) > free {
		t.workerPool <- struct{}{}
		filled++
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			for ; filled > 0; filled-- {
				<-t.workerPool
			}
		})
	}
}
