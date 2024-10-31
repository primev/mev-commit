package redisapp

import (
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/heyvito/go-leader/leader"
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
	stepsManager          *StepsManager
	leader                *Leader
	follower              *Follower
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

	// Initialize leader election
	leaderOpts := leader.Opts{
		Redis: redisClient,
		TTL:   100 * time.Millisecond,
		Wait:  200 * time.Millisecond,
		Key:   "rapp_leader_election",
	}

	procLeader, promotedCh, demotedCh, erroredCh := leader.NewLeader(leaderOpts)

	stateManager := NewRedisStateManager(instanceID, redisClient, logger, genesisBlockHash)

	stepsManager := &StepsManager{
		stateManager: stateManager,
		engineCl:     engineCL,
		logger:       logger,
		buildDelay:   buildDelay,
		buildDelayMs: uint64(buildDelay.Milliseconds()),
	}

	follower := &Follower{
		InstanceID:   instanceID,
		stateManager: stateManager,
		stepsManager: stepsManager,
		logger:       logger,
	}

	leader := &Leader{
		InstanceID:     instanceID,
		stateManager:   stateManager,
		stepsManager:   stepsManager,
		leaderElection: procLeader,
		logger:         logger,
	}

	// Initialize LeaderElectionHandler
	leaderElectionHandler := NewLeaderElectionHandler(
		instanceID,
		logger,
		procLeader,
		promotedCh,
		demotedCh,
		erroredCh,
		leader,
		follower,
		stateManager,
		stepsManager,
	)

	app := &MevCommitChain{
		InstanceID:            instanceID,
		stateManager:          stateManager,
		stepsManager:          stepsManager,
		engineCl:              engineCL,
		genesisBlockHash:      genesisBlockHash,
		logger:                logger,
		cancel:                cancel,
		leader:                leader,
		follower:              follower,
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
