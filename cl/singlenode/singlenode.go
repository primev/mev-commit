package singlenode

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/ethclient"
	localstate "github.com/primev/mev-commit/cl/singlenode/state"
)

const (
	// Stop Function
	shutdownTimeout = 5 * time.Second
)

// Config holds the configuration for the SingleNodeApp.
type Config struct {
	InstanceID               string
	EthClientURL             string
	JWTSecret                string
	EVMBuildDelay            time.Duration
	EVMBuildDelayEmptyBlocks time.Duration
	PriorityFeeReceipt       string
}

// SingleNodeApp orchestrates block production for a single node.
type SingleNodeApp struct {
	logger       *slog.Logger
	cfg          Config
	blockBuilder *blockbuilder.BlockBuilder
	stateManager *localstate.LocalStateManager
	appCtx       context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewSingleNodeApp creates and initializes a new SingleNodeApp.
func NewSingleNodeApp(
	parentCtx context.Context,
	cfg Config,
	logger *slog.Logger,
) (*SingleNodeApp, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	jwtBytes, err := hex.DecodeString(cfg.JWTSecret)
	if err != nil {
		cancel() // Cancel the derived context
		logger.Error(
			"failed to decode JWT secret",
			"error", err,
		)
		return nil, err
	}

	engineCL, err := ethclient.NewAuthClient(ctx, cfg.EthClientURL, jwtBytes)
	if err != nil {
		cancel() // Cancel the derived context
		logger.Error(
			"failed to create Ethereum engine client",
			"error", err,
		)
		return nil, err
	}

	stateMgr := localstate.NewLocalStateManager(logger.With("component", "LocalStateManager"))
	bb := blockbuilder.NewBlockBuilder(
		stateMgr,
		engineCL,
		logger.With("component", "BlockBuilder"),
		cfg.EVMBuildDelay,
		cfg.EVMBuildDelayEmptyBlocks,
		cfg.PriorityFeeReceipt,
	)

	return &SingleNodeApp{
		logger:       logger,
		cfg:          cfg,
		blockBuilder: bb,
		stateManager: stateMgr,
		appCtx:       ctx,
		cancel:       cancel,
	}, nil
}

// Start begins the main block production loop.
func (app *SingleNodeApp) Start() {
	app.logger.Info("Starting SingleNodeApp...")
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer app.logger.Info("SingleNodeApp run loop finished.")
		app.runLoop()
	}()
}

// shutdownWithError handles errors during the run loop and initiates a shutdown.
func (app *SingleNodeApp) shutdownWithError(err error, message string, args ...any) {
	// slog handles key-value pairs directly
	logArgs := append(args, "error", err)
	app.logger.Error(message, logArgs...)
	app.cancel()
}

// resetBlockProduction clears state and prepares for a new block production cycle.
// It returns true if a shutdown is initiated due to a reset failure.
func (app *SingleNodeApp) resetBlockProduction(logMessage string, logArgs ...interface{}) (shutdownInitiated bool) {
	app.logger.Info(logMessage, logArgs...)
	if err := app.stateManager.ResetBlockState(app.appCtx); err != nil {
		app.shutdownWithError(err, "Failed to reset block state during run loop operations")
		return true
	}
	return false
}

func (app *SingleNodeApp) runLoop() {
	app.logger.Info("SingleNodeApp run loop started", "instanceID", app.cfg.InstanceID)

	// Make sure we're starting with a clean state
	if app.resetBlockProduction("Initializing block production state") {
		return // Shutdown initiated by resetBlockProduction
	}

	for {
		select {
		case <-app.appCtx.Done():
			app.logger.Info("SingleNodeApp run loop stopping due to context cancellation.")
			return
		default:
			// Directly run the block production cycle without steps
			if err := app.produceBlock(); err != nil {
				if errors.Is(err, blockbuilder.ErrEmptyBlock) {
					// Handle empty block error
					app.logger.Info("empty block produced, waiting for new payload")
					continue
				}
				// Handle errors but continue the loop
				app.logger.Error(
					"block production cycle failed",
					"error", err,
				)
			}
			// Successful block production, reset for the next block
			if app.resetBlockProduction("Block production successful. Resetting state for next block.") {
				// 0 chance to happen, if in-memory store is used
				return // Shutdown initiated by resetBlockProduction
			}

		}
	}
}

// produceBlock handles the entire block production cycle in a direct, procedural manner
func (app *SingleNodeApp) produceBlock() error {
	// Step 1: Build the block
	if err := app.blockBuilder.GetPayload(app.appCtx); err != nil {
		return fmt.Errorf("failed to get payload: %w", err)
	}

	// Retrieve the current state after payload creation
	currentState := app.stateManager.GetBlockBuildState(app.appCtx)
	if currentState.PayloadID == "" {
		return errors.New("payload ID is empty after GetPayload call")
	}

	// Step 2: Finalize the block
	app.logger.Info("Finalizing block", "payload_id", currentState.PayloadID)
	if err := app.blockBuilder.FinalizeBlock(app.appCtx, currentState.PayloadID, currentState.ExecutionPayload, ""); err != nil {
		return fmt.Errorf("failed to finalize block: %w", err)
	}

	return nil
}

// Stop signals the application to shut down and waits for goroutines to finish.
func (app *SingleNodeApp) Stop() {
	app.logger.Info("stopping SingleNodeApp...")
	app.cancel()

	waitCh := make(chan struct{})
	go func() {
		app.wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		app.logger.Info("SingleNodeApp run loop shut down gracefully.")
	case <-time.After(shutdownTimeout):
		app.logger.Warn("SingleNodeApp shutdown timed out waiting for run loop.")
	}
	app.logger.Info("SingleNodeApp stopped.")
}
