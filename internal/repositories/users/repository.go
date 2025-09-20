package users

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}
