package evmclient

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultMetricsNamespace = "mev_commit"
)

type metrics struct {
	AttemptedTxCount          prometheus.Counter
	SentTxCount               prometheus.Counter
	SuccessfulTxCount         prometheus.Counter
	CancelledTxCount          prometheus.Counter
	FailedTxCount             prometheus.Counter
	NotFoundDuringCancelCount prometheus.Counter

	LastUsedNonce                  prometheus.Gauge
	LastConfirmedNonce             prometheus.Gauge
	CurrentBlockNumber             prometheus.Gauge
	GetReceiptBatchOperationTimeMs prometheus.Gauge
}

func newMetrics() *metrics {
	m := &metrics{
		AttemptedTxCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "attempted_tx_count",
			Help:      "Number of attempted transactions",
		}),
		SentTxCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "sent_tx_count",
			Help:      "Number of sent transactions",
		}),
		SuccessfulTxCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "successful_tx_count",
			Help:      "Number of successful transactions",
		}),
		CancelledTxCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "cancelled_tx_count",
			Help:      "Number of cancelled transactions",
		}),
		FailedTxCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "failed_tx_count",
			Help:      "Number of failed transactions",
		}),
		NotFoundDuringCancelCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "not_found_during_cancel_count",
			Help:      "Number of transactions not found during cancel",
		}),
		LastUsedNonce: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "last_used_nonce",
			Help:      "Last used nonce",
		}),
		LastConfirmedNonce: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "last_confirmed_nonce",
			Help:      "Last confirmed nonce",
		}),
		CurrentBlockNumber: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "current_block_number",
			Help:      "Current block number at which the node is checking",
		}),
		GetReceiptBatchOperationTimeMs: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultMetricsNamespace,
			Name:      "get_receipt_batch_operation_time_ms",
			Help:      "Time taken to get receipts in a batch",
		}),
	}

	return m
}

func (m *metrics) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.AttemptedTxCount,
		m.SentTxCount,
		m.SuccessfulTxCount,
		m.CancelledTxCount,
		m.FailedTxCount,
		m.NotFoundDuringCancelCount,
		m.LastUsedNonce,
		m.LastConfirmedNonce,
		m.CurrentBlockNumber,
		m.GetReceiptBatchOperationTimeMs,
	}
}
