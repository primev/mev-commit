package topology

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "topology"
)

type metrics struct {
	ConnectedBiddersCount   prometheus.Gauge
	ConnectedProvidersCount prometheus.Gauge
}

func newMetrics() *metrics {
	return &metrics{
		ConnectedBiddersCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "connected_bidders_count",
			Help:      "Number of connected bidders",
		}),
		ConnectedProvidersCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "connected_providers_count",
			Help:      "Number of connected providers",
		}),
	}
}

func (t *Topology) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		t.metrics.ConnectedBiddersCount,
		t.metrics.ConnectedProvidersCount,
	}
}
