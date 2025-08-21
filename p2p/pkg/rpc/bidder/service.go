package bidderapi

import (
	"context"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	depositmanager "github.com/primev/mev-commit/contracts-abi/clients/DepositManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfirmationv1 "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	preconfstore "github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Service struct {
	bidderapiv1.UnimplementedBidderServer
	owner                  common.Address
	sender                 PreconfSender
	registryContract       BidderRegistryContract
	providerRegistry       ProviderRegistryContract
	watcher                TxWatcher
	optsGetter             OptsGetter
	cs                     CommitmentStore
	logger                 *slog.Logger
	metrics                *metrics
	validator              *protovalidate.Validator
	bidTimeout             time.Duration
	setCodeHelper          SetCodeHelper
	depositManager         DepositManagerContract
	backend                Backend
	topology               *topology.Topology
	depositManagerImplAddr common.Address
}

func NewService(
	owner common.Address,
	sender PreconfSender,
	registryContract BidderRegistryContract,
	providerRegistry ProviderRegistryContract,
	validator *protovalidate.Validator,
	watcher TxWatcher,
	optsGetter OptsGetter,
	cs CommitmentStore,
	bidderBidTimeout time.Duration,
	logger *slog.Logger,
	setCodeHelper SetCodeHelper,
	depositManager DepositManagerContract,
	backend Backend,
	topology *topology.Topology,
	depositManagerImplAddr common.Address,
) *Service {
	return &Service{
		owner:                  owner,
		sender:                 sender,
		registryContract:       registryContract,
		providerRegistry:       providerRegistry,
		cs:                     cs,
		watcher:                watcher,
		optsGetter:             optsGetter,
		logger:                 logger,
		metrics:                newMetrics(),
		validator:              validator,
		bidTimeout:             bidderBidTimeout,
		setCodeHelper:          setCodeHelper,
		depositManager:         depositManager,
		backend:                backend,
		topology:               topology,
		depositManagerImplAddr: depositManagerImplAddr,
	}
}

type PreconfSender interface {
	SendBid(ctx context.Context, bid *preconfirmationv1.Bid) (chan *preconfirmationv1.PreConfirmation, error)
}

type BidderRegistryContract interface {
	DepositAsBidder(*bind.TransactOpts, common.Address) (*types.Transaction, error)
	DepositEvenlyAsBidder(*bind.TransactOpts, []common.Address) (*types.Transaction, error)
	RequestWithdrawalsAsBidder(*bind.TransactOpts, []common.Address) (*types.Transaction, error)
	WithdrawAsBidder(*bind.TransactOpts, []common.Address) (*types.Transaction, error)
	GetDeposit(*bind.CallOpts, common.Address, common.Address) (*big.Int, error)
	ParseBidderDeposited(types.Log) (*bidderregistry.BidderregistryBidderDeposited, error)
	ParseWithdrawalRequested(types.Log) (*bidderregistry.BidderregistryWithdrawalRequested, error)
	ParseBidderWithdrawal(types.Log) (*bidderregistry.BidderregistryBidderWithdrawal, error)
	FilterBidderDeposited(opts *bind.FilterOpts, bidder []common.Address, provider []common.Address, depositedAmount []*big.Int) (*bidderregistry.BidderregistryBidderDepositedIterator, error)
	OpenBid(
		opts *bind.TransactOpts,
		commitmentDigest [32]byte,
		bidAmt *big.Int,
		bidderAddr common.Address,
		providerAddr common.Address,
	) (*types.Transaction, error)
}

