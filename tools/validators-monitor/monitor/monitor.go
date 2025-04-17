package monitor

import (
	"context"
	"log/slog"
	"math/big"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/validators-monitor/api"
	"github.com/primev/mev-commit/tools/validators-monitor/config"
	"github.com/primev/mev-commit/tools/validators-monitor/contract"
	"github.com/primev/mev-commit/tools/validators-monitor/database"
	"github.com/primev/mev-commit/tools/validators-monitor/epoch"
	"github.com/primev/mev-commit/tools/validators-monitor/notification"
)

// DutyMonitor coordinates validator monitoring operations
type DutyMonitor struct {
	logger       *slog.Logger
	config       *config.Config
	calculator   *epoch.Calculator
	runningEpoch uint64

	// Service clients
	beacon     *api.BeaconClient
	relay      *api.RelayClient
	dashboard  *api.DashboardClient
	notifier   *notification.SlackNotifier
	optChecker *contract.ValidatorOptInChecker

	// Caches
	dutiesCache     map[uint64][]api.ProposerDutyInfo
	processedBlocks map[uint64]struct{}

	// DB
	db *database.PostgresDB
}

// New creates a new duty monitor with all required dependencies
func New(cfg *config.Config, log *slog.Logger) (*DutyMonitor, error) {
	// Initialize HTTP client for all services
	httpClient := createRetryableHTTPClient(log)

	// Initialize services
	beaconClient, err := api.NewBeaconClient(cfg.BeaconNodeURL, log, httpClient)
	if err != nil {
		return nil, err
	}

	optInChecker, err := contract.NewValidatorOptInChecker(cfg.EthereumRPCURL, cfg.ValidatorOptInContract)
	if err != nil {
		log.Error("Failed to initialize validator opt-in checker", "error", err)
		return nil, err
	}

	dashboardClient, err := api.NewDashboardClient(cfg.DashboardApiUrl, log, httpClient)
	if err != nil {
		log.Error("Failed to initialize dashboard client", "error", err)
		return nil, err
	}

	// Create calculator with mainnet default values
	calculator := epoch.NewCalculator(
		epoch.MainnetGenesisTime,
		12, // seconds per slot
		32, // slots per epoch
		3,  // epochs to look behind
	)

	var db *database.PostgresDB
	if cfg.DB.Enabled {
		dbConfig := database.Config{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			DBName:   cfg.DB.DBName,
			SSLMode:  cfg.DB.SSLMode,
		}

		db, err = database.NewPostgresDB(dbConfig, log.With("component", "database"))
		if err != nil {
			log.Error("Failed to connect to database", "error", err)
			return nil, err
		}

		// Initialize database schema
		if err := db.InitSchema(context.Background()); err != nil {
			log.Error("Failed to initialize database schema", "error", err)
			db.Close()
			return nil, err
		}
	} else {
		log.Info("Database is disabled, relay data will not be saved")
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
		dutiesCache:     make(map[uint64][]api.ProposerDutyInfo),
		processedBlocks: make(map[uint64]struct{}),
		db:              db,
	}, nil
}

// createRetryableHTTPClient creates a configured HTTP client with retry logic
func createRetryableHTTPClient(log *slog.Logger) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 5
	client.RetryWaitMin = 200 * time.Millisecond
	client.RetryWaitMax = 3 * time.Second
	client.Backoff = retryablehttp.DefaultBackoff
	client.HTTPClient.Timeout = 20 * time.Second
	client.Logger = log
	return client
}

// Start begins the monitoring process
func (m *DutyMonitor) Start(ctx context.Context) {
	// Initialize state
	m.runningEpoch = m.calculator.CurrentEpoch()
	m.logger.Info("Duty monitor starting",
		"current_epoch", m.runningEpoch,
		"fetch_interval", m.config.FetchIntervalSec)

	// Initial fetch
	m.fetchAndProcessDuties(ctx)

	// Set up ticker for periodic tasks
	ticker := time.NewTicker(time.Duration(m.config.FetchIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Duty monitor stopping", "reason", ctx.Err())
			return

		case <-ticker.C:
			m.checkEpochTransition()
			m.fetchAndProcessDuties(ctx)
		}
	}
}

