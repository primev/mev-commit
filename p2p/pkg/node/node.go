package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	depositmanagercontract "github.com/primev/mev-commit/contracts-abi/clients/DepositManager"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	contracts "github.com/primev/mev-commit/contracts-abi/config"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/apiserver"
	preconfstore "github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	bidderapi "github.com/primev/mev-commit/p2p/pkg/rpc/bidder"
	"github.com/primev/mev-commit/p2p/pkg/setcode"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	pebblestorage "github.com/primev/mev-commit/p2p/pkg/storage/pebble"
	"github.com/primev/mev-commit/p2p/pkg/txnstore"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type L1URLs struct {
	L1RPCURL     string
	BeaconAPIURL string
}

type noOpPreconfSender struct{}

func (noOpPreconfSender) SendBid(ctx context.Context, _ *preconfpb.Bid) (chan *preconfpb.PreConfirmation, error) {
	return nil, fmt.Errorf("preconfirmation disabled in minimal local mode")
}

type Options struct {
	Version                  string
	DataDir                  string
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
	BlockTrackerContract     string
	ProviderRegistryContract string
	BidderRegistryContract   string
	OracleContract           string
	ValidatorRouterContract  string
	EnableDepositManager     bool
	TargetDepositAmount      *big.Int
	RPCEndpoint              string
	WSRPCEndpoint            string
	NatAddr                  string
	TLSCertificateFile       string
	TLSPrivateKeyFile        string
	ProviderWhitelist        []common.Address
	DefaultGasLimit          uint64
	DefaultGasTipCap         *big.Int
	DefaultGasFeeCap         *big.Int
	BeaconAPIURL             string
	L1RPCURL                 string
	LaggardMode              *big.Int
	BidderBidTimeout         time.Duration
	ProviderDecisionTimeout  time.Duration
	NotificationsBufferCap   int
	ProposerNotifyOffset     time.Duration
	SlotDuration             time.Duration
	SlotsPerEpoch            uint64
}

type Node struct {
	cancelFunc context.CancelFunc
	closers    []io.Closer
}

