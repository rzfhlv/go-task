package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/presenter/console"
	"github.com/rzfhlv/go-task/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg := config.All()
	logger.SetDefault(cfg.App.LogLevel)

	if err := console.Execute(); err != nil {
		slog.ErrorContext(ctx, "failed init console", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
