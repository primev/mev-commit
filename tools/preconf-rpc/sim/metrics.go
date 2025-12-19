package sim

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	attempts prometheus.Counter
	success  prometheus.Counter
	fail     prometheus.Counter
	latency  prometheus.Histogram
}

func newMetrics() *metrics {
	return &metrics{
		attempts: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "simulator",
			Name:      "attempts_total",
			Help:      "Total number of simulation attempts.",
		}),
		success: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "simulator",
			Name:      "success_total",
			Help:      "Total number of successful simulations.",
		}),
		fail: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "fastrpc",
			Subsystem: "simulator",
			Name:      "fail_total",
			Help:      "Total number of failed simulations.",
		}),
		latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "fastrpc",
			Subsystem: "simulator",
			Name:      "latency_ms",
			Help:      "Histogram of simulation latencies in milliseconds.",
			Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
		}),
	}
}
