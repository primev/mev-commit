package debugapi

import (
	"log/slog"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primevprotocol/mev-commit/p2p/pkg/apiserver"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p/libp2p"
	"github.com/primevprotocol/mev-commit/p2p/pkg/topology"
)

type APIServer interface {
	ChainHandlers(string, http.Handler, ...func(http.Handler) http.Handler)
}

type Topology interface {
	GetPeers(topology.Query) []p2p.Peer
}

func RegisterAPI(
	srv APIServer,
	topo Topology,
	p2pSvc *libp2p.Service,
	logger *slog.Logger,
) {
	d := &debugapi{
		topo:   topo,
		p2p:    p2pSvc,
		logger: logger,
	}

	srv.ChainHandlers(
		"/topology",
		apiserver.MethodHandler("GET", d.handleTopology),
	)
}

type debugapi struct {
	topo   Topology
	p2p    *libp2p.Service
	logger *slog.Logger
}

type topologyResponse struct {
	Self           map[string]interface{}      `json:"self"`
	ConnectedPeers map[string][]common.Address `json:"connected_peers"`
	BlockedPeers   []p2p.BlockedPeerInfo       `json:"blocked_peers"`
}

func (d *debugapi) handleTopology(w http.ResponseWriter, r *http.Request) {
	logger := d.logger.With("method", "handleTopology")
	providers := d.topo.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	bidders := d.topo.GetPeers(topology.Query{Type: p2p.PeerTypeBidder})

	topoResp := topologyResponse{
		Self:           d.p2p.Self(),
		ConnectedPeers: make(map[string][]common.Address),
	}

	if len(providers) > 0 {
		connectedProviders := make([]common.Address, len(providers))
		for idx, provider := range providers {
			connectedProviders[idx] = provider.EthAddress
		}
		topoResp.ConnectedPeers["providers"] = connectedProviders
	}
	if len(bidders) > 0 {
		connectedBidders := make([]common.Address, len(bidders))
		for idx, bidder := range bidders {
			connectedBidders[idx] = bidder.EthAddress
		}
		topoResp.ConnectedPeers["bidders"] = connectedBidders
	}

	topoResp.BlockedPeers = d.p2p.BlockedPeers()

	err := apiserver.WriteResponse(w, http.StatusOK, topoResp)
	if err != nil {
		logger.Error("error writing response", "err", err)
	}
}
