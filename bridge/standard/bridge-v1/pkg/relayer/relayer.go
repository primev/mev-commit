package relayer

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	l1g "github.com/primevprotocol/contracts-abi/clients/L1Gateway"
	sg "github.com/primevprotocol/contracts-abi/clients/SettlementGateway"
	"github.com/primevprotocol/mev-commit/bridge/standard/bridge-v1/pkg/shared"
	"golang.org/x/crypto/sha3"
)

type Options struct {
	Ctx                    context.Context
	Logger                 *slog.Logger
	PrivateKey             *ecdsa.PrivateKey
	SettlementRPCUrl       string
	L1RPCUrl               string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
}

type Relayer struct {
	logger *slog.Logger
	// Closes ctx's Done channel and waits for all goroutines to close.
	waitOnCloseRoutines func()
	db                  *sql.DB
}

func NewRelayer(opts *Options) (r *Relayer, err error) {
	r = &Relayer{logger: opts.Logger}

	pubKey := &opts.PrivateKey.PublicKey
	pubKeyBytes := crypto.FromECDSAPub(pubKey)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKeyBytes[1:])
	address := hash.Sum(nil)[12:]

	r.logger.Info("relayer signing address", "address", common.BytesToAddress(address).Hex())

	l1Client, err := ethclient.DialContext(opts.Ctx, opts.L1RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial l1 rpc: %w", err)
	}

	l1ChainID, err := l1Client.ChainID(opts.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get l1 chain id: %w", err)
	}
	r.logger.Info("L1 chain id", "chain_id", l1ChainID)

	settlementClient, err := ethclient.Dial(opts.SettlementRPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial settlement rpc: %w", err)
	}

	settlementChainID, err := settlementClient.ChainID(opts.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to dial settlement rpc: %w", err)
	}
	r.logger.Info("settlement chain id", "chain_id", settlementChainID)

	sFilterer, err := shared.NewSettlementFilterer(opts.SettlementContractAddr, settlementClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement filterer: %w", err)
	}

	ctx, cancel := context.WithCancel(opts.Ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	sListener := NewListener(r.logger.With("component", "settlement_listener"), settlementClient, sFilterer, false)
	sListenerClosed, settlementEventChan, err := sListener.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start settlement listener: %w", err)
	}

	l1Filterer, err := shared.NewL1Filterer(opts.L1ContractAddr, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 filterer: %w", err)
	}

	l1Listener := NewListener(r.logger.With("component", "l1_listener"), l1Client, l1Filterer, true)
	l1ListenerClosed, l1EventChan, err := l1Listener.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start l1 listener: %w", err)
	}

	st, err := sg.NewSettlementgatewayTransactor(opts.SettlementContractAddr, settlementClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement gateway transactor: %w", err)
	}
	settlementTransactor := NewTransactor(
		r.logger.With("component", "settlement_transactor"),
		opts.PrivateKey,
		opts.SettlementContractAddr,
		settlementClient,
		st,
		sFilterer,
		l1EventChan, // L1 transfer initiations result in settlement finalizations
	)
	stClosed, err := settlementTransactor.Start(ctx)
	if err != nil {
		return nil, err
	}

	l1t, err := l1g.NewL1gatewayTransactor(opts.L1ContractAddr, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 gateway transactor: %w", err)
	}
	l1Transactor := NewTransactor(
		r.logger.With("component", "l1_transactor"),
		opts.PrivateKey,
		opts.L1ContractAddr,
		l1Client,
		l1t,
		l1Filterer,
		settlementEventChan, // Settlement transfer initiations result in L1 finalizations
	)
	l1tClosed, err := l1Transactor.Start(ctx)
	if err != nil {
		return nil, err
	}

	r.waitOnCloseRoutines = func() {
		// Close ctx's Done channel
		cancel()

		allClosed := make(chan struct{})
		go func() {
			defer close(allClosed)
			<-sListenerClosed
			<-l1ListenerClosed
			<-stClosed
			<-l1tClosed
		}()
		<-allClosed
	}
	return r, nil
}

// TryCloseAll attempts to close all workers and the database connection.
func (r *Relayer) TryCloseAll() (err error) {
	r.logger.Debug("closing all workers and db connection")
	defer func() {
		if r.db == nil {
			return
		}
		if err2 := r.db.Close(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	workersClosed := make(chan struct{})
	go func() {
		defer close(workersClosed)
		r.waitOnCloseRoutines()
	}()

	select {
	case <-workersClosed:
		r.logger.Info("all workers closed")
		return nil
	case <-time.After(10 * time.Second):
		msg := "failed to close all workers in 10 sec"
		r.logger.Error(msg)
		return errors.New(msg)
	}
}
