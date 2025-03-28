package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

const (
	baseURL        = "https://ethereum-beacon-api.publicnode.com"
	execRPC        = "https://ethereum-rpc.publicnode.com"
	secondsPerSlot = 12
	slotsPerEpoch  = 32
	epochsPerDay   = 225 // Approximately 225 epochs per day (86400 / (12 * 32))
	epochsPerHour  = 9   // Approximately 9 epochs per hour (3600 / (12 * 32))
	headStatePath  = "/eth/v1/beacon/headers/head"
	proposerPath   = "/eth/v1/validator/duties/proposer"
	blockPath      = "/eth/v2/beacon/blocks"

	// Contract address for validator opt-in checking
	contractAddress = "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64"

	// Batch size for contract calls
	batchSize = 50

	// Maximum retries for API calls
	maxRetries = 3

	// Maximum retries for relay calls (increased)
	maxRelayRetries = 5

	// URLs for block explorers
	etherscanURL   = "https://etherscan.io/block/%d"
	beaconchainURL = "https://beaconcha.in/slot/%d"

	// Excel file name
	defaultExcelFile = "validator_data.xlsx"

	// Default run interval in seconds
	defaultRunInterval = 300 // 5 minutes
)

// Default relay URLs
var defaultRelays = []string{
	"https://mainnet.aestus.live",
	"https://mainnet.titanrelay.xyz",
	"https://bloxroute.max-profit.blxrbdn.com",
}

// Cache for block numbers by slot
var blockNumberCache sync.Map

// Contract ABI for the IValidatorOptInRouter
const contractABIJSON = `[
	{
		"inputs": [
			{
				"internalType": "bytes[]",
				"name": "valBLSPubKeys",
				"type": "bytes[]"
			}
		],
		"name": "areValidatorsOptedIn",
		"outputs": [
			{
				"components": [
					{
						"internalType": "bool",
						"name": "isVanillaOptedIn",
						"type": "bool"
					},
					{
						"internalType": "bool",
						"name": "isAvsOptedIn",
						"type": "bool"
					},
					{
						"internalType": "bool",
						"name": "isMiddlewareOptedIn",
						"type": "bool"
					}
				],
				"internalType": "struct IValidatorOptInRouter.OptInStatus[]",
				"name": "",
				"type": "tuple[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

type HeadResponse struct {
	Data struct {
		Header struct {
			Message struct {
				Slot string `json:"slot"`
			} `json:"message"`
		} `json:"header"`
	} `json:"data"`
}

type ProposerDutiesResponse struct {
	Data []struct {
		Pubkey string `json:"pubkey"`
		Slot   string `json:"slot"`
	} `json:"data"`
}

type BlockResponse struct {
	Data struct {
		Message struct {
			Body struct {
				ExecutionPayload struct {
					BlockNumber string `json:"block_number"`
				} `json:"execution_payload"`
			} `json:"body"`
		} `json:"message"`
	} `json:"data"`
}

// ValidationDuty holds information about a validator's proposer duty
type ValidationDuty struct {
	Pubkey           string
	Slot             uint64
	BlockNumber      uint64
	IsOptedIn        bool
	RelayData        map[string]RelayResult
	AllRelaysEmpty   bool            // New field to track if all relays returned empty results
	EtherscanURL     string          // URL to view the block on Etherscan
	BeaconchainURL   string          // URL to view the slot on Beaconchain
	RelayEmptyStatus map[string]bool // Track empty status per relay
}

// OptInStatus matches the struct in the contract
type OptInStatus struct {
	IsVanillaOptedIn    bool
	IsAvsOptedIn        bool
	IsMiddlewareOptedIn bool
}

// RelayResult represents the result of querying a relay
type RelayResult struct {
	Relay      string      `json:"relay"`
	StatusCode int         `json:"status_code"`
	Response   interface{} `json:"response,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// BidTrace represents a bid trace from a relay
type BidTrace struct {
	Slot                 string `json:"slot"`
	ParentHash           string `json:"parent_hash"`
	BlockHash            string `json:"block_hash"`
	BuilderPubkey        string `json:"builder_pubkey"`
	ProposerPubkey       string `json:"proposer_pubkey"`
	ProposerFeeRecipient string `json:"proposer_fee_recipient"`
	GasLimit             string `json:"gas_limit"`
	GasUsed              string `json:"gas_used"`
	Value                string `json:"value"`
	NumTx                string `json:"num_tx"`
	BlockNumber          string `json:"block_number"`
}

// Global HTTP client with connection pooling
var httpClient = &http.Client{
	Timeout: 30 * time.Second, // Increased from 10 to 30 seconds
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		// Add more robust transport settings
		DisableKeepAlives:     false,
		MaxConnsPerHost:       100,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	},
}

// Timeouts for specific endpoints (in seconds)
const (
	defaultTimeout  = 30
	headTimeout     = 30 // Increased from 20 to 30
	proposerTimeout = 60 // Increased from 45 to 60
	blockTimeout    = 45 // Increased from 30 to 45
	relayTimeout    = 30 // Increased from 20 to 30
)

