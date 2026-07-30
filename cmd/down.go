package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [service-version]",
	Short: "Stop a service (e.g., dockkit down postgresql-16)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Layer 2 (Docker Core)
		return fmt.Errorf("not yet implemented — coming in Layer 2")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
