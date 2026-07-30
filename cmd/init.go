package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dockkit configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Create ~/.config/dockkit/ structure
		fmt.Println("Initializing dockkit configuration...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
