package validatorapi

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "validator_api"
)

type metrics struct {
	FetchedEpochDataCount   prometheus.Counter
	FetchedValidatorsCount  prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		FetchedEpochDataCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "fetched_epoch_data_count",
			Help:      "Number of fetched epoch data",
		}),
		FetchedValidatorsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "fetched_validators_count",
			Help:      "Number of fetched validators",
		}),
	}
}

func (s *Service) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.FetchedEpochDataCount,
		s.metrics.FetchedValidatorsCount,
	}
}