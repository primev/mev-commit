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

func TestBidderDedupChecksCorrectMap(t *testing.T) {
	t.Parallel()

	notifier := &testNotifier{}
	topo := topology.New(&testAddressbook{}, notifier, util.NewTestLogger(io.Discard))

	// Provider and bidder with the same eth address
	addr := common.HexToAddress("0xABC")

	provider := p2p.Peer{EthAddress: addr, Type: p2p.PeerTypeProvider}
	bidder := p2p.Peer{EthAddress: addr, Type: p2p.PeerTypeBidder}

	topo.Connected(provider)

	// Before the fix, this bidder would be silently dropped because `add`
	// checked t.providers (wrong map) for the bidder dedup.
	topo.Connected(bidder)

	providers := topo.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	bidders := topo.GetPeers(topology.Query{Type: p2p.PeerTypeBidder})

	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}
	if len(bidders) != 1 {
		t.Fatalf("expected 1 bidder, got %d", len(bidders))
	}

	// Duplicate bidder should be deduped
	topo.Connected(bidder)
	bidders = topo.GetPeers(topology.Query{Type: p2p.PeerTypeBidder})
	if len(bidders) != 1 {
		t.Fatalf("expected 1 bidder after duplicate connect, got %d", len(bidders))
	}
}

func TestNewProviderBroadcastsToAllBidders(t *testing.T) {
	t.Parallel()

	ann := &announcer{}
	notifier := &testNotifier{}
	topo := topology.New(&testAddressbook{}, notifier, util.NewTestLogger(io.Discard))
	topo.SetAnnouncer(ann)

	// Connect 3 bidders first
	bidders := []p2p.Peer{
		{EthAddress: common.HexToAddress("0xB1"), Type: p2p.PeerTypeBidder},
		{EthAddress: common.HexToAddress("0xB2"), Type: p2p.PeerTypeBidder},
		{EthAddress: common.HexToAddress("0xB3"), Type: p2p.PeerTypeBidder},
	}
	for _, b := range bidders {
		topo.Connected(b)
	}

	ann.mu.Lock()
	ann.broadcasts = nil // clear broadcasts from bidder connections
	ann.mu.Unlock()

	// Now connect a new provider — it should be broadcast to all 3 bidders
	provider := p2p.Peer{EthAddress: common.HexToAddress("0xP1"), Type: p2p.PeerTypeProvider}
	topo.Connected(provider)

	ann.mu.Lock()
	defer ann.mu.Unlock()

	if len(ann.broadcasts) != 3 {
		t.Fatalf("expected new provider to be broadcast to 3 bidders, got %d", len(ann.broadcasts))
	}

	broadcastAddrs := make(map[common.Address]bool)
	for _, b := range ann.broadcasts {
		broadcastAddrs[b.EthAddress] = true
	}
	for _, b := range bidders {
		if !broadcastAddrs[b.EthAddress] {
			t.Fatalf("bidder %s did not receive provider broadcast", b.EthAddress.Hex())
		}
	}
}
