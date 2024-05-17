package discovery

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	discoverypb "github.com/primev/mev-commit/p2p/gen/go/discovery/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProtocolName    = "discovery"
	ProtocolVersion = "2.0.0"
	checkWorkers    = 10
)

type P2PService interface {
	p2p.Streamer
	Connect(context.Context, []byte) (p2p.Peer, error)
}

type Topology interface {
	AddPeers(...p2p.Peer)
	IsConnected(common.Address) bool
}

type Discovery struct {
	topo       Topology
	streamer   P2PService
	logger     *slog.Logger
	checkPeers chan *discoverypb.PeerInfo
	sem        *semaphore.Weighted
	quit       chan struct{}
}

func New(
	topo Topology,
	streamer P2PService,
	logger *slog.Logger,
) *Discovery {
	d := &Discovery{
		topo:       topo,
		streamer:   streamer,
		logger:     logger.With("protocol", ProtocolName),
		sem:        semaphore.NewWeighted(checkWorkers),
		checkPeers: make(chan *discoverypb.PeerInfo),
		quit:       make(chan struct{}),
	}
	go d.checkAndAddPeers()
	return d
}

func (d *Discovery) peerListStream() p2p.StreamDesc {
	return p2p.StreamDesc{
		Name:    ProtocolName,
		Version: ProtocolVersion,
		Handler: d.handlePeersList,
	}
}

func (d *Discovery) Streams() []p2p.StreamDesc {
	return []p2p.StreamDesc{d.peerListStream()}
}

func (d *Discovery) handlePeersList(ctx context.Context, peer p2p.Peer, s p2p.Stream) error {
	peers := new(discoverypb.PeerList)
	err := s.ReadMsg(ctx, peers)
	if err != nil {
		d.logger.Error("failed to read peers list", "err", err, "from_peer", peer)
		return status.Errorf(codes.InvalidArgument, "failed to read peers list: %v", err)
	}

	for _, p := range peers.Peers {
		if d.topo.IsConnected(common.BytesToAddress(p.EthAddress)) {
			continue
		}
		select {
		case d.checkPeers <- p:
		case <-ctx.Done():
			d.logger.Error("failed to add peer", "err", ctx.Err(), "from_peer", peer)
			return ctx.Err()
		}
	}

	d.logger.Debug("added peers", "peers", len(peers.Peers), "from_peer", peer)
	return nil
}

func (d *Discovery) BroadcastPeers(
	ctx context.Context,
	peer p2p.Peer,
	peers []p2p.PeerInfo,
) error {
	stream, err := d.streamer.NewStream(ctx, peer, nil, d.peerListStream())
	if err != nil {
		d.logger.Error("failed to create stream", "err", err, "to_peer", peer)
		return err
	}
	defer stream.Close()

	peersToSend := make([]*discoverypb.PeerInfo, 0, len(peers))
	for _, p := range peers {
		peersToSend = append(peersToSend, &discoverypb.PeerInfo{
			EthAddress: p.EthAddress.Bytes(),
			Underlay:   p.Underlay,
		})
	}

	if err := stream.WriteMsg(ctx, &discoverypb.PeerList{Peers: peersToSend}); err != nil {
		d.logger.Error("failed to write peers list", "err", err, "to_peer", peer)
		return err
	}

	d.logger.Debug("sent peers list", "peers", len(peers), "to_peer", peer)
	return nil
}

func (d *Discovery) Close() error {
	close(d.quit)
	return nil
}

func (d *Discovery) checkAndAddPeers() {
	for {
		select {
		case <-d.quit:
			return
		case peer := <-d.checkPeers:
			_ = d.sem.Acquire(context.Background(), 1)
			go func() {
				defer d.sem.Release(1)

				p, err := d.streamer.Connect(context.Background(), peer.Underlay)
				if err != nil {
					d.logger.Error("failed to connect to peer", "err", err, "peer", peer)
					return
				}
				d.topo.AddPeers(p)
			}()
		}
	}
}
