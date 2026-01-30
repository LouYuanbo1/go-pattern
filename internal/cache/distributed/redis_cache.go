package distributedCache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache[T any] struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewRedisCache[T any](client *redis.Client, defaultTTL time.Duration) DistributedCache[T] {
	return &redisCache[T]{client: client, defaultTTL: defaultTTL}
}

func (r *redisCache[T]) SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return fmt.Errorf("json marshal error: %w", err)
	}
	err = r.client.Set(ctx, key, jsonValue, ttl).Err()
	if err != nil {
		log.Printf("redis set error: %v", err)
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

func (r *redisCache[T]) SetWithDefaultTTL(ctx context.Context, key string, value T) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return fmt.Errorf("json marshal error: %w", err)
	}
	err = r.client.Set(ctx, key, jsonValue, r.defaultTTL).Err()
	if err != nil {
		log.Printf("redis set error: %v", err)
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

func (r *redisCache[T]) Get(ctx context.Context, key string) (T, error) {
	var result T
	jsonValue, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("redis get error: %v", err)
		return result, fmt.Errorf("redis get error: %w", err)
	}
	err = json.Unmarshal(jsonValue, &result)
	if err != nil {
		log.Printf("json unmarshal error: %v", err)
		return result, fmt.Errorf("json unmarshal error: %w", err)
	}
	return result, nil
}

func (r *redisCache[T]) GetPointer(ctx context.Context, key string) (*T, error) {
	var result T
	jsonValue, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("redis get error: %v", err)
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	err = json.Unmarshal(jsonValue, &result)
	if err != nil {
		log.Printf("json unmarshal error: %v", err)
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}
	return &result, nil
}

func (r *redisCache[T]) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("redis del error: %v", err)
		return fmt.Errorf("redis del error: %w", err)
	}
	return nil
}
