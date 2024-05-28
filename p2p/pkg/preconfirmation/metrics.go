package preconfirmation

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "preconfirmation"
)

type metrics struct {
	SentBidsCount                   prometheus.Counter
	ReceivedPreconfsCount           prometheus.Counter
	ConstructPreconfDurationSummary prometheus.Summary
	VerifyPreconfDurationSummary    prometheus.Summary
	BidConstructDurationSummary     prometheus.Summary
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
		ConstructPreconfDurationSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  defaultNamespace,
			Subsystem:  subsystem,
			Name:       "encrypted_preconfirmation_construct_duration_seconds",
			Help:       "Duration taken to construct encrypted preconfirmation in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		VerifyPreconfDurationSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  defaultNamespace,
			Subsystem:  subsystem,
			Name:       "encrypted_preconfirmation_verify_duration_seconds",
			Help:       "Duration taken to verify encrypted preconfirmation in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		BidConstructDurationSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  defaultNamespace,
			Subsystem:  subsystem,
			Name:       "encrypted_bid_construct_duration_seconds",
			Help:       "Duration taken to construct encrypted bid in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
}

func (p *Preconfirmation) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		p.metrics.SentBidsCount,
		p.metrics.ReceivedPreconfsCount,
		p.metrics.ConstructPreconfDurationSummary,
		p.metrics.VerifyPreconfDurationSummary,
		p.metrics.BidConstructDurationSummary,
	}
}
