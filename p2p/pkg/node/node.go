package node

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	bidderapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/apiserver"
	bidder_registrycontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/bidder_registry"
	preconfcontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/preconf"
	provider_registrycontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/provider_registry"
	"github.com/primevprotocol/mev-commit/p2p/pkg/debugapi"
	"github.com/primevprotocol/mev-commit/p2p/pkg/discovery"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	"github.com/primevprotocol/mev-commit/p2p/pkg/keysigner"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p/libp2p"
	"github.com/primevprotocol/mev-commit/p2p/pkg/preconfirmation"
	bidderapi "github.com/primevprotocol/mev-commit/p2p/pkg/rpc/bidder"
	providerapi "github.com/primevprotocol/mev-commit/p2p/pkg/rpc/provider"
	"github.com/primevprotocol/mev-commit/p2p/pkg/signer/preconfsigner"
	"github.com/primevprotocol/mev-commit/p2p/pkg/topology"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerDialTimeout = 5 * time.Second
)

type Options struct {
	Version                  string
	KeySigner                keysigner.KeySigner
	Secret                   string
	PeerType                 string
	Logger                   *slog.Logger
	P2PPort                  int
	P2PAddr                  string
	HTTPAddr                 string
	RPCAddr                  string
	Bootnodes                []string
	PreconfContract          string
	ProviderRegistryContract string
	BidderRegistryContract   string
	RPCEndpoint              string
	NatAddr                  string
	TLSCertificateFile       string
	TLSPrivateKeyFile        string
}

type Node struct {
	closers []io.Closer
}

