package service

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/instant-bridge/api"
	"github.com/primev/mev-commit/x/accountsync"
	"github.com/primev/mev-commit/x/contracts/ethwrapper"
	"github.com/primev/mev-commit/x/health"
	"github.com/primev/mev-commit/x/keysigner"
	bidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"github.com/primev/mev-commit/x/transfer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Logger                 *slog.Logger
	Signer                 keysigner.KeySigner
	BidderRPC              string
	AutoDepositAmount      *big.Int
	L1RPCUrls              []string
	SettlementRPCUrl       string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
	SettlementThreshold    *big.Int
	SettlementTopup        *big.Int
	HTTPPort               int
	MinServiceFee          *big.Int
	GasTipCap              *big.Int
	GasFeeCap              *big.Int
}

type Service struct {
	cancel  context.CancelFunc
	closers []io.Closer
}

func New(config *Config) (*Service, error) {
	s := &Service{}

	conn, err := grpc.NewClient(
		config.BidderRPC,
		grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{InsecureSkipVerify: true},
		)),
	)
	if err != nil {
		return nil, err
	}

	s.closers = append(s.closers, conn)

	l1RPCClient, err := ethwrapper.NewClient(
		config.Logger.With("module", "ethwrapper"),
		config.L1RPCUrls,
		ethwrapper.EthClientWithMaxRetries(5),
	)
	if err != nil {
		return nil, err
	}
	l1ChainID, err := l1RPCClient.RawClient().ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	settlementClient, err := ethclient.Dial(config.SettlementRPCUrl)
	if err != nil {
		return nil, err
	}
	settlementChainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	bidderCli := bidderapiv1.NewBidderClient(conn)
	topologyCli := debugapiv1.NewDebugServiceClient(conn)
	notificationsCli := notificationsapiv1.NewNotificationsClient(conn)

	// TODO: set code to deposit manager here, set min deposit for every provider

	// status, err := bidderCli.AutoDepositStatus(context.Background(), &bidderapiv1.EmptyMessage{})
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !status.IsAutodepositEnabled {
	// 	_, err := bidderCli.AutoDeposit(
	// 		context.Background(),
	// 		&bidderapiv1.DepositRequest{
	// 			Amount: config.AutoDepositAmount.String(),
	// 		},
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	bridgeConfig := transfer.BridgeConfig{
		Signer:                 config.Signer,
		L1ContractAddr:         config.L1ContractAddr,
		SettlementContractAddr: config.SettlementContractAddr,
		L1RPCUrl:               config.L1RPCUrls[0],
		SettlementRPCUrl:       config.SettlementRPCUrl,
	}

	syncer := accountsync.NewAccountSync(config.Signer.GetAddress(), settlementClient)
	bridger := transfer.NewBridger(
		config.Logger.With("module", "bridger"),
		syncer,
		bridgeConfig,
		config.SettlementThreshold,
		config.SettlementTopup,
	)

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

	bridgerDone := bridger.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("Bridger", bridgerDone))
	s.closers = append(s.closers, channelCloser(bridgerDone))

	bidderDone := bidderClient.Start(ctx)
	healthChecker.Register(health.CloseChannelHealthCheck("BidderService", bidderDone))
	s.closers = append(s.closers, channelCloser(bidderDone))

	transferer := transfer.NewTransferer(
		config.Logger.With("module", "transferer"),
		settlementClient,
		config.Signer,
		config.GasTipCap,
		config.GasFeeCap,
	)

	apiService := api.NewAPI(
		config.Logger.With("module", "api"),
		config.HTTPPort,
		healthChecker,
		bidderClient,
		transferer,
		config.MinServiceFee,
		config.Signer.GetAddress(),
		l1RPCClient.RawClient(),
		settlementClient,
		l1ChainID,
		settlementChainID,
	)

	apiService.Start()
	s.closers = append(s.closers, apiService)

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
