package state

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/primev/mev-commit/cl/types"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const blockStreamName = "mevcommit_block_stream"

type RedisClient interface {
	redis.Cmdable
	Close() error
}

type PipelineOperation func(redis.Pipeliner) error

type StateManager interface {
	ResetBlockState(ctx context.Context) error
	GetBlockBuildState(ctx context.Context) types.BlockBuildState
	Stop()
}

type StreamManager interface {
	ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error)
	AckMessage(ctx context.Context, messageID string) error
	Stop()
}

type Coordinator interface {
	StreamManager
	StateManager
	SaveBlockStateAndPublishToStream(ctx context.Context, bsState *types.BlockBuildState) error
	Stop()
}

type RedisStateManager struct {
	instanceID  string
	redisClient RedisClient
	logger      *slog.Logger

	blockStateKey   string
	blockBuildState *types.BlockBuildState
}

type RedisStreamManager struct {
	instanceID  string
	redisClient RedisClient
	logger      *slog.Logger

	groupName    string
	consumerName string
}

type RedisCoordinator struct {
	stateMgr    *RedisStateManager
	streamMgr   *RedisStreamManager
	redisClient RedisClient // added to close the client
	logger      *slog.Logger
}

func NewRedisStateManager(
	instanceID string,
	redisClient RedisClient,
	logger *slog.Logger,
) *RedisStateManager {
	return &RedisStateManager{
		instanceID:    instanceID,
		redisClient:   redisClient,
		logger:        logger,
		blockStateKey: fmt.Sprintf("blockBuildState:%s", instanceID),
	}
}

func NewRedisStreamManager(
	instanceID string,
	redisClient RedisClient,
	logger *slog.Logger,
) *RedisStreamManager {
	return &RedisStreamManager{
		instanceID:   instanceID,
		redisClient:  redisClient,
		logger:       logger,
		groupName:    fmt.Sprintf("mevcommit_consumer_group:%s", instanceID),
		consumerName: fmt.Sprintf("follower:%s", instanceID),
	}
}

func NewRedisCoordinator(
	instanceID string,
	redisClient RedisClient,
	logger *slog.Logger,
) (*RedisCoordinator, error) {
	stateMgr := NewRedisStateManager(instanceID, redisClient, logger)
	streamMgr := NewRedisStreamManager(instanceID, redisClient, logger)

	coordinator := &RedisCoordinator{
		stateMgr:    stateMgr,
		streamMgr:   streamMgr,
		redisClient: redisClient,
		logger:      logger,
	}

	if err := streamMgr.createConsumerGroup(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return coordinator, nil
}

func (s *RedisStateManager) executeTransaction(ctx context.Context, ops ...PipelineOperation) error {
	pipe := s.redisClient.TxPipeline()

	for _, op := range ops {
		if err := op(pipe); err != nil {
			return fmt.Errorf("failed to execute operation: %w", err)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("state transaction failed: %w", err)
	}

	return nil
}

func (s *RedisStateManager) loadOrInitializeBlockState(ctx context.Context) error {
	data, err := s.redisClient.Get(ctx, s.blockStateKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			s.blockBuildState = &types.BlockBuildState{
				CurrentStep: types.StepBuildBlock,
			}
			s.logger.Info("Leader block build state not found in Redis, initializing with default values")
			return s.saveBlockState(ctx)
		}
		return fmt.Errorf("failed to retrieve leader block build state: %w", err)
	}

	var state types.BlockBuildState
	if err := msgpack.Unmarshal([]byte(data), &state); err != nil {
		return fmt.Errorf("failed to deserialize leader block build state: %w", err)
	}

	s.logger.Info(
		"Loaded leader block build state",
		"CurrentStep", state.CurrentStep.String(),
	)
	s.blockBuildState = &state
	return nil
}

func (s *RedisStateManager) saveBlockState(ctx context.Context) error {
	return s.executeTransaction(ctx, s.saveBlockStateFunc(ctx, s.blockBuildState))
}

func (s *RedisStateManager) saveBlockStateFunc(ctx context.Context, bsState *types.BlockBuildState) PipelineOperation {
	return func(pipe redis.Pipeliner) error {
		data, err := msgpack.Marshal(bsState)
		if err != nil {
			return fmt.Errorf("failed to serialize block build state: %w", err)
		}

		pipe.Set(ctx, s.blockStateKey, data, 0)
		return nil
	}
}

func (s *RedisStateManager) ResetBlockState(ctx context.Context) error {
	s.blockBuildState = &types.BlockBuildState{
		CurrentStep: types.StepBuildBlock,
	}

	if err := s.saveBlockState(ctx); err != nil {
		return fmt.Errorf("failed to reset leader state: %w", err)
	}

	return nil
}

func (s *RedisStateManager) GetBlockBuildState(ctx context.Context) types.BlockBuildState {
	if s.blockBuildState == nil {
		s.logger.Error("Leader blockBuildState is not initialized")
		if err := s.loadOrInitializeBlockState(ctx); err != nil {
			s.logger.Warn(
				"Failed to load/init state",
				"error", err,
			)
			return types.BlockBuildState{}
		}
	}

	if s.blockBuildState == nil {
		s.logger.Error("Leader blockBuildState is still not initialized")
		return types.BlockBuildState{}
	}

	s.logger.Info(
		"Leader blockBuildState retrieved",
		"CurrentStep", s.blockBuildState.CurrentStep.String(),
	)
	// Return a copy of the state to prevent external modification
	return *s.blockBuildState
}

func (s *RedisStateManager) Stop() {
	if err := s.redisClient.Close(); err != nil {
		s.logger.Error("Error closing Redis client in StateManager", "error", err)
	}
}

func (s *RedisStreamManager) executeTransaction(ctx context.Context, ops ...PipelineOperation) error {
	pipe := s.redisClient.TxPipeline()

	for _, op := range ops {
		if err := op(pipe); err != nil {
			return fmt.Errorf("failed to execute operation: %w", err)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("stream transaction failed: %w", err)
	}

	return nil
}

func (s *RedisStreamManager) createConsumerGroup(ctx context.Context) error {
	if err := s.redisClient.XGroupCreateMkStream(ctx, blockStreamName, s.groupName, "0").Err(); err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			return fmt.Errorf("failed to create consumer group '%s': %w", s.groupName, err)
		}
	}
	return nil
}

