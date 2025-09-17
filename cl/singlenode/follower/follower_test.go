package follower_test

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/primev/mev-commit/cl/singlenode/follower"
	"github.com/primev/mev-commit/cl/types"
	"github.com/primev/mev-commit/x/util"
)

type mockPayloadDB struct {
	GetPayloadsSinceFunc func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetLatestHeightFunc  func(ctx context.Context) (uint64, error)
}

func (m *mockPayloadDB) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
	return m.GetPayloadsSinceFunc(ctx, sinceHeight, limit)
}

func (m *mockPayloadDB) GetLatestHeight(ctx context.Context) (uint64, error) {
	return m.GetLatestHeightFunc(ctx)
}

type mockBlockBuilder struct {
	executionHead               *types.ExecutionHead
	SetExecutionHeadFromRPCFunc func(ctx context.Context) error
	FinalizeBlockFunc           func(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error
}

func (m *mockBlockBuilder) GetExecutionHead() *types.ExecutionHead {
	return m.executionHead
}

func (m *mockBlockBuilder) SetExecutionHeadFromRPC(ctx context.Context) error {
	return m.SetExecutionHeadFromRPCFunc(ctx)
}

func (m *mockBlockBuilder) FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	return m.FinalizeBlockFunc(ctx, payloadIDStr, executionPayloadStr, msgID)
}

func newMockBlockBuilder() *mockBlockBuilder {
	return &mockBlockBuilder{
		executionHead: nil,
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
		GetLatestHeightFunc: func(ctx context.Context) (uint64, error) {
			return latest, nil
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

	bb := newMockBlockBuilder()
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, bb, ":8080")
	if err != nil {
		t.Fatal(err)
	}

	bb.SetExecutionHeadFromRPCFunc = func(ctx context.Context) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: lastProcessed}
		return nil
	}

	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := follower.SyncFromSharedDB(ctx)
		if err != nil {
			errCh <- err
		}
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
	case err := <-errCh:
		t.Fatal(err)
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
		GetLatestHeightFunc: func(ctx context.Context) (uint64, error) {
			if attempts < 3 {
				attempts++
				return 0, nil
			}
			return 15, nil
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

	bb := newMockBlockBuilder()
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, bb, ":8080")
	if err != nil {
		t.Fatal(err)
	}

	bb.SetExecutionHeadFromRPCFunc = func(ctx context.Context) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: 0} // Only genesis block is available
		return nil
	}

	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := follower.SyncFromSharedDB(ctx)
		if err != nil {
			errCh <- err
		}
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
	case err := <-errCh:
		t.Fatal(err)
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
		GetLatestHeightFunc: func(ctx context.Context) (uint64, error) {
			numGetLatestHeightCalls++
			if numGetLatestHeightCalls > 3 {
				return 253, nil // Simulate that DB has only been updated up to block 253
			}
			return latest + uint64(numGetLatestHeightCalls), nil
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

	bb := newMockBlockBuilder()
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, bb, ":8080")
	if err != nil {
		t.Fatal(err)
	}

	bb.SetExecutionHeadFromRPCFunc = func(ctx context.Context) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: lastProcessed}
		return nil
	}

	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := follower.SyncFromSharedDB(ctx)
		if err != nil {
			errCh <- err
		}
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
	case err := <-errCh:
		t.Fatal(err)
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
		GetLatestHeightFunc: func(ctx context.Context) (uint64, error) {
			getLatestCalls++
			if getLatestCalls <= 3 {
				return 0, nil
			}
			return 1, nil
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

	bb := newMockBlockBuilder()
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, bb, ":8080")
	if err != nil {
		t.Fatal(err)
	}

	bb.SetExecutionHeadFromRPCFunc = func(ctx context.Context) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: 0} // Only genesis block is available
		return nil
	}

	bb.FinalizeBlockFunc = func(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: 1}
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := follower.Start(ctx)

	deadline := time.Now().Add(5 * time.Second)
	for {
		lp := follower.GetExecutionHead()
		if lp == nil {
			continue
		}
		if lp.BlockHeight >= 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for first block to be processed")
		}
		time.Sleep(10 * time.Millisecond)
	}

	finalExecutionHead := follower.GetExecutionHead()
	if finalExecutionHead == nil {
		t.Fatal("execution head is nil")
	}
	if finalExecutionHead.BlockHeight != 1 {
		t.Fatalf("expected execution head block height to be 1, got %d", finalExecutionHead.BlockHeight)
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
		GetLatestHeightFunc: func(ctx context.Context) (uint64, error) {
			return latest, nil
		},
		GetPayloadsSinceFunc: func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
			toReturn := make([]types.PayloadInfo, 0, limit)
			for i := uint64(0); i < uint64(limit); i++ {
				toReturn = append(toReturn, types.PayloadInfo{
					BlockHeight: sinceHeight + i,
					// Encode just the block height
					ExecutionPayload: fmt.Sprintf("%d", sinceHeight+i),
				})
			}
			return toReturn, nil
		},
	}

	syncBatchSize := uint64(20)

	bb := newMockBlockBuilder()
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, bb, ":8080")
	if err != nil {
		t.Fatal(err)
	}

	bb.SetExecutionHeadFromRPCFunc = func(ctx context.Context) error {
		bb.executionHead = &types.ExecutionHead{BlockHeight: lastProcessed}
		return nil
	}

	bb.FinalizeBlockFunc = func(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
		// decode block num from executionPayloadStr
		blockNum, err := strconv.ParseUint(executionPayloadStr, 10, 64)
		if err != nil {
			return err
		}
		bb.executionHead = &types.ExecutionHead{BlockHeight: blockNum}
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := follower.Start(ctx)

	deadline := time.Now().Add(5 * time.Second)
	for {
		lp := follower.GetExecutionHead()
		if lp == nil {
			continue
		}
		if lp.BlockHeight >= 700 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for sync + steady-state; last processed: %d", lp)
		}
		time.Sleep(10 * time.Millisecond)
	}

	finalExecutionHead := follower.GetExecutionHead()
	if finalExecutionHead == nil {
		t.Fatal("execution head is nil")
	}
	if finalExecutionHead.BlockHeight != 700 {
		t.Fatalf("expected execution head block height to be 700, got %d", finalExecutionHead.BlockHeight)
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for follower to stop")
	}
}
