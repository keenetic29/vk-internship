package main

import (
	"github.com/keenetic29/vk-internship/internal/api"
	"github.com/keenetic29/vk-internship/internal/config"
	"github.com/keenetic29/vk-internship/internal/repository"
	"github.com/keenetic29/vk-internship/internal/services"
	"github.com/keenetic29/vk-internship/pkg/database"
	"github.com/keenetic29/vk-internship/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	if err := logger.Init(cfg.LogDebug, cfg.LogFile, "marketplace.go"); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting application", 
		"version", "1.0.0",
		"debug", cfg.LogDebug,
	)

	db, err := database.InitDB(cfg.GetBDCreateString(), cfg.GetDBConnectionString(), cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	userRepo := repository.NewUserRepository(db)
	adRepo := repository.NewAdvertisementRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	adService := services.NewAdvertisementService(adRepo)

	router := api.SetupRouter(authService, adService, cfg.JWTSecret)

	if err := router.Run(":"+cfg.ServerAddr); err != nil {
		log.Fatal("Failed to start server", err)
	}
}