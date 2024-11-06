package state

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/primev/mev-commit/cl/redisapp/types"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const blockStreamName = "mevcommit_block_stream"

type RedisClient interface {
	redis.Cmdable
	Close() error
}

type StateManager interface {
	SaveExecutionHead(ctx context.Context, head *types.ExecutionHead) error
	LoadExecutionHead(ctx context.Context) (*types.ExecutionHead, error)
	LoadOrInitializeBlockState(ctx context.Context) error
	SaveBlockState(ctx context.Context) error
	ResetBlockState(ctx context.Context) error
	SaveExecutionHeadAndAck(ctx context.Context, head *types.ExecutionHead, messageID string) error
	SaveBlockStateAndPublishToStream(ctx context.Context, bsState *types.BlockBuildState) error
	GetBlockBuildState(ctx context.Context) types.BlockBuildState
	CreateConsumerGroup(ctx context.Context) error
	RecoverLeaderState() error
	ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error)
	AckMessage(ctx context.Context, messageID string) error
	Stop()
}

type RedisStateManager struct {
	instanceID       string
	redisClient      RedisClient
	logger           *slog.Logger
	genesisBlockHash string
	groupName        string
	consumerName     string

	blockStateKey   string
	blockBuildState *types.BlockBuildState
}

func NewRedisStateManager(
	instanceID string,
	redisClient RedisClient,
	logger *slog.Logger,
	genesisBlockHash string,
) StateManager {
	return &RedisStateManager{
		instanceID:       instanceID,
		redisClient:      redisClient,
		logger:           logger,
		genesisBlockHash: genesisBlockHash,
		blockStateKey:    fmt.Sprintf("blockBuildState:%s", instanceID),
		groupName:        fmt.Sprintf("mevcommit_consumer_group:%s", instanceID),
		consumerName:     fmt.Sprintf("follower:%s", instanceID),
	}
}

func (s *RedisStateManager) SaveExecutionHead(ctx context.Context, head *types.ExecutionHead) error {
	data, err := msgpack.Marshal(head)
	if err != nil {
		return fmt.Errorf("failed to serialize execution head: %w", err)
	}

	key := fmt.Sprintf("executionHead:%s", s.instanceID)
	if err := s.redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save execution head to Redis: %w", err)
	}

	return nil
}

func (s *RedisStateManager) LoadExecutionHead(ctx context.Context) (*types.ExecutionHead, error) {
	key := fmt.Sprintf("executionHead:%s", s.instanceID)
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			s.logger.Info("executionHead not found in Redis, initializing with default values")
			hashBytes, decodeErr := hex.DecodeString(s.genesisBlockHash)
			if decodeErr != nil {
				s.logger.Error("Error decoding genesis block hash", "error", decodeErr)
				return nil, decodeErr
			}
			head := &types.ExecutionHead{BlockHash: hashBytes, BlockTime: uint64(time.Now().UnixMilli())}
			if saveErr := s.SaveExecutionHead(ctx, head); saveErr != nil {
				return nil, saveErr
			}
			return head, nil
		}
		return nil, fmt.Errorf("failed to retrieve execution head: %w", err)
	}

	var head types.ExecutionHead
	if err := msgpack.Unmarshal([]byte(data), &head); err != nil {
		return nil, fmt.Errorf("failed to deserialize execution head: %w", err)
	}

	return &head, nil
}

func (s *RedisStateManager) LoadOrInitializeBlockState(ctx context.Context) error {
	data, err := s.redisClient.Get(ctx, s.blockStateKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			s.blockBuildState = &types.BlockBuildState{
				CurrentStep: types.StepBuildBlock,
			}
			s.logger.Info("Leader block build state not found in Redis, initializing with default values")
			return s.SaveBlockState(ctx)
		}
		return fmt.Errorf("failed to retrieve leader block build state: %w", err)
	}

	var state types.BlockBuildState
	if err := msgpack.Unmarshal([]byte(data), &state); err != nil {
		return fmt.Errorf("failed to deserialize leader block build state: %w", err)
	}

	s.logger.Info("Loaded leader block build state", "CurrentStep", state.CurrentStep.String())
	s.blockBuildState = &state
	return nil
}

