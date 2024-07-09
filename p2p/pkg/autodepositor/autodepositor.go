package autodepositor

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type BidderRegistryContract interface {
	DepositForWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
	WithdrawFromWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
}

type BlockTrackerContract interface {
	GetCurrentWindow() (*big.Int, error)
}

type AutoDepositTracker struct {
	startMu    sync.Mutex
	isWorking  bool
	eventMgr   events.EventManager
	deposits   sync.Map
	windowChan chan *blocktracker.BlocktrackerNewWindow
	brContract BidderRegistryContract
	btContract BlockTrackerContract
	initialAmount *big.Int
	optsGetter OptsGetter
	logger     *slog.Logger
	cancelFunc context.CancelFunc
}

func New(
	evtMgr events.EventManager,
	brContract BidderRegistryContract,
	btContract BlockTrackerContract,
	optsGetter OptsGetter,
	initialAmount *big.Int,
	logger *slog.Logger,
) *AutoDepositTracker {
	return &AutoDepositTracker{
		eventMgr:   evtMgr,
		brContract: brContract,
		btContract: btContract,
		optsGetter: optsGetter,
		windowChan: make(chan *blocktracker.BlocktrackerNewWindow, 1),
		initialAmount: initialAmount,
		logger:     logger,
	}
}

func (adt *AutoDepositTracker) Start(ctx context.Context) <-chan struct{} {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	doneChan := make(chan struct{})

	if adt.isWorking {
		close(doneChan)
		adt.logger.Error("failed to subscribe to events", "error", fmt.Errorf("auto deposit tracker is already running"))
		return doneChan
	}

	currentWindow, err := adt.btContract.GetCurrentWindow()
	if err != nil {
		close(doneChan)
		adt.logger.Error("failed to get current window", "error", err)
		return doneChan
	}
	// adding +2 as oracle runs two windows behind
	currentWindow = new(big.Int).Add(currentWindow, big.NewInt(2))

	opts, err := adt.optsGetter(ctx)
	if err != nil {
		close(doneChan)
		adt.logger.Error("failed to get transact opts", "error", err)
		return doneChan
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	egCtx, cancel := context.WithCancel(egCtx)
	adt.cancelFunc = cancel

	sub, err := adt.initSub(egCtx, opts)

	if err != nil {
		close(doneChan)
		adt.logger.Error("error subscribing to event", "error", err)
		return doneChan
	}

	err = adt.doInitialDeposit(opts, currentWindow, adt.initialAmount)
	if err != nil {
		close(doneChan)
		adt.logger.Error("failed to do initial deposit", "error", err)
		return doneChan
	}

	adt.startAutodeposit(egCtx, eg, adt.initialAmount, sub)

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			adt.logger.Error("error in errgroup", "error", err)
		}
		adt.startMu.Lock()
		adt.isWorking = false
		adt.startMu.Unlock()
	}()
	
	return doneChan
}

