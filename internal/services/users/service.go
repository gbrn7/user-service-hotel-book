package users

import (
	"context"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/helpers/configs"
)

type repository interface {
	Create(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	GetAllUser(context.Context) (*[]models.User, error)
	GetAllAdmin(context.Context) (*[]models.User, error)
	GetAllCustomer(context.Context) (*[]models.User, error)
}

type service struct {
	cfg        *configs.Config
	repository repository
}

func NewUserService(cfg *configs.Config, repository repository) *service {
	return &service{
		cfg:        cfg,
		repository: repository,
	}
}
