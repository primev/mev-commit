package payloadstore

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/primev/mev-commit/cl/types"
	"github.com/redis/go-redis/v9"
)

func newTestRepo(t *testing.T) (*RedisRepository, *miniredis.Miniredis, func()) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run error: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := NewRedisRepository(rdb, logger)

	cleanup := func() {
		_ = repo.Close()
		mr.Close()
	}
	return repo, mr, cleanup
}

func TestEmptyRepoLatestHeightIsZero(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	h, err := repo.GetLatestHeight(ctx)
	if err != nil {
		t.Fatalf("GetLatestHeight error: %v", err)
	}
	if h != 0 {
		t.Fatalf("expected latest height 0 for empty repo, got %d", h)
	}
}

func TestSaveAndGetLatest(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	now := time.Now().UTC()
	p1 := &types.PayloadInfo{PayloadID: "a", ExecutionPayload: "pa", BlockHeight: 10, InsertedAt: now}
	p2 := &types.PayloadInfo{PayloadID: "b", ExecutionPayload: "pb", BlockHeight: 12, InsertedAt: now}
	p3 := &types.PayloadInfo{PayloadID: "c", ExecutionPayload: "pc", BlockHeight: 15, InsertedAt: now}

	if err := repo.SavePayload(ctx, p1); err != nil {
		t.Fatalf("SavePayload p1 error: %v", err)
	}
	if err := repo.SavePayload(ctx, p2); err != nil {
		t.Fatalf("SavePayload p2 error: %v", err)
	}
	if err := repo.SavePayload(ctx, p3); err != nil {
		t.Fatalf("SavePayload p3 error: %v", err)
	}

	latest, err := repo.GetLatestHeight(ctx)
	if err != nil {
		t.Fatalf("GetLatestHeight error: %v", err)
	}
	if latest != 15 {
		t.Fatalf("expected latest height 15, got %d", latest)
	}
}

func TestGetPayloadsSince(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	payloads := []*types.PayloadInfo{
		{PayloadID: "h10", ExecutionPayload: "p10", BlockHeight: 10, InsertedAt: now},
		{PayloadID: "h12", ExecutionPayload: "p12", BlockHeight: 12, InsertedAt: now},
		{PayloadID: "h15", ExecutionPayload: "p15", BlockHeight: 15, InsertedAt: now},
		{PayloadID: "h20", ExecutionPayload: "p20", BlockHeight: 20, InsertedAt: now},
	}
	for _, p := range payloads {
		if err := repo.SavePayload(ctx, p); err != nil {
			t.Fatalf("SavePayload error: %v", err)
		}
	}

	got, err := repo.GetPayloadsSince(ctx, 12, 100)
	if err != nil {
		t.Fatalf("GetPayloadsSince error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 payloads, got %d", len(got))
	}
	if got[0].BlockHeight != 12 || got[1].BlockHeight != 15 || got[2].BlockHeight != 20 {
		t.Fatalf("unexpected order or heights: %#v", got)
	}

	got, err = repo.GetPayloadsSince(ctx, 10, 2)
	if err != nil {
		t.Fatalf("GetPayloadsSince error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 payloads, got %d", len(got))
	}
	if got[0].BlockHeight != 10 || got[1].BlockHeight != 12 {
		t.Fatalf("unexpected order or heights with limit=2: %#v", got)
	}
	if got[0].PayloadID != "h10" || got[1].PayloadID != "h12" {
		t.Fatalf("unexpected order or payload IDs: %#v", got)
	}
	if got[0].ExecutionPayload != "p10" || got[1].ExecutionPayload != "p12" {
		t.Fatalf("unexpected order or execution payloads: %#v", got)
	}
	if got[0].InsertedAt != now || got[1].InsertedAt != now {
		t.Fatalf("unexpected order or inserted at times: %#v", got)
	}
}