type ProviderRegistryContract interface {
	BidderSlashedAmount(*bind.CallOpts, common.Address) (*big.Int, error)
	WithdrawSlashedAmount(*bind.TransactOpts) (*types.Transaction, error)
	ParseBidderWithdrawSlashedAmount(log types.Log) (*providerregistry.ProviderregistryBidderWithdrawSlashedAmount, error)
	AreProvidersValid(*bind.CallOpts, []common.Address) ([]bool, error)
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

type SetCodeHelper interface {
	SetCode(ctx context.Context, opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
}

type DepositManagerContract interface {
	SetTargetDeposits(opts *bind.TransactOpts, providers []common.Address, amounts []*big.Int) (*types.Transaction, error)
	TopUpDeposits(opts *bind.TransactOpts, providers []common.Address) (*types.Transaction, error)
	ParseTargetDepositSet(types.Log) (*depositmanager.DepositmanagerTargetDepositSet, error)
	ParseDepositToppedUp(types.Log) (*depositmanager.DepositmanagerDepositToppedUp, error)
}

type Backend interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

func (s *Service) SendBid(
	bid *bidderapiv1.Bid,
	srv bidderapiv1.Bidder_SendBidServer,
) error {
	// timeout to prevent hanging of bidder node if provider node is not responding
	ctx, cancel := context.WithTimeout(srv.Context(), s.bidTimeout)
	defer cancel()

	// Simulate preconf from provider for demo. Just open bid directly in bidder registry
	commitmentDigest := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	bidAmt, success := big.NewInt(0).SetString(bid.Amount, 10)
	if !success {
		s.logger.Error("parsing bid amount", "bid amount", bid.Amount)
		return status.Errorf(codes.InvalidArgument, "parsing bid amount: %v", bid.Amount)
	}

	bidderAddr := s.owner
	providerAddr := common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	tx, err := s.registryContract.OpenBid(opts, commitmentDigest, bidAmt, bidderAddr, providerAddr)
	if err != nil {
		s.logger.Error("opening bid", "error", err)
		return status.Errorf(codes.Internal, "opening bid: %v", err)
	}

	receipt, err := s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		s.logger.Error("waiting for receipt", "error", err)
		return status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		s.logger.Error("receipt status", "status", receipt.Status)
		return status.Errorf(codes.Internal, "receipt status: %v", receipt.Status)
	}

	return nil
}

func (s *Service) Deposit(
	ctx context.Context,
	r *bidderapiv1.DepositRequest,
) (*bidderapiv1.DepositResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("deposit validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating deposit request: %v", err)
	}

	amount, success := big.NewInt(0).SetString(r.Amount, 10)
	if !success {
		s.logger.Error("parsing amount", "amount", r.Amount)
		return nil, status.Errorf(codes.InvalidArgument, "parsing amount: %v", r.Amount)
	}

	if amount.Cmp(big.NewInt(0)) <= 0 {
		s.logger.Error("amount must be positive", "amount", r.Amount)
		return nil, status.Errorf(codes.InvalidArgument, "amount must be positive: %v", r.Amount)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}
	opts.Value = amount

	providerAddr := common.HexToAddress(r.Provider)
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
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for deposit")
}

func (s *Service) DepositEvenly(
	ctx context.Context,
	r *bidderapiv1.DepositEvenlyRequest,
) (*bidderapiv1.DepositEvenlyResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("deposit evenly validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating deposit evenly request: %v", err)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	totalAmount, success := big.NewInt(0).SetString(r.TotalAmount, 10)
	if !success {
		s.logger.Error("parsing total amount", "total amount", r.TotalAmount)
		return nil, status.Errorf(codes.InvalidArgument, "parsing total amount: %v", r.TotalAmount)
	}

	if totalAmount.Cmp(big.NewInt(0)) <= 0 {
		s.logger.Error("total amount must be positive", "total amount", r.TotalAmount)
		return nil, status.Errorf(codes.InvalidArgument, "total amount must be positive: %v", r.TotalAmount)
	}

	providers := make([]common.Address, len(r.Providers))
	for i, provider := range r.Providers {
		providers[i] = common.HexToAddress(provider)
	}

	opts.Value = totalAmount

	tx, err := s.registryContract.DepositEvenlyAsBidder(opts, providers)
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

	expectedLogs := len(r.Providers)
	receivedLogs := 0

	response := &bidderapiv1.DepositEvenlyResponse{}
	for _, log := range receipt.Logs {
		if deposited, err := s.registryContract.ParseBidderDeposited(*log); err == nil {
			receivedLogs++
			response.Providers = append(response.Providers, common.Bytes2Hex(deposited.Provider.Bytes()))
			response.Amounts = append(response.Amounts, deposited.DepositedAmount.String())
		}
		if receivedLogs == expectedLogs {
			return response, nil
		}
	}

	s.logger.Error(
		"deposit evenly successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for deposit evenly")
}