// checkEpochTransition detects and handles epoch transitions
func (m *DutyMonitor) checkEpochTransition() {
	newEpoch := m.calculator.CurrentEpoch()
	if newEpoch > m.runningEpoch {
		m.logger.Info("Epoch transition occurred",
			"previous_epoch", m.runningEpoch,
			"new_epoch", newEpoch)

		m.runningEpoch = newEpoch
		m.cleanupCache(newEpoch)
	}
}

// fetchAndProcessDuties gets and processes proposer duties
func (m *DutyMonitor) fetchAndProcessDuties(ctx context.Context) {
	currentEpoch := m.calculator.CurrentEpoch()
	m.runningEpoch = currentEpoch
	epochsToFetch := m.calculator.EpochsToFetch()

	m.logger.Debug("Starting duties fetch operation",
		"current_epoch", currentEpoch,
		"target_epochs", epochsToFetch)

	for _, targetEpoch := range epochsToFetch {
		// Skip if already cached
		if _, exists := m.dutiesCache[targetEpoch]; exists {
			m.logger.Debug("Using cached duties", "epoch", targetEpoch)
			continue
		}

		// Fetch and process new duties
		duties, err := m.fetchDutiesForEpoch(ctx, targetEpoch)
		if err != nil {
			m.logger.Error("Failed to fetch/parse duties", "epoch", targetEpoch, "error", err)
			continue
		}

		// Cache and process
		m.dutiesCache[targetEpoch] = duties
		m.processDuties(ctx, targetEpoch, duties)
	}
}

// fetchDutiesForEpoch retrieves and parses proposer duties for a specific epoch
func (m *DutyMonitor) fetchDutiesForEpoch(ctx context.Context, targetEpoch uint64) ([]api.ProposerDutyInfo, error) {
	resp, err := m.beacon.GetProposerDuties(ctx, targetEpoch)
	if err != nil {
		return nil, err
	}

	return api.ParseProposerDuties(targetEpoch, resp)
}

// processDuties handles newly fetched proposer duties
func (m *DutyMonitor) processDuties(ctx context.Context, epochNum uint64, duties []api.ProposerDutyInfo) {
	epochStr := epoch.FormatEpoch(epochNum)

	m.logger.Info("Proposer duties summary",
		"epoch", epochNum,
		"epoch_id", epochStr,
		"duties_count", len(duties),
		"start_time", m.calculator.EpochStartTime(epochNum).Format(time.RFC3339))

	// Get opt-in status for validators
	optedInStatuses := m.getValidatorOptInStatuses(ctx, duties)

	// Get block info for slots
	blocksInfo := m.getBlocksInfoForDuties(ctx, duties)

	// Process individual duties
	for _, duty := range duties {
		m.processDuty(ctx, epochNum, epochStr, duty, optedInStatuses, blocksInfo)
	}
}

// getValidatorOptInStatuses gets opt-in status for all validators
func (m *DutyMonitor) getValidatorOptInStatuses(ctx context.Context, duties []api.ProposerDutyInfo) map[string]contract.OptInStatus {
	// Extract pubkeys
	pubkeys := make([]string, len(duties))
	for i, duty := range duties {
		pubkeys[i] = duty.PubKey
	}

	// Check all pubkeys in one batch operation
	statuses, err := m.optChecker.CheckValidatorsOptedIn(ctx, pubkeys)
	if err != nil {
		m.logger.Error("Failed to check validators opt-in status", "error", err)
		return make(map[string]contract.OptInStatus)
	}

	// Map results to pubkeys
	result := make(map[string]contract.OptInStatus)
	for i, status := range statuses {
		if i < len(pubkeys) {
			result[pubkeys[i]] = status
		}
	}

	return result
}

// getBlocksInfoForDuties gets block numbers for all slots
func (m *DutyMonitor) getBlocksInfoForDuties(ctx context.Context, duties []api.ProposerDutyInfo) map[uint64]string {
	blockInfo := make(map[uint64]string) // slot -> block number

	for _, duty := range duties {
		blockNumber, err := m.beacon.GetBlockBySlot(ctx, duty.Slot)

		if err != nil {
			m.logger.Error("Failed to fetch block information",
				"slot", duty.Slot, "error", err)
		} else {
			blockInfo[duty.Slot] = blockNumber
		}
	}

	return blockInfo
}

