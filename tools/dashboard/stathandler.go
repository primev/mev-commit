package main

import (
	"fmt"
	"math/big"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/x/contracts/events"
)

type statHandler struct {
	statMu                    sync.RWMutex
	lastBlock                 uint64
	lastWindow                uint64
	blocksPerWindow           uint64
	blockStats                *lru.Cache[uint64, *BlockStats]
	providerStakes            *lru.Cache[string, *ProviderBalances]
	bidderAllowances          *lru.Cache[uint64, []*BidderAllowance]
	commitments               *lru.Cache[[32]byte, *preconf.PreconfmanagerOpenedCommitmentStored]
	commitmentsByBlock        *lru.Cache[uint64, []*preconf.PreconfmanagerOpenedCommitmentStored]
	totalEncryptedCommitments uint64
	totalOpenedCommitments    uint64
	totalRewards              uint64
	totalSlashes              uint64
	evtMgr                    events.EventManager
	sub                       events.Subscription
	unsub                     func()
}

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
	Provider                  string `json:"provider"`
	Stake                     string `json:"stake"`
	Rewards                   string `json:"rewards"`
	EncryptedCommitmentsCount uint64 `json:"encrypted_commitments_count"`
	OpenedCommitmentsCount    uint64 `json:"opened_commitments_count"`
	RewardsCount              uint64 `json:"rewards_count"`
	SlashesCount              uint64 `json:"slashes_count"`
}

type BidderAllowance struct {
	Bidder               string `json:"bidder"`
	Allowance            string `json:"allowance"`
	Refunds              string `json:"refunds"`
	Settled              string `json:"settled"`
	Withdrawn            string `json:"withdrawn"`
	OpenCommitmentsCount uint64 `json:"open_commitments_count"`
	ReturnsCount         uint64 `json:"returns_count"`
	SettledCount         uint64 `json:"settled_count"`
}

type WindowStats struct {
	Window  uint64             `json:"window"`
	Bidders []*BidderAllowance `json:"bidders"`
	Blocks  []*BlockStats      `json:"blocks"`
}

type AggregateStats struct {
	TotalEncryptedCommitments uint64 `json:"total_encrypted_commitments"`
	TotalOpenedCommitments    uint64 `json:"total_opened_commitments"`
	TotalRewards              uint64 `json:"total_rewards"`
	TotalSlashes              uint64 `json:"total_slashes"`
}

type DashboardOut struct {
	Aggregate *AggregateStats     `json:"aggregate"`
	Providers []*ProviderBalances `json:"providers"`
	Windows   []*WindowStats      `json:"windows"`
}

func newStatHandler(evtMgr events.EventManager, blocksPerWindow uint64) (*statHandler, error) {
	blockStats, err := lru.New[uint64, *BlockStats](10000)
	if err != nil {
		return nil, err
	}

	providerStakes, err := lru.New[string, *ProviderBalances](100)
	if err != nil {
		return nil, err
	}

	bidderAllowances, err := lru.New[uint64, []*BidderAllowance](1000)
	if err != nil {
		return nil, err
	}

	commitments, err := lru.New[[32]byte, *preconf.PreconfmanagerOpenedCommitmentStored](10000)
	if err != nil {
		return nil, err
	}

	commitmentsByBlock, err := lru.New[uint64, []*preconf.PreconfmanagerOpenedCommitmentStored](10000)
	if err != nil {
		return nil, err
	}

	st := &statHandler{
		blocksPerWindow:    blocksPerWindow,
		blockStats:         blockStats,
		providerStakes:     providerStakes,
		bidderAllowances:   bidderAllowances,
		commitments:        commitments,
		commitmentsByBlock: commitmentsByBlock,
		evtMgr:             evtMgr,
	}

	if err := st.configureDashboard(); err != nil {
		return nil, err
	}

	return st, nil
}

