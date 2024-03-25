package bidderapi

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "bidder_api"
)

type metrics struct {
	ReceivedBidsCount     prometheus.Counter
	ReceivedPreconfsCount prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		ReceivedBidsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "received_bids_count",
			Help:      "Number of received bids",
		}),
		ReceivedPreconfsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "received_preconfs_count",
			Help:      "Number of received preconfirmations",
		}),
	}
}

func (s *Service) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.ReceivedBidsCount,
		s.metrics.ReceivedPreconfsCount,
	}
}
