package relayer

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	initiatedTransfers  *prometheus.CounterVec
	finalizedTransfers  *prometheus.CounterVec
	failedFinalizations *prometheus.CounterVec
}

func newMetrics() *metrics {
	return &metrics{
		initiatedTransfers: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "bridge_relayer",
			Name:      "initiated_transfers",
			Help:      "Number of initiated transfers",
		}, []string{"gateway"}),
		finalizedTransfers: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "bridge_relayer",
			Name:      "finalized_transfers",
			Help:      "Number of finalized transfers",
		}, []string{"gateway"}),
		failedFinalizations: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "bridge_relayer",
			Name:      "failed_finalizations",
			Help:      "Number of failed finalizations",
		}, []string{"gateway"}),
	}
}