func (s *statHandler) configureDashboard() error {
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
				if upd.Window.Uint64() > s.lastWindow {
					s.lastWindow = upd.Window.Uint64()
				}
			},
		),
		events.NewEventHandler(
			"UnopenedCommitmentStored",
			func(upd *preconf.PreconfmanagerUnopenedCommitmentStored) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				s.totalEncryptedCommitments++
				provider, ok := s.providerStakes.Get(upd.Committer.Hex())
				if !ok {
					return
				}
				provider.EncryptedCommitmentsCount++
				_ = s.providerStakes.Add(upd.Committer.Hex(), provider)
			},
		),
		events.NewEventHandler(
			"OpenedCommitmentStored",
			func(upd *preconf.PreconfmanagerOpenedCommitmentStored) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.blockStats.Get(upd.BlockNumber)
				if !ok {
					existing = &BlockStats{
						Number: upd.BlockNumber,
					}
				}

				existing.TotalOpenedCommitments++
				s.totalOpenedCommitments++
				_ = s.blockStats.Add(upd.BlockNumber, existing)
				_ = s.commitments.Add(upd.CommitmentIndex, upd)

				blockCommitments, ok := s.commitmentsByBlock.Get(upd.BlockNumber)
				if !ok {
					blockCommitments = make([]*preconf.PreconfmanagerOpenedCommitmentStored, 0)
				}
				blockCommitments = append(blockCommitments, upd)
				_ = s.commitmentsByBlock.Add(upd.BlockNumber, blockCommitments)

				p, ok := s.providerStakes.Get(upd.Committer.Hex())
				if !ok {
					return
				}
				p.OpenedCommitmentsCount++
				_ = s.providerStakes.Add(upd.Committer.Hex(), p)
				b, ok := s.bidderAllowances.Get(uint64(existing.Window))
				if !ok {
					return
				}
				for _, bidder := range b {
					if bidder.Bidder == upd.Bidder.Hex() {
						bidder.OpenCommitmentsCount++
						break
					}
				}
				_ = s.bidderAllowances.Add(uint64(existing.Window), b)
			},
		),
		events.NewEventHandler(
			"CommitmentProcessed",
			func(upd *oracle.OracleCommitmentProcessed) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				cmt, ok := s.commitments.Get(upd.CommitmentIndex)
				if !ok {
					return
				}

				existing, ok := s.blockStats.Get(cmt.BlockNumber)
				if !ok {
					existing = &BlockStats{
						Number: cmt.BlockNumber,
					}
				}

				if upd.IsSlash {
					existing.TotalSlashes++
					s.totalSlashes++
				} else {
					existing.TotalRewards++
					s.totalRewards++
				}

				currentAmount, ok := big.NewInt(0).SetString(existing.TotalAmount, 10)
				if !ok {
					currentAmount = big.NewInt(0)
				}
				currentAmount = big.NewInt(0).Add(currentAmount, cmt.BidAmt)
				existing.TotalAmount = currentAmount.String()
				_ = s.blockStats.Add(cmt.BlockNumber, existing)
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
				fmt.Println("ProviderRegistered", existing)
			},
		),
		events.NewEventHandler(
			"FundsDeposited",
			func(upd *providerregistry.ProviderregistryFundsDeposited) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.providerStakes.Get(upd.Provider.Hex())
				if !ok {
					return
				}
				currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
				if !ok {
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
					return
				}
				currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
				if !ok {
					return
				}
				currentStake = big.NewInt(0).Sub(currentStake, upd.Amount)
				existing.Stake = currentStake.String()
				existing.SlashesCount++
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
					return
				}
				currentRewards, ok := big.NewInt(0).SetString(existing.Rewards, 10)
				if !ok {
					currentRewards = big.NewInt(0)
				}
				currentRewards = big.NewInt(0).Add(currentRewards, upd.Amount)
				existing.Rewards = currentRewards.String()
				existing.RewardsCount++
				_ = s.providerStakes.Add(upd.Provider.Hex(), existing)

				existingBidders, ok := s.bidderAllowances.Get(upd.Window.Uint64())
				if !ok {
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
						b.SettledCount++
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
						b.ReturnsCount++
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

	sub, err := s.evtMgr.Subscribe(handlers...)
	if err != nil {
		return err
	}

	s.sub = sub
	s.unsub = sub.Unsubscribe

	return nil
}

