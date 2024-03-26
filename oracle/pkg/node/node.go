package node

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	rollupclient "github.com/primevprotocol/contracts-abi/clients/Oracle"
	preconf "github.com/primevprotocol/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-commit/oracle/pkg/apiserver"
	"github.com/primevprotocol/mev-commit/oracle/pkg/keysigner"
	"github.com/primevprotocol/mev-commit/oracle/pkg/l1Listener"
	"github.com/primevprotocol/mev-commit/oracle/pkg/settler"
	"github.com/primevprotocol/mev-commit/oracle/pkg/store"
	"github.com/primevprotocol/mev-commit/oracle/pkg/updater"
)

type Options struct {
	Logger              *slog.Logger
	KeySigner           keysigner.KeySigner
	HTTPPort            int
	SettlementRPCUrl    string
	L1RPCUrl            string
	OracleContractAddr  common.Address
	PreconfContractAddr common.Address
	PgHost              string
	PgPort              int
	PgUser              string
	PgPassword          string
	PgDbname            string
	LaggerdMode         int
	OverrideWinners     []string
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

	l2Client, err := ethclient.Dial(opts.SettlementRPCUrl)
	if err != nil {
		nd.logger.Error("Failed to connect to the L2 Ethereum client", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	var listenerL1Client l1Listener.EthClient

	listenerL1Client = l1Client
	if opts.LaggerdMode > 0 {
		listenerL1Client = &laggerdL1Client{EthClient: listenerL1Client, amount: opts.LaggerdMode}
	}

	preconfContract, err := preconf.NewPreconfcommitmentstoreCaller(
		opts.PreconfContractAddr,
		settlementClient,
	)
	if err != nil {
		nd.logger.Error("failed to instantiate preconf contract", "error", err)
		cancel()
		return nil, err
	}

	oracleContract, err := rollupclient.NewOracle(opts.OracleContractAddr, settlementClient)
	if err != nil {
		nd.logger.Error("failed to instantiate oracle contract", "error", err)
		cancel()
		return nil, err
	}

	if opts.OverrideWinners != nil && len(opts.OverrideWinners) > 0 {
		listenerL1Client = &winnerOverrideL1Client{EthClient: listenerL1Client, winners: opts.OverrideWinners}
		for _, winner := range opts.OverrideWinners {
			err := setBuilderMapping(
				ctx,
				opts.KeySigner,
				chainID,
				settlementClient,
				oracleContract,
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

	l1Lis := l1Listener.NewL1Listener(nd.logger.With("component", "l1_listener"), listenerL1Client, st)
	l1LisClosed := l1Lis.Start(ctx)

	callOpts := bind.CallOpts{
		Pending: false,
		From:    owner,
		Context: ctx,
	}

	pc := &preconf.PreconfcommitmentstoreCallerSession{
		Contract: preconfContract,
		CallOpts: callOpts,
	}
	oc := &rollupclient.OracleSession{Contract: oracleContract, CallOpts: callOpts}

	updtr := updater.NewUpdater(nd.logger.With("component", "updater"), l1Client, l2Client, st, oc, pc)
	updtrClosed := updtr.Start(ctx)

	settlr := settler.NewSettler(
		nd.logger.With("component", "settler"),
		opts.KeySigner,
		chainID,
		owner,
		oracleContract,
		st,
		settlementClient,
	)
	settlrClosed := settlr.Start(ctx)

	srv := apiserver.New(nd.logger.With("component", "apiserver"), st)
	srv.RegisterMetricsCollectors(l1Lis.Metrics()...)
	srv.RegisterMetricsCollectors(updtr.Metrics()...)
	srv.RegisterMetricsCollectors(settlr.Metrics()...)

	srvClosed := srv.Start(fmt.Sprintf(":%d", opts.HTTPPort))

	nd.waitClose = func() {
		cancel()

		_ = srv.Stop()

		closeChan := make(chan struct{})
		go func() {
			defer close(closeChan)

			<-l1LisClosed
			<-updtrClosed
			<-settlrClosed
			<-srvClosed
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
	keySigner keysigner.KeySigner,
	chainID *big.Int,
	client *ethclient.Client,
	rc *rollupclient.Oracle,
	builderName string,
	builderAddress string,
) error {
	auth, err := keySigner.GetAuth(chainID)
	if err != nil {
		return err
	}
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Returns priority fee per gas
	gasTip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	// Returns priority fee per gas + base fee per gas
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTip

	txn, err := rc.AddBuilderAddress(auth, builderName, common.HexToAddress(builderAddress))
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, client, txn)
	if err != nil {
		return err
	}

	return nil
}
