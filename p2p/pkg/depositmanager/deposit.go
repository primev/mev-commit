package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	bidderregistry "github.com/primevprotocol/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	preconfcontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/preconf"
	"github.com/primevprotocol/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BidderRegistry interface {
	CheckBidderDeposit(context.Context, common.Address, *big.Int, *big.Int) bool
	GetMinDeposit(ctx context.Context) (*big.Int, error)
}

type Store interface {
	GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error)
	SetBalance(bidder common.Address, windowNumber *big.Int, balance *big.Int) error
	GetBalanceForBlock(bidder common.Address, blockNumber int64) (*big.Int, error)
	SetBalanceForBlock(bidder common.Address, balance *big.Int, blockNumber int64) error
	RefundBalanceForBlock(bidder common.Address, amount *big.Int, blockNumber int64) error
}

type BlockTracker interface {
	GetBlocksPerWindow() (*big.Int, error)
}

type DepositManager struct {
	bidderRegistry  BidderRegistry
	blockTracker    BlockTracker
	commitmentDA    preconfcontract.Interface
	store           Store
	evtMgr          events.EventManager
	blocksPerWindow atomic.Uint64 // todo: move to the store
	minDeposit      atomic.Int64  // todo: move to the store
	currentWindow   atomic.Int64  // todo: move to the store
	bidderRegs      chan *bidderregistry.BidderregistryBidderRegistered
	logger          *slog.Logger
}

func NewDepositManager(
	br BidderRegistry,
	blockTracker BlockTracker,
	commitmentDA preconfcontract.Interface,
	store Store,
	evtMgr events.EventManager,
	logger *slog.Logger,
) *DepositManager {
	return &DepositManager{
		bidderRegistry: br,
		blockTracker:   blockTracker,
		commitmentDA:   commitmentDA,
		store:          store,
		bidderRegs:     make(chan *bidderregistry.BidderregistryBidderRegistered),
		evtMgr:         evtMgr,
		logger:         logger,
	}
}

func (dm *DepositManager) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	startWg := sync.WaitGroup{}
	startWg.Add(2)

	eg.Go(func() error {
		ev1 := events.NewEventHandler(
			"NewWindow",
			func(window *blocktracker.BlocktrackerNewWindow) {
				dm.currentWindow.Store(window.Window.Int64())
			},
		)

		sub1, err := dm.evtMgr.Subscribe(ev1)
		if err != nil {
			return fmt.Errorf("failed to subscribe to NewWindow event: %w", err)
		}
		defer sub1.Unsubscribe()

		ev2 := events.NewEventHandler(
			"BidderRegistered",
			func(bidderReg *bidderregistry.BidderregistryBidderRegistered) {
				// todo: do we need to check if commiter is connected to this bidder?
				select {
				case <-egCtx.Done():
				case dm.bidderRegs <- bidderReg:
				}
			},
		)

		sub2, err := dm.evtMgr.Subscribe(ev2)
		if err != nil {
			return fmt.Errorf("failed to subscribe to BidderRegistered event: %w", err)
		}
		defer sub2.Unsubscribe()

		startWg.Done()

		select {
		case <-egCtx.Done():
			return nil
		case err := <-sub1.Err():
			return fmt.Errorf("error in NewWindow event subscription: %w", err)
		case err := <-sub2.Err():
			return fmt.Errorf("error in BidderRegistered event subscription: %w", err)
		}
	})

	eg.Go(func() error {
		startWg.Done()

		for {
			select {
			case <-egCtx.Done():
				return nil
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
			}

		}
	})
	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			dm.logger.Error("error in DepositManager", "error", err)
		}
	}()

	startWg.Wait()

	return doneChan
}

func (dm *DepositManager) CheckAndDeductDeposit(ctx context.Context, address common.Address, bidAmountStr string, blockNumber int64) (*big.Int, error) {
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

	// adding 2 to the current window, bcs oracle is 2 windows behind
	possibleWindow := big.NewInt(dm.currentWindow.Load() + 2)

	windowToCheck := new(big.Int).SetUint64((uint64(blockNumber)-1)/blocksPerWindow + 1)
	if windowToCheck.Cmp(possibleWindow) < 0 {
		dm.logger.Error("window is too old", "window", windowToCheck.Uint64(), "possibleWindow", possibleWindow.Uint64())
		return nil, status.Errorf(codes.FailedPrecondition, "window is too old")
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

	dm.logger.Info("checking bidder deposit",
		"stake", defaultBalance.Uint64(),
		"blocksPerWindow", dm.blocksPerWindow.Load(),
		"minStake", dm.minDeposit.Load(),
		"window", windowToCheck.Uint64(),
		"address", address.Hex(),
	)

	balanceForBlock, err := dm.store.GetBalanceForBlock(address, blockNumber)
	if err != nil {
		dm.logger.Error("getting balance for block", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance for block: %v", err)
	}

	var newBalance *big.Int

	balanceToCheck := balanceForBlock
	if balanceToCheck == nil {
		balanceToCheck = defaultBalance
	}

	// Check if the balance is sufficient to cover the bid amount
	if balanceToCheck != nil && balanceToCheck.Cmp(bidAmount) >= 0 {
		newBalance = new(big.Int).Sub(balanceToCheck, bidAmount)
	} else {
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
	}

	if err := dm.store.SetBalanceForBlock(address, newBalance, blockNumber); err != nil {
		dm.logger.Error("setting balance", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to set balance: %v", err)
	}

	return newBalance, nil
}

func (dm *DepositManager) RefundDeposit(address common.Address, deductedAmount *big.Int, blockNumber int64) error {
	return dm.store.RefundBalanceForBlock(address, deductedAmount, blockNumber)
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
