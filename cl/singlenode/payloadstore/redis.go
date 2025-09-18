package payloadstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/primev/mev-commit/cl/types"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewRedisRepository(redisClient *redis.Client, logger *slog.Logger) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
		logger:      logger.With("component", "RedisRepository"),
	}
}

const zKeyPayloads = "execution_payloads:z"

func (r *RedisRepository) SavePayload(ctx context.Context, info *types.PayloadInfo) error {
	if info.InsertedAt.IsZero() {
		info.InsertedAt = time.Now().UTC()
	}
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	score := float64(info.BlockHeight)

	pipe := r.redisClient.TxPipeline()
	pipe.ZRemRangeByScore(ctx, // Remove existing payload at this height with min=max=height
		zKeyPayloads,
		strconv.FormatFloat(score, 'f', -1, 64), // min
		strconv.FormatFloat(score, 'f', -1, 64), // max
	)
	pipe.ZAdd(ctx, zKeyPayloads, redis.Z{
		Score:  score,
		Member: string(data),
	})
	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Error("Failed to save payload to Redis",
			"payload_id", info.PayloadID,
			"block_height", info.BlockHeight,
			"error", err,
		)
		return fmt.Errorf("save payload: %w", err)
	}

	r.logger.Debug("Payload saved to Redis",
		"payload_id", info.PayloadID,
		"block_height", info.BlockHeight,
	)
	return nil
}

func (r *RedisRepository) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	rangeBy := &redis.ZRangeBy{
		Min:    strconv.FormatUint(sinceHeight, 10),
		Max:    "+inf",
		Offset: 0,
		Count:  int64(limit),
	}
	members, err := r.redisClient.ZRangeByScore(ctx, zKeyPayloads, rangeBy).Result()
	if err != nil {
		return nil, fmt.Errorf("ZRangeByScore: %w", err)
	}

	result := make([]types.PayloadInfo, 0, len(members))
	for _, m := range members {
		var pi types.PayloadInfo
		if err := json.Unmarshal([]byte(m), &pi); err != nil {
			return nil, fmt.Errorf("unmarshal payload: %w", err)
		}
		result = append(result, pi)
	}

	r.logger.Debug("Retrieved payloads since height",
		"since_height", sinceHeight,
		"count", len(result),
		"limit", limit,
	)
	return result, nil
}

func (r *RedisRepository) GetLatestHeight(ctx context.Context) (uint64, error) {
	startIdx := int64(0)
	stopIdx := int64(0)
	items, err := r.redisClient.ZRevRangeWithScores(ctx, zKeyPayloads, startIdx, stopIdx).Result()
	if err != nil {
		return 0, fmt.Errorf("ZRevRangeWithScores: %w", err)
	}
	if len(items) == 0 {
		return 0, nil
	}
	if items[0].Score < 0 {
		return 0, fmt.Errorf("negative height score: %v", items[0].Score)
	}
	return uint64(items[0].Score), nil
}

func (r *RedisRepository) Close() error {
	r.logger.Info("Closing Redis client")
	return r.redisClient.Close()
}
