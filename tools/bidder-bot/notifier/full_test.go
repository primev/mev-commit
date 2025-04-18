package notifier_test

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
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
			targetBlockNumChan := make(chan uint64, 1)
			notifier := notifier.NewFullNotifier(
				logger,
				nil,
				targetBlockNumChan,
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
			case targetBlockNum := <-targetBlockNumChan:
				if !test.expectToReceive {
					t.Fatal("notification sent when not expected")
				}
				if targetBlockNum != 72 {
					t.Fatalf("expected target block number %d, got %d", 72, targetBlockNum)
				}
			}
		})
	}
}

func TestChannelTimeout(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)
	ctx := context.Background()

	targetBlockNumChan := make(chan uint64, 1)
	notifier := notifier.NewFullNotifier(
		logger,
		nil,
		targetBlockNumChan,
	)

	err := notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(5)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case blockNum := <-targetBlockNumChan:
		if blockNum != 6 {
			t.Fatalf("expected target block number %d, got %d", 6, blockNum)
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
	if !strings.Contains(err.Error(), "failed to send target block number 16") {
		t.Fatalf("unexpected error, got %v", err)
	}

	select {
	case blockNum := <-targetBlockNumChan:
		if blockNum != 11 {
			t.Fatalf("expected target block number %d, got %d", 11, blockNum)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected to receive block number")
	}

	err = notifier.HandleHeader(ctx, &types.Header{Number: big.NewInt(15)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case blockNum := <-targetBlockNumChan:
		if blockNum != 16 {
			t.Fatalf("expected target block number %d, got %d", 16, blockNum)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected to receive block number")
	}

	select {
	case <-targetBlockNumChan:
		t.Fatal("expected no more blocks in channel")
	default:
	}
}
