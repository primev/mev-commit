package providerapi

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "provider_api"
)

type metrics struct {
	BidsSentToProviderCount     prometheus.Counter
	BidsAcceptedByProviderCount prometheus.Counter
	BidsRejectedByProviderCount prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		BidsSentToProviderCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "bids_sent_to_provider_count",
			Help:      "Number of bids sent to provider",
		}),
		BidsAcceptedByProviderCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "bids_accepted_by_provider_count",
			Help:      "Number of bids accepted by provider",
		}),
		BidsRejectedByProviderCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "bids_rejected_by_provider_count",
			Help:      "Number of bids rejected by provider",
		}),
	}
}

func (s *Service) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.BidsSentToProviderCount,
		s.metrics.BidsAcceptedByProviderCount,
		s.metrics.BidsRejectedByProviderCount,
	}
}
