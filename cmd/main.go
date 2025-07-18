package main

import (
	"VK/internal/api"
	"VK/internal/config"
	"VK/pkg/database"
	"VK/pkg/logger"
	"log"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// Инициализация логгера
	if err := logger.Init(cfg.LogDebug, cfg.LogFile); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting application", 
		"version", "1.0.0",
		"debug", cfg.LogDebug,
	)

	// Инициализация базы данных
	db, err := database.InitDB(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Автомиграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	// Настройка маршрутов
	router := api.SetupRouter(db, cfg.JWTSecret)

	// Запуск сервера
	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatal("Failed to start server", err)
	}
}