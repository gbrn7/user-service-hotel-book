package users

import (
	"context"
	"strings"
	"time"
	"user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/helpers"

	"go.elastic.co/apm/v2"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	span, spanctx := apm.StartSpan(ctx, "UserService.Register", "service")
	defer span.End()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindByUsername(spanctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errConstants.ErrUsernameExist
	}

	user, err = s.repository.FindByEmail(spanctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errConstants.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstants.ErrPasswordDoesNotMatch
	}

	user, err = s.repository.Create(spanctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashPassword),
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return response, nil
}

func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	span, spanctx := apm.StartSpan(ctx, "UserService.Login", "service")
	defer span.End()

	user, err := s.repository.FindByEmail(spanctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
		PhoneNumber: user.PhoneNumber,
	}

	expirationTime := time.Now().Add(time.Duration(s.cfg.JwtConfig.JwtExpirationTime) * time.Minute).Unix()

	token, err := helpers.GenerateToken(ctx, data, expirationTime)
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: token,
	}

	return response, nil
}

func (s *service) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	span, _ := apm.StartSpan(ctx, "UserService.GetUserLogin", "service")
	defer span.End()

	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
		PhoneNumber: userLogin.PhoneNumber,
	}

	return &data, nil
}

func (s *service) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "UserService.GetUserByUUID", "service")
	defer span.End()

	user, err := s.repository.FindByUUID(spanCtx, uuid)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	return data, nil
}

func (s *service) GetAllAdmin(ctx context.Context) ([]*dto.UserResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "UserService.GetAllAdmin", "service")
	defer span.End()

	users, err := s.repository.GetAllAdmin(spanCtx)
	if err != nil {
		return nil, err
	}

	var usersSlice []*dto.UserResponse

	for _, user := range *users {
		usersSlice = append(usersSlice, &dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		})
	}
	return usersSlice, nil
}

func (s *service) GetAllCustomer(ctx context.Context) ([]*dto.UserResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "UserService.GetAllCustomer", "service")
	defer span.End()

	users, err := s.repository.GetAllCustomer(spanCtx)
	if err != nil {
		return nil, err
	}

	var usersSlice []*dto.UserResponse

	for _, user := range *users {
		usersSlice = append(usersSlice, &dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		})
	}
	return usersSlice, nil
}

func (s *service) GetAllUser(ctx context.Context) ([]*dto.UserResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "UserService.GetAllUser", "service")
	defer span.End()

	users, err := s.repository.GetAllUser(spanCtx)
	if err != nil {
		return nil, err
	}

	var usersSlice []*dto.UserResponse

	for _, user := range *users {
		usersSlice = append(usersSlice, &dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role.Name,
		})
	}
	return usersSlice, nil
}

func (s *service) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "UserService.Update", "service")
	defer span.End()

	var password string
	user, err := s.GetUserByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	userByUsername, err := s.repository.FindByUsername(spanCtx, req.Username)
	if err != nil {
		return nil, err
	}

	if userByUsername != nil && user.Username != req.Username {
		return nil, errConstants.ErrUsernameExist
	}

	userByEmail, err := s.repository.FindByEmail(spanCtx, req.Email)
	if err != nil {
		return nil, err
	}

	if userByEmail != nil && user.Email != req.Email {
		return nil, errConstants.ErrEmailExist
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConstants.ErrPasswordDoesNotMatch
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		password = string(hashPassword)
	}

	userResult, err := s.repository.Update(spanCtx, &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    &password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	}, uuid)

	if err != nil {
		return nil, err
	}

	data := &dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
	}

	return data, nil
}
