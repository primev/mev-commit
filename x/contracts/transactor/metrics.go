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

type metricsWrapper struct {
	bind.ContractTransactor

	methodSuccessDurations *prometheus.GaugeVec
	methodErrorDurations   *prometheus.GaugeVec
	methodSuccessCounts    *prometheus.CounterVec
	methodErrorCounts      *prometheus.CounterVec
}

func NewMetricsWrapper(t bind.ContractTransactor) *metricsWrapper {
	return &metricsWrapper{
		ContractTransactor: t,
		methodSuccessDurations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mev_commit",
			Subsystem: "transactor",
			Name:      "method_success_durations",
			Help:      "Duration of successful method calls",
		}, []string{"method"}),
		methodErrorDurations: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mev_commit",
			Subsystem: "transactor",
			Name:      "method_error_durations",
			Help:      "Duration of errored method calls",
		}, []string{"method"}),
		methodSuccessCounts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "mev_commit",
			Subsystem: "transactor",
			Name:      "method_success_counts",
			Help:      "Count of successful method calls",
		}, []string{"method"}),
		methodErrorCounts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "mev_commit",
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

func (m *metricsWrapper) recordTime(start time.Time, label string, err error) {
	if err != nil {
		m.methodErrorDurations.WithLabelValues(label).Set(float64(time.Since(start)))
		m.methodErrorCounts.WithLabelValues(label).Inc()
		return
	}

	m.methodSuccessDurations.WithLabelValues(label).Set(float64(time.Since(start)))
	m.methodSuccessCounts.WithLabelValues(label).Inc()
}

func (m *metricsWrapper) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	defer func(start time.Time) { m.recordTime(start, "PendingNonceAt", err) }(time.Now())

	return m.ContractTransactor.PendingNonceAt(ctx, account)
}

func (m *metricsWrapper) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	defer func(start time.Time) { m.recordTime(start, "EstimateGas", err) }(time.Now())

	return m.ContractTransactor.EstimateGas(ctx, call)
}

func (m *metricsWrapper) SuggestGasPrice(ctx context.Context) (gas *big.Int, err error) {
	defer func(start time.Time) { m.recordTime(start, "SuggestGasPrice", err) }(time.Now())

	return m.ContractTransactor.SuggestGasPrice(ctx)
}

func (m *metricsWrapper) SuggestGasTipCap(ctx context.Context) (tip *big.Int, err error) {
	defer func(start time.Time) { m.recordTime(start, "SuggestGasTipCap", err) }(time.Now())

	return m.ContractTransactor.SuggestGasTipCap(ctx)
}

func (m *metricsWrapper) SendTransaction(ctx context.Context, tx *types.Transaction) (err error) {
	defer func(start time.Time) { m.recordTime(start, "SendTransaction", err) }(time.Now())

	return m.ContractTransactor.SendTransaction(ctx, tx)
}
