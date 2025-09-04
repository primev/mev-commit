package publisher_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/util"
)

func TestWSPublisher(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)

	logs := []types.Log{
		{
			BlockNumber: 1,
			Address:     common.HexToAddress("0x1"),
			Topics:      []common.Hash{common.HexToHash("0x1")},
			Data:        []byte("abcd"),
		},
		{
			BlockNumber: 2,
			Address:     common.HexToAddress("0x2"),
			Topics:      []common.Hash{common.HexToHash("0x2")},
			Data:        []byte("efgh"),
		},
	}

	// First subscription should error, second should run
	errC := make(chan error, 1)
	errC <- errors.New("test error")

	evmClient := &testWSEVMClient{
		subscribed: make(chan struct{}, 3),
		errC:       errC,
	}
	progressStore := &testStore{}
	subscriber := &testSubscriber{
		logs: make(chan types.Log),
	}

	p := publisher.NewWSPublisher(progressStore, logger, evmClient, subscriber)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneChan := p.Start(ctx)
	p.AddContracts(common.Address{})

	// Wait for first subscribe (will immediately error and cause resubscribe)
	select {
	case <-evmClient.subscribed:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for first subscribe")
	}
	// Wait for second subscribe (active)
	select {
	case <-evmClient.subscribed:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for second subscribe")
	}

	// Send two logs and expect them to be forwarded
	evmClient.SendLog(logs[0])
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, logs[0]); diff != "" {
			t.Errorf("unexpected log (-got +want):\n%s", diff)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for first log")
	}

	evmClient.SendLog(logs[1])
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, logs[1]); diff != "" {
			t.Errorf("unexpected log (-got +want):\n%s", diff)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for second log")
	}

	cancel()
	select {
	case <-doneChan:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for doneChan")
	}

	// Ensure current subscription was unsubscribed
	evmClient.mu.Lock()
	sub := evmClient.sub
	evmClient.mu.Unlock()
	select {
	case <-sub.done:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for subscription to be unsubscribed")
	}

	if bn, _ := progressStore.LastBlock(); bn != 2 {
		t.Errorf("expected block number 2, got %d", bn)
	}
}

type testSubscription struct {
	done chan struct{}
	errC chan error
}

func (s *testSubscription) Unsubscribe() {
	close(s.done)
}

func (s *testSubscription) Err() <-chan error {
	return s.errC
}

type testWSEVMClient struct {
	mu         sync.Mutex
	subscribed chan struct{}
	logs       chan<- types.Log
	sub        *testSubscription
	errC       chan error
}

func (c *testWSEVMClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, logs chan<- types.Log) (ethereum.Subscription, error) {
	c.mu.Lock()
	c.logs = logs
	var errCh chan error
	if c.errC != nil {
		errCh = c.errC
		c.errC = nil
	} else {
		errCh = make(chan error)
	}
	c.sub = &testSubscription{
		done: make(chan struct{}),
		errC: errCh,
	}
	c.mu.Unlock()

	c.subscribed <- struct{}{}
	return c.sub, nil
}

func (c *testWSEVMClient) SendLog(l types.Log) {
	c.mu.Lock()
	ch := c.logs
	c.mu.Unlock()
	ch <- l
}
