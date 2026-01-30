package initializer

import (
	"fmt"
	"go-pattern/internal/config"

	"github.com/dgraph-io/ristretto/v2"
)

func Ristretto[T any](config *config.LocalCacheConfig) (*ristretto.Cache[string, T], error) {
	if config == nil {
		return nil, fmt.Errorf("缓存配置不能为空")
	}
	// 构建Ristretto缓存
	cache, err := ristretto.NewCache(&ristretto.Config[string, T]{
		NumCounters: config.NumCounters,
		MaxCost:     config.MaxCost,
		BufferItems: config.BufferItems,
	})
	if err != nil {
		return nil, fmt.Errorf("创建Ristretto缓存失败: %w", err)
	}
	// 返回Ristretto缓存
	return cache, nil
}
