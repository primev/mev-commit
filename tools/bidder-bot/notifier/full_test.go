package notifier_test

import (
	"context"
	"math/big"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
	"github.com/primev/mev-commit/tools/bidder-bot/notifier"
	"github.com/primev/mev-commit/x/util"
)

func TestHandleHeader(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)

	recvdBufferTime := 200 * time.Millisecond

	tests := []struct {
		name                      string
		currentBlockTimeNowOffset time.Duration
		receiveSignalWithin       time.Duration
		expectToReceive           bool
	}{
		{
			name:                      "block time in past, should receive immediately",
			currentBlockTimeNowOffset: -15 * time.Second,
			receiveSignalWithin:       recvdBufferTime,
			expectToReceive:           true,
		},
		{
			name:                      "block time now, should receive immediately",
			currentBlockTimeNowOffset: 0 * time.Second,
			receiveSignalWithin:       recvdBufferTime,
			expectToReceive:           true,
		},
		{
			name:                      "block time in future, should receive immediately",
			currentBlockTimeNowOffset: 7 * time.Second,
			receiveSignalWithin:       recvdBufferTime,
			expectToReceive:           true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			targetBlockChan := make(chan bidder.TargetBlock, 1)
			notifier := notifier.NewFullNotifier(
				logger,
				nil,
				targetBlockChan,
				1, // Setting block interval to 1 ensures no headers are skipped
			)
			header := &types.Header{
				Number: big.NewInt(int64(71)),
				Time:   uint64(time.Now().Add(test.currentBlockTimeNowOffset).Unix()),
			}
			err := notifier.HandleHeader(context.Background(), header)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			select {
			case <-time.After(test.receiveSignalWithin):
				if test.expectToReceive {
					t.Fatal("notification not sent in time")
				}
			case targetBlock := <-targetBlockChan:
				if !test.expectToReceive {
					t.Fatal("notification sent when not expected")
				}
				if targetBlock.Num != 72 {
					t.Fatalf("expected target block number %d, got %d", 72, targetBlock.Num)
				}
			}
		})
	}
}

func TestChannelTimeout(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)
	ctx := context.Background()

	targetBlockChan := make(chan bidder.TargetBlock, 1)
	notifier := notifier.NewFullNotifier(
		logger,
		nil,
		targetBlockChan,
		1, // Setting block interval to 1 ensures no headers are skipped
	)

	err := notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(5)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case targetBlock := <-targetBlockChan:
		if targetBlock.Num != 6 {
			t.Fatalf("expected target block number %d, got %d", 6, targetBlock.Num)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected to receive block number")
	}

	err = notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(10)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(15)})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to send target block 16") {
		t.Fatalf("unexpected error, got %v", err)
	}

	select {
	case targetBlock := <-targetBlockChan:
		if targetBlock.Num != 11 {
			t.Fatalf("expected target block number %d, got %d", 11, targetBlock.Num)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected to receive block number")
	}

	err = notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(15)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case targetBlock := <-targetBlockChan:
		if targetBlock.Num != 16 {
			t.Fatalf("expected target block number %d, got %d", 16, targetBlock.Num)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected to receive block number")
	}

	select {
	case <-targetBlockChan:
		t.Fatal("expected no more blocks in channel")
	default:
	}
}

func TestBlockInterval(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)

	tests := []struct {
		name                 string
		blockInterval        uint64
		headers              []types.Header
		expectedTargetBlocks []uint64
	}{
		{
			name:          "block interval 1, should receive all headers",
			blockInterval: 1,
			headers: []types.Header{
				{Number: big.NewInt(1)},
				{Number: big.NewInt(2)},
				{Number: big.NewInt(777)},
				{Number: big.NewInt(778)},
			},
			expectedTargetBlocks: []uint64{2, 3, 778, 779}, // target = header + 1
		},

		{
			name:          "block interval 2, should receive 2nd, 4th, 444th headers",
			blockInterval: 2,
			headers: []types.Header{
				{Number: big.NewInt(1)},
				{Number: big.NewInt(2)},
				{Number: big.NewInt(3)},
				{Number: big.NewInt(4)},
				{Number: big.NewInt(444)},
			},
			expectedTargetBlocks: []uint64{3, 5, 445}, // target = header + 1
		},

		{
			name:          "block interval 7, should receive 7th, 14th, 700th headers",
			blockInterval: 7,
			headers: []types.Header{
				{Number: big.NewInt(1)},
				{Number: big.NewInt(2)},
				{Number: big.NewInt(3)},
				{Number: big.NewInt(4)},
				{Number: big.NewInt(7)},
				{Number: big.NewInt(14)},
				{Number: big.NewInt(700)},
				{Number: big.NewInt(701)},
			},
			expectedTargetBlocks: []uint64{8, 15, 701}, // target = header + 1
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			targetBlockChan := make(chan bidder.TargetBlock, len(test.headers))
			notifier := notifier.NewFullNotifier(
				logger,
				nil,
				targetBlockChan,
				test.blockInterval,
			)
			for _, header := range test.headers {
				err := notifier.HandleHeader(context.Background(), &header)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
			receivedTargetBlocks := make([]uint64, 0, len(test.expectedTargetBlocks))
			for i := 0; i < len(test.expectedTargetBlocks); i++ {
				select {
				case targetBlock := <-targetBlockChan:
					receivedTargetBlocks = append(receivedTargetBlocks, targetBlock.Num)
				case <-time.After(1 * time.Second):
					t.Fatal("timeout waiting for target block")
				}
			}
			if !slices.Equal(receivedTargetBlocks, test.expectedTargetBlocks) {
				t.Fatalf("expected to receive blocks %v, but got %v", test.expectedTargetBlocks, receivedTargetBlocks)
			}
		})
	}
}
