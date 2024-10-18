package redisapp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/primev/mev-commit/cl/redisapp/types"
	"github.com/redis/go-redis/v9"
)

const redisStreamName = "mevcommit_block_stream"

type RedisClient interface {
	redis.Cmdable
	Close() error
}

type RedisStateManager struct {
	InstanceID       string
	redisClient      RedisClient
	logger           Logger
	genesisBlockHash string
	groupName        string
	consumerName     string

	blockStateKey   string
	blockBuildState *types.BlockBuildState
	blockStateMutex sync.Mutex
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

func NewRedisStateManager(
	instanceID string,
	redisClient RedisClient,
	logger Logger,
	genesisBlockHash string,
) StateManager {
	return &RedisStateManager{
		InstanceID:       instanceID,
		redisClient:      redisClient,
		logger:           logger,
		genesisBlockHash: genesisBlockHash,
		blockStateKey:    fmt.Sprintf("blockBuildState:%s", instanceID),
		groupName:        fmt.Sprintf("mevcommit_consumer_group:%s", instanceID),
		consumerName:     fmt.Sprintf("follower:%s", instanceID),
	}
}

func (s *RedisStateManager) SaveExecutionHead(ctx context.Context, head *types.ExecutionHead) error {
	data, err := json.Marshal(head)
	if err != nil {
		return fmt.Errorf("failed to serialize execution head: %w", err)
	}

	key := fmt.Sprintf("executionHead:%s", s.InstanceID)
	err = s.redisClient.Set(ctx, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save execution head to Redis: %w", err)
	}

	return nil
}

func (s *RedisStateManager) LoadExecutionHead(ctx context.Context) (*types.ExecutionHead, error) {
	key := fmt.Sprintf("executionHead:%s", s.InstanceID)
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			s.logger.Info("executionHead not found in Redis, initializing with default values")
			hashBytes, err := hex.DecodeString(s.genesisBlockHash)
			if err != nil {
				s.logger.Error("Error decoding genesis block hash", "error", err)
				return nil, err
			}
			s.SaveExecutionHead(ctx, &types.ExecutionHead{BlockHash: hashBytes})
			return &types.ExecutionHead{BlockHash: hashBytes}, nil
		}
		return nil, fmt.Errorf("failed to retrieve execution head: %w", err)
	}

	var head types.ExecutionHead
	err = json.Unmarshal([]byte(data), &head)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize execution head: %w", err)
	}

	return &head, nil
}

func (s *RedisStateManager) LoadOrInitializeBlockState(ctx context.Context) error {
	data, err := s.redisClient.Get(ctx, s.blockStateKey).Result()
	if err != nil {
		if err == redis.Nil {
			s.blockBuildState = &types.BlockBuildState{
				CurrentStep: types.StepBuildBlock,
			}
			return s.SaveBlockState(ctx)
		}
		return fmt.Errorf("failed to retrieve leader block build state: %w", err)
	}

	var state types.BlockBuildState
	err = json.Unmarshal([]byte(data), &state)
	if err != nil {
		return fmt.Errorf("failed to deserialize leader block build state: %w", err)
	}

	s.logger.Info("Loaded leader block build state", "CurrentStep", state.CurrentStep.String())

	s.blockBuildState = &state
	return nil
}

func (s *RedisStateManager) SaveBlockState(ctx context.Context) error {
	s.blockStateMutex.Lock()
	defer s.blockStateMutex.Unlock()

	data, err := json.Marshal(s.blockBuildState)
	if err != nil {
		return fmt.Errorf("failed to serialize leader block build state: %w", err)
	}

	_, err = s.redisClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.Set(ctx, s.blockStateKey, data, 0).Err(); err != nil {
			return fmt.Errorf("failed to save leader block build state to Redis: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to save leader block build state to Redis: %w", err)
	}

	return nil
}

func (s *RedisStateManager) ResetBlockState(ctx context.Context) error {
	s.blockStateMutex.Lock()
	s.blockBuildState = &types.BlockBuildState{
		CurrentStep: types.StepBuildBlock,
	}
	s.blockStateMutex.Unlock()

	err := s.SaveBlockState(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset leader state: %w", err)
	}

	return nil
}

func (s *RedisStateManager) SaveExecutionHeadAndAck(ctx context.Context, head *types.ExecutionHead, messageID string) error {
	data, err := json.Marshal(head)
	if err != nil {
		return fmt.Errorf("failed to serialize execution head: %w", err)
	}

	key := fmt.Sprintf("executionHead:%s", s.InstanceID)

	_, err = s.redisClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.Set(ctx, key, data, 0).Err(); err != nil {
			return fmt.Errorf("failed to queue SET command: %w", err)
		}

		if err := pipe.XAck(ctx, redisStreamName, s.groupName, messageID).Err(); err != nil {
			return fmt.Errorf("failed to queue XACK command: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("Follower: Execution head saved and message acknowledged", "MessageID", messageID)
	return nil
}

func (s *RedisStateManager) SaveBlockStateAndPublishToStream(ctx context.Context, bsState *types.BlockBuildState) error {
	s.blockStateMutex.Lock()
	defer s.blockStateMutex.Unlock()

	s.blockBuildState = bsState
	data, err := json.Marshal(s.blockBuildState)
	if err != nil {
		return fmt.Errorf("failed to serialize leader block build state: %w", err)
	}

	_, err = s.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.Set(ctx, s.blockStateKey, data, 0).Err(); err != nil {
			return fmt.Errorf("failed to save leader block build state to Redis: %w", err)
		}

		message := map[string]interface{}{
			"payload_id":         bsState.PayloadID,
			"execution_payload":  bsState.ExecutionPayload,
			"timestamp":          time.Now().UnixNano(),
			"sender_instance_id": s.InstanceID,
		}

		if _, err := pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: redisStreamName,
			Values: message,
		}).Result(); err != nil {
			return fmt.Errorf("failed to add message to Redis Stream: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("pipeline failed: %w", err)
	}
	return nil
}

func (s *RedisStateManager) GetBlockBuildState(ctx context.Context) types.BlockBuildState {
	s.blockStateMutex.Lock()
	defer s.blockStateMutex.Unlock()

	// return a copy of the state
	return *s.blockBuildState
}

func (s *RedisStateManager) CreateConsumerGroup(ctx context.Context) error {
	groupName := fmt.Sprintf("mevcommit_consumer_group:%s", s.InstanceID)

	err := s.redisClient.XGroupCreateMkStream(ctx, redisStreamName, groupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("failed to create consumer group '%s': %w", groupName, err)
	}
	return nil
}

func (s *RedisStateManager) RecoverLeaderState() error {
	s.blockStateMutex.Lock()
	defer s.blockStateMutex.Unlock()

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
	messages, err := s.redisClient.XReadGroup(
		ctx,
		&redis.XReadGroupArgs{
			Group:    s.groupName,
			Consumer: s.consumerName,
			Streams:  []string{redisStreamName, string(msgType)},
			Count:    1,
			Block:    1 * time.Second,
		},
	).Result()

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("error reading messages: %w", err)
	}

	return messages, nil
}

func (s *RedisStateManager) AckMessage(ctx context.Context, messageID string) error {
	err := s.redisClient.XAck(ctx, redisStreamName, s.groupName, messageID).Err()
	if err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}
	return nil
}

func (s *RedisStateManager) Stop() {
	err := s.redisClient.Close()
	if err != nil {
		s.logger.Error("Error closing Redis client", "error", err)
	}
}
