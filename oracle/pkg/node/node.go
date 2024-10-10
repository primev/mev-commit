package node

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	rollupclient "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/oracle/pkg/apiserver"
	"github.com/primev/mev-commit/oracle/pkg/l1Listener"
	"github.com/primev/mev-commit/oracle/pkg/store"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
)

const defaultMetricsNamespace = "mev_commit_oracle"

func init() {
	setupMetricsNamespace(defaultMetricsNamespace)
}

type Options struct {
	Logger                       *slog.Logger
	KeySigner                    keysigner.KeySigner
	HTTPPort                     int
	SettlementRPCUrl             string
	RelayUrls                    []string
	L1RPCUrls                    []string
	OracleContractAddr           common.Address
	PreconfContractAddr          common.Address
	BlockTrackerContractAddr     common.Address
	ProviderRegistryContractAddr common.Address
	BidderRegistryContractAddr   common.Address
	PgHost                       string
	PgPort                       int
	PgUser                       string
	PgPassword                   string
	PgDbname                     string
	LaggerdMode                  int
	OverrideWinners              []string
	RegistrationAuthToken        string
	DefaultGasLimit              uint64
	DefaultGasTipCap             *big.Int
	DefaultGasFeeCap             *big.Int
}

type Node struct {
	logger    *slog.Logger
	waitClose func()
	dbCloser  io.Closer
}

func NewNode(opts *Options) (*Node, error) {
	nd := &Node{logger: opts.Logger}
	healthChecker := health.New()

	db, err := initDB(opts)
	if err != nil {
		opts.Logger.Error("failed initializing DB", "error", err)
		return nil, err
	}
	nd.dbCloser = db

	st, err := store.NewStore(db)
	if err != nil {
		nd.logger.Error("failed initializing store", "error", err)
		return nil, err
	}

	owner := opts.KeySigner.GetAddress()

	settlementClient, err := ethclient.Dial(opts.SettlementRPCUrl)
	if err != nil {
		nd.logger.Error("failed to connect to the settlement layer", "error", err)
		return nil, err
	}

	chainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		nd.logger.Error("failed getting chain ID", "error", err)
		return nil, err
	}

	contracts, err := getContractABIs(opts)
	if err != nil {
		nd.logger.Error("failed to get contract ABIs", "error", err)
		return nil, err
	}

	monitor := txmonitor.New(
		owner,
		settlementClient,
		txmonitor.NewEVMHelperWithLogger(settlementClient, nd.logger, contracts),
		st,
		nd.logger.With("component", "tx_monitor"),
		1024,
	)

	ctx, cancel := context.WithCancel(context.Background())
	monitorClosed := monitor.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("txmonitor", monitorClosed))

	txnMgr := transactor.NewTransactor(
		settlementClient,
		monitor,
	)
	settlementRPC := transactor.NewMetricsWrapper(txnMgr)

	abis := make([]*abi.ABI, 0, len(contracts))
	contractAddrs := make([]common.Address, 0, len(contracts))

	for addr, abi := range contracts {
		abis = append(abis, abi)
		contractAddrs = append(contractAddrs, addr)
	}

	evtMgr := events.NewListener(
		nd.logger.With("component", "events"),
		abis...,
	)

	var eventsPublisher publisherStartable
	if u, err := url.Parse(opts.SettlementRPCUrl); err == nil && u.Scheme == "ws" {
		eventsPublisher = publisher.NewWSPublisher(
			st,
			nd.logger.With("component", "ws_publisher"),
			settlementClient,
			evtMgr,
		)
	} else {
		eventsPublisher = publisher.NewHTTPPublisher(
			st,
			nd.logger.With("component", "http_publisher"),
			settlementClient,
			evtMgr,
		)
	}

	blockTracker, err := blocktracker.NewBlocktrackerTransactor(
		opts.BlockTrackerContractAddr,
		settlementRPC,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate block tracker contract", "error", err)
		cancel()
		return nil, err
	}

	oracleTransactor, err := rollupclient.NewOracleTransactor(
		opts.OracleContractAddr,
		settlementRPC,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate oracle transactor", "error", err)
		cancel()
		return nil, err
	}

	tOpts, err := opts.KeySigner.GetAuth(chainID)
	if err != nil {
		nd.logger.Error("failed to get auth", "error", err)
		cancel()
		return nil, err
	}

	// Set default gas values
	tOpts.GasLimit = opts.DefaultGasLimit
	tOpts.GasTipCap = opts.DefaultGasTipCap
	tOpts.GasFeeCap = opts.DefaultGasFeeCap

	blockTrackerTransactor := &blocktracker.BlocktrackerTransactorSession{
		Contract:     blockTracker,
		TransactOpts: *tOpts,
	}

	oracleTransactorSession := &rollupclient.OracleTransactorSession{
		Contract:     oracleTransactor,
		TransactOpts: *tOpts,
	}

	l1ClientOpts := []l1ClientOptions{
		l1ClientWithMaxRetries(30),
	}
	if opts.LaggerdMode > 0 {
		l1ClientOpts = append(
			l1ClientOpts,
			l1ClientWithBlockNumberDrift(opts.LaggerdMode),
		)
	}
	if len(opts.OverrideWinners) > 0 {
		l1ClientOpts = append(
			l1ClientOpts,
			l1ClientWithWinnersOverride(opts.OverrideWinners),
		)
		for _, winner := range opts.OverrideWinners {
			nd.logger.Info("setting builder mapping", "builderName", winner, "builderAddress", winner)
			err := setBuilderMapping(
				ctx,
				blockTrackerTransactor,
				settlementClient,
				winner,
				winner,
				nd.logger,
			)
			if err != nil {
				nd.logger.Error("failed to set builder mapping", "error", err)
				cancel()
				return nil, err
			}
		}
	}
	l1Client, err := newL1Client(
		nd.logger,
		opts.L1RPCUrls,
		l1ClientOpts...,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate L1 client", "error", err)
		cancel()
		return nil, err
	}

	l1Lis := l1Listener.NewL1Listener(
		nd.logger.With("component", "l1_listener"),
		l1Client,
		st,
		evtMgr,
		blockTrackerTransactor,
		opts.RelayUrls,
	)
	l1LisClosed := l1Lis.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("l1_listener", l1LisClosed))

	updtr, err := updater.NewUpdater(
		nd.logger.With("component", "updater"),
		l1Client,
		st,
		evtMgr,
		oracleTransactorSession,
		txmonitor.NewEVMHelperWithLogger(l1Client.clients[0].cli.(*ethclient.Client), nd.logger, contracts),
	)
	if err != nil {
		nd.logger.Error("failed to instantiate updater", "error", err)
		cancel()
		return nil, err
	}

	updtrClosed := updtr.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("updater", updtrClosed))

	providerRegistry, err := providerregistry.NewProviderregistryCaller(
		opts.ProviderRegistryContractAddr,
		settlementClient,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate provider registry contract", "error", err)
		cancel()
		return nil, err
	}

	providerRegistryCaller := &providerregistry.ProviderregistryCallerSession{
		Contract: providerRegistry,
		CallOpts: bind.CallOpts{
			From:    opts.KeySigner.GetAddress(),
			Pending: false,
		},
	}

	srv := apiserver.New(
		nd.logger.With("component", "apiserver"),
		evtMgr,
		st,
		opts.RegistrationAuthToken,
		blockTrackerTransactor,
		providerRegistryCaller,
		monitor,
	)

	pubDone := eventsPublisher.Start(ctx, contractAddrs...)
	healthChecker.Register(health.CloseChannelHealthCheck("events_publisher", pubDone))

	srv.RegisterMetricsCollectors(l1Lis.Metrics()...)
	srv.RegisterMetricsCollectors(updtr.Metrics()...)
	srv.RegisterMetricsCollectors(monitor.Metrics()...)
	srv.RegisterMetricsCollectors(evtMgr.Metrics()...)
	srv.RegisterMetricsCollectors(settlementRPC.Metrics()...)
	srv.RegisterHealthCheck(healthChecker)

	srvClosed := srv.Start(fmt.Sprintf(":%d", opts.HTTPPort))

	nd.waitClose = func() {
		cancel()

		_ = srv.Stop()

		closeChan := make(chan struct{})
		go func() {
			defer close(closeChan)

			<-l1LisClosed
			<-updtrClosed
			<-srvClosed
			<-pubDone
			<-monitorClosed
		}()

		<-closeChan
	}

	return nd, nil
}

