package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start REST API server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Start the server")
		},
	})
}
