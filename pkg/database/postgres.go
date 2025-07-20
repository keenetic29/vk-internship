package database

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func InitDB(dsnCreateDB, dsnConnect, dbName string) (*gorm.DB, error) {
    // Сначала подключаемся к БД postgres (которая всегда есть)
    adminDB, err := gorm.Open(postgres.Open(dsnCreateDB), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to admin DB: %w", err)
    }

    // Проверяем существование указанной БД
    var count int64
    adminDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
    
    if count == 0 {
        // Создаём БД если её нет
        if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
            return nil, fmt.Errorf("failed to create database: %w", err)
        }
    }

    // Теперь подключаемся к нужной БД
    db, err := gorm.Open(postgres.Open(dsnConnect), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to target DB: %w", err)
    }

    return db, nil
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Advertisement{},
	)
}