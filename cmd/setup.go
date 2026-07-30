package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [service]",
	Short: "Setup a new service interactively (e.g., dockkit setup postgresql)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Layer 6 (CLI Commands)
		return fmt.Errorf("not yet implemented — coming in Layer 6")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
