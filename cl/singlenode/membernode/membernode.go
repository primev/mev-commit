package membernode

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/singlenode/api"
)

const (
	shutdownTimeout      = 5 * time.Second
	maxConsecutiveErrors = 5
	batchSize            = 10
	maxCatchupPayloads   = 100
	// Timeout for API calls to leader and Geth
	apiCallTimeout = 30 * time.Second
	// Interval for retrying initialization steps
	initRetryInterval = 2 * time.Second
	// Threshold to exit catch-up mode
	catchUpExitThreshold = batchSize / 2
)

// Config holds the configuration for the MemberNodeApp
type Config struct {
	InstanceID   string
	LeaderAPIURL string
	EthClientURL string
	JWTSecret    string
	HealthAddr   string
	PollInterval time.Duration
}

// MemberNodeApp represents a member node that follows the leader sequentially
type MemberNodeApp struct {
	logger           *slog.Logger
	cfg              Config
	blockBuilder     *blockbuilder.BlockBuilder
	payloadClient    *api.PayloadClient
	engineClient     blockbuilder.EngineClient
	appCtx           context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	connectionStatus sync.Mutex
	leaderAvailable  bool
	initializedCh    chan struct{}

	// Sequential processing state
	processingMutex     sync.RWMutex
	lastProcessedHeight uint64
	isCatchingUp        bool
	isInitialized       bool
}

// NewMemberNodeApp creates and initializes a new MemberNodeApp
func NewMemberNodeApp(
	parentCtx context.Context,
	cfg Config,
	logger *slog.Logger,
) (*MemberNodeApp, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	// Decode JWT secret
	jwtBytes, err := hex.DecodeString(cfg.JWTSecret)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to decode JWT secret: %w", err)
	}

	// Create Ethereum engine client
	engineClient, err := ethclient.NewAuthClient(ctx, cfg.EthClientURL, jwtBytes)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Ethereum engine client: %w", err)
	}

	// Create block builder for member node
	bb := blockbuilder.NewMemberBlockBuilder(engineClient, logger.With("component", "BlockBuilder"))

	// Create payload client
	payloadClient := api.NewPayloadClient(cfg.LeaderAPIURL, logger)

	return &MemberNodeApp{
		logger:              logger,
		cfg:                 cfg,
		blockBuilder:        bb,
		payloadClient:       payloadClient,
		engineClient:        engineClient,
		appCtx:              ctx,
		cancel:              cancel,
		initializedCh:       make(chan struct{}),
		leaderAvailable:     false,
		lastProcessedHeight: 0,
		isCatchingUp:        false,
		isInitialized:       false,
	}, nil
}

// getLocalGethHeight gets the current block height from local geth
func (app *MemberNodeApp) getLocalGethHeight(ctx context.Context) (uint64, error) {
	header, err := app.engineClient.HeaderByNumber(ctx, nil) // nil = latest
	if err != nil {
		return 0, fmt.Errorf("failed to get latest header from local geth: %w", err)
	}

	height := header.Number.Uint64()
	app.logger.Debug("Retrieved local geth height", "height", height)
	return height, nil
}

// initializeStartingHeight determines the starting height from local geth
func (app *MemberNodeApp) initializeStartingHeight() {
	defer close(app.initializedCh) // Signal completion regardless of outcome (or handle errors preventing it)

	app.logger.Info("Detecting starting height from local geth...")

	for {
		select {
		case <-app.appCtx.Done():
			app.logger.Info("Initialization cancelled.")
			return
		default:
			ctx, cancelTimeout := context.WithTimeout(app.appCtx, apiCallTimeout)

			// Check leader availability first
			if err := app.payloadClient.CheckHealth(ctx); err != nil {
				cancelTimeout()
				app.logger.Warn("Leader not available during initialization, retrying...", "error", err)
				select {
				case <-app.appCtx.Done():
					return
				case <-time.After(initRetryInterval):
					continue
				}
			}

			// Get local geth's current height
			localHeight, err := app.getLocalGethHeight(ctx)
			cancelTimeout() // Release timeout context

			if err != nil {
				app.logger.Warn("Failed to get local geth height, retrying...", "error", err)
				select {
				case <-app.appCtx.Done():
					return
				case <-time.After(initRetryInterval):
					continue
				}
			}

			// Set lastProcessedHeight to current local height
			// The processing loop will request localHeight + 1
			app.processingMutex.Lock()
			app.lastProcessedHeight = localHeight
			app.isInitialized = true
			app.processingMutex.Unlock()

			app.logger.Info(
				"Successfully detected starting height from local geth",
				"local_height", localHeight,
				"will_start_from", localHeight+1,
			)
			return // Initialization successful
		}
	}
}

