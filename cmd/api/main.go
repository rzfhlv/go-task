package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/presenter/console"
)

func main() {
	ctx := context.Background()
	config.All()

	if err := console.Execute(); err != nil {
		slog.ErrorContext(ctx, "failed init console", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
