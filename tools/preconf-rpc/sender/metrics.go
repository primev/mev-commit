package sender

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	connectedProviders          prometheus.Gauge
	queuedTransactions          prometheus.Gauge
	inflightTransactions        prometheus.Gauge
	preconfDurationsProvider    *prometheus.GaugeVec
	preconfCountsProvider       *prometheus.CounterVec
	blockAttemptsToConfirmation prometheus.Histogram
	totalAttemptsToConfirmation prometheus.Histogram
	timeToConfirmation          prometheus.Histogram
	timeToFirstPreconfirmation  prometheus.Histogram
	bidPriorityFee              prometheus.Gauge
}

func newMetrics() *metrics {
	return &metrics{
		connectedProviders: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "connected_providers",
			Help:      "Number of currently connected providers.",
		}),
		queuedTransactions: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "queued_transactions",
			Help:      "Number of transactions currently queued for sending.",
		}),
		inflightTransactions: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "inflight_transactions",
			Help:      "Number of transactions currently in-flight.",
		}),
		preconfDurationsProvider: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "preconfirmation_durations_provider_ms",
			Help:      "Duration taken for pre-confirmation by provider in milliseconds.",
		}, []string{"provider", "name"}),
		preconfCountsProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "preconfirmation_counts_provider_total",
			Help:      "Total number of pre-confirmations by provider.",
		}, []string{"provider", "name"}),
		blockAttemptsToConfirmation: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "block_attempts_to_confirmation",
			Help:      "Histogram of block attempts taken to confirm transactions.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}),
		totalAttemptsToConfirmation: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "total_attempts_to_confirmation",
			Help:      "Histogram of total attempts taken to confirm transactions.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}),
		timeToConfirmation: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "time_to_confirmation_ms",
			Help:      "Histogram of time taken to confirm transactions in milliseconds.",
			Buckets:   prometheus.ExponentialBuckets(50, 2, 15),
		}),
		timeToFirstPreconfirmation: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "time_to_first_preconfirmation_ms",
			Help:      "Histogram of time taken to first pre-confirmation in milliseconds.",
			Buckets:   prometheus.ExponentialBuckets(20, 1.7, 15),
		}),
		bidPriorityFee: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "sender",
			Name:      "bid_priority_fee_gwei",
			Help:      "The priority fee (in Gwei) being bid for transactions.",
		}),
	}
}
