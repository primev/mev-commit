package libp2p

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
)

func getPeerIDFromMultiaddr(maddr multiaddr.Multiaddr) (peer.ID, error) {
	// Split the multiaddress into its components
	components := maddr.Protocols()

	// Iterate over the components to find the p2p part which contains the peer ID
	for _, value := range components {
		if value.Code == multiaddr.P_P2P {
			// Extract the peer ID part of the multiaddress
			peerIDStr, err := maddr.ValueForProtocol(value.Code)
			if err != nil {
				return "", err
			}
			// Convert the extracted peer ID string into a peer.ID
			pid, err := peer.Decode(peerIDStr)
			if err != nil {
				return "", err
			}
			return pid, nil
		}
	}

	return "", fmt.Errorf("no peer ID found in multiaddress")
}

func (p *Service) startBootstrapper(addrs []string) {
	for {
		for _, addr := range addrs {
			mAddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				p.logger.Error("failed to parse bootstrap address", "addr", addr, "err", err)
				continue
			}

			var addrInfo *peer.AddrInfo

			proto, _ := multiaddr.SplitFirst(mAddr)
			if proto.Protocol().Code == multiaddr.P_DNSADDR {
				resolvedAddrs, err := madns.DefaultResolver.Resolve(context.Background(), mAddr)
				if err != nil {
					p.logger.Error("failed to resolve bootstrap address", "addr", addr, "err", err)
					continue
				}
				if len(resolvedAddrs) == 0 {
					p.logger.Error("failed to resolve bootstrap address", "addr", addr, "err", "no addresses found")
					continue
				}
				pID, err := getPeerIDFromMultiaddr(resolvedAddrs[0])
				if err != nil {
					p.logger.Error("failed to parse peer ID", "addr", addr, "err", err)
					continue
				}
				p.logger.Debug("resolved bootstrap address", "addr", addr, "addrs", resolvedAddrs, "peerID", pID.String())
				addrInfo = &peer.AddrInfo{
					ID:    pID,
					Addrs: resolvedAddrs,
				}
			} else {
				addrInfo, err = peer.AddrInfoFromString(addr)
				if err != nil {
					p.logger.Error("failed to parse bootstrap address", "addr", addr, "err", err)
					continue
				}
			}

			if _, connected := p.peers.isConnected(addrInfo.ID); connected {
				p.logger.Debug("already connected to bootstrap peer", "peer", addrInfo.ID)
				continue
			}

			addrInfoBytes, err := addrInfo.MarshalJSON()
			if err != nil {
				p.logger.Error("failed to marshal bootstrap peer", "addr", addr, "err", err)
				continue
			}

			peer, err := p.Connect(context.Background(), addrInfoBytes)
			if err != nil {
				p.logger.Error("failed to connect to bootstrap peer", "addr", addr, "err", err)
				continue
			}

			p.logger.Info("connected to bootstrap peer", "peer", peer)
		}
		time.Sleep(1 * time.Minute)
	}
}
