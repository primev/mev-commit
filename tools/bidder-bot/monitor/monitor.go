package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
)

type SentBid struct {
	TxHash            common.Hash
	TargetBlockNumber uint64
}

type BidStream interface {
	Recv() (*bidderapiv1.Commitment, error)
}

type Monitor struct {
	logger                   *slog.Logger
	l1Client                 L1Client
	monitorTxLandingTimeout  time.Duration
	monitorTxLandingInterval time.Duration
	sentBidChan              <-chan *SentBid
}

type L1Client interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func NewMonitor(
	logger *slog.Logger,
	l1Client L1Client,
	sentBidChan <-chan *SentBid,
	monitorTxLandingTimeout time.Duration,
	monitorTxLandingInterval time.Duration,
) *Monitor {
	return &Monitor{
		logger:                   logger.With("component", "bid_monitor"),
		l1Client:                 l1Client,
		sentBidChan:              sentBidChan,
		monitorTxLandingTimeout:  monitorTxLandingTimeout,
		monitorTxLandingInterval: monitorTxLandingInterval,
	}
}

func (m *Monitor) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				m.logger.Info("monitor context done")
				return
			case sentBid := <-m.sentBidChan:
				m.logger.Info("monitoring sent bid", "tx_hash", sentBid.TxHash.Hex())
				go m.monitorSentBid(ctx, sentBid)
			}
		}
	}()
	return done
}

func (m *Monitor) monitorSentBid(ctx context.Context, sentBid *SentBid) {
	landedInTargetBlock := m.monitorTxLanding(ctx, sentBid)
	if !landedInTargetBlock {
		m.logger.Error("transaction did not land in target block", "tx_hash", sentBid.TxHash.Hex())
		return
	}
	m.logger.Info("sent bid was successful",
		"tx_hash", sentBid.TxHash.Hex(),
		"landed_in_target_block", landedInTargetBlock)
}

func (m *Monitor) monitorTxLanding(ctx context.Context, sentBid *SentBid) bool {
	txLandingCtx, cancel := context.WithTimeout(ctx, m.monitorTxLandingTimeout)
	defer cancel()
	ticker := time.NewTicker(m.monitorTxLandingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-txLandingCtx.Done():
			m.logger.Warn("tx landing monitoring timeout", "tx_hash", sentBid.TxHash.Hex())
			return false
		case <-ticker.C:
			receipt, err := m.l1Client.TransactionReceipt(txLandingCtx, sentBid.TxHash)
			if err == nil && receipt != nil {
				actualBlock := receipt.BlockNumber.Uint64()
				if actualBlock == sentBid.TargetBlockNumber {
					m.logger.Info("transaction landed in the target block",
						"tx_hash", sentBid.TxHash.Hex(),
						"block", actualBlock)
					return true
				}
				m.logger.Warn("transaction landed in non-target block",
					"tx_hash", sentBid.TxHash.Hex(),
					"actual_block", actualBlock,
					"target_block", sentBid.TargetBlockNumber)
				return false
			}
		}
	}
}
