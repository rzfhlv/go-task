package console

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Service CLI",
}

func Execute() error {
	return rootCmd.Execute()
}