func (adt *AutoDepositTracker) StartFromApi(
	ctx context.Context,
	startWindow, amount *big.Int,
) error {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	if adt.isWorking {
		return fmt.Errorf("auto deposit tracker is already running")
	}

	opts, err := adt.optsGetter(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transact opts, err: %w", err)
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	egCtx, cancel := context.WithCancel(egCtx)
	adt.cancelFunc = cancel

	sub, err := adt.initSub(egCtx, opts)

	if err != nil {
		return fmt.Errorf("error subscribing to event: %w", err)
	}

	err = adt.doInitialDeposit(opts, startWindow, amount)
	if err != nil {
		return fmt.Errorf("failed to do initial deposit, err: %w", err)
	}

	adt.startAutodeposit(egCtx, eg, amount, sub)

	started := make(chan struct{})
	go func() {
		close(started)
		if err := eg.Wait(); err != nil {
			adt.logger.Error("error in errgroup", "err", err)
		}
		adt.startMu.Lock()
		adt.isWorking = false
		adt.startMu.Unlock()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-started:
		adt.isWorking = true
	}
	return nil
}

func (adt *AutoDepositTracker) doInitialDeposit(opts *bind.TransactOpts, startWindow, amount *big.Int) error {
	nextWindow := new(big.Int).Add(startWindow, big.NewInt(1))

	opts.Value = big.NewInt(0).Mul(amount, big.NewInt(2))

	// Make initial deposit for the first two windows
	_, err := adt.brContract.DepositForWindows(opts, []*big.Int{startWindow, nextWindow})
	if err != nil {
		return err
	}

	return nil
}

func (adt *AutoDepositTracker) initSub(egCtx context.Context, opts *bind.TransactOpts) (events.Subscription, error) {
	evt1 := events.NewEventHandler(
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
	evt2 := events.NewEventHandler(
		"BidderRegistered",
		func(bidderReg *bidderregistry.BidderregistryBidderRegistered) {
			if bidderReg.Bidder.Cmp(opts.From) != 0 {
				return
			}
			adt.logger.Info(
				"bidder registered event",
				"bidder", bidderReg.Bidder.String(),
				"window", bidderReg.WindowNumber,
			)
			adt.deposits.Store(bidderReg.WindowNumber.Uint64(), true)
		},
	)

	evt3 := events.NewEventHandler(
		"BidderWithdrawal",
		func(bidderReg *bidderregistry.BidderregistryBidderWithdrawal) {
			if bidderReg.Bidder.Cmp(opts.From) != 0 {
				return
			}
			adt.logger.Info(
				"bidder withdrawal event",
				"bidder", bidderReg.Bidder.String(),
				"window", bidderReg.Window,
			)
			adt.deposits.Delete(bidderReg.Window.Uint64())
		},
	)

	sub, err := adt.eventMgr.Subscribe(evt1, evt2, evt3)
	if err != nil {
		return nil, fmt.Errorf("error subscribing to event: %w", err)
	}
	return sub, nil
}

func (adt *AutoDepositTracker) startAutodeposit(egCtx context.Context, eg *errgroup.Group, amount *big.Int, sub events.Subscription) {
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				adt.logger.Info("auto deposit tracker context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("error in autodeposit event subscription: %w", err)
			case window := <-adt.windowChan:
				withdrawWindows := make([]*big.Int, 0)
				adt.deposits.Range(func(key, value interface{}) bool {
					if key.(uint64) < window.Window.Uint64() {
						withdrawWindows = append(withdrawWindows, new(big.Int).SetUint64(key.(uint64)))
					}
					return true
				})

				if len(withdrawWindows) > 0 {
					opts, err := adt.optsGetter(egCtx)
					if err != nil {
						return err
					}
					txn, err := adt.brContract.WithdrawFromWindows(opts, withdrawWindows)
					if err != nil {
						return err
					}
					adt.logger.Info("withdraw from windows", "hash", txn.Hash(), "windows", withdrawWindows)
				}

				// Make deposit for the next window. The window event is 2 windows
				// behind the current window in progress. So we need to make deposit
				// for the next window.
				nextWindow := new(big.Int).Add(window.Window, big.NewInt(3))
				if _, ok := adt.deposits.Load(nextWindow.Uint64()); ok {
					continue
				}

				opts, err := adt.optsGetter(egCtx)
				if err != nil {
					return err
				}
				opts.Value = amount

				txn, err := adt.brContract.DepositForWindows(opts, []*big.Int{nextWindow})
				if err != nil {
					return err
				}
				adt.logger.Info(
					"deposited to next window",
					"hash", txn.Hash(),
					"window", nextWindow,
					"amount", amount,
				)
			}
		}
	})
}

func (adt *AutoDepositTracker) Stop() ([]*big.Int, error) {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	if !adt.isWorking {
		return nil, fmt.Errorf("auto deposit tracker is not running")
	}
	if adt.cancelFunc != nil {
		adt.cancelFunc()
	}
	var windowNumbers []*big.Int

	adt.deposits.Range(func(key, value interface{}) bool {
		windowNumbers = append(windowNumbers, new(big.Int).SetUint64(key.(uint64)))
		adt.deposits.Delete(key)
		return true
	})

	slices.SortFunc(windowNumbers, func(i, j *big.Int) int {
		return i.Cmp(j)
	})

	adt.isWorking = false

	adt.logger.Info("stop auto deposit tracker", "windows", windowNumbers)
	return windowNumbers, nil
}

func (adt *AutoDepositTracker) IsWorking() bool {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	return adt.isWorking
}

func (adt *AutoDepositTracker) GetStatus() (map[uint64]bool, bool) {
	adt.startMu.Lock()
	isWorking := adt.isWorking
	adt.startMu.Unlock()

	deposits := make(map[uint64]bool)
	adt.deposits.Range(func(key, value interface{}) bool {
		deposits[key.(uint64)] = value.(bool)
		return true
	})
	return deposits, isWorking
}