// Start begins the member node operation
func (app *MemberNodeApp) Start() {
	app.logger.Info("Starting MemberNodeApp...")

	// Launch health server
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/health", app.healthHandler)
		addr := app.cfg.HealthAddr
		server := &http.Server{Addr: addr, Handler: mux}
		app.logger.Info("Health endpoint listening", "address", addr)

		go func() {
			<-app.appCtx.Done()
			ctx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancelShutdown()
			if err := server.Shutdown(ctx); err != nil {
				// ErrServerClosed is expected on graceful shutdown,
				// context.DeadlineExceeded if shutdownTimeout is reached.
				if !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, context.DeadlineExceeded) {
					app.logger.Warn("Health server shutdown error", "error", err)
				}
			}
		}()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Health server error", "error", err)
			app.cancel() // Trigger app shutdown if health server fails critically
		}
	}()

	// Initialize starting height from local geth
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.initializeStartingHeight()
	}()

	// Start sequential payload processing loop
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer app.logger.Info("MemberNodeApp run loop finished.")
		app.runSequentialLoop()
	}()
}

// healthHandler responds on /health
func (app *MemberNodeApp) healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.appCtx.Err(); err != nil {
		http.Error(w, "unavailable (shutting down)", http.StatusServiceUnavailable)
		return
	}

	app.connectionStatus.Lock()
	leaderAvailable := app.leaderAvailable
	app.connectionStatus.Unlock()

	if !leaderAvailable {
		app.logger.Warn("Health check failed: leader node is not available")
		http.Error(w, "leader node is not available", http.StatusServiceUnavailable)
		return
	}

	// Optionally, check if initialized
	app.processingMutex.RLock()
	initialized := app.isInitialized
	app.processingMutex.RUnlock()
	if !initialized {
		http.Error(w, "initializing", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		app.logger.Error("Failed to write health response", "error", err)
		// No return here, header already sent
	}
}

// runSequentialLoop continuously processes payloads in sequential order
func (app *MemberNodeApp) runSequentialLoop() {
	// Wait for initialization
	select {
	case <-app.appCtx.Done():
		app.logger.Info("Run loop stopping before initialization due to context cancellation.")
		return
	case <-app.initializedCh:
		app.logger.Info("Initialization complete, starting main processing.")
	}

	// Check if initialization actually completed successfully (isInitialized flag)
	// This is a safeguard in case initializedCh was closed due to appCtx.Done() during init.
	app.processingMutex.RLock()
	if !app.isInitialized {
		app.processingMutex.RUnlock()
		app.logger.Error("Initialization failed or was cancelled, run loop cannot start.")
		return
	}
	startingHeight := app.lastProcessedHeight + 1
	app.processingMutex.RUnlock()

	app.logger.Info(
		"MemberNodeApp sequential run loop started",
		"instanceID", app.cfg.InstanceID,
		"starting_height", startingHeight,
	)

	consecutiveErrors := 0
	ticker := time.NewTicker(app.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-app.appCtx.Done():
			app.logger.Info("MemberNodeApp run loop stopping due to context cancellation.")
			return
		case <-ticker.C:
			err := app.processSequentialPayloads()

			if err != nil {
				consecutiveErrors++
				app.setLeaderAvailability(false)

				if consecutiveErrors >= maxConsecutiveErrors {
					app.logger.Error(
						"Too many consecutive errors, member node may be unstable. Check leader connection and local Geth.",
						"error", err,
						"consecutive_errors", consecutiveErrors,
					)
				} else {
					app.logger.Warn(
						"Failed to process sequential payloads",
						"error", err,
						"consecutive_errors", consecutiveErrors,
					)
				}
			} else {
				if consecutiveErrors > 0 {
					app.logger.Info(
						"Recovered from errors",
						"previous_consecutive_errors", consecutiveErrors,
					)
				}
				consecutiveErrors = 0
				app.setLeaderAvailability(true) // Assuming success implies leader is available
			}
		}
	}
}

