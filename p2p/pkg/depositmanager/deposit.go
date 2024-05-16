package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	bidderregistry "github.com/primevprotocol/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Store interface {
	GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error)
	SetBalance(bidder common.Address, windowNumber *big.Int, balance *big.Int) error
	GetBalanceForBlock(
		bidder common.Address,
		windowNumber *big.Int,
		blockNumber int64,
	) (*big.Int, error)
	SetBalanceForBlock(
		bidder common.Address,
		windowNumber *big.Int,
		balance *big.Int,
		blockNumber int64,
	) error
	RefundBalanceForBlock(
		bidder common.Address,
		windowNumber *big.Int,
		amount *big.Int,
		blockNumber int64,
	) error
	ClearBalances(windowNumber *big.Int) ([]*big.Int, error)
}

type BlockTracker interface {
	GetBlocksPerWindow() (*big.Int, error)
}

type DepositManager struct {
	blockTracker    BlockTracker
	store           Store
	evtMgr          events.EventManager
	blocksPerWindow atomic.Uint64 // todo: move to the store
	currentWindow   atomic.Int64  // todo: move to the store
	bidderRegs      chan *bidderregistry.BidderregistryBidderRegistered
	windowChan      chan *blocktracker.BlocktrackerNewWindow
	logger          *slog.Logger
}

func NewDepositManager(
	blockTracker BlockTracker,
	store Store,
	evtMgr events.EventManager,
	logger *slog.Logger,
) *DepositManager {
	return &DepositManager{
		blockTracker: blockTracker,
		store:        store,
		bidderRegs:   make(chan *bidderregistry.BidderregistryBidderRegistered),
		windowChan:   make(chan *blocktracker.BlocktrackerNewWindow),
		evtMgr:       evtMgr,
		logger:       logger,
	}
}

func (dm *DepositManager) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	ev1 := events.NewEventHandler(
		"NewWindow",
		func(window *blocktracker.BlocktrackerNewWindow) {
			dm.currentWindow.Store(window.Window.Int64())
			select {
			case <-egCtx.Done():
			case dm.windowChan <- window:
			}
		},
	)

	ev2 := events.NewEventHandler(
		"BidderRegistered",
		func(bidderReg *bidderregistry.BidderregistryBidderRegistered) {
			select {
			case <-egCtx.Done():
			case dm.bidderRegs <- bidderReg:
			}
		},
	)

	sub, err := dm.evtMgr.Subscribe(ev1, ev2)
	if err != nil {
		close(doneChan)
		return doneChan
	}

	eg.Go(func() error {
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("error in event subscription: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return nil
			case window := <-dm.windowChan:
				windowToClear := new(big.Int).Sub(window.Window, big.NewInt(2))
				windows, err := dm.store.ClearBalances(windowToClear)
				if err != nil {
					dm.logger.Error("failed to clear balances", "error", err, "window", windowToClear)
					return err
				}
				dm.logger.Info("cleared balances", "windows", windows)
			case bidderReg := <-dm.bidderRegs:
				blocksPerWindow, err := dm.getOrSetBlocksPerWindow()
				if err != nil {
					dm.logger.Error("failed to get blocks per window", "error", err)
					return err
				}

				effectiveStake := new(big.Int).Div(bidderReg.DepositedAmount, new(big.Int).SetUint64(blocksPerWindow))
				if err := dm.store.SetBalance(bidderReg.Bidder, bidderReg.WindowNumber, effectiveStake); err != nil {
					return err
				}
				dm.logger.Info("set balance", "bidder", bidderReg.Bidder, "window", bidderReg.WindowNumber, "amount", effectiveStake)
			}
		}
	})
	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			dm.logger.Error("error in DepositManager", "error", err)
		}
	}()

	return doneChan
}

func (dm *DepositManager) CheckAndDeductDeposit(
	ctx context.Context,
	address common.Address,
	bidAmountStr string,
	blockNumber int64,
) (func() error, error) {
	blocksPerWindow, err := dm.getOrSetBlocksPerWindow()
	if err != nil {
		dm.logger.Error("failed to get blocks per window", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	bidAmount, ok := new(big.Int).SetString(bidAmountStr, 10)
	if !ok {
		dm.logger.Error("parsing bid amount", "amount", bidAmountStr)
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse bid amount")
	}

	windowToCheck := new(big.Int).SetUint64((uint64(blockNumber)-1)/blocksPerWindow + 1)

	balanceForBlock, err := dm.store.GetBalanceForBlock(address, windowToCheck, blockNumber)
	if err != nil {
		dm.logger.Error("getting balance for block", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance for block: %v", err)
	}

	if balanceForBlock != nil {
		newBalance := new(big.Int).Sub(balanceForBlock, bidAmount)
		if newBalance.Cmp(big.NewInt(0)) < 0 {
			dm.logger.Error("insufficient balance", "balance", balanceForBlock.Uint64(), "bidAmount", bidAmount.Uint64())
			return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
		}

		if err := dm.store.SetBalanceForBlock(address, windowToCheck, newBalance, blockNumber); err != nil {
			dm.logger.Error("setting balance for block", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to set balance for block: %v", err)
		}
		return func() error {
			return dm.store.RefundBalanceForBlock(address, windowToCheck, bidAmount, blockNumber)
		}, nil
	}

	defaultBalance, err := dm.store.GetBalance(address, windowToCheck)
	if err != nil {
		dm.logger.Error("getting balance", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	if defaultBalance == nil {
		dm.logger.Error("bidder balance not found", "address", address.Hex(), "window", windowToCheck)
		return nil, status.Errorf(codes.FailedPrecondition, "balance not found")
	}

	if defaultBalance.Cmp(bidAmount) < 0 {
		dm.logger.Error("insufficient balance", "balance", defaultBalance, "bidAmount", bidAmount)
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
	}

	newBalance := new(big.Int).Sub(defaultBalance, bidAmount)
	if err := dm.store.SetBalanceForBlock(address, windowToCheck, newBalance, blockNumber); err != nil {
		dm.logger.Error("setting balance for block", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to set balance for block: %v", err)
	}

	return func() error {
		return dm.store.RefundBalanceForBlock(address, windowToCheck, bidAmount, blockNumber)
	}, nil
}

func (dm *DepositManager) getOrSetBlocksPerWindow() (uint64, error) {
	bpwCache := dm.blocksPerWindow.Load()

	if bpwCache != 0 {
		return bpwCache, nil
	}

	blocksPerWindow, err := dm.blockTracker.GetBlocksPerWindow()
	if err != nil {
		return 0, fmt.Errorf("failed to get blocks per window: %w", err)
	}

	blocksPerWindowUint64 := blocksPerWindow.Uint64()
	dm.blocksPerWindow.Store(blocksPerWindowUint64)

	return blocksPerWindowUint64, nil
}
