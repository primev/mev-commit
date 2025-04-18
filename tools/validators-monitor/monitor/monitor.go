package monitor

import (
	"context"
	"log/slog"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/validators-monitor/api"
	"github.com/primev/mev-commit/tools/validators-monitor/config"
	"github.com/primev/mev-commit/tools/validators-monitor/contract"
	"github.com/primev/mev-commit/tools/validators-monitor/database"
	"github.com/primev/mev-commit/tools/validators-monitor/epoch"
	"github.com/primev/mev-commit/tools/validators-monitor/notification"
)

const (
	dutiesCacheTTL        = 6 * time.Hour // proposer‑duties stay hot this long
	processedBlockTTL     = 2 * time.Hour // keep “already‑seen” blocks this long
	processedBlocksTarget = 500           // hard cap after TTL pruning
	maxConcurrentFetches  = 5             // max concurrent fetches for relay data
)

type cachedDuties struct {
	duties   []api.ProposerDutyInfo
	cachedAt time.Time
}

/* external service client interfaces */

type BeaconClient interface {
	GetProposerDuties(ctx context.Context, epoch uint64) (*api.ProposerDutiesResponse, error)
	GetBlockBySlot(ctx context.Context, slot uint64) (string, error)
}

type RelayClient interface {
	QueryRelayData(ctx context.Context, blockNumber uint64) map[string]api.RelayResult
}

type DashboardClient interface {
	GetBlockInfo(ctx context.Context, blockNumber uint64) (*api.DashboardResponse, error)
}

type ValidatorOptInChecker interface {
	CheckValidatorsOptedIn(ctx context.Context, pubkeys []string) ([]contract.OptInStatus, error)
}

type SlackNotifier interface {
	NotifyRelayData(ctx context.Context,
		validatorPubkey string,
		validatorIndex uint64,
		blockNumber uint64,
		slot uint64,
		mevReward *big.Int,
		feeReceipient string,
		relaysWithData []string,
		totalRelays []string,
		blockInfo *api.DashboardResponse) error
}

type Database interface {
	SaveRelayData(ctx context.Context, record *database.RelayRecord) error
	InitSchema(ctx context.Context) error
	Close() error
}

type EpochCalculator interface {
	CurrentSlot() uint64
	CurrentEpoch() uint64
	TimeUntilNextEpoch() time.Duration
	EpochStartTime(epoch uint64) time.Time
	TargetEpoch() uint64
	EpochsToFetch() []uint64
	SlotToEpoch(slot uint64) uint64
	SetLookbackMonths(months int)
	SetMaxEpochsToFetch(max int)
}

type HistoricalState struct {
	processedEpochs  map[uint64]bool // Tracks which epochs we've fully processed
	nextEpochToFetch uint64          // Next epoch to fetch in historical range
	earliestEpoch    uint64          // Earliest epoch to consider (X months ago)
	latestEpoch      uint64          // Latest epoch to consider (usually current - offset)
	mutex            sync.Mutex      // Protects the state
}

// DutyMonitor coordinates validator monitoring operations.
type DutyMonitor struct {
	logger     *slog.Logger
	config     *config.Config
	calculator EpochCalculator

	beacon     BeaconClient
	relay      RelayClient
	dashboard  DashboardClient
	notifier   SlackNotifier
	optChecker ValidatorOptInChecker
	db         Database

	runningEpoch    uint64
	dutiesCache     map[uint64]cachedDuties // epoch → duties (+TS)
	processedBlocks map[uint64]time.Time    // blockNumber → first‑seen

	cacheMutex sync.RWMutex

	// New fields for historical data processing
	historicalState HistoricalState
	historicalMode  bool           // Whether we're in historical data collection mode
	processingWg    sync.WaitGroup // For tracking ongoing processing
}

