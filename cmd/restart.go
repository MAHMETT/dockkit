package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service]",
	Short: "Restart a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Restart service
		fmt.Printf("Restarting service: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
