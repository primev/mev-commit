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
	"sync/atomic"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	preconfstore "github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProcessedBidResponse struct {
	Status            providerapiv1.BidResponse_Status
	DispatchTimestamp int64
}

type Service struct {
	providerapiv1.UnimplementedProviderServer
	receiver               chan *providerapiv1.Bid
	bidsInProcess          map[string]func(ProcessedBidResponse)
	bidsMu                 sync.Mutex
	logger                 *slog.Logger
	owner                  common.Address
	registryContract       ProviderRegistryContract
	bidderRegistryContract BidderRegistryContract
	watcher                Watcher
	cs                     CommitmentStore
	optsGetter             OptsGetter
	metrics                *metrics
	validator              *protovalidate.Validator
	activeReceivers        atomic.Int32
}

type BidderRegistryContract interface {
	GetProviderAmount(*bind.CallOpts, common.Address) (*big.Int, error)
	WithdrawProviderAmount(*bind.TransactOpts, common.Address) (*types.Transaction, error)
}

type ProviderRegistryContract interface {
	ProviderRegistered(opts *bind.CallOpts, address common.Address) (bool, error)
	Stake(opts *bind.TransactOpts) (*types.Transaction, error)
	RegisterAndStake(opts *bind.TransactOpts) (*types.Transaction, error)
	AddVerifiedBLSKey(opts *bind.TransactOpts, blsPublicKey []byte, signature []byte) (*types.Transaction, error)
	GetProviderStake(*bind.CallOpts, common.Address) (*big.Int, error)
	GetBLSKeys(*bind.CallOpts, common.Address) ([][]byte, error)
	MinStake(*bind.CallOpts) (*big.Int, error)
	ParseProviderRegistered(types.Log) (*providerregistry.ProviderregistryProviderRegistered, error)
	ParseFundsDeposited(types.Log) (*providerregistry.ProviderregistryFundsDeposited, error)
	ParseUnstake(types.Log) (*providerregistry.ProviderregistryUnstake, error)
	ParseWithdraw(types.Log) (*providerregistry.ProviderregistryWithdraw, error)
	ParseBLSKeyAdded(types.Log) (*providerregistry.ProviderregistryBLSKeyAdded, error)
	Withdraw(opts *bind.TransactOpts) (*types.Transaction, error)
	Unstake(opts *bind.TransactOpts) (*types.Transaction, error)
}

type Watcher interface {
	WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error)
}

type CommitmentStore interface {
	GetCommitments(blockNumber int64) ([]*preconfstore.Commitment, error)
	ListCommitments(opts *preconfstore.ListOpts) ([]*preconfstore.Commitment, error)
}

type OptsGetter func(ctx context.Context) (*bind.TransactOpts, error)

