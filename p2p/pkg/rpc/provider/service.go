package providerapi

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProcessedBidResponse struct {
	Status            providerapiv1.BidResponse_Status
	DispatchTimestamp int64
}

type Service struct {
	providerapiv1.UnimplementedProviderServer
	receiver         chan *providerapiv1.Bid
	bidsInProcess    map[string]func(ProcessedBidResponse)
	bidsMu           sync.Mutex
	logger           *slog.Logger
	owner            common.Address
	registryContract ProviderRegistryContract
	watcher          Watcher
	optsGetter       OptsGetter
	metrics          *metrics
	validator        *protovalidate.Validator
}

type ProviderRegistryContract interface {
	RegisterAndStake(opts *bind.TransactOpts, blsPublicKey []byte) (*types.Transaction, error)
	CheckStake(*bind.CallOpts, common.Address) (*big.Int, error)
	MinStake(*bind.CallOpts) (*big.Int, error)
	ParseFundsDeposited(types.Log) (*providerregistry.ProviderregistryFundsDeposited, error)
}

type Watcher interface {
	WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error)
}

type OptsGetter func(ctx context.Context) (*bind.TransactOpts, error)

func NewService(
	logger *slog.Logger,
	registryContract ProviderRegistryContract,
	owner common.Address,
	watcher Watcher,
	optsGetter OptsGetter,
	validator *protovalidate.Validator,
) *Service {
	return &Service{
		receiver:         make(chan *providerapiv1.Bid),
		bidsInProcess:    make(map[string]func(ProcessedBidResponse)),
		registryContract: registryContract,
		owner:            owner,
		logger:           logger,
		watcher:          watcher,
		optsGetter:       optsGetter,
		metrics:          newMetrics(),
		validator:        validator,
	}
}

func toString(bid *providerapiv1.Bid) string {
	return fmt.Sprintf(
		"{TxHash: %v, BidAmount: %s, BlockNumber: %d, BidDigest: %x}",
		bid.TxHashes, bid.BidAmount, bid.BlockNumber, bid.BidDigest,
	)
}

func (s *Service) ProcessBid(
	ctx context.Context,
	bid *preconfpb.Bid,
) (chan ProcessedBidResponse, error) {
	bidMsg := &providerapiv1.Bid{
		TxHashes:            strings.Split(bid.TxHash, ","),
		BidAmount:           bid.BidAmount,
		BlockNumber:         bid.BlockNumber,
		BidDigest:           bid.Digest,
		DecayStartTimestamp: bid.DecayStartTimestamp,
		DecayEndTimestamp:   bid.DecayEndTimestamp,
	}

	err := s.validator.Validate(bidMsg)
	if err != nil {
		return nil, err
	}

	respC := make(chan ProcessedBidResponse, 1)
	s.bidsMu.Lock()
	s.bidsInProcess[string(bid.Digest)] = func(bidResponse ProcessedBidResponse) {
		respC <- ProcessedBidResponse{
			Status:            bidResponse.Status,
			DispatchTimestamp: bidResponse.DispatchTimestamp,
		}
		close(respC)
	}
	s.bidsMu.Unlock()

	select {
	case <-ctx.Done():
		s.bidsMu.Lock()
		delete(s.bidsInProcess, string(bid.Digest))
		s.bidsMu.Unlock()

		s.logger.Error("context cancelled for sending bid", "err", ctx.Err())
		return nil, ctx.Err()
	case s.receiver <- bidMsg:
	}
	s.logger.Info("sent bid to provider node", "bid", bid)

	return respC, nil
}

func (s *Service) ReceiveBids(
	_ *providerapiv1.EmptyMessage,
	srv providerapiv1.Provider_ReceiveBidsServer,
) error {
	for {
		select {
		case <-srv.Context().Done():
			s.logger.Error("context cancelled for receiving bid", "err", srv.Context().Err())
			return srv.Context().Err()
		case bid := <-s.receiver:
			s.logger.Info("received bid from node", "bid", toString(bid))
			err := srv.Send(bid)
			if err != nil {
				return err
			}
			s.metrics.BidsSentToProviderCount.Inc()
		}
	}
}

func (s *Service) SendProcessedBids(srv providerapiv1.Provider_SendProcessedBidsServer) error {
	for {
		status, err := srv.Recv()
		if err != nil {
			s.logger.Error("bid status", "err", err)
			return err
		}

		err = s.validator.Validate(status)
		if err != nil {
			s.logger.Error("bid status validation", "err", err)
			return err
		}

		s.bidsMu.Lock()
		callback, ok := s.bidsInProcess[string(status.BidDigest)]
		delete(s.bidsInProcess, string(status.BidDigest))
		s.bidsMu.Unlock()

		if ok {
			s.logger.Info(
				"received bid status from node",
				"bidDigest", hex.EncodeToString(status.BidDigest),
				"status", status.Status.String(),
			)
			callback(ProcessedBidResponse{
				Status:            status.Status,
				DispatchTimestamp: status.DispatchTimestamp,
			})
			if status.Status == providerapiv1.BidResponse_STATUS_ACCEPTED {
				s.metrics.BidsAcceptedByProviderCount.Inc()
			} else {
				s.metrics.BidsRejectedByProviderCount.Inc()
			}
		}
	}
}

var ErrInvalidAmount = errors.New("invalid amount for stake")

func (s *Service) RegisterStake(
	ctx context.Context,
	stake *providerapiv1.StakeRequest,
) (*providerapiv1.StakeResponse, error) {
	err := s.validator.Validate(stake)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate stake request: %v", err)
	}

	amount, success := big.NewInt(0).SetString(stake.Amount, 10)
	if !success {
		return nil, status.Errorf(codes.InvalidArgument, "parsing amount: %v", stake.Amount)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}
	opts.Value = amount

	tx, err := s.registryContract.RegisterAndStake(opts, []byte(stake.BlsPublicKey))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registering stake: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if registration, err := s.registryContract.ParseFundsDeposited(*log); err == nil {
			s.logger.Info("stake registered", "amount", registration.Amount)
			return &providerapiv1.StakeResponse{Amount: registration.Amount.String()}, nil
		}
	}

	s.logger.Error("no registration event found")
	return nil, status.Error(codes.Internal, "no registration event found")
}

func (s *Service) GetStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.StakeResponse, error) {
	stakeAmount, err := s.registryContract.CheckStake(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting stake: %v", err)
	}

	return &providerapiv1.StakeResponse{Amount: stakeAmount.String()}, nil
}

func (s *Service) GetMinStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.StakeResponse, error) {
	stakeAmount, err := s.registryContract.MinStake(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting min stake: %v", err)
	}

	return &providerapiv1.StakeResponse{Amount: stakeAmount.String()}, nil
}
