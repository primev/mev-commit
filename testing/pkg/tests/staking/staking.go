package staking

import (
	"context"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	// Get all providers
	providers := cluster.Providers()
	providerRegistry := cluster.ProviderRegistry()
	logger := cluster.Logger().With("test", "staking")

	registrations := make(chan *providerregistry.ProviderregistryProviderRegistered)
	sub, err := providerRegistry.WatchProviderRegistered(&bind.WatchOpts{
		Context: ctx,
	}, registrations, nil)
	if err != nil {
		logger.Error("failed to watch provider registered", "error", err)
		return err
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		defer sub.Unsubscribe()

		count := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-sub.Err():
				logger.Error("provider registered subscription failed", "error", err)
				return err
			case registration := <-registrations:
				logger.Info("provider registered", "provider", registration.Provider)
				if !slices.ContainsFunc(providers, func(p orchestrator.Provider) bool {
					return p.EthAddress() == registration.Provider.String()
				}) {
					logger.Error("provider not found", "provider", registration.Provider)
					return fmt.Errorf("provider not found: %s", registration.Provider)
				}
				count++
			}
			if count == len(providers) {
				break
			}
		}
		logger.Info("all providers registered events received")
		return nil
	})

	for _, p := range providers {
		// Get the provider API
		providerAPI := p.ProviderAPI()

		l := logger.With("provider", p.EthAddress())

		minStake, err := providerAPI.GetMinStake(ctx, &providerapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get min stake", "error", err)
			return fmt.Errorf("failed to get min stake: %w", err)
		}

		amount, set := big.NewInt(0).SetString(minStake.Amount, 10)
		if !set {
			l.Error("failed to parse min stake", "amount", minStake.Amount)
			return fmt.Errorf("failed to parse min stake: %s", minStake.Amount)
		}

		stakeAmount := big.NewInt(0).Mul(amount, big.NewInt(10))

		// Register a provider
		resp, err := providerAPI.RegisterStake(ctx, &providerapiv1.StakeRequest{
			Amount: stakeAmount.String(),
		})
		if err != nil {
			l.Error("failed to register stake", "error", err)
			return fmt.Errorf("failed to register stake: %w", err)
		}

		if resp.Amount != stakeAmount.String() {
			l.Error("invalid stake amount returned", "returned", resp.Amount, "expected", stakeAmount.String())
			return fmt.Errorf("invalid stake amount returned: %s", resp.Amount)
		}

		getStakeResp, err := providerAPI.GetStake(ctx, &providerapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get stake", "error", err)
			return fmt.Errorf("failed to get stake: %w", err)
		}

		if getStakeResp.Amount != stakeAmount.String() {
			l.Error("invalid stake amount returned on get", "returned", getStakeResp.Amount, "expected", stakeAmount.String())
			return fmt.Errorf("invalid stake amount returned: %s", getStakeResp.Amount)
		}
	}

	logger.Info("all providers registered stakes")

	if err := eg.Wait(); err != nil {
		logger.Error("failed to wait for provider registrations", "error", err)
		return err
	}

	logger.Info("test passed")

	return nil
}
