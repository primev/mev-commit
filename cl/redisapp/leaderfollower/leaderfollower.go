package leaderfollower

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heyvito/go-leader/leader"
	"github.com/primev/mev-commit/cl/redisapp/blockbuilder"
	"github.com/primev/mev-commit/cl/redisapp/state"
	"github.com/primev/mev-commit/cl/redisapp/types"
	"github.com/primev/mev-commit/cl/redisapp/util"
	"github.com/redis/go-redis/v9"
)

type LeaderFollowerManager struct {
	isLeader              atomic.Bool
	isFollowerInitialized atomic.Bool
	stateManager          state.StateManager
	blockBuilder          *blockbuilder.BlockBuilder
	leaderProc            leader.Leader
	logger                *slog.Logger
	instanceID            string

	wg         sync.WaitGroup
	promotedCh <-chan time.Time
	demotedCh  <-chan time.Time
	erroredCh  <-chan error
}

func NewLeaderFollowerManager(
	instanceID string,
	logger *slog.Logger,
	redisClient *redis.Client,
	stateManager state.StateManager,
	blockBuilder *blockbuilder.BlockBuilder,
) (*LeaderFollowerManager, error) {
	// Initialize leader election
	leaderOpts := leader.Opts{
		Redis: redisClient,
		TTL:   100 * time.Millisecond,
		Wait:  200 * time.Millisecond,
		Key:   "rapp_leader_election",
	}

	leaderProc, promotedCh, demotedCh, erroredCh := leader.NewLeader(leaderOpts)

	worker := &LeaderFollowerManager{
		stateManager: stateManager,
		blockBuilder: blockBuilder,
		leaderProc:   leaderProc,
		logger:       logger,
		instanceID:   instanceID,
		promotedCh:   promotedCh,
		demotedCh:    demotedCh,
		erroredCh:    erroredCh,
		wg:           sync.WaitGroup{},
	}

	return worker, nil
}

func (lfm *LeaderFollowerManager) handleLeaderElection(ctx context.Context) {
	lfm.logger.Info("Starting leader election handler")
	lfm.leaderProc.Start()

	defer func() {
		err := lfm.leaderProc.Stop()
		if err != nil {
			lfm.logger.Error("Error stopping leader election", "error", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			lfm.logger.Debug("Leader election handler exiting")
			return
		case <-lfm.promotedCh:
			lfm.isLeader.Store(true)
			lfm.logger.Info("Node promoted to leader")
		case <-lfm.demotedCh:
			lfm.isLeader.Store(false)
			lfm.logger.Info("Node demoted from leader")
		case err := <-lfm.erroredCh:
			lfm.logger.Error("Leader election error", "error", err)
		}
	}
}

func (lfm *LeaderFollowerManager) Start(ctx context.Context) {
	lfm.wg.Add(2)

	go func() {
		defer lfm.wg.Done()
		lfm.handleLeaderElection(ctx)
	}()

	go func() {
		defer lfm.wg.Done()
		lfm.run(ctx)
	}()
}

func (lfm *LeaderFollowerManager) run(ctx context.Context) {
	lfm.logger.Info("LeaderFollowerManager started")
	for {
		select {
		case <-ctx.Done():
			lfm.logger.Debug("LeaderFollowerManager exiting")
			return
		default:
			lfm.logger.Info("LeaderFollowerManager running")
			if lfm.isLeader.Load() {
				lfm.isFollowerInitialized.Store(false)
				lfm.logger.Info("Leader: Starting leader work")
				if err := lfm.leaderWork(ctx); err != nil {
					lfm.logger.Error("Error in leader work", "error", err)
					if errors.Is(err, util.ErrFailedAfterNAttempts) {
						lfm.logger.Error("Leader: failed to reach geth node after max attempts, exiting")
						err := lfm.leaderProc.Stop()
						if err != nil {
							lfm.logger.Error("Leader: Failed to stop leader election", "error", err)
						}
						lfm.blockBuilder.LastCallTime = time.Time{}
					}
				}
			} else {
				if !lfm.isFollowerInitialized.Load() {
					if err := lfm.blockBuilder.ProcessLastPayload(ctx); err != nil {
						lfm.logger.Error("Error processing last payload", "error", err)
						continue
					}
					lfm.isFollowerInitialized.Store(true)
				}
				if err := lfm.followerWork(ctx); err != nil {
					lfm.logger.Error("Error in follower work", "error", err)
				}

				// will do nothing, in case election already started
				// needed, in case leader stopped being a leader bcs of geth issue
				lfm.leaderProc.Start()
			}
		}
	}
}

func (lfm *LeaderFollowerManager) leaderWork(ctx context.Context) error {
	lfm.logger.Info("Leader: Performing leader tasks")

	// Ensure state is synchronized before starting leader tasks
	isSync := lfm.HaveMessagesToProcess(ctx)

	// if not sync, return to wait for the signal to be a leader
	if !isSync {
		lfm.logger.Info("Leader: State is not synchronized, waiting for follower to catch up")
		lfm.leaderProc.Stop()
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			lfm.logger.Debug("Leader: Exiting")
			return nil
		default:
			bbState := lfm.stateManager.GetBlockBuildState(ctx)
			currentStep := bbState.CurrentStep

			switch currentStep {
			case types.StepBuildBlock:
				lfm.logger.Info("Leader: StepBuildBlock")
				if err := lfm.blockBuilder.GetPayload(ctx); err != nil {
					lfm.logger.Error("Leader: GetPayload failed", "error", err)
					if resetErr := lfm.stateManager.ResetBlockState(ctx); resetErr != nil {
						lfm.logger.Error("Leader: Failed to reset block state", "error", resetErr)
					}

					return err
				}
			case types.StepFinalizeBlock:
				lfm.logger.Info("Leader: StepFinalizeBlock")
				if err := lfm.blockBuilder.FinalizeBlock(ctx, bbState.PayloadID, bbState.ExecutionPayload, ""); err != nil {
					lfm.logger.Error("Leader: FinalizeBlock failed", "error", err)
					if resetErr := lfm.stateManager.ResetBlockState(ctx); resetErr != nil {
						lfm.logger.Error("Leader: Failed to reset block state", "error", resetErr)
					}
					return err
				}
				if err := lfm.stateManager.ResetBlockState(ctx); err != nil {
					lfm.logger.Error("Leader: Failed to reset block state", "error", err)
					return err
				}
			default:
				lfm.logger.Warn("Leader: Unknown current step", "current_step", currentStep.String())
				if err := lfm.stateManager.ResetBlockState(ctx); err != nil {
					lfm.logger.Error("Leader: Failed to reset block state", "error", err)
					return err
				}
			}
		}
	}
}