func New(cfg *config.Config, log *slog.Logger) (*DutyMonitor, error) {
	httpClient := createRetryableHTTPClient(log)

	beaconClient, err := api.NewBeaconClient(cfg.BeaconNodeURL, log, httpClient)
	if err != nil {
		return nil, err
	}

	optInChecker, err := contract.NewValidatorOptInChecker(cfg.EthereumRPCURL, cfg.ValidatorOptInContract)
	if err != nil {
		return nil, err
	}

	dashboardClient, err := api.NewDashboardClient(cfg.DashboardApiUrl, log, httpClient)
	if err != nil {
		return nil, err
	}

	calculator := epoch.NewCalculator(
		epoch.MainnetGenesisTime,
		12, // seconds/slot
		32, // slots/epoch
		3,  // epochs to look back
	)

	if cfg.LookbackMonths > 0 {
		calculator.SetLookbackMonths(cfg.LookbackMonths)
	}

	// Set the maximum epochs to fetch in one batch
	calculator.SetMaxEpochsToFetch(cfg.MaxEpochsPerBatch)

	var db Database
	if cfg.DB.Enabled {
		dbCfg := database.Config{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			DBName:   cfg.DB.DBName,
			SSLMode:  cfg.DB.SSLMode,
		}
		db, err = database.NewPostgresDB(dbCfg, log.With("component", "database"))
		if err != nil {
			return nil, err
		}
		if err := db.InitSchema(context.Background()); err != nil {
			db.Close()
			return nil, err
		}
	}

	return &DutyMonitor{
		logger:          log,
		config:          cfg,
		calculator:      calculator,
		beacon:          beaconClient,
		relay:           api.NewRelayClient(cfg.RelayURLs, log, httpClient),
		dashboard:       dashboardClient,
		notifier:        notification.NewSlackNotifier(cfg.SlackWebhookURL, log),
		optChecker:      optInChecker,
		dutiesCache:     make(map[uint64]cachedDuties),
		processedBlocks: make(map[uint64]time.Time),
		db:              db,
		historicalMode:  cfg.LookbackMonths > 0,
		historicalState: HistoricalState{
			processedEpochs: make(map[uint64]bool),
			mutex:           sync.Mutex{},
		},
		cacheMutex: sync.RWMutex{},
	}, nil
}

func (m *DutyMonitor) Start(ctx context.Context) {
	m.runningEpoch = m.calculator.CurrentEpoch()
	if m.historicalMode {
		m.initializeHistoricalState()
		m.logger.Info("historical mode enabled",
			"lookback_months", m.config.LookbackMonths,
			"earliest_epoch", m.historicalState.earliestEpoch,
			"latest_epoch", m.historicalState.latestEpoch)
	}

	m.logger.Info("duty-monitor started",
		"epoch", m.runningEpoch,
		"interval_sec", m.config.FetchIntervalSec,
		"historical_mode", m.historicalMode)

	m.fetchAndProcessDuties(ctx) // initial fetch

	ticker := time.NewTicker(time.Duration(m.config.FetchIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("duty-monitor stopping", "reason", ctx.Err())
			m.processingWg.Wait()
			return
		case <-ticker.C:
			m.checkEpochTransition()
			m.fetchAndProcessDuties(ctx)
		}
	}
}

func (m *DutyMonitor) initializeHistoricalState() {
	calculator, ok := m.calculator.(interface {
		GetEpochForMonthsAgo(months int) uint64
	})

	if !ok {
		m.logger.Error("calculator doesn't support historical lookback")
		return
	}

	m.historicalState.mutex.Lock()
	defer m.historicalState.mutex.Unlock()

	// Calculate earliest epoch (X months ago)
	m.historicalState.earliestEpoch = calculator.GetEpochForMonthsAgo(m.config.LookbackMonths)

	// Set latest epoch to current minus offset
	m.historicalState.latestEpoch = m.calculator.TargetEpoch()

	// Start from the earliest epoch
	m.historicalState.nextEpochToFetch = m.historicalState.earliestEpoch
}

func (m *DutyMonitor) checkEpochTransition() {
	if newEpoch := m.calculator.CurrentEpoch(); newEpoch > m.runningEpoch {
		m.logger.Info("epoch transition detected", "old", m.runningEpoch, "new", newEpoch)
		m.runningEpoch = newEpoch

		if m.historicalMode {
			m.historicalState.mutex.Lock()
			m.historicalState.latestEpoch = m.calculator.TargetEpoch()
			m.historicalState.mutex.Unlock()
		}

		m.cleanupCaches()
	}
}

func (m *DutyMonitor) fetchAndProcessDuties(ctx context.Context) {
	if m.historicalMode {
		m.fetchHistoricalEpochs(ctx)
		return
	}

	for _, e := range m.calculator.EpochsToFetch() {
		m.cacheMutex.RLock()
		entry, ok := m.dutiesCache[e]
		cacheHit := ok && time.Since(entry.cachedAt) < dutiesCacheTTL
		m.cacheMutex.RUnlock()

		if cacheHit {
			m.logger.Debug("duties cache-hit", "epoch", e)
			continue
		}

		duties, err := m.fetchDutiesForEpoch(ctx, e)
		if err != nil {
			m.logger.Error("fetch duties failed", "epoch", e, "err", err)
			continue
		}

		m.cacheMutex.Lock()
		m.dutiesCache[e] = cachedDuties{duties: duties, cachedAt: time.Now()}
		m.cacheMutex.Unlock()

		m.processDuties(ctx, e, duties)
	}
}