func (n *Node) Close() (err error) {
	defer func() {
		if n.dbCloser != nil {
			if err2 := n.dbCloser.Close(); err2 != nil {
				err = errors.Join(err, err2)
			}
		}
	}()
	workersClosed := make(chan struct{})
	go func() {
		defer close(workersClosed)

		if n.waitClose != nil {
			n.waitClose()
		}
	}()

	select {
	case <-workersClosed:
		n.logger.Info("all workers closed")
		return nil
	case <-time.After(10 * time.Second):
		n.logger.Error("timeout waiting for workers to close")
		return errors.New("timeout waiting for workers to close")
	}
}

func initDB(opts *Options) (db *sql.DB, err error) {
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

func getContractABIs(opts *Options) (map[common.Address]*abi.ABI, error) {
	abis := make(map[common.Address]*abi.ABI)

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		return nil, err
	}
	abis[opts.BlockTrackerContractAddr] = &btABI

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		return nil, err
	}
	abis[opts.PreconfContractAddr] = &pcABI

	bidderRegistry, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[opts.BidderRegistryContractAddr] = &bidderRegistry

	providerRegistry, err := abi.JSON(strings.NewReader(providerregistry.ProviderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[opts.ProviderRegistryContractAddr] = &providerRegistry

	orABI, err := abi.JSON(strings.NewReader(rollupclient.OracleABI))
	if err != nil {
		return nil, err
	}
	abis[opts.OracleContractAddr] = &orABI

	return abis, nil
}

func setBuilderMapping(
	ctx context.Context,
	bt *blocktracker.BlocktrackerTransactorSession,
	client *ethclient.Client,
	builderName string,
	builderAddress string,
	logger *slog.Logger,
) error {
	logger.Info("setting builder mapping", "builderName", builderName, "builderAddress", builderAddress)

	txn, err := bt.AddBuilderAddress(builderName, common.HexToAddress(builderAddress))
	if err != nil {
		return fmt.Errorf("unable to add builder address: %w", err)
	}

	logger.Info("waiting for tx to be mined", "txHash", txn.Hash().Hex(), "nonce", txn.Nonce())
	_, err = bind.WaitMined(ctx, client, txn)
	if err != nil {
		return fmt.Errorf("unable to wait for tx to be minted: %w", err)
	}

	return nil
}

func setupMetricsNamespace(ns string) {
	transactor.Namespace = ns
	txmonitor.Namespace = ns
	events.Namespace = ns
}

type publisherStartable interface {
	Start(context.Context, ...common.Address) <-chan struct{}
}
