package node

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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

	ctx, cancel := context.WithCancel(context.Background())
	err := n.createGatewayContract(
		ctx,
		"l1",
		opts.Logger,
		opts.L1RPCURL,
		opts.Signer,
		opts.L1GatewayContractAddr,
		nil,
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
		nil,
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
			w.Write([]byte(err.Error()))
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
	db *sql.DB,
) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum node: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

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

	st, err := store.NewStore(db, component)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	monitor := txmonitor.New(
		signer.GetAddress(),
		client,
		txmonitor.NewEVMHelperWithLogger(
			client,
			logger.With("component", fmt.Sprintf("%s/evmhelper")),
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

	evtMgr := events.NewListener(
		logger,
		&parsedABI,
	)
	n.metrics.MustRegister(evtMgr.Metrics()...)

	parsedURL, err := url.Parse(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	switch parsedURL.Scheme {
	case "ws", "wss":
		p := publisher.NewWSPublisher(
			st,
			logger,
			client,
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
	case "http", "https":
		p := publisher.NewHTTPPublisher(
			st,
			logger,
			client,
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
	default:
		return fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

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
