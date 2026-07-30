package templates

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtins/*.yaml
var builtinFS embed.FS

type Template struct {
	Name         string         `yaml:"name"`
	Description  string         `yaml:"description"`
	Category     string         `yaml:"category"`
	Icon         string         `yaml:"icon"`
	Versions     []VersionEntry `yaml:"versions"`
	ConfigFields []ConfigField  `yaml:"config_fields"`
}

type VersionEntry struct {
	Key         string            `yaml:"key"`
	Image       string            `yaml:"image"`
	DefaultPort int               `yaml:"default_port"`
	Healthcheck *HealthcheckConf  `yaml:"healthcheck,omitempty"`
	EnvVars     map[string]string `yaml:"env_vars"`
	Volumes     []string          `yaml:"volumes"`
	Networks    []string          `yaml:"networks"`
	Ports       []string          `yaml:"ports"`
}

type HealthcheckConf struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval"`
	Timeout  string   `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}

type ConfigField struct {
	Key      string   `yaml:"key"`
	Label    string   `yaml:"label"`
	Type     string   `yaml:"type"`
	Default  string   `yaml:"default"`
	Required bool     `yaml:"required"`
	Options  []string `yaml:"options,omitempty"`
}

func LoadBuiltin(name string) (*Template, error) {
	data, err := builtinFS.ReadFile(fmt.Sprintf("builtins/%s.yaml", name))
	if err != nil {
		return nil, fmt.Errorf("template %s not found: %w", name, err)
	}

	var tmpl Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("invalid template %s: %w", name, err)
	}

	return &tmpl, nil
}

func ListBuiltin() ([]string, error) {
	entries, err := builtinFS.ReadDir("builtins")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".yaml") {
				names = append(names, strings.TrimSuffix(name, ".yaml"))
			}
		}
	}
	return names, nil
}
