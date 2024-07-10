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
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfirmationv1 "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
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
	blockTrackerContract BlockTrackerContract
	watcher              TxWatcher
	optsGetter           OptsGetter
	autoDepositTracker   AutoDepositTracker
	logger               *slog.Logger
	metrics              *metrics
	validator            *protovalidate.Validator
}

func NewService(
	owner common.Address,
	blocksPerWindow uint64,
	sender PreconfSender,
	registryContract BidderRegistryContract,
	blockTrackerContract BlockTrackerContract,
	validator *protovalidate.Validator,
	watcher TxWatcher,
	optsGetter OptsGetter,
	autoDepositTracker AutoDepositTracker,
	logger *slog.Logger,
) *Service {
	return &Service{
		owner:                owner,
		blocksPerWindow:      blocksPerWindow,
		sender:               sender,
		registryContract:     registryContract,
		blockTrackerContract: blockTrackerContract,
		watcher:              watcher,
		optsGetter:           optsGetter,
		logger:               logger,
		metrics:              newMetrics(),
		autoDepositTracker:   autoDepositTracker,
		validator:            validator,
	}
}

type AutoDepositTracker interface {
	Start(context.Context, *big.Int, *big.Int) error
	Stop() ([]*big.Int, error)
	IsWorking() bool
	GetStatus() (map[uint64]bool, bool)
}

type PreconfSender interface {
	SendBid(ctx context.Context, txnsStr string, amount string, blockNumber int64, decayStartTimestamp int64, decayEndTimestamp int64, revertingTxHashes string) (chan *preconfirmationv1.PreConfirmation, error)
}

type BidderRegistryContract interface {
	DepositForWindow(*bind.TransactOpts, *big.Int) (*types.Transaction, error)
	WithdrawBidderAmountFromWindow(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)
	GetDeposit(*bind.CallOpts, common.Address, *big.Int) (*big.Int, error)
	WithdrawFromWindows(*bind.TransactOpts, []*big.Int) (*types.Transaction, error)
	ParseBidderRegistered(types.Log) (*bidderregistry.BidderregistryBidderRegistered, error)
	ParseBidderWithdrawal(types.Log) (*bidderregistry.BidderregistryBidderWithdrawal, error)
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
	ctx, cancel := context.WithTimeout(srv.Context(), 10*time.Second)
	defer cancel()

	s.metrics.ReceivedBidsCount.Inc()

	err := s.validator.Validate(bid)
	if err != nil {
		s.logger.Error("bid validation", "error", err)
		return status.Errorf(codes.InvalidArgument, "validating bid: %v", err)
	}

	txnsStr := strings.Join(bid.TxHashes, ",")
	revertingTxHashesStr := strings.Join(bid.RevertingTxHashes, ",")

	respC, err := s.sender.SendBid(
		ctx,
		txnsStr,
		bid.Amount,
		bid.BlockNumber,
		bid.DecayStartTimestamp,
		bid.DecayEndTimestamp,
		revertingTxHashesStr,
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

func (s *Service) Deposit(
	ctx context.Context,
	r *bidderapiv1.DepositRequest,
) (*bidderapiv1.DepositResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validating deposit request: %v", err)
	}

	if s.autoDepositTracker.IsWorking() {
		return nil, status.Error(codes.FailedPrecondition, "auto deposit is already running, stop and then deposit")
	}

	currentWindow, err := s.blockTrackerContract.GetCurrentWindow()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
	}

	windowToDeposit, err := s.calculateWindowToDeposit(ctx, r, currentWindow.Uint64())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "calculating window to deposit: %v", err)
	}

	amount, success := big.NewInt(0).SetString(r.Amount, 10)
	if !success {
		return nil, status.Errorf(codes.InvalidArgument, "parsing amount: %v", r.Amount)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}
	opts.Value = amount

	tx, err := s.registryContract.DepositForWindow(opts, windowToDeposit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deposit: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if registration, err := s.registryContract.ParseBidderRegistered(*log); err == nil {
			s.logger.Info("deposit successful", "amount", registration.DepositedAmount, "window", registration.WindowNumber)
			return &bidderapiv1.DepositResponse{
				Amount:       registration.DepositedAmount.String(),
				WindowNumber: wrapperspb.UInt64(registration.WindowNumber.Uint64()),
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
	// Default to two windows ahead of the current window if no specific block or window is given.
	// This is for the case where the oracle works 2 windows behind the current window.
	return new(big.Int).SetUint64(currentWindow + 2), nil
}

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
			return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
		}
		// as oracle working 2 windows behind the current window, we add + 2 here
		window = new(big.Int).Add(window, big.NewInt(2))
	} else {
		window = new(big.Int).SetUint64(r.WindowNumber.Value)
	}
	stakeAmount, err := s.registryContract.GetDeposit(&bind.CallOpts{
		From:    s.owner,
		Context: ctx,
	}, s.owner, window)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting deposit: %v", err)
	}

	return &bidderapiv1.DepositResponse{Amount: stakeAmount.String(), WindowNumber: wrapperspb.UInt64(window.Uint64())}, nil
}

