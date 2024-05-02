package settler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-commit/oracle/pkg/keysigner"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var (
	allowedPendingTxnCount = 128
	batchSize              = 10
)

type SettlementType string

const (
	SettlementTypeReward SettlementType = "reward"
	SettlementTypeSlash  SettlementType = "slash"
	SettlementTypeReturn SettlementType = "return"
)

type Settlement struct {
	CommitmentIdx   []byte
	TxHash          string
	BlockNum        int64
	Builder         string
	Amount          uint64
	BidID           []byte
	Type            SettlementType
	DecayPercentage int64
}

type Return struct {
	BidIDs [][32]byte
}

func (r Return) String() string {
	strs := make([]string, len(r.BidIDs))
	for idx, bidID := range r.BidIDs {
		strs[idx] = fmt.Sprintf("%x", bidID)
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

type SettlerRegister interface {
	LastNonce() (int64, error)
	PendingTxnCount() (int, error)
	SubscribeSettlements(ctx context.Context) <-chan Settlement
	SubscribeReturns(ctx context.Context, limit int) <-chan Return
	SettlementInitiated(ctx context.Context, commitmentIdx [][]byte, txHash common.Hash, nonce uint64) error
	MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error)
}

type Oracle interface {
	ProcessBuilderCommitmentForBlockNumber(
		opts *bind.TransactOpts,
		commitmentIdx [32]byte,
		blockNum *big.Int,
		builder string,
		isSlash bool,
		residualDecay *big.Int,
	) (*types.Transaction, error)
	UnlockFunds(opts *bind.TransactOpts, bidIDs [][32]byte) (*types.Transaction, error)
}

type Transactor interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BlockNumber(ctx context.Context) (uint64, error)
}

type Settler struct {
	logger          *slog.Logger
	keySigner       keysigner.KeySigner
	chainID         *big.Int
	owner           common.Address
	rollupClient    Oracle
	settlerRegister SettlerRegister
	client          Transactor
	metrics         *metrics
	txMtx           sync.Mutex
}

func NewSettler(
	logger *slog.Logger,
	keySigner keysigner.KeySigner,
	chainID *big.Int,
	owner common.Address,
	rollupClient Oracle,
	settlerRegister SettlerRegister,
	client Transactor,
) *Settler {
	return &Settler{
		logger:          logger,
		rollupClient:    rollupClient,
		settlerRegister: settlerRegister,
		owner:           owner,
		client:          client,
		keySigner:       keySigner,
		chainID:         chainID,
		metrics:         newMetrics(),
	}
}

func (s *Settler) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := s.keySigner.GetAuth(s.chainID)
	if err != nil {
		return nil, err
	}
	nonce, err := s.client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return nil, err
	}
	usedNonce, err := s.settlerRegister.LastNonce()
	if err != nil {
		return nil, err
	}
	if nonce <= uint64(usedNonce) {
		nonce = uint64(usedNonce + 1)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Returns priority fee per gas
	gasTip, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	// Returns priority fee per gas + base fee per gas
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTip

	return auth, nil
}

func (s *Settler) Metrics() []prometheus.Collector {
	return s.metrics.Collectors()
}

func (s *Settler) settlementUpdater(ctx context.Context) error {
	queryTicker := time.NewTicker(500 * time.Millisecond)
	defer queryTicker.Stop()

	lastBlock := uint64(0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-queryTicker.C:
		}

		currentBlock, err := s.client.BlockNumber(ctx)
		if err != nil {
			s.logger.Error("failed to get block number", "error", err)
			continue
		}

		if currentBlock <= lastBlock {
			continue
		}

		lastNonce, err := s.client.NonceAt(
			ctx,
			s.owner,
			new(big.Int).SetUint64(currentBlock),
		)
		if err != nil {
			s.logger.Error("failed to get nonce", "error", err)
			continue
		}

		count, err := s.settlerRegister.MarkSettlementComplete(ctx, lastNonce)
		if err != nil {
			s.logger.Error("failed to mark settlement complete", "error", err)
			continue
		}

		s.metrics.LastConfirmedNonce.Set(float64(lastNonce))
		s.metrics.LastConfirmedBlock.Set(float64(currentBlock))
		s.metrics.SettlementsConfirmedCount.Add(float64(count))

		if count > 0 {
			s.logger.Info("marked settlement complete", "count", count)
		}

		lastBlock = currentBlock
	}
}