// processDuty handles a single proposer duty
func (m *DutyMonitor) processDuty(
	ctx context.Context,
	epochNum uint64,
	epochStr string,
	duty api.ProposerDutyInfo,
	optedInStatuses map[string]contract.OptInStatus,
	blockInfo map[uint64]string,
) {
	// Check if we need to process this validator
	status, hasStatus := optedInStatuses[duty.PubKey]
	if !hasStatus || (!status.IsAvsOptedIn && !status.IsMiddlewareOptedIn && !status.IsVanillaOptedIn) {
		return // Skip validators that haven't opted in
	}

	// Prepare log data
	logData := []interface{}{
		"epoch", epochNum,
		"epoch_id", epochStr,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"pubkey", duty.PubKey,
		"vanilla_opted_in", hasStatus && status.IsVanillaOptedIn,
		"avs_opted_in", hasStatus && status.IsAvsOptedIn,
		"middleware_opted_in", hasStatus && status.IsMiddlewareOptedIn,
	}

	// Add block information if available
	if blockNumberStr, ok := blockInfo[duty.Slot]; ok {
		if blockNumberStr == "" {
			logData = append(logData, "block_missed", true)
		} else {
			logData = append(logData, "block_number", blockNumberStr, "block_missed", false)

			// Process block if it exists and we haven't seen it before
			bn, _ := strconv.ParseUint(blockNumberStr, 10, 64)
			if bn > 0 {
				if _, alreadyProcessed := m.processedBlocks[bn]; !alreadyProcessed {
					m.processBlockData(ctx, bn, duty)
					m.processedBlocks[bn] = struct{}{}
				} else {
					m.logger.Debug("Skipping already processed block",
						"block_number", bn, "slot", duty.Slot)
				}
			}
		}
	}

	m.logger.Info("Proposer duty", logData...)
}

// processBlockData gets and processes relay data for a block
func (m *DutyMonitor) processBlockData(ctx context.Context, blockNumber uint64, duty api.ProposerDutyInfo) {
	m.logger.Info("Querying relays for block data",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"pubkey", duty.PubKey)

	// Query all relays in parallel
	relayResults := m.relay.QueryRelayData(ctx, blockNumber)

	// Extract relevant relay data
	relaysWithData := []string{}
	mevReward := new(big.Int)

	for relayURL, result := range relayResults {
		if result.Error != "" {
			m.logger.Warn("Relay query error",
				"relay", relayURL, "error", result.Error, "block", blockNumber)
			continue
		}

		bidTraces, ok := result.Response.([]api.BidTrace)
		if !ok {
			m.logger.Error("Unexpected response type from relay",
				"relay", relayURL, "block", blockNumber)
			continue
		}

		for _, trace := range bidTraces {
			if trace.ProposerPubkey == duty.PubKey {
				relaysWithData = append(relaysWithData, relayURL)

				m.logger.Info("Found relay data for validator",
					"relay", relayURL,
					"block", blockNumber,
					"slot", duty.Slot,
					"validator_pubkey", duty.PubKey,
					"bid_value", trace.Value,
					"num_tx", trace.NumTx)

				if _, ok := mevReward.SetString(trace.Value, 10); !ok {
					m.logger.Error("Failed to parse MEV reward value",
						"relay", relayURL,
						"block", blockNumber,
						"validator_pubkey", duty.PubKey,
						"bid_value", trace.Value)
				}
				break
			}
		}
	}

	// Get additional block info from dashboard
	blockInfo := m.fetchBlockInfoFromDashboard(ctx, blockNumber)

	// Send notification
	m.sendNotification(ctx, duty, blockNumber, mevReward, relaysWithData, blockInfo)

	if m.db != nil {
		m.saveRelayData(ctx, duty, blockNumber, mevReward, relaysWithData, blockInfo)
	}

	m.logger.Info("Relay data query complete",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"relays_with_data", len(relaysWithData),
		"total_relays_queried", len(relayResults))
}

