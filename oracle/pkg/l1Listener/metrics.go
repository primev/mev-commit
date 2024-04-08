package l1Listener

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "l1_listener"
)

type metrics struct {
	WinnerRoundCount *prometheus.CounterVec
	WinnerCount      prometheus.Counter
}

func newMetrics() *metrics {
	m := &metrics{}
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
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.WinnerRoundCount,
		m.WinnerCount,
	}
}
