package stakemanager

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

type ProviderRegistryContract interface {
	GetProviderStake(*bind.CallOpts, common.Address) (*big.Int, error)
	MinStake(*bind.CallOpts) (*big.Int, error)
}

type StakeManager struct {
	owner            common.Address
	evtMgr           events.EventManager
	providerRegistry ProviderRegistryContract
	notifier         notifications.Notifier
	stakeMu          sync.RWMutex
	stakes           *lru.Cache[common.Address, *big.Int]
	unstakeReqs      *lru.Cache[common.Address, struct{}]
	minStake         atomic.Pointer[big.Int]
	logger           *slog.Logger
}

func NewStakeManager(
	logger *slog.Logger,
	owner common.Address,
	evtMgr events.EventManager,
	providerRegistry ProviderRegistryContract,
	notifier notifications.Notifier,
) (*StakeManager, error) {
	minStake, err := providerRegistry.MinStake(&bind.CallOpts{
		From: owner,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get min stake: %w", err)
	}
	stakes, err := lru.New[common.Address, *big.Int](1000)
	if err != nil {
		return nil, fmt.Errorf("failed to create stakes cache: %w", err)
	}
	unstakeReqs, err := lru.New[common.Address, struct{}](1000)
	if err != nil {
		return nil, fmt.Errorf("failed to create unstake requests cache: %w", err)
	}
	sm := &StakeManager{
		providerRegistry: providerRegistry,
		evtMgr:           evtMgr,
		logger:           logger,
		owner:            owner,
		stakes:           stakes,
		unstakeReqs:      unstakeReqs,
		notifier:         notifier,
	}
	sm.minStake.Store(minStake)
	return sm, nil
}

func (sm *StakeManager) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	ch1 := make(chan *providerregistry.ProviderregistryProviderRegistered, 10)
	ev1 := events.NewChannelEventHandler(egCtx, "ProviderRegistered", ch1)

	ch2 := make(chan *providerregistry.ProviderregistryFundsSlashed, 10)
	ev2 := events.NewChannelEventHandler(egCtx, "FundsSlashed", ch2)

	ch3 := make(chan *providerregistry.ProviderregistryFundsDeposited, 10)
	ev3 := events.NewChannelEventHandler(egCtx, "FundsDeposited", ch3)

	ch4 := make(chan *providerregistry.ProviderregistryUnstake, 10)
	ev4 := events.NewChannelEventHandler(egCtx, "Unstake", ch4)

	ch5 := make(chan *providerregistry.ProviderregistryMinStakeUpdated, 10)
	ev5 := events.NewChannelEventHandler(egCtx, "MinStakeUpdated", ch5)

	sub, err := sm.evtMgr.Subscribe(ev1, ev2, ev3, ev4, ev5)
	if err != nil {
		close(doneChan)
		return doneChan
	}

	eg.Go(func() error {
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			sm.logger.Info("event subscription context done")
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("error in event subscription: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				sm.logger.Info("clear balances set balances context done")
				return nil
			case evt := <-ch1:
				// handle ProviderRegistered event
				sm.stakeMu.RLock()
				_, found := sm.stakes.Get(evt.Provider)
				sm.stakeMu.RUnlock()

				if !found {
					sm.stakeMu.Lock()
					_ = sm.stakes.Add(evt.Provider, evt.StakedAmount)
					_ = sm.unstakeReqs.Remove(evt.Provider)
					sm.stakeMu.Unlock()

					sm.notifier.Notify(
						notifications.NewNotification(
							notifications.TopicProviderRegistered,
							map[string]any{
								"provider": evt.Provider.Hex(),
							},
						),
					)
				}
			case evt := <-ch2:
				// handle FundsSlashed event
				sm.stakeMu.RLock()
				stake, found := sm.stakes.Get(evt.Provider)
				sm.stakeMu.RUnlock()

				hasEnoughStake := false

				if !found {
					// if not tracked locally, get the latest value from on-chain
					s, err := sm.providerRegistry.GetProviderStake(&bind.CallOpts{
						From:    sm.owner,
						Context: egCtx,
					}, evt.Provider)
					if err != nil {
						sm.logger.Error("failed to get provider stake", "error", err)
						continue
					}

					sm.stakeMu.Lock()
					_ = sm.stakes.Add(evt.Provider, s)
					sm.stakeMu.Unlock()
					hasEnoughStake = s.Cmp(sm.minStake.Load()) >= 0
				} else {
					// update the local value
					newStake := new(big.Int).Sub(stake, evt.Amount)
					sm.stakeMu.Lock()
					_ = sm.stakes.Add(evt.Provider, newStake)
					sm.stakeMu.Unlock()
					hasEnoughStake = newStake.Cmp(sm.minStake.Load()) >= 0
				}

				sm.notifier.Notify(
					notifications.NewNotification(
						notifications.TopicProviderSlashed,
						map[string]any{
							"provider":       evt.Provider.Hex(),
							"amount":         evt.Amount,
							"hasEnoughStake": hasEnoughStake,
						},
					),
				)
			case evt := <-ch3:
				// handle FundsDeposited event
				sm.stakeMu.RLock()
				stake, found := sm.stakes.Get(evt.Provider)
				sm.stakeMu.RUnlock()

				hasEnoughStake := false

				if !found {
					// if not tracked locally, get the latest value from on-chain
					s, err := sm.providerRegistry.GetProviderStake(&bind.CallOpts{
						From:    sm.owner,
						Context: egCtx,
					}, evt.Provider)
					if err != nil {
						sm.logger.Error("failed to get provider stake", "error", err)
						continue
					}

					sm.stakeMu.Lock()
					_ = sm.stakes.Add(evt.Provider, s)
					sm.stakeMu.Unlock()
					hasEnoughStake = s.Cmp(sm.minStake.Load()) >= 0
				} else {
					// update the local value
					newStake := new(big.Int).Add(stake, evt.Amount)
					sm.stakeMu.Lock()
					_ = sm.stakes.Add(evt.Provider, newStake)
					sm.stakeMu.Unlock()
					hasEnoughStake = newStake.Cmp(sm.minStake.Load()) >= 0
				}

				sm.notifier.Notify(
					notifications.NewNotification(
						notifications.TopicProviderDeposit,
						map[string]any{
							"provider":       evt.Provider.Hex(),
							"amount":         evt.Amount,
							"hasEnoughStake": hasEnoughStake,
						},
					),
				)
			case evt := <-ch4:
				// handle Unstake event
				// even after unstaking, the provider stake is still non-zero till
				// the withdraw request is processed. So, we keep track of the unstake
				// requests in order to avoid querying the on-chain value for the provider
				// stake.
				sm.stakeMu.RLock()
				_, found := sm.unstakeReqs.Get(evt.Provider)
				sm.stakeMu.RUnlock()

				if !found {
					sm.stakeMu.Lock()
					_ = sm.unstakeReqs.Add(evt.Provider, struct{}{})
					_ = sm.stakes.Remove(evt.Provider)
					sm.stakeMu.Unlock()

					sm.notifier.Notify(
						notifications.NewNotification(
							notifications.TopicProviderDeregistered,
							map[string]any{
								"provider": evt.Provider.Hex(),
							},
						),
					)
				}
			case evt := <-ch5:
				// handle MinStakeUpdated event
				sm.minStake.Store(evt.NewMinStake)
			}
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			sm.logger.Error("error in StakeManager", "error", err)
		}
	}()

	return doneChan
}

