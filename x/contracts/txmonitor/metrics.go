package txmonitor

import "github.com/prometheus/client_golang/prometheus"

var Namespace = "mev_commit"

type metrics struct {
	lastUsedNonce      prometheus.Gauge
	lastConfirmedNonce prometheus.Gauge
	lastBlockNumber    prometheus.Gauge
	lastUsedGas        prometheus.Gauge
	lastUsedGasPrice   prometheus.Gauge
	lastUsedGasTip     prometheus.Gauge
}

func newMetrics() *metrics {
	return &metrics{
		lastUsedNonce: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_used_nonce",
			Help:      "Last used nonce",
		}),
		lastConfirmedNonce: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_confirmed_nonce",
			Help:      "Last confirmed nonce",
		}),
		lastBlockNumber: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_block_number",
			Help:      "Last block number",
		}),
		lastUsedGas: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_used_gas",
			Help:      "Last used gas",
		}),
		lastUsedGasPrice: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_used_gas_price",
			Help:      "Last used gas price",
		}),
		lastUsedGasTip: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "txmonitor",
			Name:      "last_used_gas_tip",
			Help:      "Last used gas tip",
		}),
	}
}

func (m *metrics) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		m.lastUsedNonce,
		m.lastConfirmedNonce,
		m.lastBlockNumber,
		m.lastUsedGas,
		m.lastUsedGasPrice,
		m.lastUsedGasTip,
	}
}