func (s *RedisStateManager) SaveBlockState(ctx context.Context) error {
	data, err := msgpack.Marshal(s.blockBuildState)
	if err != nil {
		return fmt.Errorf("failed to serialize leader block build state: %w", err)
	}

	if err := s.redisClient.Set(ctx, s.blockStateKey, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save leader block build state to Redis: %w", err)
	}

	return nil
}

func (s *RedisStateManager) ResetBlockState(ctx context.Context) error {
	s.blockBuildState = &types.BlockBuildState{
		CurrentStep: types.StepBuildBlock,
	}

	if err := s.SaveBlockState(ctx); err != nil {
		return fmt.Errorf("failed to reset leader state: %w", err)
	}

	return nil
}

func (s *RedisStateManager) SaveExecutionHeadAndAck(ctx context.Context, head *types.ExecutionHead, messageID string) error {
	data, err := msgpack.Marshal(head)
	if err != nil {
		return fmt.Errorf("failed to serialize execution head: %w", err)
	}

	key := fmt.Sprintf("executionHead:%s", s.instanceID)
	pipe := s.redisClient.TxPipeline()

	pipe.Set(ctx, key, data, 0)
	pipe.XAck(ctx, blockStreamName, s.groupName, messageID)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("Follower: Execution head saved and message acknowledged", "MessageID", messageID)
	return nil
}

func (s *RedisStateManager) SaveBlockStateAndPublishToStream(ctx context.Context, bsState *types.BlockBuildState) error {
	s.blockBuildState = bsState
	data, err := msgpack.Marshal(bsState)
	if err != nil {
		return fmt.Errorf("failed to serialize leader block build state: %w", err)
	}

	pipe := s.redisClient.Pipeline()
	pipe.Set(ctx, s.blockStateKey, data, 0)

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

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline failed: %w", err)
	}

	return nil
}

func (s *RedisStateManager) GetBlockBuildState(ctx context.Context) types.BlockBuildState {
	if s.blockBuildState == nil {
		s.logger.Error("Leader blockBuildState is not initialized")
		if err := s.LoadOrInitializeBlockState(ctx); err != nil {
			s.logger.Warn("Failed to load/init state", "error", err)
			return types.BlockBuildState{}
		}
	}

	if s.blockBuildState == nil {
		s.logger.Error("Leader blockBuildState is still not initialized")
		return types.BlockBuildState{}
	}

	s.logger.Info("Leader blockBuildState retrieved", "CurrentStep", s.blockBuildState.CurrentStep.String())
	// Return a copy of the state to prevent external modification
	return *s.blockBuildState
}

func (s *RedisStateManager) CreateConsumerGroup(ctx context.Context) error {
	if err := s.redisClient.XGroupCreateMkStream(ctx, blockStreamName, s.groupName, "0").Err(); err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			return fmt.Errorf("failed to create consumer group '%s': %w", s.groupName, err)
		}
	}
	return nil
}

func (s *RedisStateManager) RecoverLeaderState() error {
	if s.blockBuildState == nil {
		return errors.New("leader blockBuildState is not initialized")
	}

	switch s.blockBuildState.CurrentStep {
	case types.StepBuildBlock:
		s.logger.Info("Leader: Starting block build process")
	case types.StepFinalizeBlock:
		s.logger.Info("Leader: Resuming from FinalizeBlock", "PayloadID", s.blockBuildState.PayloadID)
	default:
		return fmt.Errorf("leader: unknown build step: %d", s.blockBuildState.CurrentStep)
	}

	return nil
}

func (s *RedisStateManager) ReadMessagesFromStream(ctx context.Context, msgType types.RedisMsgType) ([]redis.XStream, error) {
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

func (s *RedisStateManager) AckMessage(ctx context.Context, messageID string) error {
	if err := s.redisClient.XAck(ctx, blockStreamName, s.groupName, messageID).Err(); err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}
	return nil
}

func (s *RedisStateManager) Stop() {
	if err := s.redisClient.Close(); err != nil {
		s.logger.Error("Error closing Redis client", "error", err)
	}
}