func main() {
	// Parse command line arguments
	hoursFlag := flag.Int("hours", 24, "Number of hours to look back")
	verboseFlag := flag.Bool("verbose", true, "Enable verbose output (default true)")
	quietFlag := flag.Bool("quiet", false, "Reduce output verbosity")
	checkRelaysFlag := flag.Bool("check-relays", true, "Check relay data for blocks")
	skipRelaysFlag := flag.Bool("skip-relays", false, "Skip checking relay data")
	maxWorkersFlag := flag.Int("max-workers", runtime.NumCPU()*2, "Maximum number of concurrent workers")
	excelFlag := flag.String("excel", defaultExcelFile, "Path to Excel output file (set to 'none' to disable)")
	intervalFlag := flag.Int("interval", defaultRunInterval, "Interval between runs in seconds (0 for one-time run)")
	prettyFlag := flag.Bool("pretty", true, "Use pretty console output")
	jsonFlag := flag.Bool("json", false, "Use JSON formatted logs instead of pretty console output")
	flag.Parse()

	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if *jsonFlag {
		// Use JSON logging
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else if *prettyFlag {
		// Use pretty console logging
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}

	// Set log level based on verbose/quiet flags
	if *quietFlag {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if *verboseFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	hours := *hoursFlag
	verbose := *verboseFlag && !*quietFlag
	checkRelays := *checkRelaysFlag && !*skipRelaysFlag
	maxWorkers := *maxWorkersFlag
	excelFile := *excelFlag
	interval := *intervalFlag

	// Check if Excel export is disabled
	excelEnabled := excelFile != "none" && excelFile != ""

	// Cap max workers to avoid overwhelming APIs
	if maxWorkers > 20 {
		maxWorkers = 20
	}

	if hours <= 0 {
		log.Error().Msg("Error: hours must be positive")
		os.Exit(1)
	}

	// Set up context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signal interrupts
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Info().Msg("Received interrupt signal. Shutting down gracefully...")
		cancel()
		// Allow some time for cleanup before forcing exit
		time.Sleep(2 * time.Second)
		os.Exit(1)
	}()

	log.Info().
		Int("hours", hours).
		Int("workers", maxWorkers).
		Int("interval", interval).
		Bool("check_relays", checkRelays).
		Msg("Starting validator monitoring")

	// Run continuously if interval > 0, otherwise run once
	for {
		runStart := time.Now()
		log.Info().Msg("Starting new monitoring run")

		if err := runMonitoring(ctx, hours, verbose, checkRelays, maxWorkers, excelEnabled, excelFile); err != nil {
			log.Error().Err(err).Msg("Error during monitoring run")
		}

		runDuration := time.Since(runStart)
		log.Info().
			Dur("duration", runDuration).
			Msg("Monitoring run completed")

		if interval <= 0 {
			// One-time run, exit after completion
			break
		}

		// Calculate wait time (accounting for the time spent on this run)
		waitTime := time.Duration(interval)*time.Second - runDuration
		if waitTime <= 0 {
			waitTime = 1 * time.Second // Minimum wait time
		}

		log.Info().
			Dur("wait_time", waitTime).
			Time("next_run", time.Now().Add(waitTime)).
			Msg("Waiting for next run")

		// Wait for the next run or until interrupted
		select {
		case <-ctx.Done():
			log.Info().Msg("Monitoring stopped by request")
			return
		case <-time.After(waitTime):
			// Time to run again
		}
	}
}

// runMonitoring contains the main monitoring logic extracted from main()
func runMonitoring(ctx context.Context, hours int, verbose bool, checkRelays bool, maxWorkers int, excelEnabled bool, excelFile string) error {
	// Get current slot and epoch
	currentSlot, err := getCurrentSlot(ctx)
	if err != nil {
		return fmt.Errorf("error getting current slot: %v", err)
	}

	log.Debug().Uint64("slot", currentSlot).Msg("Current head slot")

	currentEpoch := currentSlot / slotsPerEpoch

	// Calculate start epoch based on hours instead of days
	startEpoch := currentEpoch - uint64(hours*epochsPerHour)
	if startEpoch > currentEpoch { // Overflow check
		startEpoch = 0
	}

	log.Info().
		Uint64("current_slot", currentSlot).
		Uint64("current_epoch", currentEpoch).
		Uint64("start_epoch", startEpoch).
		Int("epoch_count", int(currentEpoch-startEpoch+1)).
		Int("hours_lookback", hours).
		Msg("Processing epochs")

	// OPTIMIZATION 1: Process epochs in parallel to get validator duties
	duties, err := processEpochsInParallel(startEpoch, currentEpoch, ctx, maxWorkers, verbose)

	// Handle errors, but continue with partial data if possible
	if err != nil {
		// If context was canceled, but we have some duties, log warning and continue
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			if len(duties) > 0 {
				log.Warn().
					Err(err).
					Int("validator_count", len(duties)).
					Msg("Context canceled during epoch processing, continuing with partial validator data")
			} else {
				return fmt.Errorf("error processing epochs and no data retrieved: %v", err)
			}
		} else if strings.Contains(err.Error(), "context canceled") && len(duties) > 0 {
			// Handle string-based context canceled errors
			log.Warn().
				Err(err).
				Int("validator_count", len(duties)).
				Msg("Context cancellation error during epoch processing, continuing with partial validator data")
		} else {
			// Other errors we still return
			return fmt.Errorf("error processing epochs: %v", err)
		}
	}

	// If we have no validator duties at all, report error
	if len(duties) == 0 {
		return fmt.Errorf("no validator duties found in the specified time range")
	}

	log.Info().Int("validator_count", len(duties)).Msg("Fetched proposer duties")
	log.Info().Msg("Checking which validators are opted in to the contract...")

	// Connect to Ethereum RPC
	client, err := ethclient.Dial(execRPC)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// Parse contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	// Process validators in batches
	pubkeys := make([]string, 0, len(duties))
	for pubkey := range duties {
		pubkeys = append(pubkeys, pubkey)
	}

	// OPTIMIZATION 2: Process validators in concurrent batches
	optedInValidators := processValidatorsInBatches(pubkeys, duties, client, contractAbi, ctx, maxWorkers/2, verbose)

	log.Info().Int("opted_in_count", len(optedInValidators)).Msg("Found opted-in validators")

	// Check relay data for each opted-in validator if requested
	if checkRelays && len(optedInValidators) > 0 {
		log.Info().Int("validator_count", len(optedInValidators)).Msg("Checking relay data for opted-in validators")

		// OPTIMIZATION 3: Group validators by block number for efficient relay checks
		processRelayChecksEfficiently(optedInValidators, ctx, maxWorkers/2, verbose)

		// After relay checks, determine which validators have empty responses from all relays
		determineEmptyRelayResponses(optedInValidators)
	}

	// Output the opted-in validators
	log.Info().Msg("===== VALIDATOR REPORT =====")

	// Log header only if we have validators
	if len(optedInValidators) > 0 {
		log.Info().Msg("Validator Public Key | Slot | Block | Empty Relays | Etherscan | Beaconchain")
		log.Info().Msg("-------------------------------------------------------------------------------")
	} else {
		log.Info().Msg("No opted-in validators found in the specified time range")
	}

	for _, validator := range optedInValidators {
		emptyStatus := "No"
		if validator.AllRelaysEmpty {
			emptyStatus = "Yes"
		}

		log.Info().
			Str("pubkey", validator.Pubkey).
			Uint64("slot", validator.Slot).
			Uint64("block", validator.BlockNumber).
			Str("empty_relays", emptyStatus).
			Str("etherscan", validator.EtherscanURL).
			Str("beaconchain", validator.BeaconchainURL).
			Msg("Validator")

		// Log relay status details for each validator
		if verbose && len(validator.RelayEmptyStatus) > 0 {
			for relay, isEmpty := range validator.RelayEmptyStatus {
				status := "Has Data"
				if isEmpty {
					status = "Empty"
				}
				log.Debug().
					Str("pubkey", validator.Pubkey).
					Str("relay", relay).
					Str("status", status).
					Msg("Relay Status")
			}
		}
	}

	log.Info().Msg("===== END REPORT =====")

	// Export to Excel if enabled
	if excelEnabled && len(optedInValidators) > 0 {
		if err := exportToExcel(optedInValidators, excelFile); err != nil {
			log.Error().Err(err).Str("file", excelFile).Msg("Error exporting to Excel")
		} else {
			log.Info().Str("file", excelFile).Msg("Successfully exported data to Excel")
		}
	}

	return nil
}

