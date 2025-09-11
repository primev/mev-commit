package follower

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/primev/mev-commit/cl/types"
	"golang.org/x/sync/errgroup"
)

type Follower struct {
	logger             *slog.Logger
	sharedDB           payloadDB
	syncBatchSize      uint64
	store              *Store
	storeMu            sync.RWMutex
	payloadCh          chan types.PayloadInfo
	lastSignalledBlock atomic.Uint64 // Last block num signalled through payloadCh
	bb                 blockBuilder
}

const (
	payloadBufferSize = 100
	defaultBackoff    = 200 * time.Millisecond
)

type payloadDB interface {
	GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetLatestHeight(ctx context.Context) (*uint64, error)
}

type blockBuilder interface {
	GetExecutionHead() *types.ExecutionHead
	FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error
	SetExecutionHeadFromRPC(ctx context.Context) error
}

func NewFollower(
	logger *slog.Logger,
	sharedDB payloadDB,
	syncBatchSize uint64,
	store *Store,
	bb blockBuilder,
) (*Follower, error) {
	if sharedDB == nil {
		return nil, errors.New("payload repository not provided")
	}
	if syncBatchSize == 0 {
		return nil, errors.New("sync batch size must be greater than 0")
	}
	return &Follower{
		logger:        logger,
		sharedDB:      sharedDB,
		syncBatchSize: syncBatchSize,
		store:         store,
		payloadCh:     make(chan types.PayloadInfo, payloadBufferSize),
		bb:            bb,
	}, nil
}

func (f *Follower) Start(ctx context.Context) <-chan struct{} {

	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		f.handlePayloads(egCtx)
		return nil
	})

	eg.Go(func() error {
		f.logger.Info("Starting sync from shared DB")
		f.syncFromSharedDB(egCtx)
		return nil
	})

	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			f.logger.Error("follower failed, exiting", "error", err)
		}
	}()

	return done
}

func (f *Follower) syncFromSharedDB(ctx context.Context) {
	var err error
	// lastSignalledBlock is only set from disk here
	lastProcessed, err := f.getLastProcessed(ctx)
	if err != nil {
		f.logger.Error("failed to get last processed block", "error", err)
		return
	}
	f.lastSignalledBlock.Store(lastProcessed)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		targetBlock, err := f.getLatestHeightWithBackoff(ctx)
		if err != nil {
			f.sleepRespectingContext(ctx, defaultBackoff)
			continue
		}

		blocksRemaining := targetBlock - f.lastSignalledBlock.Load()

		if blocksRemaining == 0 {
			f.sleepRespectingContext(ctx, time.Millisecond) // New payload will likely be available within milliseconds
			continue
		}

		limit := min(f.syncBatchSize, blocksRemaining)

		innerCtx, innerCancel := context.WithTimeout(ctx, 10*time.Second)
		payloads, err := f.sharedDB.GetPayloadsSince(innerCtx, f.lastSignalledBlock.Load()+1, int(limit))
		innerCancel()
		if err != nil {
			f.logger.Error("failed to get payloads since", "error", err)
			f.sleepRespectingContext(ctx, defaultBackoff)
			continue
		}
		if len(payloads) == 0 {
			f.logger.Error("no payloads returned from valid query")
			f.sleepRespectingContext(ctx, defaultBackoff)
			continue
		}

		for i := range payloads {
			p := payloads[i]
			select {
			case <-ctx.Done():
				return
			case f.payloadCh <- p:
			}
			f.lastSignalledBlock.Store(p.BlockHeight)
		}
	}
}

func (f *Follower) getLatestHeightWithBackoff(ctx context.Context) (uint64, error) {
	const maxRetries = 10
	for attempt := range maxRetries {
		lctx, cancel := context.WithTimeout(ctx, time.Second)
		latest, err := f.sharedDB.GetLatestHeight(lctx)
		cancel()
		if err == nil {
			if latest != nil {
				return *latest, nil
			}
			return 0, errors.New("nil height returned")
		}
		if err != sql.ErrNoRows {
			return 0, err
		}

		backoff := defaultBackoff * time.Duration(attempt+1)
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(backoff):
		}
	}
	return 0, errors.New("failed to get latest payload after retries")
}

func (f *Follower) sleepRespectingContext(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(duration):
	}
}

func (f *Follower) handlePayloads(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-f.payloadCh:
			if err := f.handlePayload(ctx, p); err != nil {
				f.logger.Error("Failed to process payload", "height", p.BlockHeight, "error", err)
				continue
			}
			if err := f.setLastProcessed(ctx, p.BlockHeight); err != nil {
				f.logger.Error("Failed to persist last processed height", "height", p.BlockHeight, "error", err)
				continue
			}
		}
	}
}

func (f *Follower) handlePayload(ctx context.Context, payload types.PayloadInfo) error {
	if f.bb.GetExecutionHead() == nil {
		if err := f.bb.SetExecutionHeadFromRPC(ctx); err != nil {
			return fmt.Errorf("failed to set execution head from rpc: %w", err)
		}
	}
	return f.bb.FinalizeBlock(ctx, payload.PayloadID, payload.ExecutionPayload, "")
}

func (f *Follower) getLastProcessed(ctx context.Context) (uint64, error) {
	f.storeMu.RLock()
	defer f.storeMu.RUnlock()
	return f.store.GetLastProcessed()
}

func (f *Follower) setLastProcessed(ctx context.Context, height uint64) error {
	f.storeMu.Lock()
	defer f.storeMu.Unlock()
	return f.store.SetLastProcessed(height)
}
