package libp2p_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	mockkeysigner "github.com/primevprotocol/mev-commit/p2p/pkg/keysigner/mock"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p/libp2p"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type testRegistry struct{}

func (t *testRegistry) CheckProviderRegistered(
	_ context.Context,
	_ common.Address,
) bool {
	return true
}

func newTestLogger(t *testing.T, w io.Writer) *slog.Logger {
	t.Helper()

	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func newTestService(t *testing.T) *libp2p.Service {
	t.Helper()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	ks := mockkeysigner.NewMockKeySigner(privKey, address)
	svc, err := libp2p.New(&libp2p.Options{
		KeySigner:  ks,
		Secret:     "test",
		ListenPort: 0,
		ListenAddr: "0.0.0.0",
		PeerType:   p2p.PeerTypeProvider,
		Register:   &testRegistry{},
		Logger:     newTestLogger(t, os.Stdout),
	})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestP2PService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping protocols test in short mode")
	}

	t.Run("new and close", func(t *testing.T) {
		svc := newTestService(t)

		err := svc.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("add protocol and connect", func(t *testing.T) {
		svc := newTestService(t)
		client := newTestService(t)

		t.Cleanup(func() {
			err := errors.Join(svc.Close(), client.Close())
			if err != nil {
				t.Fatal(err)
			}
		})

		done := make(chan struct{})
		stream := p2p.StreamDesc{
			Name:    "test",
			Version: "1.0.0",
			Handler: func(ctx context.Context, peer p2p.Peer, str p2p.Stream) error {
				if peer.EthAddress.Cmp(client.Peer().EthAddress) != 0 {
					t.Fatalf(
						"expected eth address %s, got %s",
						client.Peer().EthAddress, peer.EthAddress,
					)
				}

				if peer.Type != client.Peer().Type {
					t.Fatalf(
						"expected peer type %s, got %s",
						client.Peer().Type, peer.Type,
					)
				}

				strMsg := new(wrapperspb.StringValue)

				err := str.ReadMsg(ctx, strMsg)
				if err != nil {
					t.Fatal(err)
				}
				if strMsg.Value != "test" {
					t.Fatalf("expected message %s, got %s", "test", strMsg.Value)
				}
				close(done)
				return nil
			},
		}

		svc.AddStreamHandlers(stream)

		svAddr, err := svc.Addrs()
		if err != nil {
			t.Fatal(err)
		}

		p, err := client.Connect(context.Background(), svAddr)
		if err != nil {
			t.Fatal(err)
		}

		if p.EthAddress.Hex() != svc.Peer().EthAddress.Hex() {
			t.Fatalf(
				"expected eth address %s, got %s",
				svc.Peer().EthAddress.Hex(), p.EthAddress.Hex(),
			)
		}

		if p.Type != svc.Peer().Type {
			t.Fatalf(
				"expected peer type %s, got %s",
				svc.Peer().Type.String(), p.Type.String(),
			)
		}

		str, err := client.NewStream(context.Background(), p, nil, stream)
		if err != nil {
			t.Fatal(err)
		}

		err = str.WriteMsg(context.Background(), &wrapperspb.StringValue{Value: "test"})
		if err != nil {
			t.Fatal(err)
		}

		<-done

		err = str.Close()
		if err != nil {
			t.Fatal(err)
		}

		svcInfo, err := client.GetPeerInfo(svc.Peer())
		if err != nil {
			t.Fatal(err)
		}

		var svcAddr peer.AddrInfo
		err = svcAddr.UnmarshalJSON(svcInfo)
		if err != nil {
			t.Fatal(err)
		}

		if svcAddr.ID != svc.HostID() {
			t.Fatalf("expected host id %s, got %s", svc.HostID(), svcAddr.ID)
		}

		clientInfo, err := svc.GetPeerInfo(client.Peer())
		if err != nil {
			t.Fatal(err)
		}

		var clientAddr peer.AddrInfo
		err = clientAddr.UnmarshalJSON(clientInfo)
		if err != nil {
			t.Fatal(err)
		}

		if clientAddr.ID != client.HostID() {
			t.Fatalf("expected host id %s, got %s", client.HostID(), clientAddr.ID)
		}
	})
}

type testNotifier struct {
	mu    sync.Mutex
	peers []p2p.Peer
}

func (t *testNotifier) Connected(p p2p.Peer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peers = append(t.peers, p)
}

func (t *testNotifier) Disconnected(p p2p.Peer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, peer := range t.peers {
		if peer.EthAddress == p.EthAddress {
			t.peers = append(t.peers[:i], t.peers[i+1:]...)
			return
		}
	}
}

func (t *testNotifier) Peers() []p2p.Peer {
	t.mu.Lock()
	defer t.mu.Unlock()

	peers := make([]p2p.Peer, len(t.peers))
	copy(peers, t.peers)
	return peers
}

func TestBootstrap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping bootstrap test in short mode")
	}

	testDefaultOptions := libp2p.Options{
		Secret:     "test",
		ListenPort: 0,
		ListenAddr: "0.0.0.0",
		PeerType:   p2p.PeerTypeProvider,
		Register:   &testRegistry{},
		Logger:     newTestLogger(t, os.Stdout),
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	ks := mockkeysigner.NewMockKeySigner(privKey, address)

	bnOpts := testDefaultOptions
	bnOpts.KeySigner = ks
	bnOpts.PeerType = p2p.PeerTypeBootnode

	bootnode, err := libp2p.New(&bnOpts)
	if err != nil {
		t.Fatal(err)
	}

	notifier := &testNotifier{}
	bootnode.SetNotifier(notifier)

	privKey, err = crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	address = crypto.PubkeyToAddress(privKey.PublicKey)
	ks = mockkeysigner.NewMockKeySigner(privKey, address)

	n1Opts := testDefaultOptions
	n1Opts.BootstrapAddrs = []string{bootnode.AddrString()}
	n1Opts.KeySigner = ks

	p1, err := libp2p.New(&n1Opts)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for {
		if time.Since(start) > 10*time.Second {
			t.Fatal("timed out waiting for peers to connect")
		}

		if p1.PeerCount() == 1 {
			peers := notifier.Peers()
			if len(peers) != 1 {
				t.Fatalf("expected 1 peer, got %d", len(peers))
			}
			if peers[0].Type != p2p.PeerTypeProvider {
				t.Fatalf(
					"expected peer type %s, got %s",
					p2p.PeerTypeProvider, peers[0].Type,
				)
			}
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestHandlerError(t *testing.T) {
	svc := newTestService(t)
	client := newTestService(t)

	t.Cleanup(func() {
		err := errors.Join(svc.Close(), client.Close())
		if err != nil {
			t.Fatal(err)
		}
	})

	stream := p2p.StreamDesc{
		Name:    "test",
		Version: "1.0.0",
		Handler: func(ctx context.Context, peer p2p.Peer, str p2p.Stream) error {
			return status.Error(codes.Internal, "test error")
		},
	}

	svc.AddStreamHandlers(stream)

	svAddr, err := svc.Addrs()
	if err != nil {
		t.Fatal(err)
	}

	p, err := client.Connect(context.Background(), svAddr)
	if err != nil {
		t.Fatal(err)
	}

	if p.EthAddress.Hex() != svc.Peer().EthAddress.Hex() {
		t.Fatalf(
			"expected eth address %s, got %s",
			svc.Peer().EthAddress.Hex(), p.EthAddress.Hex(),
		)
	}

	if p.Type != svc.Peer().Type {
		t.Fatalf(
			"expected peer type %s, got %s",
			svc.Peer().Type.String(), p.Type.String(),
		)
	}

	str, err := client.NewStream(context.Background(), p, nil, stream)
	if err != nil {
		t.Fatal(err)
	}

	err = str.WriteMsg(context.Background(), &wrapperspb.StringValue{Value: "test"})
	if err != nil {
		t.Fatal(err)
	}

	err = str.ReadMsg(context.Background(), &wrapperspb.StringValue{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	assert.Equal(t, codes.Internal, status.Convert(err).Code())
	assert.Equal(t, "test error", status.Convert(err).Message())
}
