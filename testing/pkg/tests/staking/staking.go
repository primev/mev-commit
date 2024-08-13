package staking

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"

	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	// Get all providers
	providers := cluster.Providers()
	logger := cluster.Logger().With("test", "staking")

	registrations := make(chan *providerregistry.ProviderregistryProviderRegistered)
	sub, err := cluster.Events().Subscribe(
		events.NewEventHandler(
			"ProviderRegistered",
			func(p *providerregistry.ProviderregistryProviderRegistered) {
				registrations <- p
			},
		),
	)
	if err != nil {
		logger.Error("failed to subscribe to provider registered events", "error", err)
		return fmt.Errorf("failed to subscribe to provider registered events: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
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

		blsPubkeyBytes := make([]byte, 48)
		_, err = rand.Read(blsPubkeyBytes)
		if err != nil {
			l.Error("failed to generate mock BLS public key", "err", err)
			return fmt.Errorf("failed to generate mock BLS public key: %w", err)
		}

		// Register a provider
		resp, err := providerAPI.RegisterStake(ctx, &providerapiv1.StakeRequest{
			Amount:       stakeAmount.String(),
			BlsPublicKey: hex.EncodeToString(blsPubkeyBytes),
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
