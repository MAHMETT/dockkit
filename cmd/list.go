package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: List services from config
		fmt.Println("Listing services...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