// processSequentialPayloads fetches and processes payloads in sequential order
func (app *MemberNodeApp) processSequentialPayloads() error {
	ctx, cancel := context.WithTimeout(app.appCtx, apiCallTimeout) // Overall timeout for this processing cycle
	defer cancel()

	// Check leader health first
	if err := app.payloadClient.CheckHealth(ctx); err != nil {
		return fmt.Errorf("leader health check failed: %w", err)
	}
	app.setLeaderAvailability(true) // Leader is reachable

	app.processingMutex.RLock()
	lastProcessedHeight := app.lastProcessedHeight
	isCatchingUp := app.isCatchingUp
	app.processingMutex.RUnlock()

	// Determine how many payloads to request
	var requestLimit int
	if isCatchingUp {
		requestLimit = maxCatchupPayloads // e.g., 100
		app.logger.Debug("In catch-up mode, requesting more payloads", "limit", requestLimit)
	} else {
		requestLimit = batchSize // e.g., 10
		app.logger.Debug("In normal mode, requesting standard batch of payloads", "limit", requestLimit)
	}

	// Get payloads since our last processed height
	nextHeightToRequest := lastProcessedHeight + 1
	payloadsResponse, err := app.payloadClient.GetPayloadsSince(ctx, nextHeightToRequest, requestLimit)
	if err != nil {
		return fmt.Errorf("failed to get payloads since height %d: %w", nextHeightToRequest, err)
	}

	if len(payloadsResponse.Payloads) == 0 {
		app.logger.Debug("No new payloads available", "waiting_for_height", nextHeightToRequest)
		return nil
	}

	// Update catch-up mode status
	currentlyCatchingUp := isCatchingUp
	if !currentlyCatchingUp && len(payloadsResponse.Payloads) >= batchSize {
		app.processingMutex.Lock()
		app.isCatchingUp = true
		app.processingMutex.Unlock()
		app.logger.Info(
			"Entering catch-up mode",
			"current_height", lastProcessedHeight,
			"available_payloads", len(payloadsResponse.Payloads),
		)
	} else if currentlyCatchingUp && len(payloadsResponse.Payloads) < catchUpExitThreshold {
		app.processingMutex.Lock()
		app.isCatchingUp = false
		app.processingMutex.Unlock()
		app.logger.Info(
			"Exiting catch-up mode",
			"current_height", lastProcessedHeight,
			"available_payloads", len(payloadsResponse.Payloads),
		)
	}

	// Process payloads sequentially
	processedCount := 0
	for _, payload := range payloadsResponse.Payloads {
		select {
		case <-app.appCtx.Done(): // Check for shutdown signal before processing each payload
			return nil
		default:
		}

		// Get the most up-to-date lastProcessedHeight for sequence check
		app.processingMutex.RLock()
		currentSystemHeight := app.lastProcessedHeight
		app.processingMutex.RUnlock()

		expectedHeightForThisPayload := currentSystemHeight + 1

		// Case 1: Gap detected (payload is for a future height)
		if payload.BlockHeight > expectedHeightForThisPayload {
			app.logger.Warn(
				"Gap detected in payload sequence, attempting to fill",
				"expected_height", expectedHeightForThisPayload,
				"received_payload_height", payload.BlockHeight,
				"gap_size", payload.BlockHeight-expectedHeightForThisPayload,
			)
			// Try to fill the gap from expectedHeightForThisPayload up to payload.BlockHeight - 1
			if err := app.fillPayloadGap(ctx, expectedHeightForThisPayload, payload.BlockHeight-1); err != nil {
				return fmt.Errorf("failed to fill payload gap from %d to %d: %w",
					expectedHeightForThisPayload, payload.BlockHeight-1, err)
			}
			// After gap fill, lastProcessedHeight should be (payload.BlockHeight - 1)
			// Re-fetch currentSystemHeight to ensure the next check is correct
			app.processingMutex.RLock()
			currentSystemHeight = app.lastProcessedHeight
			app.processingMutex.RUnlock()
			expectedHeightForThisPayload = currentSystemHeight + 1
		}

		// Case 2: Payload is for an already processed or an older, unexpected height
		if payload.BlockHeight < expectedHeightForThisPayload {
			app.logger.Debug(
				"Skipping already processed or out-of-order (older) payload",
				"payload_height", payload.BlockHeight,
				"expected_at_least", expectedHeightForThisPayload,
				"current_system_height", currentSystemHeight,
			)
			continue
		}

		// Case 3: Payload is for the expected next height (critical check)
		if payload.BlockHeight != expectedHeightForThisPayload {
			// This should ideally not be reached if gap filling and previous checks are correct
			return fmt.Errorf("critical sequence error: payload height %d does not match expected next height %d after potential gap fill. Current system height: %d",
				payload.BlockHeight, expectedHeightForThisPayload, currentSystemHeight)
		}

		// Process the payload
		if err := app.processPayload(ctx, &payload); err != nil {
			return fmt.Errorf("failed to process payload at height %d: %w", payload.BlockHeight, err)
		}

		// Update processed height (this is critical)
		app.processingMutex.Lock()
		app.lastProcessedHeight = payload.BlockHeight
		app.processingMutex.Unlock()

		processedCount++

		// In catch-up mode, limit processing per cycle to avoid holding locks for too long
		// or starving other operations, and to allow context cancellation checks.
		app.processingMutex.RLock()
		stillCatchingUp := app.isCatchingUp
		currentHeightAfterProcess := app.lastProcessedHeight
		app.processingMutex.RUnlock()

		if stillCatchingUp && processedCount >= maxCatchupPayloads {
			app.logger.Info(
				"Processed maximum catch-up payloads in this cycle, will continue in next cycle",
				"processed_count", processedCount,
				"current_height", currentHeightAfterProcess,
			)
			break // Exit the loop for this batch, will fetch new batch in next ticker
		}
	}

	if processedCount > 0 {
		app.processingMutex.RLock()
		finalHeight := app.lastProcessedHeight
		catchUpMode := app.isCatchingUp
		app.processingMutex.RUnlock()

		app.logger.Info(
			"Successfully processed sequential payloads batch",
			"processed_count", processedCount,
			"final_height", finalHeight,
			"catch_up_mode", catchUpMode,
		)
	}
	return nil
}

