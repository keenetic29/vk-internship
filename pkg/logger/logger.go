package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	Log  *slog.Logger
	file *os.File
	mu   sync.Mutex
)

// Init инициализирует логгер с записью в файл и консоль
func Init(debug bool, logDir string) error {
	mu.Lock()
	defer mu.Unlock()

	// Создаем директорию для логов, если ее нет
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Открываем файл логов
	logPath := filepath.Join(logDir, "marketplace.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	file = f

	// Настраиваем уровень логирования
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	// Создаем мультиписатель (файл + консоль)
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Настраиваем обработчик
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	// Инициализируем логгер
	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}

// Sync гарантирует, что все логи будут записаны
func Sync() {
	mu.Lock()
	defer mu.Unlock()
	
	if file != nil {
		file.Sync()
		file.Close()
	}
}