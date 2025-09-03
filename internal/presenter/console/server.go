package console

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/infrastructure"
	"github.com/rzfhlv/go-task/internal/presenter/rest"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start REST API server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cfg := config.Get()
			infra, err := infrastructure.New(ctx, cfg)
			if err != nil {
				log.Fatalf("fail to load infrastructure: %v", err)
			}

			e := rest.Init(infra, cfg)

			// start server
			go func() {
				if err := e.Start(fmt.Sprintf(":%s", cfg.App.Port)); err != nil && err != http.ErrServerClosed {
					e.Logger.Fatal("shutting down the server")
				}
			}()

			// graceful shutdown
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}

			if err := infra.SQLStore().Close(); err != nil {
				e.Logger.Fatal(err)
			}

			if err := infra.Redis().Close(); err != nil {
				e.Logger.Fatal(err)
			}
		},
	})
}
