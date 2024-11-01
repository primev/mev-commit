package redisapp

import (
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/redisapp/blockbuilder"
	"github.com/primev/mev-commit/cl/redisapp/leaderelection"
	"github.com/primev/mev-commit/cl/redisapp/state"
	"github.com/redis/go-redis/v9"
)

type MevCommitChain struct {
	logger           *slog.Logger

	cancel context.CancelFunc

	// Managers and components
	stateManager          state.StateManager
	blockBuilder          *blockbuilder.BlockBuilder
	leaderElectionHandler *leaderelection.LeaderElectionHandler
}

func NewMevCommitChain(instanceID, ecURL, jwtSecret, genesisBlockHash string, logger *slog.Logger, redisAddr string, buildDelay time.Duration) (*MevCommitChain, error) {
	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// JWT secret decoding
	bytes, err := hex.DecodeString(jwtSecret)
	if err != nil {
		cancel()
		logger.Error("Error decoding JWT secret", "error", err)
		return nil, err
	}

	engineCL, err := ethclient.NewAuthClient(ctx, ecURL, bytes)
	if err != nil {
		cancel()
		logger.Error("Error creating engine client", "error", err)
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	err = redisClient.ConfigSet(ctx, "min-replicas-to-write", "1").Err()
	if err != nil {
		cancel()
		logger.Error("Error setting min-replicas-to-write", "error", err)
		return nil, err
	}

	stateManager := state.NewRedisStateManager(instanceID, redisClient, logger, genesisBlockHash)

	blockBuilder := blockbuilder.NewBlockBuilder(stateManager, engineCL, logger, buildDelay)

	leaderElectionHandler := leaderelection.NewLeaderElectionHandler(
		instanceID,
		logger,
		redisClient,
		stateManager,
		blockBuilder,
	)

	app := &MevCommitChain{
		stateManager:          stateManager,
		blockBuilder:          blockBuilder,
		logger:                logger,
		cancel:                cancel,
		leaderElectionHandler: leaderElectionHandler,
	}

	logger.Info("MevCommitChain initialized")

	err = app.stateManager.LoadOrInitializeBlockState(ctx)
	if err != nil {
		cancel()
		logger.Error("Failed to load or initialize build state", "error", err)
		return nil, err
	}

	// Start leader election handling
	app.leaderElectionHandler.HandleLeadershipEvents()

	return app, nil
}

func (app *MevCommitChain) Stop() {
	// Cancel the context to signal all goroutines to stop
	app.cancel()
	app.leaderElectionHandler.Stop()
	// Wait for all goroutines to finish
	app.logger.Info("Waiting for goroutines to finish")
	app.stateManager.Stop()
	app.logger.Info("MevCommitChain stopped gracefully")
}