// fetchHistoricalEpochs fetches and processes epochs from the historical range
func (m *DutyMonitor) fetchHistoricalEpochs(ctx context.Context) {
	m.historicalState.mutex.Lock()

	// If we've processed all epochs in our range, we're done
	if m.historicalState.nextEpochToFetch > m.historicalState.latestEpoch {
		m.logger.Info("historical data collection complete",
			"earliest_epoch", m.historicalState.earliestEpoch,
			"latest_epoch", m.historicalState.latestEpoch)
		m.historicalState.mutex.Unlock()
		return
	}

	// Get the next batch of epochs to fetch
	epochs := make([]uint64, 0, maxConcurrentFetches)
	for i := 0; i < maxConcurrentFetches &&
		m.historicalState.nextEpochToFetch+uint64(i) <= m.historicalState.latestEpoch; i++ {

		epoch := m.historicalState.nextEpochToFetch + uint64(i)
		// Skip already processed epochs
		if !m.historicalState.processedEpochs[epoch] {
			epochs = append(epochs, epoch)
		}
	}

	// Update next epoch to fetch
	if len(epochs) > 0 {
		m.historicalState.nextEpochToFetch += uint64(len(epochs))
	} else {
		// If no epochs were selected, move to the next batch
		m.historicalState.nextEpochToFetch += uint64(maxConcurrentFetches)
	}

	m.historicalState.mutex.Unlock()

	// Process each epoch concurrently
	for _, epoch := range epochs {
		m.processingWg.Add(1)
		go func(e uint64) {
			defer m.processingWg.Done()

			m.logger.Info("processing historical epoch", "epoch", e)

			// Check if we have this epoch in cache
			m.processingEpoch(ctx, e)

			// Mark as processed
			m.historicalState.mutex.Lock()
			m.historicalState.processedEpochs[e] = true
			m.historicalState.mutex.Unlock()
		}(epoch)
	}
}

func (m *DutyMonitor) processingEpoch(ctx context.Context, epochNum uint64) {
	// Check cache first
	m.cacheMutex.RLock()
	entry, ok := m.dutiesCache[epochNum]
	cacheHit := ok && time.Since(entry.cachedAt) < dutiesCacheTTL
	m.cacheMutex.RUnlock()

	if cacheHit {
		m.logger.Debug("historical duties cache-hit", "epoch", epochNum)
		m.processDuties(ctx, epochNum, entry.duties)
		return
	}

	// Fetch duties for this epoch
	duties, err := m.fetchDutiesForEpoch(ctx, epochNum)
	if err != nil {
		m.logger.Error("fetch historical duties failed", "epoch", epochNum, "err", err)
		return
	}

	// Cache the result
	m.cacheMutex.Lock()
	m.dutiesCache[epochNum] = cachedDuties{duties: duties, cachedAt: time.Now()}
	m.cacheMutex.Unlock()

	// Process the duties
	m.processDuties(ctx, epochNum, duties)
}

func (m *DutyMonitor) fetchDutiesForEpoch(ctx context.Context, epoch uint64) ([]api.ProposerDutyInfo, error) {
	resp, err := m.beacon.GetProposerDuties(ctx, epoch)
	if err != nil {
		return nil, err
	}
	return api.ParseProposerDuties(epoch, resp)
}

func (m *DutyMonitor) processDuties(ctx context.Context, epochNum uint64, duties []api.ProposerDutyInfo) {
	m.logger.Info("duties fetched",
		"epoch", epochNum,
		"count", len(duties),
		"start_time", m.calculator.EpochStartTime(epochNum).Format(time.RFC3339))

	opted := m.getValidatorOptInStatuses(ctx, duties)
	blocks := m.getBlocksInfoForDuties(ctx, duties)

	for _, d := range duties {
		m.processDuty(ctx, d, opted, blocks)
	}
}

func (m *DutyMonitor) getValidatorOptInStatuses(ctx context.Context, duties []api.ProposerDutyInfo) map[string]contract.OptInStatus {
	pubkeys := make([]string, len(duties))
	for i, d := range duties {
		pubkeys[i] = d.PubKey
	}
	statuses, err := m.optChecker.CheckValidatorsOptedIn(ctx, pubkeys)
	if err != nil {
		m.logger.Error("opt-in checker error", "err", err)
		return nil
	}
	out := make(map[string]contract.OptInStatus, len(pubkeys))
	for i, s := range statuses {
		out[pubkeys[i]] = s
	}
	return out
}

