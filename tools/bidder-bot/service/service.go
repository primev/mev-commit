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

	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
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
	BeaconApiUrls     []string
	SettlementRPCUrl  string
	GasTipCap         *big.Int
	GasFeeCap         *big.Int
	BidAmount         *big.Int
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
	notificationsCli := notificationsapiv1.NewNotificationsClient(conn)
	config.Logger.Debug("created notifications client")

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

	// Only a single upcomingProposer can be buffered, the notifier overwrites if the buffer is full
	proposerChan := make(chan *notifier.UpcomingProposer, 1)

	notifier := notifier.NewNotifier(
		config.Logger.With("module", "notifier"),
		notificationsCli,
		proposerChan, // send-and-receive for draining capability
	)

	if len(config.BeaconApiUrls) == 0 {
		return nil, fmt.Errorf("no beacon API URLs provided")
	}

	bidder := bidder.NewBidder(
		config.Logger.With("module", "bidder"),
		bidderCli,
		topologyCli,
		l1RPCClient,
		config.BeaconApiUrls[0],
		config.Signer,
		config.GasTipCap,
		config.GasFeeCap,
		config.BidAmount,
		proposerChan, // receive-only
	)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	err = s.checkBalances(ctx, config.Signer, l1RPCClient, settlementRPCClient)
	if err != nil {
		return nil, err
	}
	config.Logger.Info("keystore account has enough balance on L1 and mev-commit chain")

	healthChecker := health.New()

	notifierDone := notifier.Start(ctx)
	bidderDone := bidder.Start(ctx)

	healthChecker.Register(health.CloseChannelHealthCheck("NotifierService", notifierDone))
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))

	s.closers = append(s.closers,
		channelCloser(notifierDone),
		channelCloser(bidderDone),
	)

	return s, nil
}

func (s *Service) checkBalances(ctx context.Context, signer keysigner.KeySigner, l1RPCClient *ethwrapper.Client, settlementRPCClient *ethwrapper.Client) error {
	l1Balance, err := l1RPCClient.RawClient().BalanceAt(ctx, signer.GetAddress(), nil)
	if err != nil {
		return err
	}
	pointZeroFiveEth := big.NewInt(50000000000000000)
	if l1Balance.Cmp(pointZeroFiveEth) < 0 {
		return fmt.Errorf("keystore account has less than 0.05 eth on L1")
	}

	settlementBalance, err := settlementRPCClient.RawClient().BalanceAt(ctx, signer.GetAddress(), nil)
	if err != nil {
		return err
	}
	if settlementBalance.Cmp(pointZeroFiveEth) < 0 {
		return fmt.Errorf("keystore account has less than 0.05 eth on mev-commit chain")
	}
	return nil
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