func NewService(
	logger *slog.Logger,
	registryContract ProviderRegistryContract,
	bidderRegistryContract BidderRegistryContract,
	owner common.Address,
	watcher Watcher,
	cs CommitmentStore,
	optsGetter OptsGetter,
	validator *protovalidate.Validator,
) *Service {
	return &Service{
		receiver:               make(chan *providerapiv1.Bid),
		bidsInProcess:          make(map[string]func(ProcessedBidResponse)),
		registryContract:       registryContract,
		bidderRegistryContract: bidderRegistryContract,
		owner:                  owner,
		logger:                 logger,
		watcher:                watcher,
		cs:                     cs,
		optsGetter:             optsGetter,
		metrics:                newMetrics(),
		validator:              validator,
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
	if s.activeReceivers.Load() == 0 {
		return nil, status.Error(codes.Internal, "no active receivers")
	}
	var revertingTxnHashes []string
	if bid.RevertingTxHashes != "" {
		revertingTxnHashes = strings.Split(bid.RevertingTxHashes, ",")
	}
	bidMsg := &providerapiv1.Bid{
		TxHashes:            strings.Split(bid.TxHash, ","),
		BidAmount:           bid.BidAmount,
		SlashAmount:         bid.SlashAmount,
		BlockNumber:         bid.BlockNumber,
		BidDigest:           bid.Digest,
		DecayStartTimestamp: bid.DecayStartTimestamp,
		DecayEndTimestamp:   bid.DecayEndTimestamp,
		RevertingTxHashes:   revertingTxnHashes,
		RawTransactions:     bid.RawTransactions,
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
	s.activeReceivers.Add(1)
	defer s.activeReceivers.Add(-1)

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

func (s *Service) Stake(
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

	var (
		tx    *types.Transaction
		txErr error
	)

	registered, err := s.registryContract.ProviderRegistered(&bind.CallOpts{Context: ctx, From: s.owner}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "checking registration: %v", err)
	}

	if !registered {
		if len(stake.BlsPublicKeys) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing BLS keys")
		}
		if len(stake.BlsSignatures) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing BLS signatures")
		}
		tx, txErr = s.registryContract.RegisterAndStake(opts)
	} else {
		tx, txErr = s.registryContract.Stake(opts)
	}
	if txErr != nil {
		return nil, status.Errorf(codes.Internal, "failed to stake: %v", txErr)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt for registration: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	var stakeResponse *providerapiv1.StakeResponse
	for _, log := range receipt.Logs {
		if registration, err := s.registryContract.ParseProviderRegistered(*log); err == nil {
			s.logger.Info("stake registered", "amount", registration.StakedAmount)
			stakeResponse = &providerapiv1.StakeResponse{
				Amount: registration.StakedAmount.String(),
			}
		} else if deposit, err := s.registryContract.ParseFundsDeposited(*log); err == nil {
			s.logger.Info("stake deposited", "amount", deposit.Amount)
			stakeResponse = &providerapiv1.StakeResponse{Amount: deposit.Amount.String()}
		}
	}

	for i := range stake.BlsPublicKeys {
		blsPublicKey, err := hex.DecodeString(strings.TrimPrefix(stake.BlsPublicKeys[i], "0x"))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "decoding bls public key: %v", err)
		}
		blsSignature, err := hex.DecodeString(strings.TrimPrefix(stake.BlsSignatures[i], "0x"))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "decoding bls signature: %v", err)
		}

		opts, err = s.optsGetter(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "getting transact opts for adding BLS key: %v", err)
		}

		s.logger.Info("adding verified bls key", "blsPublicKey", hex.EncodeToString(blsPublicKey), "blsSignature", hex.EncodeToString(blsSignature))
		tx, txErr = s.registryContract.AddVerifiedBLSKey(opts, blsPublicKey, blsSignature)
		if txErr != nil {
			return nil, status.Errorf(codes.Internal, "adding verified bls key: %v", txErr)
		}
		receipt, err = s.watcher.WaitForReceipt(ctx, tx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "waiting for receipt for adding verified bls key: %v", err)
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
		}

		s.logger.Info("verified bls key added", "key", stake.BlsPublicKeys[i])

		for _, log := range receipt.Logs {
			if blsKeyEvent, err := s.registryContract.ParseBLSKeyAdded(*log); err == nil {
				s.logger.Info("verified bls key added", "key", blsKeyEvent.BlsPublicKey)
				stakeResponse.BlsPublicKeys = append(stakeResponse.BlsPublicKeys, hex.EncodeToString(blsKeyEvent.BlsPublicKey))
			}
		}
	}

	return stakeResponse, nil
}

func (s *Service) GetStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.StakeResponse, error) {
	stakeAmount, err := s.registryContract.GetProviderStake(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting stake: %v", err)
	}

	blsPubkey, err := s.registryContract.GetBLSKeys(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting bls public key: %v", err)
	}

	encodedKeys := make([]string, len(blsPubkey))
	for i, key := range blsPubkey {
		encodedKeys[i] = hex.EncodeToString(key)
	}
	return &providerapiv1.StakeResponse{Amount: stakeAmount.String(), BlsPublicKeys: encodedKeys}, nil
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

func (s *Service) WithdrawStake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.WithdrawalResponse, error) {
	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	tx, err := s.registryContract.Withdraw(opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "withdrawing stake: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseWithdraw(*log); err == nil {
			s.logger.Info("stake withdrawn", "amount", withdrawal.Amount)
			return &providerapiv1.WithdrawalResponse{Amount: withdrawal.Amount.String()}, nil
		}
	}

	s.logger.Error("no withdrawal event found")
	return nil, status.Error(codes.Internal, "no withdrawal event found")
}

