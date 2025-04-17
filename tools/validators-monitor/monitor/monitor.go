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
	"github.com/primev/mev-commit/tools/validators-monitor/epoch"
	"github.com/primev/mev-commit/tools/validators-monitor/notification"
)

// DutyMonitor is responsible for periodically fetching and logging proposer duties
type DutyMonitor struct {
	logger          *slog.Logger
	client          *api.BeaconClient
	calculator      *epoch.Calculator
	config          *config.Config
	dutiesCache     map[uint64][]api.ProposerDutyInfo
	optInChecker    *contract.ValidatorOptInChecker
	relayClient     *api.RelayClient
	slackNotifier   *notification.SlackNotifier
	dashboardClient *api.DashboardClient
	runningEpoch    uint64
	processedBlocks map[uint64]struct{}
}

// New creates a new duty monitor
func New(cfg *config.Config) (*DutyMonitor, error) {
	// Use mainnet genesis time for now, could be made configurable
	calculator := epoch.NewCalculator(
		epoch.MainnetGenesisTime,
		12,
		32,
		3,
	)

	rc := retryablehttp.NewClient()
	rc.RetryMax = 5
	rc.RetryWaitMin = 200 * time.Millisecond
	rc.RetryWaitMax = 3 * time.Second
	rc.Backoff = retryablehttp.DefaultBackoff
	rc.HTTPClient.Timeout = 20 * time.Second
	rc.Logger = cfg.Logger
	// Create API client for the beacon node
	client, err := api.NewBeaconClient(cfg.BeaconNodeURL, cfg.Logger, rc)
	if err != nil {
		return nil, err
	}

	optInChecker, err := contract.NewValidatorOptInChecker(cfg.EthereumRPCURL, cfg.ValidatorOptInContract)
	if err != nil {
		cfg.Logger.Error("Failed to initialize validator opt-in checker", "error", err)
		return nil, err
	}

	relayClient := api.NewRelayClient(cfg.RelayURLs, cfg.Logger, rc)

	slackNotifier := notification.NewSlackNotifier(cfg.SlackWebhookURL, cfg.Logger)

	dashboardClient, err := api.NewDashboardClient(cfg.DashboardApiUrl, cfg.Logger, rc)
	if err != nil {
		cfg.Logger.Error("Failed to initialize dashboard client", "error", err)
		return nil, err
	}
	return &DutyMonitor{
		client:          client,
		calculator:      calculator,
		config:          cfg,
		logger:          cfg.Logger,
		optInChecker:    optInChecker,
		relayClient:     relayClient,
		slackNotifier:   slackNotifier,
		dashboardClient: dashboardClient,
		dutiesCache:     make(map[uint64][]api.ProposerDutyInfo),
		processedBlocks: make(map[uint64]struct{}),
	}, nil
}

// Start begins the monitoring process
func (m *DutyMonitor) Start(ctx context.Context) {
	// Log our starting state
	currentEpoch := m.calculator.CurrentEpoch()
	m.runningEpoch = currentEpoch
	m.logger.Info("Duty monitor starting",
		"current_epoch", currentEpoch,
		"fetch_interval", m.config.FetchIntervalSec)

	// Initial fetch of proposer duties
	m.fetchAndLogDuties(ctx)

	// Set up single ticker for all operations
	ticker := time.NewTicker(time.Duration(m.config.FetchIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Duty monitor stopping", "reason", ctx.Err())
			return

		case <-ticker.C:
			// Check for epoch transition
			newEpoch := m.calculator.CurrentEpoch()
			if newEpoch > m.runningEpoch {
				// Log the epoch transition
				m.logger.Info("Epoch transition occurred",
					"previous_epoch", m.runningEpoch,
					"new_epoch", newEpoch)

				// Update running epoch
				m.runningEpoch = newEpoch

				// Clean up older cached duties
				m.cleanupCache(newEpoch)
			}

			// Regular interval fetch (happens on every tick)
			m.fetchAndLogDuties(ctx)
		}
	}
}

