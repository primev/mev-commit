package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/backrunner"
	bidder "github.com/primev/mev-commit/tools/preconf-rpc/bidder"
	"github.com/primev/mev-commit/tools/preconf-rpc/blocktracker"
	"github.com/primev/mev-commit/tools/preconf-rpc/handlers"
	"github.com/primev/mev-commit/tools/preconf-rpc/notifier"
	"github.com/primev/mev-commit/tools/preconf-rpc/points"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	tracker "github.com/primev/mev-commit/tools/preconf-rpc/settlement-tracker"
	"github.com/primev/mev-commit/tools/preconf-rpc/sim"
	"github.com/primev/mev-commit/tools/preconf-rpc/store"
	"github.com/primev/mev-commit/x/accountsync"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/transfer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultMaxBodySize = 1 * 1024 * 1024 // 1 MB

type Config struct {
	Logger                 *slog.Logger
	PgHost                 string
	PgPort                 int
	PgUser                 string
	PgPassword             string
	PgDbname               string
	PgSSL                  bool
	Signer                 keysigner.KeySigner
	BidderRPC              string
	TargetDepositAmount    *big.Int
	L1RPCHTTPUrl           string
	L1RPCWSUrl             string
	SettlementRPCUrl       string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
	DepositAddress         common.Address
	BridgeAddress          common.Address
	SettlementThreshold    *big.Int
	SettlementTopup        *big.Int
	BidderThreshold        *big.Int
	BidderTopup            *big.Int
	HTTPPort               int
	GasTipCap              *big.Int
	GasFeeCap              *big.Int
	PricerAPIKey           string
	Webhooks               []string
	Token                  string
	SimulatorURL           string
	BackrunnerRPC          string
	BackrunnerAPIURL       string
	BackrunnerAPIKey       string
	PointsAPIURL           string
}

type Service struct {
	cancel  context.CancelFunc
	closers []io.Closer
}