func (s *Settler) settlementExecutor(ctx context.Context) error {
RESTART:
	cctx, unsub := context.WithCancel(ctx)
	settlementChan := s.settlerRegister.SubscribeSettlements(cctx)

	for {
		select {
		case <-ctx.Done():
			unsub()
			return ctx.Err()
		case settlement, more := <-settlementChan:
			if !more {
				unsub()
				goto RESTART
			}

			err := func() error {
				s.txMtx.Lock()
				defer s.txMtx.Unlock()

				if settlement.Type == SettlementTypeReturn {
					s.logger.Warn("return settlement", "commitmentIdx", fmt.Sprintf("%x", settlement.CommitmentIdx))
					return nil
				}

				pendingTxns, err := s.settlerRegister.PendingTxnCount()
				if err != nil {
					return err
				}

				if pendingTxns > allowedPendingTxnCount {
					return errors.New("too many pending txns")
				}

				var (
					commitmentIdx [32]byte
				)

				opts, err := s.getTransactOpts(ctx)
				if err != nil {
					return err
				}

				copy(commitmentIdx[:], settlement.CommitmentIdx)

				commitmentPostingTxn, err := s.rollupClient.ProcessBuilderCommitmentForBlockNumber(
					opts,
					commitmentIdx,
					big.NewInt(settlement.BlockNum),
					settlement.Builder,
					settlement.Type == SettlementTypeSlash,
					big.NewInt(settlement.DecayPercentage),
				)
				if err != nil {
					return fmt.Errorf("process commitment: %w nonce %d", err, opts.Nonce.Uint64())
				}

				err = s.settlerRegister.SettlementInitiated(
					ctx,
					[][]byte{settlement.BidID},
					commitmentPostingTxn.Hash(),
					commitmentPostingTxn.Nonce(),
				)
				if err != nil {
					return fmt.Errorf("failed to mark settlement initiated: %w", err)
				}

				s.metrics.LastUsedNonce.Set(float64(commitmentPostingTxn.Nonce()))
				s.metrics.SettlementsPostedCount.Inc()
				s.metrics.CurrentSettlementL1Block.Set(float64(settlement.BlockNum))

				s.logger.Info(
					"builder commitment processed",
					"blockNum", settlement.BlockNum,
					"txHash", commitmentPostingTxn.Hash().Hex(),
					"builder", settlement.Builder,
					"settlementType", string(settlement.Type),
					"nonce", commitmentPostingTxn.Nonce(),
					"decay", settlement.DecayPercentage,
				)

				return nil
			}()
			if err != nil {
				s.logger.Error("failed to process builder commitment", "error", err)
				unsub()
				time.Sleep(5 * time.Second)
				goto RESTART
			}
		}
	}
}

func (s *Settler) returnExecutor(ctx context.Context) error {
RESTART:
	cctx, unsub := context.WithCancel(ctx)
	returnsChan := s.settlerRegister.SubscribeReturns(cctx, batchSize)

	for {
		select {
		case <-ctx.Done():
			unsub()
			return ctx.Err()
		case returns, more := <-returnsChan:
			if !more {
				unsub()
				goto RESTART
			}

			err := func() error {
				s.txMtx.Lock()
				defer s.txMtx.Unlock()

				pendingTxns, err := s.settlerRegister.PendingTxnCount()
				if err != nil {
					return err
				}

				if pendingTxns > allowedPendingTxnCount {
					return errors.New("too many pending txns")
				}

				opts, err := s.getTransactOpts(ctx)
				if err != nil {
					return err
				}

				bidIDs := make([][]byte, 0, len(returns.BidIDs))
				for _, bidID := range returns.BidIDs {
					b := make([]byte, 32)
					copy(b, bidID[:])
					bidIDs = append(bidIDs, b)
				}

				s.logger.Debug("processing return", "bidIDs", returns, "count", len(returns.BidIDs))

				commitmentPostingTxn, err := s.rollupClient.UnlockFunds(
					opts,
					returns.BidIDs,
				)
				if err != nil {
					return fmt.Errorf("process return: %w nonce %d", err, opts.Nonce.Uint64())
				}

				err = s.settlerRegister.SettlementInitiated(
					ctx,
					bidIDs,
					commitmentPostingTxn.Hash(),
					commitmentPostingTxn.Nonce(),
				)
				if err != nil {
					return fmt.Errorf("failed to mark settlement initiated: %w", err)
				}

				s.metrics.LastUsedNonce.Set(float64(commitmentPostingTxn.Nonce()))
				s.metrics.SettlementsPostedCount.Inc()

				s.logger.Info(
					"builder return processed",
					"txHash", commitmentPostingTxn.Hash().Hex(),
					"batchSize", len(returns.BidIDs),
					"nonce", commitmentPostingTxn.Nonce(),
				)

				return nil
			}()
			if err != nil {
				s.logger.Error("failed to process return", "error", err)
				unsub()
				time.Sleep(5 * time.Second)
				goto RESTART
			}
		}
	}
}

func (s *Settler) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.settlementUpdater(egCtx)
	})

	eg.Go(func() error {
		return s.settlementExecutor(egCtx)
	})

	eg.Go(func() error {
		return s.returnExecutor(egCtx)
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			s.logger.Error("settler error", "error", err)
		}
	}()

	return doneChan
}
