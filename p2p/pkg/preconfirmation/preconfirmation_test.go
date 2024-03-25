package preconfirmation_test

import (
	"context"
	"io"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	p2ptest "github.com/primevprotocol/mev-commit/p2p/pkg/p2p/testing"
	"github.com/primevprotocol/mev-commit/p2p/pkg/preconfirmation"
	"github.com/primevprotocol/mev-commit/p2p/pkg/topology"
)

type testTopo struct {
	peer p2p.Peer
}

func (t *testTopo) GetPeers(q topology.Query) []p2p.Peer {
	return []p2p.Peer{t.peer}
}

type testBidderStore struct{}

func (t *testBidderStore) CheckBidderAllowance(_ context.Context, _ common.Address) bool {
	return true
}

type testSigner struct {
	bid                   *preconfpb.Bid
	preConfirmation       *preconfpb.PreConfirmation
	bidSigner             common.Address
	preConfirmationSigner common.Address
}

func (t *testSigner) ConstructSignedBid(_ string, _ string, _ int64, _ int64, _ int64) (*preconfpb.Bid, error) {
	return t.bid, nil
}

func (t *testSigner) ConstructPreConfirmation(_ *preconfpb.Bid) (*preconfpb.PreConfirmation, error) {
	return t.preConfirmation, nil
}

func (t *testSigner) VerifyBid(_ *preconfpb.Bid) (*common.Address, error) {
	return &t.bidSigner, nil
}

func (t *testSigner) VerifyPreConfirmation(_ *preconfpb.PreConfirmation) (*common.Address, error) {
	return &t.preConfirmationSigner, nil
}

type testProcessor struct {
	status providerapiv1.BidResponse_Status
}

func (t *testProcessor) ProcessBid(
	_ context.Context,
	_ *preconfpb.Bid,
) (chan providerapiv1.BidResponse_Status, error) {
	statusC := make(chan providerapiv1.BidResponse_Status, 1)
	statusC <- t.status
	return statusC, nil
}

type testCommitmentDA struct{}

func (t *testCommitmentDA) StoreCommitment(
	_ context.Context,
	_ *big.Int,
	_ uint64,
	_ string,
	_ uint64,
	_ uint64,
	_ []byte,
	_ []byte,
) error {
	return nil
}

func (t *testCommitmentDA) Close() error {
	return nil
}

func newTestLogger(t *testing.T, w io.Writer) *slog.Logger {
	t.Helper()

	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func TestPreconfBidSubmission(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		client := p2p.Peer{
			EthAddress: common.HexToAddress("0x1"),
			Type:       p2p.PeerTypeBidder,
		}
		server := p2p.Peer{
			EthAddress: common.HexToAddress("0x2"),
			Type:       p2p.PeerTypeProvider,
		}

		bid := &preconfpb.Bid{
			TxHash:              "test",
			BidAmount:           "10",
			BlockNumber:         10,
			DecayStartTimestamp: time.Now().UnixMilli() - 10000*time.Millisecond.Milliseconds(),
			DecayEndTimestamp:   time.Now().UnixMilli(),
			Digest:              []byte("test"),
			Signature:           []byte("test"),
		}

		preConfirmation := &preconfpb.PreConfirmation{
			Bid:       bid,
			Digest:    []byte("test"),
			Signature: []byte("test"),
		}

		svc := p2ptest.New(
			&client,
		)

		topo := &testTopo{server}
		us := &testBidderStore{}
		proc := &testProcessor{
			status: providerapiv1.BidResponse_STATUS_ACCEPTED,
		}
		signer := &testSigner{
			bid:                   bid,
			preConfirmation:       preConfirmation,
			bidSigner:             common.HexToAddress("0x1"),
			preConfirmationSigner: common.HexToAddress("0x2"),
		}

		p := preconfirmation.New(
			topo,
			svc,
			signer,
			us,
			proc,
			&testCommitmentDA{},
			newTestLogger(t, os.Stdout),
		)

		svc.SetPeerHandler(server, p.Streams()[0])

		respC, err := p.SendBid(context.Background(), bid.TxHash, bid.BidAmount, bid.BlockNumber, bid.DecayStartTimestamp, bid.DecayEndTimestamp)
		if err != nil {
			t.Fatal(err)
		}

		commitment := <-respC

		if string(commitment.Digest) != "test" {
			t.Fatalf("data hash is not equal to test")
		}

		if string(commitment.Signature) != "test" {
			t.Fatalf("preConfirmation signature is not equal to test")
		}
	})
}