func New(config *Config) (*Service, error) {
	s := &Service{}

	conn, err := grpc.NewClient(
		config.BidderRPC,
		grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{InsecureSkipVerify: true},
		)),
	)
	if err != nil {
		return nil, err
	}

	s.closers = append(s.closers, conn)

	l1RPCClient, err := ethclient.DialContext(context.Background(), config.L1RPCWSUrl)
	if err != nil {
		return nil, err
	}
	settlementClient, err := ethclient.Dial(config.SettlementRPCUrl)
	if err != nil {
		return nil, err
	}

	l1ChainID, err := l1RPCClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get L1 chain ID: %w", err)
	}

	settlementChainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement chain ID: %w", err)
	}

	bidderCli := bidderapiv1.NewBidderClient(conn)
	topologyCli := debugapiv1.NewDebugServiceClient(conn)
	notificationsCli := notificationsapiv1.NewNotificationsClient(conn)

	if err := setupDeposits(bidderCli, config.TargetDepositAmount); err != nil {
		return nil, fmt.Errorf("failed to setup deposits: %w", err)
	}

	notifier := notifier.NewNotifier(config.Webhooks, config.Logger.With("module", "notifier"))

	bridgeConfig := transfer.BridgeConfig{
		Signer:                 config.Signer,
		L1ContractAddr:         config.L1ContractAddr,
		SettlementContractAddr: config.SettlementContractAddr,
		L1RPCUrl:               config.L1RPCHTTPUrl,
		SettlementRPCUrl:       config.SettlementRPCUrl,
	}

	syncer := accountsync.NewAccountSync(config.Signer.GetAddress(), settlementClient)
	bridger := transfer.NewBridger(
		config.Logger.With("module", "bridger"),
		syncer,
		bridgeConfig,
		config.SettlementThreshold,
		config.SettlementTopup,
	)

	bidderClient := bidder.NewBidderClient(
		config.Logger.With("module", "bidder"),
		bidderCli,
		topologyCli,
		notificationsCli,
		l1RPCClient,
	)

	transferer := transfer.NewTransferer(
		config.Logger.With("module", "transferer"),
		settlementClient,
		config.Signer,
		config.GasTipCap,
		config.GasFeeCap,
	)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	healthChecker := health.New()

	balanceNotifierDone := notifier.SetupLowBalanceNotification(
		ctx,
		"RPC Operator AccountBalance Low",
		l1RPCClient,
		config.Signer.GetAddress(),
		3.0,
		5*time.Minute,
		15*time.Minute,
	)

	healthChecker.Register(health.CloseChannelHealthCheck("BalanceNotifier", balanceNotifierDone))
	s.closers = append(s.closers, channelCloser(balanceNotifierDone))

	bidderEOA, err := getBidderEOA(topologyCli)
	if err != nil {
		return nil, fmt.Errorf("failed to get bidder EOA: %w", err)
	}

	config.Logger.Info("bidder EOA", "address", bidderEOA.Hex())

	bidderFunderDone := startBidderFunder(
		ctx,
		config.Logger.With("module", "bidderfunder"),
		bidderEOA,
		accountsync.NewAccountSync(bidderEOA, settlementClient),
		transferer,
		config.BidderThreshold,
		config.BidderTopup,
		settlementClient,
		settlementChainID,
		notifier,
	)

	txnNotifierDone := notifier.StartTransactionNotifier(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("TransactionNotifier", txnNotifierDone))
	s.closers = append(s.closers, channelCloser(txnNotifierDone))

	healthChecker.Register(health.CloseChannelHealthCheck("BidderFunder", bidderFunderDone))
	s.closers = append(s.closers, channelCloser(bidderFunderDone))

	bridgerDone := bridger.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("Bridger", bridgerDone))
	s.closers = append(s.closers, channelCloser(bridgerDone))

	bidderDone := bidderClient.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))
	s.closers = append(s.closers, channelCloser(bidderDone))

	rpcServer, err := rpcserver.NewJSONRPCServer(
		config.L1RPCHTTPUrl,
		config.Logger.With("module", "rpcserver"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC server: %w", err)
	}

	metricsRegistry := prometheus.NewRegistry()

	bidpricer, err := pricer.NewPricer(config.PricerAPIKey, config.Logger.With("module", "bidpricer"))
	if err != nil {
		return nil, fmt.Errorf("failed to create bid pricer: %w", err)
	}
	metricsRegistry.MustRegister(bidpricer.Metrics()...)

	db, err := initDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	rpcstore, err := store.New(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	blockTracker, err := blocktracker.NewBlockTracker(
		l1RPCClient,
		config.Logger.With("module", "blocktracker"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create block tracker: %w", err)
	}

	blockTrackerDone := blockTracker.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("BlockTracker", blockTrackerDone))
	s.closers = append(s.closers, channelCloser(blockTrackerDone))

	simulator := sim.NewSimulator(config.SimulatorURL)
	metricsRegistry.MustRegister(simulator.Metrics()...)

	pointsTracker := points.NewPointsTracker(
		config.PointsAPIURL,
		config.Logger.With("module", "pointstracker"),
	)

	brunner, err := backrunner.New(
		config.BackrunnerAPIKey,
		config.BackrunnerAPIURL,
		config.BackrunnerRPC,
		rpcstore,
		pointsTracker,
		config.Logger.With("module", "backrunner"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create backrunner: %w", err)
	}
	backrunnerDone := brunner.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("Backrunner", backrunnerDone))
	s.closers = append(s.closers, channelCloser(backrunnerDone))
	metricsRegistry.MustRegister(brunner.Metrics()...)

	sndr, err := sender.NewTxSender(
		rpcstore,
		bidderClient,
		bidpricer,
		blockTracker,
		transferer,
		notifier,
		simulator,
		brunner,
		settlementChainID,
		config.Logger.With("module", "txsender"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction sender: %w", err)
	}

	senderDone := sndr.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("TxSender", senderDone))
	s.closers = append(s.closers, channelCloser(senderDone))
	metricsRegistry.MustRegister(sndr.Metrics()...)

	settlementTracker := tracker.NewTracker(
		bidderClient,
		rpcstore,
		config.Logger.With("module", "settlementtracker"),
	)
	settlementTrackerDone := settlementTracker.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("SettlementTracker", settlementTrackerDone))
	s.closers = append(s.closers, channelCloser(settlementTrackerDone))

	rpcHandlers := handlers.NewRPCMethodHandler(
		config.Logger.With("module", "handlers"),
		bidpricer,
		bidderClient,
		rpcstore,
		blockTracker,
		sndr,
		config.DepositAddress,
		config.BridgeAddress,
		l1ChainID,
	)

	rpcHandlers.RegisterMethods(rpcServer)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := healthChecker.Health(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.Handle("/metrics", promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/{option...}", func(w http.ResponseWriter, r *http.Request) {
		options := r.PathValue("option")

		if options != "" {
			splits := strings.Split(options, "/")
			if len(splits) != 3 {
				http.Error(w, "invalid position constraint format", http.StatusBadRequest)
				return
			}
			constraint := new(bidderapiv1.PositionConstraint)
			switch splits[0] {
			case "top":
				constraint.Anchor = bidderapiv1.PositionConstraint_ANCHOR_TOP
			case "bottom":
				constraint.Anchor = bidderapiv1.PositionConstraint_ANCHOR_BOTTOM
			default:
				http.Error(w, "invalid position constraint", http.StatusBadRequest)
				return
			}

			switch splits[1] {
			case "absolute":
				constraint.Basis = bidderapiv1.PositionConstraint_BASIS_ABSOLUTE
			case "percentile":
				constraint.Basis = bidderapiv1.PositionConstraint_BASIS_PERCENTILE
			case "gas_percentile":
				constraint.Basis = bidderapiv1.PositionConstraint_BASIS_GAS_PERCENTILE
			default:
				http.Error(w, "invalid position constraint type", http.StatusBadRequest)
				return
			}

			value, err := strconv.Atoi(splits[2])
			if err != nil {
				http.Error(w, "invalid position constraint value", http.StatusBadRequest)
				return
			}
			constraint.Value = int32(value)

			r = r.WithContext(handlers.SetPositionConstraint(r.Context(), constraint))
		}
		rpcServer.ServeHTTP(w, r)
	})

	registerAdminAPIs(mux, config.Token, sndr, rpcstore)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPPort),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			config.Logger.Error("failed to start HTTP server", "error", err)
		}
	}()

	s.closers = append(s.closers, &srv)

	return s, nil
}

