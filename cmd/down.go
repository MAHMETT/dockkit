package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [service]",
	Short: "Stop a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Stop service with docker compose
		fmt.Printf("Stopping service: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
