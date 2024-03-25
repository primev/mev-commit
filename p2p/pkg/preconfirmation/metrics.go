package preconfirmation

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "preconfirmation"
)

type metrics struct {
	SentBidsCount         prometheus.Counter
	ReceivedPreconfsCount prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		SentBidsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "sent_bids_count",
			Help:      "Number of sent bids",
		}),
		ReceivedPreconfsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "received_preconfs_count",
			Help:      "Number of received preconfirmations",
		}),
	}
}

func (p *Preconfirmation) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		p.metrics.SentBidsCount,
		p.metrics.ReceivedPreconfsCount,
	}
}
