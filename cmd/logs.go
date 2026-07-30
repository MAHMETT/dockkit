package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View service logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Stream logs from service
		fmt.Printf("Viewing logs for: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
