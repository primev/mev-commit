package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
	"github.com/primev/mev-commit/tools/bidder-bot/monitor"
	notifier "github.com/primev/mev-commit/tools/bidder-bot/notifier"
	"github.com/primev/mev-commit/x/contracts/ethwrapper"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Logger            *slog.Logger
	Signer            keysigner.KeySigner
	BidderNodeRPC     string
	AutoDepositAmount *big.Int
	L1RPCUrls         []string
	L1WsUrls          []string
	BeaconApiUrls     []string
	SettlementRPCUrl  string
	GasTipCap         *big.Int
	GasFeeCap         *big.Int
	BidAmount         *big.Int
	IsFullNotifier    bool
	BlockInterval     uint64
	CheckBalances     bool
}

type Service struct {
	cancel  context.CancelFunc
	closers []io.Closer
}

func New(config *Config) (*Service, error) {
	s := &Service{}

	opts := []grpc.DialOption{}
	if strings.HasPrefix(config.BidderNodeRPC, "https://") {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{InsecureSkipVerify: true},
		)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(config.BidderNodeRPC, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	config.Logger.Debug("created gRPC connection", "address", config.BidderNodeRPC)

	s.closers = append(s.closers, conn)

	l1RPCClient, err := ethwrapper.NewClient(
		config.Logger.With("module", "ethwrapper"),
		config.L1RPCUrls,
		ethwrapper.EthClientWithMaxRetries(5),
	)
	if err != nil {
		return nil, err
	}
	config.Logger.Debug("created ethwrapper client", "urls", config.L1RPCUrls)

	settlementRPCClient, err := ethwrapper.NewClient(
		config.Logger.With("module", "ethwrapper"),
		[]string{config.SettlementRPCUrl},
		ethwrapper.EthClientWithMaxRetries(5),
	)
	if err != nil {
		return nil, err
	}

	bidderCli := bidderapiv1.NewBidderClient(conn)
	config.Logger.Debug("created bidder client")
	topologyCli := debugapiv1.NewDebugServiceClient(conn)
	config.Logger.Debug("created topology client")

	targetBlockChan := make(chan bidder.TargetBlock, 1)
	acceptedBidChan := make(chan *monitor.AcceptedBid, 5)

	type Notifier interface {
		Start(ctx context.Context) <-chan struct{}
	}
	var notif Notifier

	if config.IsFullNotifier {
		if len(config.L1WsUrls) == 0 {
			return nil, fmt.Errorf("no L1 WebSocket URLs provided")
		}
		l1WsClient, err := ethclient.Dial(config.L1WsUrls[0])
		if err != nil {
			return nil, fmt.Errorf("failed to create L1 WebSocket client: %w", err)
		}
		config.Logger.Debug("created L1 WebSocket client", "url", config.L1WsUrls[0])

		notif = notifier.NewFullNotifier(
			config.Logger.With("module", "full_notifier"),
			l1WsClient,
			targetBlockChan, // send-and-receive for draining capability
			config.BlockInterval,
		)
	} else {
		if len(config.BeaconApiUrls) == 0 {
			return nil, fmt.Errorf("no beacon API URLs provided")
		}
		notificationsCli := notificationsapiv1.NewNotificationsClient(conn)
		config.Logger.Debug("created notifications client")

		beaconClient := notifier.NewBeaconClient(config.BeaconApiUrls[0], config.Logger.With("component", "beacon_client"))

		notif = notifier.NewSelectiveNotifier(
			config.Logger.With("module", "selective_notifier"),
			notificationsCli,
			beaconClient,
			targetBlockChan, // send-and-receive for draining capability
		)
	}

	bidder := bidder.NewBidder(
		config.Logger.With("module", "bidder"),
		bidderCli,
		topologyCli,
		l1RPCClient,
		config.Signer,
		config.GasTipCap,
		config.GasFeeCap,
		config.BidAmount,
		targetBlockChan, // receive-only
		acceptedBidChan, // send-only
	)

	monitorTxLandingTimeout := 15 * time.Minute
	monitorTxLandingInterval := 30 * time.Second

	monitor := monitor.NewMonitor(
		config.Logger.With("module", "monitor"),
		l1RPCClient,
		acceptedBidChan, // receive-only
		monitorTxLandingTimeout,
		monitorTxLandingInterval,
	)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	balanceChecker := NewBalanceChecker(
		config.Logger.With("module", "balance_checker"),
		config.Signer,
		l1RPCClient,
		settlementRPCClient,
	)

	if config.CheckBalances {
		err = balanceChecker.CheckBalances(ctx)
		if err != nil {
			return nil, err
		}
		config.Logger.Info("keystore account has enough balance on L1 and mev-commit chain")
	} else {
		config.Logger.Info("balance checking disabled")
	}

	status, err := bidderCli.DepositManagerStatus(context.Background(), &bidderapiv1.DepositManagerStatusRequest{})
	if err != nil {
		return nil, err
	}
	if !status.Enabled {
		resp, err := bidderCli.EnableDepositManager(context.Background(), &bidderapiv1.EnableDepositManagerRequest{})
		if err != nil {
			return nil, err
		}
		if !resp.Success {
			return nil, errors.New("failed to enable deposit manager")
		}
	}
	config.Logger.Info("deposit manager enabled")

	validProviders, err := bidderCli.GetValidProviders(context.Background(), &bidderapiv1.GetValidProvidersRequest{})
	if err != nil {
		return nil, err
	}
	if len(validProviders.ValidProviders) == 0 {
		return nil, errors.New("no connected and valid providers found")
	}

	targetDeposits := make([]*bidderapiv1.TargetDeposit, len(validProviders.ValidProviders))
	for i, provider := range validProviders.ValidProviders {
		targetDeposits[i] = &bidderapiv1.TargetDeposit{
			Provider:      provider,
			TargetDeposit: config.AutoDepositAmount.String(),
		}
	}

	resp, err := bidderCli.SetTargetDeposits(context.Background(), &bidderapiv1.SetTargetDepositsRequest{
		TargetDeposits: targetDeposits,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.SuccessfullySetDeposits) != len(targetDeposits) {
		return nil, errors.New("failed to set target deposits")
	}

	healthChecker := health.New()

	notifierDone := notif.Start(ctx)
	bidderDone := bidder.Start(ctx)

	healthChecker.Register(health.CloseChannelHealthCheck("NotifierService", notifierDone))
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))

	s.closers = append(s.closers,
		channelCloser(notifierDone),
		channelCloser(bidderDone),
	)

	if config.CheckBalances {
		balanceCheckerDone := balanceChecker.Start(ctx)
		healthChecker.Register(health.CloseChannelHealthCheck("BalanceCheckerService", balanceCheckerDone))
		s.closers = append(s.closers, channelCloser(balanceCheckerDone))
	}

	if !config.IsFullNotifier {
		monitorDone := monitor.Start(ctx)
		healthChecker.Register(health.CloseChannelHealthCheck("MonitorService", monitorDone))
		s.closers = append(s.closers, channelCloser(monitorDone))
	}

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
