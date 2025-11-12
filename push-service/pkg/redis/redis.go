package redis

import (
	"context"
	"encoding/json"
	"push-service/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	// Use the GetRedisURL method
	redisURL := cfg.GetRedisURL()

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	zap.L().Info("Connected to Redis", zap.String("url", redisURL))
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Queue operations
func (r *RedisClient) Enqueue(ctx context.Context, queueName string, message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.Client.LPush(ctx, queueName, jsonMessage).Err()
}

func (r *RedisClient) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (string, error) {
	result, err := r.Client.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		return "", err
	}

	if len(result) < 2 {
		return "", redis.Nil
	}

	return result[1], nil
}

func (r *RedisClient) QueueLength(ctx context.Context, queueName string) (int64, error) {
	return r.Client.LLen(ctx, queueName).Result()
}
