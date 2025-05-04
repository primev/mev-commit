package notifier_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
	"github.com/primev/mev-commit/tools/bidder-bot/notifier"
	"github.com/primev/mev-commit/x/util"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type MockBeaconClient struct{}

func (m *MockBeaconClient) GetPayloadDataForSlot(ctx context.Context, slot uint64) (
	blockNumber uint64,
	timestamp uint64,
	err error,
) {
	if slot > 20 {
		return 0, 0, fmt.Errorf("slot %d is greater than 20", slot)
	}
	blockNumber = slot + 10
	timestamp = blockNumber * 12
	return blockNumber, timestamp, nil
}

func TestHandleMsg(t *testing.T) {
	testcases := []struct {
		name                string
		msg                 *notificationsapiv1.Notification
		expectedError       bool
		expectedTargetBlock *bidder.TargetBlock
	}{
		{
			name: "no missed slot, no error",
			msg: &notificationsapiv1.Notification{
				Topic: "validator_opted_in",
				Value: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"slot":    structpb.NewNumberValue(15),
						"epoch":   structpb.NewNumberValue(6),
						"bls_key": structpb.NewStringValue("0x123"),
					},
				},
			},
			expectedError: false,
			expectedTargetBlock: &bidder.TargetBlock{
				Num:  25,
				Time: time.Unix(int64(25*12), 0),
			},
		},
		{
			name: "missed slot, no error",
			msg: &notificationsapiv1.Notification{
				Topic: "validator_opted_in",
				Value: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"slot":    structpb.NewNumberValue(23), // 23 - 2 = 21, beacon client will return error for slot 21, not 20
						"epoch":   structpb.NewNumberValue(6),
						"bls_key": structpb.NewStringValue("0x123"),
					},
				},
			},
			expectedError: false,
			expectedTargetBlock: &bidder.TargetBlock{
				Num:  32, // block num is 9 higher than slot number (not 10) to account for missed slot
				Time: time.Unix(int64(32*12), 0),
			},
		},
		{
			name: "two missed slots, error",
			msg: &notificationsapiv1.Notification{
				Topic: "validator_opted_in",
				Value: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"slot":    structpb.NewNumberValue(30), // 30 - 2 = 28, beacon client will return error for slot 28 and 27
						"epoch":   structpb.NewNumberValue(6),
						"bls_key": structpb.NewStringValue("0x123"),
					},
				},
			},
			expectedError: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			logger := util.NewTestLogger(os.Stdout)
			targetBlockChan := make(chan bidder.TargetBlock, 1)
			mockBeaconClient := &MockBeaconClient{}
			notifier := notifier.NewSelectiveNotifier(logger, nil, mockBeaconClient, targetBlockChan)
			err := notifier.HandleMsg(context.Background(), testcase.msg)
			if testcase.expectedError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !testcase.expectedError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if testcase.expectedTargetBlock != nil {
				timeout := time.After(1 * time.Second)
				select {
				case targetBlock := <-targetBlockChan:
					if targetBlock != *testcase.expectedTargetBlock {
						t.Errorf("expected target block %v, got %v", testcase.expectedTargetBlock, targetBlock)
					}
				case <-timeout:
					t.Errorf("expected target block, got timeout")
				}
			}
		})
	}
}
