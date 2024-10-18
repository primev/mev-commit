package relayer

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	settlementgateway "github.com/primev/mev-commit/contracts-abi/clients/SettlementGateway"
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
				err := r.settlementGateway.FinalizeTransfer(
					egCtx,
					upd.Recipient,
					upd.Amount,
					upd.TransferIdx,
				)
				if err != nil {
					r.logger.Error("error in settlement finalization", "error", err)
					return err
				}
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
				err := r.l1Gateway.FinalizeTransfer(
					egCtx,
					upd.Recipient,
					upd.Amount,
					upd.TransferIdx,
				)
				if err != nil {
					r.logger.Error("error in l1 finalization", "error", err)
					return err
				}
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
