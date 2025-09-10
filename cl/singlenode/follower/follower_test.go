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
	GetPayloadsSinceFunc   func(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetLatestHeightFunc    func(ctx context.Context) (*uint64, error)
	GetPayloadByHeightFunc func(ctx context.Context, height uint64) (*types.PayloadInfo, error)
}

func (m *mockPayloadDB) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
	return m.GetPayloadsSinceFunc(ctx, sinceHeight, limit)
}

func (m *mockPayloadDB) GetLatestHeight(ctx context.Context) (*uint64, error) {
	return m.GetLatestHeightFunc(ctx)
}

func (m *mockPayloadDB) GetPayloadByHeight(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
	return m.GetPayloadByHeightFunc(ctx, height)
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
	caughtUpThreshold := uint64(5)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, caughtUpThreshold, st)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetLastProcessed(lastProcessed)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- follower.SyncFromSharedDB(context.Background())
	}()

	payloadCh := follower.PayloadCh()

	// expect 50 payloads
	received := 0
	expectedBlockHeight := uint64(501)
	numErrSignals := 0
	for received < 50 {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("follower failed, exiting: %v", err)
			}
			if numErrSignals > 1 {
				t.Fatalf("SyncFromSharedDB should only signal nil error once")
			}
			numErrSignals++
			continue
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
		if err != nil {
			t.Fatalf("follower failed, exiting: %v", err)
		}
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
	caughtUpThreshold := uint64(5)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, caughtUpThreshold, st)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- follower.SyncFromSharedDB(context.Background())
	}()

	payloadCh := follower.PayloadCh()

	// expect 15 payloads
	received := 0
	expectedBlockHeight := uint64(1)
	numErrSignals := 0
	for received < 15 {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("follower failed, exiting: %v", err)
			}
			if numErrSignals > 1 {
				t.Fatalf("SyncFromSharedDB should only signal nil error once")
			}
			numErrSignals++
			continue
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
		if err != nil {
			t.Fatalf("follower failed, exiting: %v", err)
		}
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
	caughtUpThreshold := uint64(10)
	st := follower.NewStore(logger, inmemstorage.New())

	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, caughtUpThreshold, st)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetLastProcessed(lastProcessed)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- follower.SyncFromSharedDB(context.Background())
	}()

	payloadCh := follower.PayloadCh()

	// expect payloads up to 253
	received := 0
	expectedBlockHeight := uint64(201)
	numErrSignals := 0
	for received < 53 {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("follower failed, exiting: %v", err)
			}
			if numErrSignals > 1 {
				t.Fatalf("SyncFromSharedDB should only signal nil error once")
			}
			numErrSignals++
			continue
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
		if err != nil {
			t.Fatalf("follower failed, exiting: %v", err)
		}
	case <-payloadCh:
		t.Fatal("received unexpected payload")
	case <-time.After(1 * time.Second):
	}
}

func TestFollower_queryPayloadsFromSharedDB(t *testing.T) {
	t.Parallel()

	lastSignalledBlock := uint64(100)
	numCalls := 0
	payloadRepo := &mockPayloadDB{
		GetPayloadByHeightFunc: func(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
			numCalls++
			if numCalls > 50 { // Simulate catching up to latest block after 50 calls
				if numCalls == 60 {
					// 60th call returns a payload again
					return &types.PayloadInfo{BlockHeight: lastSignalledBlock + 51}, nil
				}
				return nil, sql.ErrNoRows
			}
			return &types.PayloadInfo{BlockHeight: lastSignalledBlock + uint64(numCalls)}, nil
		},
	}

	syncBatchSize := uint64(20)
	caughtUpThreshold := uint64(10)
	logger := util.NewTestLogger(io.Discard)
	st := follower.NewStore(logger, inmemstorage.New())
	follower, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, caughtUpThreshold, st)
	if err != nil {
		t.Fatal(err)
	}

	follower.SetLastSignalledBlock(lastSignalledBlock)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		follower.QueryPayloadsFromSharedDB(ctx)
	}()

	// All 51 payloads should be received in little time
	received := 0
	payloadCh := follower.PayloadCh()
	deadline := time.Now().Add(time.Second)
	for received < 51 {
		select {
		case p := <-payloadCh:
			received++
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received zero payload at %d", p.BlockHeight)
			}
			if p.BlockHeight != lastSignalledBlock+uint64(received) {
				t.Fatalf("expected payload height %d, got %d", lastSignalledBlock+uint64(received), p.BlockHeight)
			}
		case <-time.After(time.Until(deadline)):
			t.Fatalf("timeout waiting for payload")
		case <-ctx.Done():
		}
	}
}

func TestFollower_Start_SimulateNewChain(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)

	getLatestCalls := 0
	getByHeightCalls := 0

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
			t.Fatalf("GetPayloadsSince should not be called, got sinceHeight=%d limit=%d", sinceHeight, limit)
			return nil, nil
		},
		GetPayloadByHeightFunc: func(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
			getByHeightCalls++
			if getByHeightCalls > 10 {
				return nil, sql.ErrNoRows
			}
			return &types.PayloadInfo{BlockHeight: uint64(getByHeightCalls)}, nil
		},
	}

	syncBatchSize := uint64(100)
	caughtUpThreshold := uint64(5)
	st := follower.NewStore(logger, inmemstorage.New())

	f, err := follower.NewFollower(logger, payloadRepo, syncBatchSize, caughtUpThreshold, st)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := f.Start(ctx)

	deadline := time.Now().Add(5 * time.Second)
	for {
		lp, err := st.GetLastProcessed()
		if err != nil {
			t.Fatal(err)
		}
		if lp >= 10 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for first block to be processed")
		}
		time.Sleep(10 * time.Millisecond)
	}

	if f.LastSignalledBlock() != 10 {
		t.Fatalf("expected last signalled block to be 10, got %d", f.LastSignalledBlock())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for follower to stop")
	}
}
