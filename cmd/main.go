package main

import (
	"VK/internal/api"
	"VK/internal/config"
	"VK/pkg/database"
	"VK/pkg/logger"
	"os"
)

func main() {
	if err := logger.Init(true, "./logs"); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Log.Info("Starting marketplace application", 
		"version", "1.0.0",
		"environment", "development")

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Log.Debug("Configuration loaded", "db_host", cfg.DBHost, "server_address", cfg.ServerAddr)

	db, err := database.InitDB(cfg.GetDBConnectionString()	)
	if err != nil {
		logger.Log.Error("Database connection failed", 
			"error", err,
			"host", cfg.DBHost,
			"port", cfg.DBPort)
		os.Exit(1)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	logger.Log.Info("Database connection established")

	if err := database.RunMigrations(db); err != nil {
		logger.Log.Error("Database migrations failed", "error", err)
		os.Exit(1)
	}
	logger.Log.Info("Database migrations completed")

	router := api.SetupRouter(db, cfg.JWTSecret)
	logger.Log.Info("Starting HTTP server", 
		"address", cfg.ServerAddr,
		"log_file", "./logs/marketplace.log")

	if err := router.Run(cfg.ServerAddr); err != nil {
		logger.Log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

