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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/blocktracker"
	"github.com/primev/mev-commit/tools/preconf-rpc/handlers"
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
	HTTPPort               int
	GasTipCap              *big.Int
	GasFeeCap              *big.Int
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

	status, err := bidderCli.AutoDepositStatus(context.Background(), &bidderapiv1.EmptyMessage{})
	if err != nil {
		return nil, err
	}

	if !status.IsAutodepositEnabled {
		_, err := bidderCli.AutoDeposit(
			context.Background(),
			&bidderapiv1.DepositRequest{
				Amount: config.AutoDepositAmount.String(),
			},
		)
		if err != nil {
			return nil, err
		}
	}

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

	bidpricer := &pricer.BidPricer{}

	db, err := initDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	rpcstore, err := store.New(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	blockTracker := blocktracker.NewBlockTracker(
		l1RPCClient,
		config.Logger.With("module", "blocktracker"),
	)
	blockTrackerDone := blockTracker.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("BlockTracker", blockTrackerDone))

	sndr := sender.NewTxSender(
		rpcstore,
		bidderClient,
		bidpricer,
		blockTracker,
		transferer,
		settlementChainID,
		config.Logger.With("module", "txsender"),
	)

	senderDone := sndr.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("TxSender", senderDone))

	handlers := handlers.NewRPCMethodHandler(
		config.Logger.With("module", "handlers"),
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
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := healthChecker.Health(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
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
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.PgHost, opts.PgPort, opts.PgUser, opts.PgPassword, opts.PgDbname,
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
