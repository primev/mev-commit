package publisher_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/primevprotocol/mev-commit/x/contracts/publisher"
	"github.com/primevprotocol/mev-commit/x/util"
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

	evmClient := &testWSEVMClient{
		subscribed: make(chan struct{}),
		sub: &testSubscription{
			done: make(chan struct{}),
		},
	}
	progressStore := &testStore{}
	subscriber := &testSubscriber{
		logs: make(chan types.Log),
	}

	publisher := publisher.NewWSPublisher(progressStore, logger, evmClient, subscriber)
	noContractsDone := publisher.Start(context.Background())
	select {
	case <-noContractsDone:
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for doneChan")
	}

	ctx, cancel := context.WithCancel(context.Background())

	doneChan := publisher.Start(ctx, common.Address{})

	<-evmClient.subscribed

	evmClient.logs <- logs[0]
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, logs[0]); diff != "" {
			t.Errorf("unexpected log (-got +want):\n%s", diff)
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for log")
	}

	evmClient.logs <- logs[1]
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, logs[1]); diff != "" {
			t.Errorf("unexpected log (-got +want):\n%s", diff)
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for log")
	}

	cancel()
	select {
	case <-doneChan:
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for doneChan")
	}

	select {
	case <-evmClient.sub.done:
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for subscription to be unsubscribed")
	}

	if progressStore.blockNumber != 2 {
		t.Errorf("expected block number 2, got %d", progressStore.blockNumber)
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
	subscribed chan struct{}
	logs       chan<- types.Log
	sub        *testSubscription
}

func (c *testWSEVMClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, logs chan<- types.Log) (ethereum.Subscription, error) {
	defer close(c.subscribed)
	c.logs = logs
	return c.sub, nil
}
