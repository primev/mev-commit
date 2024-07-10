package keyexchange_test

import (
	"bytes"
	"crypto/ecdh"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/keyexchange"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
	"github.com/primev/mev-commit/p2p/pkg/signer"
	"github.com/primev/mev-commit/p2p/pkg/store"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	mockkeysigner "github.com/primev/mev-commit/x/keysigner/mock"
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

	bidderPeer := p2p.Peer{
		EthAddress: ks.GetAddress(),
		Type:       p2p.PeerTypeBidder,
	}

	bidderStore := store.NewStore()
	providerStore := store.NewStore()

	encryptionPrivateKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = providerStore.SetECIESPrivateKey(encryptionPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	err = providerStore.SetNikePrivateKey(nikePrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	providerPeer := p2p.Peer{
		EthAddress: ks.GetAddress(),
		Type:       p2p.PeerTypeProvider,
		Keys:       &p2p.Keys{PKEPublicKey: &encryptionPrivateKey.PublicKey, NIKEPublicKey: nikePrivateKey.PublicKey()},
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

	aesKey, err := p2pcrypto.GenerateAESKey()
	if err != nil {
		t.Fatal(err)
	}
	err = bidderStore.SetAESKey(ks.GetAddress(), aesKey)
	if err != nil {
		t.Fatal(err)
	}

	ke1 := keyexchange.New(topo1, svc1, ks, aesKey, bidderStore, logger, signer, nil)
	ke2 := keyexchange.New(topo2, svc2, ks, nil, providerStore, logger, signer, nil)
	if err != nil {
		t.Fatalf("keyexchange new failed: %v", err)
	}
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

func TestKeyExchange_Whitelist(t *testing.T) {
	t.Parallel()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	ks := mockkeysigner.NewMockKeySigner(privKey, address)

	bidderPeer := p2p.Peer{
		EthAddress: ks.GetAddress(),
		Type:       p2p.PeerTypeBidder,
	}

	bidderStore := store.NewStore()
	providerStore := store.NewStore()

	encryptionPrivateKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = providerStore.SetECIESPrivateKey(encryptionPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	err = providerStore.SetNikePrivateKey(nikePrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	providerPeer := p2p.Peer{
		EthAddress: ks.GetAddress(),
		Type:       p2p.PeerTypeProvider,
		Keys:       &p2p.Keys{PKEPublicKey: &encryptionPrivateKey.PublicKey, NIKEPublicKey: nikePrivateKey.PublicKey()},
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

	aesKey, err := p2pcrypto.GenerateAESKey()
	if err != nil {
		t.Fatal(err)
	}
	err = bidderStore.SetAESKey(ks.GetAddress(), aesKey)
	if err != nil {
		t.Fatal(err)
	}

	randomWhitelistedPeerKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	randomWhitelistedPeerAddress := crypto.PubkeyToAddress(randomWhitelistedPeerKey.PublicKey)

	ke1 := keyexchange.New(topo1, svc1, ks, aesKey, bidderStore, logger, signer, []common.Address{randomWhitelistedPeerAddress})
	ke2 := keyexchange.New(topo2, svc2, ks, nil, providerStore, logger, signer, nil)
	if err != nil {
		t.Fatalf("keyexchange new failed: %v", err)
	}
	svc1.SetPeerHandler(bidderPeer, ke2.Streams()[0])

	err = ke1.SendTimestampMessage()
	if !errors.Is(err, keyexchange.ErrNoProvidersAvailable) {
		t.Fatalf("SendTimestampMessage should have failed")
	}
}
