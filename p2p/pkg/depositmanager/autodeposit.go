package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type BidderRegistryContract interface {
	MoveDepositToWindow(opts *bind.TransactOpts, fromWindow *big.Int, toWindow *big.Int) (*types.Transaction, error)
}

type AutoDepositTracker struct {
	deposits   map[uint64]bool
	windowChan chan *blocktracker.BlocktrackerNewWindow
	eventMgr   events.EventManager
	isWorking  bool
	brContract BidderRegistryContract
	optsGetter OptsGetter
	logger     *slog.Logger
	cancelFunc context.CancelFunc
}

func NewAutoDepositTracker(
	evtMgr events.EventManager,
	brContract BidderRegistryContract,
	optsGetter OptsGetter,
	logger *slog.Logger,
) *AutoDepositTracker {
	return &AutoDepositTracker{
		deposits:   make(map[uint64]bool),
		eventMgr:   evtMgr,
		brContract: brContract,
		optsGetter: optsGetter,
		windowChan: make(chan *blocktracker.BlocktrackerNewWindow, 1),
		logger:     logger,
	}
}

func (adt *AutoDepositTracker) DoAutoMoveToAnotherWindow(ads []*bidderapiv1.AutoDeposit) <-chan struct{} {
	if adt.isWorking {
		return nil
	}
	adt.isWorking = true

	for _, ad := range ads {
		adt.deposits[ad.WindowNumber.Value] = true
	}

	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(context.Background())
	egCtx, cancel := context.WithCancel(egCtx)
	adt.cancelFunc = cancel

	evt := events.NewEventHandler(
		"NewWindow",
		func(update *blocktracker.BlocktrackerNewWindow) {
			adt.logger.Info(
				"new window event",
				"window", update.Window,
			)
			select {
			case <-egCtx.Done():
			case adt.windowChan <- update:
			}
		},
	)

	sub, err := adt.eventMgr.Subscribe(evt)
	if err != nil {
		close(doneChan)
		return doneChan
	}

	eg.Go(func() error {
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			adt.logger.Info("event subscription context done")
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("error in event subscription: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				adt.logger.Info("context done")
				return nil
			case window := <-adt.windowChan:
				// logic for 3 windows for deposit
				fromWindow := new(big.Int).Sub(window.Window, big.NewInt(1))
				if _, ok := adt.deposits[fromWindow.Uint64()]; !ok {
					continue
				}
				toWindow := new(big.Int).Add(window.Window, big.NewInt(1))

				opts, err := adt.optsGetter(egCtx)
				if err != nil {
					return err
				}
				txn, err := adt.brContract.MoveDepositToWindow(opts, fromWindow, toWindow)
				if err != nil {
					return err
				}
				adt.logger.Info("move deposit to window", "hash", txn.Hash(), "from", fromWindow, "to", toWindow)
				delete(adt.deposits, fromWindow.Uint64())
				adt.deposits[toWindow.Uint64()] = true
			}
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			adt.logger.Error("error in errgroup", "err", err)
		}
		adt.isWorking = false
	}()

	return doneChan
}

func (adt *AutoDepositTracker) Stop() {
	if adt.cancelFunc != nil {
		adt.cancelFunc()
	}
}