// Determine which validators have empty responses from all relays
func determineEmptyRelayResponses(validators []ValidationDuty) {
	for i := range validators {
		validator := &validators[i]

		// Skip if relays weren't checked
		if len(validator.RelayData) == 0 {
			continue
		}

		allEmpty := true
		validator.RelayEmptyStatus = make(map[string]bool)

		for relayURL, result := range validator.RelayData {
			isEmpty := true

			// Check if the relay had a successful response
			if result.Error == "" && result.StatusCode == http.StatusOK {
				// Check if the response is an empty array
				if resp, ok := result.Response.([]interface{}); ok {
					isEmpty = len(resp) == 0
				} else if resp, ok := result.Response.([]BidTrace); ok {
					isEmpty = len(resp) == 0
				}
			}

			validator.RelayEmptyStatus[relayURL] = isEmpty
			if !isEmpty {
				allEmpty = false
			}
		}

		validator.AllRelaysEmpty = allEmpty
	}
}

// Process epochs in parallel with a worker pool
func processEpochsInParallel(startEpoch, endEpoch uint64, ctx context.Context, maxWorkers int, verbose bool) (map[string]uint64, error) {
	duties := make(map[string]uint64)
	var mu sync.Mutex

	// Add a counter for completed epochs
	var completedEpochs int64
	var progressMu sync.Mutex

	// Create a channel for epoch tasks
	tasks := make(chan uint64, endEpoch-startEpoch+1)

	// Create error channel with buffer for all potential errors
	errChan := make(chan error, endEpoch-startEpoch+1)

	// Create a channel to signal when processing is complete to stop the progress tracker
	processingDone := make(chan struct{})

	// Create a wait group to wait for all workers
	var wg sync.WaitGroup

	// Check if context is already canceled before starting
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context already canceled before starting epoch processing: %w", ctx.Err())
	}

	// Launch progress tracker
	if verbose {
		go func() {
			totalEpochs := endEpoch - startEpoch + 1
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.Debug().Msg("Progress tracker stopped due to context cancellation")
					return
				case <-processingDone:
					// Log final progress before exiting
					progressMu.Lock()
					completed := completedEpochs
					progressMu.Unlock()

					progress := float64(completed) / float64(totalEpochs) * 100

					log.Info().
						Int64("completed_epochs", completed).
						Int("total_epochs", int(totalEpochs)).
						Float64("progress_percent", progress).
						Msg("Final epoch processing progress - tracker stopping")
					return
				case <-ticker.C:
					progressMu.Lock()
					completed := completedEpochs
					progressMu.Unlock()

					progress := float64(completed) / float64(totalEpochs) * 100

					log.Info().
						Int64("completed_epochs", completed).
						Int("total_epochs", int(totalEpochs)).
						Float64("progress_percent", progress).
						Msg("Processing epochs progress")
				}
			}
		}()
	}

	// Launch workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Debug().Int("worker_id", workerID).Msg("Starting epoch processing worker")

			for epoch := range tasks {
				select {
				case <-ctx.Done():
					log.Debug().
						Int("worker_id", workerID).
						Err(ctx.Err()).
						Msg("Worker stopped due to context cancellation")
					return
				default:
					// Process the epoch
					log.Debug().
						Int("worker_id", workerID).
						Uint64("epoch", epoch).
						Msg("Worker processing epoch")

					epochDuties, err := getProposerDuties(ctx, epoch)

					// Check context again after call
					if ctxErr := ctx.Err(); ctxErr != nil {
						log.Debug().
							Int("worker_id", workerID).
							Uint64("epoch", epoch).
							Err(ctxErr).
							Msg("Context canceled after getting proposer duties")
						return
					}

					if err != nil {
						// Check if context was canceled
						if ctxErr := ctx.Err(); ctxErr != nil {
							log.Warn().
								Int("worker_id", workerID).
								Uint64("epoch", epoch).
								Err(ctxErr).
								Msg("Context canceled while processing epoch")
							return
						}

						// Only send non-context errors to the error channel
						log.Warn().
							Int("worker_id", workerID).
							Uint64("epoch", epoch).
							Err(err).
							Msg("Error getting proposer duties")

						errChan <- fmt.Errorf("error getting proposer duties for epoch %d: %w", epoch, err)

						// Still increment the counter even if there was an error
						progressMu.Lock()
						completedEpochs++
						progressMu.Unlock()
						continue
					}

					// Update shared map safely
					mu.Lock()
					for _, duty := range epochDuties.Data {
						slot, err := strconv.ParseUint(duty.Slot, 10, 64)
						if err != nil {
							log.Error().
								Int("worker_id", workerID).
								Str("slot", duty.Slot).
								Err(err).
								Msg("Error parsing slot")
							continue
						}
						duties[duty.Pubkey] = slot
					}
					mu.Unlock()

					// Increment completed epochs counter
					progressMu.Lock()
					completedEpochs++
					progressMu.Unlock()

					log.Debug().
						Int("worker_id", workerID).
						Uint64("epoch", epoch).
						Int("duties_count", len(epochDuties.Data)).
						Msg("Completed processing epoch")

					// Small delay to avoid rate limiting
					time.Sleep(25 * time.Millisecond)
				}
			}

			log.Debug().Int("worker_id", workerID).Msg("Worker finished processing all assigned epochs")
		}(i)
	}

	// Queue up all the epochs to process
	log.Info().
		Uint64("start_epoch", startEpoch).
		Uint64("end_epoch", endEpoch).
		Int("total_epochs", int(endEpoch-startEpoch+1)).
		Msg("Queueing epochs for processing")

	go func() {
		for epoch := startEpoch; epoch <= endEpoch; epoch++ {
			select {
			case <-ctx.Done():
				log.Info().Msg("Epoch queueing stopped due to context cancellation")
				close(tasks)
				return
			case tasks <- epoch:
				// Task queued
			}
		}
		log.Debug().Msg("All epochs queued for processing")
		close(tasks)
	}()

	// Wait for all workers to finish
	log.Debug().Msg("Waiting for all epoch processing workers to complete")
	wg.Wait()
	log.Debug().Msg("All epoch processing workers have completed")

	// Signal to the progress tracker that processing is done
	if verbose {
		close(processingDone)
	}

	// Check if context was canceled
	if err := ctx.Err(); err != nil {
		log.Warn().
			Err(err).
			Int("duties_collected", len(duties)).
			Msg("Epoch processing interrupted by context cancellation")
		return duties, fmt.Errorf("epoch processing interrupted: %w", err)
	}

	// Check for other errors
	select {
	case err := <-errChan:
		log.Error().
			Err(err).
			Int("duties_collected", len(duties)).
			Msg("Error occurred during epoch processing")
		return duties, err
	default:
		// No errors
		log.Info().
			Int("duties_collected", len(duties)).
			Msg("Epoch processing completed successfully")
		return duties, nil
	}
}

