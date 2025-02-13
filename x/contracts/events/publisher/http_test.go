package publisher_test

import (
	"context"
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

func TestHTTPPublisher(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)

	evmClient := &testEVMClient{
		logs: []types.Log{
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
		},
	}
	progressStore := &testStore{}
	subscriber := &testSubscriber{
		logs: make(chan types.Log),
	}

	publisher := publisher.NewHTTPPublisher(progressStore, logger, evmClient, subscriber)

	ctx, cancel := context.WithCancel(context.Background())

	doneChan := publisher.Start(ctx, common.Address{})
	publisher.AddContracts(common.HexToAddress("0x1"))

	evmClient.SetBlockNumber(1)
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, evmClient.logs[0]); diff != "" {
			t.Errorf("unexpected log (-got +want):\n%s", diff)
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for log")
	}
	publisher.AddContracts(common.HexToAddress("0x2"))

	evmClient.SetBlockNumber(2)
	select {
	case log := <-subscriber.logs:
		if diff := cmp.Diff(log, evmClient.logs[1]); diff != "" {
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

	if progressStore.blockNumber != 2 {
		t.Errorf("expected block number 2, got %d", progressStore.blockNumber)
	}
}

type testSubscriber struct {
	logs chan types.Log
}

func (t *testSubscriber) PublishLogEvent(ctx context.Context, log types.Log) {
	t.logs <- log
}

type testEVMClient struct {
	mu       sync.Mutex
	blockNum uint64
	logs     []types.Log
}

func (t *testEVMClient) SetBlockNumber(blockNum uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.blockNum = blockNum
}

func (t *testEVMClient) BlockNumber(context.Context) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.blockNum, nil
}

func (t *testEVMClient) FilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
) ([]types.Log, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	logs := make([]types.Log, 0, len(t.logs))
	for _, log := range t.logs {
		if log.BlockNumber >= q.FromBlock.Uint64() && log.BlockNumber <= q.ToBlock.Uint64() {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

type testStore struct {
	mu          sync.Mutex
	blockNumber uint64
}

func (t *testStore) LastBlock() (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.blockNumber, nil
}

func (t *testStore) SetLastBlock(blockNumber uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.blockNumber = blockNumber
	return nil
}