func (s *statHandler) healthy() bool {
	select {
	case <-s.sub.Err():
		return false
	default:
	}
	return true
}

func (s *statHandler) close() {
	s.unsub()
}

func (s *statHandler) getDashboard(page, limit int) *DashboardOut {
	s.statMu.RLock()
	start := s.lastWindow
	s.statMu.RUnlock()

	if start > uint64(limit*page) {
		start = start - uint64(limit*page)
	}

	windows := make([]*WindowStats, 0)

	for i := start; i >= 1 && len(windows) <= limit; i-- {
		window := s.getWindowStat(i)
		if window == nil {
			continue
		}
		windows = append(windows, window)
	}

	s.statMu.RLock()
	providers := s.providerStakes.Values()
	agg := &AggregateStats{
		TotalEncryptedCommitments: s.totalEncryptedCommitments,
		TotalOpenedCommitments:    s.totalOpenedCommitments,
		TotalRewards:              s.totalRewards,
		TotalSlashes:              s.totalSlashes,
	}
	s.statMu.RUnlock()

	return &DashboardOut{
		Providers: providers,
		Windows:   windows,
		Aggregate: agg,
	}
}

func (s *statHandler) getWindowStat(window uint64) *WindowStats {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	windowStats := new(WindowStats)
	windowStats.Window = window

	blockStart := (window-1)*s.blocksPerWindow + 1
	blockEnd := window * s.blocksPerWindow
	for i := blockEnd; i >= blockStart; i-- {
		stats, ok := s.blockStats.Get(i)
		if !ok {
			continue
		}
		windowStats.Blocks = append(windowStats.Blocks, stats)
	}

	bidders, ok := s.bidderAllowances.Get(window)
	if !ok {
		bidders = make([]*BidderAllowance, 0)
	}
	windowStats.Bidders = bidders

	return windowStats
}

func (s *statHandler) getProviders() []*ProviderBalances {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	return s.providerStakes.Values()
}

func (s *statHandler) getWindows(page, limit int) []*WindowStats {
	s.statMu.RLock()
	start := s.lastWindow
	s.statMu.RUnlock()

	if start > uint64(limit*page) {
		start = start - uint64(limit*page)
	}

	windows := make([]*WindowStats, 0)
	for i := start; i >= 1 && len(windows) <= limit; i-- {
		window := s.getWindowStat(i)
		if window == nil {
			continue
		}
		windows = append(windows, window)
	}

	return windows
}

func (s *statHandler) getCurrentBidders() []*BidderAllowance {
	s.statMu.RLock()
	window := s.lastWindow
	s.statMu.RUnlock()

	return s.getBidders(int(window))
}

func (s *statHandler) getBidders(window int) []*BidderAllowance {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	windowAllowances, ok := s.bidderAllowances.Get(uint64(window))
	if !ok {
		windowAllowances = make([]*BidderAllowance, 0)
	}

	return windowAllowances
}

func (s *statHandler) getBlockStats(block uint64) *BlockStats {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	stats, ok := s.blockStats.Get(block)
	if !ok {
		return nil
	}

	return stats
}

func (s *statHandler) getBlocks(page, limit int) []*BlockStats {
	s.statMu.RLock()
	start := s.lastBlock
	s.statMu.RUnlock()

	if start > uint64(limit*page) {
		start = start - uint64(limit*page)
	}

	blocks := make([]*BlockStats, 0)
	for i := start; i >= 1 && len(blocks) <= limit; i-- {
		stats, ok := s.blockStats.Get(i)
		if !ok {
			continue
		}
		blocks = append(blocks, stats)
	}

	return blocks
}

func (s *statHandler) getOpenCommitmentsByBlock(blockNumber uint64) []*preconf.PreconfmanagerOpenedCommitmentStored {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	commitments, ok := s.commitmentsByBlock.Get(blockNumber)
	if !ok {
		return make([]*preconf.PreconfmanagerOpenedCommitmentStored, 0)
	}

	return commitments
}
