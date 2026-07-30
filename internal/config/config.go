package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName = ".config/dockkit"
	ConfigFileName = "config.yaml"
	VersionCurrent = "1"
)

type Config struct {
	Version  string             `yaml:"version"`
	General  GeneralConfig      `yaml:"general"`
	Services map[string]Service `yaml:"services"`
}

type GeneralConfig struct {
	Timezone        string `yaml:"timezone"`
	DefaultNetwork  string `yaml:"default_network"`
	AutoRefresh     bool   `yaml:"auto_refresh"`
	RefreshInterval string `yaml:"refresh_interval"`
}

type Service struct {
	Prefix   string                    `yaml:"prefix"`
	Versions map[string]ServiceVersion `yaml:"versions"`
}

type ServiceVersion struct {
	Enabled       bool   `yaml:"enabled"`
	Port          int    `yaml:"port"`
	ContainerName string `yaml:"container_name"`
	Image         string `yaml:"image"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Database      string `yaml:"database"`
}

// ConfigDir returns ~/.config/dockkit/
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ConfigDirName), nil
}

// ConfigFile returns ~/.config/dockkit/config.yaml
func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

// ServicesDir returns ~/.config/dockkit/services/
func ServicesDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "services"), nil
}

// TemplatesDir returns ~/.config/dockkit/templates/
func TemplatesDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "templates"), nil
}

// CacheDir returns ~/.config/dockkit/cache/tags/
func CacheDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "cache", "tags"), nil
}

// ServiceDir returns ~/.config/dockkit/services/<name>/<version>/
func ServiceDir(name, version string) (string, error) {
	dir, err := ServicesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name, version), nil
}

// Load reads config from default path.
// Returns default config if file doesn't exist.
func Load() (*Config, error) {
	path, err := ConfigFile()
	if err != nil {
		return nil, err
	}
	return LoadFromFile(path)
}

// LoadFromFile reads config from a specific path.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes parses config from YAML bytes.
func LoadFromBytes(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	cfg.fillDefaults()
	return &cfg, nil
}

// Save writes config to default path with backup.
func Save(cfg *Config) error {
	path, err := ConfigFile()
	if err != nil {
		return err
	}
	return SaveToFile(cfg, path)
}

// SaveToFile writes config to a specific path with backup.
func SaveToFile(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	// Backup existing config
	if _, err := os.Stat(path); err == nil {
		backup := path + ".backup." + time.Now().Format("20060102-150405.000000")
		if err := copyFile(path, backup); err != nil {
			return fmt.Errorf("backing up config: %w", err)
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// EnsureDirs creates all required directories.
func EnsureDirs() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home dir: %w", err)
	}

	dirs := []string{
		filepath.Join(home, ConfigDirName),
		filepath.Join(home, ConfigDirName, "services"),
		filepath.Join(home, ConfigDirName, "templates"),
		filepath.Join(home, ConfigDirName, "cache", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("creating dir %s: %w", dir, err)
		}
	}
	return nil
}

// fillDefaults sets zero-value fields to defaults.
func (c *Config) fillDefaults() {
	if c.Version == "" {
		c.Version = VersionCurrent
	}
	if c.General.Timezone == "" {
		c.General.Timezone = "UTC"
	}
	if c.General.DefaultNetwork == "" {
		c.General.DefaultNetwork = "dockkit-network"
	}
	if c.General.RefreshInterval == "" {
		c.General.RefreshInterval = "30s"
	}
	if c.Services == nil {
		c.Services = map[string]Service{}
	}
	for name, svc := range c.Services {
		if svc.Versions == nil {
			svc.Versions = map[string]ServiceVersion{}
			c.Services[name] = svc
		}
	}
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}
