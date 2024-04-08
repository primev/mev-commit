package settler

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "settler"
)

type metrics struct {
	LastConfirmedNonce        prometheus.Gauge
	LastUsedNonce             prometheus.Gauge
	LastConfirmedBlock        prometheus.Gauge
	CurrentSettlementL1Block  prometheus.Gauge
	SettlementsPostedCount    prometheus.Counter
	SettlementsConfirmedCount prometheus.Counter
}

func newMetrics() *metrics {
	m := &metrics{}
	m.LastConfirmedNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_confirmed_nonce",
			Help:      "Last confirmed nonce (L2 block number)",
		},
	)
	m.LastUsedNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_used_nonce",
			Help:      "Last used nonce (L2 block number)",
		},
	)
	m.LastConfirmedBlock = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_confirmed_block",
			Help:      "Last confirmed block (L2 block number)",
		},
	)
	m.CurrentSettlementL1Block = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "current_settlement_l1_block",
			Help:      "Current settlement L1 block",
		},
	)
	m.SettlementsPostedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "settlements_posted_count",
			Help:      "Number of settlement transactions posted",
		},
	)
	m.SettlementsConfirmedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "settlements_confirmed_count",
			Help:      "Number of settlement transactions confirmed",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.LastConfirmedNonce,
		m.LastUsedNonce,
		m.LastConfirmedBlock,
		m.CurrentSettlementL1Block,
		m.SettlementsPostedCount,
		m.SettlementsConfirmedCount,
	}
}
