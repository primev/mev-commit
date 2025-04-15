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

	// Only a single target block number can be buffered, the notifier overwrites if the buffer is full
	targetBlockNumChan := make(chan uint64, 1)

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

		notifySecondsAhead := 7 * time.Second

		notif = notifier.NewFullNotifier(
			config.Logger.With("module", "full_notifier"),
			l1WsClient,
			notifySecondsAhead,
			targetBlockNumChan, // send-and-receive for draining capability
		)
	} else {
		if len(config.BeaconApiUrls) == 0 {
			return nil, fmt.Errorf("no beacon API URLs provided")
		}
		notificationsCli := notificationsapiv1.NewNotificationsClient(conn)
		config.Logger.Debug("created notifications client")

		notif = notifier.NewSelectiveNotifier(
			config.Logger.With("module", "selective_notifier"),
			notificationsCli,
			config.BeaconApiUrls[0],
			targetBlockNumChan, // send-and-receive for draining capability
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
		targetBlockNumChan, // receive-only
		acceptedBidChan,    // send-only
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

	err = balanceChecker.CheckBalances(ctx)
	if err != nil {
		return nil, err
	}
	config.Logger.Info("keystore account has enough balance on L1 and mev-commit chain")

	status, err := bidderCli.AutoDepositStatus(context.Background(), &bidderapiv1.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	config.Logger.Info("got auto deposit status", "enabled", status.IsAutodepositEnabled)

	if !status.IsAutodepositEnabled {
		config.Logger.Info("enabling auto deposit")
		resp, err := bidderCli.AutoDeposit(
			context.Background(),
			&bidderapiv1.DepositRequest{
				Amount: config.AutoDepositAmount.String(),
			},
		)
		if err != nil {
			return nil, err
		}
		config.Logger.Debug("auto deposit enabled", "amount", resp.AmountPerWindow, "window", resp.StartWindowNumber)
	}

	healthChecker := health.New()

	notifierDone := notif.Start(ctx)
	bidderDone := bidder.Start(ctx)
	monitorDone := monitor.Start(ctx)
	balanceCheckerDone := balanceChecker.Start(ctx)

	healthChecker.Register(health.CloseChannelHealthCheck("NotifierService", notifierDone))
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))
	healthChecker.Register(health.CloseChannelHealthCheck("MonitorService", monitorDone))
	healthChecker.Register(health.CloseChannelHealthCheck("BalanceCheckerService", balanceCheckerDone))

	s.closers = append(s.closers,
		channelCloser(notifierDone),
		channelCloser(bidderDone),
		channelCloser(monitorDone),
		channelCloser(balanceCheckerDone),
	)

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
