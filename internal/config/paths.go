package config

import (
	"os"
	"path/filepath"
)

// EnsureDirs creates all required directories.
// Deprecated: Use config.EnsureDirs() instead.
func EnsureDirectories() error {
	return EnsureDirs()
}

// ServiceDirPath returns the path for a specific service version.
// Deprecated: Use config.ServiceDir() instead.
func ServiceDirPath(name, version string) (string, error) {
	return ServiceDir(name, version)
}

// CustomTemplatesDir returns the path for custom templates.
// Deprecated: Use config.TemplatesDir() instead.
func CustomTemplatesDir() (string, error) {
	return TemplatesDir()
}

// EnsureDir creates a single directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0700)
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// HomeDir returns the user's home directory.
func HomeDir() (string, error) {
	return os.UserHomeDir()
}

// JoinPath joins path elements.
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}