// Process validators in concurrent batches
func processValidatorsInBatches(pubkeys []string, duties map[string]uint64, client *ethclient.Client, contractAbi abi.ABI, ctx context.Context, maxWorkers int, verbose bool) []ValidationDuty {
	var optedInValidators []ValidationDuty
	var mu sync.Mutex

	// Add counters for progress tracking
	var completedBatches int64
	var optedInCount int64
	var progressMu sync.Mutex

	// Create a channel for batch tasks
	numBatches := (len(pubkeys) + batchSize - 1) / batchSize
	tasks := make(chan int, numBatches)

	// Create a wait group to wait for all workers
	var wg sync.WaitGroup

	// Launch progress tracker
	if verbose {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					progressMu.Lock()
					completed := completedBatches
					foundOptIns := optedInCount
					progressMu.Unlock()

					progress := float64(completed) / float64(numBatches) * 100

					log.Info().
						Int64("completed_batches", completed).
						Int("total_batches", numBatches).
						Float64("progress_percent", progress).
						Int64("opted_in_found", foundOptIns).
						Msg("Checking validator opt-ins progress")
				}
			}
		}()
	}

	// Launch workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batchIndex := range tasks {
				select {
				case <-ctx.Done():
					return
				default:
					start := batchIndex * batchSize
					end := start + batchSize
					if end > len(pubkeys) {
						end = len(pubkeys)
					}

					batch := pubkeys[start:end]
					if verbose {
						log.Debug().
							Int("batch", batchIndex+1).
							Int("total_batches", numBatches).
							Int("batch_size", len(batch)).
							Msg("Checking opt-in status for batch")
					}

					// Prepare blsPubKeys for contract call
					blsPubKeys := make([][]byte, len(batch))
					for j, pubkey := range batch {
						// Convert from hex string (with or without 0x prefix) to bytes
						pubkeyStr := strings.TrimPrefix(pubkey, "0x")
						pubkeyBytes, err := hex.DecodeString(pubkeyStr)
						if err != nil {
							log.Error().Str("pubkey", pubkey).Err(err).Msg("Error decoding pubkey")
							continue
						}
						blsPubKeys[j] = pubkeyBytes
					}

					// Create the contract call data
					contractAddr := common.HexToAddress(contractAddress)
					callData, err := contractAbi.Pack("areValidatorsOptedIn", blsPubKeys)
					if err != nil {
						log.Error().Err(err).Msg("Failed to pack contract call data")
						continue
					}

					// Implement retry for contract calls
					var result []byte
					var callErr error
					for attempt := 0; attempt < maxRetries; attempt++ {
						result, callErr = client.CallContract(ctx, ethereum.CallMsg{
							To:   &contractAddr,
							Data: callData,
						}, nil) // nil means latest block

						if callErr == nil {
							break
						}

						backoff := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
						if backoff > 2*time.Second {
							backoff = 2 * time.Second
						}

						log.Warn().
							Int("attempt", attempt+1).
							Int("max_retries", maxRetries).
							Err(callErr).
							Dur("backoff", backoff).
							Msg("Contract call failed, retrying")

						time.Sleep(backoff)
					}

					if callErr != nil {
						log.Error().
							Int("max_retries", maxRetries).
							Err(callErr).
							Msg("Failed to call contract after retries")

						progressMu.Lock()
						completedBatches++
						progressMu.Unlock()
						continue
					}

					// Unpack the result
					var optInStatuses []OptInStatus
					err = contractAbi.UnpackIntoInterface(&optInStatuses, "areValidatorsOptedIn", result)
					if err != nil {
						log.Error().Err(err).Msg("Failed to unpack contract result")

						progressMu.Lock()
						completedBatches++
						progressMu.Unlock()
						continue
					}

					// Process opted-in validators
					var batchResults []ValidationDuty
					var batchOptedIn int64

					for j, status := range optInStatuses {
						if j >= len(batch) {
							break // Safety check
						}

						// A validator is considered opted in if any of the flags are true
						isOptedIn := status.IsVanillaOptedIn || status.IsAvsOptedIn || status.IsMiddlewareOptedIn

						if isOptedIn {
							batchOptedIn++
							slot := duties[batch[j]]

							// Get block number for this slot with caching
							blockNumber, err := getBlockNumberForSlotWithCache(ctx, slot)
							if err != nil {
								log.Error().
									Uint64("slot", slot).
									Err(err).
									Msg("Error getting block number for slot")
								continue
							}

							// Create URLs for Etherscan and Beaconchain
							etherscanURL := fmt.Sprintf(etherscanURL, blockNumber)
							beaconchainURL := fmt.Sprintf(beaconchainURL, slot)

							batchResults = append(batchResults, ValidationDuty{
								Pubkey:         batch[j],
								Slot:           slot,
								BlockNumber:    blockNumber,
								IsOptedIn:      true,
								RelayData:      make(map[string]RelayResult),
								EtherscanURL:   etherscanURL,
								BeaconchainURL: beaconchainURL,
							})
						}
					}

					// Update shared results safely
					if len(batchResults) > 0 {
						mu.Lock()
						optedInValidators = append(optedInValidators, batchResults...)
						mu.Unlock()

						// Log details about opt-ins found in this batch
						log.Info().
							Int("batch", batchIndex+1).
							Int("validators_checked", len(batch)).
							Int("opted_in_found", len(batchResults)).
							Uint64("starting_slot", batchResults[0].Slot).
							Msg("Found opted-in validators in batch")
					}

					// Update progress counters
					progressMu.Lock()
					completedBatches++
					optedInCount += batchOptedIn
					progressMu.Unlock()

					// Small delay to avoid rate limiting
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}

	// Queue up all the batches to process
	for i := 0; i < numBatches; i++ {
		tasks <- i
	}
	close(tasks)

	// Wait for all workers to finish
	wg.Wait()

	return optedInValidators
}

