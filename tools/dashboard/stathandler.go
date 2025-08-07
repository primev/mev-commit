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
	blocksPerWindow           uint64
	blockStats                *lru.Cache[uint64, *BlockStats]
	providerStakes            *lru.Cache[string, *ProviderBalances]
	bidderDeposits            *lru.Cache[depositKey, []*BidderDeposit]
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

type BidderDeposit struct {
	Bidder               string `json:"bidder"`
	Provider             string `json:"provider"`
	Amount               string `json:"amount"`
	Refunds              string `json:"refunds"`
	Settled              string `json:"settled"`
	Withdrawn            string `json:"withdrawn"`
	OpenCommitmentsCount uint64 `json:"open_commitments_count"`
	ReturnsCount         uint64 `json:"returns_count"`
	SettledCount         uint64 `json:"settled_count"`
}

type depositKey struct {
	bidder   string
	provider string
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

	bidderDeposits, err := lru.New[depositKey, []*BidderDeposit](1000)
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
		bidderDeposits:     bidderDeposits,
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
				b, ok := s.bidderDeposits.Get(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Committer.Hex(),
				})
				if !ok {
					return
				}
				for _, bidder := range b {
					if bidder.Bidder == upd.Bidder.Hex() {
						bidder.OpenCommitmentsCount++
						break
					}
				}
				_ = s.bidderDeposits.Add(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Committer.Hex(),
				}, b)
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

				existingBidders, ok := s.bidderDeposits.Get(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				})
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
				_ = s.bidderDeposits.Add(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				}, existingBidders)
			},
		),
		events.NewEventHandler(
			"BidderDeposited",
			func(upd *bidderregistry.BidderregistryBidderDeposited) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderDeposits.Get(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				})
				if !ok {
					existing = make([]*BidderDeposit, 0)
				}

				for _, b := range existing {
					if b.Bidder == upd.Bidder.Hex() {
						return
					}
				}

				existing = append(existing, &BidderDeposit{
					Bidder:   upd.Bidder.Hex(),
					Provider: upd.Provider.Hex(),
					Amount:   upd.DepositedAmount.String(),
				})
				_ = s.bidderDeposits.Add(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				}, existing)
			},
		),
		events.NewEventHandler(
			"FundsUnlocked",
			func(upd *bidderregistry.BidderregistryFundsUnlocked) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderDeposits.Get(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				})
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
				_ = s.bidderDeposits.Add(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				}, existing)
			},
		),
		events.NewEventHandler(
			"BidderWithdrawal",
			func(upd *bidderregistry.BidderregistryBidderWithdrawal) {
				s.statMu.Lock()
				defer s.statMu.Unlock()

				existing, ok := s.bidderDeposits.Get(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				})
				if !ok {
					return
				}

				for idx, b := range existing {
					if b.Bidder == upd.Bidder.Hex() {
						existing[idx].Withdrawn = upd.AmountWithdrawn.String()
						break
					}
				}

				_ = s.bidderDeposits.Add(depositKey{
					bidder:   upd.Bidder.Hex(),
					provider: upd.Provider.Hex(),
				}, existing)
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
		Aggregate: agg,
	}
}

func (s *statHandler) getProviders() []*ProviderBalances {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	return s.providerStakes.Values()
}

func (s *statHandler) getBidders() []*BidderDeposit {
	s.statMu.RLock()
	defer s.statMu.RUnlock()

	all := make([]*BidderDeposit, 0, s.bidderDeposits.Len())

	for _, deposits := range s.bidderDeposits.Values() {
		all = append(all, deposits...)
	}
	return all
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
