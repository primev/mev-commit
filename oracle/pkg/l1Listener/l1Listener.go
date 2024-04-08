package l1Listener

import (
	"bytes"
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
)

var checkInterval = 2 * time.Second

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner string) error
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type L1Listener struct {
	logger         *slog.Logger
	l1Client       EthClient
	winnerRegister WinnerRegister
	metrics        *metrics
}

func NewL1Listener(
	logger *slog.Logger,
	l1Client EthClient,
	winnerRegister WinnerRegister,
) *L1Listener {
	return &L1Listener{
		logger:         logger,
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		metrics:        newMetrics(),
	}
}

func (l *L1Listener) Metrics() []prometheus.Collector {
	return l.metrics.Collectors()
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		currentBlockNo := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blockNum, err := l.l1Client.BlockNumber(ctx)
				if err != nil {
					l.logger.Error("failed to get block number", "error", err)
					continue
				}

				if blockNum <= uint64(currentBlockNo) {
					continue
				}

				header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				if err != nil {
					l.logger.Error("failed to get header", "block", blockNum, "error", err)
					continue
				}

				winner := string(bytes.ToValidUTF8(header.Extra, []byte("�")))
				if len(winner) == 0 {
					l.logger.Warn("no winner registered", "block", header.Number.Int64())
					continue
				} else {
					err = l.winnerRegister.RegisterWinner(ctx, int64(blockNum), winner)
					if err != nil {
						l.logger.Error("failed to register winner for block", "block", blockNum, "error", err)
						return
					}

					l.metrics.WinnerRoundCount.WithLabelValues(winner).Inc()
					l.metrics.WinnerCount.Inc()

					l.logger.Info("registered winner", "winner", winner, "block", header.Number.Int64())
				}
				currentBlockNo = int(blockNum)
			}
		}

	}()

	return doneChan
}
