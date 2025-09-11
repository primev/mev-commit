package follower_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/primev/mev-commit/cl/singlenode/follower"
	"github.com/primev/mev-commit/cl/types"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	"github.com/primev/mev-commit/x/util"
)

type mockPayloadDB struct {
	GetPayloadsSinceFunc func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetLatestHeightFunc  func(ctx context.Context) (*uint64, error)
}

func (m *mockPayloadDB) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
	return m.GetPayloadsSinceFunc(ctx, sinceHeight, limit)
}

func (m *mockPayloadDB) GetLatestHeight(ctx context.Context) (*uint64, error) {
	return m.GetLatestHeightFunc(ctx)
}

type mockBlockBuilder struct {
	GetExecutionHeadFunc        func() *types.ExecutionHead
	SetExecutionHeadFromRPCFunc func(ctx context.Context) error
	FinalizeBlockFunc           func(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error
}

func (m *mockBlockBuilder) GetExecutionHead() *types.ExecutionHead {
	return m.GetExecutionHeadFunc()
}

func (m *mockBlockBuilder) SetExecutionHeadFromRPC(ctx context.Context) error {
	return m.SetExecutionHeadFromRPCFunc(ctx)
}

func (m *mockBlockBuilder) FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	return m.FinalizeBlockFunc(ctx, payloadIDStr, executionPayloadStr, msgID)
}

func newNoopBlockBuilder() *mockBlockBuilder {
	return &mockBlockBuilder{
		GetExecutionHeadFunc: func() *types.ExecutionHead {
			return nil
		},
		SetExecutionHeadFromRPCFunc: func(ctx context.Context) error {
			return nil
		},
		FinalizeBlockFunc: func(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
			return nil
		},
	}
}

func TestFollower_syncFromSharedDB(t *testing.T) {
	t.Parallel()

	lastProcessed := uint64(500)
	latest := uint64(550)

	logger := util.NewTestLogger(io.Discard)
	payloadRepo := &mockPayloadDB{
		GetLatestHeightFunc: func(ctx context.Context) (*uint64, error) {
			return &latest, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			if sinceHeight != 501 {
				t.Fatal("unexpected sinceHeight", sinceHeight)
			}
			toReturn := []types.PayloadInfo{}
			for i := sinceHeight; i <= latest; i++ {
				toReturn = append(toReturn, types.PayloadInfo{BlockHeight: i})
			}
			return toReturn, nil
		},
	}
	syncBatchSize := uint64(100)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, st, newNoopBlockBuilder())
	if err != nil {
		t.Fatal(err)
	}

	err = follower.SetLastProcessed(context.Background(), lastProcessed)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		follower.SyncFromSharedDB(ctx)
	}()

	payloadCh := follower.PayloadCh()

	// expect 50 payloads
	received := 0
	expectedBlockHeight := uint64(501)
	for received < 50 {
		select {
		case p := <-payloadCh:
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received zero payload for expected block height %d", expectedBlockHeight)
			}
			if p.BlockHeight != expectedBlockHeight {
				t.Fatalf("expected payload height %d, got %d", expectedBlockHeight, p.BlockHeight)
			}
			expectedBlockHeight++
			received++
		case <-time.After(1 * time.Second):
			t.Fatalf("timeout waiting for payload for expected block height %d", expectedBlockHeight)
		}
	}
	if received != 50 {
		t.Fatalf("expected 50 payloads, got %d", received)
	}

	// No more than 50
	select {
	case <-payloadCh:
		t.Fatal("received unexpected payload")
	case <-time.After(1 * time.Second):
	}
}

func TestFollower_syncFromSharedDB_NoRows(t *testing.T) {
	t.Parallel()

	attempts := 0
	logger := util.NewTestLogger(io.Discard)
	payloadRepo := &mockPayloadDB{
		GetLatestHeightFunc: func(ctx context.Context) (*uint64, error) {
			if attempts < 3 {
				attempts++
				return nil, sql.ErrNoRows
			}
			block15 := uint64(15)
			return &block15, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			if sinceHeight != 1 {
				return nil, fmt.Errorf("unexpected sinceHeight %d", sinceHeight)
			}
			if limit != 15 {
				return nil, fmt.Errorf("unexpected limit %d", limit)
			}
			toReturn := []types.PayloadInfo{}
			for i := 1; i <= 15; i++ {
				toReturn = append(toReturn, types.PayloadInfo{BlockHeight: uint64(i)})
			}
			return toReturn, nil
		},
	}
	syncBatchSize := uint64(100)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, st, newNoopBlockBuilder())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		follower.SyncFromSharedDB(ctx)
	}()

	payloadCh := follower.PayloadCh()

	// expect 15 payloads
	received := 0
	expectedBlockHeight := uint64(1)
	for received < 15 {
		select {
		case p := <-payloadCh:
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received zero payload at %d", expectedBlockHeight)
			}
			if p.BlockHeight != expectedBlockHeight {
				t.Fatalf("expected payload height %d, got %d", expectedBlockHeight, p.BlockHeight)
			}
			expectedBlockHeight++
			received++
		case <-time.After(10 * time.Second):
			t.Fatalf("timeout waiting for payload %d", expectedBlockHeight)
		}
	}
	if received != 15 {
		t.Fatalf("expected 15 payloads, got %d", received)
	}

	// No more than 15
	select {
	case <-payloadCh:
		t.Fatal("received unexpected payload")
	case <-time.After(1 * time.Second):
	}
}

