package database

import (
	"VK/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Advertisement{},
		&domain.AuthToken{},
	)
}