// fetchAndLogDuties fetches proposer duties for upcoming epochs and logs them
func (m *DutyMonitor) fetchAndLogDuties(ctx context.Context) {
	// Get the current epoch and determine which epochs to fetch
	currentEpoch := m.calculator.CurrentEpoch()
	m.runningEpoch = currentEpoch
	epochsToFetch := m.calculator.EpochsToFetch()

	// Log start of fetch operation
	m.logger.Debug("Starting duties fetch operation",
		"current_epoch", currentEpoch,
		"target_epochs", epochsToFetch)

	// For each target epoch, fetch the proposer duties
	for _, targetEpoch := range epochsToFetch {
		// Check if we already have this epoch's duties in cache
		if _, exists := m.dutiesCache[targetEpoch]; exists {
			m.logger.Debug("Using cached duties", "epoch", targetEpoch)
			continue
		}

		// Fetch duties for this epoch
		resp, err := m.client.GetProposerDuties(ctx, targetEpoch)
		if err != nil {
			m.logger.Error("Failed to fetch proposer duties",
				"epoch", targetEpoch,
				"error", err)
			continue
		}

		// Parse the duties into a more usable format
		duties, err := api.ParseProposerDuties(targetEpoch, resp)
		if err != nil {
			m.logger.Error("Failed to parse proposer duties",
				"epoch", targetEpoch,
				"error", err)
			continue
		}

		// Cache the duties
		m.dutiesCache[targetEpoch] = duties

		// Log the duties
		m.logDuties(targetEpoch, duties)
	}
}

// logDuties logs the proposer duties for an epoch
func (m *DutyMonitor) logDuties(epochNum uint64, duties []api.ProposerDutyInfo) {
	epochStr := epoch.FormatEpoch(epochNum)

	// Log summary information
	m.logger.Info("Proposer duties summary",
		"epoch", epochNum,
		"epoch_id", epochStr,
		"duties_count", len(duties),
		"start_time", m.calculator.EpochStartTime(epochNum).Format(time.RFC3339))

	// Extract pubkeys
	pubkeys := make([]string, len(duties))
	for i, duty := range duties {
		pubkeys[i] = duty.PubKey
	}

	// Check opt-in status for all validators in this epoch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	statuses, err := m.optInChecker.CheckValidatorsOptedIn(ctx, pubkeys)
	cancel()

	optedInStatuses := make(map[string]contract.OptInStatus)
	if err != nil {
		m.logger.Error("Failed to check validators opt-in status", "error", err)
	} else {
		// Map results to pubkeys
		for i, status := range statuses {
			if i < len(pubkeys) {
				optedInStatuses[pubkeys[i]] = status
			}
		}
	}

	blockInfo := make(map[uint64]string) // slot -> block number
	for _, duty := range duties {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		blockNumber, err := m.client.GetBlockBySlot(ctx, duty.Slot)
		cancel()

		if err != nil {
			m.logger.Error("Failed to fetch block information",
				"slot", duty.Slot,
				"error", err)
		} else {
			blockInfo[duty.Slot] = blockNumber
		}
	}

	// Log individual duties with detailed information
	for _, duty := range duties {
		status, hasStatus := optedInStatuses[duty.PubKey]
		if status.IsAvsOptedIn || status.IsMiddlewareOptedIn || status.IsVanillaOptedIn {
			logValues := []interface{}{
				"epoch", epochNum,
				"epoch_id", epochStr,
				"slot", duty.Slot,
				"validator_index", duty.ValidatorIndex,
				"pubkey", duty.PubKey,
				"vanilla_opted_in", hasStatus && status.IsVanillaOptedIn,
				"avs_opted_in", hasStatus && status.IsAvsOptedIn,
				"middleware_opted_in", hasStatus && status.IsMiddlewareOptedIn,
			}

			if blockNumber, ok := blockInfo[duty.Slot]; ok {
				if blockNumber == "" {
					logValues = append(logValues, "block_missed", true)
				} else {
					logValues = append(logValues, "block_number", blockNumber, "block_missed", false)

					bn, _ := strconv.ParseUint(blockNumber, 10, 64)
					if bn > 0 {
						// Only check relay data if we haven't processed this block before
						if _, alreadyProcessed := m.processedBlocks[bn]; !alreadyProcessed {
							// Check relay data for this block/slot
							m.checkRelayDataForBlock(context.Background(), bn, duty)
							// Mark this block as processed to avoid duplicates
							m.processedBlocks[bn] = struct{}{}
						} else {
							m.logger.Debug("Skipping already processed block",
								"block_number", bn,
								"slot", duty.Slot,
								"validator_index", duty.ValidatorIndex)
						}
					}
				}
			}
			m.logger.Info("Proposer duty", logValues...)
		}
	}
}