func (s *Service) Withdraw(
	ctx context.Context,
	r *bidderapiv1.WithdrawRequest,
) (*bidderapiv1.WithdrawResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validating withdraw request: %v", err)
	}

	if s.autoDepositTracker.IsWorking() {
		return nil, status.Error(codes.FailedPrecondition, "auto deposit is already running, stop and then withdraw")
	}

	var window *big.Int
	if r.WindowNumber == nil {
		window, err = s.blockTrackerContract.GetCurrentWindow()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
		}
		window = new(big.Int).Sub(window, big.NewInt(1))
	} else {
		window = new(big.Int).SetUint64(r.WindowNumber.Value)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	tx, err := s.registryContract.WithdrawBidderAmountFromWindow(opts, s.owner, window)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "withdrawing deposit: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseBidderWithdrawal(*log); err == nil {
			s.logger.Info("withdrawal successful", "amount", withdrawal.Amount.String(), "window", withdrawal.Window.String())
			return &bidderapiv1.WithdrawResponse{
				Amount:       withdrawal.Amount.String(),
				WindowNumber: wrapperspb.UInt64(withdrawal.Window.Uint64()),
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
	err := s.validator.Validate(r)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validating auto deposit request: %v", err)
	}

	if s.autoDepositTracker.IsWorking() {
		return nil, status.Error(codes.FailedPrecondition, "auto deposit is already running")
	}

	currentWindow, err := s.blockTrackerContract.GetCurrentWindow()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting current window: %v", err)
	}

	windowToDeposit, err := s.calculateWindowToDeposit(ctx, r, currentWindow.Uint64())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "calculating window to deposit: %v", err)
	}

	amount, success := big.NewInt(0).SetString(r.Amount, 10)
	if !success {
		return nil, status.Errorf(codes.InvalidArgument, "parsing amount: %v", r.Amount)
	}

	err = s.autoDepositTracker.Start(ctx, windowToDeposit, amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "starting auto deposit: %v", err)
	}

	s.logger.Error(
		"autodeposit enabled",
		"window", windowToDeposit,
		"amount", amount.String(),
	)

	return &bidderapiv1.AutoDepositResponse{
		StartBlockNumber: wrapperspb.UInt64(windowToDeposit.Uint64()),
		AmountPerWindow:  amount.String(),
	}, nil
}

