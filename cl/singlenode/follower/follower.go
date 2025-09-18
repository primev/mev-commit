package follower

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/primev/mev-commit/cl/singlenode/payloadstore"
	"github.com/primev/mev-commit/cl/types"
	"golang.org/x/sync/errgroup"
)

type Follower struct {
	logger                *slog.Logger
	sharedDB              payloadDB
	syncBatchSize         uint64
	payloadCh             chan types.PayloadInfo
	bbMutex               sync.RWMutex
	bb                    blockBuilder
	healthAddr            string
	syncStopped           chan struct{}
	handlePayloadsStopped chan struct{}
	healthStopped         chan struct{}
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
	ctx context.Context,
	logger *slog.Logger,
	postgresDSN string,
	redisURL string,
	syncBatchSize uint64,
	bb blockBuilder,
	healthAddr string,
) (*Follower, error) {
	var sharedDB payloadDB
	if postgresDSN != "" {
		pgRepo, err := payloadstore.NewPostgresRepository(ctx, postgresDSN, logger)
		if err != nil {
			return nil, err
		}
		sharedDB = pgRepo
	} else if redisURL != "" {
		redisRepo, err := payloadstore.NewRedisRepositoryFromURL(ctx, redisURL, logger)
		if err != nil {
			return nil, err
		}
		sharedDB = redisRepo
	} else {
		return nil, errors.New("postgresDSN or redisURL must be provided")
	}
	if syncBatchSize == 0 {
		return nil, errors.New("sync batch size must be greater than 0")
	}
	return &Follower{
		logger:                logger,
		sharedDB:              sharedDB,
		syncBatchSize:         syncBatchSize,
		payloadCh:             make(chan types.PayloadInfo),
		bb:                    bb,
		healthAddr:            healthAddr,
		healthStopped:         make(chan struct{}),
		syncStopped:           make(chan struct{}),
		handlePayloadsStopped: make(chan struct{}),
	}, nil
}

func (f *Follower) Start(ctx context.Context) <-chan struct{} {

	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer close(f.healthStopped)
		mux := http.NewServeMux()
		mux.HandleFunc("/health", f.healthHandler)
		server := &http.Server{Addr: f.healthAddr, Handler: mux}
		f.logger.Info("Health endpoint listening", "address", f.healthAddr)

		go func() {
			<-egCtx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		defer close(f.handlePayloadsStopped)
		return f.handlePayloads(egCtx)
	})

	eg.Go(func() error {
		defer close(f.syncStopped)
		f.logger.Info("Starting sync from shared DB")
		return f.syncFromSharedDB(egCtx)
	})

	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			f.logger.Error("follower failed, exiting", "error", err)
		}
	}()

	return done
}

func (f *Follower) healthHandler(w http.ResponseWriter, r *http.Request) {

	select {
	case <-f.healthStopped:
		http.Error(w, "health loop has stopped", http.StatusServiceUnavailable)
		return
	case <-f.syncStopped:
		http.Error(w, "sync from shared DB has stopped", http.StatusServiceUnavailable)
		return
	case <-f.handlePayloadsStopped:
		http.Error(w, "handle payloads loop has stopped", http.StatusServiceUnavailable)
		return
	default:
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (f *Follower) syncFromSharedDB(ctx context.Context) error {
	if f.getExecutionHead() == nil {
		if err := f.setExecutionHeadFromRPC(ctx); err != nil {
			f.logger.Error("failed to set execution head from rpc", "error", err)
			return err
		}
		f.logger.Debug("set execution head from rpc")
	}

	lastSignalledBlock := f.getExecutionHead().BlockHeight
	f.logger.Debug("lastSignalledBlock set from execution head", "block height", lastSignalledBlock)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
			return fmt.Errorf("internal invariant has been broken. Follower EL is ahead of signer")
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
				return ctx.Err()
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

func (f *Follower) handlePayloads(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
