package redisapp

import (
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/redis/go-redis/v9"
)

type EngineClient interface {
	NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
		beaconRoot *common.Hash) (engine.PayloadStatusV1, error)

	ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
		payloadAttributes *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)

	GetPayloadV3(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)
}

type MevCommitChain struct {
	InstanceID       string
	engineCl         EngineClient
	genesisBlockHash string
	logger           *slog.Logger

	cancel context.CancelFunc

	// Managers and components
	stateManager          StateManager
	blockBuilder          *BlockBuilder
	leaderElectionHandler *LeaderElectionHandler
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

	stateManager := NewRedisStateManager(instanceID, redisClient, logger, genesisBlockHash)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     engineCL,
		logger:       logger,
		buildDelay:   buildDelay,
		buildDelayMs: uint64(buildDelay.Milliseconds()),
	}

	// Initialize LeaderElectionHandler
	leaderElectionHandler := NewLeaderElectionHandler(
		instanceID,
		logger,
		redisClient,
		stateManager,
		blockBuilder,
	)

	app := &MevCommitChain{
		InstanceID:            instanceID,
		stateManager:          stateManager,
		blockBuilder:          blockBuilder,
		engineCl:              engineCL,
		genesisBlockHash:      genesisBlockHash,
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
	app.leaderElectionHandler.handleLeadershipEvents()

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
