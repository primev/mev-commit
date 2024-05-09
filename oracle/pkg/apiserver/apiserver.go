package apiserver

import (
	"bufio"
	"context"
	"expvar"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/primevprotocol/mev-commit/oracle/pkg/updater"
	"github.com/primevprotocol/mev-commit/x/contracts/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultNamespace = "mev_commit_oracle"
)

type Store interface {
	Settlement(context.Context, []byte) (updater.Settlement, error)
}

// Service wraps http.Server with additional functionality for metrics and
// other common middlewares.
type Service struct {
	logger           *slog.Logger
	metricsRegistry  *prometheus.Registry
	router           *http.ServeMux
	srv              *http.Server
	evtMgr           events.EventManager
	store            Store
	statMu           sync.RWMutex
	blockStats       *lru.Cache[uint64, *BlockStats]
	providerStakes   *lru.Cache[string, *ProviderBalances]
	bidderAllowances *lru.Cache[uint64, []*BidderAllowance]
	lastBlock        uint64
	shutdown         chan struct{}
}

// New creates a new Service.
func New(
	logger *slog.Logger,
	evm events.EventManager,
	store Store,
) *Service {
	blockStats, _ := lru.New[uint64, *BlockStats](10000)
	providerStakes, _ := lru.New[string, *ProviderBalances](1000)
	bidderAllowances, _ := lru.New[uint64, []*BidderAllowance](1000)

	srv := &Service{
		logger:           logger,
		router:           http.NewServeMux(),
		metricsRegistry:  newMetrics(),
		evtMgr:           evm,
		store:            store,
		shutdown:         make(chan struct{}),
		blockStats:       blockStats,
		providerStakes:   providerStakes,
		bidderAllowances: bidderAllowances,
	}

	err := srv.configureDashboard()
	if err != nil {
		logger.Error("failed to configure dashboard", "error", err)
	}

	srv.registerDebugEndpoints()
	return srv
}

func (s *Service) registerDebugEndpoints() {
	// register metrics handler
	s.router.Handle("/metrics", promhttp.HandlerFor(s.metricsRegistry, promhttp.HandlerOpts{}))

	// register pprof handlers
	s.router.Handle(
		"/debug/pprof",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Path += "/"
			http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
		}),
	)
	s.router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	s.router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	s.router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	s.router.Handle("/debug/pprof/{profile}", http.HandlerFunc(pprof.Index))
	s.router.Handle("/debug/vars", expvar.Handler())
}

func newMetrics() (r *prometheus.Registry) {
	r = prometheus.NewRegistry()

	// register standard metrics
	r.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			Namespace: defaultNamespace,
		}),
		collectors.NewGoCollector(),
	)

	return r
}

func (s *Service) Start(addr string) <-chan struct{} {
	s.logger.Info("starting api server")

	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			recorder := &responseStatusRecorder{ResponseWriter: w}

			start := time.Now()
			s.router.ServeHTTP(recorder, req)
			s.logger.Info(
				"api access",
				slog.Int("status", recorder.status),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Duration("duration", time.Since(start)),
			)
		}),
	}
	s.srv = srv

	done := make(chan struct{})
	go func() {
		defer close(done)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("api server failed", "error", err)
		}
	}()

	return done
}

func (s *Service) Stop() error {
	s.logger.Info("stopping api server")
	if s.srv == nil {
		return nil
	}
	close(s.shutdown)
	return s.srv.Shutdown(context.Background())
}

// RegisterMetricsCollectors registers prometheus collectors.
func (s *Service) RegisterMetricsCollectors(cs ...prometheus.Collector) {
	s.metricsRegistry.MustRegister(cs...)
}

type responseStatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseStatusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Hijack implements http.Hijacker.
func (r *responseStatusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush implements http.Flusher.
func (r *responseStatusRecorder) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Push implements http.Pusher.
func (r *responseStatusRecorder) Push(target string, opts *http.PushOptions) error {
	return r.ResponseWriter.(http.Pusher).Push(target, opts)
}
