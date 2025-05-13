package singlenode

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/singlenode/state"
	"github.com/primev/mev-commit/cl/types"
	"github.com/primev/mev-commit/cl/util"
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
	engineClient blockbuilder.EngineClient // Keep a reference if needed for direct calls, though BB handles most
	appCtx       context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewSingleNodeApp creates and initializes a new SingleNodeApp.
func NewSingleNodeApp(
	appCtx context.Context, // Parent context
	cfg Config,
	logger *slog.Logger,
) (*SingleNodeApp, error) {
	ctx, cancel := context.WithCancel(appCtx)

	jwtBytes, err := hex.DecodeString(cfg.JWTSecret)
	if err != nil {
		cancel()
		logger.Error("Failed to decode JWT secret", "error", err)
		return nil, err
	}

	engineCL, err := ethclient.NewAuthClient(ctx, cfg.EthClientURL, jwtBytes)
	if err != nil {
		cancel()
		logger.Error("Failed to create Ethereum engine client", "error", err)
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
		engineClient: engineCL, // Stored if needed
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

func (app *SingleNodeApp) runLoop() {
	app.logger.Info("SingleNodeApp run loop started", "instanceID", app.cfg.InstanceID)

	// On startup, process any pending block from a previous (crashed) session.
	// With LocalStateManager (in-memory), this is mostly for conceptual completeness
	// unless persistence is added to LocalStateManager.
	// if err := app.blockBuilder.ProcessLastPayload(app.appCtx); err != nil {
	// 	// If ProcessLastPayload fails (e.g., after retries on FinalizeBlock),
	// 	// it might indicate a persistent issue with Geth.
	// 	if errors.Is(err, util.ErrFailedAfterNAttempts) || errors.Is(err, context.Canceled) {
	// 		app.logger.Error("Critical error processing last payload at startup. Shutting down.", "error", err)
	// 		app.cancel() // Trigger shutdown for the rest of the app
	// 		return
	// 	}
	// 	app.logger.Warn("Non-critical error processing last payload at startup. State might be reset.", "error", err)
	// 	// Attempt to reset state to ensure clean start if ProcessLastPayload had an issue but wasn't critical.
	// 	if resetErr := app.stateManager.ResetBlockState(app.appCtx); resetErr != nil {
	// 		app.logger.Error("Failed to reset state after ProcessLastPayload error. Shutting down.", "error", resetErr)
	// 		app.cancel()
	// 		return
	// 	}
	// }

	// Main production loop
	for {
		select {
		case <-app.appCtx.Done():
			app.logger.Info("SingleNodeApp run loop stopping due to context cancellation.")
			return
		default:
			// Determine current step from state
			currentState := app.stateManager.GetBlockBuildState(app.appCtx)
			var err error

			switch currentState.CurrentStep {
			case types.StepBuildBlock:
				app.logger.Info("RunLoop: StepBuildBlock")
				err = app.blockBuilder.GetPayload(app.appCtx)
				if err != nil {
					app.logger.Error("RunLoop: GetPayload failed", "error", err)
					// If GetPayload fails, state remains StepBuildBlock.
					// Retry is handled within GetPayload. If it returns ErrFailedAfterNAttempts,
					// it's a critical error.
					if errors.Is(err, util.ErrFailedAfterNAttempts) {
						app.logger.Error("RunLoop: GetPayload failed critically after retries. Shutting down.", "error", err)
						app.cancel()
						return
					}
					// For other errors, or if GetPayload skipped an empty block (err == nil but state not changed),
					// the loop will retry. Add a small delay to prevent tight spin on persistent non-critical errors.
					time.Sleep(500 * time.Millisecond)
				}
				// If GetPayload was successful and built a block, stateManager would have updated
				// CurrentStep to StepFinalizeBlock. If it skipped an empty block, CurrentStep remains StepBuildBlock.

			case types.StepFinalizeBlock:
				app.logger.Info("RunLoop: StepFinalizeBlock", "payload_id", currentState.PayloadID)
				err = app.blockBuilder.FinalizeBlock(app.appCtx, currentState.PayloadID, currentState.ExecutionPayload, "")
				if err != nil {
					app.logger.Error("RunLoop: FinalizeBlock failed", "error", err)
					// If FinalizeBlock fails, state remains StepFinalizeBlock.
					// Retry is handled within FinalizeBlock. If it returns ErrFailedAfterNAttempts, critical.
					if errors.Is(err, util.ErrFailedAfterNAttempts) {
						app.logger.Error("RunLoop: FinalizeBlock failed critically after retries. Shutting down.", "error", err)
						app.cancel()
						return
					}
					// If error is due to payload validation (permanent), reset state to avoid retrying same bad payload.
					var bErr *backoff.PermanentError
					if errors.As(err, &bErr) && strings.Contains(bErr.Err.Error(), "execution payload validation failed") {
						app.logger.Error("RunLoop: FinalizeBlock failed due to permanent payload validation error. Resetting state.", "original_error", bErr.Err)
						if resetErr := app.stateManager.ResetBlockState(app.appCtx); resetErr != nil {
							app.logger.Error("RunLoop: Failed to reset state after permanent FinalizeBlock error. Shutting down.", "error", resetErr)
							app.cancel()
						}
						return // return from switch case, not entire loop, to re-evaluate state
					}
					// For other errors, loop will retry. Add delay.
					time.Sleep(500 * time.Millisecond)
				} else {
					// FinalizeBlock successful, reset state for the next block.
					app.logger.Info("RunLoop: FinalizeBlock successful. Resetting state for next block.")
					if resetErr := app.stateManager.ResetBlockState(app.appCtx); resetErr != nil {
						app.logger.Error("RunLoop: Failed to reset state after successful FinalizeBlock. Shutting down.", "error", resetErr)
						app.cancel()
						return
					}
				}

			default:
				app.logger.Warn("RunLoop: Unknown current step in state", "step", currentState.CurrentStep.String())
				if resetErr := app.stateManager.ResetBlockState(app.appCtx); resetErr != nil {
					app.logger.Error("RunLoop: Failed to reset state from unknown step. Shutting down.", "error", resetErr)
					app.cancel()
					return
				}
				time.Sleep(1 * time.Second) // Pause if in unknown state before retrying
			}
			// A short general delay to prevent extremely tight loops if conditions lead to no state change
			// For example, if GetPayload consistently skips empty blocks very fast.
			// This can be tied to blockBuilder.GetBuildDelay() or a fixed minimum.
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Stop signals the application to shut down and waits for goroutines to finish.
func (app *SingleNodeApp) Stop() {
	app.logger.Info("Stopping SingleNodeApp...")
	app.cancel() // Signal all operations using app.appCtx to stop

	// Wait for the main run loop to finish
	// Set a timeout for waiting to prevent indefinite blocking.
	waitCh := make(chan struct{})
	go func() {
		app.wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		app.logger.Info("SingleNodeApp run loop shut down gracefully.")
	case <-time.After(5 * time.Second): // Timeout for shutdown
		app.logger.Warn("SingleNodeApp shutdown timed out waiting for run loop.")
	}
	app.logger.Info("SingleNodeApp stopped.")
}
