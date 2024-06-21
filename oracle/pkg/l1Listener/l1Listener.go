package l1Listener

import (
	"bytes"
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var checkInterval = 2 * time.Second

type L1Recorder interface {
	RecordL1Block(blockNum *big.Int, winner string) (*types.Transaction, error)
}

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner []byte, window int64) error
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type L1Listener struct {
	logger         *slog.Logger
	l1Client       EthClient
	winnerRegister WinnerRegister
	eventMgr       events.EventManager
	recorder       L1Recorder
	metrics        *metrics
}

func NewL1Listener(
	logger *slog.Logger,
	l1Client EthClient,
	winnerRegister WinnerRegister,
	evtMgr events.EventManager,
	recorder L1Recorder,
) *L1Listener {
	return &L1Listener{
		logger:         logger,
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		eventMgr:       evtMgr,
		recorder:       recorder,
		metrics:        newMetrics(),
	}
}

func (l *L1Listener) Metrics() []prometheus.Collector {
	return l.metrics.Collectors()
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return l.watchL1Block(egCtx)
	})

	evt := events.NewEventHandler(
		"NewL1Block",
		func(update *blocktracker.BlocktrackerNewL1Block) {
			l.logger.Info(
				"new L1 block event",
				"block", update.BlockNumber,
				"winner", update.Winner.String(),
				"window", update.Window,
			)
			err := l.winnerRegister.RegisterWinner(
				ctx,
				update.BlockNumber.Int64(),
				update.Winner.Bytes(),
				update.Window.Int64(),
			)
			if err != nil {
				l.logger.Error(
					"failed to register winner",
					"block", update.BlockNumber,
					"winner", update.Winner.String(),
					"error", err,
				)
				return
			}
			l.metrics.WinnerCount.Inc()
			l.metrics.WinnerRoundCount.WithLabelValues(update.Winner.String()).Inc()
		},
	)

	sub, err := l.eventMgr.Subscribe(evt)
	if err != nil {
		close(doneChan)
		return doneChan
	}

	eg.Go(func() error {
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			return egCtx.Err()
		case err := <-sub.Err():
			return err
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			l.logger.Error("L1listener error", "error", err)
		}
	}()

	l.logger.Info("L1Listener started")

	return doneChan
}

func (l *L1Listener) watchL1Block(ctx context.Context) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// TODO: change it to the store to not miss blocks, if oracle is down
	currentBlockNo, err := l.l1Client.BlockNumber(ctx)
	if err != nil {
		l.logger.Error("failed to get block number", "error", err)
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			blockNum, err := l.l1Client.BlockNumber(ctx)
			if err != nil {
				l.logger.Error("failed to get block number", "error", err)
				continue
			}

			if blockNum <= uint64(currentBlockNo) {
				continue
			}

			for b := uint64(currentBlockNo) + 1; b <= blockNum; b++ {
				header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(b)))
				if err != nil {
					l.logger.Error("failed to get header", "block", b, "error", err)
					continue
				}

				winner := string(bytes.ToValidUTF8(header.Extra, []byte("�")))

				l.logger.Info(
					"new L1 winner",
					"winner", winner,
					"block", header.Number.Int64(),
				)

				winnerPostingTxn, err := l.recorder.RecordL1Block(
					big.NewInt(0).SetUint64(b),
					winner,
				)
				if err != nil {
					l.logger.Error("failed to register winner for block", "block", b, "error", err)
					continue
				}

				l.metrics.WinnerPostedCount.Inc()
				l.metrics.LastSentNonce.Set(float64(winnerPostingTxn.Nonce()))

				l.logger.Info(
					"registered winner",
					"winner", winner,
					"block", header.Number.Int64(),
					"txn", winnerPostingTxn.Hash().String(),
				)
			}

			currentBlockNo = blockNum
		}
	}
}
