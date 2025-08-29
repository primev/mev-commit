package preconfirmation_test

import (
	"context"
	"crypto/rand"
	"io"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	dm "github.com/primev/mev-commit/p2p/pkg/depositmanager"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	providerapi "github.com/primev/mev-commit/p2p/pkg/rpc/provider"
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
	nikePrivateKey           *fr.Element
	preConfirmation          *preconfpb.PreConfirmation
	sharedSecretKey          []byte
	bidSigner                common.Address
	preConfirmationSigner    common.Address
}

func (t *testEncryptor) ConstructEncryptedBid(_ *preconfpb.Bid) (*preconfpb.EncryptedBid, *fr.Element, error) {
	return t.encryptedBid, t.nikePrivateKey, nil
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

func (t *testEncryptor) VerifyEncryptedPreConfirmation(
	_ *preconfpb.Bid,
	_ *bn254.G1Affine,
	_ *fr.Element,
	_ *preconfpb.EncryptedPreConfirmation,
) ([]byte, *common.Address, error) {
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

func (t *testCommitmentDA) StoreUnopenedCommitment(
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
	bidderAddr common.Address,
	providerAddr common.Address,
	bidAmountStr string,
) (func() error, error) {
	return func() error { return nil }, nil
}

func (t *testDepositManager) AddPendingRefund(
	commitmentDigest dm.CommitmentDigest,
	bidder common.Address,
	provider common.Address,
	amount *big.Int,
) {

}

type testTracker struct{}

func (t *testTracker) TrackCommitment(
	ctx context.Context,
	cm *store.Commitment,
	txn *types.Transaction,
) error {
	return nil
}

func TestPreconfBidSubmission(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		client := p2p.Peer{
			EthAddress: common.HexToAddress("0x1"),
			Type:       p2p.PeerTypeBidder,
		}

		encryptionPrivateKey, err := ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
		if err != nil {
			t.Fatal(err)
		}

		_, nikePublicKey, err := p2pcrypto.GenerateKeyPairBN254()
		if err != nil {
			t.Fatal(err)
		}
		server := p2p.Peer{
			EthAddress: common.HexToAddress("0x2"),
			Type:       p2p.PeerTypeProvider,
			Keys: &p2p.Keys{
				PKEPublicKey:  &encryptionPrivateKey.PublicKey,
				NIKEPublicKey: nikePublicKey,
			},
		}

		bid := &preconfpb.Bid{
			TxHash:              "test",
			BidAmount:           "10",
			SlashAmount:         "0",
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
			30*time.Second,
			newTestLogger(t, os.Stdout),
		)

		svc.SetPeerHandler(server, p.Streams()[0])

		respC, err := p.SendBid(context.Background(), bid)
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
