package apiserver

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconfcommitmentstore "github.com/primev/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/x/contracts/events"
)

type BlockStats struct {
	Number                 uint64 `json:"number"`
	Winner                 string `json:"winner"`
	Window                 int64  `json:"window"`
	TotalOpenedCommitments int    `json:"total_opened_commitments"`
	TotalRewards           int    `json:"total_rewards"`
	TotalSlashes           int    `json:"total_slashes"`
	TotalAmount            string `json:"total_amount"`
}

type ProviderBalances struct {
	Provider string `json:"provider"`
	Stake    string `json:"stake"`
	Rewards  string `json:"rewards"`
}

type BidderAllowance struct {
	Bidder    string `json:"bidder"`
	Allowance string `json:"allowance"`
	Refunds   string `json:"refunds"`
	Settled   string `json:"settled"`
	Withdrawn string `json:"withdrawn"`
}

type DashboardOut struct {
	Block     *BlockStats         `json:"block"`
	Providers []*ProviderBalances `json:"providers"`
	Bidders   []*BidderAllowance  `json:"bidders"`
}

func (s *Service) configureDashboard() error {
	handlers := []events.EventHandler{
		events.NewEventHandler(
			"NewL1Block",
			func(upd *blocktracker.BlocktrackerNewL1Block) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.blockStats.Get(upd.BlockNumber.Uint64())
				if !ok {
					existing = &BlockStats{
						Number: upd.BlockNumber.Uint64(),
					}
				}

				existing.Winner = upd.Winner.Hex()
				existing.Window = upd.Window.Int64()
				_ = s.blockStats.Add(upd.BlockNumber.Uint64(), existing)
				if upd.BlockNumber.Uint64() > s.lastBlock {
					s.lastBlock = upd.BlockNumber.Uint64()
				}
			},
		),
		events.NewEventHandler(
			"CommitmentStored",
			func(upd *preconfcommitmentstore.PreconfcommitmentstoreCommitmentStored) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.blockStats.Get(upd.BlockNumber)
				if !ok {
					existing = &BlockStats{
						Number: upd.BlockNumber,
					}
				}

				existing.TotalOpenedCommitments++
				_ = s.blockStats.Add(upd.BlockNumber, existing)
			},
		),
		events.NewEventHandler(
			"CommitmentProcessed",
			func(upd *oracle.OracleCommitmentProcessed) {
				cmt, err := s.store.Settlement(context.Background(), upd.CommitmentIndex[:])
				if err != nil {
					s.logger.Error("failed to get settlement", "error", err)
					return
				}

				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.blockStats.Get(uint64(cmt.BlockNum))
				if !ok {
					existing = &BlockStats{
						Number: uint64(cmt.BlockNum),
					}
				}

				if upd.IsSlash {
					existing.TotalSlashes++
				} else {
					existing.TotalRewards++
				}
				currentAmount, ok := big.NewInt(0).SetString(existing.TotalAmount, 10)
				if !ok {
					currentAmount = big.NewInt(0)
				}
				currentAmount = big.NewInt(0).Add(currentAmount, cmt.Amount)
				existing.TotalAmount = currentAmount.String()
				_ = s.blockStats.Add(uint64(cmt.BlockNum), existing)
			},
		),
		events.NewEventHandler(
			"ProviderRegistered",
			func(upd *providerregistry.ProviderregistryProviderRegistered) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.providerStakes.Get(upd.Provider.Hex())
				if !ok {
					existing = &ProviderBalances{
						Provider: upd.Provider.Hex(),
					}
				}
				existing.Stake = upd.StakedAmount.String()
				_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			},
		),
		events.NewEventHandler(
			"FundsDeposited",
			func(upd *providerregistry.ProviderregistryFundsDeposited) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.providerStakes.Get(upd.Provider.Hex())
				if !ok {
					s.logger.Error("provider not found", "provider", upd.Provider.Hex())
					return
				}
				currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
				if !ok {
					s.logger.Error("failed to parse stake", "stake", existing.Stake)
					return
				}
				currentStake = big.NewInt(0).Add(currentStake, upd.Amount)
				existing.Stake = currentStake.String()
				_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			},
		),
		events.NewEventHandler(
			"FundsSlashed",
			func(upd *providerregistry.ProviderregistryFundsSlashed) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.providerStakes.Get(upd.Provider.Hex())
				if !ok {
					s.logger.Error("provider not found", "provider", upd.Provider.Hex())
					return
				}
				currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
				if !ok {
					s.logger.Error("failed to parse stake", "stake", existing.Stake)
					return
				}
				currentStake = big.NewInt(0).Sub(currentStake, upd.Amount)
				existing.Stake = currentStake.String()
				_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			},
		),
		events.NewEventHandler(
			"FundsRewarded",
			func(upd *bidderregistry.BidderregistryFundsRewarded) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.providerStakes.Get(upd.Provider.Hex())
				if !ok {
					s.logger.Error("provider not found", "provider", upd.Provider.Hex())
					return
				}
				currentRewards, ok := big.NewInt(0).SetString(existing.Rewards, 10)
				if !ok {
					currentRewards = big.NewInt(0)
				}
				currentRewards = big.NewInt(0).Add(currentRewards, upd.Amount)
				existing.Rewards = currentRewards.String()
				_ = s.providerStakes.Add(upd.Provider.Hex(), existing)

				existingBidders, ok := s.bidderAllowances.Get(upd.Window.Uint64())
				if !ok {
					s.logger.Error("window not found", "window", upd.Window)
					return
				}
				for _, b := range existingBidders {
					if b.Bidder == upd.Bidder.Hex() {
						currentSettled, ok := big.NewInt(0).SetString(b.Settled, 10)
						if !ok {
							currentSettled = big.NewInt(0)
						}
						currentSettled = big.NewInt(0).Add(currentSettled, upd.Amount)
						b.Settled = currentSettled.String()
						break
					}
				}
				_ = s.bidderAllowances.Add(upd.Window.Uint64(), existingBidders)
			},
		),
		events.NewEventHandler(
			"BidderRegistered",
			func(upd *bidderregistry.BidderregistryBidderRegistered) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderAllowances.Get(upd.WindowNumber.Uint64())
				if !ok {
					existing = make([]*BidderAllowance, 0)
				}

				for _, b := range existing {
					if b.Bidder == upd.Bidder.Hex() {
						s.logger.Error("bidder already registered", "bidder", upd.Bidder.Hex())
						return
					}
				}

				existing = append(existing, &BidderAllowance{
					Bidder:    upd.Bidder.Hex(),
					Allowance: upd.DepositedAmount.String(),
				})
				_ = s.bidderAllowances.Add(upd.WindowNumber.Uint64(), existing)
			},
		),
		events.NewEventHandler(
			"FundsRetrieved",
			func(upd *bidderregistry.BidderregistryFundsRetrieved) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderAllowances.Get(upd.Window.Uint64())
				if !ok {
					s.logger.Error("window not found", "window", upd.Window)
					return
				}

				for _, b := range existing {
					if b.Bidder == upd.Bidder.Hex() {
						currentReturned, ok := big.NewInt(0).SetString(b.Refunds, 10)
						if !ok {
							currentReturned = big.NewInt(0)
						}
						currentReturned = big.NewInt(0).Add(currentReturned, upd.Amount)
						b.Refunds = currentReturned.String()
						break
					}
				}
				_ = s.bidderAllowances.Add(upd.Window.Uint64(), existing)
			},
		),
		events.NewEventHandler(
			"BidderWithdrawal",
			func(upd *bidderregistry.BidderregistryBidderWithdrawal) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderAllowances.Get(upd.Window.Uint64())
				if !ok {
					s.logger.Error("window not found", "window", upd.Window)
					return
				}

				for idx, b := range existing {
					if b.Bidder == upd.Bidder.Hex() {
						existing[idx].Withdrawn = upd.Amount.String()
						break
					}
				}

				_ = s.bidderAllowances.Add(upd.Window.Uint64(), existing)
			},
		),
	}

	subs := make([]events.Subscription, 0, len(handlers))
	unsub := func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}

	for _, h := range handlers {
		sub, err := s.evtMgr.Subscribe(h)
		if err != nil {
			unsub()
			return err
		}
		subs = append(subs, sub)
	}

	closed := make(chan struct{})
	go func() {
		<-s.shutdown
		close(closed)
		unsub()
	}()

	s.router.Handle("/dashboard", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-closed:
				http.Error(w, "listener closed", http.StatusServiceUnavailable)
			default:
			}

			limit := 10
			limitStr := r.URL.Query().Get("limit")
			if limitStr != "" {
				l, err := strconv.Atoi(limitStr)
				if err == nil {
					limit = l
				}
			}

			page := 0
			pageStr := r.URL.Query().Get("page")
			if pageStr != "" {
				p, err := strconv.Atoi(pageStr)
				if err == nil {
					page = p
				}
			}

			s.statMu.RLock()
			defer s.statMu.RUnlock()

			lastBlock := s.lastBlock
			lastBlockStr := r.URL.Query().Get("last_block")
			if lastBlockStr != "" {
				lb, err := strconv.ParseUint(lastBlockStr, 10, 64)
				if err == nil {
					lastBlock = lb
				}
			}

			start := lastBlock
			if start > uint64(limit*page) {
				start = start - uint64(limit*page)
			}

			dash := make([]*DashboardOut, 0)

			for i := start; i > 0 && len(dash) <= limit; i-- {
				stats, ok := s.blockStats.Get(i)
				if !ok {
					continue
				}
				bidders, ok := s.bidderAllowances.Get(uint64(stats.Window))
				if !ok {
					bidders = make([]*BidderAllowance, 0)
				}

				providers := s.providerStakes.Values()

				dashEntry := &DashboardOut{
					Block:     stats,
					Providers: providers,
					Bidders:   bidders,
				}
				dash = append(dash, dashEntry)
			}

			if err := json.NewEncoder(w).Encode(dash); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		}),
	)

	return nil
}
