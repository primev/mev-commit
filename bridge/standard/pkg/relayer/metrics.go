package relayer

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	initiatedTransfers       *prometheus.CounterVec
	finalizedTransfers       *prometheus.CounterVec
	failedFinalizations      *prometheus.CounterVec
	initiatedTransfersValue  *prometheus.CounterVec
	finalizedTransfersValue  *prometheus.CounterVec
	failedFinalizationsValue *prometheus.CounterVec
}

func newMetrics() *metrics {
	return &metrics{
		initiatedTransfers: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "initiated_transfers",
			Help:      "Number of initiated transfers",
		}, []string{"gateway"}),
		finalizedTransfers: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "finalized_transfers",
			Help:      "Number of finalized transfers",
		}, []string{"gateway"}),
		failedFinalizations: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "failed_finalizations",
			Help:      "Number of failed finalizations",
		}, []string{"gateway"}),
		initiatedTransfersValue: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "initiated_transfers_value",
			Help:      "Value of initiated transfers",
		}, []string{"gateway"}),
		finalizedTransfersValue: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "finalized_transfers_value",
			Help:      "Value of finalized transfers",
		}, []string{"gateway"}),
		failedFinalizationsValue: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "relayer",
			Name:      "failed_finalizations_value",
			Help:      "Value of failed finalizations",
		}, []string{"gateway"}),
	}
}
