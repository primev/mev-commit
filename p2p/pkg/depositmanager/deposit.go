package depositmanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BidderRegistryContract interface {
	GetDepositConsideringWithdrawalRequest(opts *bind.CallOpts, bidder common.Address, provider common.Address) (*big.Int, error)
}

type Store interface {
	GetBalance(bidder common.Address, provider common.Address) (*big.Int, error)
	SetBalance(bidder common.Address, provider common.Address, balance *big.Int) error
	DeleteBalance(bidder common.Address, provider common.Address) error
	RefundBalanceIfExists(bidder common.Address, provider common.Address, amount *big.Int) error
}

type DepositManager struct {
	store            Store
	evtMgr           events.EventManager
	bidderRegistry   BidderRegistryContract
	deposits         chan *bidderregistry.BidderregistryBidderDeposited
	withdrawRequests chan *bidderregistry.BidderregistryWithdrawalRequested
	withdrawals      chan *bidderregistry.BidderregistryBidderWithdrawal
	logger           *slog.Logger
}

func NewDepositManager(
	store Store,
	evtMgr events.EventManager,
	bidderRegistry BidderRegistryContract,
	logger *slog.Logger,
) *DepositManager {
	return &DepositManager{
		store:            store,
		bidderRegistry:   bidderRegistry,
		deposits:         make(chan *bidderregistry.BidderregistryBidderDeposited),
		withdrawRequests: make(chan *bidderregistry.BidderregistryWithdrawalRequested),
		withdrawals:      make(chan *bidderregistry.BidderregistryBidderWithdrawal),
		evtMgr:           evtMgr,
		logger:           logger,
	}
}

