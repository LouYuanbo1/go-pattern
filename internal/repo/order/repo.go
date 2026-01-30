package repo

import (
	"go-pattern/internal/model"

	genericRepo "go-pattern/internal/repo/generic"

	"gorm.io/gorm"
)

type OrderRepo interface {
	genericRepo.GenericRepo[model.Order, *model.Order]
}

type orderRepo struct {
	genericRepo.GenericRepo[model.Order, *model.Order]
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{
		GenericRepo: genericRepo.NewGenericRepo[model.Order](db),
	}
}
