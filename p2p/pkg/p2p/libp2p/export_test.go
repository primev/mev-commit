package libp2p

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
)

var (
	NewStream         = newStream
	NewMetadataStream = newMetadataStream
)

func (s *Service) Addrs() ([]byte, error) {
	info := s.host.Peerstore().PeerInfo(s.host.ID())
	return info.MarshalJSON()
}

func (s *Service) Peer() p2p.Peer {
	return p2p.Peer{
		EthAddress: s.ethAddress,
		Type:       s.peerType,
	}
}

func (s *Service) HostID() peer.ID {
	return s.host.ID()
}

func (s *Service) AddrString() string {
	return s.host.Addrs()[0].String() + "/p2p/" + s.host.ID().String()
}

func (s *Service) PeerCount() int {
	return len(s.host.Network().Peers())
}
