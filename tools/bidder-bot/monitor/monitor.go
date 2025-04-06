package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type AcceptedBid struct {
	TxHash            common.Hash
	TargetBlockNumber uint64
}

type Monitor struct {
	logger                   *slog.Logger
	l1Client                 L1Client
	monitorTxLandingTimeout  time.Duration
	monitorTxLandingInterval time.Duration
	acceptedBidChan          <-chan *AcceptedBid
}

type L1Client interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func NewMonitor(
	logger *slog.Logger,
	l1Client L1Client,
	acceptedBidChan <-chan *AcceptedBid,
	monitorTxLandingTimeout time.Duration,
	monitorTxLandingInterval time.Duration,
) *Monitor {
	return &Monitor{
		logger:                   logger.With("component", "bid_monitor"),
		l1Client:                 l1Client,
		acceptedBidChan:          acceptedBidChan,
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
			case acceptedBid := <-m.acceptedBidChan:
				m.logger.Info("monitoring accepted bid", "tx_hash", acceptedBid.TxHash.Hex())
				go m.monitorAcceptedBid(ctx, acceptedBid)
			}
		}
	}()
	return done
}

func (m *Monitor) monitorAcceptedBid(ctx context.Context, acceptedBid *AcceptedBid) {
	landedInTargetBlock := m.monitorTxLanding(ctx, acceptedBid)
	if !landedInTargetBlock {
		m.logger.Error("transaction did not land in target block", "tx_hash", acceptedBid.TxHash.Hex())
		return
	}
	m.logger.Info("accepted bid landed in target block",
		"tx_hash", acceptedBid.TxHash.Hex(),
		"target_block_number", acceptedBid.TargetBlockNumber)
}

func (m *Monitor) monitorTxLanding(ctx context.Context, acceptedBid *AcceptedBid) bool {
	txLandingCtx, cancel := context.WithTimeout(ctx, m.monitorTxLandingTimeout)
	defer cancel()
	ticker := time.NewTicker(m.monitorTxLandingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-txLandingCtx.Done():
			m.logger.Warn("tx landing monitoring timeout", "tx_hash", acceptedBid.TxHash.Hex())
			return false
		case <-ticker.C:
			receipt, err := m.l1Client.TransactionReceipt(txLandingCtx, acceptedBid.TxHash)
			if err == nil && receipt != nil {
				actualBlock := receipt.BlockNumber.Uint64()
				if actualBlock == acceptedBid.TargetBlockNumber {
					m.logger.Info("transaction landed in the target block",
						"tx_hash", acceptedBid.TxHash.Hex(),
						"block", actualBlock)
					return true
				}
				m.logger.Warn("transaction landed in non-target block",
					"tx_hash", acceptedBid.TxHash.Hex(),
					"actual_block", actualBlock,
					"target_block", acceptedBid.TargetBlockNumber)
				return false
			}
		}
	}
}
