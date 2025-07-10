package singlenode

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	gethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/singlenode/api"
	"github.com/primev/mev-commit/cl/singlenode/payloadstore"
	localstate "github.com/primev/mev-commit/cl/singlenode/state"
	"github.com/primev/mev-commit/cl/types"
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
	HealthAddr               string
	PostgresDSN              string
	APIAddr                  string
	NonAuthEthClientURL      string
	TxPoolPollingInterval    time.Duration
}

type BlockBuilder interface {
	GetPayload(ctx context.Context) error
	FinalizeBlock(ctx context.Context, payloadID string, executionPayload string, extraData string) error
	GetExecutionHead() *types.ExecutionHead
}

// SingleNodeApp orchestrates block production for a single node.
type SingleNodeApp struct {
	logger       *slog.Logger
	cfg          Config
	blockBuilder BlockBuilder
	// stateManager is a local state manager for block production
	// it's not anticipated to use DB as all the state already in geth client
	stateManager      *localstate.LocalStateManager
	payloadRepo       types.PayloadRepository
	payloadServer     *api.PayloadServer
	appCtx            context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	connectionStatus  sync.Mutex
	connectionRefused bool
	ethClient         *gethclient.Client
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
		cancel()
		logger.Error(
			"failed to decode JWT secret",
			"error", err,
		)
		return nil, err
	}

	engineCL, err := ethclient.NewAuthClient(ctx, cfg.EthClientURL, jwtBytes)
	if err != nil {
		cancel()
		logger.Error(
			"failed to create Ethereum engine client",
			"error", err,
		)
		return nil, err
	}

	gethClient, err := gethclient.DialContext(ctx, cfg.NonAuthEthClientURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	stateMgr := localstate.NewLocalStateManager(logger.With("component", "LocalStateManager"))
	bb := blockbuilder.NewBlockBuilder(
		stateMgr,
		engineCL,
		logger.With("component", "BlockBuilder"),
		cfg.EVMBuildDelay,
		cfg.EVMBuildDelayEmptyBlocks,
		cfg.PriorityFeeReceipt,
		gethClient,
	)

	var pRepo types.PayloadRepository
	if cfg.PostgresDSN != "" {
		repo, err := payloadstore.NewPostgresRepository(ctx, cfg.PostgresDSN, logger)
		if err != nil {
			cancel()
			logger.Error(
				"failed to create payload repository",
				"error", err,
			)
			return nil, fmt.Errorf("failed to initialize payload repository: %w", err)
		}
		pRepo = repo
		logger.Info("Payload repository initialized, payloads will be saved to PostgreSQL.")
	} else {
		logger.Info("PostgresDSN not provided, payload saving to DB is disabled.")
	}

	var payloadServer *api.PayloadServer
	if cfg.APIAddr != "" {
		payloadServer = api.NewPayloadServer(
			cfg.APIAddr,
			stateMgr,
			pRepo,
			logger.With("component", "APIServer"),
		)
		logger.Info("API server initialized for member nodes", "addr", cfg.APIAddr)
	} else {
		logger.Info("API address not provided, member node API is disabled.")
	}

	return &SingleNodeApp{
		logger:            logger,
		cfg:               cfg,
		blockBuilder:      bb,
		stateManager:      stateMgr,
		payloadRepo:       pRepo,
		payloadServer:     payloadServer,
		appCtx:            ctx,
		cancel:            cancel,
		connectionRefused: false,
		ethClient:         gethClient,
	}, nil
}

// isConnectionRefused checks if the error is a connection refused error
func isConnectionRefused(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}

// setConnectionStatus updates the connection status based on the error
func (app *SingleNodeApp) setConnectionStatus(err error) {
	app.connectionStatus.Lock()
	defer app.connectionStatus.Unlock()

	if err == nil {
		// Reset the connection refused flag if the operation was successful
		app.connectionRefused = false
		return
	}

	// Check if the error indicates a connection refused
	if isConnectionRefused(err) {
		app.connectionRefused = true
		app.logger.Warn(
			"Connection refused detected, Ethereum might be unavailable",
			"error", err,
		)
	}
}

