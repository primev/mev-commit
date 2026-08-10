package bidder_test

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
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
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

	bidder.SetNowFunc(timeSetter.Now)

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
	bidderClient := bidder.NewBidderClient(
		util.NewTestLogger(os.Stdout),
		rpcServices,
		rpcServices,
		rpcServices,
		blockNumberGetter,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := bidderClient.Start(ctx)

	_, err = bidderClient.Estimate()
	if err != bidder.ErrNoEpochInfo {
		t.Fatalf("expected error %v, got %v", bidder.ErrNoEpochInfo, err)
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

	buf2 := make([]byte, 32)
	_, _ = rand.Read(buf2)
	txString2 := hex.EncodeToString(buf2)

	rpcServices.topo = &debugapiv1.TopologyResponse{
		Topology: topoVal,
	}

	providers, err := bidderClient.ConnectedProviders(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}

	statusC, err := bidderClient.Bid(ctx, big.NewInt(1), big.NewInt(1), []string{txString, txString2}, nil)
	if err != nil {
		t.Fatal(err)
	}

	commitments := 0
waitLoop:
	for {
		select {
		case status, more := <-statusC:
			if !more {
				break waitLoop
			}
			switch status.Type {
			case bidder.BidStatusWaitSecs:
				if status.Arg.(int) != 2 {
					t.Fatalf("expected 2 seconds, got %d", status.Arg)
				}
			case bidder.BidStatusAttempted:
				if status.Arg.(uint64) != 11 {
					t.Fatalf("expected 11, got %d", status.Arg)
				}
			case bidder.BidStatusCommitment:
				if status.Arg.(*bidderapiv1.Commitment).BlockNumber != 11 {
					t.Fatalf("expected block number 11, got %d", status.Arg.(*bidderapiv1.Commitment).BlockNumber)
				}
				commitments++
			}
		case bid := <-rpcServices.bidChan:
			if bid.Amount != big.NewInt(1).String() {
				t.Fatalf("expected amount 1, got %s", bid.Amount)
			}
			if bid.BlockNumber != 11 {
				t.Fatalf("expected block number 11, got %d", bid.BlockNumber)
			}
			if len(bid.RawTransactions) != 2 {
				t.Fatalf("expected 2 raw transactions, got %d", len(bid.RawTransactions))
			}
			if bid.RawTransactions[0] != txString {
				t.Fatalf("expected raw transaction %x, got %s", buf, bid.RawTransactions[0])
			}
			if bid.RawTransactions[1] != txString2 {
				t.Fatalf("expected raw transaction %x, got %s", buf2, bid.RawTransactions[1])
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

func TestSubscribeSettlementsSplit(t *testing.T) {
	t.Parallel()

	rpcServices := &testRPCServices{
		notificationChan: make(chan *notificationsapiv1.Notification),
	}
	bidderClient := bidder.NewBidderClient(
		util.NewTestLogger(os.Stdout),
		rpcServices,
		rpcServices,
		rpcServices,
		&testBlockNumberGetter{blockNumber: 10},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	settlements := bidderClient.SubscribeSettlements(ctx)

	val, err := structpb.NewStruct(map[string]interface{}{
		"transaction_hashes": "aaa,bbb",
		"provider":           "provider1",
		"is_slashed":         true,
	})
	if err != nil {
		t.Fatal(err)
	}

	rpcServices.notificationChan <- &notificationsapiv1.Notification{
		Topic: "transaction_settled",
		Value: val,
	}

	for _, want := range []string{"aaa", "bbb"} {
		msg := <-settlements
		if msg.TransactionHash != want {
			t.Fatalf("expected transaction hash %s, got %s", want, msg.TransactionHash)
		}
		if msg.Provider != "provider1" {
			t.Fatalf("expected provider provider1, got %s", msg.Provider)
		}
		if !msg.IsSlash {
			t.Fatal("expected is slashed to be true")
		}
	}
}

func TestSubscribePaymentsSplit(t *testing.T) {
	t.Parallel()

	rpcServices := &testRPCServices{
		notificationChan: make(chan *notificationsapiv1.Notification),
	}
	bidderClient := bidder.NewBidderClient(
		util.NewTestLogger(os.Stdout),
		rpcServices,
		rpcServices,
		rpcServices,
		&testBlockNumberGetter{blockNumber: 10},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payments := bidderClient.SubscribePayments(ctx)

	val, err := structpb.NewStruct(map[string]interface{}{
		"transaction_hashes": "aaa,bbb",
		"payment":            "100",
		"refund":             "40",
	})
	if err != nil {
		t.Fatal(err)
	}

	rpcServices.notificationChan <- &notificationsapiv1.Notification{
		Topic: "transaction_payment",
		Value: val,
	}

	for _, want := range []string{"aaa", "bbb"} {
		msg := <-payments
		if msg.TransactionHash != want {
			t.Fatalf("expected transaction hash %s, got %s", want, msg.TransactionHash)
		}
		if msg.Payment.Cmp(big.NewInt(100)) != 0 {
			t.Fatalf("expected payment 100, got %s", msg.Payment)
		}
		if msg.Refund.Cmp(big.NewInt(40)) != 0 {
			t.Fatalf("expected refund 40, got %s", msg.Refund)
		}
		if len(msg.BundleTxnHashes) != 2 || msg.BundleTxnHashes[0] != "aaa" || msg.BundleTxnHashes[1] != "bbb" {
			t.Fatalf("expected bundle hashes [aaa bbb], got %v", msg.BundleTxnHashes)
		}
	}
}
