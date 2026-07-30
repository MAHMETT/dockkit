package templates

import (
	"fmt"
	"strings"

	"github.com/MAHMETT/dockkit/internal/config"
	"gopkg.in/yaml.v3"
)

// ComposeConfig represents a docker-compose.yml structure.
type ComposeConfig struct {
	Version  string                    `yaml:"version"`
	Services map[string]ComposeService `yaml:"services"`
	Networks map[string]ComposeNetwork `yaml:"networks,omitempty"`
	Volumes  map[string]ComposeVolume  `yaml:"volumes,omitempty"`
}

type ComposeService struct {
	Image         string              `yaml:"image"`
	ContainerName string              `yaml:"container_name"`
	Restart       string              `yaml:"restart"`
	Environment   map[string]string   `yaml:"environment"`
	Ports         []string            `yaml:"ports"`
	Volumes       []string            `yaml:"volumes"`
	Healthcheck   *ComposeHealthcheck `yaml:"healthcheck,omitempty"`
	Networks      []string            `yaml:"networks"`
	Command       []string            `yaml:"command,omitempty"`
	ShmSize       string              `yaml:"shm_size,omitempty"`
}

type ComposeHealthcheck struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval"`
	Timeout  string   `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}

type ComposeNetwork struct {
	Name   string `yaml:"name"`
	Driver string `yaml:"driver"`
}

type ComposeVolume struct {
	Name   string `yaml:"name"`
	Driver string `yaml:"driver,omitempty"`
}

// RenderOptions configures how a template is rendered.
type RenderOptions struct {
	ServiceName    string
	Version        string
	Port           int
	User           string
	Password       string
	Database       string
	ContainerName  string
	Timezone       string
	Network        string
	RestartPolicy  string
	Command        []string
	ExtraEnv       map[string]string
}

// Render generates a docker-compose.yml from a template and options.
func Render(tmpl *Template, opts RenderOptions) (*ComposeConfig, error) {
	if tmpl == nil {
		return nil, fmt.Errorf("template is nil")
	}

	version := findVersion(tmpl, opts.Version)
	if version == nil {
		return nil, fmt.Errorf("version %q not found in template %s", opts.Version, tmpl.Name)
	}

	interpolator := NewInterpolator(nil, opts.ServiceName, opts.Version)
	interpolator.SetServiceConfig(&config.ServiceVersion{
		Port:          opts.Port,
		User:          opts.User,
		Password:      opts.Password,
		Database:      opts.Database,
		ContainerName: opts.ContainerName,
		Image:         version.Image,
	})

	if opts.Timezone != "" {
		interpolator.SetVar("GENERAL_TIMEZONE", opts.Timezone)
	}
	if opts.Network != "" {
		interpolator.SetVar("GENERAL_DEFAULT_NETWORK", opts.Network)
	}

	env := interpolator.InterpolateMap(version.EnvVars)
	if opts.ExtraEnv != nil {
		for k, v := range opts.ExtraEnv {
			env[k] = v
		}
	}

	restart := opts.RestartPolicy
	if restart == "" {
		restart = "unless-stopped"
	}

	svc := ComposeService{
		Image:         version.Image,
		ContainerName: opts.ContainerName,
		Restart:       restart,
		Environment:   env,
		Ports:         interpolator.InterpolateSlice(version.Ports),
		Volumes:       interpolator.InterpolateSlice(version.Volumes),
		Networks:      interpolator.InterpolateSlice(version.Networks),
	}

	if version.Healthcheck != nil {
		svc.Healthcheck = &ComposeHealthcheck{
			Test:     interpolator.InterpolateSlice(version.Healthcheck.Test),
			Interval: version.Healthcheck.Interval,
			Timeout:  version.Healthcheck.Timeout,
			Retries:  version.Healthcheck.Retries,
		}
	}

	if len(opts.Command) > 0 {
		svc.Command = opts.Command
	}

	serviceKey := sanitizeServiceKey(tmpl.Name)
	compose := &ComposeConfig{
		Version:  "3.8",
		Services: map[string]ComposeService{serviceKey: svc},
	}

	if len(version.Networks) > 0 {
		compose.Networks = map[string]ComposeNetwork{}
		seen := map[string]bool{}
		for _, n := range version.Networks {
			name := interpolator.Interpolate(n)
			if seen[name] {
				continue
			}
			seen[name] = true
			compose.Networks[name] = ComposeNetwork{
				Name:   name,
				Driver: "bridge",
			}
		}
	}

	return compose, nil
}

// RenderYAML renders a docker-compose config as YAML bytes.
func RenderYAML(tmpl *Template, opts RenderOptions) ([]byte, error) {
	compose, err := Render(tmpl, opts)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(compose)
}

// RenderToString renders a docker-compose config as a YAML string.
func RenderToString(tmpl *Template, opts RenderOptions) (string, error) {
	data, err := RenderYAML(tmpl, opts)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// findVersion finds a version entry by key.
func findVersion(tmpl *Template, key string) *VersionEntry {
	for i := range tmpl.Versions {
		if tmpl.Versions[i].Key == key {
			return &tmpl.Versions[i]
		}
	}
	return nil
}

// sanitizeServiceKey converts template name to a valid compose service key.
func sanitizeServiceKey(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
