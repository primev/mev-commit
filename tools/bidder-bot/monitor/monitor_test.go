package monitor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type mockBidStream struct {
	commitments []*bidderapiv1.Commitment
	current     int
	delay       time.Duration
}

func (m *mockBidStream) Recv() (*bidderapiv1.Commitment, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	if m.current >= len(m.commitments) {
		return nil, io.EOF
	}

	commitment := m.commitments[m.current]
	m.current++
	return commitment, nil
}

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

type mockTopologyClient struct {
	providerCount int
	err           error
	responseDelay time.Duration
}

func (m *mockTopologyClient) GetTopology(ctx context.Context, req *debugapiv1.EmptyMessage, opts ...grpc.CallOption) (*debugapiv1.TopologyResponse, error) {
	if m.responseDelay > 0 {
		time.Sleep(m.responseDelay)
	}

	if m.err != nil {
		return nil, m.err
	}

	providersCount := m.providerCount
	if providersCount < 0 {
		providersCount = 0
	}

	providers := make([]*structpb.Value, providersCount)
	for i := 0; i < providersCount; i++ {
		providers[i] = structpb.NewStringValue(fmt.Sprintf("provider%d", i+1))
	}

	listValue := structpb.ListValue{
		Values: providers,
	}

	fields := map[string]*structpb.Value{
		"connected_providers": structpb.NewListValue(&listValue),
	}

	topologyStruct := &structpb.Struct{
		Fields: fields,
	}

	return &debugapiv1.TopologyResponse{
		Topology: topologyStruct,
	}, nil
}

func (m *mockTopologyClient) CancelTransaction(ctx context.Context, req *debugapiv1.CancelTransactionReq, opts ...grpc.CallOption) (*debugapiv1.CancelTransactionResponse, error) {
	return &debugapiv1.CancelTransactionResponse{}, nil
}

func (m *mockTopologyClient) GetPendingTransactions(ctx context.Context, req *debugapiv1.EmptyMessage, opts ...grpc.CallOption) (*debugapiv1.PendingTransactionsResponse, error) {
	return &debugapiv1.PendingTransactionsResponse{}, nil
}

func TestMonitorCommitments(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(os.Stdout)

	tests := []struct {
		name              string
		providerCount     int
		commitments       []*bidderapiv1.Commitment
		topologyErr       error
		topologyDelay     time.Duration
		commitmentDelay   time.Duration
		expectAllReceived bool
		description       string
	}{
		{
			name:              "all commitments received quickly",
			providerCount:     3,
			commitments:       []*bidderapiv1.Commitment{{}, {}, {}},
			expectAllReceived: true,
			description:       "All commitments are received before timeout",
		},
		{
			name:              "all commitments received with delay",
			providerCount:     2,
			commitments:       []*bidderapiv1.Commitment{{}, {}},
			commitmentDelay:   10 * time.Millisecond,
			expectAllReceived: true,
			description:       "All commitments are received but with some network delay",
		},
		{
			name:              "all commitments received with commitment delay",
			providerCount:     3,
			commitments:       []*bidderapiv1.Commitment{{}, {}, {}},
			commitmentDelay:   30 * time.Millisecond,
			expectAllReceived: true,
			description:       "All commitments are received but with some network delay",
		},
		{
			name:              "topology delay but all commitments received",
			providerCount:     2,
			commitments:       []*bidderapiv1.Commitment{{}, {}},
			topologyDelay:     15 * time.Millisecond,
			expectAllReceived: true,
			description:       "There's a delay in getting topology but all commitments are received",
		},
		{
			name:              "topology error",
			topologyErr:       errors.New("failed to fetch topology"),
			expectAllReceived: false,
			description:       "Error when fetching topology information",
		},
		{
			name:              "no commitments available",
			providerCount:     2,
			commitments:       []*bidderapiv1.Commitment{},
			expectAllReceived: false,
			description:       "No commitments are available from the stream",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			topologyClient := &mockTopologyClient{
				providerCount: tc.providerCount,
				err:           tc.topologyErr,
				responseDelay: tc.topologyDelay,
			}

			bidStream := &mockBidStream{
				commitments: tc.commitments,
				delay:       tc.commitmentDelay,
			}

			m := &Monitor{
				logger:         logger,
				topologyClient: topologyClient,
			}

			ctx := context.Background()

			txHash := common.HexToHash("0x123")

			sentBid := &SentBid{
				TxHash:            txHash,
				TargetBlockNumber: 12345,
				BidStream:         bidStream,
			}

			result := m.monitorCommitments(ctx, sentBid)

			if result != tc.expectAllReceived {
				t.Errorf("monitorCommitments() for '%s' = %v, want %v", tc.name, result, tc.expectAllReceived)
				t.Logf("Description: %s", tc.description)
				if tc.providerCount > 0 {
					t.Logf("Expected commitments: %d, Available commitments: %d", tc.providerCount, len(tc.commitments))
				}
			}
		})
	}
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

	topologyClient := &mockTopologyClient{
		providerCount: 2,
	}

	l1Client := &mockL1Client{
		receipts: make(map[common.Hash]*types.Receipt),
	}

	bidChan := make(chan *SentBid, 5)

	monitorTxLandingTimeout := 15 * time.Minute
	monitorTxLandingInterval := 30 * time.Second

	m := NewMonitor(
		logger,
		topologyClient,
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
	bidStream := &mockBidStream{
		commitments: []*bidderapiv1.Commitment{{}, {}},
	}

	l1Client.receipts[txHash] = &types.Receipt{
		BlockNumber: new(big.Int).SetUint64(targetBlockNumber),
	}

	sentBid := &SentBid{
		TxHash:            txHash,
		TargetBlockNumber: targetBlockNumber,
		BidStream:         bidStream,
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
