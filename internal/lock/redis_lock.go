package lock

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisLock struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) Lock {
	return &redisLock{
		client: client,
	}
}

func (r *redisLock) Acquire(ctx context.Context, key string, expiration time.Duration) (string, bool, error) {
	lockID := uuid.New().String()
	success, err := r.client.SetNX(ctx, key, lockID, expiration).Result()
	if err != nil {
		return "", false, err
	}
	return lockID, success, nil
}

func (r *redisLock) Release(ctx context.Context, key string, lockID string) error {
	luaScript := `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
    `
	script := redis.NewScript(luaScript)
	_, err := script.Run(ctx, r.client, []string{key}, lockID).Result()
	if err != nil {
		log.Printf("redis unlock error: %v", err)
		return fmt.Errorf("redis unlock error: %w", err)
	}
	return nil
}
