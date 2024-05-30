package transactor

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
)

var Namespace = "mev_commit"

type metricsWrapper struct {
	bind.ContractBackend

	methodSuccessDurations *prometheus.GaugeVec
	methodErrorDurations   *prometheus.GaugeVec
	methodSuccessCounts    *prometheus.CounterVec
	methodErrorCounts      *prometheus.CounterVec
}

func NewMetricsWrapper(t bind.ContractBackend) *metricsWrapper {
	return &metricsWrapper{
		ContractBackend: t,
		methodSuccessDurations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "transactor",
			Name:      "method_success_durations",
			Help:      "Duration of successful method calls",
		}, []string{"method"}),
		methodErrorDurations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "transactor",
			Name:      "method_error_durations",
			Help:      "Duration of errored method calls",
		}, []string{"method"}),
		methodSuccessCounts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "transactor",
			Name:      "method_success_counts",
			Help:      "Count of successful method calls",
		}, []string{"method"}),
		methodErrorCounts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "transactor",
			Name:      "method_error_counts",
			Help:      "Count of errored method calls",
		}, []string{"method"}),
	}
}

func (m *metricsWrapper) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		m.methodSuccessDurations,
		m.methodErrorDurations,
		m.methodSuccessCounts,
		m.methodErrorCounts,
	}
}

func (m *metricsWrapper) recordMetrics(start time.Time, label string, err error) {
	if err != nil {
		m.methodErrorDurations.WithLabelValues(label).Set(float64(time.Since(start)))
		m.methodErrorCounts.WithLabelValues(label).Inc()
		return
	}

	m.methodSuccessDurations.WithLabelValues(label).Set(float64(time.Since(start)))
	m.methodSuccessCounts.WithLabelValues(label).Inc()
}

func (m *metricsWrapper) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	defer func(start time.Time) { m.recordMetrics(start, "PendingNonceAt", err) }(time.Now())

	return m.ContractBackend.PendingNonceAt(ctx, account)
}

func (m *metricsWrapper) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	defer func(start time.Time) { m.recordMetrics(start, "EstimateGas", err) }(time.Now())

	return m.ContractBackend.EstimateGas(ctx, call)
}

func (m *metricsWrapper) SuggestGasPrice(ctx context.Context) (gas *big.Int, err error) {
	defer func(start time.Time) { m.recordMetrics(start, "SuggestGasPrice", err) }(time.Now())

	return m.ContractBackend.SuggestGasPrice(ctx)
}

func (m *metricsWrapper) SuggestGasTipCap(ctx context.Context) (tip *big.Int, err error) {
	defer func(start time.Time) { m.recordMetrics(start, "SuggestGasTipCap", err) }(time.Now())

	return m.ContractBackend.SuggestGasTipCap(ctx)
}

func (m *metricsWrapper) SendTransaction(ctx context.Context, tx *types.Transaction) (err error) {
	defer func(start time.Time) { m.recordMetrics(start, "SendTransaction", err) }(time.Now())

	return m.ContractBackend.SendTransaction(ctx, tx)
}
