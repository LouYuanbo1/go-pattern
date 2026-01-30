package repo

import (
	"context"
	orderRepo "go-pattern/internal/repo/order"
	productRepo "go-pattern/internal/repo/product"
	userRepo "go-pattern/internal/repo/user"
)

type RepoFactory interface {
	Transaction(ctx context.Context, fn func(factory RepoFactory) error) error
	User() userRepo.UserRepo
	Order() orderRepo.OrderRepo
	Product() productRepo.ProductRepo
}
