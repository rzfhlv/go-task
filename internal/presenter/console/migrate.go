package console

import (
	"fmt"
	"log/slog"

	"github.com/rzfhlv/go-task/internal/infrastructure/sqlstore"
	"github.com/spf13/cobra"
)

var (
	step int
	name string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all available migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := sqlstore.NewMigrator()
		if err != nil {
			slog.Error("error when initaite migrator", slog.String("error", err.Error()))
		}

		err = migrator.Migrate()
		if err != nil {
			slog.Error("error when migrate", slog.String("error", err.Error()))
		}

		fmt.Println("Migration success")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recent migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := sqlstore.NewMigrator()
		if err != nil {
			slog.Error("error when initaite migrator", slog.String("error", err.Error()))
		}

		err = migrator.Rollback(step)
		if err != nil {
			slog.Error("error when rollback", slog.String("error", err.Error()))
		}

		fmt.Println("Rollback success")
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create migration file",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := sqlstore.NewMigrator()
		if err != nil {
			slog.Error("error when initaite migrator", slog.String("error", err.Error()))
		}

		err = migrator.Create(name)
		if err != nil {
			slog.Error("error when create migration", slog.String("error", err.Error()))
		}

		fmt.Println("Create migration success")
	},
}

func init() {
	migrateDownCmd.Flags().IntVar(&step, "step", 1, "number of rollback to apply")
	migrateCreateCmd.Flags().StringVar(&name, "name", "", "name of migration file")

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)

	rootCmd.AddCommand(migrateCmd)
}