func (s *Service) GetDeposit(
	ctx context.Context,
	r *bidderapiv1.GetDepositRequest,
) (*bidderapiv1.DepositResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("get deposit validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating get deposit request: %v", err)
	}

	providerAddr := common.HexToAddress(r.Provider)
	deposit, err := s.registryContract.GetDeposit(&bind.CallOpts{
		From:    s.owner,
		Context: ctx,
	}, s.owner, providerAddr)
	if err != nil {
		s.logger.Error("getting deposit", "error", err)
		return nil, status.Errorf(codes.Internal, "getting deposit: %v", err)
	}

	return &bidderapiv1.DepositResponse{Amount: deposit.String(), Provider: r.Provider}, nil
}

func (s *Service) GetAllDeposits(
	ctx context.Context,
	r *bidderapiv1.GetAllDepositsRequest,
) (*bidderapiv1.GetAllDepositsResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("get all deposits validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating get all deposits request: %v", err)
	}

	deposits, err := s.registryContract.FilterBidderDeposited(
		&bind.FilterOpts{
			Context: ctx,
			Start:   0,
			End:     nil,
		},
		[]common.Address{s.owner}, // This bidder
		[]common.Address{},        // all providers
		[]*big.Int{},              // all amounts
	)
	if err != nil {
		s.logger.Error("filtering bidder deposited", "error", err)
		return nil, status.Errorf(codes.Internal, "filtering bidder deposited: %v", err)
	}

	providersToQuery := make(map[common.Address]bool)
	for deposits.Next() {
		providersToQuery[deposits.Event.Provider] = true
	}
	if err := deposits.Error(); err != nil {
		s.logger.Error("error iterating over deposits", "error", err)
		return nil, status.Errorf(codes.Internal, "error iterating over deposits: %v", err)
	}

	response := &bidderapiv1.GetAllDepositsResponse{}
	for provider := range providersToQuery {
		deposit, err := s.registryContract.GetDeposit(&bind.CallOpts{
			From:    s.owner,
			Context: ctx,
		}, s.owner, provider)
		if err != nil {
			s.logger.Error("getting deposit", "error", err)
			return nil, status.Errorf(codes.Internal, "getting deposit: %v", err)
		}
		if deposit.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		response.Deposits = append(response.Deposits, &bidderapiv1.DepositInfo{
			Provider: provider.Hex(),
			Amount:   deposit.String(),
		})
	}

	balance, err := s.backend.BalanceAt(ctx, s.owner, nil)
	if err != nil {
		s.logger.Error("getting bidder balance", "error", err)
		return nil, status.Errorf(codes.Internal, "getting bidder balance: %v", err)
	}
	response.BidderBalance = balance.String()

	return response, nil
}

