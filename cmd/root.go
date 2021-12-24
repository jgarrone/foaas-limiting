package cmd

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "foaas",
		Short: "FOAAS Limiting command-line interface",
	}

	// Add new commands here.
	rootCmd.AddCommand(
		ServeCommand(),
	)

	return rootCmd
}
