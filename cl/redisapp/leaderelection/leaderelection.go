package leaderelection

import (
	"context"
	"log/slog"
	"time"

	"github.com/heyvito/go-leader/leader"
	"github.com/primev/mev-commit/cl/redisapp/blockbuilder"
	"github.com/primev/mev-commit/cl/redisapp/state"
	"github.com/primev/mev-commit/cl/redisapp/util"
	"github.com/redis/go-redis/v9"
)

type LeaderElectionHandler struct {
	// Leader election components
	leaderElection leader.Leader
	promotedCh     <-chan time.Time
	demotedCh      <-chan time.Time
	erroredCh      <-chan error

	// Context and cancellation
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// Dependencies
	logger       *slog.Logger
	instanceID   string
	stateManager state.StateManager
	blockBuilder *blockbuilder.BlockBuilder
	leader       *leaderManager
	follower     *followerManager
}

func NewLeaderElectionHandler(
	instanceID string,
	logger *slog.Logger,
	redisClient *redis.Client,
	stateManager state.StateManager,
	blockBuilder *blockbuilder.BlockBuilder,
) *LeaderElectionHandler {
	leaderCtx, cancel := context.WithCancel(context.Background())

	// Initialize leader election
	leaderOpts := leader.Opts{
		Redis: redisClient,
		TTL:   100 * time.Millisecond,
		Wait:  200 * time.Millisecond,
		Key:   "rapp_leader_election",
	}

	procLeader, promotedCh, demotedCh, erroredCh := leader.NewLeader(leaderOpts)

	follower := &followerManager{
		instanceID:   instanceID,
		stateManager: stateManager,
		blockBuilder: blockBuilder,
		logger:       logger,
	}

	leader := &leaderManager{
		instanceID:     instanceID,
		stateManager:   stateManager,
		blockBuilder:   blockBuilder,
		leaderElection: procLeader,
		logger:         logger,
	}

	return &LeaderElectionHandler{
		ctx:            leaderCtx,
		cancel:         cancel,
		logger:         logger,
		instanceID:     instanceID,
		leaderElection: procLeader,
		promotedCh:     promotedCh,
		demotedCh:      demotedCh,
		erroredCh:      erroredCh,
		leader:         leader,
		follower:       follower,
		stateManager:   stateManager,
		blockBuilder:   blockBuilder,
	}
}

func (leh *LeaderElectionHandler) HandleLeadershipEvents() {
	leh.logger.Info("Starting leader election event handler")
	leh.leaderElection.Start()

	leh.done = make(chan struct{})
	go func() {
		defer close(leh.done)
		if err := leh.initializeFollower(); err != nil {
			return
		}

		for {
			select {
			case <-leh.ctx.Done():
				leh.stopLeaderAndFollower()
				return
			case <-leh.promotedCh:
				continueLoop, err := leh.handlePromotion()
				if err != nil {
					leh.logger.Error("Error handling promotion", "error", err)
				}
				if !continueLoop {
					return
				}
			case <-leh.demotedCh:
				continueLoop, err := leh.handleDemotion()
				if err != nil {
					leh.logger.Error("Error handling demotion", "error", err)
					return
				}
				if !continueLoop {
					return
				}
			case err := <-leh.erroredCh:
				leh.logger.Error("Leader election error", "error", err)
			}
		}
	}()
}

func (leh *LeaderElectionHandler) initializeFollower() error {
	if err := leh.blockBuilder.ProcessLastPayload(leh.ctx); err != nil {
		leh.logger.Error("Error processing last payload", "error", err)
		return err
	}
	leh.follower.startFollowerLoop()
	return nil
}

func (leh *LeaderElectionHandler) stopLeaderAndFollower() {
	leh.leader.stopLeaderLoop()
	leh.follower.stopFollowerLoop()
}

func (leh *LeaderElectionHandler) handlePromotion() (bool, error) {
	leh.logger.Info("Promoting to leader")
	if !leh.follower.isSynced() {
		leh.logger.Info("Follower not synced, skipping promotion")
		leh.leaderElection.Stop()
		leh.logger.Info("Waiting for follower sync...")
		fswChannel := leh.follower.initSyncChannel()
		select {
		case <-fswChannel:
			leh.logger.Info("Sync finished, restarting leader election")
			leh.leaderElection.Start()
		case <-leh.ctx.Done():
			leh.logger.Info("App stopped, exiting")
			leh.follower.stopFollowerLoop()
			return false, nil
		}
	} else {
		leh.logger.Info("Promoted to leader")
		leh.follower.stopFollowerLoop()
		if err := leh.stateManager.RecoverLeaderState(); err != nil {
			leh.logger.Error("Failed to recover leader state", "error", err)
		}
		leh.leader.startLeaderLoop()
	}
	return true, nil
}

func (leh *LeaderElectionHandler) handleDemotion() (bool, error) {
	leh.logger.Info("Demoted from leader")

	if leh.follower.isRunning() {
		leh.logger.Info("Follower loop already running")
		return true, nil
	}

	leh.leader.stopLeaderLoop()

	if err := util.RetryWithInfiniteBackoff(leh.ctx, leh.logger, func() error {
		if err := leh.stateManager.LoadOrInitializeBlockState(leh.ctx); err != nil {
			leh.logger.Warn("Failed to load/init state, retrying...", "error", err)
			return err // will retry
		}
		return nil
	}); err != nil {
		leh.logger.Error("Failed to load/init state with retry, exiting")
		return false, err
	}

	if err := leh.blockBuilder.ProcessLastPayload(leh.ctx); err != nil {
		leh.logger.Error("Error processing last payload", "error", err)
		return false, err
	}

	leh.follower.startFollowerLoop()

	// In case leader demoted when geth is unreachable
	leh.leaderElection.Start()

	return true, nil
}

func (leh *LeaderElectionHandler) Stop() {
	leh.cancel()
	err := leh.leaderElection.Stop()
	if err != nil {
		leh.logger.Error("Error stopping leader election", "error", err)
	}
	<-leh.done
	leh.logger.Info("Leader election handler stopped")
}
