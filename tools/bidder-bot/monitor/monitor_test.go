package monitor

import (
	"context"
	"errors"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/util"
)

type mockL1Client struct {
	receipts       map[common.Hash]*types.Receipt
	receiptDelay   time.Duration
	receiptRetries int
	callCount      int
}

func (m *mockL1Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.callCount++

	if m.receiptDelay > 0 {
		time.Sleep(m.receiptDelay)
	}

	if m.receiptRetries > 0 && m.callCount <= m.receiptRetries {
		return nil, errors.New("transaction not yet mined")
	}

	receipt, exists := m.receipts[txHash]
	if !exists {
		return nil, errors.New("receipt not found")
	}
	return receipt, nil
}

func TestMonitorTxLanding(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(os.Stdout)

	targetBlockNumber := uint64(12345)
	txHash := common.HexToHash("0x123")

	tests := []struct {
		name                      string
		receipt                   *types.Receipt
		receiptDelay              time.Duration
		receiptRetries            int
		contextTimeout            time.Duration
		checkInterval             time.Duration
		expectLandedInTargetBlock bool
	}{
		{
			name: "transaction lands in target block",
			receipt: &types.Receipt{
				BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
			},
			contextTimeout:            250 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: true,
		},
		{
			name: "transaction lands in different block",
			receipt: &types.Receipt{
				BlockNumber: new(big.Int).SetUint64(targetBlockNumber + 1),
			},
			contextTimeout:            250 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: false,
		},
		{
			name:                      "transaction does not land (no receipt)",
			receipt:                   nil,
			contextTimeout:            50 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: false,
		},
		{
			name: "transaction does not land (delayed past timeout)",
			receipt: &types.Receipt{
				BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
			},
			receiptRetries:            2,
			receiptDelay:              100 * time.Millisecond,
			contextTimeout:            50 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: false,
		},
		{
			name: "transaction lands after retries",
			receipt: &types.Receipt{
				BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
			},
			receiptRetries:            2,
			contextTimeout:            250 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: true,
		},
		{
			name: "transaction lands with delay",
			receipt: &types.Receipt{
				BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
			},
			receiptDelay:              10 * time.Millisecond,
			contextTimeout:            250 * time.Millisecond,
			checkInterval:             10 * time.Millisecond,
			expectLandedInTargetBlock: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l1Client := &mockL1Client{
				receipts:       make(map[common.Hash]*types.Receipt),
				receiptDelay:   tc.receiptDelay,
				receiptRetries: tc.receiptRetries,
				callCount:      0,
			}
			if tc.receipt != nil {
				l1Client.receipts[txHash] = tc.receipt
			}

			m := &Monitor{
				logger:                   logger,
				l1Client:                 l1Client,
				monitorTxLandingTimeout:  tc.contextTimeout,
				monitorTxLandingInterval: tc.checkInterval,
			}

			sentBid := &SentBid{
				TxHash:            txHash,
				TargetBlockNumber: targetBlockNumber,
			}

			t.Logf("Test case: %s", tc.name)
			t.Logf("Context timeout: %v, Receipt delay: %v, Retries: %d",
				tc.contextTimeout, tc.receiptDelay, tc.receiptRetries)

			result := m.monitorTxLanding(context.Background(), sentBid)

			if result != tc.expectLandedInTargetBlock {
				t.Errorf("monitorTxLanding() for %s = %v, want %v", tc.name, result, tc.expectLandedInTargetBlock)
				t.Logf("Call count: %d", l1Client.callCount)
			}
		})
	}
}

func TestStartMonitor(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(os.Stdout)

	l1Client := &mockL1Client{
		receipts: make(map[common.Hash]*types.Receipt),
	}

	bidChan := make(chan *SentBid, 5)

	monitorTxLandingTimeout := 15 * time.Minute
	monitorTxLandingInterval := 30 * time.Second

	m := NewMonitor(
		logger,
		l1Client,
		bidChan,
		monitorTxLandingTimeout,
		monitorTxLandingInterval,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := m.Start(ctx)

	txHash := common.HexToHash("0x123")
	targetBlockNumber := uint64(12345)

	l1Client.receipts[txHash] = &types.Receipt{
		BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
	}

	sentBid := &SentBid{
		TxHash:            txHash,
		TargetBlockNumber: targetBlockNumber,
	}

	bidChan <- sentBid

	time.Sleep(50 * time.Millisecond)

	cancel()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Monitor failed to stop")
	}
}
