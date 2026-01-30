package repo

import (
	"context"

	orderRepo "go-pattern/internal/repo/order"
	productRepo "go-pattern/internal/repo/product"
	userRepo "go-pattern/internal/repo/user"

	"gorm.io/gorm"
)

type repoFactory struct {
	db *gorm.DB
}

func NewRepoFactory(db *gorm.DB) *repoFactory {
	return &repoFactory{db: db}
}

func (f *repoFactory) withTransaction(tx *gorm.DB) *repoFactory {
	return &repoFactory{
		db: tx,
	}
}

// Transaction 执行事务操作
func (f *repoFactory) Transaction(ctx context.Context, fn func(factory RepoFactory) error) error {
	// 使用gorm事务,自动控制事务提交和回滚
	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建事务工厂
		txFactory := f.withTransaction(tx)
		// 执行用户逻辑
		return fn(txFactory)
	})
}

func (f *repoFactory) User() userRepo.UserRepo {
	// 这里可以添加缓存，避免重复创建
	return userRepo.NewUserRepo(f.db)
}

func (f *repoFactory) Order() orderRepo.OrderRepo {
	// 这里可以添加缓存，避免重复创建
	return orderRepo.NewOrderRepo(f.db)
}

func (f *repoFactory) Product() productRepo.ProductRepo {
	// 这里可以添加缓存，避免重复创建
	return productRepo.NewProductRepo(f.db)
}
