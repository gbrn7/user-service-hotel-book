package users

import (
	"context"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_repository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		model dto.RegisterRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				model: dto.RegisterRequest{
					Email:       "test@gmail.com",
					Name:        "testUser",
					Username:    "usertest",
					Password:    "password",
					PhoneNumber: "0988338",
					RoleID:      1,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					args.model.Name,
					args.model.Username,
					args.model.Password,
					args.model.Email,
					args.model.RoleID,
					args.model.PhoneNumber,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				model: dto.RegisterRequest{
					Email:           "test@gmail.com",
					Name:            "testUser",
					Username:        "usertest",
					Password:        "password",
					PhoneNumber:     "0988338",
					RoleID:          1,
					ConfirmPassword: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					args.model.Name,
					args.model.Username,
					args.model.Password,
					args.model.Email,
					args.model.RoleID,
					args.model.PhoneNumber,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				db: gormDB,
			}
			if _, err := r.Create(context.Background(), &tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("repository.Create() error %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		model dto.UpdateRequest
		uuid  uuid.UUID
	}

	password := "password"

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				uuid: uuid.New(),
				model: dto.UpdateRequest{
					Email:           "update@gmail.com",
					Name:            "update",
					Username:        "update",
					Password:        &password,
					ConfirmPassword: &password,
					PhoneNumber:     "0988338",
					RoleID:          1,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "users" (.+)`).
					WithArgs(
						args.model.Name,
						args.model.Username,
						args.model.Password,
						args.model.Email,
						args.model.RoleID,
						args.model.PhoneNumber,
						sqlmock.AnyArg(),
						args.uuid.String(),
					).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				uuid: uuid.New(),
				model: dto.UpdateRequest{
					Email:           "test@gmail.com",
					Name:            "testUser",
					Username:        "usertest",
					Password:        &password,
					PhoneNumber:     "0988338",
					RoleID:          1,
					ConfirmPassword: &password,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "users" (.+)`).
					WithArgs(
						args.model.Name,
						args.model.Username,
						args.model.Password,
						args.model.Email,
						args.model.RoleID,
						args.model.PhoneNumber,
						sqlmock.AnyArg(),
						args.uuid.String(),
					).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				db: gormDB,
			}
			if _, err := r.Update(context.Background(), &tt.args.model, tt.args.uuid.String()); (err != nil) != tt.wantErr {
				t.Errorf("repository.Update() error %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_FindByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		uuid uuid.UUID
	}

	uuid := uuid.New()

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *models.User
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				uuid: uuid,
			},
			want: &models.User{
				ID:          1,
				Name:        "name",
				Username:    "username",
				Password:    "password",
				Email:       "email@gmail.com",
				RoleID:      1,
				PhoneNumber: "08332323",
				CreatedAt:   now,
				UpdatedAt:   now,
				Role: models.Role{
					ID:        1,
					Code:      "ADMIN",
					Name:      "Administrator",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						uuid,
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRow(1, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1`).
					WithArgs(
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "ADMIN", "Administrator", now, now))
			},
		},
		{
			name: "error",
			want: nil,
			args: args{
				uuid: uuid,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						uuid,
						1,
					).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				db: gormDB,
			}
			got, err := r.FindByUUID(context.Background(), uuid.String())

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindByUUID() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindByUUID() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		email string
	}

	email := "email@gmail.com"

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *models.User
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				email: email,
			},
			want: &models.User{
				ID:          1,
				Name:        "name",
				Username:    "username",
				Password:    "password",
				Email:       "email@gmail.com",
				RoleID:      1,
				PhoneNumber: "08332323",
				CreatedAt:   now,
				UpdatedAt:   now,
				Role: models.Role{
					ID:        1,
					Code:      "ADMIN",
					Name:      "Administrator",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						email,
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRow(1, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1`).
					WithArgs(
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "ADMIN", "Administrator", now, now))
			},
		},
		{
			name: "error",
			want: nil,
			args: args{
				email: email,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						email,
						1,
					).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				db: gormDB,
			}
			got, err := r.FindByEmail(context.Background(), email)

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindByEmail() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindByEmail() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		username string
	}

	username := "username@gmail.com"

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *models.User
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				username: username,
			},
			want: &models.User{
				ID:          1,
				Name:        "name",
				Username:    "username",
				Password:    "password",
				Email:       "email@gmail.com",
				RoleID:      1,
				PhoneNumber: "08332323",
				CreatedAt:   now,
				UpdatedAt:   now,
				Role: models.Role{
					ID:        1,
					Code:      "ADMIN",
					Name:      "Administrator",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						username,
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRow(1, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1`).
					WithArgs(
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "ADMIN", "Administrator", now, now))
			},
		},
		{
			name: "error",
			want: nil,
			args: args{
				username: username,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
					WithArgs(
						username,
						1,
					).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				db: gormDB,
			}
			got, err := r.FindByUsername(context.Background(), username)

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindByUsername() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindByUsername() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetAllUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name    string
		wantErr bool
		want    *[]models.User
		mockFn  func()
	}{
		{
			name: "success",
			want: &[]models.User{
				models.User{
					ID:          1,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "ADMIN",
						Name:      "Administrator",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				models.User{
					ID:          3,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "ADMIN",
						Name:      "Administrator",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			wantErr: false,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRows([]driver.Value{
					1, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
					[]driver.Value{
						3, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
				))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1`).
					WithArgs(
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "ADMIN", "Administrator", now, now))
			},
		},
		{
			name:    "error",
			want:    nil,
			wantErr: true,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := &repository{
				db: gormDB,
			}
			got, err := r.GetAllUser(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetAllUser() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetAllUser() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetAllAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name    string
		wantErr bool
		want    *[]models.User
		mockFn  func()
	}{
		{
			name: "success",
			want: &[]models.User{
				models.User{
					ID:          1,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "ADMIN",
						Name:      "Administrator",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				models.User{
					ID:          3,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "ADMIN",
						Name:      "Administrator",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			wantErr: false,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRows([]driver.Value{
					1, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
					[]driver.Value{
						3, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
				))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1 AND id = \$2`).
					WithArgs(
						1,
						1,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "ADMIN", "Administrator", now, now))
			},
		},
		{
			name:    "error",
			want:    nil,
			wantErr: true,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := &repository{
				db: gormDB,
			}
			got, err := r.GetAllAdmin(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetAllAdmin() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetAllAdmin() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetAllCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name    string
		wantErr bool
		want    *[]models.User
		mockFn  func()
	}{
		{
			name: "success",
			want: &[]models.User{
				models.User{
					ID:          2,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "CUSTOMER",
						Name:      "Customeer",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				models.User{
					ID:          4,
					Name:        "name",
					Username:    "username",
					Password:    "password",
					Email:       "email@gmail.com",
					RoleID:      1,
					PhoneNumber: "08332323",
					CreatedAt:   now,
					UpdatedAt:   now,
					Role: models.Role{
						ID:        1,
						Code:      "CUSTOMER",
						Name:      "Customeer",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			wantErr: false,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "role_id", "phone_number", "created_at", "updated_at"}).AddRows([]driver.Value{
					2, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
					[]driver.Value{
						4, "name", "username", "password", "email@gmail.com", 1, "08332323", now, now},
				))

				mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id" = \$1 AND id = \$2`).
					WithArgs(
						1,
						2,
					).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).AddRow(1, "CUSTOMER", "Customeer", now, now))
			},
		},
		{
			name:    "error",
			want:    nil,
			wantErr: true,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			r := &repository{
				db: gormDB,
			}
			got, err := r.GetAllCustomer(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetAllCustomer() error %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetAllCustomer() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
