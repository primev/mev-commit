package topology

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
)

type Query struct {
	Type p2p.PeerType
}

type Announcer interface {
	BroadcastPeers(context.Context, p2p.Peer, []p2p.PeerInfo) error
}

type Topology struct {
	mu          sync.RWMutex
	providers   map[common.Address]p2p.Peer
	bidders     map[common.Address]p2p.Peer
	logger      *slog.Logger
	addressbook p2p.Addressbook
	announcer   Announcer
	metrics     *metrics
	subs        []func(p2p.Peer)
}

func (t *Topology) SubscribePeer(handler func(p2p.Peer)) {
	t.subs = append(t.subs, handler)
}

func New(a p2p.Addressbook, logger *slog.Logger) *Topology {
	return &Topology{
		providers:   make(map[common.Address]p2p.Peer),
		bidders:     make(map[common.Address]p2p.Peer),
		addressbook: a,
		logger:      logger,
		metrics:     newMetrics(),
	}
}

func (t *Topology) SetAnnouncer(a Announcer) {
	t.announcer = a
}

func (t *Topology) Connected(p p2p.Peer) {
	t.add(p)

	if t.announcer != nil {
		// Whether its a provider or bidder, we want to broadcast the provider peers
		peersToBroadcast := t.GetPeers(Query{Type: p2p.PeerTypeProvider})
		var underlays []p2p.PeerInfo
		for _, peer := range peersToBroadcast {
			if peer.EthAddress == p.EthAddress {
				continue
			}
			u, err := t.addressbook.GetPeerInfo(peer)
			if err != nil {
				t.logger.Error("failed to get peer info", "err", err, "peer", peer)
				continue
			}
			underlays = append(underlays, p2p.PeerInfo{
				EthAddress: peer.EthAddress,
				Underlay:   u,
			})
		}

		if len(underlays) > 0 {
			err := t.announcer.BroadcastPeers(context.Background(), p, underlays)
			if err != nil {
				t.logger.Error("failed to broadcast peers", "err", err, "peer", p)
			}
		}

		if p.Type == p2p.PeerTypeProvider {
			t.logger.Info("provider connected broadcasting to previous bidders", "peer", p)
			// If the peer is a provider, we want to broadcast to the bidder peers
			peersToBroadcastTo := t.GetPeers(Query{Type: p2p.PeerTypeBidder})
			providerUnderlay, err := t.addressbook.GetPeerInfo(p)
			if err != nil {
				t.logger.Error("failed to get peer info", "err", err, "peer", p)
				return
			}
			for _, peer := range peersToBroadcastTo {
				err := t.announcer.BroadcastPeers(context.Background(), peer, []p2p.PeerInfo{
					{
						EthAddress: p.EthAddress,
						Underlay:   providerUnderlay,
					},
				})
				if err != nil {
					t.logger.Error("failed to broadcast peer", "err", err, "peer", peer)
				}
			}
		}
	}
}

func (t *Topology) add(p p2p.Peer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch p.Type {
	case p2p.PeerTypeProvider:
		t.providers[p.EthAddress] = p
		t.metrics.ConnectedProvidersCount.Inc()
	case p2p.PeerTypeBidder:
		t.bidders[p.EthAddress] = p
		t.metrics.ConnectedBiddersCount.Inc()
	}
}

func (t *Topology) Disconnected(p p2p.Peer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.logger.Info("disconnected", "peer", p)

	switch p.Type {
	case p2p.PeerTypeProvider:
		delete(t.providers, p.EthAddress)
		t.metrics.ConnectedProvidersCount.Dec()
	case p2p.PeerTypeBidder:
		delete(t.bidders, p.EthAddress)
		t.metrics.ConnectedBiddersCount.Dec()
	}
}

func (t *Topology) AddPeers(peers ...p2p.Peer) {
	for _, p := range peers {
		t.add(p)
		for _, sub := range t.subs {
			sub(p)
		}
	}
}

func (t *Topology) GetPeers(q Query) []p2p.Peer {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var peers []p2p.Peer

	switch q.Type {
	case p2p.PeerTypeProvider:
		for _, p := range t.providers {
			peers = append(peers, p)
		}
	case p2p.PeerTypeBidder:
		for _, p := range t.bidders {
			peers = append(peers, p)
		}
	}

	return peers
}

func (t *Topology) IsConnected(addr common.Address) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.providers[addr]; ok {
		return true
	}

	if _, ok := t.bidders[addr]; ok {
		return true
	}

	return false
}
