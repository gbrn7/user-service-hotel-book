package users

import (
	"context"
	"testing"
	"time"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/helpers/configs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_service_Register(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		request dto.RegisterRequest
	}

	now := time.Now()

	password := "password"

	tests := []struct {
		name    string // description of this test case
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				request: dto.RegisterRequest{
					Name:            "name",
					Username:        "username",
					Password:        password,
					ConfirmPassword: password,
					Email:           "email@gmail.com",
					PhoneNumber:     "0898989",
					RoleID:          2,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUsername(context.Background(), args.request.Username).Return(nil, nil)
				mockRepo.EXPECT().FindByEmail(context.Background(), args.request.Email).Return(nil, nil)
				mockRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(&models.User{
					ID:          1,
					UUID:        uuid.New(),
					Name:        args.request.Name,
					Username:    args.request.Username,
					Password:    args.request.Password,
					Email:       args.request.Email,
					RoleID:      args.request.RoleID,
					PhoneNumber: args.request.PhoneNumber,
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        2,
						Code:      "CUSTOMER",
						Name:      "Customeer",
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
		},
		{
			name: "fail data username is existed",
			args: args{
				request: dto.RegisterRequest{
					Name:            "name",
					Username:        "username",
					Password:        password,
					ConfirmPassword: password,
					Email:           "email@gmail.com",
					PhoneNumber:     "0898989",
					RoleID:          2,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUsername(context.Background(), args.request.Username).Return(&models.User{
					Name:        args.request.Name,
					Username:    args.request.Username,
					Password:    args.request.Password,
					Email:       args.request.Email,
					RoleID:      args.request.RoleID,
					PhoneNumber: args.request.PhoneNumber,
				}, nil)
			},
		},
		{
			name: "fail data email is existed",
			args: args{
				request: dto.RegisterRequest{
					Name:            "name",
					Username:        "username",
					Password:        password,
					ConfirmPassword: password,
					Email:           "email@gmail.com",
					PhoneNumber:     "0898989",
					RoleID:          2,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUsername(context.Background(), args.request.Username).Return(nil, nil)
				mockRepo.EXPECT().FindByEmail(context.Background(), args.request.Email).Return(&models.User{
					Name:        args.request.Name,
					Username:    args.request.Username,
					Password:    args.request.Password,
					Email:       args.request.Email,
					RoleID:      args.request.RoleID,
					PhoneNumber: args.request.PhoneNumber,
				}, nil)
			},
		},
		{
			name: "fail password not match",
			args: args{
				request: dto.RegisterRequest{
					Name:            "name",
					Username:        "username",
					Password:        password,
					ConfirmPassword: "pass",
					Email:           "email@gmail.com",
					PhoneNumber:     "0898989",
					RoleID:          2,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUsername(context.Background(), args.request.Username).Return(nil, nil)
				mockRepo.EXPECT().FindByEmail(context.Background(), args.request.Email).Return(nil, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.Register(context.Background(), &tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		request dto.LoginRequest
	}

	password := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	assert.NoError(t, err)

	tests := []struct {
		name    string // description of this test case
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				request: dto.LoginRequest{
					Email:    "email@gmail.com",
					Password: password,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), args.request.Email).Return(&models.User{
					Name:        "name",
					Username:    "username",
					Password:    string(hashedPassword),
					Email:       "email",
					RoleID:      1,
					PhoneNumber: "08213243",
				}, nil)
			},
		},
		{
			name: "fail data not found",
			args: args{
				request: dto.LoginRequest{
					Email:    "email@gmail.com",
					Password: password,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), args.request.Email).Return(nil, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.Login(context.Background(), &tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_get_user_by_uuid(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		uuid uuid.UUID
	}

	tests := []struct {
		name    string // description of this test case
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				uuid: uuid.New(),
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(&models.User{
					Name:        "name",
					Username:    "username",
					Email:       "email",
					RoleID:      1,
					PhoneNumber: "08213243",
				}, nil)
			},
		},
		{
			name: "fail data not found",
			args: args{
				uuid: uuid.New(),
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.GetUserByUUID(context.Background(), tt.args.uuid.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("service.getUserByUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_get_all_admin(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	tests := []struct {
		name    string // description of this test case
		wantErr bool
		mockFn  func()
	}{
		{
			name:    "success",
			wantErr: false,
			mockFn: func() {
				mockRepo.EXPECT().GetAllAdmin(gomock.Any()).Return(&[]models.User{
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
				}, nil)
			},
		},
		{
			name:    "fail data not found",
			wantErr: true,
			mockFn: func() {
				mockRepo.EXPECT().GetAllAdmin(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.GetAllAdmin(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("service.getAllAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_get_all_customer(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	tests := []struct {
		name    string // description of this test case
		wantErr bool
		mockFn  func()
	}{
		{
			name:    "success",
			wantErr: false,
			mockFn: func() {
				mockRepo.EXPECT().GetAllCustomer(gomock.Any()).Return(&[]models.User{
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
				}, nil)
			},
		},
		{
			name:    "fail data not found",
			wantErr: true,
			mockFn: func() {
				mockRepo.EXPECT().GetAllCustomer(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.GetAllCustomer(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("service.getAllCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_get_all_user(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	tests := []struct {
		name    string // description of this test case
		wantErr bool
		mockFn  func()
	}{
		{
			name:    "success",
			wantErr: false,
			mockFn: func() {
				mockRepo.EXPECT().GetAllUser(gomock.Any()).Return(&[]models.User{
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
					models.User{
						Name:        "name",
						Username:    "username",
						Email:       "email",
						RoleID:      1,
						PhoneNumber: "08213243",
					},
				}, nil)
			},
		},
		{
			name:    "fail data not found",
			wantErr: true,
			mockFn: func() {
				mockRepo.EXPECT().GetAllUser(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.GetAllUser(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetAllUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_update(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		uuid    uuid.UUID
		request dto.UpdateRequest
	}

	tests := []struct {
		name    string // description of this test case
		wantErr bool
		args    args
		mockFn  func(args args)
	}{
		{
			name:    "success",
			wantErr: false,
			args: args{
				uuid: uuid.New(),
				request: dto.UpdateRequest{
					Name:     "haloo",
					Username: "updateusername",
					Email:    "email@gmail.com",
				}},
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(&models.User{
					Name:     args.request.Name,
					Username: args.request.Username,
					Email:    args.request.Email,
				}, nil)

				mockRepo.EXPECT().FindByUsername(gomock.Any(), args.request.Username).Return(nil, nil)

				mockRepo.EXPECT().FindByEmail(gomock.Any(), args.request.Email).Return(nil, nil)

				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), args.uuid.String()).Return(&models.User{
					UUID:        args.uuid,
					Name:        "name",
					Username:    "username",
					Email:       "email",
					RoleID:      1,
					PhoneNumber: "08213243",
				}, nil)
			},
		},
		{
			name:    "fail fetch data by uuid not found",
			wantErr: true,
			args: args{
				uuid: uuid.New(),
				request: dto.UpdateRequest{
					Name:     "haloo",
					Username: "updateusername",
					Email:    "email@gmail.com",
				},
			},
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(nil, assert.AnError)
			},
		},
		{
			name:    "fail fetch data by username not found",
			wantErr: true,
			args: args{
				uuid: uuid.New(),
				request: dto.UpdateRequest{
					Name:     "haloo",
					Username: "updateusername",
					Email:    "email@gmail.com",
				},
			},
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(&models.User{
					Name:     args.request.Name,
					Username: args.request.Username,
					Email:    args.request.Email,
				}, nil)

				mockRepo.EXPECT().FindByUsername(gomock.Any(), args.request.Username).Return(nil, assert.AnError)
			},
		},
		{
			name:    "fail fetch data by email not found",
			wantErr: true,
			args: args{
				uuid: uuid.New(),
				request: dto.UpdateRequest{
					Name:     "haloo",
					Username: "updateusername",
					Email:    "email@gmail.com",
				},
			},
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(&models.User{
					Name:     args.request.Name,
					Username: args.request.Username,
					Email:    args.request.Email,
				}, nil)

				mockRepo.EXPECT().FindByUsername(gomock.Any(), args.request.Username).Return(nil, nil)

				mockRepo.EXPECT().FindByEmail(gomock.Any(), args.request.Email).Return(nil, assert.AnError)
			},
		},
		{
			name:    "fail update",
			wantErr: true,
			args: args{
				uuid: uuid.New(),
				request: dto.UpdateRequest{
					Name:     "haloo",
					Username: "updateusername",
					Email:    "email@gmail.com",
				}},
			mockFn: func(args args) {
				mockRepo.EXPECT().FindByUUID(gomock.Any(), args.uuid.String()).Return(&models.User{
					Name:     args.request.Name,
					Username: args.request.Username,
					Email:    args.request.Email,
				}, nil)

				mockRepo.EXPECT().FindByUsername(gomock.Any(), args.request.Username).Return(nil, nil)

				mockRepo.EXPECT().FindByEmail(gomock.Any(), args.request.Email).Return(nil, nil)

				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), args.uuid.String()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			_, err := s.Update(context.Background(), &tt.args.request, tt.args.uuid.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
