package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.Version != VersionCurrent {
		t.Errorf("Version = %q, want %q", cfg.Version, VersionCurrent)
	}
	if cfg.General.Timezone == "" {
		t.Error("General.Timezone is empty")
	}
	if cfg.General.DefaultNetwork == "" {
		t.Error("General.DefaultNetwork is empty")
	}
	if cfg.Services == nil {
		t.Error("Services is nil")
	}
}

func TestLoadFromBytes(t *testing.T) {
	data := []byte(`
version: "1"
general:
  timezone: Asia/Jakarta
  default_network: my-network
  auto_refresh: true
  refresh_interval: 30s
services:
  postgresql:
    prefix: PG
    versions:
      "16":
        enabled: true
        port: 5432
        container_name: my-pg-16
        image: postgres:16-alpine
        user: postgres
        password: secret
        database: mydb
`)
	cfg, err := LoadFromBytes(data)
	if err != nil {
		t.Fatalf("LoadFromBytes() error = %v", err)
	}
	if cfg.General.Timezone != "Asia/Jakarta" {
		t.Errorf("Timezone = %q, want %q", cfg.General.Timezone, "Asia/Jakarta")
	}
	if cfg.General.DefaultNetwork != "my-network" {
		t.Errorf("DefaultNetwork = %q, want %q", cfg.General.DefaultNetwork, "my-network")
	}
	svc, ok := cfg.Services["postgresql"]
	if !ok {
		t.Fatal("postgresql service not found")
	}
	v, ok := svc.Versions["16"]
	if !ok {
		t.Fatal("version 16 not found")
	}
	if v.Port != 5432 {
		t.Errorf("Port = %d, want 5432", v.Port)
	}
	if v.ContainerName != "my-pg-16" {
		t.Errorf("ContainerName = %q, want %q", v.ContainerName, "my-pg-16")
	}
}

func TestLoadFromBytes_Invalid(t *testing.T) {
	_, err := LoadFromBytes([]byte("not: valid: yaml: [[["))
	if err == nil {
		t.Error("LoadFromBytes() should fail on invalid YAML")
	}
}

func TestLoadFromBytes_Empty(t *testing.T) {
	cfg, err := LoadFromBytes([]byte{})
	if err != nil {
		t.Fatalf("LoadFromBytes() error = %v", err)
	}
	// Should return default config with defaults filled
	if cfg.General.DefaultNetwork == "" {
		t.Error("DefaultNetwork should be filled")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := DefaultConfig()
	cfg.General.Timezone = "America/New_York"

	if err := SaveToFile(cfg, path); err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if loaded.General.Timezone != "America/New_York" {
		t.Errorf("Timezone = %q, want %q", loaded.General.Timezone, "America/New_York")
	}
}

func TestLoadFromFile_NotExist(t *testing.T) {
	cfg, err := LoadFromFile("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("LoadFromFile() should return default on missing file, got error = %v", err)
	}
	if cfg.Version != VersionCurrent {
		t.Errorf("Version = %q, want %q", cfg.Version, VersionCurrent)
	}
}

func TestSaveToFile_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// Write initial config
	cfg := DefaultConfig()
	if err := SaveToFile(cfg, path); err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// Write again — should create backup
	cfg.General.Timezone = "UTC"
	if err := SaveToFile(cfg, path); err != nil {
		t.Fatalf("SaveToFile() second call error = %v", err)
	}

	// Check backup exists
	entries, _ := filepath.Glob(filepath.Join(dir, "config.yaml.backup.*"))
	if len(entries) == 0 {
		t.Error("backup file not created")
	}
}

func TestEnsureDirs(t *testing.T) {
	dir := t.TempDir()

	// Override home dir by testing the logic
	_ = dir

	// Test that EnsureDirs doesn't error
	// (It uses real home dir, so just check it doesn't panic)
	err := EnsureDirs()
	if err != nil {
		t.Errorf("EnsureDirs() error = %v", err)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "missing timezone",
			cfg: &Config{
				Version: "1",
				General: GeneralConfig{
					DefaultNetwork: "net",
				},
				Services: map[string]Service{},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			cfg: &Config{
				Version: "1",
				General: GeneralConfig{
					Timezone:       "UTC",
					DefaultNetwork: "net",
				},
				Services: map[string]Service{
					"pg": {
						Prefix: "PG",
						Versions: map[string]ServiceVersion{
							"16": {Port: 100, ContainerName: "c", Image: "img"},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.cfg)
			if (len(errs) > 0) != tt.wantErr {
				t.Errorf("Validate() errors = %v, wantErr %v", errs, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		input string
		want  int
		err   bool
	}{
		{"5432", 5432, false},
		{"80", 0, true},    // too low
		{"abc", 0, true},   // not a number
		{"99999", 0, true}, // too high
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ValidatePort(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ValidatePort(%q) error = %v, wantErr %v", tt.input, err, tt.err)
			}
			if got != tt.want {
				t.Errorf("ValidatePort(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateContainerName(t *testing.T) {
	tests := []struct {
		name string
		err  bool
	}{
		{"my-container", false},
		{"container_123", false},
		{"container.v2", false},
		{"", true},                    // empty
		{"-starts-with-dash", true},   // starts with dash
		{"has spaces", true},          // spaces
		{"has@special", true},         // special chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContainerName(tt.name)
			if (err != nil) != tt.err {
				t.Errorf("ValidateContainerName(%q) error = %v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestDefaultServiceVersions(t *testing.T) {
	services := DefaultServiceVersions()
	if len(services) == 0 {
		t.Error("DefaultServiceVersions() returned empty")
	}
	pg, ok := services["postgresql"]
	if !ok {
		t.Error("postgresql not in default services")
	}
	if pg.Prefix != "PG" {
		t.Errorf("PG prefix = %q, want %q", pg.Prefix, "PG")
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("hello"), 0600)
	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}

	data, _ := os.ReadFile(dst)
	if string(data) != "hello" {
		t.Errorf("copied content = %q, want %q", string(data), "hello")
	}
}
