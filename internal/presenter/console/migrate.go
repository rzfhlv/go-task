package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all available migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migrations success")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recent migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rollback success")
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create migration file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create Migration")
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)

	rootCmd.AddCommand(migrateCmd)
}
