package users

import (
	"context"
	"errors"
	"user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *repository) Create(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      req.RoleID,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Password:    *req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      req.RoleID,
	}

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Updates(&user).
		Error

	if err != nil {
		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &user, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("uuid = ?", uuid).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &user, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &user, nil
}

func (r *repository) GetAllUser(ctx context.Context) (*[]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Preload("Role").
		Find(&users).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &users, nil
}

func (r *repository) GetAllAdmin(ctx context.Context) (*[]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Preload("Role", "id = ?", constants.Admin).
		Find(&users).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &users, nil
}

func (r *repository) GetAllCustomer(ctx context.Context) (*[]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Preload("Role", "id = ?", constants.Customer).
		Find(&users).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.WrapError(errConstants.ErrUserNotFound)
		}

		return nil, helpers.WrapError(errConstants.ErrSQLError)
	}

	return &users, nil
}
