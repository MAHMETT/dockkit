package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MAHMETT/dockkit/internal/config"
	"github.com/MAHMETT/dockkit/internal/templates"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [service]",
	Short: "Setup a new service interactively (e.g., dockkit setup postgresql)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]

		// Load template
		tmpl, err := templates.LoadBuiltin(serviceName)
		if err != nil {
			return fmt.Errorf("loading template %q: %w", serviceName, err)
		}

		fmt.Printf("Setting up %s %s\n", tmpl.Icon, tmpl.Name)
		fmt.Println()

		// Interactive setup
		reader := bufio.NewReader(os.Stdin)

		// Select version
		fmt.Println("Available versions:")
		for i, v := range tmpl.Versions {
			fmt.Printf("  %d) %s\n", i+1, v.Key)
		}
		fmt.Print("Select version (number): ")
		versionIdx, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading input: %w", err)
		}
		versionIdx = strings.TrimSpace(versionIdx)
		idx, err := strconv.Atoi(versionIdx)
		if err != nil || idx < 1 || idx > len(tmpl.Versions) {
			return fmt.Errorf("invalid selection: %s", versionIdx)
		}
		selectedVersion := tmpl.Versions[idx-1]

		// Get port
		fmt.Printf("Port (default %d): ", selectedVersion.DefaultPort)
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
		port := selectedVersion.DefaultPort
		if portStr != "" {
			p, err := strconv.Atoi(portStr)
			if err == nil && p > 0 {
				port = p
			}
		}

		// Get username
		defaultUser := "postgres"
		if tmpl.Name == "MySQL" || tmpl.Name == "MariaDB" {
			defaultUser = "root"
		}
		fmt.Printf("Username (default %s): ", defaultUser)
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)
		if user == "" {
			user = defaultUser
		}

		// Get password
		fmt.Printf("Password (default %s): ", user)
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if password == "" {
			password = user
		}

		// Get database
		defaultDB := user
		fmt.Printf("Database (default %s): ", defaultDB)
		database, _ := reader.ReadString('\n')
		database = strings.TrimSpace(database)
		if database == "" {
			database = defaultDB
		}

		// Generate container name
		containerName := fmt.Sprintf("dockkit-%s-%s", serviceName, selectedVersion.Key)

		// Update config
		if cfg.Services == nil {
			cfg.Services = map[string]config.Service{}
		}

		svc, ok := cfg.Services[serviceName]
		if !ok {
			svc = config.Service{
				Prefix:   strings.ToUpper(serviceName[:3]),
				Versions: map[string]config.ServiceVersion{},
			}
		}

		if svc.Versions == nil {
			svc.Versions = map[string]config.ServiceVersion{}
		}

		svc.Versions[selectedVersion.Key] = config.ServiceVersion{
			Enabled:       true,
			Port:          port,
			ContainerName: containerName,
			Image:         selectedVersion.Image,
			User:          user,
			Password:      password,
			Database:      database,
		}

		cfg.Services[serviceName] = svc

		// Save config
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		// Create service directory
		serviceDir, err := config.ServiceDir(serviceName, selectedVersion.Key)
		if err != nil {
			return fmt.Errorf("getting service dir: %w", err)
		}
		if err := os.MkdirAll(serviceDir, 0700); err != nil {
			return fmt.Errorf("creating service dir: %w", err)
		}

		// Render docker-compose.yml
		opts := templates.RenderOptions{
			ServiceName:   serviceName,
			Version:       selectedVersion.Key,
			Port:          port,
			User:          user,
			Password:      password,
			Database:      database,
			ContainerName: containerName,
			Timezone:      cfg.General.Timezone,
			Network:       cfg.General.DefaultNetwork,
		}

		composeYAML, err := templates.RenderToString(tmpl, opts)
		if err != nil {
			return fmt.Errorf("rendering docker-compose: %w", err)
		}

		// Write docker-compose.yml
		composePath := serviceDir + "/docker-compose.yml"
		if err := os.WriteFile(composePath, []byte(composeYAML), 0644); err != nil {
			return fmt.Errorf("writing docker-compose.yml: %w", err)
		}

		// Write .env file
		envContent := fmt.Sprintf("%s_PORT=%d\n%s_USER=%s\n%s_PASSWORD=%s\n%s_DATABASE=%s\n",
			strings.ToUpper(serviceName), port,
			strings.ToUpper(serviceName), user,
			strings.ToUpper(serviceName), password,
			strings.ToUpper(serviceName), database,
		)
		envPath := serviceDir + "/.env"
		if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
			return fmt.Errorf("writing .env: %w", err)
		}

		fmt.Println()
		fmt.Printf("✓ %s %s configured successfully!\n", tmpl.Name, selectedVersion.Key)
		fmt.Printf("  Config:    %s\n", serviceDir)
		fmt.Printf("  Compose:   %s/docker-compose.yml\n", serviceDir)
		fmt.Printf("  Port:      %d\n", port)
		fmt.Printf("  Container: %s\n", containerName)
		fmt.Println()
		fmt.Printf("Start with: dockkit up %s-%s\n", serviceName, selectedVersion.Key)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
