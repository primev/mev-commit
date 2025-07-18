package redisapp

import (
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/redisapp/leaderfollower"
	"github.com/primev/mev-commit/cl/redisapp/state"
	"github.com/redis/go-redis/v9"
)

type MevCommitChain struct {
	logger *slog.Logger

	cancel context.CancelFunc

	// Managers and components
	stateManager state.StateManager
	blockBuilder *blockbuilder.BlockBuilder
	lfm          *leaderfollower.LeaderFollowerManager
}

func NewMevCommitChain(
	instanceID, ecURL, jwtSecret, redisAddr, feeReceipt string,
	logger *slog.Logger,
	buildDelay, buildDelayEmptyBlocks time.Duration,
) (*MevCommitChain, error) {
	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// JWT secret decoding
	bytes, err := hex.DecodeString(jwtSecret)
	if err != nil {
		cancel()
		logger.Error(
			"Error decoding JWT secret",
			"error", err,
		)
		return nil, err
	}

	engineCL, err := ethclient.NewAuthClient(ctx, ecURL, bytes)
	if err != nil {
		cancel()
		logger.Error(
			"Error creating engine client",
			"error", err,
		)
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	err = redisClient.ConfigSet(ctx, "min-replicas-to-write", "1").Err()
	if err != nil {
		cancel()
		logger.Error(
			"Error setting min-replicas-to-write",
			"error", err,
		)
		return nil, err
	}

	coordinator, err := state.NewRedisCoordinator(instanceID, redisClient, logger)
	if err != nil {
		cancel()
		logger.Error(
			"Error creating state manager",
			"error", err,
		)
		return nil, err
	}
	blockBuilder := blockbuilder.NewBlockBuilder(coordinator, engineCL, logger, buildDelay, buildDelayEmptyBlocks, feeReceipt, nil)

	lfm, err := leaderfollower.NewLeaderFollowerManager(
		instanceID,
		logger,
		redisClient,
		coordinator,
		blockBuilder,
	)
	if err != nil {
		cancel()
		logger.Error(
			"Error creating lfm",
			"error", err,
		)
		return nil, err
	}
	app := &MevCommitChain{
		stateManager: coordinator,
		blockBuilder: blockBuilder,
		logger:       logger,
		cancel:       cancel,
		lfm:          lfm,
	}

	logger.Info("MevCommitChain initialized")

	// Start leader election handling
	app.lfm.Start(ctx)

	return app, nil
}

func (app *MevCommitChain) Stop() {
	// Cancel the context to signal all goroutines to stop
	app.cancel()
	app.stateManager.Stop()
	err := app.lfm.WaitForGoroutinesToStop()
	if err != nil {
		app.logger.Error(
			"Error waiting for goroutines to stop",
			"error", err,
		)
	}
	app.logger.Info("MevCommitChain stopped gracefully")
}
