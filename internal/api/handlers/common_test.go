package handlers_test

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/keenetic29/vk-internship/pkg/logger"
	"os"
	"path/filepath"
)

// Вынес функции в отдельный файл, поскольку, находясь в одном пакете handlers_test, требуются в обоих тестах

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestMain(m *testing.M) {
	// Инициализация логгера для всех тестов
	logDir := filepath.Join("..", "..", "..", "logs")
	if err := logger.Init("true", logDir, "test.log"); err != nil {
		panic(err)
	}
	defer logger.Sync()
	// Запуск тестов
	code := m.Run()
	os.Exit(code)
}