func TestFollower_syncFromSharedDB_MultipleIterations(t *testing.T) {
	t.Parallel()

	lastProcessed := uint64(200)
	latest := uint64(250)

	logger := util.NewTestLogger(io.Discard)

	numGetLatestHeightCalls := 0
	numGetPayloadsCalls := 0
	payloadRepo := &mockPayloadDB{
		GetLatestHeightFunc: func(ctx context.Context) (*uint64, error) {
			numGetLatestHeightCalls++
			if numGetLatestHeightCalls > 3 {
				toReturn := uint64(253) // Simulate that DB has only been updated up to block 253
				return &toReturn, nil
			}
			toReturn := latest + uint64(numGetLatestHeightCalls)
			return &toReturn, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			numGetPayloadsCalls++
			switch numGetPayloadsCalls {
			case 1:
				// First iteration should request payloads from 201 to 220
				if sinceHeight != 201 {
					t.Fatal("unexpected sinceHeight", sinceHeight)
				}
				if limit != 20 {
					t.Fatal("unexpected limit", limit)
				}
			case 2:
				// Second iteration should request payloads from 221 to 240
				if sinceHeight != 221 {
					t.Fatal("unexpected sinceHeight", sinceHeight)
				}
				if limit != 20 {
					t.Fatal("unexpected limit", limit)
				}
			case 3:
				// Third iteration should request payloads from 241 to 253
				if sinceHeight != 241 {
					t.Fatal("unexpected sinceHeight", sinceHeight)
				}
				if limit != 13 {
					t.Fatal("unexpected limit", limit)
				}
			default:
				t.Fatal("unexpected numGetPayloadsCalls", numGetPayloadsCalls)
				return nil, nil
			}
			toReturn := []types.PayloadInfo{}
			for i := sinceHeight; i < sinceHeight+uint64(limit); i++ {
				toReturn = append(toReturn, types.PayloadInfo{BlockHeight: i})
			}
			return toReturn, nil
		},
	}
	syncBatchSize := uint64(20)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, st, newNoopBlockBuilder())
	if err != nil {
		t.Fatal(err)
	}

	err = follower.SetLastProcessed(context.Background(), lastProcessed)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		follower.SyncFromSharedDB(ctx)
	}()

	payloadCh := follower.PayloadCh()

	// expect payloads up to 253
	received := 0
	expectedBlockHeight := uint64(201)
	for received < 53 {
		select {
		case p := <-payloadCh:
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received zero payload at %d", expectedBlockHeight)
			}
			if p.BlockHeight != expectedBlockHeight {
				t.Fatalf("expected payload height %d, got %d", expectedBlockHeight, p.BlockHeight)
			}
			expectedBlockHeight++
			received++
		case <-time.After(10 * time.Second):
			t.Fatalf("timeout waiting for payload %d", expectedBlockHeight)
		}
	}
	if received != 53 {
		t.Fatalf("expected 53 payloads, got %d", received)
	}

	// No more than 53
	select {
	case <-payloadCh:
		t.Fatal("received unexpected payload")
	case <-time.After(1 * time.Second):
	}
}

func TestFollower_Start_SimulateNewChain(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)

	getLatestCalls := 0

	payloadRepo := &mockPayloadDB{
		GetLatestHeightFunc: func(ctx context.Context) (*uint64, error) {
			getLatestCalls++
			if getLatestCalls <= 3 {
				return nil, sql.ErrNoRows
			}
			h := uint64(1)
			return &h, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			if sinceHeight != 1 {
				t.Fatalf("unexpected sinceHeight %d", sinceHeight)
			}
			if limit != 1 {
				t.Fatalf("unexpected limit %d", limit)
			}
			return []types.PayloadInfo{{BlockHeight: 1}}, nil
		},
	}

	syncBatchSize := uint64(100)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, st, newNoopBlockBuilder())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := follower.Start(ctx)

	deadline := time.Now().Add(5 * time.Second)
	for {
		lp, err := follower.GetLastProcessed(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if lp >= 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for first block to be processed")
		}
		time.Sleep(10 * time.Millisecond)
	}

	if follower.LastSignalledBlock() != 1 {
		t.Fatalf("expected last signalled block to be 1, got %d", follower.LastSignalledBlock())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for follower to stop")
	}
}

func TestFollower_Start_SyncExistingChain(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)

	lastProcessed := uint64(450)
	latest := uint64(700)

	payloadRepo := &mockPayloadDB{
		GetLatestHeightFunc: func(ctx context.Context) (*uint64, error) {
			return &latest, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			toReturn := make([]types.PayloadInfo, 0, limit)
			for i := uint64(0); i < uint64(limit); i++ {
				toReturn = append(toReturn, types.PayloadInfo{BlockHeight: sinceHeight + i})
			}
			return toReturn, nil
		},
	}

	syncBatchSize := uint64(20)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, st, newNoopBlockBuilder())
	if err != nil {
		t.Fatal(err)
	}

	if err := follower.SetLastProcessed(context.Background(), lastProcessed); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := follower.Start(ctx)

	deadline := time.Now().Add(5 * time.Second)
	for {
		lp, err := follower.GetLastProcessed(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if lp >= 700 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for sync + steady-state; last processed: %d", lp)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if follower.LastSignalledBlock() != 700 {
		t.Fatalf("expected last signalled to be %d, got %d", 700, follower.LastSignalledBlock())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for follower to stop")
	}
}