func (m *DutyMonitor) getBlocksInfoForDuties(ctx context.Context, duties []api.ProposerDutyInfo) map[uint64]string {
	out := make(map[uint64]string, len(duties))
	for _, d := range duties {
		bn, err := m.beacon.GetBlockBySlot(ctx, d.Slot)
		if err != nil {
			m.logger.Warn("blockBySlot error", "slot", d.Slot, "err", err)
			continue
		}
		if bn != "" {
			out[d.Slot] = bn
		}
	}
	return out
}

func (m *DutyMonitor) processDuty(
	ctx context.Context,
	duty api.ProposerDutyInfo,
	optedIn map[string]contract.OptInStatus,
	blockInfo map[uint64]string,
) {
	status, ok := optedIn[duty.PubKey]
	if !ok || (!status.IsAvsOptedIn && !status.IsMiddlewareOptedIn && !status.IsVanillaOptedIn) {
		return
	}

	blockStr, ok := blockInfo[duty.Slot]
	if !ok || blockStr == "" {
		return
	}

	bn, _ := strconv.ParseUint(blockStr, 10, 64)
	if bn == 0 {
		return
	}

	if !m.historicalMode {
		m.cacheMutex.RLock()
		_, done := m.processedBlocks[bn]
		m.cacheMutex.RUnlock()

		if done {
			return // already handled in normal mode
		}
	}

	m.processBlockData(ctx, bn, duty)

	// Only track processed blocks in normal mode
	if !m.historicalMode {
		m.cacheMutex.Lock()
		m.processedBlocks[bn] = time.Now()
		m.cacheMutex.Unlock()
	}
}

func (m *DutyMonitor) processBlockData(ctx context.Context, blockNumber uint64, duty api.ProposerDutyInfo) {
	m.logger.Info("querying relays for block",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"pubkey", duty.PubKey)

	/* query all relays */
	relayResults := m.relay.QueryRelayData(ctx, blockNumber)

	relaysWithData := []string{}
	mevReward := new(big.Int)
	var feeReceipient string
	for relayURL, result := range relayResults {
		if result.Error != "" {
			m.logger.Warn("relay query error",
				"relay", relayURL, "error", result.Error, "block", blockNumber)
			continue
		}

		bidTraces, ok := result.Response.([]api.BidTrace)
		if !ok {
			m.logger.Error("unexpected relay response type",
				"relay", relayURL, "block", blockNumber)
			continue
		}

		for _, trace := range bidTraces {
			if trace.ProposerPubkey == duty.PubKey {
				relaysWithData = append(relaysWithData, relayURL)

				m.logger.Info("relay bid for validator",
					"relay", relayURL,
					"block", blockNumber,
					"slot", duty.Slot,
					"validator_pubkey", duty.PubKey,
					"bid_value", trace.Value,
					"num_tx", trace.NumTx)

				if _, ok := mevReward.SetString(trace.Value, 10); !ok {
					m.logger.Error("parse MEV reward",
						"relay", relayURL,
						"block", blockNumber,
						"bid_value", trace.Value)
				}

				feeReceipient = trace.ProposerFeeRecipient
				break
			}
		}
	}

	/* dashboard info (optional) */
	blockInfo := m.fetchBlockInfoFromDashboard(ctx, blockNumber)

	/* notifications & DB */
	if !m.historicalMode {
		m.sendNotification(ctx, duty, blockNumber, mevReward, feeReceipient, relaysWithData, blockInfo)
	}

	if m.db != nil {
		m.saveRelayData(ctx, duty, blockNumber, mevReward, feeReceipient, relaysWithData, blockInfo)
	}

	m.logger.Info("relay data processed",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"relays_with_data", len(relaysWithData),
		"total_relays_queried", len(relayResults))
}

func (m *DutyMonitor) fetchBlockInfoFromDashboard(ctx context.Context, blockNumber uint64) *api.DashboardResponse {
	if m.dashboard == nil {
		return nil
	}

	info, err := m.dashboard.GetBlockInfo(ctx, blockNumber)
	if err != nil {
		m.logger.Error("dashboard query error",
			"block_number", blockNumber, "err", err)
		return nil
	}

	m.logger.Info("dashboard block info",
		"block_number", blockNumber,
		"winner", info.Winner,
		"total_commitments", info.TotalOpenedCommitments,
		"total_rewards", info.TotalRewards,
		"total_slashes", info.TotalSlashes,
		"total_amount", info.TotalAmount)

	return info
}

