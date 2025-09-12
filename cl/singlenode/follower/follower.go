package follower

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/primev/mev-commit/cl/types"
	"golang.org/x/sync/errgroup"
)

type Follower struct {
	logger        *slog.Logger
	sharedDB      payloadDB
	syncBatchSize uint64
	payloadCh     chan types.PayloadInfo
	bbMutex       sync.RWMutex
	bb            blockBuilder
}

const (
	defaultBackoff = 200 * time.Millisecond
)

type payloadDB interface {
	GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error)
	GetLatestHeight(ctx context.Context) (uint64, error)
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
		payloadCh:     make(chan types.PayloadInfo),
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
	if f.getExecutionHead() == nil {
		if err := f.setExecutionHeadFromRPC(ctx); err != nil {
			f.logger.Error("failed to set execution head from rpc", "error", err)
			return
		}
		f.logger.Debug("set execution head from rpc")
	}

	lastSignalledBlock := f.getExecutionHead().BlockHeight
	f.logger.Debug("lastSignalledBlock set from execution head", "block height", lastSignalledBlock)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		targetBlock, err := f.sharedDB.GetLatestHeight(cctx)
		cancel()
		if err != nil {
			f.sleepRespectingContext(ctx, defaultBackoff)
			continue
		}

		if lastSignalledBlock > targetBlock {
			f.logger.Error("internal invariant has been broken. Follower EL is ahead of signer")
			return
		}

		blocksRemaining := targetBlock - lastSignalledBlock

		if blocksRemaining == 0 {
			f.sleepRespectingContext(ctx, time.Millisecond) // New payload will likely be available within milliseconds
			continue
		}
		f.logger.Debug("non-zero blocksRemaining", "blocksRemaining", blocksRemaining)

		limit := min(f.syncBatchSize, blocksRemaining)

		cctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		payloads, err := f.sharedDB.GetPayloadsSince(cctx, lastSignalledBlock+1, int(limit))
		cancel()
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
		f.logger.Debug("number of payloads returned", "number of payloads", len(payloads))

		for i := range payloads {
			p := payloads[i]
			select {
			case <-ctx.Done():
				return
			case f.payloadCh <- p:
				lastSignalledBlock = p.BlockHeight
			}
		}
	}
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
			if err := f.finalizeBlock(ctx, p.PayloadID, p.ExecutionPayload, ""); err != nil {
				f.logger.Error("Failed to process payload", "height", p.BlockHeight, "error", err)
				continue
			}
			f.logger.Info("Successfully processed payload", "height", p.BlockHeight)
		}
	}
}

func (f *Follower) getExecutionHead() *types.ExecutionHead {
	f.bbMutex.RLock()
	defer f.bbMutex.RUnlock()
	return f.bb.GetExecutionHead()
}

func (f *Follower) setExecutionHeadFromRPC(ctx context.Context) error {
	f.bbMutex.Lock()
	defer f.bbMutex.Unlock()
	return f.bb.SetExecutionHeadFromRPC(ctx)
}

func (f *Follower) finalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	f.bbMutex.Lock()
	defer f.bbMutex.Unlock()
	return f.bb.FinalizeBlock(ctx, payloadIDStr, executionPayloadStr, msgID)
}
