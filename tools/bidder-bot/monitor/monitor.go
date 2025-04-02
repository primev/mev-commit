package monitor

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
)

type SentBid struct {
	TxHash            common.Hash
	TargetBlockNumber uint64
	BidStream         BidStream
}

type BidStream interface {
	Recv() (*bidderapiv1.Commitment, error)
}

type Monitor struct {
	logger                   *slog.Logger
	topologyClient           debugapiv1.DebugServiceClient
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
	topologyClient debugapiv1.DebugServiceClient,
	l1Client L1Client,
	sentBidChan <-chan *SentBid,
	monitorTxLandingTimeout time.Duration,
	monitorTxLandingInterval time.Duration,
) *Monitor {
	return &Monitor{
		logger:                   logger.With("component", "bid_monitor"),
		topologyClient:           topologyClient,
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
	expectedCommitmentsReceived := m.monitorCommitments(ctx, sentBid)
	if !expectedCommitmentsReceived {
		m.logger.Error("expected commitments not received", "tx_hash", sentBid.TxHash.Hex())
		return
	}

	landedInTargetBlock := m.monitorTxLanding(ctx, sentBid)
	if !landedInTargetBlock {
		m.logger.Error("transaction did not land in target block", "tx_hash", sentBid.TxHash.Hex())
		return
	}
	m.logger.Info("sent bid was successful",
		"tx_hash", sentBid.TxHash.Hex(),
		"expected_commitments_received", expectedCommitmentsReceived,
		"landed_in_target_block", landedInTargetBlock)
}

func (m *Monitor) monitorCommitments(ctx context.Context, sentBid *SentBid) bool {
	topo, err := m.topologyClient.GetTopology(ctx, &debugapiv1.EmptyMessage{})
	if err != nil {
		m.logger.Error("failed to get topology", "tx_hash", sentBid.TxHash.Hex(), "error", err)
		return false
	}

	providers := topo.Topology.Fields["connected_providers"].GetListValue()
	if providers == nil || len(providers.Values) == 0 {
		m.logger.Error("no connected providers", "tx_hash", sentBid.TxHash.Hex())
		return false
	}

	expectedCommitments := len(providers.Values)
	commitments := make([]*bidderapiv1.Commitment, 0)
	for {
		select {
		case <-ctx.Done():
			if len(commitments) == expectedCommitments {
				m.logger.Info("all commitments received", "tx_hash", sentBid.TxHash.Hex())
				return true
			}
			m.logger.Warn("commitment timeout",
				"tx_hash", sentBid.TxHash.Hex(),
				"received", commitments,
				"expected", expectedCommitments)
			return false
		default:
		}

		// 12 second timeout is already set for bid stream in bidder.go
		msg, err := sentBid.BidStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			m.logger.Error("failed to receive commitment", "tx_hash", sentBid.TxHash.Hex(), "error", err)
			return false
		}

		commitments = append(commitments, msg)
		m.logger.Debug("received commitment",
			"tx_hash", sentBid.TxHash.Hex(),
			"count", len(commitments),
			"expected", expectedCommitments)

		if len(commitments) == expectedCommitments {
			m.logger.Info("all commitments received", "tx_hash", sentBid.TxHash.Hex())
			return true
		}
	}

	return len(commitments) == expectedCommitments
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
