package database

import (
	"fmt"
	"net/url"
	"time"
	"user-service/helpers/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(config configs.Config) (*gorm.DB, error) {
	encodePassword := url.QueryEscape(config.Database.Password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.Username,
		encodePassword,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)

	db, err := gorm.Open(postgres.Open(uri))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)

	return db, nil
}
