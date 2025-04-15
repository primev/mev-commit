package notifier_test

import (
	"context"
	"math/big"
	"os"
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
		notifySecondsAhead        time.Duration
		receiveSignalWithin       time.Duration
		expectToReceive           bool
	}{
		{
			name:                      "block time in past, should receive immediately",
			currentBlockTimeNowOffset: -15 * time.Second,
			notifySecondsAhead:        10 * time.Second,
			receiveSignalWithin:       recvdBufferTime,
			expectToReceive:           true,
		},
		{
			name:                      "block time in past, don't wait long enough",
			currentBlockTimeNowOffset: -5 * time.Second,
			notifySecondsAhead:        2 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 4*time.Second, // -5 + 12 - 2 = 5
			expectToReceive:           false,
		},
		{
			name:                      "block time in past, wait long enough",
			currentBlockTimeNowOffset: -5 * time.Second,
			notifySecondsAhead:        2 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 5*time.Second,
			expectToReceive:           true,
		},
		{
			name:                      "block time now, 12 sec notif ahead",
			currentBlockTimeNowOffset: 0 * time.Second,
			notifySecondsAhead:        12 * time.Second,
			receiveSignalWithin:       recvdBufferTime,
			expectToReceive:           true,
		},
		{
			name:                      "block time now, 10 sec notif ahead",
			currentBlockTimeNowOffset: 0 * time.Second,
			notifySecondsAhead:        10 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 2*time.Second, // 12 second block time - 10 second = 2
			expectToReceive:           true,
		},
		{
			name:                      "block time now, 8 sec notif ahead, don't wait long enough",
			currentBlockTimeNowOffset: 0 * time.Second,
			notifySecondsAhead:        8 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 3*time.Second, // 12 - 8 = 4
			expectToReceive:           false,
		},
		{
			name:                      "block time now, 8 sec notif ahead, wait long enough",
			currentBlockTimeNowOffset: 0 * time.Second,
			notifySecondsAhead:        8 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 4*time.Second, // 12 - 8 = 4
			expectToReceive:           true,
		},
		{
			name:                      "block time in future, 11 sec notif ahead, don't wait long enough",
			currentBlockTimeNowOffset: 2 * time.Second,
			notifySecondsAhead:        11 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 2*time.Second, // 12 + 2 - 11 = 3
			expectToReceive:           false,
		},
		{
			name:                      "block time in future, 11 sec notif ahead, wait long enough",
			currentBlockTimeNowOffset: 2 * time.Second,
			notifySecondsAhead:        11 * time.Second,
			receiveSignalWithin:       recvdBufferTime + 3*time.Second,
			expectToReceive:           true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			targetBlockNumChan := make(chan uint64, 1)
			notifier := notifier.NewFullNotifier(
				logger,
				nil,
				test.notifySecondsAhead,
				targetBlockNumChan,
			)
			header := &types.Header{
				Number: big.NewInt(int64(71)),
				Time:   uint64(time.Now().Add(test.currentBlockTimeNowOffset).Unix()),
			}
			notifier.HandleHeader(context.Background(), header)

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

func TestDrain(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)

	targetBlockNumChan := make(chan uint64, 1)
	notifier := notifier.NewFullNotifier(
		logger,
		nil,
		10*time.Second,
		targetBlockNumChan,
	)
	notifier.HandleHeader(context.Background(), &types.Header{Number: big.NewInt(5)})
	targetBlockNum := <-targetBlockNumChan
	if targetBlockNum != 6 {
		t.Fatalf("expected target block number %d, got %d", 6, targetBlockNum)
	}
	notifier.HandleHeader(context.Background(), &types.Header{Number: big.NewInt(15)})
	targetBlockNum = <-targetBlockNumChan
	if targetBlockNum != 16 {
		t.Fatalf("expected target block number %d, got %d", 16, targetBlockNum)
	}

	pastTime := time.Now().Add(-100 * time.Second) // To ensure sendTargetBlockNotification is called immediately

	notifier.HandleHeader(context.Background(), &types.Header{Number: big.NewInt(25), Time: uint64(pastTime.Unix())})
	// draining starts here
	notifier.HandleHeader(context.Background(), &types.Header{Number: big.NewInt(35), Time: uint64(pastTime.Unix())})
	notifier.HandleHeader(context.Background(), &types.Header{Number: big.NewInt(45), Time: uint64(pastTime.Unix())})
	targetBlockNum = <-targetBlockNumChan
	if targetBlockNum != 46 {
		t.Fatalf("expected target block number %d, got %d", 46, targetBlockNum)
	}

	select {
	case <-targetBlockNumChan:
		t.Fatal("expected no more in channel")
	default:
	}
}
