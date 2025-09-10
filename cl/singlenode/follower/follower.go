package follower

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/primev/mev-commit/cl/types"
	"golang.org/x/sync/errgroup"
)

type Follower struct {
	logger            *slog.Logger
	sharedDB          payloadDB
	syncBatchSize     uint64
	caughtUpThreshold uint64
	store             *Store
	payloadCh         chan *types.PayloadInfo
}

const (
	payloadBufferSize = 100
	defaultBackoff    = 200 * time.Millisecond
)

type payloadDB interface {
	GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetPayloadByHeight(ctx context.Context, height uint64) (*types.PayloadInfo, error)
	GetLatestHeight(ctx context.Context) (*uint64, error)
}

func NewFollower(
	logger *slog.Logger,
	sharedDB payloadDB,
	syncBatchSize uint64,
	caughtUpThreshold uint64,
	store *Store,
) (*Follower, error) {
	if sharedDB == nil {
		return nil, errors.New("payload repository not provided")
	}
	if syncBatchSize == 0 {
		return nil, errors.New("sync batch size must be greater than 0")
	}
	if caughtUpThreshold == 0 {
		return nil, errors.New("caught up threshold must be greater than 0")
	}
	if caughtUpThreshold > syncBatchSize {
		return nil, errors.New("caught up threshold must be less than sync batch size")
	}
	return &Follower{
		logger:            logger,
		sharedDB:          sharedDB,
		syncBatchSize:     syncBatchSize,
		caughtUpThreshold: caughtUpThreshold,
		store:             store,
		payloadCh:         make(chan *types.PayloadInfo, payloadBufferSize),
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
		f.logger.Info("Starting initial sync from shared DB")
		if err := f.syncFromSharedDB(egCtx); err != nil {
			f.logger.Error("Failed during initial sync", "error", err)
			return err
		}
		f.logger.Info("Entering steady-state querying of shared DB")
		f.queryPayloadsFromSharedDB(egCtx)
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

func (f *Follower) syncFromSharedDB(ctx context.Context) error {
	lastProcessedBlock, err := f.store.GetLastProcessed()
	if err != nil {
		return err
	}
	lastSignalledBlock := lastProcessedBlock

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		targetBlock, err := f.getLatestHeightWithBackoff(ctx)
		if err != nil {
			return err
		}

		blocksRemaining := targetBlock - lastSignalledBlock

		if blocksRemaining <= f.caughtUpThreshold {
			f.logger.Info("Sync complete", "last_processed", lastProcessedBlock, "target", targetBlock)
			return nil
		}

		limit := min(f.syncBatchSize, blocksRemaining)

		innerCtx, innerCancel := context.WithTimeout(ctx, 10*time.Second)
		payloads, err := f.sharedDB.GetPayloadsSince(innerCtx, lastProcessedBlock+1, int(limit)) // TODO: confirm no off-by-one issue
		innerCancel()
		if err != nil {
			return err
		}
		if len(payloads) == 0 {
			return errors.New("no payloads returned from valid query")
		}

		for _, p := range payloads {
			f.payloadCh <- &p // Non-blocking up to payloadBufferSize
			lastSignalledBlock = p.BlockHeight
		}
	}
}

func (f *Follower) getLatestHeightWithBackoff(ctx context.Context) (uint64, error) {
	const maxRetries = 10
	const base = 5 * time.Second
	for attempt := range maxRetries {
		lctx, cancel := context.WithTimeout(ctx, base)
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

		backoff := base * time.Duration(attempt+1)
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(backoff):
		}
	}
	return 0, errors.New("failed to get latest payload after retries")
}

func (f *Follower) queryPayloadsFromSharedDB(ctx context.Context) {
	for {
		lastProcessed, err := f.store.GetLastProcessed()
		if err != nil {
			f.logger.Error("Failed to read last processed height", "error", err)
			time.Sleep(defaultBackoff)
			continue
		}

		lctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		payload, err := f.sharedDB.GetPayloadByHeight(lctx, lastProcessed+1)
		cancel()

		if err != nil {
			if err == sql.ErrNoRows {
				time.Sleep(time.Millisecond) // New payload will likely be available within milliseconds
				continue
			}
			f.logger.Error("Failed to fetch next payload by height with unexpected error", "height", lastProcessed+1, "error", err)
			time.Sleep(defaultBackoff)
			continue
		}

		if payload == nil {
			f.logger.Error("Received nil payload from valid query")
			time.Sleep(defaultBackoff)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case f.payloadCh <- payload:
			f.logger.Debug("Sent payload to channel", "height", lastProcessed+1)
		default:
			f.logger.Error("Payload channel buffer is full", "height", lastProcessed+1)
			time.Sleep(defaultBackoff)
			continue
		}
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
			if err := f.store.SetLastProcessed(p.BlockHeight); err != nil {
				f.logger.Error("Failed to persist last processed height", "height", p.BlockHeight, "error", err)
				continue
			}
		}
	}
}

// TODO: confirm nothing would be broken if the service crashes mid processing of a payload.
// How does the node recover from this w.r.t engine api? Might need to save FSM state in additon to block number in kv store.

// TODO: Or w.r.t above could the engine api just handle 1-2 duplicate calls and gracefully continue? Need to confirm.

func (f *Follower) handlePayload(ctx context.Context, payload *types.PayloadInfo) error {
	// TODO: Apply the payload to follower's EL via Engine API in later steps.
	f.logger.Info("Processing payload",
		"payload_id", payload.PayloadID,
		"block_height", payload.BlockHeight,
	)
	return nil
}
