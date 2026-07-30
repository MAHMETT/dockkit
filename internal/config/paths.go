package config

import (
	"os"
	"path/filepath"
)

func EnsureDirectories() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(home, ConfigDir),
		filepath.Join(home, ConfigDir, "cache", "tags"),
		filepath.Join(home, ConfigDir, "templates"),
		filepath.Join(home, ConfigDir, "services"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	return nil
}
