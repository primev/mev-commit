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
	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	registrycontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/provider_registry"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	providerapiv1.UnimplementedProviderServer
	receiver         chan *providerapiv1.Bid
	bidsInProcess    map[string]func(providerapiv1.ProcessedBidResponse)
	bidsMu           sync.Mutex
	logger           *slog.Logger
	owner            common.Address
	registryContract registrycontract.Interface
	evmClient        EvmClient
	metrics          *metrics
	validator        *protovalidate.Validator
}

type EvmClient interface {
	PendingTxns() []evmclient.TxnInfo
	CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error)
}

func NewService(
	logger *slog.Logger,
	registryContract registrycontract.Interface,
	owner common.Address,
	e EvmClient,
	validator *protovalidate.Validator,
) *Service {
	return &Service{
		receiver:         make(chan *providerapiv1.Bid),
		bidsInProcess:    make(map[string]func(providerapiv1.ProcessedBidResponse)),
		registryContract: registryContract,
		owner:            owner,
		logger:           logger,
		evmClient:        e,
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
) (chan providerapiv1.ProcessedBidResponse, error) {
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

	respC := make(chan providerapiv1.ProcessedBidResponse, 1)
	s.bidsMu.Lock()
	s.bidsInProcess[string(bid.Digest)] = func(bidResponse providerapiv1.ProcessedBidResponse) {
		respC <- providerapiv1.ProcessedBidResponse{
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
			callback(providerapiv1.ProcessedBidResponse{
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

	err = s.registryContract.RegisterProvider(ctx, amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registering stake: %v", err)
	}

	stakeAmount, err := s.registryContract.GetStake(ctx, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting stake: %v", err)
	}

	return &providerapiv1.StakeResponse{Amount: stakeAmount.String()}, nil
}

func (s *Service) GetStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.StakeResponse, error) {
	stakeAmount, err := s.registryContract.GetStake(ctx, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting stake: %v", err)
	}

	return &providerapiv1.StakeResponse{Amount: stakeAmount.String()}, nil
}

func (s *Service) GetMinStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.StakeResponse, error) {
	stakeAmount, err := s.registryContract.GetMinStake(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting min stake: %v", err)
	}

	return &providerapiv1.StakeResponse{Amount: stakeAmount.String()}, nil
}

func (s *Service) GetPendingTxns(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.PendingTxnsResponse, error) {
	txns := s.evmClient.PendingTxns()

	txnsMsg := make([]*providerapiv1.TransactionInfo, len(txns))
	for i, txn := range txns {
		txnsMsg[i] = &providerapiv1.TransactionInfo{
			TxHash:  txn.Hash,
			Nonce:   int64(txn.Nonce),
			Created: txn.Created,
		}
	}

	return &providerapiv1.PendingTxnsResponse{PendingTxns: txnsMsg}, nil
}

func (s *Service) CancelTransaction(
	ctx context.Context,
	cancel *providerapiv1.CancelReq,
) (*providerapiv1.CancelResponse, error) {
	txHash := common.HexToHash(cancel.TxHash)
	cHash, err := s.evmClient.CancelTx(ctx, txHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cancelling transaction: %v", err)
	}

	return &providerapiv1.CancelResponse{TxHash: cHash.Hex()}, nil
}