func (s *Service) CancelAutoDeposit(
	ctx context.Context,
	r *bidderapiv1.CancelAutoDepositRequest,
) (*bidderapiv1.CancelAutoDepositResponse, error) {
	if !s.autoDepositTracker.IsWorking() {
		return nil, status.Error(codes.FailedPrecondition, "auto deposit is not running")
	}
	windows, err := s.autoDepositTracker.Stop()
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "cancel auto deposit: %v", err)
	}
	if r.Withdraw {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				currentWindow, err := s.blockTrackerContract.GetCurrentWindow()
				if err != nil {
					s.logger.Error("getting current window", "error", err)
					continue
				}
				doWithdraw := true
				for _, w := range windows {
					if w.Uint64() >= currentWindow.Uint64() {
						doWithdraw = false
						break
					}
				}
				if doWithdraw {
					opts, err := s.optsGetter(context.Background())
					if err != nil {
						s.logger.Error("getting transact opts", "error", err)
						continue
					}
					txn, err := s.registryContract.WithdrawFromWindows(opts, windows)
					if err != nil {
						s.logger.Error("withdraw from windows", "error", err)
						return
					}
					receipt, err := s.watcher.WaitForReceipt(context.Background(), txn)
					if err != nil {
						s.logger.Error("waiting for receipt", "error", err)
						return
					}
					if receipt.Status != types.ReceiptStatusSuccessful {
						s.logger.Error("receipt status", "status", receipt.Status)
					}
					return
				}
				s.logger.Info("waiting for windows to be in the past before withdrawing", "currentWindow", currentWindow, "windows", windows)
			}
		}()
		return &bidderapiv1.CancelAutoDepositResponse{}, nil
	}

	withdrawWindows := []*wrapperspb.UInt64Value{}
	for _, w := range windows {
		withdrawWindows = append(withdrawWindows, wrapperspb.UInt64(w.Uint64()))
	}

	return &bidderapiv1.CancelAutoDepositResponse{
		WindowNumbers: withdrawWindows,
	}, nil
}

func (s *Service) WithdrawFromWindows(
	ctx context.Context,
	r *bidderapiv1.WithdrawFromWindowsRequest,
) (*bidderapiv1.WithdrawFromWindowsResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validating withdraw from n windows request: %v", err)
	}

	if s.autoDepositTracker.IsWorking() {
		return nil, status.Error(codes.FailedPrecondition, "auto deposit is already running, stop and then withdraw")
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	windows := make([]*big.Int, len(r.WindowNumbers))
	for i, w := range r.WindowNumbers {
		windows[i] = new(big.Int).SetUint64(w.Value)
	}

	tx, err := s.registryContract.WithdrawFromWindows(opts, windows)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "withdrawing deposit: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	var amountsAndWindows []*bidderapiv1.WithdrawResponse
	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseBidderWithdrawal(*log); err == nil {
			s.logger.Info("withdrawal successful", "amount", withdrawal.Amount.String(), "window", withdrawal.Window.String())
			amountsAndWindows = append(amountsAndWindows, &bidderapiv1.WithdrawResponse{
				Amount:       withdrawal.Amount.String(),
				WindowNumber: wrapperspb.UInt64(withdrawal.Window.Uint64()),
			})
		}
	}

	if len(amountsAndWindows) > 0 {
		return &bidderapiv1.WithdrawFromWindowsResponse{
			WithdrawResponses: amountsAndWindows,
		}, nil
	}

	s.logger.Error(
		"withdraw successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"window", r.WindowNumbers,
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for withdrawal")
}

func (s *Service) AutoDepositStatus(
	ctx context.Context,
	_ *bidderapiv1.EmptyMessage,
) (*bidderapiv1.AutoDepositStatusResponse, error) {
	deposits, isWorking := s.autoDepositTracker.GetStatus()
	var autoDeposits []*bidderapiv1.AutoDeposit
	for window, ok := range deposits {
		if ok {
			stakeAmount, err := s.registryContract.GetDeposit(&bind.CallOpts{
				From:    s.owner,
				Context: ctx,
			}, s.owner, new(big.Int).SetUint64(window))
			if err != nil {
				return nil, status.Errorf(codes.Internal, "getting deposit: %v", err)
			}
			autoDeposits = append(autoDeposits, &bidderapiv1.AutoDeposit{
				WindowNumber: wrapperspb.UInt64(window),
				Amount:       stakeAmount.String(),
			})
		}
	}

	return &bidderapiv1.AutoDepositStatusResponse{
		WindowBalances: autoDeposits,
		IsWorking:      isWorking,
	}, nil
}
