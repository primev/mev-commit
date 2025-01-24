package topology_test

import (
	"context"
	"errors"
	"io"
	"slices"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"github.com/primev/mev-commit/x/util"
)

type testAddressbook struct{}

func (t *testAddressbook) GetPeerInfo(p p2p.Peer) ([]byte, error) {
	return []byte("test"), nil
}

type announcer struct {
	mu         sync.Mutex
	broadcasts []p2p.Peer
}

func (a *announcer) BroadcastPeers(_ context.Context, p p2p.Peer, peers []p2p.PeerInfo) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.broadcasts = append(a.broadcasts, p)

	if len(peers) != 1 {
		return errors.New("wrong number of peers")
	}

	if string(peers[0].Underlay) != "test" {
		return errors.New("wrong peer underlay")
	}

	return nil
}

type testNotifier struct {
	mu           sync.Mutex
	connected    []string
	disconnected []string
}

func (t *testNotifier) Notify(n *notifications.Notification) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch n.Topic() {
	case notifications.TopicPeerConnected:
		t.connected = append(t.connected, n.Value()["ethAddress"].(string))
	case notifications.TopicPeerDisconnected:
		t.disconnected = append(t.disconnected, n.Value()["ethAddress"].(string))
	}
}

func TestTopology(t *testing.T) {
	t.Parallel()

	announcer := &announcer{}
	notifier := &testNotifier{}
	topo := topology.New(&testAddressbook{}, notifier, util.NewTestLogger(io.Discard))
	topo.SetAnnouncer(announcer)

	p1 := p2p.Peer{
		EthAddress: common.HexToAddress("0x1"),
		Type:       p2p.PeerTypeProvider,
	}

	s1 := p2p.Peer{
		EthAddress: common.HexToAddress("0x2"),
		Type:       p2p.PeerTypeBidder,
	}

	topo.Connected(p1)

	topo.Connected(s1)

	if len(announcer.broadcasts) != 1 {
		t.Fatal("expected one broadcast")
	}

	if announcer.broadcasts[0].EthAddress != s1.EthAddress {
		t.Fatal("wrong peer")
	}

	p2 := p2p.Peer{
		EthAddress: common.HexToAddress("0x3"),
		Type:       p2p.PeerTypeProvider,
	}

	topo.AddPeers(p2)

	for _, p := range []p2p.Peer{p1, s1, p2} {
		if !topo.IsConnected(p.EthAddress) {
			t.Fatal("peer not connected")
		}
		if !slices.Contains(notifier.connected, p.EthAddress.Hex()) {
			t.Fatal("peer connected notification not found", p)
		}
	}

	peers := topo.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	if len(peers) != 2 {
		t.Fatal("wrong number of peers")
	}

	if len(notifier.connected) != 3 {
		t.Fatal("wrong number of peers in notifier")
	}

	for _, p := range peers {
		if p.Type != p2p.PeerTypeProvider {
			t.Fatal("wrong peer type")
		}
		if p.EthAddress != p1.EthAddress && p.EthAddress != p2.EthAddress {
			t.Fatal("wrong peer")
		}
	}

	peers = topo.GetPeers(topology.Query{Type: p2p.PeerTypeBidder})
	if len(peers) != 1 {
		t.Fatal("wrong number of peers")
	}

	if peers[0].Type != p2p.PeerTypeBidder {
		t.Fatal("wrong peer type")
	}

	if peers[0].EthAddress != s1.EthAddress {
		t.Fatal("wrong peer")
	}

	topo.Disconnected(p1)

	if topo.IsConnected(p1.EthAddress) {
		t.Fatal("peer still connected")
	}

	if len(notifier.disconnected) != 1 {
		t.Fatal("disconnect notification not found")
	}

	if notifier.disconnected[0] != p1.EthAddress.Hex() {
		t.Fatal("wrong peer in disconnect notification")
	}
}
