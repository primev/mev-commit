package autodepositor

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

var ErrNotRunning = fmt.Errorf("auto deposit tracker is not running")

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type BidderRegistryContract interface {
	DepositForWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
	DepositForWindow(opts *bind.TransactOpts, window *big.Int) (*types.Transaction, error)
	WithdrawFromWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
	GetDeposit(opts *bind.CallOpts, bidder common.Address, window *big.Int) (*big.Int, error)
}

type BlockTrackerContract interface {
	GetCurrentWindow() (*big.Int, error)
}

type DepositStore interface {
	// StoreDeposits stores the deposited windows.
	StoreDeposits(ctx context.Context, windows []*big.Int) error
	// ListDeposits lists the deposited windows upto and including the lastWindow.
	// If lastWindow is nil, it lists all deposits.
	ListDeposits(ctx context.Context, lastWindow *big.Int) ([]*big.Int, error)
	// ClearDeposits clears the deposits for the given windows.
	ClearDeposits(ctx context.Context, windows []*big.Int) error
	// IsDepositMade checks if the deposit is already made for the given window.
	IsDepositMade(ctx context.Context, window *big.Int) bool
}

type AutoDepositTracker struct {
	startMu             sync.Mutex
	isWorking           bool
	eventMgr            events.EventManager
	windowChan          chan *blocktracker.BlocktrackerNewWindow
	brContract          BidderRegistryContract
	btContract          BlockTrackerContract
	store               DepositStore
	optsGetter          OptsGetter
	currentOracleWindow atomic.Value
	oracleWindowOffset  *big.Int
	logger              *slog.Logger
	cancelFunc          context.CancelFunc
}

func New(
	evtMgr events.EventManager,
	brContract BidderRegistryContract,
	btContract BlockTrackerContract,
	optsGetter OptsGetter,
	store DepositStore,
	oracleWindowOffset *big.Int,
	logger *slog.Logger,
) *AutoDepositTracker {
	return &AutoDepositTracker{
		eventMgr:           evtMgr,
		brContract:         brContract,
		btContract:         btContract,
		optsGetter:         optsGetter,
		store:              store,
		oracleWindowOffset: oracleWindowOffset,
		windowChan:         make(chan *blocktracker.BlocktrackerNewWindow, 1),
		logger:             logger,
	}
}

