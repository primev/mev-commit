package libp2p

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	core "github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
)

type peerRegistry struct {
	// overlays maps peer IDs to peer info relevant to the primev client
	overlays map[core.PeerID]*p2p.Peer
	// underlays maps Ethereum addresses to peer IDs used by libp2p
	underlays map[common.Address]core.PeerID
	// connections maps peer IDs to their open connections, this is used to
	// track all connections to a peer.
	connections map[core.PeerID]map[network.Conn]struct{}
	// streams maps peer IDs to their open streams, this is used to track all
	// streams to a peer and cancel the contexts passed to the handlers when
	// the stream is closed.
	streams map[core.PeerID]map[network.Stream]context.CancelFunc
	mu      sync.RWMutex

	disconnector disconnector
	network.Notifiee
}

type disconnector interface {
	disconnected(p2p.Peer)
}

func newPeerRegistry() *peerRegistry {
	return &peerRegistry{
		overlays:    make(map[core.PeerID]*p2p.Peer),
		underlays:   make(map[common.Address]core.PeerID),
		connections: make(map[core.PeerID]map[network.Conn]struct{}),
		streams:     make(map[core.PeerID]map[network.Stream]context.CancelFunc),
		Notifiee:    new(network.NoopNotifiee),
	}
}

func (r *peerRegistry) setDisconnector(d disconnector) {
	r.disconnector = d
}

func (r *peerRegistry) Disconnected(_ network.Network, c network.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	peerID := c.RemotePeer()
	if _, ok := r.connections[peerID]; !ok {
		return
	}

	delete(r.connections[peerID], c)
	if len(r.connections[peerID]) > 0 {
		// if there are still connections, don't remove the peer
		return
	}

	delete(r.connections, peerID)
	peerInfo := r.overlays[peerID]
	delete(r.overlays, peerID)
	delete(r.underlays, peerInfo.EthAddress)
	for _, cancel := range r.streams[peerID] {
		cancel()
	}
	delete(r.streams, peerID)
	r.disconnector.disconnected(*peerInfo)
}

func (r *peerRegistry) addPeer(c network.Conn, p *p2p.Peer) (exists bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.connections[c.RemotePeer()]; !ok {
		r.connections[c.RemotePeer()] = make(map[network.Conn]struct{})
	}
	r.connections[c.RemotePeer()][c] = struct{}{}

	if _, exists := r.underlays[p.EthAddress]; exists {
		return true
	}

	r.overlays[c.RemotePeer()] = p
	r.underlays[p.EthAddress] = c.RemotePeer()
	r.streams[c.RemotePeer()] = make(map[network.Stream]context.CancelFunc)
	return false
}

func (r *peerRegistry) removePeer(peer *p2p.Peer) (found bool, peerID core.PeerID) { //nolint:unused
	r.mu.Lock()
	defer r.mu.Unlock()

	peerID, found = r.underlays[peer.EthAddress]
	delete(r.overlays, peerID)
	delete(r.underlays, peer.EthAddress)
	delete(r.connections, peerID)
	for _, cancel := range r.streams[peerID] {
		cancel()
	}
	delete(r.streams, peerID)
	return
}

func (r *peerRegistry) getPeer(peerID core.PeerID) (*p2p.Peer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.overlays[peerID]
	return p, ok
}

func (r *peerRegistry) getPeerID(ethAddress common.Address) (core.PeerID, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peerID, ok := r.underlays[ethAddress]
	return peerID, ok
}

func (r *peerRegistry) getPeers() []*p2p.Peer { //nolint:unused
	r.mu.RLock()
	defer r.mu.RUnlock()

	peers := make([]*p2p.Peer, 0, len(r.overlays))
	for _, peer := range r.overlays {
		peers = append(peers, peer)
	}
	return peers
}

func (r *peerRegistry) addStream(
	peerID core.PeerID,
	stream network.Stream,
	cancel context.CancelFunc,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.streams[peerID]; !ok {
		return
	}

	r.streams[peerID][stream] = cancel
}

func (r *peerRegistry) removeStream(peerID core.PeerID, stream network.Stream) {
	r.mu.Lock()
	defer r.mu.Unlock()

	peerStreams, ok := r.streams[peerID]
	if !ok {
		return
	}

	cancel, ok := peerStreams[stream]
	if !ok {
		return
	}

	cancel()
	delete(r.streams[peerID], stream)
}

func (r *peerRegistry) isConnected(peerID core.PeerID) (*p2p.Peer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peer, ok := r.overlays[peerID]
	if !ok {
		return nil, false
	}

	_, ok = r.connections[peerID]
	return peer, ok
}
