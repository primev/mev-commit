package connectivity

import (
	"context"
	"fmt"
	"slices"
	"time"

	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
)

const (
	// this is based on the blocklist duration specified in p2p pkg
	connectedTimeout = 6 * time.Minute
)

func Run(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	providers := cluster.Providers()
	bidders := cluster.Bidders()

	logger := cluster.Logger().With("test", "connectivity")

	start := time.Now()
	bootnodeConnected := false
	for {
		// first check if all nodes are connected to the bootnode
		for _, b := range cluster.Bootnodes() {
			l := cluster.Logger().With("bootnode", b.EthAddress())
			topo, err := b.DebugAPI().GetTopology(ctx, &debugapiv1.EmptyMessage{})
			if err != nil {
				l.Error("failed to get topology", "error", err)
				return fmt.Errorf("failed to get topology: %s", err)
			}

			connectedBidders := topo.Topology.GetFields()["connected_bidders"].GetListValue()
			if len(connectedBidders.Values) != len(bidders) {
				break
			}

			if topo.Topology.GetFields()["connected_providers"] == nil {
				break
			}

			connectedProviders := topo.Topology.GetFields()["connected_providers"].GetListValue()
			if len(connectedProviders.Values) != len(providers) {
				break
			}

			for _, p := range connectedProviders.Values {
				if !slices.ContainsFunc(providers, func(p1 orchestrator.Provider) bool {
					return p1.EthAddress() == p.GetStringValue()
				}) {
					l.Error("provider not connected", "provider", p.GetStringValue())
					return fmt.Errorf("provider not connected: %s", p.GetStringValue())
				}
			}

			for _, b := range connectedBidders.Values {
				if !slices.ContainsFunc(bidders, func(b1 orchestrator.Bidder) bool {
					return b1.EthAddress() == b.GetStringValue()
				}) {
					l.Error("bidder not connected", "bidder", b.GetStringValue())
					return fmt.Errorf("bidder not connected: %s", b.GetStringValue())
				}
			}

			logger.Info("all connected to bootnode")
			bootnodeConnected = true
		}

		if bootnodeConnected {
			break
		}

		if time.Since(start) > connectedTimeout {
			logger.Error("timeout waiting for all nodes to connect to bootnode")
			return fmt.Errorf("timeout waiting for all nodes to connect")
		}
	}

	// check if all bidders are connected to all providers
	for _, b := range bidders {
		l := cluster.Logger().With("bidder", b.EthAddress())

		topo, err := b.DebugAPI().GetTopology(ctx, &debugapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get topology", "error", err)
			return fmt.Errorf("failed to get topology: %s", err)
		}

		connectedProviders := topo.Topology.GetFields()["connected_providers"].GetListValue()
		if len(connectedProviders.Values) != len(providers) {
			l.Error("bidder not connected to all providers")
			return fmt.Errorf("bidder not connected to all providers: %s", b.EthAddress())
		}

		for _, p := range connectedProviders.Values {
			if !slices.ContainsFunc(providers, func(p1 orchestrator.Provider) bool {
				return p1.EthAddress() == p.GetStringValue()
			}) {
				l.Error("bidder connected to unknown provider", "provider", p.GetStringValue())
				return fmt.Errorf("bidder connected to unknown provider: %s", p.GetStringValue())
			}
		}

		l.Info("bidder connected to all providers")
	}

	logger.Info("test passed")

	return nil
}