func (s *RedisStreamManager) ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error) {
	args := &redis.XReadGroupArgs{
		Group:    s.groupName,
		Consumer: s.consumerName,
		Streams:  []string{blockStreamName, string(msgType)},
		Count:    1,
		Block:    time.Second,
	}

	messages, err := s.redisClient.XReadGroup(ctx, args).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error reading messages: %w", err)
	}

	return messages, nil
}

func (s *RedisStreamManager) AckMessage(ctx context.Context, messageID string) error {
	return s.executeTransaction(ctx, s.ackMessageFunc(ctx, messageID))
}

func (s *RedisStreamManager) ackMessageFunc(ctx context.Context, messageID string) PipelineOperation {
	return func(pipe redis.Pipeliner) error {
		pipe.XAck(ctx, blockStreamName, s.groupName, messageID)
		return nil
	}
}

func (s *RedisStreamManager) publishToStreamFunc(ctx context.Context, bsState *types.BlockBuildState) PipelineOperation {
	return func(pipe redis.Pipeliner) error {
		message := map[string]interface{}{
			"payload_id":         bsState.PayloadID,
			"execution_payload":  bsState.ExecutionPayload,
			"timestamp":          time.Now().UnixNano(),
			"sender_instance_id": s.instanceID,
		}

		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: blockStreamName,
			Values: message,
		})
		return nil
	}
}

func (s *RedisStreamManager) Stop() {
	if err := s.redisClient.Close(); err != nil {
		s.logger.Error(
			"Error closing Redis client in StreamManager",
			"error", err,
		)
	}
}

func (c *RedisCoordinator) SaveBlockStateAndPublishToStream(ctx context.Context, bsState *types.BlockBuildState) error {
	c.stateMgr.blockBuildState = bsState

	err := c.stateMgr.executeTransaction(
		ctx,
		c.stateMgr.saveBlockStateFunc(ctx, bsState),
		c.streamMgr.publishToStreamFunc(ctx, bsState),
	)
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (c *RedisCoordinator) ResetBlockState(ctx context.Context) error {
	return c.stateMgr.ResetBlockState(ctx)
}

func (c *RedisCoordinator) GetBlockBuildState(ctx context.Context) types.BlockBuildState {
	return c.stateMgr.GetBlockBuildState(ctx)
}

func (c *RedisCoordinator) ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error) {
	return c.streamMgr.ReadMessagesFromStream(ctx, msgType)
}

func (c *RedisCoordinator) AckMessage(ctx context.Context, messageID string) error {
	return c.streamMgr.AckMessage(ctx, messageID)
}

func (c *RedisCoordinator) Stop() {
	if err := c.redisClient.Close(); err != nil {
		c.logger.Error(
			"Error closing Redis client in StateManager",
			"error", err,
		)
	}
}
