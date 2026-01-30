package cache

import (
	distributedCache "go-pattern/internal/cache/distributed"
	localCache "go-pattern/internal/cache/local"
	"go-pattern/internal/config"
	"go-pattern/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type multiLevelCacheFactory struct {
	redisClient *redis.Client
}

func NewMultiLevelCacheFactory(redisClient *redis.Client) *multiLevelCacheFactory {
	return &multiLevelCacheFactory{
		redisClient: redisClient,
	}
}

func (f *multiLevelCacheFactory) User(localCacheConfig *config.LocalCacheConfig, defaultTTLDistributedCache time.Duration) MultiLevelCache[model.User] {
	distributedCache := distributedCache.NewRedisCache[model.User](f.redisClient, defaultTTLDistributedCache)
	ristrettoCache, err := localCache.NewRistrettoCache[model.User](localCacheConfig)
	if err != nil {
		panic(err)
	}
	return NewMultiLevelCache(
		ristrettoCache,
		distributedCache,
	)
}

func (f *multiLevelCacheFactory) Order(localCacheConfig *config.LocalCacheConfig, defaultTTLDistributedCache time.Duration) MultiLevelCache[model.Order] {
	distributedCache := distributedCache.NewRedisCache[model.Order](f.redisClient, defaultTTLDistributedCache)
	ristrettoCache, err := localCache.NewRistrettoCache[model.Order](localCacheConfig)
	if err != nil {
		panic(err)
	}
	return NewMultiLevelCache(
		ristrettoCache,
		distributedCache,
	)
}

func (f *multiLevelCacheFactory) Product(localCacheConfig *config.LocalCacheConfig, defaultTTLDistributedCache time.Duration) MultiLevelCache[model.Product] {
	distributedCache := distributedCache.NewRedisCache[model.Product](f.redisClient, defaultTTLDistributedCache)
	ristrettoCache, err := localCache.NewRistrettoCache[model.Product](localCacheConfig)
	if err != nil {
		panic(err)
	}
	return NewMultiLevelCache(
		ristrettoCache,
		distributedCache,
	)
}