func (s *Service) RequestWithdrawals(
	ctx context.Context,
	r *bidderapiv1.RequestWithdrawalsRequest,
) (*bidderapiv1.RequestWithdrawalsResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("request withdrawals validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating request withdrawals request: %v", err)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	providers := make([]common.Address, len(r.Providers))
	for i, provider := range r.Providers {
		providers[i] = common.HexToAddress(provider)
	}

	tx, err := s.registryContract.RequestWithdrawalsAsBidder(opts, providers)
	if err != nil {
		s.logger.Error("requesting withdrawals", "error", err)
		return nil, status.Errorf(codes.Internal, "requesting withdrawals: %v", err)
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

	expectedLogs := len(r.Providers)
	receivedLogs := 0

	response := &bidderapiv1.RequestWithdrawalsResponse{}
	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseWithdrawalRequested(*log); err == nil {
			receivedLogs++
			response.Providers = append(response.Providers, common.Bytes2Hex(withdrawal.Provider.Bytes()))
			response.Amounts = append(response.Amounts, withdrawal.AvailableAmount.String())
		}
		if receivedLogs == expectedLogs {
			return response, nil
		}
	}

	s.logger.Error(
		"request withdrawals successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for request withdrawals")
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

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	providers := make([]common.Address, len(r.Providers))
	for i, provider := range r.Providers {
		providers[i] = common.HexToAddress(provider)
	}

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

	expectedLogs := len(r.Providers)
	receivedLogs := 0

	response := &bidderapiv1.WithdrawResponse{}
	for _, log := range receipt.Logs {
		if withdrawal, err := s.registryContract.ParseBidderWithdrawal(*log); err == nil {
			receivedLogs++
			response.Amounts = append(response.Amounts, withdrawal.AmountWithdrawn.String())
			response.Providers = append(response.Providers, common.Bytes2Hex(withdrawal.Provider.Bytes()))
		}
		if receivedLogs == expectedLogs {
			return response, nil
		}
	}

	s.logger.Error(
		"withdraw successful but missing log",
		"txHash", receipt.TxHash.Hex(),
		"logs", receipt.Logs,
	)

	return nil, status.Errorf(codes.Internal, "missing log for withdraw")
}

func (s *Service) EnableDepositManager(
	ctx context.Context,
	r *bidderapiv1.EnableDepositManagerRequest,
) (*bidderapiv1.EnableDepositManagerResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("enable deposit manager validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating enable deposit manager request: %v", err)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	depositManagerEnabled, err := s.DepositManagerStatus(ctx, &bidderapiv1.DepositManagerStatusRequest{})
	if err != nil {
		s.logger.Error("checking deposit manager status", "error", err)
		return nil, status.Errorf(codes.Internal, "checking deposit manager status: %v", err)
	}

	if depositManagerEnabled.Enabled {
		s.logger.Error("EnableDepositManager failed: deposit manager is already enabled")
		return nil, status.Errorf(codes.FailedPrecondition, "EnableDepositManager failed: deposit manager is already enabled")
	}

	tx, err := s.setCodeHelper.SetCode(ctx, opts, s.depositManagerImplAddr)
	if err != nil {
		s.logger.Error("setting code", "error", err)
		return nil, status.Errorf(codes.Internal, "setting code: %v", err)
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

	return &bidderapiv1.EnableDepositManagerResponse{Success: true}, nil
}

func (s *Service) SetTargetDeposits(
	ctx context.Context,
	r *bidderapiv1.SetTargetDepositsRequest,
) (*bidderapiv1.SetTargetDepositsResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("set target deposits validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating set target deposits request: %v", err)
	}

	opts, err := s.optsGetter(ctx)
	if err != nil {
		s.logger.Error("getting transact opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting transact opts: %v", err)
	}

	depositManagerEnabled, err := s.DepositManagerStatus(ctx, &bidderapiv1.DepositManagerStatusRequest{})
	if err != nil {
		s.logger.Error("checking deposit manager status", "error", err)
		return nil, status.Errorf(codes.Internal, "checking deposit manager status: %v", err)
	}

	if !depositManagerEnabled.Enabled {
		s.logger.Error("SetTargetDeposits failed: deposit manager is not enabled")
		return nil, status.Errorf(codes.FailedPrecondition, "SetTargetDeposits failed: deposit manager is not enabled")
	}

	if len(r.TargetDeposits) == 0 {
		s.logger.Error("SetTargetDeposits failed: no target deposits provided")
		return nil, status.Errorf(codes.InvalidArgument, "SetTargetDeposits failed: no target deposits provided")
	}

	providers := make([]common.Address, len(r.TargetDeposits))
	amounts := make([]*big.Int, len(r.TargetDeposits))
	for i, targetDeposit := range r.TargetDeposits {
		providers[i] = common.HexToAddress(targetDeposit.Provider)
		amounts[i] = big.NewInt(0)
		amounts[i].SetString(targetDeposit.TargetDeposit, 10)
	}
	tx, err := s.depositManager.SetTargetDeposits(opts, providers, amounts)
	if err != nil {
		s.logger.Error("setting target deposits", "error", err)
		return nil, status.Errorf(codes.Internal, "setting target deposits: %v", err)
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

	response := &bidderapiv1.SetTargetDepositsResponse{}
	for _, log := range receipt.Logs {
		if targetDeposit, err := s.depositManager.ParseTargetDepositSet(*log); err == nil {
			response.SuccessfullySetDeposits = append(response.SuccessfullySetDeposits, &bidderapiv1.TargetDeposit{
				Provider:      common.Bytes2Hex(targetDeposit.Provider.Bytes()),
				TargetDeposit: targetDeposit.Amount.String(),
			})
		}
	}

	tx, err = s.depositManager.TopUpDeposits(opts, providers)
	if err != nil {
		s.logger.Error("topping up deposits", "error", err)
		return nil, status.Errorf(codes.Internal, "topping up deposits: %v", err)
	}

	receipt, err = s.watcher.WaitForReceipt(ctx, tx)
	if err != nil {
		s.logger.Error("waiting for receipt", "error", err)
		return nil, status.Errorf(codes.Internal, "waiting for receipt: %v", err)
	}

	for _, log := range receipt.Logs {
		if depositToppedUp, err := s.depositManager.ParseDepositToppedUp(*log); err == nil {
			response.SuccessfullyToppedUpProviders = append(
				response.SuccessfullyToppedUpProviders,
				depositToppedUp.Provider.Hex(),
			)
		}
	}

	return response, nil
}

func (s *Service) DepositManagerStatus(
	ctx context.Context,
	r *bidderapiv1.DepositManagerStatusRequest,
) (*bidderapiv1.DepositManagerStatusResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("deposit manager status validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating deposit manager status request: %v", err)
	}

	code, err := s.backend.CodeAt(ctx, s.owner, nil)
	if err != nil {
		s.logger.Error("getting code", "error", err)
		return nil, status.Errorf(codes.Internal, "getting code: %v", err)
	}
	if len(code) == 0 {
		s.logger.Info("deposit manager not enabled")
		return &bidderapiv1.DepositManagerStatusResponse{Enabled: false}, nil
	}

	codehash := crypto.Keccak256Hash(code)
	expectedCodehash := crypto.Keccak256Hash(common.FromHex("0xef0100"), s.depositManagerImplAddr.Bytes())
	if codehash != expectedCodehash {
		s.logger.Error("codehash is not correct", "actual", codehash, "expected", expectedCodehash)
		return nil, status.Errorf(codes.Internal, "codehash is not correct")
	}

	return &bidderapiv1.DepositManagerStatusResponse{Enabled: true}, nil
}

