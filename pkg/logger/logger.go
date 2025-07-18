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

func Init(debug string, logDir string) error {
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "marketplace.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	level := slog.LevelInfo
	if debug == "true" {
		level = slog.LevelDebug
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}

func Sync() {
	mu.Lock()
	defer mu.Unlock()
	
	if file != nil {
		file.Sync()
		file.Close()
	}
}