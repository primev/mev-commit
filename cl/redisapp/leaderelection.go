package redisapp

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/heyvito/go-leader/leader"
)

type LeaderElectionHandler struct {
	// Leader election components
	leaderElection leader.Leader
	promotedCh     <-chan time.Time
	demotedCh      <-chan time.Time
	erroredCh      <-chan error

	// Leader state
	isLeaderMutex sync.RWMutex
	isLeader      bool

	// Context and WaitGroup for goroutines
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// Dependencies
	logger       *slog.Logger
	instanceID   string
	stateManager StateManager
	stepsManager *StepsManager
	leader       *Leader
	follower     *Follower
}

func NewLeaderElectionHandler(
	instanceID string,
	logger *slog.Logger,
	procLeader leader.Leader,
	promotedCh <-chan time.Time,
	demotedCh <-chan time.Time,
	erroredCh <-chan error,
	leader *Leader,
	follower *Follower,
	stateManager StateManager,
	stepsManager *StepsManager,
) *LeaderElectionHandler {
	leaderCtx, cancel := context.WithCancel(context.Background())

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
		stepsManager:   stepsManager,
	}
}

func (leh *LeaderElectionHandler) handleLeadershipEvents() {
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
	if err := leh.stepsManager.processLastPayload(leh.ctx); err != nil {
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
	if !leh.follower.IsSynced() {
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
		leh.setIsLeader(true)
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
	leh.setIsLeader(false)
	leh.logger.Info("Demoted from leader")

	if leh.follower.isRunning() {
		leh.logger.Info("Follower loop already running")
		return true, nil
	}

	leh.leader.stopLeaderLoop()

	if err := retryWithInfiniteBackoff(leh.ctx, leh.logger, func() error {
		if err := leh.stateManager.LoadOrInitializeBlockState(leh.ctx); err != nil {
			leh.logger.Warn("Failed to load/init state, retrying...", "error", err)
			return err // will retry
		}
		return nil
	}); err != nil {
		leh.logger.Error("Failed to load/init state with retry, exiting")
		return false, err
	}
	
	if err := leh.stepsManager.processLastPayload(leh.ctx); err != nil {
		leh.logger.Error("Error processing last payload", "error", err)
		return false, err
	}

	leh.follower.startFollowerLoop()

	// In case leader demoted when geth is unreachable
	leh.leaderElection.Start()

	return true, nil
}

func (leh *LeaderElectionHandler) setIsLeader(value bool) {
	leh.isLeaderMutex.Lock()
	defer leh.isLeaderMutex.Unlock()
	leh.isLeader = value
}

func (leh *LeaderElectionHandler) IsLeader() bool {
	leh.isLeaderMutex.RLock()
	defer leh.isLeaderMutex.RUnlock()
	return leh.isLeader
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
