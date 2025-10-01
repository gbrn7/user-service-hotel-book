package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/helpers"
	"user-service/internal/controllers/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserController_SignUp(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	type args struct {
		ctx context.Context
		req dto.RegisterRequest
	}

	newUUID := uuid.New()

	tests := []struct {
		name               string
		args               args
		mockFn             func(args args)
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Name:            "name",
					Password:        "password",
					ConfirmPassword: "password",
					RoleID:          1,
					Username:        "username",
					Email:           "email@gmail.com",
					PhoneNumber:     "phone_number",
				},
			},
			expectedStatusCode: http.StatusCreated,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: map[string]interface{}{
					"uuid":         newUUID.String(),
					"name":         "name",
					"username":     "username",
					"email":        "email@gmail.com",
					"phone_number": "phone_number",
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func(args args) {
				mockSvc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&dto.RegisterResponse{
					User: dto.UserResponse{
						UUID:        newUUID,
						Name:        args.req.Name,
						Username:    args.req.Username,
						Email:       args.req.Email,
						PhoneNumber: args.req.PhoneNumber,
					},
				}, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Name:            "name",
					Password:        "password",
					ConfirmPassword: "password",
					RoleID:          1,
					Username:        "username",
					Email:           "email@gmail.com",
					PhoneNumber:     "phone_number",
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			wantErr:            true,
			ExpectedBody: helpers.Response{
				Status:  constants.Error,
				Message: errConstants.ErrInternalServerError.Error(),
			},
			mockFn: func(args args) {
				mockSvc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/signup"

			userV1 := r.Group("/user/v1")
			userV1.POST("auth/signup", controller.SignUp)

			registerReq := dto.RegisterRequest{
				Name:            tt.args.req.Name,
				Username:        tt.args.req.Username,
				Password:        tt.args.req.Password,
				ConfirmPassword: tt.args.req.ConfirmPassword,
				Email:           tt.args.req.Email,
				PhoneNumber:     tt.args.req.PhoneNumber,
				RoleID:          tt.args.req.RoleID,
			}

			w := httptest.NewRecorder()

			val, err := json.Marshal(registerReq)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endpoint, body)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_SignIn(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	token := uuid.New().String()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: map[string]interface{}{
					"uuid":         newUUID.String(),
					"name":         "name",
					"username":     "username",
					"email":        "email@gmail.com",
					"phone_number": "phone_number",
				},
				Message: http.StatusText(http.StatusOK),
				Token:   &token,
			},
			mockFn: func() {
				mockSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&dto.LoginResponse{
					User: dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
					Token: token,
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/signin"

			r.POST(endpoint, controller.SignIn)

			loginReq := dto.LoginRequest{
				Email:    "email@gmail.com",
				Password: "password",
			}

			w := httptest.NewRecorder()

			val, err := json.Marshal(loginReq)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endpoint, body)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}
func TestUserController_GetUserLogin(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: map[string]interface{}{
					"uuid":         newUUID.String(),
					"name":         "name",
					"username":     "username",
					"email":        "email@gmail.com",
					"phone_number": "phone_number",
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func() {
				mockSvc.EXPECT().GetUserLogin(gomock.Any()).Return(&dto.UserResponse{
					UUID:        newUUID,
					Name:        "name",
					Username:    "username",
					Email:       "email@gmail.com",
					PhoneNumber: "phone_number",
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().GetUserLogin(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/user"

			r.GET(endpoint, controller.GetUserLogin)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_GetUserByUUID(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: map[string]interface{}{
					"uuid":         newUUID.String(),
					"name":         "name",
					"username":     "username",
					"email":        "email@gmail.com",
					"phone_number": "phone_number",
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func() {
				mockSvc.EXPECT().GetUserByUUID(gomock.Any(), gomock.Any()).Return(&dto.UserResponse{
					UUID:        newUUID,
					Name:        "name",
					Username:    "username",
					Email:       "email@gmail.com",
					PhoneNumber: "phone_number",
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().GetUserByUUID(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/1"

			r.GET("/user/v1/auth/:uuid", controller.GetUserByUUID)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_GetAllCustomer(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: []interface{}{
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func() {
				mockSvc.EXPECT().GetAllCustomer(gomock.Any()).Return([]*dto.UserResponse{
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().GetAllCustomer(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/cust"

			r.GET(endpoint, controller.GetAllCustomer)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_GetAllAdmin(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: []interface{}{
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func() {
				mockSvc.EXPECT().GetAllAdmin(gomock.Any()).Return([]*dto.UserResponse{
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().GetAllAdmin(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/admin"

			r.GET(endpoint, controller.GetAllAdmin)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_GetAllUser(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	newUUID := uuid.New()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: []interface{}{
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
					map[string]interface{}{
						"uuid":         newUUID.String(),
						"name":         "name",
						"username":     "username",
						"email":        "email@gmail.com",
						"phone_number": "phone_number",
					},
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func() {
				mockSvc.EXPECT().GetAllUser(gomock.Any()).Return([]*dto.UserResponse{
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
					&dto.UserResponse{
						UUID:        newUUID,
						Name:        "name",
						Username:    "username",
						Email:       "email@gmail.com",
						PhoneNumber: "phone_number",
					},
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().GetAllUser(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/users"

			r.GET(endpoint, controller.GetAllUser)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}

func TestUserController_Update(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := users.NewMockservice(ctrlMock)

	type args struct {
		ctx context.Context
		req dto.UpdateRequest
	}

	newUUID := uuid.New()

	tests := []struct {
		name               string
		args               args
		mockFn             func(args args)
		expectedStatusCode int
		ExpectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: dto.UpdateRequest{
					Name:        "name",
					RoleID:      1,
					Username:    "username",
					Email:       "email@gmail.com",
					PhoneNumber: "phone_number",
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			ExpectedBody: helpers.Response{
				Status: constants.Success,
				Data: map[string]interface{}{
					"uuid":         newUUID.String(),
					"name":         "name",
					"username":     "username",
					"email":        "email@gmail.com",
					"phone_number": "phone_number",
				},
				Message: http.StatusText(http.StatusOK),
			},
			mockFn: func(args args) {
				mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&dto.UserResponse{
					UUID:        newUUID,
					Name:        args.req.Name,
					Username:    args.req.Username,
					Email:       args.req.Email,
					PhoneNumber: args.req.PhoneNumber,
				}, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				req: dto.UpdateRequest{
					Name:        "name",
					RoleID:      1,
					Username:    "username",
					Email:       "email@gmail.com",
					PhoneNumber: "phone_number",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			ExpectedBody: helpers.Response{
				Status:  constants.Error,
				Message: errConstants.ErrInternalServerError.Error(),
			},
			mockFn: func(args args) {
				mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := gin.New()
			controller := users.NewUserController(r, mockSvc)
			endpoint := "/user/v1/auth/update"

			userV1 := r.Group("/user/v1")
			userV1.POST("auth/update", controller.Update)

			updateReq := dto.UpdateRequest{
				Name:        tt.args.req.Name,
				Username:    tt.args.req.Username,
				Email:       tt.args.req.Email,
				PhoneNumber: tt.args.req.PhoneNumber,
				RoleID:      tt.args.req.RoleID,
			}

			w := httptest.NewRecorder()

			val, err := json.Marshal(updateReq)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endpoint, body)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.ExpectedBody, response)
			}

		})
	}
}
