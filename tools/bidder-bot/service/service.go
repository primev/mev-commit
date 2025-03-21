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
	SettlementRPCUrl  string
	GasTipCap         *big.Int
	GasFeeCap         *big.Int
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
	if status == nil {
		for i := 0; i < 5; i++ {
			config.Logger.Debug("got empty auto deposit status, trying again", "attempt", i+1)
			time.Sleep(10 * time.Second)
			status, err = bidderCli.AutoDepositStatus(context.Background(), &bidderapiv1.EmptyMessage{})
			if err != nil {
				return nil, err
			}
			if status != nil {
				break
			}
			if i == 4 {
				return nil, errors.New("got empty auto deposit status")
			}
		}
	}
	config.Logger.Debug("got auto deposit status", "status", status)

	if !status.IsAutodepositEnabled {
		_, err := bidderCli.AutoDeposit(
			context.Background(),
			&bidderapiv1.DepositRequest{
				Amount: config.AutoDepositAmount.String(),
			},
		)
		if err != nil {
			return nil, err
		}
	}

	bidderClient := bidder.NewBidderClient(
		config.Logger.With("module", "bidder"),
		bidderCli,
		topologyCli,
		notificationsCli,
		l1RPCClient,
	)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	healthChecker := health.New()

	bidderDone := bidderClient.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))
	s.closers = append(s.closers, channelCloser(bidderDone))

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