// healthHandler responds on /health with 200 OK if the app context is active and no connection refusal, or 503 otherwise.
func (app *SingleNodeApp) healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.appCtx.Err(); err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}

	// Check connection status
	app.connectionStatus.Lock()
	connectionRefused := app.connectionRefused
	app.connectionStatus.Unlock()

	if connectionRefused {
		app.logger.Warn("Health check failed: ethereum is not available (connection refused)")
		http.Error(w, "ethereum is not available", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Start begins the main block production loop and launches the health endpoint.
func (app *SingleNodeApp) Start() {
	app.logger.Info("Starting SingleNodeApp...")

	// Launch health server
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/health", app.healthHandler)
		addr := app.cfg.HealthAddr
		server := &http.Server{Addr: addr, Handler: mux}
		app.logger.Info("Health endpoint listening", "address", addr)

		// Shutdown server when app context is done
		go func() {
			<-app.appCtx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Health server error", "error", err)
		}
	}()

	if app.payloadServer != nil {
		app.wg.Add(1)
		go func() {
			defer app.wg.Done()
			if err := app.payloadServer.Start(app.appCtx); err != nil {
				app.logger.Error("API server error", "error", err)
			}
		}()
	}

	// Start block production loop
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer app.logger.Info("SingleNodeApp run loop finished.")
		app.runLoop()
	}()
}

// shutdownWithError handles errors during the run loop and initiates a shutdown.
func (app *SingleNodeApp) shutdownWithError(err error, message string, args ...any) {
	logArgs := append(args, "error", err)
	app.logger.Error(message, logArgs...)
	app.cancel()
}

// resetBlockProduction clears state and prepares for a new block production cycle.
func (app *SingleNodeApp) resetBlockProduction(logMessage string, logArgs ...interface{}) bool {
	app.logger.Info(logMessage, logArgs...)
	if err := app.stateManager.ResetBlockState(app.appCtx); err != nil {
		app.shutdownWithError(err, "Failed to reset block state during run loop operations")
		return true
	}
	return false
}

func (app *SingleNodeApp) runLoop() {
	app.logger.Info("SingleNodeApp run loop started", "instanceID", app.cfg.InstanceID)
	if app.resetBlockProduction("Initializing block production state") {
		return
	}

	for {
		select {
		case <-app.appCtx.Done():
			app.logger.Info("SingleNodeApp run loop stopping due to context cancellation.")
			return
		default:
			err := app.produceBlock()
			// Update connection status based on the error
			app.setConnectionStatus(err)

			if err != nil {
				if errors.Is(err, blockbuilder.ErrEmptyBlock) {
					app.logger.Debug("no pending transactions, will try again in: %s", "timeout", app.cfg.TxPoolPollingInterval)
					time.Sleep(app.cfg.TxPoolPollingInterval)
					continue
				} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					app.logger.Info("context canceled or deadline exceeded, stopping block production")
					return
				}
				app.logger.Error(
					"block production cycle failed",
					"error", err,
				)
			}
			if app.resetBlockProduction("Block production successful. Resetting state for next block.") {
				// if state is local, it couldn't happen
				return
			}
		}
	}
}

// produceBlock handles the entire block production cycle.
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

	// Get current block height from the execution head
	executionHead := app.blockBuilder.GetExecutionHead()
	var blockHeight uint64
	if executionHead != nil {
		blockHeight = executionHead.BlockHeight + 1 // Next block height
	} else {
		app.logger.Warn("No execution head available, using height 0")
		blockHeight = 0
	}

	if app.payloadRepo != nil {
		payloadInfo := &types.PayloadInfo{
			PayloadID:        currentState.PayloadID,
			ExecutionPayload: currentState.ExecutionPayload,
			BlockHeight:      blockHeight,
		}
		saveCtx, saveCancel := context.WithTimeout(app.appCtx, 200*time.Millisecond)
		defer saveCancel()

		if err := app.payloadRepo.SavePayload(saveCtx, payloadInfo); err != nil {
			app.logger.Error(
				"Failed to save payload to database",
				"payload_id", currentState.PayloadID,
				"error", err,
			)
			return fmt.Errorf("failed to save payload to database: %w", err)
		} else {
			app.logger.Info("Payload details submitted to database for saving", "payload_id", currentState.PayloadID)
		}
	}

	// Step 2: Finalize the block
	app.logger.Info(
		"finalizing block",
		"payload_id", currentState.PayloadID,
		"block_height", blockHeight,
	)
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

	if app.payloadRepo != nil {
		if err := app.payloadRepo.Close(); err != nil {
			app.logger.Error("Error closing payload repository", "error", err)
		} else {
			app.logger.Info("Payload repository closed.")
		}
	}

	if app.ethClient != nil {
		app.ethClient.Close()
		app.logger.Info("Ethereum client closed.")
	}

	app.logger.Info("SingleNodeApp stopped.")
}
