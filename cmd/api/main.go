package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rzfhlv/go-test/internal/presenter/console"
)

func main() {
	ctx := context.Background()

	if err := console.Execute(); err != nil {
		slog.ErrorContext(ctx, "failed init console", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
