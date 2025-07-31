package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BidderRegistryContract interface {
	GetDeposit(opts *bind.CallOpts, bidder common.Address, window *big.Int) (*big.Int, error)
}

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

type DepositManager struct {
	store           Store
	evtMgr          events.EventManager
	bidderRegistry  BidderRegistryContract
	blocksPerWindow uint64
	bidderRegs      chan *bidderregistry.BidderregistryBidderRegistered
	windowChan      chan *blocktracker.BlocktrackerNewWindow
	logger          *slog.Logger
}

func NewDepositManager(
	blocksPerWindow uint64,
	store Store,
	evtMgr events.EventManager,
	bidderRegistry BidderRegistryContract,
	logger *slog.Logger,
) *DepositManager {
	return &DepositManager{
		store:           store,
		blocksPerWindow: blocksPerWindow,
		bidderRegistry:  bidderRegistry,
		bidderRegs:      make(chan *bidderregistry.BidderregistryBidderRegistered),
		windowChan:      make(chan *blocktracker.BlocktrackerNewWindow),
		evtMgr:          evtMgr,
		logger:          logger,
	}
}

func (dm *DepositManager) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	ev1 := events.NewEventHandler(
		"NewWindow",
		func(window *blocktracker.BlocktrackerNewWindow) {
			select {
			case <-egCtx.Done():
				dm.logger.Info("new window context done")
			case dm.windowChan <- window:
			}
		},
	)

	ev2 := events.NewEventHandler(
		"BidderRegistered",
		func(bidderReg *bidderregistry.BidderregistryBidderRegistered) {
			select {
			case <-egCtx.Done():
				dm.logger.Info("bidder registered context done")
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
			dm.logger.Info("event subscription context done")
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("error in event subscription: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				dm.logger.Info("clear balances set balances context done")
				return nil
			case window := <-dm.windowChan:
				windowToClear := new(big.Int).Sub(window.Window, big.NewInt(1))
				windows, err := dm.store.ClearBalances(windowToClear)
				if err != nil {
					dm.logger.Error("failed to clear balances", "error", err, "window", windowToClear)
					return err
				}
				dm.logger.Info("cleared balances", "windows", windows)
			case bidderReg := <-dm.bidderRegs:
				effectiveStake := new(big.Int).Div(bidderReg.DepositedAmount, new(big.Int).SetUint64(dm.blocksPerWindow))
				if err := dm.store.SetBalance(bidderReg.Bidder, bidderReg.WindowNumber, effectiveStake); err != nil {
					dm.logger.Error("setting balance", "error", err)
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

// TODO: Add check in provider node to see if bidder has requested a withdrawal.
func (dm *DepositManager) CheckAndDeductDeposit(
	ctx context.Context,
	address common.Address,
	bidAmountStr string,
	blockNumber int64,
) (func() error, error) {
	bidAmount, ok := new(big.Int).SetString(bidAmountStr, 10)
	if !ok {
		dm.logger.Error("parsing bid amount", "amount", bidAmountStr)
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse bid amount")
	}

	windowToCheck := new(big.Int).SetUint64((uint64(blockNumber)-1)/dm.blocksPerWindow + 1)

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

	defaultBalance, err := dm.getBalanceForWindow(ctx, address, windowToCheck)
	if err != nil {
		return nil, err
	}

	if defaultBalance == nil {
		dm.logger.Error("bidder balance not found", "address", address.Hex(), "window", windowToCheck)
		return nil, status.Errorf(codes.FailedPrecondition, "balance not found for window %s", windowToCheck.String())
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

// fallback to contract if balance not found in store
func (dm *DepositManager) getBalanceForWindow(
	ctx context.Context,
	address common.Address,
	windowNumber *big.Int,
) (*big.Int, error) {
	balance, err := dm.store.GetBalance(address, windowNumber)
	if err != nil {
		dm.logger.Error("getting balance", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	if balance == nil {
		dm.logger.Info("balance not found in store", "address", address.Hex(), "window", windowNumber)
		balance, err = dm.bidderRegistry.GetDeposit(&bind.CallOpts{
			Context: ctx,
		}, address, windowNumber)
		if err != nil {
			dm.logger.Error("getting deposit from contract", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get deposit: %v", err)
		}

		// The set balance will only set the max amount that can be used in a block for
		// the given window. The actual balance will be deducted in the CheckAndDeductDeposit.
		// This prevents the need to synchrnoize the balance update from the events. They
		// update the same value.
		effectiveBalance := new(big.Int).Div(balance, new(big.Int).SetUint64(dm.blocksPerWindow))
		if err := dm.store.SetBalance(address, windowNumber, effectiveBalance); err != nil {
			dm.logger.Error("setting balance", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to set balance: %v", err)
		}
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil, nil
	}

	return balance, nil
}
