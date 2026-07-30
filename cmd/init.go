package cmd

import (
	"fmt"

	"github.com/MAHMETT/dockkit/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dockkit configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.EnsureDirs(); err != nil {
			return fmt.Errorf("creating directories: %w", err)
		}
		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving default config: %w", err)
		}
		fmt.Println("dockkit initialized at ~/.config/dockkit/")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
