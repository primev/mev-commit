package node

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/bridge/standard/pkg/gwcontract"
	"github.com/primev/mev-commit/bridge/standard/pkg/relayer"
	"github.com/primev/mev-commit/bridge/standard/pkg/store"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	settlementgateway "github.com/primev/mev-commit/contracts-abi/clients/SettlementGateway"
	"github.com/primev/mev-commit/x/contracts/ethwrapper"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Options struct {
	Logger                 *slog.Logger
	HTTPPort               int
	Signer                 keysigner.KeySigner
	L1RPCURL               string
	L1GatewayContractAddr  common.Address
	SettlementRPCURL       string
	SettlementContractAddr common.Address
	PgHost                 string
	PgPort                 int
	PgUser                 string
	PgPassword             string
	PgDB                   string
}

type StartableObjWithDesc struct {
	Startable Startable
	Desc      string
}

type Startable interface {
	Start(ctx context.Context) <-chan struct{}
}

type StartableFunc func(ctx context.Context) <-chan struct{}

func (f StartableFunc) Start(ctx context.Context) <-chan struct{} {
	return f(ctx)
}

type Node struct {
	metrics           *prometheus.Registry
	startables        []StartableObjWithDesc
	l1Gateway         relayer.L1Gateway
	settlementGateway relayer.SettlementGateway
	closeFn           func() error
}

func NewNode(opts *Options) (*Node, error) {
	n := &Node{
		metrics: prometheus.NewRegistry(),
	}

	db, err := initDB(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}

	l1Store, err := store.NewStore(db, "l1")
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 store: %w", err)
	}

	settlementStore, err := store.NewStore(db, "settlement")
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement store: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = n.createGatewayContract(
		ctx,
		"l1",
		opts.Logger,
		opts.L1RPCURL,
		opts.Signer,
		opts.L1GatewayContractAddr,
		l1Store,
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create L1 gateway contract: %w", err)
	}

	err = n.createGatewayContract(
		ctx,
		"settlement",
		opts.Logger,
		opts.SettlementRPCURL,
		opts.Signer,
		opts.SettlementContractAddr,
		settlementStore,
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create settlement gateway contract: %w", err)
	}

	r := relayer.NewRelayer(
		opts.Logger.With("component", "relayer"),
		n.l1Gateway,
		n.settlementGateway,
	)
	n.metrics.MustRegister(r.Metrics()...)

	n.startables = append(n.startables, StartableObjWithDesc{Startable: r, Desc: "relayer"})

	h := health.New()
	waitChan := make([]<-chan struct{}, 0, len(n.startables))
	for _, s := range n.startables {
		closeChan := s.Startable.Start(ctx)
		h.Register(health.CloseChannelHealthCheck(s.Desc, closeChan))
		waitChan = append(waitChan, closeChan)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(n.metrics, promhttp.HandlerOpts{}))
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.Health(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/pending_l1_transfers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txns, err := l1Store.PendingTxns()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if err := json.NewEncoder(w).Encode(txns); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/pending_settlement_transfers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txns, err := settlementStore.PendingTxns()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if err := json.NewEncoder(w).Encode(txns); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", opts.HTTPPort),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			opts.Logger.Error("failed to start HTTP server", "error", err)
		}
	}()

	n.closeFn = func() error {
		cancel()
		_ = db.Close()

		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		for _, c := range waitChan {
			select {
			case <-c:
			case <-closeCtx.Done():
				return fmt.Errorf("failed to close in time")
			}
		}
		return server.Shutdown(closeCtx)
	}

	return n, nil
}

func (n *Node) Close() error {
	return n.closeFn()
}

func (n *Node) createGatewayContract(
	ctx context.Context,
	component string,
	logger *slog.Logger,
	rpcURL string,
	signer keysigner.KeySigner,
	contractAddr common.Address,
	st *store.Store,
) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum node: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	setupMetricsNamespace(fmt.Sprintf("bridge_%s", component))

	var contractABI string
	switch component {
	case "l1":
		contractABI = l1gateway.L1gatewayABI
	case "settlement":
		contractABI = settlementgateway.SettlementgatewayABI
	default:
		return fmt.Errorf("unknown component: %s", component)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	wrappedClient, err := ethwrapper.NewClient(
		logger.With("component", fmt.Sprintf("%s/ethwrapper", component)),
		[]string{rpcURL},
		ethwrapper.EthClientWithBlockNumberDrift(2*32),
	)
	if err != nil {
		return fmt.Errorf("failed to create wrapped client: %w", err)
	}

	monitor := txmonitor.New(
		signer.GetAddress(),
		wrappedClient,
		txmonitor.NewEVMHelperWithLogger(
			client,
			logger.With("component", fmt.Sprintf("%s/evmhelper", component)),
			map[common.Address]*abi.ABI{contractAddr: &parsedABI},
		),
		st,
		logger.With("component", fmt.Sprintf("%s/txmonitor", component)),
		1024,
	)

	n.startables = append(
		n.startables,
		StartableObjWithDesc{
			Startable: monitor,
			Desc:      fmt.Sprintf("%s/txmonitor", component),
		},
	)
	n.metrics.MustRegister(monitor.Metrics()...)

	txtor := transactor.NewMetricsWrapper(
		transactor.NewTransactor(
			client,
			monitor,
		),
	)
	n.metrics.MustRegister(txtor.Metrics()...)

	var gatewayTxtor gwcontract.GatewayTransactor
	switch component {
	case "l1":
		gatewayTxtor, err = l1gateway.NewL1gatewayTransactor(contractAddr, txtor)
	case "settlement":
		gatewayTxtor, err = settlementgateway.NewSettlementgatewayTransactor(contractAddr, txtor)
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
	if err != nil {
		return fmt.Errorf("failed to create gateway transactor: %w", err)
	}

	evtMgr := events.NewListener(
		logger,
		&parsedABI,
	)
	n.metrics.MustRegister(evtMgr.Metrics()...)

	p := publisher.NewHTTPPublisher(
		st,
		logger,
		wrappedClient,
		evtMgr,
	)
	n.startables = append(
		n.startables,
		StartableObjWithDesc{
			Startable: StartableFunc(
				func(ctx context.Context) <-chan struct{} {
					return p.Start(ctx, contractAddr)
				},
			),
			Desc: fmt.Sprintf("%s/publisher", component),
		},
	)

	switch component {
	case "l1":
		n.l1Gateway = gwcontract.NewGateway[l1gateway.L1gatewayTransferInitiated](
			logger,
			monitor,
			evtMgr,
			gatewayTxtor,
			func(ctx context.Context) (*bind.TransactOpts, error) {
				return signer.GetAuthWithCtx(ctx, chainID)
			},
			st,
		)
	case "settlement":
		n.settlementGateway = gwcontract.NewGateway[settlementgateway.SettlementgatewayTransferInitiated](
			logger,
			monitor,
			evtMgr,
			gatewayTxtor,
			func(ctx context.Context) (*bind.TransactOpts, error) {
				return signer.GetAuthWithCtx(ctx, chainID)
			},
			st,
		)
	default:
		return fmt.Errorf("unknown component: %s", component)
	}

	return nil
}

func initDB(opts *Options) (db *sql.DB, err error) {
	// Connection string
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.PgHost, opts.PgPort, opts.PgUser, opts.PgPassword, opts.PgDB,
	)

	// Open a connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db, err
}

func setupMetricsNamespace(namespace string) {
	txmonitor.Namespace = namespace
	events.Namespace = namespace
	transactor.Namespace = namespace
}
