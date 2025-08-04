package leaderfollower

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heyvito/go-leader/leader"
	"github.com/primev/mev-commit/cl/blockbuilder"
	"github.com/primev/mev-commit/cl/types"
	"github.com/primev/mev-commit/cl/util"
	"github.com/redis/go-redis/v9"
)

type LeaderFollowerManager struct {
	isLeader              atomic.Bool
	isFollowerInitialized atomic.Bool
	stateManager          stateManager
	blockBuilder          blockBuilder
	leaderProc            leader.Leader
	logger                *slog.Logger
	instanceID            string

	wg         sync.WaitGroup
	promotedCh <-chan time.Time
	demotedCh  <-chan time.Time
	erroredCh  <-chan error
}

type blockBuilder interface {
	// Retrieves the latest payload and ensures it meets necessary conditions
	GetPayload(ctx context.Context) error

	// Finalizes a block, pushing it to the EVM and updating the execution state
	FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error

	// Processes any unfinished payload from a previous session
	ProcessLastPayload(ctx context.Context) error
}

// todo: work with block state through block builder, not directly
type stateManager interface {
	// state related methods
	GetBlockBuildState(ctx context.Context) types.BlockBuildState
	ResetBlockState(ctx context.Context) error

	// stream related methods
	AckMessage(ctx context.Context, messageID string) error
	ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error)
}

func NewLeaderFollowerManager(
	instanceID string,
	logger *slog.Logger,
	redisClient *redis.Client,
	stateManager stateManager,
	blockBuilder blockBuilder,
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
			lfm.logger.Error(
				"Error stopping leader election",
				"error", err,
			)
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
			lfm.logger.Error(
				"Leader election error",
				"error", err,
			)
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
					lfm.logger.Error(
						"Error in leader work",
						"error", err,
					)
				}
			} else {
				if !lfm.isFollowerInitialized.Load() {
					if err := lfm.blockBuilder.ProcessLastPayload(ctx); err != nil {
						lfm.logger.Error(
							"Error processing last payload",
							"error", err,
						)
						continue
					}
					lfm.isFollowerInitialized.Store(true)
				}
				if err := lfm.followerWork(ctx); err != nil {
					lfm.logger.Error(
						"Error in follower work",
						"error", err,
					)
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
	isHaveMessagesToProcess := lfm.HaveMessagesToProcess(ctx)

	// if has messages to process, return to wait for the signal to be a leader
	if isHaveMessagesToProcess {
		lfm.logger.Info("Leader: State is not synchronized, waiting for follower to catch up")
		err := lfm.leaderProc.Stop()
		if err != nil {
			lfm.logger.Error(
				"Leader: Failed to stop leader election",
				"error", err,
			)
			return err
		}
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
			err := func() error {
				switch currentStep {
				case types.StepBuildBlock:
					lfm.logger.Info("Leader: StepBuildBlock")
					if err := lfm.blockBuilder.GetPayload(ctx); err != nil {
						if errors.Is(err, blockbuilder.ErrEmptyBlock) {
							lfm.logger.Info("Leader: Empty block, skipping")
							return nil
						}
						lfm.logger.Error(
							"Leader: GetPayload failed",
							"error", err,
						)
						if resetErr := lfm.stateManager.ResetBlockState(ctx); resetErr != nil {
							lfm.logger.Error(
								"Leader: Failed to reset block state",
								"error", resetErr,
							)
						}

						return err
					}
				case types.StepFinalizeBlock:
					lfm.logger.Info("Leader: StepFinalizeBlock")
					if err := lfm.blockBuilder.FinalizeBlock(ctx, bbState.PayloadID, bbState.ExecutionPayload, ""); err != nil {
						lfm.logger.Error(
							"Leader: FinalizeBlock failed",
							"error", err,
						)
						return err
					}
					if err := lfm.stateManager.ResetBlockState(ctx); err != nil {
						lfm.logger.Error(
							"Leader: Failed to reset block state",
							"error", err,
						)
						return err
					}
				default:
					lfm.logger.Warn("Leader: Unknown current step", "current_step", currentStep.String())
					if err := lfm.stateManager.ResetBlockState(ctx); err != nil {
						lfm.logger.Error(
							"Leader: Failed to reset block state",
							"error", err,
						)
						return err
					}
				}
				return nil
			}()
			if err != nil {
				// in that case app will stop being a leader so we need to return
				if errors.Is(err, util.ErrFailedAfterNAttempts) {
					lfm.logger.Error("Leader: failed to reach geth node after max attempts, exiting")
					stopElecErr := lfm.leaderProc.Stop()
					// todo: refactor to generate timestamp outside blockbuilder
					if stopElecErr != nil {
						lfm.logger.Error(
							"Leader: Failed to stop leader election",
							"error", stopElecErr,
						)
						return stopElecErr
					}
					return err
				}
				// otherwise there is a problem with redis/payload, so we just log it and continue
				lfm.logger.Error(
					"Leader: Error in leader work",
					"error", err,
				)
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
				lfm.logger.Error(
					"Follower: Error reading messages",
					"error", err,
				)
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
							lfm.logger.Error(
								"Follower: Failed to acknowledge message",
								"error", ackErr,
							)
						}
						continue
					}

					// Ignore messages sent by self
					if senderInstanceID == lfm.instanceID {
						lfm.logger.Info(
							"Follower: Ignoring own message",
							"PayloadID", payloadIDStr,
						)
						if ackErr := lfm.stateManager.AckMessage(ctx, field.ID); ackErr != nil {
							lfm.logger.Error(
								"Follower: Failed to acknowledge own message",
								"error", ackErr,
							)
						}
						continue
					}

					lfm.logger.Info(
						"Follower: Processing message",
						"PayloadID", payloadIDStr,
					)

					// Finalize block
					// msg will be acknowledged inside of state manager with execution head saved
					if err := lfm.blockBuilder.FinalizeBlock(ctx, payloadIDStr, executionPayloadStr, field.ID); err != nil {
						lfm.logger.Error(
							"Follower: Failed to finalize block",
							"error", err,
						)
						continue
					}

					err = lfm.stateManager.AckMessage(ctx, field.ID)
					if err != nil {
						lfm.logger.Error(
							"Follower: Failed to acknowledge message",
							"error", err,
						)
					} else {
						lfm.logger.Info(
							"Follower: Successfully acknowledged message",
							"PayloadID", payloadIDStr,
						)
					}

					lfm.logger.Info(
						"Follower: Successfully finalized block",
						"PayloadID", payloadIDStr,
					)
				}
			}
		}
	}
}

func (lfm *LeaderFollowerManager) readMessages(ctx context.Context) ([]redis.XStream, error) {
	// Try to read pending messages first
	messages, err := lfm.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypePending)
	if err != nil {
		lfm.logger.Error(
			"Follower: Error reading pending messages",
			"error", err,
		)
		return nil, err
	}

	// If no pending messages, read new messages
	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		messages, err = lfm.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypeNew)
		if err != nil {
			lfm.logger.Error(
				"Follower: Error reading new messages",
				"error", err,
			)
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
		lfm.logger.Error(
			"Error reading messages",
			"error", err,
		)
		return false
	}
	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		return false
	}
	return true
}
