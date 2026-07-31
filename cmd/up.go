package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MAHMETT/dockkit/internal/docker"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [service-version]",
	Short: "Start a service (e.g., dockkit up postgresql-16)",
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
			fmt.Printf("Service %s-%s is disabled. Enable it in config first.\n", name, version)
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
			fmt.Printf("Run 'dockkit setup %s' first to generate config.\n", name)
			return nil
		}

		// Create Docker client for pre-flight checks
		client, err := docker.NewClient()
		if err != nil {
			fmt.Printf("Warning: cannot connect to Docker: %v\n", err)
			fmt.Println("Attempting to start service anyway...")
		} else {
			defer client.Close()
			ctx := context.Background()
			if err := client.Ping(ctx); err != nil {
				fmt.Printf("Warning: Docker is not running: %v\n", err)
				fmt.Println("Attempting to start service anyway...")
			}
		}

		// Start service with docker compose
		fmt.Printf("Starting %s %s...\n", name, version)
		if err := docker.ComposeUp(context.Background(), dir); err != nil {
			return fmt.Errorf("starting service: %w", err)
		}

		fmt.Printf("✓ %s %s started on port %d\n", name, version, ver.Port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
