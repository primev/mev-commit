package updater

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "updater"
)

type metrics struct {
	UpdaterTriggerCount   prometheus.Counter
	CommitmentsCount      prometheus.Counter
	RewardsCount          prometheus.Counter
	SlashesCount          prometheus.Counter
	BlockCommitmentsCount prometheus.Counter
}

func newMetrics() *metrics {
	m := &metrics{}
	m.UpdaterTriggerCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "updater_trigger_count",
			Help:      "Number of times the updater was triggered",
		},
	)
	m.CommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "commitments_count",
			Help:      "Number of commitments processed",
		},
	)
	m.RewardsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "rewards_count",
			Help:      "Number of rewards",
		},
	)
	m.SlashesCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "slashes_count",
			Help:      "Number of slashes",
		},
	)
	m.BlockCommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "block_commitments_count",
			Help:      "Number of blocks for which commitments were processed",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.UpdaterTriggerCount,
		m.CommitmentsCount,
		m.RewardsCount,
		m.SlashesCount,
		m.BlockCommitmentsCount,
	}
}
