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
	DeductAndCheckBalanceForBlock(bidder common.Address, defaultAmount, bidAmount *big.Int, blockNumber int64) (*big.Int, error)
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
				if err := dm.store.SetBalance(bidderReg.Bidder, bidderReg.WindowNumber, bidderReg.DepositedAmount); err != nil {
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
	if dm.blocksPerWindow.Load() == 0 {
		blocksPerWindow, err := dm.blockTracker.GetBlocksPerWindow()
		if err != nil {
			dm.logger.Error("getting blocks per window", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get blocks per window: %v", err)
		}
		dm.blocksPerWindow.Store(blocksPerWindow.Uint64())
	}

	bidAmount, ok := new(big.Int).SetString(bidAmountStr, 10)
	if !ok {
		dm.logger.Error("parsing bid amount", "amount", bidAmountStr)
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse bid amount")
	}

	// adding 2 to the current window, bcs oracle is 2 windows behind
	windowToCheck := big.NewInt(dm.currentWindow.Load() + 2)

	balance, err := dm.store.GetBalance(address, windowToCheck)
	if err != nil {
		dm.logger.Error("getting balance", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	if balance == nil {
		dm.logger.Error("bidder balance not found", "address", address.Hex(), "window", windowToCheck)
		return nil, status.Errorf(codes.FailedPrecondition, "balance not found")
	}

	dm.logger.Info("checking bidder deposit",
		"stake", balance.Uint64(),
		"blocksPerWindow", dm.blocksPerWindow.Load(),
		"minStake", dm.minDeposit.Load(),
		"window", windowToCheck.Uint64(),
		"address", address.Hex(),
	)

	blocksPerWindow := new(big.Int).SetUint64(dm.blocksPerWindow.Load())

	// todo: make sense to do division only once, when bidder deposit funds,
	// not everytime, when checking deposit
	effectiveStake := new(big.Int).Div(new(big.Int).Set(balance), blocksPerWindow)

	deductedBalance, err := dm.store.DeductAndCheckBalanceForBlock(address, effectiveStake, bidAmount, blockNumber)
	if err != nil {
		dm.logger.Error("deducting balance", "error", err)
		return nil, status.Errorf(codes.FailedPrecondition, "failed to deduct balance: %v", err)
	}
	return deductedBalance, nil
}

func (dm *DepositManager) RefundDeposit(address common.Address, deductedAmount *big.Int, blockNumber int64) error {
	return dm.store.RefundBalanceForBlock(address, deductedAmount, blockNumber)
}
