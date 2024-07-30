package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/fishus/go-advanced-gophkeeper/internal/adapter/config"
)

// logger is the default logger used by the application
var logger *slog.Logger

// Set sets the logger configuration based on the environment
func Set(config *config.App) {
	var level slog.Level
	switch strings.ToUpper(config.LogLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelWarn
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, opts),
	)

	slog.SetDefault(logger)
}
