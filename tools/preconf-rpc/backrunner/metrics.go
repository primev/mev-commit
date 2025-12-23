package backrunner

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	attempts     prometheus.Counter
	success      prometheus.Counter
	fail         prometheus.Counter
	rewards      prometheus.Counter
	rewardsTotal prometheus.Gauge
	latency      prometheus.Histogram
}

func newMetrics() *metrics {
	return &metrics{
		attempts: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "attempts_total",
			Help:      "Total number of backrun attempts.",
		}),
		success: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "success_total",
			Help:      "Total number of successful backruns.",
		}),
		fail: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "fail_total",
			Help:      "Total number of failed backruns.",
		}),
		rewards: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "rewards_total",
			Help:      "Total number of backrun rewards collected.",
		}),
		rewardsTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "rewards_amount_wei",
			Help:      "Total amount of backrun rewards collected in WEI.",
		}),
		latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "backrunner",
			Name:      "latency_ms",
			Help:      "Histogram of backrun latencies in milliseconds.",
			Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
		}),
	}
}