func NewNode(opts *Options) (*Node, error) {
	nd := &Node{
		closers: make([]io.Closer, 0),
	}

	srv := apiserver.New(opts.Version, opts.Logger.With("component", "apiserver"))

	var (
		contractRPC *ethclient.Client
		err         error
	)
	if opts.WSRPCEndpoint != "" {
		contractRPC, err = ethclient.Dial(opts.WSRPCEndpoint)
		if err != nil {
			opts.Logger.Error("failed to connect to ws rpc", "error", err)
			return nil, err
		}
	} else {
		contractRPC, err = ethclient.Dial(opts.RPCEndpoint)
		if err != nil {
			opts.Logger.Error("failed to connect to rpc", "error", err)
			return nil, err
		}
	}

	chainID, err := contractRPC.ChainID(context.Background())
	if err != nil {
		opts.Logger.Error("failed to get chain ID", "error", err)
		return nil, err
	}

	if defaults, ok := contracts.DefaultsContracts[chainID.String()]; ok {
		setDefault(&opts.PreconfContract, defaults.PreconfManager)
		setDefault(&opts.BlockTrackerContract, defaults.BlockTracker)
		setDefault(&opts.ProviderRegistryContract, defaults.ProviderRegistry)
		setDefault(&opts.BidderRegistryContract, defaults.BidderRegistry)
		setDefault(&opts.OracleContract, defaults.Oracle)
	}

	var store storage.Storage
	if opts.DataDir != "" {
		store, err = pebblestorage.New(opts.DataDir)
		if err != nil {
			opts.Logger.Error("failed to create storage", "error", err)
			return nil, err
		}
	} else {
		store = inmem.New()
	}
	nd.closers = append(nd.closers, store)

	contractsABIs, err := getContractABIs(opts)
	if err != nil {
		opts.Logger.Error("failed to get contract ABIs", "error", err)
		return nil, err
	}

	txnStore := txnstore.New(store)
	monitor := txmonitor.New(
		opts.KeySigner.GetAddress(),
		contractRPC,
		txmonitor.NewEVMHelperWithLogger(contractRPC, opts.Logger.With("component", "txmonitor"), contractsABIs),
		txnStore,
		opts.Logger.With("component", "txmonitor"),
		1024,
	)
	srv.RegisterMetricsCollectors(monitor.Metrics()...)

	var startables []StartableObjWithDesc
	startables = append(
		startables,
		StartableObjWithDesc{
			Desc:      "txmonitor",
			Startable: monitor,
		},
	)

	backend := transactor.NewMetricsWrapper(
		transactor.NewTransactor(
			contractRPC,
			monitor,
		),
	)
	srv.RegisterMetricsCollectors(backend.Metrics()...)

	providerRegistry, err := providerregistry.NewProviderregistry(
		common.HexToAddress(opts.ProviderRegistryContract),
		backend,
	)
	if err != nil {
		opts.Logger.Error("failed to instantiate provider registry contract", "error", err)
		return nil, err
	}

	bidderRegistry, err := bidderregistry.NewBidderregistry(
		common.HexToAddress(opts.BidderRegistryContract),
		backend,
	)
	if err != nil {
		opts.Logger.Error("failed to instantiate bidder registry contract", "error", err)
		return nil, err
	}

	depositManagerContract, err := depositmanagercontract.NewDepositmanager(
		opts.KeySigner.GetAddress(), // EOA bound
		backend,
	)
	if err != nil {
		opts.Logger.Error("creating deposit manager", "error", err)
		return nil, err
	}

	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		tOpts, err := opts.KeySigner.GetAuthWithCtx(ctx, chainID)
		if err == nil {
			tOpts.GasLimit = opts.DefaultGasLimit
			tOpts.GasTipCap = opts.DefaultGasTipCap
			tOpts.GasFeeCap = opts.DefaultGasFeeCap
		}
		return tOpts, err
	}

	lis, err := net.Listen("tcp", opts.RPCAddr)
	if err != nil {
		opts.Logger.Error("failed to listen", "error", err)
		return nil, errors.Join(err, nd.Close())
	}

	grpcServer := grpc.NewServer()

	validator, err := protovalidate.New()
	if err != nil {
		opts.Logger.Error("failed to create proto validator", "error", err)
		return nil, errors.Join(err, nd.Close())
	}

	preconfStore := preconfstore.New(store)

	setCodeHelper := setcode.NewSetCodeHelper(
		opts.Logger.With("component", "setcode_helper"),
		opts.KeySigner,
		backend,
		chainID,
	)

	depositManagerImplAddr, err := bidderRegistry.DepositManagerImpl(nil)
	if err != nil {
		opts.Logger.Error("failed to get deposit manager implementation address", "error", err)
		return nil, err
	}

	bidderAPI := bidderapi.NewService(
		opts.KeySigner.GetAddress(),
		noOpPreconfSender{}, // no preconfirmation in minimal mode
		bidderRegistry,
		providerRegistry,
		validator,
		monitor, // tx watcher
		optsGetter,
		preconfStore,
		opts.BidderBidTimeout,
		opts.Logger.With("component", "bidderapi"),
		setCodeHelper,
		depositManagerContract,
		contractRPC, // backend (balance/code queries)
		nil,         // topology (nil: avoid GetValidProviders)
		depositManagerImplAddr,
	)
	bidderapiv1.RegisterBidderServer(grpcServer, bidderAPI)
	srv.RegisterMetricsCollectors(bidderAPI.Metrics()...)

	if opts.EnableDepositManager {
		if err := handleEnableDepositManager(bidderAPI, opts); err != nil {
			opts.Logger.Error("failed to enable deposit manager", "error", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	nd.cancelFunc = cancel

	healthChecker := health.New()

	for _, s := range startables {
		closeChan := s.Startable.Start(ctx)
		healthChecker.Register(health.CloseChannelHealthCheck(s.Desc, closeChan))
		nd.closers = append(nd.closers, channelCloserFunc(closeChan))
	}

	started := make(chan struct{})
	go func() {
		close(started)
		if err := grpcServer.Serve(lis); err != nil {
			opts.Logger.Error("failed to start grpc server", "err", err)
		}
	}()
	nd.closers = append(nd.closers, lis)
	<-started

	grpcConn, err := grpc.DialContext(
		context.Background(),
		opts.RPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("dialing grpc server failed: %w", err), nd.Close())
	}
	healthChecker.Register(health.GrpcGatewayHealthCheck(grpcConn))

	handlerCtx, handlerCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer handlerCancel()

	gatewayMux := runtime.NewServeMux()
	if err := bidderapiv1.RegisterBidderHandler(handlerCtx, gatewayMux, grpcConn); err != nil {
		opts.Logger.Error("failed to register bidder handler", "err", err)
		return nil, errors.Join(err, nd.Close())
	}

	srv.ChainHandlers("/", gatewayMux)
	srv.ChainHandlers(
		"/health",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				if err := healthChecker.Health(); err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintln(w, "ok")
			},
		),
	)

	server := &http.Server{
		Addr:    opts.HTTPAddr,
		Handler: srv.Router(),
	}

	go func() {
		opts.Logger.Info("starting to listen", "tls", false)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			opts.Logger.Error("failed to start server", "err", err)
		}
	}()
	nd.closers = append(nd.closers, server)

	return nd, nil
}

