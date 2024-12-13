package staking

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"slices"

	"github.com/cloudflare/circl/sign/bls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

type StakingConfig struct {
	RelayURL string
}

func Run(ctx context.Context, cluster orchestrator.Orchestrator, cfg any) error {
	// Get all providers
	providers := cluster.Providers()
	logger := cluster.Logger().With("test", "staking")

	stakingCfg, ok := cfg.(StakingConfig)
	if !ok {
		return fmt.Errorf("invalid staking config type")
	}

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
		// Generate a BLS signature to verify
		message := common.HexToAddress(p.EthAddress()).Bytes()
		hashedMessage := crypto.Keccak256(message)
		ikm := make([]byte, 32)
		privateKey, err := bls.KeyGen[bls.G1](ikm, nil, nil)
		if err != nil {
			l.Error("failed to generate private key", "error", err)
			return fmt.Errorf("failed to generate private key: %w", err)
		}
		publicKey := privateKey.PublicKey()
		signature := bls.Sign(privateKey, hashedMessage)

		// Verify the signature
		if !bls.Verify(publicKey, hashedMessage, signature) {
			l.Error("failed to verify generated BLS signature")
			return fmt.Errorf("failed to verify generated BLS signature")
		}

		pubkeyb, err := publicKey.MarshalBinary()
		if err != nil {
			l.Error("failed to marshal public key", "error", err)
			return fmt.Errorf("failed to marshal public key: %w", err)
		}

		// Register a provider
		resp, err := providerAPI.Stake(ctx, &providerapiv1.StakeRequest{
			Amount:        stakeAmount.String(),
			BlsPublicKeys: []string{hex.EncodeToString(pubkeyb)},
			BlsSignatures: []string{hex.EncodeToString(signature)},
		})

		if err != nil {
			l.Error("failed to register stake", "error", err)
			return fmt.Errorf("failed to register stake: %w", err)
		}

		if resp.Amount != stakeAmount.String() {
			l.Error("invalid stake amount returned", "returned", resp.Amount, "expected", stakeAmount.String())
			return fmt.Errorf("invalid stake amount returned: %s", resp.Amount)
		}

		reqBytes, err := json.Marshal([]string{hex.EncodeToString(pubkeyb)})
		if err != nil {
			l.Error("failed to create bls key upload req", "error", err)
			return fmt.Errorf("failed to marshal bls keys: %w", err)
		}
		relayResp, err := http.Post(
			fmt.Sprintf("%s/register", stakingCfg.RelayURL),
			"application/json",
			bytes.NewReader(reqBytes),
		)
		if err != nil {
			l.Error("failed to post to relay", "error", err)
			return fmt.Errorf("failed to post to relay: %w", err)
		}

		if relayResp.StatusCode != http.StatusOK {
			l.Error("failed to post to relay", "status", relayResp.Status)
			return fmt.Errorf("invalid status code from relay: %d", relayResp.StatusCode)
		}

		getStakeResp, err := providerAPI.GetStake(ctx, &providerapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get stake", "error", err)
			return fmt.Errorf("failed to get stake: %w", err)
		}

		if getStakeResp.Amount != stakeAmount.String() {
			l.Error(
				"invalid stake amount returned on get",
				"returned", getStakeResp.Amount,
				"expected", stakeAmount.String(),
			)
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

func RunAddDeposit(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	// Get all providers
	providers := cluster.Providers()
	logger := cluster.Logger().With("test", "staking_add_deposit")

	deposits := make(chan *providerregistry.ProviderregistryFundsDeposited)
	sub, err := cluster.Events().Subscribe(
		events.NewEventHandler(
			"FundsDeposited",
			func(p *providerregistry.ProviderregistryFundsDeposited) {
				deposits <- p
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
			case deposit := <-deposits:
				logger.Info("provider deposited", "provider", deposit.Provider)
				if !slices.ContainsFunc(providers, func(p orchestrator.Provider) bool {
					return p.EthAddress() == deposit.Provider.String()
				}) {
					logger.Error("provider not found", "provider", deposit.Provider)
					return fmt.Errorf("provider not found: %s", deposit.Provider)
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

		getStakeResp, err := providerAPI.GetStake(ctx, &providerapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get stake", "error", err)
			return fmt.Errorf("failed to get stake: %w", err)
		}

		amount, set := big.NewInt(0).SetString(getStakeResp.Amount, 10)
		if !set {
			l.Error("failed to parse stake", "amount", getStakeResp.Amount)
			return fmt.Errorf("failed to parse stake: %s", getStakeResp.Amount)
		}

		// Register a provider
		resp, err := providerAPI.Stake(ctx, &providerapiv1.StakeRequest{
			Amount: amount.String(),
		})
		if err != nil {
			l.Error("failed to register stake", "error", err)
			return fmt.Errorf("failed to register stake: %w", err)
		}

		if resp.Amount != amount.String() {
			l.Error("invalid stake amount returned", "returned", resp.Amount, "expected", amount.String())
			return fmt.Errorf("invalid stake amount returned: %s", resp.Amount)
		}

		getStakeResp, err = providerAPI.GetStake(ctx, &providerapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get stake", "error", err)
			return fmt.Errorf("failed to get stake: %w", err)
		}

		expStake := big.NewInt(0).Mul(amount, big.NewInt(2))
		if getStakeResp.Amount != expStake.String() {
			l.Error("invalid stake amount returned on get", "returned", getStakeResp.Amount, "expected", expStake.String())
			return fmt.Errorf("invalid stake amount returned: %s", getStakeResp.Amount)
		}
	}

	logger.Info("all providers deposited stakes")

	if err := eg.Wait(); err != nil {
		logger.Error("failed to wait for provider registrations", "error", err)
		return err
	}

	logger.Info("test passed")

	return nil
}
