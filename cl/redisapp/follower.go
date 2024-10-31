package redisapp

import (
	"context"
	"log/slog"
	"sync"

	"github.com/primev/mev-commit/cl/redisapp/types"
)

type Follower struct {
	InstanceID   string
	stateManager StateManager
	blockBuilder *BlockBuilder
	logger       *slog.Logger
	cancel       context.CancelFunc
	done         chan struct{}

	fMutex          sync.Mutex
	isSynced        bool
	syncWaitChannel chan struct{}
}

func (f *Follower) startFollowerLoop() {
	f.logger.Info("Starting follower loop")
	followerCtx, followerCancel := context.WithCancel(context.Background())
	f.cancel = followerCancel
	f.done = make(chan struct{})

	go func() {
		defer close(f.done)
		f.followerLoop(followerCtx)
	}()
}

func (f *Follower) stopFollowerLoop() {
	if f.cancel != nil {
		f.cancel()
		<-f.done
		f.cancel = nil
		f.logger.Info("Follower loop stopped")
	}
}

func (f *Follower) followerLoop(ctx context.Context) {
	err := f.stateManager.CreateConsumerGroup(ctx)
	if err != nil {
		f.logger.Error("Failed to create consumer group", "error", err)
		return
	}

	for {
		f.logger.Info("Follower: loop is running")
		select {
		case <-ctx.Done():
			f.logger.Info("Follower: loop exiting")
			return
		default:
			f.fMutex.Lock()
			// First, try to read pending messages (unacknowledged entries)
			messages, err := f.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypePending)
			if err != nil {
				f.logger.Error("Error reading pending messages", "error", err)
				f.fMutex.Unlock()
				continue
			}

			if len(messages) == 0 || len(messages[0].Messages) == 0 {
				// Listen to the Redis Stream using the consumer group
				messages, err = f.stateManager.ReadMessagesFromStream(ctx, types.RedisMsgTypeNew)
				if err != nil {
					f.logger.Error("Error reading new messages", "error", err)
					f.fMutex.Unlock()
					continue
				}
			}

			f.isSynced = len(messages) == 0 || len(messages[0].Messages) == 0
			if f.isSynced && f.syncWaitChannel != nil && len(f.syncWaitChannel) != 1 {
				select {
				case f.syncWaitChannel <- struct{}{}:
					f.logger.Info("follower is synced")
				case <-ctx.Done():
					f.logger.Debug("ctx done")
				}
				f.syncWaitChannel = nil
			}
			for _, msg := range messages {
				for _, field := range msg.Messages {
					// Extract the PayloadID and ExecutionPayload
					payloadIDStr, ok := field.Values["payload_id"].(string)
					executionPayloadStr, okPayload := field.Values["execution_payload"].(string)
					senderInstanceID, okSenderID := field.Values["sender_instance_id"].(string)
					if !ok || !okPayload || !okSenderID || payloadIDStr == "" || executionPayloadStr == "" || senderInstanceID == "" {
						f.logger.Error("Follower: Invalid message format: missing payload_id or execution_payload")
						err = f.stateManager.AckMessage(ctx, field.ID)
						if err != nil {
							f.logger.Error("Failed to acknowledge invalid message", "error", err)
						}
						continue
					}

					if senderInstanceID == f.InstanceID {
						f.logger.Info("Follower: Received own message", "PayloadID", payloadIDStr)
						err = f.stateManager.AckMessage(ctx, field.ID)
						if err != nil {
							f.logger.Error("Failed to acknowledge own message", "error", err)
						}
						continue
					}

					f.logger.Info("Follower: Received message", "PayloadID", payloadIDStr)

					err := f.blockBuilder.finalizeBlock(ctx, payloadIDStr, executionPayloadStr, field.ID)
					if err != nil {
						f.logger.Error("Failed to finalize block", "error", err)
						continue
					}

					f.logger.Info("Follower: Finalized block", "PayloadID", payloadIDStr)
				}
			}
			f.fMutex.Unlock()
		}
	}
}

func (f *Follower) IsSynced() bool {
	f.fMutex.Lock()
	defer f.fMutex.Unlock()
	return f.isSynced
}

func (f *Follower) initSyncChannel() chan struct{} {
	if f.syncWaitChannel == nil {
		f.syncWaitChannel = make(chan struct{})
	}
	return f.syncWaitChannel
}

func (f *Follower) isRunning() bool {
	return f.cancel != nil
}
