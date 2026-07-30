package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var containerNameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []*ValidationError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}
	msg := errs[0].Message
	for _, e := range errs[1:] {
		msg += "; " + e.Message
	}
	return msg
}

func (errs ValidationErrors) HasErrors() bool {
	return len(errs) > 0
}

// Validate checks the entire config for errors.
func Validate(cfg *Config) ValidationErrors {
	var errs ValidationErrors

	if cfg.Version == "" {
		errs = append(errs, &ValidationError{Field: "version", Message: "version is required"})
	}

	if err := validateGeneral(&cfg.General); err != nil {
		errs = append(errs, err...)
	}

	for name, svc := range cfg.Services {
		if err := validateService(name, &svc); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

func validateGeneral(g *GeneralConfig) ValidationErrors {
	var errs ValidationErrors

	if g.Timezone == "" {
		errs = append(errs, &ValidationError{Field: "general.timezone", Message: "timezone is required"})
	}

	if g.DefaultNetwork == "" {
		errs = append(errs, &ValidationError{Field: "general.default_network", Message: "default_network is required"})
	}

	if g.RefreshInterval != "" {
		if _, err := time.ParseDuration(g.RefreshInterval); err != nil {
			errs = append(errs, &ValidationError{
				Field:   "general.refresh_interval",
				Message: fmt.Sprintf("invalid duration %q: %v", g.RefreshInterval, err),
			})
		}
	}

	return errs
}

func validateService(name string, svc *Service) ValidationErrors {
	var errs ValidationErrors

	if svc.Prefix == "" {
		errs = append(errs, &ValidationError{
			Field:   fmt.Sprintf("services.%s.prefix", name),
			Message: "prefix is required",
		})
	}

	if len(svc.Versions) == 0 {
		errs = append(errs, &ValidationError{
			Field:   fmt.Sprintf("services.%s.versions", name),
			Message: "at least one version is required",
		})
	}

	for ver, v := range svc.Versions {
		if err := validateServiceVersion(name, ver, &v); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

func validateServiceVersion(service, version string, v *ServiceVersion) ValidationErrors {
	var errs ValidationErrors
	prefix := fmt.Sprintf("services.%s.versions.%s", service, version)

	if v.Port < 1024 || v.Port > 65535 {
		errs = append(errs, &ValidationError{
			Field:   prefix + ".port",
			Message: fmt.Sprintf("port %d must be between 1024 and 65535", v.Port),
		})
	}

	if v.ContainerName == "" {
		errs = append(errs, &ValidationError{
			Field:   prefix + ".container_name",
			Message: "container_name is required",
		})
	} else if !containerNameRe.MatchString(v.ContainerName) {
		errs = append(errs, &ValidationError{
			Field:   prefix + ".container_name",
			Message: fmt.Sprintf("invalid container name %q", v.ContainerName),
		})
	}

	if v.Image == "" {
		errs = append(errs, &ValidationError{
			Field:   prefix + ".image",
			Message: "image is required",
		})
	}

	return errs
}

// ValidatePort checks if a port number is valid.
func ValidatePort(port string) (int, error) {
	n, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: must be a number", port)
	}
	if n < 1024 || n > 65535 {
		return 0, fmt.Errorf("port %d must be between 1024 and 65535", n)
	}
	return n, nil
}

// ValidateContainerName checks if a container name is valid.
func ValidateContainerName(name string) error {
	if name == "" {
		return fmt.Errorf("container name is required")
	}
	if !containerNameRe.MatchString(name) {
		return fmt.Errorf("invalid container name %q: must start with alphanumeric, can contain [a-zA-Z0-9_.-]", name)
	}
	if len(name) > 128 {
		return fmt.Errorf("container name too long: max 128 characters")
	}
	return nil
}