func (lfm *LeaderFollowerManager) followerWork(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			lfm.logger.Debug("Follower: Exiting")
			return nil
		default:
			lfm.logger.Info("Follower: Performing follower tasks")

			messages, err := lfm.readMessages(ctx)
			if err != nil {
				lfm.logger.Error("Follower: Error reading messages", "error", err)
				continue
			}

			if len(messages) == 0 || len(messages[0].Messages) == 0 {
				lfm.logger.Info("Follower: No messages to process")
				return nil
			}

			// Process messages
			for _, msg := range messages {
				for _, field := range msg.Messages {
					// Extract fields
					payloadIDStr, ok := field.Values["payload_id"].(string)
					executionPayloadStr, okPayload := field.Values["execution_payload"].(string)
					senderInstanceID, okSenderID := field.Values["sender_instance_id"].(string)
					if !ok || !okPayload || !okSenderID || payloadIDStr == "" || executionPayloadStr == "" || senderInstanceID == "" {
						lfm.logger.Error("Follower: Invalid message format")
						// Acknowledge the message to avoid reprocessing
						if ackErr := lfm.stateManager.AckMessage(ctx, field.ID); ackErr != nil {
							lfm.logger.Error("Follower: Failed to acknowledge message", "error", ackErr)
						}
						continue
					}

					// Ignore messages sent by self
					if senderInstanceID == lfm.instanceID {
						lfm.logger.Info("Follower: Ignoring own message", "PayloadID", payloadIDStr)
						if ackErr := lfm.stateManager.AckMessage(ctx, field.ID); ackErr != nil {
							lfm.logger.Error("Follower: Failed to acknowledge own message", "error", ackErr)
						}
						continue
					}

					lfm.logger.Info("Follower: Processing message", "PayloadID", payloadIDStr)

					// Finalize block
					if err := lfm.blockBuilder.FinalizeBlock(ctx, payloadIDStr, executionPayloadStr, field.ID); err != nil {
						lfm.logger.Error("Follower: Failed to finalize block", "error", err)
						continue
					}

					lfm.logger.Info("Follower: Successfully finalized block", "PayloadID", payloadIDStr)
				}
			}
		}
	}
}

func (lfm *LeaderFollowerManager) readMessages(ctx context.Context) ([]redis.XStream, error) {
	// Try to read pending messages first
	messages, err := lfm.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypePending)
	if err != nil {
		lfm.logger.Error("Follower: Error reading pending messages", "error", err)
		return nil, err
	}

	// If no pending messages, read new messages
	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		messages, err = lfm.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypeNew)
		if err != nil {
			lfm.logger.Error("Follower: Error reading new messages", "error", err)
			return nil, err
		}
	}

	return messages, nil
}

func (lfm *LeaderFollowerManager) WaitForGoroutinesToStop() error {
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		lfm.wg.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		lfm.logger.Error("Workers still running")
		return errors.New("workers still running")
	case <-closed:
		return nil
	}
}

func (lfm *LeaderFollowerManager) HaveMessagesToProcess(ctx context.Context) bool {
	messages, err := lfm.readMessages(ctx)
	if err != nil {
		lfm.logger.Error("Error reading messages", "error", err)
		return false
	}
	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		return false
	}
	return true
}
