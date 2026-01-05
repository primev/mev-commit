package rpcserver

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	methodSuccessCounts         *prometheus.CounterVec
	methodSuccessDurations      *prometheus.HistogramVec
	methodFailureCounts         *prometheus.CounterVec
	methodFailureDurations      *prometheus.HistogramVec
	proxyMethodSuccessCounts    *prometheus.CounterVec
	proxyMethodSuccessDurations *prometheus.HistogramVec
	proxyMethodFailureCounts    *prometheus.CounterVec
	proxyMethodFailureDurations *prometheus.HistogramVec
}

func newMetrics() *metrics {
	return &metrics{
		methodSuccessCounts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "method_success_counts",
				Help:      "Count of successful RPC method calls",
			},
			[]string{"method"},
		),
		methodSuccessDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "method_success_durations_ms",
				Help:      "Duration of successful RPC method calls",
				Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
			},
			[]string{"method"},
		),
		methodFailureCounts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "method_failure_counts",
				Help:      "Count of failed RPC method calls",
			},
			[]string{"method"},
		),
		methodFailureDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "method_failure_durations_ms",
				Help:      "Duration of failed RPC method calls",
				Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
			},
			[]string{"method"},
		),
		proxyMethodSuccessCounts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "rpc_proxy_method_success_counts",
				Help:      "Count of successful proxied RPC method calls",
			},
			[]string{"method"},
		),
		proxyMethodSuccessDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "proxy_method_success_durations_ms",
				Help:      "Duration of successful proxied RPC method calls",
				Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
			},
			[]string{"method"},
		),
		proxyMethodFailureCounts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "proxy_method_failure_counts",
				Help:      "Count of failed proxied RPC method calls",
			},
			[]string{"method"},
		),
		proxyMethodFailureDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "rpc",
				Subsystem: "server",
				Name:      "proxy_method_failure_durations_ms",
				Help:      "Duration of failed proxied RPC method calls",
				Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
			},
			[]string{"method"},
		),
	}
}
