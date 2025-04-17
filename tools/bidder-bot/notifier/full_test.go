package notifier_test

import (
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
			err := notifier.HandleHeader(header)
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

func TestDrain(t *testing.T) {
	t.Parallel()
	logger := util.NewTestLogger(os.Stdout)

	targetBlockNumChan := make(chan uint64, 1)
	notifier := notifier.NewFullNotifier(
		logger,
		nil,
		targetBlockNumChan,
	)
	err := notifier.HandleHeader(&types.Header{Number: big.NewInt(5)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	targetBlockNum := <-targetBlockNumChan
	if targetBlockNum != 6 {
		t.Fatalf("expected target block number %d, got %d", 6, targetBlockNum)
	}
	err = notifier.HandleHeader(&types.Header{Number: big.NewInt(15)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	targetBlockNum = <-targetBlockNumChan
	if targetBlockNum != 16 {
		t.Fatalf("expected target block number %d, got %d", 16, targetBlockNum)
	}

	pastTime := time.Now().Add(-100 * time.Second) // To ensure sendTargetBlockNotification is called immediately

	err = notifier.HandleHeader(&types.Header{Number: big.NewInt(25), Time: uint64(pastTime.Unix())})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// draining starts here
	err = notifier.HandleHeader(&types.Header{Number: big.NewInt(35), Time: uint64(pastTime.Unix())})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = notifier.HandleHeader(&types.Header{Number: big.NewInt(45), Time: uint64(pastTime.Unix())})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
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
