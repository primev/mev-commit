package orchestrator

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Orchestrator interface {
	Providers() []Provider
	Bidders() []Bidder
	Bootnodes() []Bootnode

	ProviderRegistry() *providerregistry.ProviderregistryFilterer
	Logger() *slog.Logger

	io.Closer
}

type BaseNode interface {
	EthAddress() string
	DebugAPI() debugapiv1.DebugServiceClient

	io.Closer
}

type Provider interface {
	BaseNode

	ProviderAPI() providerapiv1.ProviderClient
}

type Bidder interface {
	BaseNode

	BidderAPI() bidderapiv1.BidderClient
}

type Bootnode interface {
	BaseNode
}

type Options struct {
	SettlementRPCEndpoint   string
	ProviderRegistryAddress common.Address
	ProviderRPCAddresses    []string
	BidderRPCAddresses      []string
	BootnodeRPCAddresses    []string
	Logger                  *slog.Logger
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

func newNode(rpcAddr string, logger *slog.Logger) (*node, error) {
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
		// {"TLS system pool certificate", true, credentials.NewClientTLSFromCert(nil, "")},
		{"TLS skip verification", false, credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})},
		{"TLS disabled", false, insecure.NewCredentials()},
	} {
		logger.Info("dialing to grpc server", "strategy", e.strategy)
		conn, err = grpc.DialContext(
			context.Background(),
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

	providerRegistry *providerregistry.ProviderregistryFilterer
	logger           *slog.Logger
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

func (o *orchestrator) ProviderRegistry() *providerregistry.ProviderregistryFilterer {
	return o.providerRegistry
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

	return errs
}

func NewOrchestrator(opts Options) (Orchestrator, error) {
	providers := make([]Provider, 0, len(opts.ProviderRPCAddresses))
	for _, rpcAddr := range opts.ProviderRPCAddresses {
		n, err := newNode(rpcAddr, opts.Logger)
		if err != nil {
			return nil, err
		}
		providers = append(providers, n)
	}

	bidders := make([]Bidder, 0, len(opts.BidderRPCAddresses))
	for _, rpcAddr := range opts.BidderRPCAddresses {
		n, err := newNode(rpcAddr, opts.Logger)
		if err != nil {
			return nil, err
		}
		bidders = append(bidders, n)
	}

	// bootnodes := make([]Bootnode, 0, len(opts.BootnodeRPCAddresses))
	// for _, rpcAddr := range opts.BootnodeRPCAddresses {
	// 	n, err := newNode(rpcAddr, opts.Logger)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	bootnodes = append(bootnodes, n)
	// }

	ethClient, err := ethclient.Dial(opts.SettlementRPCEndpoint)
	if err != nil {
		return nil, err
	}

	providerRegistry, err := providerregistry.NewProviderregistryFilterer(opts.ProviderRegistryAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &orchestrator{
		providers: providers,
		bidders:   bidders,
		// bootnodes:        bootnodes,
		providerRegistry: providerRegistry,
		logger:           opts.Logger,
	}, nil
}
