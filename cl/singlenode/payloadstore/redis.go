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

func NewRedisRepositoryFromURL(ctx context.Context, redisURL string, logger *slog.Logger) (*RedisRepository, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return NewRedisRepository(rdb, logger), nil
}

const (
	zKeyPayloads    = "execution_payloads:z"
	keyLatestHeight = "execution_payloads:latest_height"
)

var savePayloadScript = redis.NewScript(`
-- KEYS[1] = zset key
-- KEYS[2] = latest height string key
-- ARGV[1] = score (height)
-- ARGV[2] = member (JSON string)

local zkey = KEYS[1]
local skey = KEYS[2]
local score = tonumber(ARGV[1])
local member = ARGV[2]

redis.call('ZREMRANGEBYSCORE', zkey, score, score)
redis.call('ZADD', zkey, score, member)

local cur = redis.call('GET', skey)
if not cur or score > tonumber(cur) then
  redis.call('SET', skey, tostring(score))
end

return 1
`)

func (r *RedisRepository) SavePayload(ctx context.Context, info *types.PayloadInfo) error {
	if info.InsertedAt.IsZero() {
		info.InsertedAt = time.Now().UTC()
	}
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	score := float64(info.BlockHeight)

	res := savePayloadScript.Run(ctx, r.redisClient, []string{zKeyPayloads, keyLatestHeight}, score, string(data))
	if err := res.Err(); err != nil {
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

func (r *RedisRepository) GetPayloadByHeight(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
	hStr := strconv.FormatUint(height, 10)
	members, err := r.redisClient.ZRangeByScore(ctx, zKeyPayloads, &redis.ZRangeBy{
		Min:    hStr,
		Max:    hStr,
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("ZRangeByScore: %w", err)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("payload not found")
	}
	var pi types.PayloadInfo
	if err := json.Unmarshal([]byte(members[0]), &pi); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}
	return &pi, nil
}

func (r *RedisRepository) GetLatestPayload(ctx context.Context) (*types.PayloadInfo, error) {
	h, err := r.GetLatestHeight(ctx)
	if err != nil {
		return nil, err
	}
	if h == 0 {
		return nil, nil
	}
	return r.GetPayloadByHeight(ctx, h)
}

func (r *RedisRepository) GetLatestHeight(ctx context.Context) (uint64, error) {
	s, err := r.redisClient.Get(ctx, keyLatestHeight).Result()
	if err == redis.Nil || s == "" {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("GET latest_height: %w", err)
	}
	h, perr := strconv.ParseUint(s, 10, 64)
	if perr != nil {
		return 0, nil
	}
	return h, nil
}

func (r *RedisRepository) Close() error {
	r.logger.Info("Closing Redis client")
	return r.redisClient.Close()
}
