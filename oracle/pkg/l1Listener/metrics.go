package l1Listener

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "l1_listener"
)

type metrics struct {
	WinnerPostedCount prometheus.Counter
	WinnerRoundCount  *prometheus.CounterVec
	WinnerCount       prometheus.Counter
	LastSentNonce     prometheus.Gauge
}

func newMetrics() *metrics {
	m := &metrics{}
	m.WinnerPostedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "winner_posted_count",
			Help:      "Number of winners posted",
		},
	)
	m.WinnerRoundCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "winner_round_count",
			Help:      "Number of rounds won per provider",
		},
		[]string{"builder_name"},
	)
	m.WinnerCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "winner_count",
			Help:      "Number of times a provider won",
		},
	)
	m.LastSentNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_sent_nonce",
			Help:      "Last sent nonce",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.WinnerPostedCount,
		m.WinnerRoundCount,
		m.WinnerCount,
		m.LastSentNonce,
	}
}
