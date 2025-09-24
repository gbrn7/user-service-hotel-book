package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role,omitempty"`
	PhoneNumber string    `json:"phone_number"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Email           string `json:"email" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	RoleID          uint
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type UpdateRequest struct {
	Name            string  `json:"name" validate:"omitempty"`
	Username        string  `json:"username" validate:"omitempty"`
	Password        *string `json:"password" validate:"omitempty"`
	ConfirmPassword *string `json:"confirm_password" validate:"omitempty"`
	Email           string  `json:"email" validate:"omitempty"`
	PhoneNumber     string  `json:"phone_number" validate:"omitempty"`
	RoleID          uint
}
