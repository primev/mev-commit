package updater

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "updater"
)

type metrics struct {
	CommitmentsReceivedCount  prometheus.Counter
	CommitmentsProcessedCount prometheus.Counter
	CommitmentsTooOldCount    prometheus.Counter
	DuplicateCommitmentsCount prometheus.Counter
	RewardsCount              prometheus.Counter
	SlashesCount              prometheus.Counter
	EncryptedCommitmentsCount prometheus.Counter
	NoWinnerCount             prometheus.Counter
	BlockTxnCacheHits         prometheus.Counter
	BlockTxnCacheMisses       prometheus.Counter
	BlockTimeCacheHits        prometheus.Counter
	BlockTimeCacheMisses      prometheus.Counter
	LastSentNonce             prometheus.Gauge
}

func newMetrics() *metrics {
	m := &metrics{}
	m.CommitmentsReceivedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "commitments_received_count",
			Help:      "Number of commitments received",
		},
	)
	m.CommitmentsProcessedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "commitments_processed_count",
			Help:      "Number of commitments processed",
		},
	)
	m.CommitmentsTooOldCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "commitments_too_old_count",
			Help:      "Number of commitments that are too old",
		},
	)
	m.DuplicateCommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "duplicate_commitments_count",
			Help:      "Number of duplicate commitments",
		},
	)
	m.RewardsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "rewards_count",
			Help:      "Number of rewards",
		},
	)
	m.SlashesCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "slashes_count",
			Help:      "Number of slashes",
		},
	)
	m.EncryptedCommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "encrypted_commitments_count",
			Help:      "Number of encrypted commitments",
		},
	)
	m.NoWinnerCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "no_winner_count",
			Help:      "Number of times no winner was found",
		},
	)
	m.BlockTxnCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "block_txn_cache_hits",
			Help:      "Number of block txn cache hits",
		},
	)
	m.BlockTxnCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "block_txn_cache_misses",
			Help:      "Number of block txn cache misses",
		},
	)
	m.BlockTimeCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "block_time_cache_hits",
			Help:      "Number of block time cache hits",
		},
	)
	m.BlockTimeCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "block_time_cache_misses",
			Help:      "Number of block time cache misses",
		},
	)
	m.LastSentNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_sent_nonce",
			Help:      "Last nonce sent to for settlement",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.CommitmentsReceivedCount,
		m.CommitmentsProcessedCount,
		m.CommitmentsTooOldCount,
		m.DuplicateCommitmentsCount,
		m.RewardsCount,
		m.SlashesCount,
		m.EncryptedCommitmentsCount,
		m.NoWinnerCount,
		m.BlockTxnCacheHits,
		m.BlockTxnCacheMisses,
		m.BlockTimeCacheHits,
		m.BlockTimeCacheMisses,
		m.LastSentNonce,
	}
}
