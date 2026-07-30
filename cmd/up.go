package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [service]",
	Short: "Start a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Start service with docker compose
		fmt.Printf("Starting service: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
