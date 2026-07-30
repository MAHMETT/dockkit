package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [service]",
	Short: "Setup a new service (non-TUI wizard)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Interactive CLI setup for service
		fmt.Printf("Setting up service: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
