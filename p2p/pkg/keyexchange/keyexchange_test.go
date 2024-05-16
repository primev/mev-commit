package keyexchange_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/primev/mev-commit/p2p/pkg/keyexchange"
	"github.com/primev/mev-commit/p2p/pkg/keykeeper"
	mockkeysigner "github.com/primev/mev-commit/p2p/pkg/keykeeper/keysigner/mock"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
	"github.com/primev/mev-commit/p2p/pkg/signer"
	"github.com/primev/mev-commit/p2p/pkg/store"
	"github.com/primev/mev-commit/p2p/pkg/topology"
)

type testTopology struct {
	peers []p2p.Peer
}

func (tt *testTopology) GetPeers(q topology.Query) []p2p.Peer {
	return tt.peers
}

func newTestLogger(t *testing.T, w io.Writer) *slog.Logger {
	t.Helper()

	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func TestKeyExchange_SendAndHandleTimestampMessage(t *testing.T) {
	t.Parallel()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	ks := mockkeysigner.NewMockKeySigner(privKey, address)
	bidderKK, err := keykeeper.NewBidderKeyKeeper(ks)
	if err != nil {
		t.Fatalf("Failed to create BidderKeyKeeper: %v", err)
	}

	providerKK, err := keykeeper.NewProviderKeyKeeper(ks)
	if err != nil {
		t.Fatalf("Failed to create ProviderKeyKeeper: %v", err)
	}

	bidderPeer := p2p.Peer{
		EthAddress: bidderKK.KeySigner.GetAddress(),
		Type:       p2p.PeerTypeBidder,
	}

	providerPeer := p2p.Peer{
		EthAddress: providerKK.KeySigner.GetAddress(),
		Type:       p2p.PeerTypeProvider,
		Keys:       &p2p.Keys{PKEPublicKey: providerKK.GetECIESPublicKey(), NIKEPublicKey: providerKK.GetNIKEPublicKey()},
	}
	topo1 := &testTopology{peers: []p2p.Peer{providerPeer}}
	topo2 := &testTopology{peers: []p2p.Peer{bidderPeer}}

	logger := newTestLogger(t, os.Stdout)

	signer := signer.New()
	svc1 := p2ptest.New(
		&bidderPeer,
	)

	svc2 := p2ptest.New(
		&providerPeer,
	)

	bidderStore := store.NewStore()
	providerStore := store.NewStore()

	ke1 := keyexchange.New(topo1, svc1, bidderKK, bidderStore, logger, signer)
	ke2 := keyexchange.New(topo2, svc2, providerKK, providerStore, logger, signer)

	svc1.SetPeerHandler(bidderPeer, ke2.Streams()[0])

	err = ke1.SendTimestampMessage()
	if err != nil {
		t.Fatalf("SendTimestampMessage failed: %v", err)
	}

	start := time.Now()
	for {
		if time.Since(start) > 5*time.Second {
			t.Fatal("timed out")
		}
		providerAesKey, err := providerStore.GetAESKey(bidderPeer.EthAddress)
		if err != nil {
			t.Fatal(err)
		}
		bidderAesKey, err := bidderStore.GetAESKey(bidderPeer.EthAddress)
		if err != nil {
			t.Fatal(err)
		}
		if providerAesKey != nil {
			if !bytes.Equal(providerAesKey, bidderAesKey) {
				t.Fatal("AES keys are not equal")
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