// Process relay checks efficiently by grouping validators by block number
func processRelayChecksEfficiently(validators []ValidationDuty, ctx context.Context, maxWorkers int, verbose bool) {
	// Group validators by block number
	blockGroups := make(map[uint64][]int)
	for i, v := range validators {
		blockGroups[v.BlockNumber] = append(blockGroups[v.BlockNumber], i)
	}

	if verbose {
		log.Info().
			Int("validator_count", len(validators)).
			Int("unique_blocks", len(blockGroups)).
			Msg("Grouped validators by block")
	}

	// Create a channel for block tasks
	tasks := make(chan uint64, len(blockGroups))

	// Add progress tracking
	var completedBlocks int64
	var progressMu sync.Mutex

	// Create a wait group to wait for all workers
	var wg sync.WaitGroup

	// Launch progress tracker
	if verbose && len(blockGroups) > 0 {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					progressMu.Lock()
					completed := completedBlocks
					progressMu.Unlock()

					progress := float64(completed) / float64(len(blockGroups)) * 100

					log.Info().
						Int64("completed_blocks", completed).
						Int("total_blocks", len(blockGroups)).
						Float64("progress_percent", progress).
						Msg("Checking relay data progress")
				}
			}
		}()
	}

	// Launch workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blockNumber := range tasks {
				select {
				case <-ctx.Done():
					return
				default:
					indices := blockGroups[blockNumber]
					if verbose {
						log.Debug().
							Uint64("block", blockNumber).
							Int("validator_count", len(indices)).
							Msg("Checking relay data for block")
					}

					// Query relay data once for this block
					relayData := queryRelayData(blockNumber, defaultRelays)

					// Assign to all validators that proposed this block
					for _, idx := range indices {
						validators[idx].RelayData = relayData
					}

					// Update progress counter
					progressMu.Lock()
					completedBlocks++
					progressMu.Unlock()

					// Small delay to avoid rate limiting
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}

	// Queue up all the blocks to process
	for blockNumber := range blockGroups {
		tasks <- blockNumber
	}
	close(tasks)

	// Wait for all workers to finish
	wg.Wait()
}

// Export validator data to Excel file
func exportToExcel(validators []ValidationDuty, filePath string) error {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing Excel file: %v\n", err)
		}
	}()

	// Get the default sheet name
	sheetName := f.GetSheetName(0)

	// Set headers
	headers := []string{
		"Validator Public Key",
		"Slot",
		"Block Number",
		"All Relays Empty",
		"Individual Relay Status",
		"Etherscan URL",
		"Beaconchain URL",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Add data for each validator
	for i, validator := range validators {
		row := i + 2 // Start from row 2 (after headers)

		// Build relay status string
		var relayStatus strings.Builder
		for relay, isEmpty := range validator.RelayEmptyStatus {
			relayName := filepath.Base(relay)
			status := "Has Data"
			if isEmpty {
				status = "Empty"
			}
			if relayStatus.Len() > 0 {
				relayStatus.WriteString("\n")
			}
			relayStatus.WriteString(fmt.Sprintf("%s: %s", relayName, status))
		}

		// Set cell values
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), validator.Pubkey)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), validator.Slot)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), validator.BlockNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), validator.AllRelaysEmpty)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), relayStatus.String())
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), validator.EtherscanURL)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), validator.BeaconchainURL)

		// Make URLs clickable
		f.SetCellHyperLink(sheetName, fmt.Sprintf("F%d", row), validator.EtherscanURL, "External")
		f.SetCellHyperLink(sheetName, fmt.Sprintf("G%d", row), validator.BeaconchainURL, "External")
	}

	// Auto-fit columns
	for i := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, col, col, 25) // Set a reasonable default width
	}

	// Apply formatting
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DDEBF7"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating header style: %v", err)
	}

	// Apply header style
	headerRange, _ := excelize.CoordinatesToCellName(1, 1)
	// lastCol, _ := excelize.ColumnNumberToName(len(headers))
	lastHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle(sheetName, headerRange, lastHeaderCell, style)

	// Create hyperlink style
	hyperlinkStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color:     "#0563C1",
			Underline: "single",
		},
	})
	if err != nil {
		return fmt.Errorf("error creating hyperlink style: %v", err)
	}

	// Apply hyperlink style to URL columns
	for i := 2; i <= len(validators)+1; i++ {
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", i), fmt.Sprintf("G%d", i), hyperlinkStyle)
	}

	// Save the Excel file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("error saving Excel file: %v", err)
	}

	return nil
}

// Cache-aware block number getter
func getBlockNumberForSlotWithCache(ctx context.Context, slot uint64) (uint64, error) {
	// Check cache first
	if blockNumber, exists := blockNumberCache.Load(slot); exists {
		return blockNumber.(uint64), nil
	}

	// Cache miss, fetch from API
	blockNumber, err := getBlockNumberForSlot(ctx, slot)
	if err != nil {
		return 0, err
	}

	// Update cache
	blockNumberCache.Store(slot, blockNumber)

	return blockNumber, nil
}