// fillPayloadGap attempts to fetch and process missing payloads in a range
func (app *MemberNodeApp) fillPayloadGap(ctx context.Context, startHeight, endHeight uint64) error {
	if startHeight > endHeight {
		app.logger.Info(
			"No gap to fill or invalid range",
			"start", startHeight, "end", endHeight,
		)
		return nil
	}
	app.logger.Info(
		"Filling payload gap",
		"start_height", startHeight,
		"end_height", endHeight,
		"gap_size", endHeight-startHeight+1,
	)

	for height := startHeight; height <= endHeight; height++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during gap fill at height %d: %w", height, ctx.Err())
		case <-app.appCtx.Done():
			return fmt.Errorf("application shutting down during gap fill at height %d: %w", height, app.appCtx.Err())
		default:
		}

		// Get specific payload by height
		payload, err := app.payloadClient.GetPayloadByHeight(ctx, height)
		if err != nil {
			return fmt.Errorf("failed to get payload for gap at height %d: %w", height, err)
		}

		// Process the payload
		if err := app.processPayload(ctx, payload); err != nil {
			return fmt.Errorf("failed to process gap payload at height %d: %w", height, err)
		}

		app.processingMutex.Lock()
		if app.lastProcessedHeight != height-1 {
			app.processingMutex.Unlock()
			// This indicates a severe internal inconsistency or a concurrent modification problem.
			// It means another part of the code or a previous iteration did not leave lastProcessedHeight as expected.
			return fmt.Errorf("critical sequence error during gap fill: expected lastProcessedHeight %d before processing %d, but got %d",
				height-1, height, app.lastProcessedHeight)
		}
		app.lastProcessedHeight = height
		app.processingMutex.Unlock()

		app.logger.Debug("Filled gap payload", "height", height)
	}

	app.logger.Info(
		"Successfully filled payload gap",
		"start_height", startHeight,
		"end_height", endHeight,
		"final_processed_height_after_gap_fill", endHeight,
	)
	return nil
}

// processPayload applies a single payload to the local geth client
func (app *MemberNodeApp) processPayload(ctx context.Context, payload *api.PayloadResponse) error {
	app.logger.Info(
		"Processing payload",
		"payload_id", payload.PayloadID,
		"block_height", payload.BlockHeight,
	)

	// Apply payload to local geth client
	err := app.blockBuilder.FinalizeBlock(ctx, payload.PayloadID, payload.ExecutionPayload, "")
	if err != nil {
		app.logger.Error(
			"Failed to finalize block",
			"payload_id", payload.PayloadID,
			"block_height", payload.BlockHeight,
			"error", err,
		)
		return fmt.Errorf("blockBuilder.FinalizeBlock failed for height %d: %w", payload.BlockHeight, err)
	}

	app.logger.Info(
		"Successfully applied payload",
		"payload_id", payload.PayloadID,
		"block_height", payload.BlockHeight,
	)
	return nil
}

// setLeaderAvailability updates the leader availability status
func (app *MemberNodeApp) setLeaderAvailability(available bool) {
	app.connectionStatus.Lock()
	defer app.connectionStatus.Unlock()

	if app.leaderAvailable != available {
		app.leaderAvailable = available
		app.logger.Info("Leader availability changed", "available", available)
	}
}

// GetLastProcessedHeight returns the last successfully processed block height
func (app *MemberNodeApp) GetLastProcessedHeight() uint64 {
	app.processingMutex.RLock()
	defer app.processingMutex.RUnlock()
	return app.lastProcessedHeight
}

// Stop gracefully stops the member node
func (app *MemberNodeApp) Stop() {
	app.logger.Info("Stopping MemberNodeApp...")
	app.cancel() // Signal all goroutines to stop

	waitCh := make(chan struct{})
	go func() {
		app.wg.Wait() // Wait for all primary goroutines to finish
		close(waitCh)
	}()

	select {
	case <-waitCh:
		app.logger.Info("MemberNodeApp goroutines shut down gracefully.")
	case <-time.After(shutdownTimeout + 1*time.Second):
		app.logger.Warn("MemberNodeApp shutdown timed out waiting for goroutines.")
	}

	app.processingMutex.RLock()
	finalHeight := app.lastProcessedHeight
	app.processingMutex.RUnlock()

	app.logger.Info("MemberNodeApp stopped.", "final_processed_height", finalHeight)
}