func (sm *StakeManager) GetStake(ctx context.Context, provider common.Address) (*big.Int, error) {
	sm.stakeMu.RLock()
	_, found := sm.unstakeReqs.Get(provider)
	sm.stakeMu.RUnlock()
	if found {
		return big.NewInt(0), nil
	}

	sm.stakeMu.RLock()
	stake, found := sm.stakes.Get(provider)
	sm.stakeMu.RUnlock()
	if !found {
		stake, err := sm.providerRegistry.GetProviderStake(&bind.CallOpts{
			From:    sm.owner,
			Context: ctx,
		}, provider)
		if err != nil {
			return nil, fmt.Errorf("failed to get provider stake: %w", err)
		}

		sm.stakeMu.Lock()
		_ = sm.stakes.Add(provider, stake)
		sm.stakeMu.Unlock()

		return stake, nil
	}

	return stake, nil
}

func (sm *StakeManager) MinStake() *big.Int {
	return new(big.Int).Set(sm.minStake.Load())
}

func (sm *StakeManager) CheckProviderRegistered(ctx context.Context, provider common.Address) bool {
	stake, err := sm.GetStake(ctx, provider)
	if err != nil {
		sm.logger.Error("failed to get provider stake", "error", err)
		return false
	}

	return stake.Cmp(sm.minStake.Load()) >= 0
}