func getContractABIs(opts *Options) (map[common.Address]*abi.ABI, error) {
	abis := make(map[common.Address]*abi.ABI)

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.BlockTrackerContract)] = &btABI

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.PreconfContract)] = &pcABI

	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.BidderRegistryContract)] = &brABI

	vrABI, err := abi.JSON(strings.NewReader(validatorrouter.ValidatoroptinrouterABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.ValidatorRouterContract)] = &vrABI

	orABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.OracleContract)] = &orABI

	prABI, err := abi.JSON(strings.NewReader(providerregistry.ProviderregistryABI))
	if err != nil {
		return nil, err
	}
	abis[common.HexToAddress(opts.ProviderRegistryContract)] = &prABI

	return abis, nil
}

func (n *Node) Close() error {
	if n.cancelFunc != nil {
		n.cancelFunc()
	}

	var err error
	for i := len(n.closers) - 1; i >= 0; i-- {
		err = errors.Join(err, n.closers[i].Close())
	}

	return err
}

type PublisherStartable interface {
	Start(ctx context.Context, contracts ...common.Address) <-chan struct{}
}

type StartableObjWithDesc struct {
	Startable Startable
	Desc      string
}

type Startable interface {
	Start(ctx context.Context) <-chan struct{}
}

type StartableFunc func(ctx context.Context) <-chan struct{}

func (f StartableFunc) Start(ctx context.Context) <-chan struct{} {
	return f(ctx)
}

func setDefault(field *string, defaultValue string) {
	if *field == "" {
		*field = defaultValue
	}
}

type channelCloser <-chan struct{}

func channelCloserFunc(c <-chan struct{}) io.Closer {
	return channelCloser(c)
}

func (c channelCloser) Close() error {
	select {
	case <-c:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("timeout waiting for channel to close")
	}
}

func handleEnableDepositManager(bidderAPI *bidderapi.Service, opts *Options) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	enableDepositMngrResp, err := bidderAPI.EnableDepositManager(ctx, &bidderapiv1.EnableDepositManagerRequest{})
	if err != nil && !strings.Contains(err.Error(),
		"EnableDepositManager failed: deposit manager is already enabled") {
		return fmt.Errorf("failed to enable deposit manager: %w", err)
	}
	if !enableDepositMngrResp.Success {
		return fmt.Errorf("failed to enable deposit manager: %w", err)
	}
	opts.Logger.Info("deposit manager enabled")
	if opts.TargetDepositAmount != nil {
		providers, err := bidderAPI.GetValidProviders(ctx, &bidderapiv1.GetValidProvidersRequest{})
		if err != nil {
			return fmt.Errorf("failed to get valid providers: %w", err)
		}
		opts.Logger.Info("valid providers who'll be deposited to", "providers", providers.ValidProviders)
		setTargetDepositsReq := &bidderapiv1.SetTargetDepositsRequest{}
		for _, provider := range providers.ValidProviders {
			setTargetDepositsReq.TargetDeposits = append(setTargetDepositsReq.TargetDeposits,
				&bidderapiv1.TargetDeposit{
					Provider:      provider,
					TargetDeposit: opts.TargetDepositAmount.String(),
				},
			)
		}

		setTargetDepositsResp, err := bidderAPI.SetTargetDeposits(ctx, setTargetDepositsReq)
		if err != nil {
			return fmt.Errorf("failed to set target deposit amount: %w", err)
		}
		if len(setTargetDepositsResp.SuccessfullySetDeposits) < len(providers.ValidProviders) {
			return fmt.Errorf("failed to set target deposit amount for all valid providers: %w", err)
		}
		opts.Logger.Info("target deposit amount set for all valid providers", "amount", opts.TargetDepositAmount)
		opts.Logger.Info("successfully topped up providers", "providers", setTargetDepositsResp.SuccessfullyToppedUpProviders)
	}
	return nil
}
