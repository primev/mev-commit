package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type BidderRegistryContract interface {
	MoveDepositToWindow(opts *bind.TransactOpts, fromWindow *big.Int, toWindow *big.Int) (*types.Transaction, error)
	WithdrawFromSpecificWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
}

type AutoDepositTracker struct {
	deposits   map[uint64]bool
	windowChan chan *blocktracker.BlocktrackerNewWindow
	eventMgr   events.EventManager
	isWorking  atomic.Bool
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

func (adt *AutoDepositTracker) DoAutoMoveToAnotherWindow(ctx context.Context, ads []*bidderapiv1.AutoDeposit) error {
	if adt.isWorking.Load() {
		return fmt.Errorf("auto deposit tracker is already working")
	}
	adt.isWorking.Store(true)

	for _, ad := range ads {
		adt.deposits[ad.WindowNumber.Value] = true
	}

	eg, egCtx := errgroup.WithContext(ctx)
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
		return fmt.Errorf("error subscribing to event: %w", err)
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
				toWindow := new(big.Int).Add(window.Window, big.NewInt(2))

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

	started := make(chan struct{})
	go func() {
		close(started)
		if err := eg.Wait(); err != nil {
			adt.logger.Error("error in errgroup", "err", err)
		}
		adt.isWorking.Store(false)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-started:
	}
	return nil
}

func (adt *AutoDepositTracker) Stop() (*bidderapiv1.CancelAutoDepositResponse, error) {
	if !adt.isWorking.Load() {
		return nil, fmt.Errorf("auto deposit tracker is not running")
	}
	if adt.cancelFunc != nil {
		adt.cancelFunc()
	}
	var windowNumbers []*wrapperspb.UInt64Value

	for i := range adt.deposits {
		windowNumbers = append(windowNumbers, &wrapperspb.UInt64Value{Value: i})
		delete(adt.deposits, i)
	}
	adt.logger.Info("stop auto deposit tracker", "windows", windowNumbers)
	return &bidderapiv1.CancelAutoDepositResponse{
		WindowNumbers: windowNumbers,
	}, nil
}

func (adt *AutoDepositTracker) IsWorking() bool {
	return adt.isWorking.Load()
}

func (adt *AutoDepositTracker) WithdrawAutoDeposit(ctx context.Context, windowNumbers []*wrapperspb.UInt64Value) error {
	adt.logger.Info("withdraw auto deposit")

	if len(windowNumbers) == 0 {
		return nil
	}
	var windows []*big.Int
	for _, windowNumber := range windowNumbers {
		windows = append(windows, new(big.Int).SetUint64(windowNumber.Value))
	}

	sort.Slice(windows, func(i, j int) bool {
		return windows[i].Cmp(windows[j]) < 0
	})

	lastWindowNumber := windows[len(windows)-1]
	eg, egCtx := errgroup.WithContext(ctx)
	egCtx, cancel := context.WithCancel(egCtx)

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
		cancel()
		return fmt.Errorf("error subscribing to event: %w", err)
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
			case curWindow := <-adt.windowChan:
				if curWindow.Window.Cmp(lastWindowNumber) > 0 {
					opts, err := adt.optsGetter(egCtx)
					if err != nil {
						return err
					}
					txn, err := adt.brContract.WithdrawFromSpecificWindows(opts, windows)
					if err != nil {
						return err
					}
					adt.logger.Info("withdraw from specific windows", "hash", txn.Hash(), "windows", windows)
					cancel()
					return nil
				}
				adt.logger.Info("current window is less or equal latest deposit window, waiting...", "currentWindow", curWindow.Window, "latestDepositWindow", lastWindowNumber)
			}
		}
	})

	started := make(chan struct{})
	go func() {
		close(started)
		if err := eg.Wait(); err != nil {
			adt.logger.Error("error in errgroup", "err", err)
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-started:
	}
	return nil
}

func (adt *AutoDepositTracker) GetStatus() (map[uint64]bool, bool) {
	return adt.deposits, adt.isWorking.Load()
}
