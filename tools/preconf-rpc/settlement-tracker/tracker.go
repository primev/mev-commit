package tracker

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
	"golang.org/x/sync/errgroup"
)

type BidderClient interface {
	SubscribeSettlements(ctx context.Context) <-chan bidder.SettlementMsg
	SubscribePayments(ctx context.Context) <-chan bidder.PaymentMsg
}

type Store interface {
	UpdateSettlementStatus(
		ctx context.Context,
		txnHash common.Hash,
		isSlashed bool,
		provider common.Address,
	) error
	UpdateSettlementPayment(
		ctx context.Context,
		txnHash common.Hash,
		payment *big.Int,
		refund *big.Int,
	) error
}

type tracker struct {
	client BidderClient
	store  Store
	logger *slog.Logger
}

func NewTracker(
	client BidderClient,
	store Store,
	logger *slog.Logger,
) *tracker {
	return &tracker{
		client: client,
		store:  store,
		logger: logger,
	}
}

func (t *tracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}

			sub := t.client.SubscribeSettlements(egCtx)
			for {
				select {
				case <-egCtx.Done():
					return egCtx.Err()
				case msg := <-sub:
					// Process settlement message
					if err := t.store.UpdateSettlementStatus(
						egCtx,
						common.HexToHash(msg.TransactionHash),
						msg.IsSlash,
						common.HexToAddress(msg.Provider),
					); err != nil {
						t.logger.Error(
							"Failed to update settlement status",
							"error", err,
							"txnHash", msg.TransactionHash,
							"isSlash", msg.IsSlash,
							"provider", msg.Provider,
						)
					} else {
						t.logger.Info("Updated settlement status",
							"txnHash", msg.TransactionHash,
							"isSlash", msg.IsSlash,
							"provider", msg.Provider,
						)
					}
				}
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}

			sub := t.client.SubscribePayments(egCtx)
			for {
				select {
				case <-egCtx.Done():
					return egCtx.Err()
				case msg := <-sub:
					// Process payment message
					if err := t.store.UpdateSettlementPayment(
						egCtx,
						common.HexToHash(msg.TransactionHash),
						msg.Payment,
						msg.Refund,
					); err != nil {
						t.logger.Error(
							"Failed to update settlement payment",
							"error", err,
							"txnHash", msg.TransactionHash,
							"payment", msg.Payment,
							"refund", msg.Refund,
						)
					} else {
						t.logger.Info(
							"Updated settlement payment",
							"txnHash", msg.TransactionHash,
							"payment", msg.Payment,
							"refund", msg.Refund,
						)
					}
				}
			}
		}
	})

	go func() {
		defer close(done)

		if err := eg.Wait(); err != nil {
			t.logger.Error("Tracker encountered an error", "error", err)
		}
	}()
	return done
}
