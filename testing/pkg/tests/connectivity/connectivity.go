package connectivity

import (
	"context"
	"fmt"
	"slices"
	"time"

	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// this is based on the blocklist duration specified in p2p pkg
	connectedTimeout = 6 * time.Minute
)

func Run(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	providers := cluster.Providers()
	bidders := cluster.Bidders()

	logger := cluster.Logger().With("test", "connectivity")

	cctx, cancel := context.WithTimeout(ctx, connectedTimeout)
	defer cancel()

	eg, egCtx := errgroup.WithContext(cctx)

	providerAddrs := make([]string, len(providers))
	for i, p := range providers {
		providerAddrs[i] = p.EthAddress()
	}

	for _, b := range bidders {
		eg.Go(func() error {
			notificationChan, err := b.NotificationsAPI().Subscribe(egCtx, &notificationsapiv1.SubscribeRequest{
				Topics: []string{string(notifications.TopicPeerConnected)},
			})
			if err != nil {
				return fmt.Errorf("failed to subscribe to notifications: %w", err)
			}
			count := 0
			for {
				msg, err := notificationChan.Recv()
				if err != nil {
					return fmt.Errorf("failed to receive notification: %w", err)
				}
				logger.Info("received notification", "notification", msg)
				if msg.Topic != string(notifications.TopicPeerConnected) {
					logger.Info("skipping notification", "topic", msg.Topic)
					continue
				}
				count++
				if count == len(providers) {
					logger.Info("all providers connected to bidder", "bidder", b.EthAddress())
					return nil
				}
				logger.Info("waiting for all providers to connect to bidder", "connected", count, "total", len(providers))
			}
		})
	}

	err := eg.Wait()
	if err != nil {
		logger.Error("failed to get correct notifications for node connectivity", "error", err)
		return fmt.Errorf("failed to get correct notifications for node connectivity: %w", err)
	}

	// ensure correct nodes are connected
	for _, b := range cluster.Bootnodes() {
		l := cluster.Logger().With("bootnode", b.EthAddress())
		topo, err := b.DebugAPI().GetTopology(ctx, &debugapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get topology", "error", err)
			return fmt.Errorf("failed to get topology: %s", err)
		}

		connectedBidders := getBidders(topo)
		if len(connectedBidders) != len(bidders) {
			return fmt.Errorf(
				"bootnode not connected to all bidders: connected %d, total %d",
				len(connectedBidders), len(bidders),
			)
		}

		connectedProviders := getProviders(topo)
		if len(connectedProviders) != len(providers) {
			return fmt.Errorf(
				"bootnode not connected to all providers: connected %d, total %d",
				len(connectedProviders), len(providers),
			)
		}

		for _, p := range connectedProviders {
			if !slices.ContainsFunc(providers, func(p1 orchestrator.Provider) bool {
				return p1.EthAddress() == p.GetStringValue()
			}) {
				l.Error("provider not connected", "provider", p.GetStringValue())
				return fmt.Errorf("provider not connected: %s", p.GetStringValue())
			}
		}

		for _, b := range connectedBidders {
			if !slices.ContainsFunc(bidders, func(b1 orchestrator.Bidder) bool {
				return b1.EthAddress() == b.GetStringValue()
			}) {
				l.Error("bidder not connected", "bidder", b.GetStringValue())
				return fmt.Errorf("bidder not connected: %s", b.GetStringValue())
			}
		}

		logger.Info("all connected to bootnode")
	}

	// check if all bidders are connected to all providers
	for _, b := range bidders {
		l := cluster.Logger().With("bidder", b.EthAddress())

		topo, err := b.DebugAPI().GetTopology(ctx, &debugapiv1.EmptyMessage{})
		if err != nil {
			l.Error("failed to get topology", "error", err)
			return fmt.Errorf("failed to get topology: %s", err)
		}

		connectedProviders := getProviders(topo)
		if len(connectedProviders) != len(providers) {
			l.Error(
				"bidder not connected to all providers",
				"connected", len(connectedProviders), "total", len(providers),
			)
			return fmt.Errorf(
				"bidder not connected to all providers: connected %d, total %d",
				len(connectedProviders), len(providers),
			)
		}

		for _, p := range connectedProviders {
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

func getProviders(topo *debugapiv1.TopologyResponse) []*structpb.Value {
	providerList := topo.Topology.GetFields()["connected_providers"]
	if providerList == nil {
		return nil
	}

	connectedProviders := providerList.GetListValue()
	return connectedProviders.Values
}

func getBidders(topo *debugapiv1.TopologyResponse) []*structpb.Value {
	bidderList := topo.Topology.GetFields()["connected_bidders"]
	if bidderList == nil {
		return nil
	}

	connectedBidders := bidderList.GetListValue()
	return connectedBidders.Values
}