// Enhanced with retry and context awareness, with rate limiting for 429 errors
func getCurrentSlot(ctx context.Context) (uint64, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Use a context with timeout specific for this operation
		reqCtx, cancel := context.WithTimeout(ctx, headTimeout*time.Second)

		req, err := http.NewRequestWithContext(reqCtx, "GET", baseURL+headStatePath, nil)
		if err != nil {
			cancel()
			return 0, err
		}

		resp, err = httpClient.Do(req)
		cancel() // Always cancel context to free resources

		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// Handle rate limiting (429)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()

			// Use retry headers if available, otherwise use exponential backoff
			retryAfter := 1.0 // Default 1 second
			if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
				if retrySeconds, err := strconv.ParseFloat(retryHeader, 64); err == nil && retrySeconds > 0 {
					retryAfter = retrySeconds
				}
			}

			backoff := time.Duration(retryAfter * float64(time.Second))
			log.Warn().
				Int("status_code", http.StatusTooManyRequests).
				Int("attempt", attempt+1).
				Dur("backoff", backoff).
				Msg("Rate limited by API, waiting before retry")

			time.Sleep(backoff)
			continue
		}

		if resp != nil {
			resp.Body.Close()
		}

		// Exponential backoff for other errors with increased delay
		backoff := time.Duration(math.Pow(2, float64(attempt))) * 500 * time.Millisecond
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}

		log.Warn().
			Int("attempt", attempt+1).
			Int("max_retries", maxRetries).
			Err(err).
			Dur("backoff", backoff).
			Msg("Failed to get current slot, retrying")

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(backoff):
			// Continue with retry
		}
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get current slot after %d retries: %v", maxRetries, err)
	}
	defer resp.Body.Close()

	var headResponse HeadResponse
	if err := json.NewDecoder(resp.Body).Decode(&headResponse); err != nil {
		return 0, err
	}

	slot, err := strconv.ParseUint(headResponse.Data.Header.Message.Slot, 10, 64)
	if err != nil {
		return 0, err
	}

	return slot, nil
}

func getProposerDuties(ctx context.Context, epoch uint64) (*ProposerDutiesResponse, error) {
	url := fmt.Sprintf("%s%s/%d", baseURL, proposerPath, epoch)

	var resp *http.Response
	var err error

	// First, check if context is already canceled
	if ctx.Err() != nil {
		log.Warn().
			Uint64("epoch", epoch).
			Err(ctx.Err()).
			Msg("Context already canceled before starting proposer duties request")
		return nil, ctx.Err()
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Use a context with timeout specific for proposer duties
		reqCtx, cancel := context.WithTimeout(ctx, proposerTimeout*time.Second)

		req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
		if err != nil {
			cancel()
			log.Warn().
				Uint64("epoch", epoch).
				Err(err).
				Msg("Error creating HTTP request for proposer duties")
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		// Add exponential backoff to retries
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * 500 * time.Millisecond // Increased from 100ms to 500ms
			if backoff > 5*time.Second {                                                       // Increased from 2s to 5s
				backoff = 5 * time.Second
			}

			log.Info().
				Uint64("epoch", epoch).
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Msg("Backing off before retry")

			select {
			case <-ctx.Done():
				cancel()
				log.Warn().
					Uint64("epoch", epoch).
					Err(ctx.Err()).
					Msg("Context canceled during backoff before API call")
				return nil, ctx.Err()
			case <-time.After(backoff):
				// Continue after backoff
			}
		}

		// Add request ID and user agent for better tracking
		req.Header.Set("User-Agent", "MEV-Commit-Monitor/1.0")
		req.Header.Set("X-Request-ID", fmt.Sprintf("epoch-%d-attempt-%d", epoch, attempt))

		// Log request start
		log.Debug().
			Uint64("epoch", epoch).
			Int("attempt", attempt+1).
			Str("url", url).
			Msg("Starting proposer duties request")

		startTime := time.Now()
		resp, err = httpClient.Do(req)
		requestDuration := time.Since(startTime)

		// Check context immediately after request
		if ctxErr := ctx.Err(); ctxErr != nil {
			cancel()
			if resp != nil {
				resp.Body.Close()
			}
			log.Warn().
				Uint64("epoch", epoch).
				Err(ctxErr).
				Dur("request_duration", requestDuration).
				Msg("Context canceled during HTTP request for proposer duties")
			return nil, ctxErr
		}

		cancel() // Always cancel to free resources

		// Check for successful response
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Debug().
				Uint64("epoch", epoch).
				Int("attempt", attempt+1).
				Dur("duration", requestDuration).
				Msg("Proposer duties request successful")
			break
		}

		// Handle rate limiting (429)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()

			// Use retry headers if available, otherwise use exponential backoff
			retryAfter := 2.0 // Default 2 seconds (increased from 1)
			if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
				if retrySeconds, err := strconv.ParseFloat(retryHeader, 64); err == nil && retrySeconds > 0 {
					retryAfter = retrySeconds
				}
			}

			backoff := time.Duration(retryAfter * float64(time.Second))
			log.Warn().
				Int("status_code", http.StatusTooManyRequests).
				Uint64("epoch", epoch).
				Int("attempt", attempt+1).
				Dur("backoff", backoff).
				Msg("Rate limited by API, waiting before retry")

			select {
			case <-ctx.Done():
				log.Warn().
					Uint64("epoch", epoch).
					Err(ctx.Err()).
					Msg("Context canceled while waiting after rate limit")
				return nil, ctx.Err()
			case <-time.After(backoff):
				// Continue with retry after waiting
			}
			continue
		}

		// Log detailed error information
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Warn().
				Int("status_code", resp.StatusCode).
				Str("response", string(body)).
				Dur("request_duration", requestDuration).
				Msg("Received non-200 response")
		}

		// More detailed error logging
		errorDetail := "unknown error"
		if err != nil {
			if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
				errorDetail = "timeout"
			} else {
				errorDetail = err.Error()
			}
		} else if resp != nil {
			errorDetail = fmt.Sprintf("status code %d", resp.StatusCode)
		}

		// Exponential backoff for other errors
		backoff := time.Duration(math.Pow(2, float64(attempt))) * 1 * time.Second // Increased from 500ms to 1s
		if backoff > 20*time.Second {                                             // Increased from 10s to 20s
			backoff = 20 * time.Second
		}

		log.Warn().
			Int("attempt", attempt+1).
			Int("max_retries", maxRetries).
			Str("error_type", errorDetail).
			Uint64("epoch", epoch).
			Err(err).
			Dur("backoff", backoff).
			Msg("Failed to get proposer duties, retrying")

		select {
		case <-ctx.Done():
			log.Warn().
				Uint64("epoch", epoch).
				Err(ctx.Err()).
				Msg("Context canceled while waiting to retry proposer duties request")
			return nil, ctx.Err()
		case <-time.After(backoff):
			// Continue with retry
		}
	}

	// Check context again
	if ctxErr := ctx.Err(); ctxErr != nil {
		if resp != nil {
			resp.Body.Close()
		}
		log.Warn().
			Uint64("epoch", epoch).
			Err(ctxErr).
			Msg("Context canceled after all retry attempts")
		return nil, ctxErr
	}

	// Check for HTTP client errors after all retries
	if err != nil {
		log.Error().
			Uint64("epoch", epoch).
			Err(err).
			Int("max_retries", maxRetries).
			Msg("Failed to get proposer duties after all retries")
		return nil, fmt.Errorf("failed to get proposer duties after %d retries: %w", maxRetries, err)
	}

	// Check for valid response
	if resp == nil {
		return nil, fmt.Errorf("no response received for epoch %d after %d retries", epoch, maxRetries)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %s: %s", resp.Status, string(body))
	}

	// First fully read the response body into a byte slice
	// Use local context with timeout for reading to avoid long hangs
	readCtx, readCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer readCancel()

	// Read with a separate goroutine and channel to handle context cancellation gracefully
	type readResult struct {
		data []byte
		err  error
	}
	readChan := make(chan readResult, 1)

	go func() {
		responseBody, err := io.ReadAll(resp.Body)
		readChan <- readResult{data: responseBody, err: err}
	}()

	// Wait for either the read to complete or context to be canceled
	var responseBody []byte
	select {
	case <-ctx.Done():
		// If parent context is canceled, we still want to try to read any data
		// already available before giving up completely
		select {
		case result := <-readChan:
			// If read already completed, use the data
			responseBody = result.data
			if result.err != nil {
				log.Warn().
					Uint64("epoch", epoch).
					Err(result.err).
					Msg("Error reading response body but continuing with partial data")
			}
		case <-readCtx.Done():
			// Read timeout - give up
			log.Warn().
				Uint64("epoch", epoch).
				Err(ctx.Err()).
				Msg("Context canceled and read timeout during response body read")
			return nil, ctx.Err()
		}
	case <-readCtx.Done():
		// Read timeout
		log.Warn().
			Uint64("epoch", epoch).
			Err(readCtx.Err()).
			Msg("Timeout reading response body")
		return nil, fmt.Errorf("timeout reading response body for epoch %d: %w", epoch, readCtx.Err())
	case result := <-readChan:
		// Read completed normally
		responseBody = result.data
		if result.err != nil {
			return nil, fmt.Errorf("error reading response body for epoch %d: %w", epoch, result.err)
		}
	}

	// If we got here with parent context canceled but some data, log a warning but try to continue
	if ctx.Err() != nil && len(responseBody) > 0 {
		log.Warn().
			Uint64("epoch", epoch).
			Err(ctx.Err()).
			Int("response_size", len(responseBody)).
			Msg("Context canceled but proceeding with partial response data")
	} else if ctx.Err() != nil {
		log.Warn().
			Uint64("epoch", epoch).
			Err(ctx.Err()).
			Msg("Context canceled after reading response body")
		return nil, ctx.Err()
	}

	// If we have no data even after attempting to read, fail
	if len(responseBody) == 0 {
		return nil, fmt.Errorf("empty response body for epoch %d", epoch)
	}

	// Parse response body from the byte slice (won't be affected by context cancellation)
	var duties ProposerDutiesResponse
	if err := json.Unmarshal(responseBody, &duties); err != nil {
		return nil, fmt.Errorf("error decoding response for epoch %d: %w", epoch, err)
	}

	return &duties, nil
}

