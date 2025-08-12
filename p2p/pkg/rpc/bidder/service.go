package bidderapi

import (
	"context"
	"encoding/hex"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfirmationv1 "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	preconfstore "github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Service struct {
	bidderapiv1.UnimplementedBidderServer
	owner                common.Address
	blocksPerWindow      uint64
	sender               PreconfSender
	registryContract     BidderRegistryContract
	providerRegistry     ProviderRegistryContract
	blockTrackerContract BlockTrackerContract
	watcher              TxWatcher
	optsGetter           OptsGetter
	cs                   CommitmentStore
	oracleWindowOffset   *big.Int
	logger               *slog.Logger
	metrics              *metrics
	validator            *protovalidate.Validator
	bidTimeout           time.Duration
}

func NewService(
	owner common.Address,
	sender PreconfSender,
	registryContract BidderRegistryContract,
	blockTrackerContract BlockTrackerContract,
	providerRegistry ProviderRegistryContract,
	validator *protovalidate.Validator,
	watcher TxWatcher,
	optsGetter OptsGetter,
	cs CommitmentStore,
	oracleWindowOffset *big.Int,
	bidderBidTimeout time.Duration,
	logger *slog.Logger,
) *Service {
	return &Service{
		owner:                owner,
		sender:               sender,
		registryContract:     registryContract,
		blockTrackerContract: blockTrackerContract,
		providerRegistry:     providerRegistry,
		cs:                   cs,
		watcher:              watcher,
		optsGetter:           optsGetter,
		logger:               logger,
		metrics:              newMetrics(),
		oracleWindowOffset:   oracleWindowOffset,
		validator:            validator,
		bidTimeout:           bidderBidTimeout,
	}
}

type PreconfSender interface {
	SendBid(ctx context.Context, bid *preconfirmationv1.Bid) (chan *preconfirmationv1.PreConfirmation, error)
}

type BidderRegistryContract interface {
	DepositAsBidder(*bind.TransactOpts, common.Address) (*types.Transaction, error)
	RequestWithdrawalsAsBidder(*bind.TransactOpts, []common.Address) (*types.Transaction, error)
	WithdrawAsBidder(*bind.TransactOpts, []common.Address) (*types.Transaction, error)
	GetDeposit(*bind.CallOpts, common.Address, common.Address) (*big.Int, error)
	ParseBidderDeposited(types.Log) (*bidderregistry.BidderregistryBidderDeposited, error)
	ParseWithdrawalRequested(types.Log) (*bidderregistry.BidderregistryWithdrawalRequested, error)
	ParseBidderWithdrawal(types.Log) (*bidderregistry.BidderregistryBidderWithdrawal, error)
}

type ProviderRegistryContract interface {
	BidderSlashedAmount(*bind.CallOpts, common.Address) (*big.Int, error)
	WithdrawSlashedAmount(*bind.TransactOpts) (*types.Transaction, error)
	ParseBidderWithdrawSlashedAmount(log types.Log) (*providerregistry.ProviderregistryBidderWithdrawSlashedAmount, error)
}

type CommitmentStore interface {
	GetCommitments(blockNumber int64) ([]*preconfstore.Commitment, error)
	ListCommitments(opts *preconfstore.ListOpts) ([]*preconfstore.Commitment, error)
}

type BlockTrackerContract interface {
	GetCurrentWindow() (*big.Int, error)
}

type TxWatcher interface {
	WaitForReceipt(context.Context, *types.Transaction) (*types.Receipt, error)
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

func (s *Service) SendBid(
	bid *bidderapiv1.Bid,
	srv bidderapiv1.Bidder_SendBidServer,
) error {
	// timeout to prevent hanging of bidder node if provider node is not responding
	ctx, cancel := context.WithTimeout(srv.Context(), s.bidTimeout)
	defer cancel()

	s.metrics.ReceivedBidsCount.Inc()

	err := s.validator.Validate(bid)
	if err != nil {
		s.logger.Error("bid validation", "error", err)
		return status.Errorf(codes.InvalidArgument, "validating bid: %v", err)
	}

	switch {
	case len(bid.TxHashes) == 0 && len(bid.RawTransactions) == 0:
		s.logger.Error("empty bid", "bid", bid)
		return status.Error(codes.InvalidArgument, "empty bid")
	case len(bid.TxHashes) > 0 && len(bid.RawTransactions) > 0:
		s.logger.Error("both txHashes and rawTransactions are provided", "bid", bid)
		return status.Error(codes.InvalidArgument, "both txHashes and rawTransactions are provided")
	}

	// Helper function to strip "0x" prefix
	stripPrefix := func(hashes []string) []string {
		stripped := make([]string, len(hashes))
		for i, hash := range hashes {
			stripped[i] = strings.TrimPrefix(hash, "0x")
		}
		return stripped
	}
	var (
		txnsStr string
	)
	switch {
	case len(bid.TxHashes) > 0:
		txnsStr = strings.Join(stripPrefix(bid.TxHashes), ",")
	case len(bid.RawTransactions) > 0:
		strBuilder := new(strings.Builder)
		for i, rawTx := range bid.RawTransactions {
			rawTxnBytes, err := hex.DecodeString(strings.TrimPrefix(rawTx, "0x"))
			if err != nil {
				s.logger.Error("decoding raw transaction", "error", err)
				return status.Errorf(codes.InvalidArgument, "decoding raw transaction: %v", err)
			}
			txnObj := new(types.Transaction)
			err = txnObj.UnmarshalBinary(rawTxnBytes)
			if err != nil {
				s.logger.Error("unmarshaling raw transaction", "error", err)
				return status.Errorf(codes.InvalidArgument, "unmarshaling raw transaction: %v", err)
			}
			strBuilder.WriteString(strings.TrimPrefix(txnObj.Hash().Hex(), "0x"))
			if i != len(bid.RawTransactions)-1 {
				strBuilder.WriteString(",")
			}
		}
		txnsStr = strBuilder.String()
	}

	if bid.SlashAmount == "" {
		bid.SlashAmount = "0"
	}

	respC, err := s.sender.SendBid(
		ctx,
		&preconfirmationv1.Bid{
			TxHash:              txnsStr,
			BidAmount:           bid.Amount,
			SlashAmount:         bid.SlashAmount,
			BlockNumber:         bid.BlockNumber,
			DecayStartTimestamp: bid.DecayStartTimestamp,
			DecayEndTimestamp:   bid.DecayEndTimestamp,
			RevertingTxHashes:   strings.Join(stripPrefix(bid.RevertingTxHashes), ","),
			RawTransactions:     bid.RawTransactions,
		},
	)
	if err != nil {
		s.logger.Error("sending bid", "error", err)
		return status.Errorf(codes.Internal, "error sending bid: %v", err)
	}

	for resp := range respC {
		b := resp.Bid
		err := srv.Send(&bidderapiv1.Commitment{
			TxHashes:             strings.Split(b.TxHash, ","),
			BidAmount:            b.BidAmount,
			SlashAmount:          b.SlashAmount,
			BlockNumber:          b.BlockNumber,
			ReceivedBidDigest:    hex.EncodeToString(b.Digest),
			ReceivedBidSignature: hex.EncodeToString(b.Signature),
			CommitmentDigest:     hex.EncodeToString(resp.Digest),
			CommitmentSignature:  hex.EncodeToString(resp.Signature),
			ProviderAddress:      common.Bytes2Hex(resp.ProviderAddress),
			DecayStartTimestamp:  b.DecayStartTimestamp,
			DecayEndTimestamp:    b.DecayEndTimestamp,
			DispatchTimestamp:    resp.DispatchTimestamp,
			RevertingTxHashes:    strings.Split(b.RevertingTxHashes, ","),
		})
		if err != nil {
			s.logger.Error("sending preConfirmation", "error", err)
			return err
		}
		s.metrics.ReceivedPreconfsCount.Inc()
	}

	return nil
}

// TODO: update to new semantics, needs provider addr
func (s *Service) Deposit(
	ctx context.Context,
	r *bidderapiv1.DepositRequest,
) (*bidderapiv1.DepositResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("deposit validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating deposit request: %v", err)
	}

	currentWindow, err := s.blockTrackerContract.GetCurrentWindow()
	if err != nil {
		s.logger.Error("getting current window", "error", err)
		return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
	}

	windowToDeposit, err := s.calculateWindowToDeposit(ctx, r, currentWindow.Uint64())
	if err != nil {
		s.logger.Error("calculating window to deposit", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "calculating window to deposit: %v", err)
	}

	amount, success := big.NewInt(0).SetString(r.Amount, 10)
	if !success {
		s.logger.Error("parsing amount", "amount", r.Amount)
		return nil, status.Errorf(codes.InvalidArgument, "parsing amount: %v", r.Amount)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}
	opts.Value = amount

	providerAddr := common.HexToAddress("0x0000000000000000000000000000000000000000") // TODO: get from request
	tx, err := s.registryContract.DepositAsBidder(opts, providerAddr)
	if err != nil {
		s.logger.Error("depositing", "error", err)
		return nil, status.Errorf(codes.Internal, "deposit: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		s.logger.Error("waiting for receipt", "error", err)
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		s.logger.Error("receipt status", "status", receipt.Status)
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if deposited, err := s.registryContract.ParseBidderDeposited(*log); err == nil {
			s.logger.Info("deposit successful", "amount", deposited.DepositedAmount.String())
			return &bidderapiv1.DepositResponse{
				Amount: deposited.DepositedAmount.String(),
			}, nil
		}
	}

	s.logger.Error(
		"deposit successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"window", windowToDeposit,
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for deposit")
}

func (s *Service) calculateWindowToDeposit(ctx context.Context, r *bidderapiv1.DepositRequest, currentWindow uint64) (*big.Int, error) {
	if r.WindowNumber != nil {
		// Directly use the specified window number if available.
		return new(big.Int).SetUint64(r.WindowNumber.Value), nil
	} else if r.BlockNumber != nil {
		return new(big.Int).SetUint64((r.BlockNumber.Value-1)/s.blocksPerWindow + 1), nil
	}
	// Default to N window ahead of the current window if no specific block or window is given.
	// This is for the case where the oracle works N windows behind the current window.
	return new(big.Int).SetUint64(currentWindow + s.oracleWindowOffset.Uint64()), nil
}

// TODO: update to new semantics, needs provider addr
func (s *Service) GetDeposit(
	ctx context.Context,
	r *bidderapiv1.GetDepositRequest,
) (*bidderapiv1.DepositResponse, error) {
	var (
		window *big.Int
		err    error
	)
	if r.WindowNumber == nil {
		window, err = s.blockTrackerContract.GetCurrentWindow()
		if err != nil {
			s.logger.Error("getting current window", "error", err)
			return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
		}
		// as oracle working N windows behind the current window, we add + N here
		window = new(big.Int).Add(window, s.oracleWindowOffset)
	} else {
		window = new(big.Int).SetUint64(r.WindowNumber.Value)
	}
	providerAddr := common.HexToAddress("0x0000000000000000000000000000000000000000") // TODO: get from request
	stakeAmount, err := s.registryContract.GetDeposit(&bind.CallOpts{
		From:    s.owner,
		Context: ctx,
	}, s.owner, providerAddr)
	if err != nil {
		s.logger.Error("getting deposit", "error", err)
		return nil, status.Errorf(codes.Internal, "getting deposit: %v", err)
	}

	return &bidderapiv1.DepositResponse{Amount: stakeAmount.String(), WindowNumber: wrapperspb.UInt64(window.Uint64())}, nil
}

func (s *Service) RequestWithdrawal(
	ctx context.Context,
	r *bidderapiv1.WithdrawRequest,
) (*bidderapiv1.WithdrawResponse, error) {
	// TODO:
	return nil, status.Errorf(codes.Unimplemented, "method RequestWithdrawal not implemented")
}

func (s *Service) Withdraw(
	ctx context.Context,
	r *bidderapiv1.WithdrawRequest,
) (*bidderapiv1.WithdrawResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("withdraw validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating withdraw request: %v", err)
	}

	var window *big.Int
	if r.WindowNumber == nil {
		window, err = s.blockTrackerContract.GetCurrentWindow()
		if err != nil {
			s.logger.Error("getting current window", "error", err)
			return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
		}
		window = new(big.Int).Sub(window, big.NewInt(1))
	} else {
		window = new(big.Int).SetUint64(r.WindowNumber.Value)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	providerAddr := common.HexToAddress("0x0000000000000000000000000000000000000000") // TODO: get from request
	providers := []common.Address{providerAddr}
	tx, err := s.registryContract.WithdrawAsBidder(opts, providers)
	if err != nil {
		s.logger.Error("withdrawing deposit", "error", err)
		return nil, status.Errorf(codes.Internal, "withdrawing deposit: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		s.logger.Error("waiting for receipt", "error", err)
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		s.logger.Error("receipt status", "status", receipt.Status)
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseBidderWithdrawal(*log); err == nil {
			s.logger.Info("withdrawal successful", "amount", withdrawal.AmountWithdrawn.String())
			return &bidderapiv1.WithdrawResponse{
				Amount: withdrawal.AmountWithdrawn.String(),
			}, nil
		}
	}

	s.logger.Error(
		"withdraw successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"window", window.Uint64(),
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for withdrawal")
}

func (s *Service) AutoDeposit(
	ctx context.Context,
	r *bidderapiv1.DepositRequest,
) (*bidderapiv1.AutoDepositResponse, error) {
	// TODO: Set EOA code
	return nil, status.Errorf(codes.Unimplemented, "method AutoDeposit not implemented")
}

func (s *Service) CancelAutoDeposit(
	ctx context.Context,
	r *bidderapiv1.CancelAutoDepositRequest,
) (*bidderapiv1.CancelAutoDepositResponse, error) {
	// TODO: Unset EOA code
	return nil, status.Errorf(codes.Unimplemented, "method CancelAutoDeposit not implemented")
}

func (s *Service) AutoDepositStatus(
	ctx context.Context,
	_ *bidderapiv1.EmptyMessage,
) (*bidderapiv1.AutoDepositStatusResponse, error) {
	// TODO: Query EOA code
	return nil, status.Errorf(codes.Unimplemented, "method AutoDepositStatus not implemented")
}

func (s *Service) ClaimSlashedFunds(
	ctx context.Context,
	_ *bidderapiv1.EmptyMessage,
) (*wrapperspb.StringValue, error) {
	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	amount, err := s.providerRegistry.BidderSlashedAmount(&bind.CallOpts{
		From:    s.owner,
		Context: ctx,
	}, s.owner)
	if err != nil {
		s.logger.Error("getting slashed amount", "error", err)
		return nil, status.Errorf(codes.Internal, "getting slashed amount: %v", err)
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		s.logger.Info("no slashed amount to claim")
		return &wrapperspb.StringValue{Value: "0"}, nil
	}

	tx, err := s.providerRegistry.WithdrawSlashedAmount(opts)
	if err != nil {
		s.logger.Error("withdrawing slashed amount", "error", err)
		return nil, status.Errorf(codes.Internal, "withdrawing slashed amount: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		s.logger.Error("waiting for receipt", "error", err)
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		s.logger.Error("receipt status", "status", receipt.Status)
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if withdrawal, err := s.providerRegistry.ParseBidderWithdrawSlashedAmount(*log); err == nil {
			s.logger.Info("slashed amount withdrawal successful", "amount", withdrawal.Amount.String())
			return &wrapperspb.StringValue{Value: withdrawal.Amount.String()}, nil
		}
	}

	s.logger.Error(
		"withdraw slashed amount successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"logs", receipt.Logs,
	)

	s.logger.Error("withdraw slashed amount successful but missing log", "txHash", receipt.TxHash.Hex(), "logs", receipt.Logs)
	return nil, status.Errorf(codes.Internal, "missing log for slashed amount withdrawal")
}

const (
	defaultLimit = 100
)

func (s *Service) GetBidInfo(
	ctx context.Context,
	req *bidderapiv1.GetBidInfoRequest,
) (*bidderapiv1.GetBidInfoResponse, error) {
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
		return &bidderapiv1.GetBidInfoResponse{}, nil
	}
	blockBids := make([]*bidderapiv1.GetBidInfoResponse_BlockBidInfo, 0)
LOOP:
	for _, c := range cmts {
		if len(blockBids) == 0 || blockBids[len(blockBids)-1].BlockNumber != c.Bid.BlockNumber {
			blockBids = append(blockBids, &bidderapiv1.GetBidInfoResponse_BlockBidInfo{
				BlockNumber: c.Bid.BlockNumber,
				Bids:        make([]*bidderapiv1.GetBidInfoResponse_BidInfo, 0),
			})
		}
		cmtWithStatus := &bidderapiv1.GetBidInfoResponse_CommitmentWithStatus{
			ProviderAddress:   common.Bytes2Hex(c.ProviderAddress),
			DispatchTimestamp: c.PreConfirmation.DispatchTimestamp,
			Status:            string(c.Status),
			Details:           c.Details,
			Payment:           c.Payment,
			Refund:            c.Refund,
		}
		for idx, b := range blockBids[len(blockBids)-1].Bids {
			if common.Bytes2Hex(c.Bid.Digest) == b.BidDigest {
				blockBids[len(blockBids)-1].Bids[idx].Commitments = append(blockBids[len(blockBids)-1].Bids[idx].Commitments, cmtWithStatus)
				continue LOOP
			}
		}
		blockBids[len(blockBids)-1].Bids = append(blockBids[len(blockBids)-1].Bids, &bidderapiv1.GetBidInfoResponse_BidInfo{
			BidDigest:           common.Bytes2Hex(c.Bid.Digest),
			TxnHashes:           strings.Split(c.Bid.TxHash, ","),
			RevertableTxnHashes: strings.Split(c.Bid.RevertingTxHashes, ","),
			BlockNumber:         c.Bid.BlockNumber,
			BidAmount:           c.Bid.BidAmount,
			SlashAmount:         c.Bid.SlashAmount,
			DecayStartTimestamp: c.Bid.DecayStartTimestamp,
			DecayEndTimestamp:   c.Bid.DecayEndTimestamp,
			Commitments: []*bidderapiv1.GetBidInfoResponse_CommitmentWithStatus{
				cmtWithStatus,
			},
		})
	}

	return &bidderapiv1.GetBidInfoResponse{
		BlockBidInfo: blockBids,
	}, nil
}
