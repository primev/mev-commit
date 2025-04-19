package deposit

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

const (
	noOfWindows = 2
	withdraw    = true
)

func RunAutoDeposit(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	bidders := cluster.Bidders()
	logger := cluster.Logger().With("test", "autodeposit")

	deposits := make(chan *bidderregistry.BidderregistryBidderRegistered)
	withdrawals := make(chan *bidderregistry.BidderregistryBidderWithdrawal)
	window := make(chan *blocktracker.BlocktrackerNewWindow)

	// Listen for deposits and withdrawals
	sub, err := cluster.Events().Subscribe(
		events.NewEventHandler(
			"BidderRegistered",
			func(r *bidderregistry.BidderregistryBidderRegistered) {
				deposits <- r
			},
		),
		events.NewEventHandler(
			"BidderWithdrawal",
			func(r *bidderregistry.BidderregistryBidderWithdrawal) {
				withdrawals <- r
			},
		),
		events.NewEventHandler(
			"NewWindow",
			func(w *blocktracker.BlocktrackerNewWindow) {
				window <- w
			},
		),
	)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	var start atomic.Value
	depositsRcvd := make(map[common.Address][]*bidderregistry.BidderregistryBidderRegistered)
	withdrawalsRcvd := make(map[common.Address][]*bidderregistry.BidderregistryBidderWithdrawal)

	eg := errgroup.Group{}
	egCtx, egCancel := context.WithCancel(ctx)
	defer egCancel()

	eg.Go(func() error {
		logger.Info("Starting test... waiting for new window")
		for {
			select {
			case <-egCtx.Done():
				return nil
			case r := <-deposits:
				logger.Info("Received deposit", "bidder", r.Bidder)
				depositsRcvd[r.Bidder] = append(depositsRcvd[r.Bidder], r)
			case r := <-withdrawals:
				logger.Info("Received withdrawal", "bidder", r.Bidder)
				withdrawalsRcvd[r.Bidder] = append(withdrawalsRcvd[r.Bidder], r)
			case w := <-window:
				logger.Info("Received new window", "window", w.Window)
				switch {
				case start.Load() == nil:
					for _, bidder := range bidders {
						resp, err := bidder.BidderAPI().AutoDeposit(egCtx, &bidderapiv1.DepositRequest{
							Amount: "1000000000000000000",
						})
						if err != nil {
							return err
						}
						logger.Info(
							"Auto deposit",
							"bidder", bidder.EthAddress(),
							"window", w.Window,
							"response", resp,
						)
					}
					start.Store(new(big.Int).Add(w.Window, big.NewInt(1)))
					logger.Info("Autodeposit", "start", start.Load())
				case new(big.Int).Sub(w.Window, start.Load().(*big.Int)).Cmp(big.NewInt(noOfWindows)) == 0:
					logger.Info("Finish autodeposit checker", "window", w.Window)
					return nil
				}
			}
		}
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	// ensure we have deposits from start window and withdrawals from start window
	for _, bidder := range bidders {
		depositsForBidder := depositsRcvd[common.HexToAddress(bidder.EthAddress())]
		withdrawalsForBidder := withdrawalsRcvd[common.HexToAddress(bidder.EthAddress())]

		for i, deposit := range depositsForBidder {
			expWindow := new(big.Int).Add(start.Load().(*big.Int), big.NewInt(int64(i)))
			if deposit.WindowNumber.Cmp(expWindow) != 0 {
				logger.Error(
					"Deposit received in wrong window",
					"bidder", bidder.EthAddress(),
					"expected", expWindow,
					"received", deposit.WindowNumber,
				)
				return fmt.Errorf("deposit received in wrong window")
			}
		}

		for i, withdrawal := range withdrawalsForBidder {
			expWindow := new(big.Int).Add(start.Load().(*big.Int), big.NewInt(int64(i)))
			if withdrawal.Window.Cmp(expWindow) != 0 {
				logger.Error(
					"Withdrawal received in wrong window",
					"bidder", bidder.EthAddress(),
					"expected", expWindow,
					"received", withdrawal.Window,
				)
				return fmt.Errorf("withdrawal received in wrong window")
			}
		}
	}

	return nil
}

func RunCancelAutoDeposit(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	bidders := cluster.Bidders()
	logger := cluster.Logger().With("test", "cancel_autodeposit")

	deposits := make(chan *bidderregistry.BidderregistryBidderRegistered)
	withdrawals := make(chan *bidderregistry.BidderregistryBidderWithdrawal)
	window := make(chan *blocktracker.BlocktrackerNewWindow)

	// Listen for deposits and withdrawals
	sub, err := cluster.Events().Subscribe(
		events.NewEventHandler(
			"BidderRegistered",
			func(r *bidderregistry.BidderregistryBidderRegistered) {
				deposits <- r
			},
		),
		events.NewEventHandler(
			"BidderWithdrawal",
			func(r *bidderregistry.BidderregistryBidderWithdrawal) {
				withdrawals <- r
			},
		),
		events.NewEventHandler(
			"NewWindow",
			func(w *blocktracker.BlocktrackerNewWindow) {
				window <- w
			},
		),
	)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	var stop, end atomic.Value
	depositsRcvd := make(map[common.Address][]*bidderregistry.BidderregistryBidderRegistered)
	withdrawalsRcvd := make(map[common.Address][]*bidderregistry.BidderregistryBidderWithdrawal)

	eg, ctx := errgroup.WithContext(ctx)
	egCtx, egCancel := context.WithCancel(ctx)
	defer egCancel()

	eg.Go(func() error {
		logger.Info("Starting test... waiting for new window")
		for {
			select {
			case <-egCtx.Done():
				return nil
			case r := <-deposits:
				logger.Info("Received deposit", "bidder", r.Bidder)
				depositsRcvd[r.Bidder] = append(depositsRcvd[r.Bidder], r)
				if stop.Load() == nil {
					for _, bidder := range bidders {
						resp, err := bidder.BidderAPI().CancelAutoDeposit(egCtx, &bidderapiv1.CancelAutoDepositRequest{
							Withdraw: withdraw,
						})
						if err != nil {
							return err
						}
						logger.Info(
							"Cancelled auto deposit",
							"bidder", bidder.EthAddress(),
							"window", r.WindowNumber,
							"response", resp,
						)
					}
					stop.Store(new(big.Int).Add(r.WindowNumber, big.NewInt(1)))
					end.Store(new(big.Int).Add(r.WindowNumber, big.NewInt(4)))
				}
			case r := <-withdrawals:
				logger.Info("Received withdrawal", "bidder", r.Bidder)
				withdrawalsRcvd[r.Bidder] = append(withdrawalsRcvd[r.Bidder], r)
			case w := <-window:
				logger.Info("Received new window", "window", w.Window)
				if end.Load() != nil && w.Window.Cmp(end.Load().(*big.Int)) == 0 {
					logger.Info("Finished test", "window", w.Window)
					return nil
				}
			}
		}
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	for _, bidder := range bidders {
		depositsForBidder := depositsRcvd[common.HexToAddress(bidder.EthAddress())]
		withdrawalsForBidder := withdrawalsRcvd[common.HexToAddress(bidder.EthAddress())]

		if depositsForBidder[len(depositsForBidder)-1].WindowNumber.Cmp(stop.Load().(*big.Int)) == 1 {
			logger.Error(
				"Last deposit received after stop window",
				"bidder", bidder.EthAddress(),
				"stop", stop.Load(),
				"received", depositsForBidder[len(depositsForBidder)-1].WindowNumber,
			)
			return fmt.Errorf("deposit received after stop window")
		}

		// last withdrawal should be 1 less than the stop window
		if withdrawalsForBidder[len(withdrawalsForBidder)-1].Window.Cmp(stop.Load().(*big.Int)) != -1 {
			logger.Error(
				"Last withdrawal received after stop window",
				"bidder", bidder.EthAddress(),
				"stop", stop.Load(),
				"received", withdrawalsForBidder[len(withdrawalsForBidder)-1].Window,
			)
			return fmt.Errorf("withdrawal received after stop window")
		}
	}

	return nil
}