func (m *DutyMonitor) checkRelayDataForBlock(ctx context.Context, blockNumber uint64, duty api.ProposerDutyInfo) {
	m.logger.Info("Querying relays for block data",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"pubkey", duty.PubKey)

	// Query relays for this block
	relayResults := m.relayClient.QueryRelayData(ctx, blockNumber)

	// Count relays that have data for this block
	relaysWithData := 0
	relaysWithDataList := []string{}
	mevReward := new(big.Int)

	// Process results from each relay
	for relayURL, result := range relayResults {
		if result.Error != "" {
			m.logger.Warn("Relay query error",
				"relay", relayURL,
				"error", result.Error,
				"block", blockNumber)
			continue
		}

		// Check if we got bid traces
		bidTraces, ok := result.Response.([]api.BidTrace)
		if !ok {
			m.logger.Error("Unexpected response type from relay",
				"relay", relayURL,
				"block", blockNumber)
			continue
		}

		// Check if any bid traces match our validator
		for _, trace := range bidTraces {
			if trace.ProposerPubkey == duty.PubKey {
				relaysWithData++
				relaysWithDataList = append(relaysWithDataList, relayURL)
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
						"slot", duty.Slot,
						"validator_pubkey", duty.PubKey,
						"bid_value", trace.Value,
					)
				}
				break
			}
		}
	}

	blockInfo := m.fetchBlockInfo(ctx, blockNumber)

	notifyCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.slackNotifier.NotifyRelayData(
		notifyCtx,
		duty.PubKey,
		duty.ValidatorIndex,
		blockNumber,
		duty.Slot,
		mevReward,
		relaysWithDataList,
		m.config.RelayURLs,
		blockInfo,
	)

	if err != nil {
		m.logger.Error("Failed to send relay data notification",
			"error", err,
			"validator", duty.PubKey,
			"block", blockNumber)
	}

	m.logger.Info("Relay data query complete",
		"block_number", blockNumber,
		"slot", duty.Slot,
		"validator_index", duty.ValidatorIndex,
		"relays_with_data", relaysWithData,
		"total_relays_queried", len(relayResults))
}

// cleanupCache removes old epochs from the cache
func (m *DutyMonitor) cleanupCache(currentEpoch uint64) {
	// Keep only current and future epochs
	for cachedEpoch := range m.dutiesCache {
		if cachedEpoch < currentEpoch {
			delete(m.dutiesCache, cachedEpoch)
			m.logger.Debug("Removed old epoch from cache", "epoch", cachedEpoch)
		}
	}

	// Clean up processed blocks map to prevent memory growth
	// Production approach: keep only reasonably recent blocks
	// by removing blocks that are too old
	if len(m.processedBlocks) > 500 {
		m.logger.Debug("Cleaning up processed blocks cache",
			"before_cleanup", len(m.processedBlocks))

		// Get the active duty range from cache
		minEpoch := currentEpoch
		for epoch := range m.dutiesCache {
			if epoch < minEpoch {
				minEpoch = epoch
			}
		}

		// Determine a reasonable block threshold based on epochs
		// Ethereum has ~225 blocks per epoch (12s slots, 32 slots per epoch)
		// We'll delete blocks that are more than 3 epochs before our earliest cached epoch
		minBlockThreshold := uint64(0)
		if minEpoch > 3 {
			// ~225 blocks per epoch × 3 epochs before our earliest epoch
			minBlockThreshold = (minEpoch - 3) * 225
		}

		// Remove blocks that are definitely old
		for blockNum := range m.processedBlocks {
			if blockNum < minBlockThreshold {
				delete(m.processedBlocks, blockNum)
			}
		}

		m.logger.Debug("Processed blocks cache cleaned",
			"after_cleanup", len(m.processedBlocks),
			"blocks_below", minBlockThreshold)
	}
}

func (m *DutyMonitor) fetchBlockInfo(ctx context.Context, blockNumber uint64) *api.DashboardResponse {
	if m.dashboardClient == nil {
		return nil
	}

	// Create a context with timeout for the API call
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Query dashboard service for block info
	blockInfo, err := m.dashboardClient.GetBlockInfo(queryCtx, blockNumber)
	if err != nil {
		m.logger.Error("Failed to fetch block info from dashboard service",
			"error", err,
			"block_number", blockNumber)
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
