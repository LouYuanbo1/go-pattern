package repo

import (
	"go-pattern/internal/model"
	genericRepo "go-pattern/internal/repo/generic"

	"gorm.io/gorm"
)

type UserRepo interface {
	genericRepo.GenericRepo[model.User, *model.User]
}

type userRepo struct {
	genericRepo.GenericRepo[model.User, *model.User]
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		GenericRepo: genericRepo.NewGenericRepo[model.User](db),
	}
}
