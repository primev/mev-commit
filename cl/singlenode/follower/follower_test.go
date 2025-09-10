package follower_test

import (
	"context"
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

	go func() {
		follower.SyncFromSharedDB(context.Background())
	}()

	payloadCh := follower.PayloadCh()

	// expect all 50 payload signals to be sent
	for i := 501; i <= 550; i++ {
		select {
		case p := <-payloadCh:
			if p == nil {
				t.Fatalf("received nil payload at %d", i)
			}
			if p.BlockHeight != uint64(i) {
				t.Fatalf("expected payload height %d, got %d", i, p.BlockHeight)
			}
			if err := st.SetLastProcessed(p.BlockHeight); err != nil {
				t.Fatalf("failed to set last processed: %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("timeout waiting for payload %d", i)
		}
	}

	// No more than 50
	select {
	case <-payloadCh:
		t.Fatal("received unexpected payload")
	case <-time.After(1 * time.Second):
	}
}

// TODO: test where channel becomes full and the sync thread blocks.. testing "// Non-blocking up to payloadBufferSize"

// TODO: test threshold
