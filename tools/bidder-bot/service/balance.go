package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/primev/mev-commit/x/contracts/ethwrapper"
	"github.com/primev/mev-commit/x/keysigner"
)

type BalanceChecker struct {
	logger              *slog.Logger
	signer              keysigner.KeySigner
	l1RPCClient         *ethwrapper.Client
	settlementRPCClient *ethwrapper.Client
}

func NewBalanceChecker(
	logger *slog.Logger,
	signer keysigner.KeySigner,
	l1RPCClient *ethwrapper.Client,
	settlementRPCClient *ethwrapper.Client,
) *BalanceChecker {
	return &BalanceChecker{
		logger:              logger,
		signer:              signer,
		l1RPCClient:         l1RPCClient,
		settlementRPCClient: settlementRPCClient,
	}
}

func (b *BalanceChecker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := b.CheckBalances(ctx)
				if err != nil {
					b.logger.Error("balance check failed", "error", err)
				}
			}
		}
	}()
	return done
}

func (b *BalanceChecker) CheckBalances(ctx context.Context) error {
	balanceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	l1Balance, err := b.l1RPCClient.RawClient().BalanceAt(balanceCtx, b.signer.GetAddress(), nil)
	if err != nil {
		return err
	}

	settlementBalance, err := b.settlementRPCClient.RawClient().BalanceAt(balanceCtx, b.signer.GetAddress(), nil)
	if err != nil {
		return err
	}

	pointZeroFiveEth := big.NewInt(50000000000000000)
	if l1Balance.Cmp(pointZeroFiveEth) < 0 {
		return fmt.Errorf("keystore account has less than 0.05 eth on L1")
	}

	pointFiveEth := big.NewInt(500000000000000000)
	if settlementBalance.Cmp(pointFiveEth) < 0 {
		return fmt.Errorf("keystore account has less than 0.5 eth on mev-commit chain")
	}

	return nil
}
