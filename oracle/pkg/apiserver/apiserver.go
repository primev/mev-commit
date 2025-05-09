package apiserver

import (
	"bufio"
	"context"
	"expvar"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/health"
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
	blockTracker     *blocktracker.BlocktrackerTransactorSession
	providerRegistry *providerregistry.ProviderregistryCallerSession
	monitor          *txmonitor.Monitor
	shutdown         chan struct{}
}

// New creates a new Service.
func New(
	logger *slog.Logger,
	evm events.EventManager,
	store Store,
	token string,
	blockTracker *blocktracker.BlocktrackerTransactorSession,
	providerRegistry *providerregistry.ProviderregistryCallerSession,
	monitor *txmonitor.Monitor,
) *Service {

	srv := &Service{
		logger:           logger,
		router:           http.NewServeMux(),
		metricsRegistry:  newMetrics(),
		evtMgr:           evm,
		store:            store,
		blockTracker:     blockTracker,
		providerRegistry: providerRegistry,
		monitor:          monitor,
		shutdown:         make(chan struct{}),
	}

	srv.router.Handle("/register_provider", srv.registerProvider(token))

	srv.registerDebugEndpoints()
	return srv
}

func (s *Service) registerProvider(token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Expected format "Bearer <token>"
		headerToken, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		if headerToken != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		provider := r.URL.Query().Get("provider")
		if provider == "" {
			http.Error(w, "Provider not specified", http.StatusBadRequest)
			return
		}

		if !common.IsHexAddress(provider) {
			http.Error(w, "Invalid provider address", http.StatusBadRequest)
			return
		}

		providerAddress := common.HexToAddress(provider)

		grafiti := r.URL.Query().Get("grafiti")
		if grafiti == "" {
			http.Error(w, "Grafiti not specified", http.StatusBadRequest)
			return
		}

		minStake, err := s.providerRegistry.MinStake()
		if err != nil {
			http.Error(w, "Failed to get minimum stake amount", http.StatusInternalServerError)
			return
		}

		stake, err := s.providerRegistry.GetProviderStake(providerAddress)
		if err != nil {
			http.Error(w, "Failed to check provider stake", http.StatusInternalServerError)
			return
		}

		if stake.Cmp(minStake) < 0 {
			http.Error(w, "Insufficient stake", http.StatusBadRequest)
			return
		}

		txn, err := s.blockTracker.AddBuilderAddress(grafiti, providerAddress)
		if err != nil {
			http.Error(w, "Failed to add provider mapping", http.StatusInternalServerError)
			return
		}

		receipt, err := s.monitor.WaitForReceipt(context.Background(), txn)
		if err != nil {
			http.Error(w, "Failed to get receipt for transaction", http.StatusInternalServerError)
			return
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			http.Error(w, "Transaction failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
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
				slog.Int("http_status", recorder.status),
				slog.String("http_method", req.Method),
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

	s.logger.Info("api server started", "addr", addr)

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

func (s *Service) RegisterHealthCheck(hc health.Health) {
	s.router.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			if err := hc.Health(); err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintln(w, "ok")
			if err != nil {
				s.logger.Error(
					"failed to write health check response",
					"error", err,
				)
			}
		},
	)
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
