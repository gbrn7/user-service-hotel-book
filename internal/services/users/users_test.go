package users

import (
	"context"
	"testing"
	"time"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/helpers/configs"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
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