func (s *Service) Unstake(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.EmptyMessage, error) {
	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	tx, err := s.registryContract.Unstake(opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "requesting withdrawal: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseUnstake(*log); err == nil {
			s.logger.Info("withdrawal requested", "timestamp", withdrawal.Timestamp)
			return &providerapiv1.EmptyMessage{}, nil
		}
	}

	s.logger.Error("no withdrawal request event found")
	return nil, status.Error(codes.Internal, "no withdrawal request event found")
}

func (s *Service) WithdrawProviderReward(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.WithdrawalResponse, error) {
	// Get the amount before withdrawal to report in response
	amount, err := s.bidderRegistryContract.GetProviderAmount(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "checking provider rewards: %v", err)
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return &providerapiv1.WithdrawalResponse{Amount: "0"}, nil
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	tx, err := s.bidderRegistryContract.WithdrawProviderAmount(opts, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "withdrawing provider rewards: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "transaction failed with status: %v", receipt.Status)
	}

	s.logger.Info("provider rewards withdrawn", "amount", amount.String())
	return &providerapiv1.WithdrawalResponse{Amount: amount.String()}, nil
}

func (s *Service) GetProviderReward(
	ctx context.Context,
	_ *providerapiv1.EmptyMessage,
) (*providerapiv1.RewardResponse, error) {
	// Get the provider's accumulated reward amount
	amount, err := s.bidderRegistryContract.GetProviderAmount(&bind.CallOpts{
		Context: ctx,
		From:    s.owner,
	}, s.owner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "checking provider rewards: %v", err)
	}

	s.logger.Info("retrieved provider reward amount", "amount", amount.String())
	return &providerapiv1.RewardResponse{Amount: amount.String()}, nil
}

const (
	defaultPage  = 0
	defaultLimit = 100
)

func (s *Service) GetCommitmentInfo(
	ctx context.Context,
	req *providerapiv1.GetCommitmentInfoRequest,
) (*providerapiv1.CommitmentInfoResponse, error) {
	var (
		cmts        []*preconfstore.Commitment
		err         error
		page, limit = int(req.Page), int(req.Limit)
	)

	if limit == 0 {
		limit = defaultLimit
	}

	if req.BlockNumber != 0 {
		cmts, err = s.cs.GetCommitments(req.BlockNumber)
	} else {
		cmts, err = s.cs.ListCommitments(&preconfstore.ListOpts{
			Page:  page,
			Limit: limit,
		})
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting commitments: %v", err)
	}
	if len(cmts) == 0 {
		return &providerapiv1.CommitmentInfoResponse{
			Commitments: []*providerapiv1.CommitmentInfoResponse_BlockCommitments{},
		}, nil
	}

	blockCommitments := make([]*providerapiv1.CommitmentInfoResponse_BlockCommitments, 0)
	for _, c := range cmts {
		if len(blockCommitments) == 0 || blockCommitments[len(blockCommitments)-1].BlockNumber != c.Bid.BlockNumber {
			blockCommitments = append(blockCommitments, &providerapiv1.CommitmentInfoResponse_BlockCommitments{
				BlockNumber: c.Bid.BlockNumber,
				Commitments: make([]*providerapiv1.CommitmentInfoResponse_Commitment, 0),
			})
		}
		blockCommitments[len(blockCommitments)-1].Commitments = append(blockCommitments[len(blockCommitments)-1].Commitments, &providerapiv1.CommitmentInfoResponse_Commitment{
			TxnHashes:           strings.Split(c.Bid.TxHash, ","),
			RevertableTxnHashes: strings.Split(c.Bid.RevertingTxHashes, ","),
			Amount:              c.Bid.BidAmount,
			BlockNumber:         c.Bid.BlockNumber,
			ProviderAddress:     common.Bytes2Hex(c.ProviderAddress),
			DecayStartTimestamp: c.Bid.DecayStartTimestamp,
			DecayEndTimestamp:   c.Bid.DecayEndTimestamp,
			DispatchTimestamp:   c.EncryptedPreConfirmation.DispatchTimestamp,
			SlashAmount:         c.Bid.SlashAmount,
			Status:              string(c.Status),
			Details:             c.Details,
			Payment:             c.Payment,
			Refund:              c.Refund,
		})
	}

	return &providerapiv1.CommitmentInfoResponse{
		Commitments: blockCommitments,
	}, nil
}
