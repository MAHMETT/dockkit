package cmd

import (
	"fmt"
	"strings"

	"github.com/MAHMETT/dockkit/internal/config"
)

// parseServiceArg parses "service-version" into service name and version.
// Example: "postgresql-16" -> "postgresql", "16"
// Example: "redis-7" -> "redis", "7"
// Example: "mysql-8.0" -> "mysql", "8.0"
func parseServiceArg(arg string) (name, version string, err error) {
	// Find the last dash that separates service from version
	// Handle cases like "postgresql-16", "mysql-8.0", "redis-7"
	idx := strings.LastIndex(arg, "-")
	if idx == -1 {
		return "", "", fmt.Errorf("invalid format %q: expected 'service-version' (e.g., postgresql-16)", arg)
	}

	name = arg[:idx]
	version = arg[idx+1:]

	if name == "" || version == "" {
		return "", "", fmt.Errorf("invalid format %q: expected 'service-version' (e.g., postgresql-16)", arg)
	}

	return name, version, nil
}

// findServiceVersion finds a service version in the config.
func findServiceVersion(cfg *config.Config, name, version string) (*config.ServiceVersion, error) {
	svc, ok := cfg.Services[name]
	if !ok {
		return nil, fmt.Errorf("service %q not found", name)
	}

	ver, ok := svc.Versions[version]
	if !ok {
		return nil, fmt.Errorf("version %q not found for service %q", version, name)
	}

	return &ver, nil
}

// getContainerName returns the container name for a service version.
func getContainerName(cfg *config.Config, name, version string) (string, error) {
	ver, err := findServiceVersion(cfg, name, version)
	if err != nil {
		return "", err
	}
	if ver.ContainerName == "" {
		return "", fmt.Errorf("container name not configured for %s-%s", name, version)
	}
	return ver.ContainerName, nil
}

// getServiceDir returns the service directory path.
func getServiceDir(name, version string) (string, error) {
	dir, err := config.ServiceDir(name, version)
	if err != nil {
		return "", fmt.Errorf("getting service dir: %w", err)
	}
	return dir, nil
}