func (dm *DepositManager) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	ev1 := events.NewEventHandler(
		"BidderDeposited",
		func(bidderDeposit *bidderregistry.BidderregistryBidderDeposited) {
			select {
			case <-egCtx.Done():
				dm.logger.Info("bidder deposited context done")
			case dm.deposits <- bidderDeposit:
			}
		},
	)

	ev2 := events.NewEventHandler(
		"WithdrawalRequested",
		func(withdrawalRequested *bidderregistry.BidderregistryWithdrawalRequested) {
			select {
			case <-egCtx.Done():
				dm.logger.Info("withdrawal requested context done")
			case dm.withdrawRequests <- withdrawalRequested:
			}
		},
	)

	ev3 := events.NewEventHandler(
		"BidderWithdrawal",
		func(bidderWithdrawal *bidderregistry.BidderregistryBidderWithdrawal) {
			select {
			case <-egCtx.Done():
				dm.logger.Info("bidder withdrawal context done")
			case dm.withdrawals <- bidderWithdrawal:
			}
		},
	)

	sub, err := dm.evtMgr.Subscribe(ev1, ev2, ev3)
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
				dm.logger.Info("deposit manager context done")
				return nil

			case deposit := <-dm.deposits:
				currentBalance, err := dm.store.GetBalance(deposit.Bidder, deposit.Provider)
				if err != nil {
					dm.logger.Error("getting balance", "error", err)
					return err
				}
				if currentBalance == nil {
					dm.logger.Debug("balance not found in store, using default from contract",
						"bidder", deposit.Bidder,
						"provider", deposit.Provider,
					)
					blockBeforeDeposit := new(big.Int).SetUint64(deposit.Raw.BlockNumber - 1)
					currentBalance, err = dm.getDefaultBalance(egCtx, deposit.Bidder, deposit.Provider, blockBeforeDeposit)
					if err != nil {
						dm.logger.Error("getting default balance", "error", err)
						return err
					}
					if currentBalance == nil {
						dm.logger.Error("No balance found in contract. Assuming zero")
						currentBalance = big.NewInt(0)
					}
				}
				newBalance := new(big.Int).Add(currentBalance, deposit.DepositedAmount)
				if err := dm.store.SetBalance(deposit.Bidder, deposit.Provider, newBalance); err != nil {
					dm.logger.Error("setting balance", "error", err)
					return err
				}
				dm.logger.Info("set balance from bidder deposit event",
					"bidder", deposit.Bidder,
					"provider", deposit.Provider,
					"new balance", newBalance,
				)

			case withdrawalRequest := <-dm.withdrawRequests:
				if err := dm.store.DeleteBalance(withdrawalRequest.Bidder, withdrawalRequest.Provider); err != nil {
					dm.logger.Error("deleting balance", "error", err)
					return err
				}
				dm.logger.Info("deleted balance from withdrawal request event",
					"bidder", withdrawalRequest.Bidder,
					"provider", withdrawalRequest.Provider,
				)

			case withdrawal := <-dm.withdrawals:
				if err := dm.store.DeleteBalance(withdrawal.Bidder, withdrawal.Provider); err != nil {
					dm.logger.Error("deleting balance", "error", err)
					return err
				}
				dm.logger.Info("deleted balance from withdrawal event",
					"bidder", withdrawal.Bidder,
					"provider", withdrawal.Provider,
				)
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
	bidderAddr common.Address,
	providerAddr common.Address,
	bidAmountStr string,
) (func() error, error) {
	bidAmount, ok := new(big.Int).SetString(bidAmountStr, 10)
	if !ok {
		dm.logger.Error("parsing bid amount", "amount", bidAmountStr)
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse bid amount")
	}

	balance, err := dm.store.GetBalance(bidderAddr, providerAddr)
	if err != nil {
		dm.logger.Error("getting balance", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	if balance != nil {
		newBalance := new(big.Int).Sub(balance, bidAmount)
		if newBalance.Cmp(big.NewInt(0)) < 0 {
			dm.logger.Error("insufficient balance", "balance", balance.Uint64(), "bidAmount", bidAmount.Uint64())
			return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
		}

		if err := dm.store.SetBalance(bidderAddr, providerAddr, newBalance); err != nil {
			dm.logger.Error("setting balance", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to set balance: %v", err)
		}
		return func() error {
			return dm.store.RefundBalanceIfExists(bidderAddr, providerAddr, bidAmount)
		}, nil
	}
	dm.logger.Info("balance not found in store, defaulting to contract call",
		"bidder", bidderAddr.Hex(),
		"provider", providerAddr.Hex(),
	)

	defaultBalance, err := dm.getDefaultBalance(ctx, bidderAddr, providerAddr, nil) // nil for latest block
	if err != nil {
		return nil, err
	}

	if defaultBalance == nil {
		dm.logger.Error("bidder balance not found", "bidder", bidderAddr.Hex(), "provider", providerAddr.Hex())
		return nil, status.Errorf(codes.FailedPrecondition,
			"balance not found for bidder %s and provider %s", bidderAddr.Hex(), providerAddr.Hex())
	}

	if defaultBalance.Cmp(bidAmount) < 0 {
		dm.logger.Error("insufficient balance", "balance", defaultBalance, "bidAmount", bidAmount)
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
	}

	newBalance := new(big.Int).Sub(defaultBalance, bidAmount)
	if err := dm.store.SetBalance(bidderAddr, providerAddr, newBalance); err != nil {
		dm.logger.Error("setting balance for block", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to set balance for block: %v", err)
	}

	return func() error {
		return dm.store.RefundBalanceIfExists(bidderAddr, providerAddr, bidAmount)
	}, nil
}

// fallback to contract if balance not found in store
func (dm *DepositManager) getDefaultBalance(
	ctx context.Context,
	bidderAddr common.Address,
	providerAddr common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {

	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}

	balance, err := dm.bidderRegistry.GetDepositConsideringWithdrawalRequest(callOpts, bidderAddr, providerAddr)
	if err != nil {
		dm.logger.Error("getting deposit from contract", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get deposit: %v", err)
	}

	if balance.Cmp(big.NewInt(0)) > 0 {
		if err := dm.store.SetBalance(bidderAddr, providerAddr, balance); err != nil {
			dm.logger.Error("setting balance", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to set balance: %v", err)
		}
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil, nil
	}

	return balance, nil
}
