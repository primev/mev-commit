package orchestrator

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Orchestrator interface {
	io.Closer

	// Providers returns the list of providers
	Providers() []Provider
	// Bidders returns the list of bidders
	Bidders() []Bidder
	// Bootnodes returns the list of bootnodes
	Bootnodes() []Bootnode

	// Events returns the event manager used to listen to contract events
	Events() events.EventManager
	// Logger returns the logger used by the orchestrator
	Logger() *slog.Logger

	// L1Client returns the L1 RPC client
	L1Client() *ethclient.Client
}

type Node interface {
	io.Closer

	EthAddress() string
	DebugAPI() debugapiv1.DebugServiceClient
}

type Provider interface {
	Node

	ProviderAPI() providerapiv1.ProviderClient
}

type Bidder interface {
	Node

	BidderAPI() bidderapiv1.BidderClient
}

type Bootnode interface {
	Node
}

type Options struct {
	L1RPCEndpoint               string
	SettlementRPCEndpoint       string
	ProviderRegistryAddress     common.Address
	BlockTrackerContractAddress common.Address
	PreconfContractAddress      common.Address
	BidderRegistryAddress       common.Address
	OracleContractAddress       common.Address
	ProviderRPCAddresses        []string
	BidderRPCAddresses          []string
	BootnodeRPCAddresses        []string
	Logger                      *slog.Logger
}

type node struct {
	ethAddr string
	conn    *grpc.ClientConn
}

func (n *node) EthAddress() string {
	return n.ethAddr
}

func (n *node) DebugAPI() debugapiv1.DebugServiceClient {
	return debugapiv1.NewDebugServiceClient(n.conn)
}

func (n *node) ProviderAPI() providerapiv1.ProviderClient {
	return providerapiv1.NewProviderClient(n.conn)
}

func (n *node) BidderAPI() bidderapiv1.BidderClient {
	return bidderapiv1.NewBidderClient(n.conn)
}

func (n *node) Close() error {
	return n.conn.Close()
}

func newNode(rpcAddr string, logger *slog.Logger) (any, error) {
	// Since we don't know if the server has TLS enabled on its rpc
	// endpoint, we try different strategies from most secure to
	// least secure. In the future, when only TLS-enabled servers
	// are allowed, only the TLS system pool certificate strategy
	// should be used.
	var conn *grpc.ClientConn
	var err error

	for _, e := range []struct {
		strategy   string
		isSecure   bool
		credential credentials.TransportCredentials
	}{
		{"TLS skip verification", false, credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})},
		{"TLS disabled", false, insecure.NewCredentials()},
	} {
		logger.Info("dialing to grpc server", "strategy", e.strategy)
		conn, err = grpc.NewClient(
			rpcAddr,
			grpc.WithTransportCredentials(e.credential),
		)
		if err != nil {
			logger.Error("failed to dial grpc server", "error", err)
			continue
		}

		if !e.isSecure {
			logger.Warn("established connection with the grpc server has potential security risk")
		}
		break
	}
	if conn == nil {
		logger.Error("dialing of grpc server failed")
		return nil, fmt.Errorf("dialing of grpc server failed")
	}

	topo, err := debugapiv1.NewDebugServiceClient(conn).GetTopology(context.Background(), &debugapiv1.EmptyMessage{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s topology: %w", rpcAddr, err)
	}

	ethAddr := topo.Topology.Fields["self"].GetStructValue().Fields["Ethereum Address"].GetStringValue()
	if ethAddr == "" {
		return nil, fmt.Errorf("ethereum address not found in topology")
	}

	return &node{
		ethAddr: ethAddr,
		conn:    conn,
	}, nil
}

type orchestrator struct {
	providers []Provider
	bidders   []Bidder
	bootnodes []Bootnode

	l1RPC      *ethclient.Client
	evtMgr     events.EventManager
	logger     *slog.Logger
	pubCancel  context.CancelFunc
	pubStopped <-chan struct{}
}

func (o *orchestrator) Providers() []Provider {
	return o.providers
}

func (o *orchestrator) Bidders() []Bidder {
	return o.bidders
}

func (o *orchestrator) Bootnodes() []Bootnode {
	return o.bootnodes
}

func (o *orchestrator) L1Client() *ethclient.Client {
	return o.l1RPC
}

func (o *orchestrator) Events() events.EventManager {
	return o.evtMgr
}

func (o *orchestrator) Logger() *slog.Logger {
	return o.logger
}

func (o *orchestrator) Close() error {
	var errs error
	for _, p := range o.providers {
		if err := p.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	for _, b := range o.bidders {
		if err := b.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	for _, b := range o.bootnodes {
		if err := b.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	o.pubCancel()
	<-o.pubStopped

	return errs
}

func createNodes[T any](logger *slog.Logger, rpcAddrs []string) ([]T, error) {
	nodes := make([]T, 0, len(rpcAddrs))
	for _, rpcAddr := range rpcAddrs {
		n, err := newNode(rpcAddr, logger)
		if err != nil {
			return nil, err
		}
		tn, ok := n.(T)
		if !ok {
			return nil, fmt.Errorf("unexpected node type")
		}
		nodes = append(nodes, tn)
	}
	return nodes, nil
}

func NewOrchestrator(opts Options) (Orchestrator, error) {
	providers, err := createNodes[Provider](opts.Logger, opts.ProviderRPCAddresses)
	if err != nil {
		return nil, err
	}

	bidders, err := createNodes[Bidder](opts.Logger, opts.BidderRPCAddresses)
	if err != nil {
		return nil, err
	}

	bootnodes, err := createNodes[Bootnode](opts.Logger, opts.BootnodeRPCAddresses)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethclient.Dial(opts.SettlementRPCEndpoint)
	if err != nil {
		return nil, err
	}

	contracts, err := getContractABIs(opts)
	if err != nil {
		return nil, err
	}

	abis := make([]*abi.ABI, 0, len(contracts))
	contractAddrs := make([]common.Address, 0, len(contracts))

	for addr, abi := range contracts {
		abis = append(abis, abi)
		contractAddrs = append(contractAddrs, addr)
	}

	evtMgr := events.NewListener(
		opts.Logger.With("component", "events"),
		abis...,
	)

	evtPublisher := publisher.NewWSPublisher(
		dummyStore{},
		opts.Logger.With("component", "ws_publisher"),
		ethClient,
		evtMgr,
	)

	l1RPC, err := ethclient.Dial(opts.L1RPCEndpoint)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	stopped := evtPublisher.Start(ctx, contractAddrs...)

	return &orchestrator{
		providers:  providers,
		bidders:    bidders,
		bootnodes:  bootnodes,
		l1RPC:      l1RPC,
		evtMgr:     evtMgr,
		logger:     opts.Logger,
		pubCancel:  cancel,
		pubStopped: stopped,
	}, nil
}

func getContractABIs(opts Options) (map[common.Address]*abi.ABI, error) {
	abis := make(map[common.Address]*abi.ABI)

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		return nil, err
	}
	abis[opts.BlockTrackerContractAddress] = &btABI

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		return nil, err
	}
	abis[opts.PreconfContractAddress] = &pcABI

	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[opts.BidderRegistryAddress] = &brABI

	prABI, err := abi.JSON(strings.NewReader(providerregistry.ProviderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[opts.ProviderRegistryAddress] = &prABI

	orABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
	if err != nil {
		return nil, err
	}
	abis[opts.OracleContractAddress] = &orABI

	return abis, nil
}

type dummyStore struct{}

func (dummyStore) SetLastBlock(block uint64) error {
	return nil
}

func (dummyStore) LastBlock() (uint64, error) {
	return 0, nil
}
