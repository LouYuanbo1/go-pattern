package lock

import (
	"context"
	"time"
)

type Lock interface {
	Acquire(ctx context.Context, key string, expiration time.Duration) (string, bool, error)
	Release(ctx context.Context, key string, lockID string) error
}
