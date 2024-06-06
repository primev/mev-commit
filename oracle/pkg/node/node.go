package node

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	rollupclient "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/oracle/pkg/apiserver"
	"github.com/primev/mev-commit/oracle/pkg/l1Listener"
	"github.com/primev/mev-commit/oracle/pkg/store"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
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
	L1RPCUrl                     string
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
}

type Node struct {
	logger    *slog.Logger
	waitClose func()
	dbCloser  io.Closer
}

func NewNode(opts *Options) (*Node, error) {
	nd := &Node{logger: opts.Logger}

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

	l1Client, err := ethclient.Dial(opts.L1RPCUrl)
	if err != nil {
		nd.logger.Error("Failed to connect to the L1 Ethereum client", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := txmonitor.New(
		owner,
		settlementClient,
		txmonitor.NewEVMHelper(settlementClient.Client()),
		st,
		nd.logger.With("component", "tx_monitor"),
		1024,
	)

	monitorClosed := monitor.Start(ctx)

	txnMgr := transactor.NewTransactor(
		settlementClient,
		monitor,
	)
	settlementRPC := transactor.NewMetricsWrapper(txnMgr)

	contracts, err := getContractABIs(opts)
	if err != nil {
		nd.logger.Error("failed to get contract ABIs", "error", err)
		cancel()
		return nil, err
	}

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

	httpPub := publisher.NewHTTPPublisher(
		st,
		nd.logger.With("component", "http_publisher"),
		settlementClient,
		evtMgr,
	)

	var listenerL1Client l1Listener.EthClient

	listenerL1Client = l1Client
	if opts.LaggerdMode > 0 {
		listenerL1Client = &laggerdL1Client{EthClient: listenerL1Client, amount: opts.LaggerdMode}
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

	blockTrackerTransactor := &blocktracker.BlocktrackerTransactorSession{
		Contract:     blockTracker,
		TransactOpts: *tOpts,
	}

	oracleTransactorSession := &rollupclient.OracleTransactorSession{
		Contract:     oracleTransactor,
		TransactOpts: *tOpts,
	}

	if opts.OverrideWinners != nil && len(opts.OverrideWinners) > 0 {
		listenerL1Client = &winnerOverrideL1Client{EthClient: listenerL1Client, winners: opts.OverrideWinners}
		for _, winner := range opts.OverrideWinners {
			err := setBuilderMapping(
				ctx,
				blockTrackerTransactor,
				settlementClient,
				winner,
				winner,
			)
			if err != nil {
				nd.logger.Error("failed to set builder mapping", "error", err)
				cancel()
				return nil, err
			}
		}
	}

	l1Lis := l1Listener.NewL1Listener(
		nd.logger.With("component", "l1_listener"),
		listenerL1Client,
		st,
		evtMgr,
		blockTrackerTransactor,
	)
	l1LisClosed := l1Lis.Start(ctx)

	updtr, err := updater.NewUpdater(
		nd.logger.With("component", "updater"),
		l1Client,
		st,
		evtMgr,
		oracleTransactorSession,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate updater", "error", err)
		cancel()
		return nil, err
	}

	updtrClosed := updtr.Start(ctx)

	srv := apiserver.New(
		nd.logger.With("component", "apiserver"),
		evtMgr,
		st,
	)

	httpPubDone := httpPub.Start(ctx, contractAddrs...)

	srv.RegisterMetricsCollectors(l1Lis.Metrics()...)
	srv.RegisterMetricsCollectors(updtr.Metrics()...)
	srv.RegisterMetricsCollectors(monitor.Metrics()...)
	srv.RegisterMetricsCollectors(evtMgr.Metrics()...)
	srv.RegisterMetricsCollectors(settlementRPC.Metrics()...)

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
			<-httpPubDone
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

	return db, err
}

func getContractABIs(opts *Options) (map[common.Address]*abi.ABI, error) {
	abis := make(map[common.Address]*abi.ABI)

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		return nil, err
	}
	abis[opts.BlockTrackerContractAddr] = &btABI

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfcommitmentstoreABI))
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

type laggerdL1Client struct {
	l1Listener.EthClient
	amount int
}

func (l *laggerdL1Client) BlockNumber(ctx context.Context) (uint64, error) {
	blkNum, err := l.EthClient.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blkNum - uint64(l.amount), nil
}

type winnerOverrideL1Client struct {
	l1Listener.EthClient
	winners []string
}

func (w *winnerOverrideL1Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	hdr, err := w.EthClient.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	idx := number.Int64() % int64(len(w.winners))
	hdr.Extra = []byte(w.winners[idx])

	return hdr, nil
}

func setBuilderMapping(
	ctx context.Context,
	bt *blocktracker.BlocktrackerTransactorSession,
	client *ethclient.Client,
	builderName string,
	builderAddress string,
) error {
	txn, err := bt.AddBuilderAddress(builderName, common.HexToAddress(builderAddress))
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, client, txn)
	if err != nil {
		return err
	}

	return nil
}

func setupMetricsNamespace(ns string) {
	transactor.Namespace = ns
	txmonitor.Namespace = ns
	events.Namespace = ns
}
