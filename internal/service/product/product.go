package service

import (
	"context"
	"go-pattern/internal/model"
	repo "go-pattern/internal/repo/factory"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	CreateProducts(ctx context.Context, products []*model.Product, batchSize int) error
	GetProduct(ctx context.Context, id uint64) (*model.Product, error)
	GetProducts(ctx context.Context, ids []uint64) ([]*model.Product, error)
	GetProductsByPage(ctx context.Context, page, pageSize uint64) ([]*model.Product, error)
	GetProductsByCursor(ctx context.Context, cursor, pageSize uint64) ([]*model.Product, uint64, bool, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id uint64) error
	DeleteProducts(ctx context.Context, ids []uint64) error
	ReduceQuantity(ctx context.Context, productID, count uint64) error
}

type productService struct {
	repoFactory repo.RepoFactory
}

func NewProductService(repoFactory repo.RepoFactory) ProductService {
	return &productService{repoFactory: repoFactory}
}

func (p *productService) CreateProduct(ctx context.Context, Product *model.Product) error {
	return p.repoFactory.Product().Create(ctx, Product)
}

func (p *productService) CreateProducts(ctx context.Context, Products []*model.Product, batchSize int) error {
	return p.repoFactory.Product().CreateInBatches(ctx, Products, batchSize)
}

func (p *productService) GetProduct(ctx context.Context, id uint64) (*model.Product, error) {
	return p.repoFactory.Product().GetByID(ctx, id)
}

func (p *productService) GetProducts(ctx context.Context, ids []uint64) ([]*model.Product, error) {
	return p.repoFactory.Product().GetByIDs(ctx, ids)
}

func (p *productService) GetProductsByPage(ctx context.Context, page, pageSize uint64) ([]*model.Product, error) {
	return p.repoFactory.Product().GetByPage(ctx, page, pageSize)
}

func (p *productService) GetProductsByCursor(ctx context.Context, cursor, pageSize uint64) ([]*model.Product, uint64, bool, error) {
	return p.repoFactory.Product().GetByCursor(ctx, cursor, pageSize)
}

func (p *productService) UpdateProduct(ctx context.Context, Product *model.Product) error {
	return p.repoFactory.Product().Update(ctx, Product)
}

func (p *productService) DeleteProduct(ctx context.Context, id uint64) error {
	return p.repoFactory.Product().DeleteByID(ctx, id)
}

func (p *productService) DeleteProducts(ctx context.Context, ids []uint64) error {
	return p.repoFactory.Product().DeleteByIDs(ctx, ids)
}

func (p *productService) ReduceQuantity(ctx context.Context, productID, count uint64) error {
	return p.repoFactory.Product().ReduceQuantity(ctx, productID, count)
}
