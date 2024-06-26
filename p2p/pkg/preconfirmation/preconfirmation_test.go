package preconfirmation_test

import (
	"context"
	"crypto/ecdh"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation"
	providerapi "github.com/primev/mev-commit/p2p/pkg/rpc/provider"
	"github.com/primev/mev-commit/p2p/pkg/store"
	"github.com/primev/mev-commit/p2p/pkg/topology"
)

type testTopo struct {
	peer p2p.Peer
}

func (t *testTopo) GetPeers(q topology.Query) []p2p.Peer {
	return []p2p.Peer{t.peer}
}

type testEncryptor struct {
	bidHash                  []byte
	encryptedBid             *preconfpb.EncryptedBid
	bid                      *preconfpb.Bid
	encryptedPreConfirmation *preconfpb.EncryptedPreConfirmation
	nikePrivateKey           *ecdh.PrivateKey
	preConfirmation          *preconfpb.PreConfirmation
	sharedSecretKey          []byte
	bidSigner                common.Address
	preConfirmationSigner    common.Address
}

func (t *testEncryptor) ConstructEncryptedBid(_ string, _ string, _ int64, _ int64, _ int64, _ string) (*preconfpb.Bid, *preconfpb.EncryptedBid, *ecdh.PrivateKey, error) {
	return t.bid, t.encryptedBid, t.nikePrivateKey, nil
}

func (t *testEncryptor) ConstructEncryptedPreConfirmation(_ *preconfpb.Bid) (*preconfpb.PreConfirmation, *preconfpb.EncryptedPreConfirmation, error) {
	return t.preConfirmation, t.encryptedPreConfirmation, nil
}

func (t *testEncryptor) VerifyBid(_ *preconfpb.Bid) (*common.Address, error) {
	return &t.bidSigner, nil
}

func (t *testEncryptor) VerifyPreConfirmation(_ *preconfpb.PreConfirmation) (*common.Address, error) {
	return &t.preConfirmationSigner, nil
}

func (t *testEncryptor) DecryptBidData(_ common.Address, _ *preconfpb.EncryptedBid) (*preconfpb.Bid, error) {
	return t.bid, nil
}

func (t *testEncryptor) VerifyEncryptedPreConfirmation(*ecdh.PublicKey, *ecdh.PrivateKey, []byte, *preconfpb.EncryptedPreConfirmation) ([]byte, *common.Address, error) {
	return t.sharedSecretKey, &t.preConfirmationSigner, nil
}

type testProcessor struct {
	BidResponse providerapi.ProcessedBidResponse
}

func (t *testProcessor) ProcessBid(
	_ context.Context,
	_ *preconfpb.Bid,
) (chan providerapi.ProcessedBidResponse, error) {
	statusC := make(chan providerapi.ProcessedBidResponse, 1)
	statusC <- providerapi.ProcessedBidResponse{Status: t.BidResponse.Status, DispatchTimestamp: t.BidResponse.DispatchTimestamp}
	return statusC, nil
}

type testCommitmentDA struct{}

func (t *testCommitmentDA) StoreEncryptedCommitment(
	_ *bind.TransactOpts,
	_ [32]byte,
	_ []byte,
	_ uint64,
) (*types.Transaction, error) {
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}

func newTestLogger(t *testing.T, w io.Writer) *slog.Logger {
	t.Helper()

	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

type testDepositManager struct{}

func (t *testDepositManager) CheckAndDeductDeposit(
	ctx context.Context,
	address common.Address,
	bidAmountStr string,
	blockNumber int64,
) (func() error, error) {
	return func() error { return nil }, nil
}

type testTracker struct{}

func (t *testTracker) TrackCommitment(ctx context.Context, cm *store.EncryptedPreConfirmationWithDecrypted) error {
	return nil
}

func TestPreconfBidSubmission(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		client := p2p.Peer{
			EthAddress: common.HexToAddress("0x1"),
			Type:       p2p.PeerTypeBidder,
		}

		encryptionPrivateKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
		if err != nil {
			t.Fatal(err)
		}

		nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		server := p2p.Peer{
			EthAddress: common.HexToAddress("0x2"),
			Type:       p2p.PeerTypeProvider,
			Keys: &p2p.Keys{
				PKEPublicKey:  &encryptionPrivateKey.PublicKey,
				NIKEPublicKey: nikePrivateKey.PublicKey(),
			},
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

		encryptedBid := &preconfpb.EncryptedBid{
			Ciphertext: []byte("test"),
		}

		preConfirmation := &preconfpb.PreConfirmation{
			Bid:       bid,
			Digest:    []byte("test"),
			Signature: []byte("test"),
		}

		encryptedPreConfirmation := &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("test"),
			Signature:  []byte("test"),
		}
		svc := p2ptest.New(
			&client,
		)

		topo := &testTopo{server}
		proc := &testProcessor{
			BidResponse: providerapi.ProcessedBidResponse{
				Status:            providerapiv1.BidResponse_STATUS_ACCEPTED,
				DispatchTimestamp: 10,
			},
		}
		signer := &testEncryptor{
			bidHash:                  bid.Digest,
			encryptedBid:             encryptedBid,
			bid:                      bid,
			preConfirmation:          preConfirmation,
			encryptedPreConfirmation: encryptedPreConfirmation,
			bidSigner:                common.HexToAddress("0x1"),
			preConfirmationSigner:    common.HexToAddress("0x2"),
			revertingTxHashes:        "",
		}

		depositMgr := &testDepositManager{}
		p := preconfirmation.New(
			topo,
			svc,
			signer,
			depositMgr,
			proc,
			&testCommitmentDA{},
			&testTracker{},
			func(context.Context) (*bind.TransactOpts, error) {
				return &bind.TransactOpts{
					From: client.EthAddress,
				}, nil
			},
			newTestLogger(t, os.Stdout),
		)

		svc.SetPeerHandler(server, p.Streams()[0])

		respC, err := p.SendBid(context.Background(), bid.TxHash, bid.BidAmount, bid.BlockNumber, bid.DecayStartTimestamp, bid.DecayEndTimestamp, "")
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