func TestUpsertByHeight(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	orig := &types.PayloadInfo{PayloadID: "orig", ExecutionPayload: "p1", BlockHeight: 12, InsertedAt: now}
	if err := repo.SavePayload(ctx, orig); err != nil {
		t.Fatalf("SavePayload orig error: %v", err)
	}

	updated := &types.PayloadInfo{PayloadID: "updated", ExecutionPayload: "p2", BlockHeight: 12, InsertedAt: now.Add(time.Second)}
	if err := repo.SavePayload(ctx, updated); err != nil {
		t.Fatalf("SavePayload updated error: %v", err)
	}

	got, err := repo.GetPayloadsSince(ctx, 12, 10)
	if err != nil {
		t.Fatalf("GetPayloadsSince error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 payload at height 12, got %d", len(got))
	}
	if got[0].PayloadID != "updated" || got[0].ExecutionPayload != "p2" {
		t.Fatalf("upsert failed, got %#v", got[0])
	}
}

func TestClose(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()
	if err := repo.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}
}

func TestGetPayloadByHeight_Found(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	p5 := &types.PayloadInfo{PayloadID: "id5", ExecutionPayload: "p5", BlockHeight: 5, InsertedAt: now}
	p7 := &types.PayloadInfo{PayloadID: "id7", ExecutionPayload: "p7", BlockHeight: 7, InsertedAt: now}

	if err := repo.SavePayload(ctx, p5); err != nil {
		t.Fatalf("SavePayload p5 error: %v", err)
	}
	if err := repo.SavePayload(ctx, p7); err != nil {
		t.Fatalf("SavePayload p7 error: %v", err)
	}

	got, err := repo.GetPayloadByHeight(ctx, 7)
	if err != nil {
		t.Fatalf("GetPayloadByHeight error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected payload, got nil")
	}
	if got.PayloadID != "id7" || got.ExecutionPayload != "p7" || got.BlockHeight != 7 || got.InsertedAt != now {
		t.Fatalf("unexpected payload: %#v", got)
	}

	got, err = repo.GetPayloadByHeight(ctx, 5)
	if err != nil {
		t.Fatalf("GetPayloadByHeight error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected payload, got nil")
	}
	if got.PayloadID != "id5" || got.ExecutionPayload != "p5" || got.BlockHeight != 5 || got.InsertedAt != now {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestGetPayloadByHeight_NotFound(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	p5 := &types.PayloadInfo{PayloadID: "id5", ExecutionPayload: "p5", BlockHeight: 5, InsertedAt: now}
	if err := repo.SavePayload(ctx, p5); err != nil {
		t.Fatalf("SavePayload p5 error: %v", err)
	}

	got, err := repo.GetPayloadByHeight(ctx, 6)
	if err == nil {
		t.Fatalf("expected error for missing height, got nil (payload: %#v)", got)
	}
	if got != nil {
		t.Fatalf("expected nil payload on not found, got: %#v", got)
	}
}

func TestGetLatestPayload_Empty(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	got, err := repo.GetLatestPayload(ctx)
	if err != nil {
		t.Fatalf("GetLatestPayload error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil payload for empty repo, got %#v", got)
	}
}

func TestGetLatestPayload_ReturnsHighest(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	p1 := &types.PayloadInfo{PayloadID: "a", ExecutionPayload: "pa", BlockHeight: 10, InsertedAt: now}
	p2 := &types.PayloadInfo{PayloadID: "b", ExecutionPayload: "pb", BlockHeight: 12, InsertedAt: now}
	p3 := &types.PayloadInfo{PayloadID: "c", ExecutionPayload: "pc", BlockHeight: 15, InsertedAt: now}

	if err := repo.SavePayload(ctx, p1); err != nil {
		t.Fatalf("SavePayload p1 error: %v", err)
	}
	if err := repo.SavePayload(ctx, p2); err != nil {
		t.Fatalf("SavePayload p2 error: %v", err)
	}
	if err := repo.SavePayload(ctx, p3); err != nil {
		t.Fatalf("SavePayload p3 error: %v", err)
	}

	got, err := repo.GetLatestPayload(ctx)
	if err != nil {
		t.Fatalf("GetLatestPayload error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected payload, got nil")
	}
	if got.BlockHeight != 15 || got.PayloadID != "c" || got.ExecutionPayload != "pc" || got.InsertedAt != now {
		t.Fatalf("unexpected latest payload: %#v", got)
	}
}