func getBlockNumberForSlot(ctx context.Context, slot uint64) (uint64, error) {
	url := fmt.Sprintf("%s%s/%d", baseURL, blockPath, slot)

	var resp *http.Response
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Use a context with timeout specific for this operation
		reqCtx, cancel := context.WithTimeout(ctx, blockTimeout*time.Second)

		req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
		if err != nil {
			cancel()
			return 0, err
		}

		resp, err = httpClient.Do(req)
		cancel() // Always cancel to free resources

		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// Handle rate limiting (429)
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()

			// Use retry headers if available, otherwise use exponential backoff
			retryAfter := 1.0 // Default 1 second
			if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
				if retrySeconds, err := strconv.ParseFloat(retryHeader, 64); err == nil && retrySeconds > 0 {
					retryAfter = retrySeconds
				}
			} else {
				// Exponential backoff if no Retry-After header
				retryAfter = math.Pow(2, float64(attempt)) * 0.5 // 0.5s, 1s, 2s, 4s, 8s
				if retryAfter > 10 {
					retryAfter = 10 // Cap at 10 seconds
				}
			}

			backoff := time.Duration(retryAfter * float64(time.Second))
			log.Warn().
				Int("status_code", http.StatusTooManyRequests).
				Uint64("slot", slot).
				Int("attempt", attempt+1).
				Dur("backoff", backoff).
				Msg("Rate limited by API, waiting before retry")

			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(backoff):
				// Continue with retry after waiting
			}
			continue
		}

		if resp != nil {
			resp.Body.Close()
		}

		// Exponential backoff for other errors with increased delay
		backoff := time.Duration(math.Pow(2, float64(attempt))) * 500 * time.Millisecond
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}

		log.Warn().
			Int("attempt", attempt+1).
			Int("max_retries", maxRetries).
			Uint64("slot", slot).
			Err(err).
			Dur("backoff", backoff).
			Msg("Failed to get block for slot, retrying")

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(backoff):
			// Continue with retry
		}
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get block for slot %d after %d retries: %v", slot, maxRetries, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("HTTP error %s: %s", resp.Status, string(body))
	}

	// Use separate goroutine for reading body to handle context cancellation gracefully
	readCtx, readCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer readCancel()

	// Read with a separate goroutine and channel to handle context cancellation gracefully
	type readResult struct {
		data []byte
		err  error
	}
	readChan := make(chan readResult, 1)

	go func() {
		responseBody, err := io.ReadAll(resp.Body)
		readChan <- readResult{data: responseBody, err: err}
	}()

	// Wait for either the read to complete or context to be canceled
	var responseBody []byte
	select {
	case <-ctx.Done():
		// If parent context is canceled, we still want to try to read any data
		// already available before giving up completely
		select {
		case result := <-readChan:
			// If read already completed, use the data
			responseBody = result.data
			if result.err != nil {
				log.Warn().
					Uint64("slot", slot).
					Err(result.err).
					Msg("Error reading response body but continuing with partial data")
			}
		case <-readCtx.Done():
			// Read timeout - give up
			log.Warn().
				Uint64("slot", slot).
				Err(ctx.Err()).
				Msg("Context canceled and read timeout during response body read")
			return 0, ctx.Err()
		}
	case <-readCtx.Done():
		// Read timeout
		log.Warn().
			Uint64("slot", slot).
			Err(readCtx.Err()).
			Msg("Timeout reading response body")
		return 0, fmt.Errorf("timeout reading response body for slot %d: %w", slot, readCtx.Err())
	case result := <-readChan:
		// Read completed normally
		responseBody = result.data
		if result.err != nil {
			return 0, fmt.Errorf("error reading response body for slot %d: %w", slot, result.err)
		}
	}

	// If we got here with parent context canceled but some data, log a warning but try to continue
	if ctx.Err() != nil && len(responseBody) > 0 {
		log.Warn().
			Uint64("slot", slot).
			Err(ctx.Err()).
			Int("response_size", len(responseBody)).
			Msg("Context canceled but proceeding with partial response data")
	} else if ctx.Err() != nil {
		log.Warn().
			Uint64("slot", slot).
			Err(ctx.Err()).
			Msg("Context canceled after reading response body")
		return 0, ctx.Err()
	}

	// If we have no data even after attempting to read, fail
	if len(responseBody) == 0 {
		return 0, fmt.Errorf("empty response body for slot %d", slot)
	}

	// Parse response body from the byte slice
	var blockResponse BlockResponse
	if err := json.Unmarshal(responseBody, &blockResponse); err != nil {
		return 0, fmt.Errorf("error decoding response for slot %d: %w", slot, err)
	}

	blockNumber, err := strconv.ParseUint(blockResponse.Data.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// queryRelayData queries multiple relays for proposer payload delivered bid traces for a block
// With enhanced retry logic for problematic relays and zerolog logging
func queryRelayData(blockNumber uint64, relays []string) map[string]RelayResult {
	log.Info().Uint64("block_number", blockNumber).Msg("Querying relays for block")

	results := make(map[string]RelayResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, relay := range relays {
		wg.Add(1)
		go func(relayURL string) {
			defer wg.Done()

			// Construct API path
			apiPath := "/relay/v1/data/bidtraces/proposer_payload_delivered"

			// Handle URLs with or without trailing slash
			baseURL := strings.TrimSuffix(relayURL, "/")
			fullURL := baseURL + apiPath

			// Add query parameters
			params := url.Values{}
			params.Add("block_number", strconv.FormatUint(blockNumber, 10))
			queryURL := fullURL + "?" + params.Encode()

			// Initialize result
			result := RelayResult{
				Relay: relayURL,
			}

			// Make request with enhanced retry logic
			var resp *http.Response
			var lastErr error
			var responseBody []byte
			var statusCode int
			success := false

			// Improved retry logic specifically for problematic relays
			for attempt := 0; attempt < maxRelayRetries; attempt++ {
				// Create a new context with timeout for each attempt - longer for relay requests
				ctx, cancel := context.WithTimeout(context.Background(), relayTimeout*time.Second)

				req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
				if err != nil {
					cancel()
					lastErr = err
					continue
				}

				// Add request headers that might help with stability
				req.Header.Set("Accept", "application/json")
				req.Header.Set("Connection", "close") // Try to avoid connection reuse issues
				req.Header.Set("User-Agent", "MEV-Commit-Monitor/1.0")
				req.Header.Set("X-Request-ID", fmt.Sprintf("block-%d-relay-%s-attempt-%d",
					blockNumber, filepath.Base(relayURL), attempt))

				resp, err = httpClient.Do(req)
				cancel() // Always cancel context to free resources

				// Handle rate limiting (429)
				if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
					statusCode = resp.StatusCode
					resp.Body.Close()

					// Use retry headers if available, otherwise use exponential backoff
					retryAfter := 1.0 // Default 1 second
					if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
						if retrySeconds, err := strconv.ParseFloat(retryHeader, 64); err == nil && retrySeconds > 0 {
							retryAfter = retrySeconds
						}
					} else {
						// Exponential backoff if no Retry-After header
						retryAfter = math.Pow(2, float64(attempt)) * 0.5 // 0.5s, 1s, 2s, 4s, 8s
						if retryAfter > 10 {
							retryAfter = 10 // Cap at 10 seconds
						}
					}

					backoff := time.Duration(retryAfter * float64(time.Second))
					log.Warn().
						Int("status_code", http.StatusTooManyRequests).
						Str("relay", relayURL).
						Uint64("block", blockNumber).
						Int("attempt", attempt+1).
						Dur("backoff", backoff).
						Msg("Rate limited by relay API")

					time.Sleep(backoff)
					continue
				}

				if err == nil && resp.StatusCode == http.StatusOK {
					statusCode = resp.StatusCode
					// Successfully got response, now read body
					responseBody, err = io.ReadAll(resp.Body)
					resp.Body.Close()

					if err == nil {
						success = true
						break
					}

					lastErr = fmt.Errorf("error reading response body: %v", err)
					log.Warn().
						Str("relay", relayURL).
						Err(lastErr).
						Int("attempt", attempt+1).
						Int("max_retries", maxRelayRetries).
						Msg("Error reading relay response")
				} else {
					if resp != nil {
						statusCode = resp.StatusCode
						resp.Body.Close()
					}

					if err != nil {
						lastErr = err
					} else {
						lastErr = fmt.Errorf("HTTP status %d", statusCode)
					}

					log.Warn().
						Str("relay", relayURL).
						Err(lastErr).
						Int("attempt", attempt+1).
						Int("max_retries", maxRelayRetries).
						Int("status_code", statusCode).
						Msg("Relay request failed, retrying")
				}

				// Progressive backoff, longer for each attempt
				backoff := time.Duration(math.Pow(2, float64(attempt))) * 750 * time.Millisecond
				if backoff > 8*time.Second {
					backoff = 8 * time.Second
				}
				time.Sleep(backoff)
			}

			if !success {
				result.Error = fmt.Sprintf("Failed after %d retries: %v", maxRelayRetries, lastErr)
				log.Error().
					Str("relay", relayURL).
					Str("error", result.Error).
					Msg("Relay query failed")
			} else {
				result.StatusCode = statusCode

				// Parse response
				var response []BidTrace
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					result.Error = fmt.Sprintf("Error parsing JSON: %v", err)
					log.Error().
						Str("relay", relayURL).
						Str("error", result.Error).
						Msg("Failed to parse relay response")
				} else {
					result.Response = response

					isEmpty := len(response) == 0
					log.Info().
						Str("relay", relayURL).
						Int("status_code", statusCode).
						Int("response_items", len(response)).
						Bool("is_empty", isEmpty).
						Msg("Relay response received")
				}
			}

			// Store result
			mu.Lock()
			results[relayURL] = result
			mu.Unlock()
		}(relay)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return results
}
