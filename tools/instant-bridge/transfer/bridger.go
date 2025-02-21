package transfer

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/bridge/standard/pkg/transfer"
	"github.com/primev/mev-commit/x/keysigner"
)

type AccountSyncer interface {
	Subscribe(ctx context.Context, threshold *big.Int) <-chan struct{}
}

type BridgeConfig struct {
	Signer                 keysigner.KeySigner
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
	L1RPCUrl               string
	SettlementRPCUrl       string
}

type Bridger struct {
	syncer    AccountSyncer
	threshold *big.Int
	topup     *big.Int
	config    BridgeConfig
	logger    *slog.Logger
}

func NewBridger(
	logger *slog.Logger,
	syncer AccountSyncer,
	config BridgeConfig,
	threshold *big.Int,
	topup *big.Int,
) *Bridger {
	return &Bridger{
		syncer:    syncer,
		threshold: threshold,
		topup:     topup,
		config:    config,
		logger:    logger,
	}
}

func (b *Bridger) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			thresholdCrossed := b.syncer.Subscribe(ctx, b.threshold)
			select {
			case <-ctx.Done():
				return
			case <-thresholdCrossed:
				// Bridge funds
				tx, err := transfer.NewTransferToSettlement(
					b.topup,
					b.config.Signer.GetAddress(),
					b.config.Signer,
					b.config.SettlementRPCUrl,
					b.config.L1RPCUrl,
					b.config.L1ContractAddr,
					b.config.SettlementContractAddr,
				)
				if err != nil {
					b.logger.Error("failed to create transfer", "error", err)
					continue
				}
				statusC := tx.Do(ctx)
				for status := range statusC {
					if status.Error != nil {
						b.logger.Error("transfer failed", "error", status.Error)
						continue
					}
					b.logger.Info("transfer progress", "message", status.Message)
				}
			}
		}
	}()
	return done
}
