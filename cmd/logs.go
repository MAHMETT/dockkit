package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MAHMETT/dockkit/internal/docker"
	"github.com/spf13/cobra"
)

var (
	logsFollow bool
	logsTail   string
)

var logsCmd = &cobra.Command{
	Use:   "logs [service-version]",
	Short: "View service logs (e.g., dockkit logs postgresql-16)",
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

		// Get container name
		containerName, err := getContainerName(cfg, name, version)
		if err != nil {
			return err
		}

		// Create Docker client
		client, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("creating docker client: %w", err)
		}
		defer client.Close()

		// Check Docker is running
		ctx := context.Background()
		if err := client.Ping(ctx); err != nil {
			return fmt.Errorf("docker is not running: %w", err)
		}

		// Stream logs
		fmt.Printf("Logs for %s %s (container: %s)\n", name, version, containerName)
		fmt.Println("---")
		if logsFollow {
			fmt.Println("(Press Ctrl+C to stop following)")
			fmt.Println("---")
		}

		reader, err := client.StreamLogs(ctx, containerName, logsFollow)
		if err != nil {
			return fmt.Errorf("streaming logs: %w", err)
		}
		defer reader.Close()

		// Copy logs to stdout
		if _, err := os.Stdout.ReadFrom(reader); err != nil {
			return fmt.Errorf("reading logs: %w", err)
		}

		return nil
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&logsTail, "tail", "100", "Number of lines to show")
	rootCmd.AddCommand(logsCmd)
}
