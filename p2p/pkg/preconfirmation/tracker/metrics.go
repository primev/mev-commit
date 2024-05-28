package preconftracker

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	totalEncryptedCommitments  prometheus.Counter
	totalCommitmentsToOpen     prometheus.Counter
	totalOpenedCommitments     prometheus.Counter
	blockCommitmentProcessTime prometheus.Gauge
}

func newMetrics() *metrics {
	return &metrics{
		totalEncryptedCommitments: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mev_commit",
			Subsystem: "preconftracker",
			Name:      "total_encrypted_commitments",
			Help:      "Total number of encrypted commitments",
		}),
		totalCommitmentsToOpen: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mev_commit",
			Subsystem: "preconftracker",
			Name:      "total_commitments_to_open",
			Help:      "Total number of commitments to open",
		}),
		totalOpenedCommitments: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mev_commit",
			Subsystem: "preconftracker",
			Name:      "total_opened_commitments",
			Help:      "Total number of opened commitments",
		}),
		blockCommitmentProcessTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mev_commit",
			Subsystem: "preconftracker",
			Name:      "block_commitment_process_time",
			Help:      "Time taken to process commitments in a block",
		}),
	}
}

func (m *metrics) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		m.totalEncryptedCommitments,
		m.totalCommitmentsToOpen,
		m.totalOpenedCommitments,
		m.blockCommitmentProcessTime,
	}
}