func (m *DutyMonitor) sendNotification(
	ctx context.Context,
	duty api.ProposerDutyInfo,
	blockNumber uint64,
	mevReward *big.Int,
	feeReceipient string,
	relaysWithData []string,
	blockInfo *api.DashboardResponse,
) {
	if err := m.notifier.NotifyRelayData(
		ctx,
		duty.PubKey,
		duty.ValidatorIndex,
		blockNumber,
		duty.Slot,
		mevReward,
		feeReceipient,
		relaysWithData,
		m.config.RelayURLs,
		blockInfo,
	); err != nil {
		m.logger.Error("slack notification error",
			"validator", duty.PubKey,
			"block", blockNumber,
			"err", err)
	}
}

func (m *DutyMonitor) saveRelayData(
	ctx context.Context,
	duty api.ProposerDutyInfo,
	blockNumber uint64,
	mevReward *big.Int,
	feeReceipient string,
	relaysWithData []string,
	blockInfo *api.DashboardResponse,
) {
	if m.db == nil {
		return
	}

	record := &database.RelayRecord{
		Slot:               duty.Slot,
		BlockNumber:        blockNumber,
		ValidatorIndex:     duty.ValidatorIndex,
		ValidatorPubkey:    duty.PubKey,
		MEVReward:          mevReward,
		MEVRewardRecipient: feeReceipient,
		RelaysWithData:     relaysWithData,
	}

	if blockInfo != nil {
		record.Winner = blockInfo.Winner
		record.TotalCommitments = blockInfo.TotalOpenedCommitments
		record.TotalRewards = blockInfo.TotalRewards
		record.TotalSlashes = blockInfo.TotalSlashes
		record.TotalAmount = blockInfo.TotalAmount
	}

	if err := m.db.SaveRelayData(ctx, record); err != nil {
		m.logger.Error("DB save error",
			"validator", duty.PubKey,
			"block", blockNumber,
			"err", err)
	} else {
		m.logger.Debug("relay data saved", "id", record.ID)
	}
}

func (m *DutyMonitor) cleanupCaches() {
	now := time.Now()

	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	/* duties TTL */
	for ep, entry := range m.dutiesCache {
		if now.Sub(entry.cachedAt) > dutiesCacheTTL {
			delete(m.dutiesCache, ep)
		}
	}

	/* processed blocks TTL */
	for bn, ts := range m.processedBlocks {
		if now.Sub(ts) > processedBlockTTL {
			delete(m.processedBlocks, bn)
		}
	}

	/* hard cap */
	if len(m.processedBlocks) > processedBlocksTarget {
		type kv struct {
			bn uint64
			ts time.Time
		}
		lst := make([]kv, 0, len(m.processedBlocks))
		for bn, ts := range m.processedBlocks {
			lst = append(lst, kv{bn, ts})
		}
		sort.Slice(lst, func(i, j int) bool { return lst[i].ts.Before(lst[j].ts) })
		for i := 0; len(m.processedBlocks) > processedBlocksTarget && i < len(lst); i++ {
			delete(m.processedBlocks, lst[i].bn)
		}
	}

	if m.historicalMode {
		m.historicalState.mutex.Lock()
		totalEpochs := m.historicalState.latestEpoch - m.historicalState.earliestEpoch + 1
		processedCount := len(m.historicalState.processedEpochs)
		percentComplete := float64(processedCount) / float64(totalEpochs) * 100.0
		m.logger.Info("historical processing status",
			"processed_epochs", processedCount,
			"total_epochs", totalEpochs,
			"percent_complete", percentComplete,
			"next_epoch_to_fetch", m.historicalState.nextEpochToFetch)
		m.historicalState.mutex.Unlock()
	}

	m.logger.Debug("cache sizes",
		"duties", len(m.dutiesCache),
		"blocks", len(m.processedBlocks))
}

/* helper to expose DB in tests */
func (m *DutyMonitor) GetDB() Database { return m.db }

// createRetryableHTTPClient creates a retryable HTTP client with custom settings
func createRetryableHTTPClient(log *slog.Logger) *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.RetryMax = 5
	c.RetryWaitMin = 200 * time.Millisecond
	c.RetryWaitMax = 5 * time.Second
	c.HTTPClient.Timeout = 20 * time.Second
	c.Logger = log
	return c
}
