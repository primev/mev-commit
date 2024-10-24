package relayer

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	settlementgateway "github.com/primev/mev-commit/contracts-abi/clients/SettlementGateway"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type L1Gateway interface {
	Subscribe(ctx context.Context) (<-chan *l1gateway.L1gatewayTransferInitiated, <-chan error)
	FinalizeTransfer(ctx context.Context, recipient common.Address, amount *big.Int, transferIdx *big.Int) error
}

type SettlementGateway interface {
	Subscribe(ctx context.Context) (<-chan *settlementgateway.SettlementgatewayTransferInitiated, <-chan error)
	FinalizeTransfer(ctx context.Context, recipient common.Address, amount *big.Int, transferIdx *big.Int) error
}

type Relayer struct {
	logger            *slog.Logger
	l1Gateway         L1Gateway
	settlementGateway SettlementGateway
	metrics           *metrics
}

func NewRelayer(
	logger *slog.Logger,
	l1Gateway L1Gateway,
	settlementGateway SettlementGateway,
) *Relayer {
	return &Relayer{
		logger:            logger,
		l1Gateway:         l1Gateway,
		settlementGateway: settlementGateway,
		metrics:           newMetrics(),
	}
}

func (r *Relayer) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		r.metrics.initiatedTransfers,
		r.metrics.finalizedTransfers,
		r.metrics.failedFinalizations,
		r.metrics.initiatedTransfersValue,
		r.metrics.finalizedTransfersValue,
		r.metrics.failedFinalizationsValue,
	}
}

func (r *Relayer) Start(ctx context.Context) <-chan struct{} {
	l1Transfers, l1Err := r.l1Gateway.Subscribe(ctx)
	settlementTransfers, settlementErr := r.settlementGateway.Subscribe(ctx)

	done := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				r.logger.Info("relayer context done")
				return nil
			case err := <-l1Err:
				return err
			case upd := <-l1Transfers:
				r.metrics.initiatedTransfers.WithLabelValues("l1").Inc()
				r.metrics.initiatedTransfersValue.WithLabelValues("l1").Add(float64(upd.Amount.Int64()))
				err := r.settlementGateway.FinalizeTransfer(
					egCtx,
					upd.Recipient,
					upd.Amount,
					upd.TransferIdx,
				)
				if err != nil {
					r.logger.Error("error in settlement finalization", "recipient", upd.Recipient, "amount", upd.Amount, "transferIdx", upd.TransferIdx, "error", err)
					r.metrics.failedFinalizations.WithLabelValues("settlement").Inc()
					r.metrics.failedFinalizationsValue.WithLabelValues("settlement").Add(float64(upd.Amount.Int64()))
					continue
				}
				r.metrics.finalizedTransfers.WithLabelValues("settlement").Inc()
				r.metrics.finalizedTransfersValue.WithLabelValues("settlement").Add(float64(upd.Amount.Int64()))
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				r.logger.Info("relayer context done")
				return nil
			case err := <-settlementErr:
				return err
			case upd := <-settlementTransfers:
				r.metrics.initiatedTransfers.WithLabelValues("settlement").Inc()
				r.metrics.initiatedTransfersValue.WithLabelValues("settlement").Add(float64(upd.Amount.Int64()))
				err := r.l1Gateway.FinalizeTransfer(
					egCtx,
					upd.Recipient,
					upd.Amount,
					upd.TransferIdx,
				)
				if err != nil {
					r.logger.Error("error in l1 finalization", "recipient", upd.Recipient, "amount", upd.Amount, "transferIdx", upd.TransferIdx, "error", err)
					r.metrics.failedFinalizations.WithLabelValues("l1").Inc()
					r.metrics.failedFinalizationsValue.WithLabelValues("l1").Add(float64(upd.Amount.Int64()))
					continue
				}
				r.metrics.finalizedTransfers.WithLabelValues("l1").Inc()
				r.metrics.finalizedTransfersValue.WithLabelValues("l1").Add(float64(upd.Amount.Int64()))
			}
		}
	})

	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			r.logger.Error("error in relayer", "error", err)
		}
	}()

	return done
}
