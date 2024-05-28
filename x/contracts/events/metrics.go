package events

import "github.com/prometheus/client_golang/prometheus"

var Namespace = "mev_commit"

type metrics struct {
	totalLogs             prometheus.Counter
	totalEvents           prometheus.Counter
	eventHandlerDurations *prometheus.GaugeVec
	eventCounts           *prometheus.CounterVec
}

func newMetrics() *metrics {
	return &metrics{
		totalLogs: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "events",
			Name:      "total_logs",
			Help:      "Total number of logs",
		}),
		totalEvents: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "events",
			Name:      "total_events",
			Help:      "Total number of events",
		}),
		eventHandlerDurations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "events",
			Name:      "event_handler_durations",
			Help:      "Duration of event handler",
		}, []string{"event_name"}),
		eventCounts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "events",
			Name:      "event_counts",
			Help:      "Count of events",
		}, []string{"event_name"}),
	}
}

func (m *metrics) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		m.totalLogs,
		m.totalEvents,
		m.eventHandlerDurations,
		m.eventCounts,
	}
}
