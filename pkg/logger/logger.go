package logger

import (
	"log/slog"
	"os"
)

func SetDefault(level string) {
	var slogLevel slog.Leveler

	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	}))

	slog.SetDefault(logger)
}
