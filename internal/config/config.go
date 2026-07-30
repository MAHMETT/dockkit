package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".config/dockkit"
	ConfigFile = "config.yaml"
)

type Config struct {
	Version  string             `yaml:"version"`
	General  GeneralConfig      `yaml:"general"`
	Services map[string]Service `yaml:"services"`
}

type GeneralConfig struct {
	Timezone       string `yaml:"timezone"`
	DefaultNetwork string `yaml:"default_network"`
	AutoRefresh    bool   `yaml:"auto_refresh"`
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

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir), nil
}

func ConfigFilePath() (string, error) {
	dir, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFile), nil
}

func Load() (*Config, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create config dir: %w", err)
	}

	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Version: "1",
		General: GeneralConfig{
			Timezone:        "Asia/Jakarta",
			DefaultNetwork:  "dockkit-network",
			AutoRefresh:     true,
			RefreshInterval: "30s",
		},
		Services: map[string]Service{},
	}
}
