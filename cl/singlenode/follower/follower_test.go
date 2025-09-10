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
	fmt.Println(follower)

	st.SetLastProcessed(lastProcessed)

	errCh := make(chan error, 1)

	go func() {
		errCh <- follower.SyncFromSharedDB(context.Background())
	}()

	payloadCh := follower.PayloadCh()

	// expect all 50 payload signals to be sent
	for i := 501; i <= 550; i++ {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("follower failed, exiting: %v", err)
			}
		case p := <-payloadCh:
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received nil payload at %d", i)
			}
			if p.BlockHeight != uint64(i) {
				t.Fatalf("expected payload height %d, got %d", i, p.BlockHeight)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("timeout waiting for payload %d", i)
		}
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
		GetPayloadByHeightFunc: func(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
			return nil, nil
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
	for i := 1; i <= 15; i++ {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("follower failed, exiting: %v", err)
			}
		case p := <-payloadCh:
			if p == (types.PayloadInfo{}) {
				t.Fatalf("received nil payload at %d", i)
			}
			if p.BlockHeight != uint64(i) {
				t.Fatalf("expected payload height %d, got %d", i, p.BlockHeight)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("timeout waiting for payload %d", i)
		}
	}
}

// TODO: test with two for loop iterations to test last processed vs last signalled block

// TODO: test syncFromSharedDB leading into steady state with sql.ErrNoRows returned a couple times.
// When it does come online.. test from block 1, simulating new chain.

// TODO: test where channel becomes full and the sync thread blocks.. testing "// Non-blocking up to payloadBufferSize"

// TODO: test threshold