func (adt *AutoDepositTracker) Start(
	ctx context.Context,
	startWindow, amount *big.Int,
) error {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	if adt.isWorking {
		return fmt.Errorf("auto deposit tracker is already running")
	}

	currentOracleWindow, err := adt.btContract.GetCurrentWindow()
	if err != nil {
		return fmt.Errorf("failed to get current window: %w", err)
	}
	adt.currentOracleWindow.Store(currentOracleWindow)

	// Withdraw from past windows first
	err = adt.withdrawFromPastWindows(ctx, currentOracleWindow)
	if err != nil {
		adt.logger.Error("failed to withdraw from past windows", "err", err)
		// Don't return error here, we can still continue with deposits
	}

	if startWindow == nil {
		startWindow = currentOracleWindow
		// adding + N as oracle runs N window behind
		startWindow = new(big.Int).Add(startWindow, adt.oracleWindowOffset)
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	egCtx, cancel := context.WithCancel(egCtx)
	adt.cancelFunc = cancel

	sub, err := adt.initSub(egCtx)

	if err != nil {
		return fmt.Errorf("error subscribing to event: %w", err)
	}

	err = adt.doInitialDeposit(ctx, startWindow, amount)
	if err != nil {
		return fmt.Errorf("failed to do initial deposit: %w", err)
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

func (adt *AutoDepositTracker) withdrawFromPastWindows(ctx context.Context, currentWindow *big.Int) error {
	// Find windows to withdraw from (all windows before the current one)
	withdrawWindows, err := adt.store.ListDeposits(ctx, new(big.Int).Sub(currentWindow, big.NewInt(1)))
	switch {
	case err != nil:
		adt.logger.Error("failed to list deposits", "err", err)
		return err
	case len(withdrawWindows) == 0:
		adt.logger.Info("no deposits to withdraw")
		return nil
	}

	adt.logger.Info("deposits to withdraw", "windows", withdrawWindows)
	opts, err := adt.optsGetter(ctx)
	if err != nil {
		return err
	}

	txn, err := adt.brContract.WithdrawFromWindows(opts, withdrawWindows)
	if err != nil {
		return err
	}

	adt.logger.Info("withdraw from windows", "hash", txn.Hash(), "windows", withdrawWindows)
	err = adt.store.ClearDeposits(ctx, withdrawWindows)
	if err != nil {
		return fmt.Errorf("failed to clear deposits: %w", err)
	}

	return nil
}

func (adt *AutoDepositTracker) doInitialDeposit(ctx context.Context, startWindow, amount *big.Int) error {
	nextWindow := new(big.Int).Add(startWindow, big.NewInt(1))
	newDeposits := []*big.Int{startWindow, nextWindow}

	// Check if the deposit is already made. If the nodes was down for a short period
	// and the deposits were already made, we should not make the deposit again.
	newDeposits = slices.DeleteFunc(newDeposits, func(i *big.Int) bool {
		return adt.isDeposited(ctx, i)
	})

	if len(newDeposits) == 0 {
		return nil
	}

	opts, err := adt.optsGetter(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transact opts: %w", err)
	}
	opts.Value = big.NewInt(0).Mul(amount, big.NewInt(int64(len(newDeposits))))

	// Make initial deposit for the first two windows
	_, err = adt.brContract.DepositForWindows(opts, newDeposits)
	if err != nil {
		return fmt.Errorf("failed to deposit for windows: %w", err)
	}

	return adt.store.StoreDeposits(ctx, newDeposits)
}

func (adt *AutoDepositTracker) initSub(egCtx context.Context) (events.Subscription, error) {
	evt := events.NewEventHandler(
		"NewWindow",
		func(update *blocktracker.BlocktrackerNewWindow) {
			adt.logger.Debug(
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
				adt.currentOracleWindow.Store(window.Window)
				err := adt.withdrawFromPastWindows(egCtx, window.Window)
				if err != nil {
					return fmt.Errorf("failed to withdraw from past windows: %w", err)
				}

				// Make deposit for the next window. The window event is N windows
				// behind the current window in progress.
				nextWindow := new(big.Int).Add(window.Window, adt.oracleWindowOffset)
				nextWindow = new(big.Int).Add(nextWindow, big.NewInt(1))
				if adt.isDeposited(egCtx, nextWindow) {
					continue
				}

				opts, err := adt.optsGetter(egCtx)
				if err != nil {
					return err
				}
				opts.Value = amount

				txn, err := adt.brContract.DepositForWindow(opts, nextWindow)
				if err != nil {
					return err
				}
				adt.logger.Info(
					"deposited to next window",
					"hash", txn.Hash(),
					"window", nextWindow,
					"amount", amount,
				)
				err = adt.store.StoreDeposits(egCtx, []*big.Int{nextWindow})
				if err != nil {
					return fmt.Errorf("failed to store deposits: %w", err)
				}
			}
		}
	})
}

func (adt *AutoDepositTracker) Stop() ([]*big.Int, error) {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	if !adt.isWorking {
		return nil, ErrNotRunning
	}
	if adt.cancelFunc != nil {
		adt.cancelFunc()
	}

	windowNumbers, err := adt.store.ListDeposits(context.Background(), nil)
	if err != nil {
		adt.logger.Error("failed to list deposits", "err", err)
	}

	adt.isWorking = false

	adt.logger.Info("stop auto deposit tracker", "windowsToWithdraw", windowNumbers)
	return windowNumbers, nil
}

func (adt *AutoDepositTracker) IsWorking() bool {
	adt.startMu.Lock()
	defer adt.startMu.Unlock()

	return adt.isWorking
}

func (adt *AutoDepositTracker) GetStatus() (map[uint64]bool, bool, *big.Int) {
	adt.startMu.Lock()
	isWorking := adt.isWorking
	adt.startMu.Unlock()

	windows, err := adt.store.ListDeposits(context.Background(), nil)
	if err != nil {
		adt.logger.Error("failed to list deposits", "err", err)
	}
	deposits := make(map[uint64]bool)
	for _, w := range windows {
		deposits[w.Uint64()] = true
	}

	var currentOracleWindow *big.Int
	if val := adt.currentOracleWindow.Load(); val != nil {
		currentOracleWindow = val.(*big.Int)
	}

	return deposits, isWorking, currentOracleWindow
}

func (adt *AutoDepositTracker) isDeposited(ctx context.Context, window *big.Int) bool {
	if adt.store.IsDepositMade(ctx, window) {
		return true
	}

	opts, err := adt.optsGetter(ctx)
	if err != nil {
		adt.logger.Error("failed to get transact opts", "err", err)
		return false
	}

	// fallback to contract call if the local state was not flushed properly
	deposit, err := adt.brContract.GetDeposit(&bind.CallOpts{
		Context: ctx,
		From:    opts.From,
	}, opts.From, window)
	if err != nil {
		adt.logger.Error("failed to get deposit", "err", err)
		return false
	}

	return deposit.Cmp(big.NewInt(0)) > 0
}
