package repo

import (
	"context"
	"fmt"
	"go-pattern/internal/model"

	genericRepo "go-pattern/internal/repo/generic"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo interface {
	genericRepo.GenericRepo[model.Product, *model.Product]
	ReduceQuantity(ctx context.Context, productID, count uint64) error
}

type productRepo struct {
	genericRepo.GenericRepo[model.Product, *model.Product]
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{
		GenericRepo: genericRepo.NewGenericRepo[model.Product](db),
		db:          db,
	}
}

func (p *productRepo) ReduceQuantity(ctx context.Context, productID, count uint64) error {
	result := p.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Product{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND quantity >= ?", productID, count).
			Update("quantity", gorm.Expr("quantity - ?", count)).Error
	})
	if result != nil {
		return fmt.Errorf("reduce quantity failed: %w", result)
	}
	return nil
}