type RPCStore interface {
	AddSubsidy(ctx context.Context, account common.Address, amount *big.Int) error
	GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*sender.Transaction, error)
	GetUserTransactions(ctx context.Context, user common.Address) (store.UserTxnsResponse, error)
}

func registerAdminAPIs(mux *http.ServeMux, token string, sndr *sender.TxSender, rpcstore RPCStore) {
	checkAuthorization := func(r *http.Request) error {
		if token == "" {
			return errors.New("server not configured with authorization token")
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return errors.New("authorization header missing")
		}

		// Expected format "Bearer <token>"
		headerToken, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			return errors.New("invalid authorization header format")
		}

		if headerToken != token {
			return errors.New("unauthorized: invalid token")
		}

		return nil
	}

	var fastTrackMutex sync.RWMutex
	allSlots := false
	providers := []common.Address{}

	fastTrackFn := func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool {
		fastTrackMutex.RLock()
		defer fastTrackMutex.RUnlock()

		if !allSlots && !optedInSlot {
			return false
		}
		if len(providers) == 0 {
			return false
		}
		for _, p := range providers {
			if !slices.ContainsFunc(cmts, func(cmt *bidderapiv1.Commitment) bool {
				return common.HexToAddress(cmt.ProviderAddress).Cmp(p) == 0
			}) {
				return false
			}
		}
		return true
	}

	sndr.SetFastTrackFunc(fastTrackFn)

	mux.HandleFunc("POST /fast-track/enable", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		type fastTrackReq struct {
			AllSlots  bool
			Providers []string
		}

		var req fastTrackReq

		r.Body = http.MaxBytesReader(w, r.Body, defaultMaxBodySize)
		defer func() {
			_ = r.Body.Close()
		}()

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		fastTrackMutex.Lock()
		defer fastTrackMutex.Unlock()

		allSlots = req.AllSlots
		providers = make([]common.Address, 0, len(req.Providers))
		for _, p := range req.Providers {
			if !common.IsHexAddress(p) {
				http.Error(w, fmt.Sprintf("invalid provider address: %s", p), http.StatusBadRequest)
				return
			}
			if slices.ContainsFunc(providers, func(addr common.Address) bool {
				return addr.Cmp(common.HexToAddress(p)) == 0
			}) {
				continue
			}
			providers = append(providers, common.HexToAddress(p))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /fast-track/disable", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		allSlots = false
		providers = []common.Address{}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /fast-track/status", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		type fastTrackStatus struct {
			AllSlots  bool
			Providers []string
		}
		resp := fastTrackStatus{
			AllSlots:  allSlots,
			Providers: make([]string, 0, len(providers)),
		}
		for _, p := range providers {
			resp.Providers = append(resp.Providers, p.Hex())
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /subsidize", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Get account address and amount from URL params
		account := r.URL.Query().Get("account")
		if account == "" || !common.IsHexAddress(account) {
			http.Error(w, "invalid or missing account address", http.StatusBadRequest)
			return
		}

		amountStr := r.URL.Query().Get("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		if err := rpcstore.AddSubsidy(r.Context(), common.HexToAddress(account), amount); err != nil {
			http.Error(w, fmt.Sprintf("failed to add subsidy: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /update_bid_timeout", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		timeoutStr := r.URL.Query().Get("timeout")
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			http.Error(w, "invalid timeout", http.StatusBadRequest)
			return
		}
		sndr.UpdateBidTimeout(timeout)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /status/{txnHash}", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		txnHash := common.HexToHash(r.PathValue("txnHash"))

		txn, err := rpcstore.GetTransactionByHash(r.Context(), txnHash)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get transaction: %v", err), http.StatusNotFound)
			return
		}

		status := struct {
			Hash    string
			Status  string
			Details string
		}{
			Hash:    txn.Hash().Hex(),
			Status:  string(txn.Status),
			Details: txn.Details,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /user-transactions", func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userAddress := r.URL.Query().Get("address")
		if userAddress == "" || !common.IsHexAddress(userAddress) {
			http.Error(w, "invalid or missing user address", http.StatusBadRequest)
			return
		}

		txnList, err := rpcstore.GetUserTransactions(r.Context(), common.HexToAddress(userAddress))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user transactions: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(txnList); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) Close() error {
	s.cancel()

	for _, c := range s.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

type channelCloser <-chan struct{}

func (c channelCloser) Close() error {
	select {
	case <-c:
	case <-time.After(5 * time.Second):
		return errors.New("timed out waiting for channel to close")
	}
	return nil
}

func initDB(opts *Config) (db *sql.DB, err error) {
	// Connection string
	sslMode := "disable"
	if opts.PgSSL {
		sslMode = "require"
	}
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		opts.PgHost, opts.PgPort, opts.PgUser, opts.PgPassword, opts.PgDbname, sslMode,
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

func startBidderFunder(
	ctx context.Context,
	logger *slog.Logger,
	bidderAccount common.Address,
	syncer *accountsync.AccountSync,
	transferer *transfer.Transferer,
	settlementThreshold *big.Int,
	settlementTopup *big.Int,
	settlementClient *ethclient.Client,
	settlementChainID *big.Int,
	notifier *notifier.Notifier,
) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			sub := syncer.Subscribe(ctx, settlementThreshold)
			select {
			case <-ctx.Done():
				return
			case <-sub:
				logger.Info("bidder account balance below threshold")
				err := transferer.Transfer(ctx, bidderAccount, settlementChainID, settlementTopup)
				if err != nil {
					logger.Error("failed to transfer funds to bidder account", "error", err)
				} else {
					logger.Info("successfully transferred funds to bidder account")
					if err := notifier.SendBidderFundedNotification(
						ctx,
						bidderAccount,
						settlementTopup,
					); err != nil {
						logger.Error("failed to send bidder funded notification", "error", err)
					}
				}
				time.Sleep(1 * time.Minute) // Prevent rapid retries
			}
		}
	}()

	return done
}

func setupDeposits(bidderCli bidderapiv1.BidderClient, amount *big.Int) error {
	status, err := bidderCli.DepositManagerStatus(context.Background(), &bidderapiv1.DepositManagerStatusRequest{})
	if err != nil {
		return fmt.Errorf("failed to get deposit manager status: %w", err)
	}
	if !status.Enabled {
		resp, err := bidderCli.EnableDepositManager(context.Background(), &bidderapiv1.EnableDepositManagerRequest{})
		if err != nil {
			return fmt.Errorf("failed to enable deposit manager: %w", err)
		}
		if !resp.Success {
			return fmt.Errorf("failed to enable deposit manager")
		}
	}

	validProviders, err := bidderCli.GetValidProviders(context.Background(), &bidderapiv1.GetValidProvidersRequest{})
	if err != nil {
		return fmt.Errorf("failed to get valid providers: %w", err)
	}
	if len(validProviders.ValidProviders) == 0 {
		return fmt.Errorf("no valid providers found")
	}

	targetDeposits := make([]*bidderapiv1.TargetDeposit, 0, len(validProviders.ValidProviders))
	for _, provider := range validProviders.ValidProviders {
		if status.Enabled && slices.ContainsFunc(status.TargetDeposits, func(td *bidderapiv1.TargetDeposit) bool {
			if td.Provider == provider && td.TargetDeposit == amount.String() {
				return true
			}
			return false
		}) {
			continue
		}
		targetDeposits = append(targetDeposits, &bidderapiv1.TargetDeposit{
			Provider:      provider,
			TargetDeposit: amount.String(),
		})
	}

	if len(targetDeposits) > 0 {
		resp, err := bidderCli.SetTargetDeposits(context.Background(), &bidderapiv1.SetTargetDepositsRequest{
			TargetDeposits: targetDeposits,
		})
		if err != nil {
			return fmt.Errorf("failed to set target deposits: %w", err)
		}
		if len(resp.SuccessfullySetDeposits) != len(targetDeposits) {
			return fmt.Errorf("failed to set target deposits")
		}
	}

	return nil
}

func getBidderEOA(debugClient debugapiv1.DebugServiceClient) (common.Address, error) {
	info, err := debugClient.GetTopology(context.Background(), &debugapiv1.EmptyMessage{})
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get node info: %w", err)
	}
	self := info.Topology.Fields["self"].GetStructValue()
	if self == nil {
		return common.Address{}, fmt.Errorf("self field not found in topology")
	}
	addressHex := self.Fields["Ethereum Address"].GetStringValue()
	if addressHex == "" {
		return common.Address{}, fmt.Errorf("ethereum address not found in topology self field")
	}
	return common.HexToAddress(addressHex), nil
}
