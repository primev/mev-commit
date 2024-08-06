package validatorapi

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit"
	subsystem        = "validator_api"
)

type metrics struct {
	FetchedValidatorsCount prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		FetchedValidatorsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "fetched_validators_count",
			Help:      "Number of fetched validators requests",
		}),
	}
}

func (s *Service) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.FetchedValidatorsCount,
	}
}
