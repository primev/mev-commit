package gwcontract

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
)

type GatewayTransactor interface {
	FinalizeTransfer(
		opts *bind.TransactOpts,
		_recipient common.Address,
		_amount *big.Int,
		_counterpartyIdx *big.Int,
	) (*types.Transaction, error)
}

type Storage interface {
	StoreTransfer(
		ctx context.Context,
		transferIdx *big.Int,
		amount *big.Int,
		recipient common.Address,
		nonce uint64,
		chainHash common.Hash,
	) error
	MarkTransferSettled(ctx context.Context, transferIdx *big.Int) error
	IsSettled(ctx context.Context, transferIdx *big.Int) (bool, error)
}

type Monitor interface {
	WatchTx(hash common.Hash, nonce uint64) <-chan txmonitor.Result
}

type Gateway[EventType any] struct {
	monitor           Monitor
	listener          events.EventManager
	gatewayTransactor GatewayTransactor
	store             Storage
	optsGetter        func(context.Context) (*bind.TransactOpts, error)
	logger            *slog.Logger
}

func NewGateway[EventType any](
	logger *slog.Logger,
	monitor Monitor,
	listener events.EventManager,
	transactor GatewayTransactor,
	optsGetter func(context.Context) (*bind.TransactOpts, error),
	store Storage,
) *Gateway[EventType] {
	return &Gateway[EventType]{
		monitor:           monitor,
		listener:          listener,
		gatewayTransactor: transactor,
		store:             store,
		optsGetter:        optsGetter,
		logger:            logger,
	}
}

func (g *Gateway[EventType]) FinalizeTransfer(
	ctx context.Context,
	recipient common.Address,
	amount *big.Int,
	counterpartyIdx *big.Int,
) error {
	switch settled, err := g.store.IsSettled(ctx, counterpartyIdx); {
	case err != nil:
		g.logger.Error("failed to check if transfer is settled", "counterpartyIdx", counterpartyIdx, "error", err)
		return err
	case settled:
		g.logger.Debug("transfer already settled", "counterpartyIdx", counterpartyIdx)
		return nil
	}

	opts, err := g.optsGetter(ctx)
	if err != nil {
		g.logger.Error("failed to get auth", "error", err)
		return err
	}

	tx, err := g.gatewayTransactor.FinalizeTransfer(
		opts,
		recipient,
		amount,
		counterpartyIdx,
	)
	if err != nil {
		g.logger.Error(
			"failed to send transaction",
			"sender", opts.From,
			"receipient", recipient,
			"amount", amount,
			"error", err,
		)
		return err
	}

	err = g.store.StoreTransfer(
		ctx,
		counterpartyIdx,
		amount,
		recipient,
		tx.Nonce(),
		tx.Hash(),
	)
	if err != nil {
		g.logger.Error("failed to store transfer", "error", err)
		return err
	}

	res := g.monitor.WatchTx(tx.Hash(), tx.Nonce())
	select {
	case <-ctx.Done():
		return ctx.Err()
	case r := <-res:
		if r.Err != nil {
			g.logger.Error("transaction failed", "error", r.Err, "txHash", tx.Hash())
			return r.Err
		}
		if r.Receipt.Status != types.ReceiptStatusSuccessful {
			g.logger.Error("transaction status unsuccessful", "status", r.Receipt.Status, "txHash", tx.Hash())
			return fmt.Errorf("transaction status unsuccessful: %d", r.Receipt.Status)
		}
		return g.store.MarkTransferSettled(ctx, counterpartyIdx)
	}
}

func (g *Gateway[EventType]) Subscribe(ctx context.Context) (<-chan *EventType, <-chan error) {
	gatewayTransfers := make(chan *EventType)
	sub, err := g.listener.Subscribe(
		events.NewEventHandler(
			"TransferInitiated",
			func(upd *EventType) {
				select {
				case <-ctx.Done():
				case gatewayTransfers <- upd:
				}
			},
		),
	)

	errs := make(chan error, 1)
	if err != nil {
		errs <- err
		close(gatewayTransfers)
		return gatewayTransfers, errs
	}

	go func() {
		select {
		case <-ctx.Done():
		case err := <-sub.Err():
			errs <- err
			close(gatewayTransfers)
		}
		sub.Unsubscribe()
		close(errs)
	}()

	return gatewayTransfers, errs
}