// fetchBlockInfoFromDashboard gets block information from dashboard service
func (m *DutyMonitor) fetchBlockInfoFromDashboard(ctx context.Context, blockNumber uint64) *api.DashboardResponse {
	if m.dashboard == nil {
		return nil
	}

	blockInfo, err := m.dashboard.GetBlockInfo(ctx, blockNumber)
	if err != nil {
		m.logger.Error("Failed to fetch block info from dashboard service",
			"error", err, "block_number", blockNumber)
		return nil
	}

	m.logger.Info("Block info fetched from dashboard service",
		"block_number", blockNumber,
		"winner", blockInfo.Winner,
		"total_commitments", blockInfo.TotalOpenedCommitments,
		"total_rewards", blockInfo.TotalRewards,
		"total_slashes", blockInfo.TotalSlashes,
		"total_amount", blockInfo.TotalAmount)

	return blockInfo
}

// sendNotification sends information about relay data
func (m *DutyMonitor) sendNotification(
	ctx context.Context,
	duty api.ProposerDutyInfo,
	blockNumber uint64,
	mevReward *big.Int,
	relaysWithData []string,
	blockInfo *api.DashboardResponse,
) {
	err := m.notifier.NotifyRelayData(
		ctx,
		duty.PubKey,
		duty.ValidatorIndex,
		blockNumber,
		duty.Slot,
		mevReward,
		relaysWithData,
		m.config.RelayURLs,
		blockInfo,
	)

	if err != nil {
		m.logger.Error("Failed to send relay data notification",
			"error", err,
			"validator", duty.PubKey,
			"block", blockNumber)
	}
}

// cleanupCache removes old epochs and blocks from the cache
func (m *DutyMonitor) cleanupCache(currentEpoch uint64) {
	// Remove old epochs
	for cachedEpoch := range m.dutiesCache {
		if cachedEpoch < currentEpoch {
			delete(m.dutiesCache, cachedEpoch)
			m.logger.Debug("Removed old epoch from cache", "epoch", cachedEpoch)
		}
	}

	// Clean up processed blocks if we have too many
	if len(m.processedBlocks) <= 500 {
		return
	}

	m.logger.Debug("Cleaning up processed blocks cache",
		"before_cleanup", len(m.processedBlocks))

	// Find minimum epoch in cache
	minEpoch := currentEpoch
	for epoch := range m.dutiesCache {
		if epoch < minEpoch {
			minEpoch = epoch
		}
	}

	// Delete blocks older than 3 epochs before our earliest cached epoch
	// (Ethereum has ~225 blocks per epoch - 12s slots, 32 slots per epoch)
	minBlockThreshold := uint64(0)
	if minEpoch > 3 {
		minBlockThreshold = (minEpoch - 3) * 225
	}

	// Remove old blocks
	for blockNum := range m.processedBlocks {
		if blockNum < minBlockThreshold {
			delete(m.processedBlocks, blockNum)
		}
	}

	m.logger.Debug("Processed blocks cache cleaned",
		"after_cleanup", len(m.processedBlocks),
		"blocks_below", minBlockThreshold)
}

func (m *DutyMonitor) saveRelayData(
	ctx context.Context,
	duty api.ProposerDutyInfo,
	blockNumber uint64,
	mevReward *big.Int,
	relaysWithData []string,
	blockInfo *api.DashboardResponse,
) {
	if m.db == nil {
		return
	}

	record := &database.RelayRecord{
		Slot:            duty.Slot,
		BlockNumber:     blockNumber,
		ValidatorIndex:  duty.ValidatorIndex,
		ValidatorPubkey: duty.PubKey,
		MEVReward:       mevReward,
		RelaysWithData:  relaysWithData,
	}

	// Add dashboard info if available
	if blockInfo != nil {
		record.Winner = blockInfo.Winner
		record.TotalCommitments = blockInfo.TotalOpenedCommitments
		record.TotalRewards = blockInfo.TotalRewards
		record.TotalSlashes = blockInfo.TotalSlashes
		record.TotalAmount = blockInfo.TotalAmount
	}

	err := m.db.SaveRelayData(ctx, record)
	if err != nil {
		m.logger.Error("Failed to save relay data to database",
			"error", err,
			"validator", duty.PubKey,
			"block", blockNumber)
	} else {
		m.logger.Debug("Saved relay data to database",
			"id", record.ID,
			"validator", duty.PubKey,
			"block", blockNumber)
	}
}

func (m *DutyMonitor) GetDB() *database.PostgresDB {
	return m.db
}
