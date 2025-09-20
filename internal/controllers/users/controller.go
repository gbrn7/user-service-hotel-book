package users

import (
	"context"
	"user-service/domain/dto"

	"github.com/gin-gonic/gin"
)

type service interface {
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetAllAdmin(context.Context) ([]*dto.UserResponse, error)
	GetAllCustomer(context.Context) ([]*dto.UserResponse, error)
	GetAllUser(context.Context) ([]*dto.UserResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type UserController struct {
	*gin.Engine
	service service
}

func NewUserController(api *gin.Engine, service service) *UserController {
	return &UserController{
		Engine:  api,
		service: service,
	}
}
