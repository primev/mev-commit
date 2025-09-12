package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/blocktracker"
	"github.com/primev/mev-commit/tools/preconf-rpc/handlers"
	"github.com/primev/mev-commit/tools/preconf-rpc/notifier"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/tools/preconf-rpc/rpcserver"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"github.com/primev/mev-commit/tools/preconf-rpc/store"
	"github.com/primev/mev-commit/x/accountsync"
	"github.com/primev/mev-commit/x/contracts/ethwrapper"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	bidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"github.com/primev/mev-commit/x/transfer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
	AutoDepositAmount      *big.Int
	L1RPCUrls              []string
	SettlementRPCUrl       string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
	DepositAddress         common.Address
	BridgeAddress          common.Address
	SettlementThreshold    *big.Int
	SettlementTopup        *big.Int
	BidderTopup            *big.Int
	HTTPPort               int
	GasTipCap              *big.Int
	GasFeeCap              *big.Int
	PricerAPIKey           string
	Webhooks               []string
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

	l1RPCClient, err := ethwrapper.NewClient(
		config.Logger.With("module", "ethwrapper"),
		config.L1RPCUrls,
		ethwrapper.EthClientWithMaxRetries(5),
	)
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

	if err := setupDeposits(bidderCli, config.AutoDepositAmount); err != nil {
		return nil, fmt.Errorf("failed to setup deposits: %w", err)
	}

	notifier := notifier.NewNotifier(config.Webhooks, config.Logger.With("module", "notifier"))

	bridgeConfig := transfer.BridgeConfig{
		Signer:                 config.Signer,
		L1ContractAddr:         config.L1ContractAddr,
		SettlementContractAddr: config.SettlementContractAddr,
		L1RPCUrl:               config.L1RPCUrls[0],
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
		l1RPCClient.RawClient(),
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

	bidderFunderDone := startBidderFunder(
		ctx,
		config.Logger.With("module", "bidderfunder"),
		bidderEOA,
		accountsync.NewAccountSync(bidderEOA, settlementClient),
		transferer,
		config.SettlementThreshold,
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

	rpcServer := rpcserver.NewJSONRPCServer(
		config.L1RPCUrls[0],
		config.Logger.With("module", "rpcserver"),
	)

	bidpricer, err := pricer.NewPricer(config.PricerAPIKey, config.Logger.With("module", "bidpricer"))
	if err != nil {
		return nil, fmt.Errorf("failed to create bid pricer: %w", err)
	}

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

	sndr, err := sender.NewTxSender(
		rpcstore,
		bidderClient,
		bidpricer,
		blockTracker,
		transferer,
		notifier,
		settlementChainID,
		config.Logger.With("module", "txsender"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction sender: %w", err)
	}

	senderDone := sndr.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("TxSender", senderDone))

	handlers := handlers.NewRPCMethodHandler(
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

	handlers.RegisterMethods(rpcServer)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := healthChecker.Health(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.Handle("/", rpcServer)

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