func NewNode(opts *Options) (*Node, error) {
	nd := &Node{
		closers: make([]io.Closer, 0),
	}

	srv := apiserver.New(opts.Version, opts.Logger.With("component", "apiserver"))
	peerType := p2p.FromString(opts.PeerType)

	contractRPC, err := ethclient.Dial(opts.RPCEndpoint)
	if err != nil {
		return nil, err
	}
	evmClient, err := evmclient.New(
		opts.KeySigner,
		evmclient.WrapEthClient(contractRPC),
		opts.Logger.With("component", "evmclient"),
	)
	if err != nil {
		return nil, err
	}
	nd.closers = append(nd.closers, evmClient)

	srv.MetricsRegistry().MustRegister(evmClient.Metrics()...)

	bidderRegistryContractAddr := common.HexToAddress(opts.BidderRegistryContract)

	bidderRegistry := bidder_registrycontract.New(
		bidderRegistryContractAddr,
		evmClient,
		opts.Logger.With("component", "bidderregistry"),
	)

	providerRegistryContractAddr := common.HexToAddress(opts.ProviderRegistryContract)

	providerRegistry := provider_registrycontract.New(
		providerRegistryContractAddr,
		evmClient,
		opts.Logger.With("component", "providerregistry"),
	)

	p2pSvc, err := libp2p.New(&libp2p.Options{
		KeySigner:      opts.KeySigner,
		Secret:         opts.Secret,
		PeerType:       peerType,
		Register:       providerRegistry,
		Logger:         opts.Logger.With("component", "p2p"),
		ListenPort:     opts.P2PPort,
		ListenAddr:     opts.P2PAddr,
		MetricsReg:     srv.MetricsRegistry(),
		BootstrapAddrs: opts.Bootnodes,
		NatAddr:        opts.NatAddr,
	})
	if err != nil {
		return nil, err
	}
	nd.closers = append(nd.closers, p2pSvc)

	topo := topology.New(p2pSvc, opts.Logger.With("component", "topology"))
	disc := discovery.New(topo, p2pSvc, opts.Logger.With("component", "discovery_protocol"))
	nd.closers = append(nd.closers, disc)

	srv.RegisterMetricsCollectors(topo.Metrics()...)

	// Set the announcer for the topology service
	topo.SetAnnouncer(disc)
	// Set the notifier for the p2p service
	p2pSvc.SetNotifier(topo)

	// Register the discovery protocol with the p2p service
	p2pSvc.AddStreamHandlers(disc.Streams()...)

	debugapi.RegisterAPI(srv, topo, p2pSvc, opts.Logger.With("component", "debugapi"))

	if opts.PeerType != p2p.PeerTypeBootnode.String() {
		lis, err := net.Listen("tcp", opts.RPCAddr)
		if err != nil {
			return nil, errors.Join(err, nd.Close())
		}

		var tlsCredentials credentials.TransportCredentials
		if opts.TLSCertificateFile != "" && opts.TLSPrivateKeyFile != "" {
			tlsCredentials, err = credentials.NewServerTLSFromFile(
				opts.TLSCertificateFile,
				opts.TLSPrivateKeyFile,
			)
			if err != nil {
				return nil, fmt.Errorf("unable to load TLS credentials: %w", err)
			}
		}

		grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
		preconfSigner := preconfsigner.NewSigner(opts.KeySigner)
		validator, err := protovalidate.New()
		if err != nil {
			return nil, errors.Join(err, nd.Close())
		}

		var (
			bidProcessor preconfirmation.BidProcessor = noOpBidProcessor{}
			commitmentDA preconfcontract.Interface    = noOpCommitmentDA{}
		)

		switch opts.PeerType {
		case p2p.PeerTypeProvider.String():
			providerAPI := providerapi.NewService(
				opts.Logger.With("component", "providerapi"),
				providerRegistry,
				opts.KeySigner.GetAddress(),
				evmClient,
				validator,
			)
			providerapiv1.RegisterProviderServer(grpcServer, providerAPI)
			bidProcessor = providerAPI
			srv.RegisterMetricsCollectors(providerAPI.Metrics()...)

			preconfContractAddr := common.HexToAddress(opts.PreconfContract)

			commitmentDA = preconfcontract.New(
				preconfContractAddr,
				evmClient,
				opts.Logger.With("component", "preconfcontract"),
			)

			preconfProto := preconfirmation.New(
				topo,
				p2pSvc,
				preconfSigner,
				bidderRegistry,
				bidProcessor,
				commitmentDA,
				opts.Logger.With("component", "preconfirmation_protocol"),
			)
			// Only register handler for provider
			p2pSvc.AddStreamHandlers(preconfProto.Streams()...)
			srv.RegisterMetricsCollectors(preconfProto.Metrics()...)

		case p2p.PeerTypeBidder.String():
			preconfProto := preconfirmation.New(
				topo,
				p2pSvc,
				preconfSigner,
				bidderRegistry,
				bidProcessor,
				commitmentDA,
				opts.Logger.With("component", "preconfirmation_protocol"),
			)
			srv.RegisterMetricsCollectors(preconfProto.Metrics()...)

			bidderAPI := bidderapi.NewService(
				preconfProto,
				opts.KeySigner.GetAddress(),
				bidderRegistry,
				validator,
				opts.Logger.With("component", "bidderapi"),
			)
			bidderapiv1.RegisterBidderServer(grpcServer, bidderAPI)
			srv.RegisterMetricsCollectors(bidderAPI.Metrics()...)
		}

		started := make(chan struct{})
		go func() {
			// signal that the server has started
			close(started)

			err := grpcServer.Serve(lis)
			if err != nil {
				opts.Logger.Error("failed to start grpc server", "err", err)
			}
		}()
		nd.closers = append(nd.closers, lis)

		// Wait for the server to start
		<-started

		// Since we don't know if the server has TLS enabled on its rpc
		// endpoint, we try different strategies from most secure to
		// least secure. In the future, when only TLS-enabled servers
		// are allowed, only the TLS system pool certificate strategy
		// should be used.
		var grpcConn *grpc.ClientConn
		for _, e := range []struct {
			strategy   string
			isSecure   bool
			credential credentials.TransportCredentials
		}{
			{"TLS system pool certificate", true, credentials.NewClientTLSFromCert(nil, "")},
			{"TLS skip verification", false, credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})},
			{"TLS disabled", false, insecure.NewCredentials()},
		} {
			ctx, cancel := context.WithTimeout(context.Background(), grpcServerDialTimeout)
			opts.Logger.Info("dialing to grpc server", "strategy", e.strategy)
			grpcConn, err = grpc.DialContext(
				ctx,
				opts.RPCAddr,
				grpc.WithBlock(),
				grpc.WithTransportCredentials(e.credential),
			)
			if err != nil {
				opts.Logger.Error("failed to dial grpc server", "error", err)
				cancel()
				continue
			}

			cancel()
			if !e.isSecure {
				opts.Logger.Warn("established connection with the grpc server has potential security risk")
			}
			break
		}
		if grpcConn == nil {
			return nil, errors.New("dialing of grpc server failed")
		}

		gatewayMux := runtime.NewServeMux()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		switch opts.PeerType {
		case p2p.PeerTypeProvider.String():
			err := providerapiv1.RegisterProviderHandler(ctx, gatewayMux, grpcConn)
			if err != nil {
				opts.Logger.Error("failed to register provider handler", "err", err)
				return nil, errors.Join(err, nd.Close())
			}
		case p2p.PeerTypeBidder.String():
			err := bidderapiv1.RegisterBidderHandler(ctx, gatewayMux, grpcConn)
			if err != nil {
				opts.Logger.Error("failed to register bidder handler", "err", err)
				return nil, errors.Join(err, nd.Close())
			}
		}

		srv.ChainHandlers("/", gatewayMux)
		srv.ChainHandlers(
			"/health",
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					if s := grpcConn.GetState(); s != connectivity.Ready {
						http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
						return
					}
					fmt.Fprintln(w, "ok")
				},
			),
		)
	}

	server := &http.Server{
		Addr:    opts.HTTPAddr,
		Handler: srv.Router(),
	}

	go func() {
		var (
			err        error
			tlsEnabled = opts.TLSCertificateFile != "" && opts.TLSPrivateKeyFile != ""
		)
		opts.Logger.Info("starting to listen", "tls", tlsEnabled)
		if tlsEnabled {
			err = server.ListenAndServeTLS(
				opts.TLSCertificateFile,
				opts.TLSPrivateKeyFile,
			)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			opts.Logger.Error("failed to start server", "err", err)
		}
	}()
	nd.closers = append(nd.closers, server)

	return nd, nil
}

func (n *Node) Close() error {
	var err error
	for _, c := range n.closers {
		err = errors.Join(err, c.Close())
	}

	return err
}

type noOpBidProcessor struct{}

// ProcessBid auto accepts all bids sent.
func (noOpBidProcessor) ProcessBid(
	_ context.Context,
	_ *preconfpb.Bid,
) (chan providerapiv1.BidResponse_Status, error) {
	statusC := make(chan providerapiv1.BidResponse_Status, 5)
	statusC <- providerapiv1.BidResponse_STATUS_ACCEPTED
	close(statusC)

	return statusC, nil
}

type noOpCommitmentDA struct{}

func (noOpCommitmentDA) StoreCommitment(
	_ context.Context,
	_ *big.Int,
	_ uint64,
	_ string,
	_ uint64,
	_ uint64,
	_ []byte,
	_ []byte,
) error {
	return nil
}

func (noOpCommitmentDA) Close() error {
	return nil
}
