package cache

import (
	"context"
	"fmt"
	distributedCache "go-pattern/internal/cache/distributed"
	localCache "go-pattern/internal/cache/local"
	"go-pattern/internal/model"
	"log"
	"time"
)

type MultiLevelCache[T any] interface {
	SetWithTTL(ctx context.Context, key string, value T, l1Expiration time.Duration, l2Expiration time.Duration) error
	SetWithDefaultTTL(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (T, error)
	GetPointer(ctx context.Context, key string) (*T, error)
	Del(ctx context.Context, key string) error
}

type multiLevelCache[T any] struct {
	localCache       localCache.LocalCache[T]
	distributedCache distributedCache.DistributedCache[T]
}

func NewMultiLevelCache[T any, PT model.PointerModel[T]](
	localCache localCache.LocalCache[T],
	distributedCache distributedCache.DistributedCache[T],
) MultiLevelCache[T] {
	return &multiLevelCache[T]{
		localCache:       localCache,
		distributedCache: distributedCache,
	}
}

func (m *multiLevelCache[T]) SetWithTTL(ctx context.Context, key string, value T, l1Expiration time.Duration, l2Expiration time.Duration) error {
	err := m.distributedCache.SetWithTTL(ctx, key, value, l2Expiration)
	if err != nil {
		return fmt.Errorf("distributed cache set failed: %w", err)
	}
	isSuccess := m.localCache.SetWithTTL(ctx, key, value, l1Expiration)
	if !isSuccess {
		log.Printf("warning: local cache set failed, key: %s", key)
	}
	return nil
}

func (m *multiLevelCache[T]) SetWithDefaultTTL(ctx context.Context, key string, value T) error {
	err := m.distributedCache.SetWithDefaultTTL(ctx, key, value)
	if err != nil {
		return fmt.Errorf("distributed cache set failed: %w", err)
	}
	isSuccess := m.localCache.SetWithDefaultTTL(ctx, key, value)
	if !isSuccess {
		log.Printf("warning: local cache set failed, key: %s", key)
	}
	return nil
}

func (m *multiLevelCache[T]) Get(ctx context.Context, key string) (T, error) {
	value, isExist := m.localCache.Get(ctx, key)
	if isExist {
		return value, nil
	}
	log.Printf("local cache get failed, key: %s", key)
	value, err := m.distributedCache.Get(ctx, key)
	if err != nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("distributed cache get failed: %w", err)
	}
	isSuccess := m.localCache.SetWithDefaultTTL(ctx, key, value)
	if !isSuccess {
		log.Printf("warning: local cache set failed, key: %s", key)
	}
	return value, nil
}

func (m *multiLevelCache[T]) GetPointer(ctx context.Context, key string) (*T, error) {
	value, isExist := m.localCache.GetPointer(ctx, key)
	if isExist {
		return value, nil
	}
	log.Printf("local cache get failed, key: %s", key)
	value, err := m.distributedCache.GetPointer(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("distributed cache get failed: %w", err)
	}
	return value, nil
}

func (m *multiLevelCache[T]) Del(ctx context.Context, key string) error {
	err := m.distributedCache.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("distributed cache del failed: %w", err)
	}
	m.localCache.Del(ctx, key)
	return nil
}
