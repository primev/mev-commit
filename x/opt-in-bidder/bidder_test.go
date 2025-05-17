package optinbidder_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder/bidder"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type testRPCServices struct {
	bidderapiv1.BidderClient
	debugapiv1.DebugServiceClient
	notificationsapiv1.NotificationsClient

	notificationChan chan *notificationsapiv1.Notification
	bidChan          chan *bidderapiv1.Bid
	commitmentChan   chan *bidderapiv1.Commitment
	topo             *debugapiv1.TopologyResponse
}

type testNotificationStream struct {
	grpc.ClientStream

	ctx              context.Context
	notificationChan chan *notificationsapiv1.Notification
}

func (t *testNotificationStream) Recv() (*notificationsapiv1.Notification, error) {
	select {
	case <-t.ctx.Done():
		return nil, io.EOF
	case n := <-t.notificationChan:
		return n, nil
	}
}

type testCommitmentStream struct {
	grpc.ClientStream

	ctx            context.Context
	commitmentChan chan *bidderapiv1.Commitment
}

func (t *testCommitmentStream) Recv() (*bidderapiv1.Commitment, error) {
	select {
	case <-t.ctx.Done():
		return nil, io.EOF
	case c, more := <-t.commitmentChan:
		if !more {
			return nil, io.EOF
		}
		return c, nil
	}
}

func (t *testRPCServices) SendBid(
	ctx context.Context,
	in *bidderapiv1.Bid,
	_ ...grpc.CallOption,
) (grpc.ServerStreamingClient[bidderapiv1.Commitment], error) {
	select {
	case t.bidChan <- in:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return &testCommitmentStream{ctx: ctx, commitmentChan: t.commitmentChan}, nil
}

func (t *testRPCServices) GetTopology(
	_ context.Context,
	_ *debugapiv1.EmptyMessage,
	_ ...grpc.CallOption,
) (*debugapiv1.TopologyResponse, error) {
	return t.topo, nil
}

func (t *testRPCServices) Subscribe(
	ctx context.Context,
	in *notificationsapiv1.SubscribeRequest,
	_ ...grpc.CallOption,
) (grpc.ServerStreamingClient[notificationsapiv1.Notification], error) {
	return &testNotificationStream{ctx: ctx, notificationChan: t.notificationChan}, nil
}

type testBlockNumberGetter struct {
	blockNumber uint64
}

func (t *testBlockNumberGetter) BlockNumber(ctx context.Context) (uint64, error) {
	return t.blockNumber, nil
}

type testTimeSetter struct {
	now time.Time
}

func (t *testTimeSetter) Now() time.Time {
	return t.now
}

func TestBidderClient(t *testing.T) {
	t.Parallel()

	clock := time.Now()
	timeSetter := &testTimeSetter{
		now: clock,
	}

	optinbidder.SetNowFunc(timeSetter.Now)

	topoVal, err := structpb.NewStruct(map[string]interface{}{
		"connected_providers": []any{"provider1", "provider2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Create a new test RPC services.
	rpcServices := &testRPCServices{
		notificationChan: make(chan *notificationsapiv1.Notification),
		bidChan:          make(chan *bidderapiv1.Bid),
		commitmentChan:   make(chan *bidderapiv1.Commitment),
		topo: &debugapiv1.TopologyResponse{
			Topology: &structpb.Struct{},
		},
	}

	blockNumberGetter := &testBlockNumberGetter{blockNumber: 10}
	bidderClient := optinbidder.NewBidderClient(
		util.NewTestLogger(os.Stdout),
		rpcServices,
		rpcServices,
		rpcServices,
		blockNumberGetter,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := bidderClient.Start(ctx)

	_, err = bidderClient.Estimate()
	if err != optinbidder.ErrNoEpochInfo {
		t.Fatalf("expected error %v, got %v", optinbidder.ErrNoEpochInfo, err)
	}

	// Send a notification.
	nVal, err := structpb.NewStruct(map[string]interface{}{
		"epoch":            1,
		"epoch_start_time": clock.Add(2 * time.Second).Unix(),
		"slots": []any{
			map[string]interface{}{
				"slot":       33,
				"start_time": clock.Add(14 * time.Second).Unix(),
				"bls_key":    "key2",
				"opted_in":   true,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rpcServices.notificationChan <- &notificationsapiv1.Notification{
		Topic: "epoch_validators_opted_in",
		Value: nVal,
	}

	for {
		if _, err := bidderClient.Estimate(); err == nil {
			break
		}
	}

	estimate, err := bidderClient.Estimate()
	if err != nil {
		t.Fatal(err)
	}
	if estimate != 13 {
		t.Fatalf("expected estimate 13, got %d", estimate)
	}

	timeSetter.now = clock.Add(10 * time.Second)

	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	txString := hex.EncodeToString(buf)

	_, err = bidderClient.Bid(ctx, big.NewInt(1), big.NewInt(1), big.NewInt(1), txString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	rpcServices.topo = &debugapiv1.TopologyResponse{
		Topology: topoVal,
	}

	statusC, err := bidderClient.Bid(ctx, big.NewInt(1), big.NewInt(1), big.NewInt(1), txString)
	if err != nil {
		t.Fatal(err)
	}

	commitments := 0
waitLoop:
	for {
		select {
		case status := <-statusC:
			switch {
			case status.Type == optinbidder.BidStatusNoOfProviders:
				if status.Arg.(int) != 2 {
					t.Fatalf("expected 2 providers, got %d", status.Arg)
				}
			case status.Type == optinbidder.BidStatusWaitSecs:
				if status.Arg.(int) != 2 {
					t.Fatalf("expected 2 seconds, got %d", status.Arg)
				}
			case status.Type == optinbidder.BidStatusAttempted:
				if status.Arg.(int) != 11 {
					t.Fatalf("expected 11, got %d", status.Arg)
				}
			case status.Type == optinbidder.BidStatusCommitment:
				if status.Arg.(*bidderapiv1.Commitment).BlockNumber != 11 {
					t.Fatalf("expected block number 11, got %d", status.Arg.(*bidderapiv1.Commitment).BlockNumber)
				}
				commitments++
			case status.Type == optinbidder.BidStatusDone:
				break waitLoop
			}
		case bid := <-rpcServices.bidChan:
			if bid.Amount != big.NewInt(1).String() {
				t.Fatalf("expected amount 1, got %s", bid.Amount)
			}
			if bid.BlockNumber != 11 {
				t.Fatalf("expected block number 11, got %d", bid.BlockNumber)
			}
			if bid.RawTransactions[0] != txString {
				t.Fatalf("expected raw transaction %x, got %s", buf, bid.RawTransactions[0])
			}
			rpcServices.commitmentChan <- &bidderapiv1.Commitment{
				BlockNumber: 11,
			}
			rpcServices.commitmentChan <- &bidderapiv1.Commitment{
				BlockNumber: 11,
			}
			close(rpcServices.commitmentChan)
		}
	}

	if commitments != 2 {
		t.Fatalf("expected 2 commitments, got %d", commitments)
	}

	cancel()
	<-done
}
