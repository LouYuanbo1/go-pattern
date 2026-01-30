package distributedCache

import (
	"context"
	"time"
)

type DistributedCache[T any] interface {
	SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration) error
	SetWithDefaultTTL(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (T, error)
	GetPointer(ctx context.Context, key string) (*T, error)
	Del(ctx context.Context, key string) error
}
