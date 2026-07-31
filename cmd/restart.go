package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MAHMETT/dockkit/internal/docker"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service-version]",
	Short: "Restart a service (e.g., dockkit restart postgresql-16)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, version, err := parseServiceArg(args[0])
		if err != nil {
			return err
		}

		// Find service in config
		ver, err := findServiceVersion(cfg, name, version)
		if err != nil {
			return err
		}

		if !ver.Enabled {
			fmt.Printf("Service %s-%s is disabled.\n", name, version)
			return nil
		}

		// Get service directory
		dir, err := getServiceDir(name, version)
		if err != nil {
			return err
		}

		// Check if docker-compose.yml exists
		if _, statErr := os.Stat(dir + "/docker-compose.yml"); os.IsNotExist(statErr) {
			fmt.Printf("No docker-compose.yml found at %s\n", dir)
			return nil
		}

		// Restart service with docker compose
		fmt.Printf("Restarting %s %s...\n", name, version)
		if err := docker.ComposeRestart(context.Background(), dir); err != nil {
			return fmt.Errorf("restarting service: %w", err)
		}

		fmt.Printf("✓ %s %s restarted\n", name, version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