// TODO: api/handling for a bidder removing set code auth

func (s *Service) GetValidProviders(
	ctx context.Context,
	r *bidderapiv1.GetValidProvidersRequest,
) (*bidderapiv1.GetValidProvidersResponse, error) {
	err := s.validator.Validate(r)
	if err != nil {
		s.logger.Error("get valid providers validation", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validating get valid providers request: %v", err)
	}

	// Simulate provider p2p connection
	connectedProviders := []p2p.Peer{
		p2p.Peer{
			EthAddress: common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"),
			Type:       p2p.PeerTypeProvider,
		},
	}

	providerAddrs := make([]common.Address, len(connectedProviders))
	for i, provider := range connectedProviders {
		providerAddrs[i] = provider.EthAddress
	}

	validProviders := make([]string, 0)
	areValid, err := s.providerRegistry.AreProvidersValid(&bind.CallOpts{
		Context: ctx,
	}, providerAddrs)
	if err != nil {
		s.logger.Error("checking if providers are valid", "error", err)
		return nil, status.Errorf(codes.Internal, "checking if providers are valid: %v", err)
	}

	for i, isValid := range areValid {
		if isValid {
			validProviders = append(validProviders, connectedProviders[i].EthAddress.Hex())
		}
	}

	return &bidderapiv1.GetValidProvidersResponse{ValidProviders: validProviders}, nil
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
