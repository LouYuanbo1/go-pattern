package localCache

import (
	"context"
	"go-pattern/internal/config"
	"go-pattern/internal/initializer"
	"log"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type ristrettoCache[T any] struct {
	cache      *ristretto.Cache[string, T]
	defaultTTL time.Duration
}

func NewRistrettoCache[T any](localCacheConfig *config.LocalCacheConfig) (LocalCache[T], error) {
	localCache, err := initializer.Ristretto[T](localCacheConfig)
	if err != nil {
		return nil, err
	}
	return &ristrettoCache[T]{cache: localCache, defaultTTL: time.Duration(localCacheConfig.DefaultTTL) * time.Second}, nil
}
func (r *ristrettoCache[T]) SetWithTTL(ctx context.Context, key string, value T, ttl time.Duration) bool {
	isSuccess := r.cache.SetWithTTL(key, value, 1, ttl)
	if !isSuccess {
		log.Printf("ristretto set drop key: %s", key)
		return false
	}
	return true
}

func (r *ristrettoCache[T]) SetWithDefaultTTL(ctx context.Context, key string, value T) bool {
	isSuccess := r.cache.SetWithTTL(key, value, 1, r.defaultTTL)
	if !isSuccess {
		log.Printf("ristretto set drop key: %s", key)
		return false
	}
	return true
}

func (r *ristrettoCache[T]) Get(ctx context.Context, key string) (T, bool) {
	value, isExist := r.cache.Get(key)
	if !isExist {
		log.Printf("ristretto get not exist key: %s", key)
		var zeroValue T
		return zeroValue, false
	}
	return value, true
}

func (r *ristrettoCache[T]) GetPointer(ctx context.Context, key string) (*T, bool) {
	value, isExist := r.cache.Get(key)
	if !isExist {
		log.Printf("ristretto get not exist key: %s", key)
		return nil, false
	}
	return &value, true
}

func (r *ristrettoCache[T]) Del(ctx context.Context, key string) {
	r.cache.Del(key)
}
