package templates

import (
	"regexp"
	"strconv"

	"github.com/MAHMETT/dockkit/internal/config"
)

var varRe = regexp.MustCompile(`\$\{([A-Z_]+)\}`)

// Interpolator resolves ${VAR} placeholders in template strings.
type Interpolator struct {
	vars map[string]string
}

// NewInterpolator creates an interpolator with config values.
func NewInterpolator(cfg *config.Config, serviceName, version string) *Interpolator {
	vars := map[string]string{}

	if cfg != nil {
		vars["GENERAL_TIMEZONE"] = cfg.General.Timezone
		vars["GENERAL_DEFAULT_NETWORK"] = cfg.General.DefaultNetwork
	}

	vars["SERVICE_NAME"] = serviceName
	vars["VERSION_NUMBER"] = version

	return &Interpolator{vars: vars}
}

// SetVar adds or overrides a variable.
func (i *Interpolator) SetVar(key, value string) {
	i.vars[key] = value
}

// SetServiceConfig sets CONFIG_* variables from service version config.
func (i *Interpolator) SetServiceConfig(sv *config.ServiceVersion) {
	if sv == nil {
		return
	}
	if sv.Port > 0 {
		i.vars["CONFIG_PORT"] = strconv.Itoa(sv.Port)
	}
	i.vars["CONFIG_USER"] = sv.User
	i.vars["CONFIG_PASSWORD"] = sv.Password
	i.vars["CONFIG_DATABASE"] = sv.Database
	i.vars["CONFIG_IMAGE"] = sv.Image
	i.vars["CONFIG_CONTAINER_NAME"] = sv.ContainerName
}

// Interpolate replaces ${VAR} in s with resolved values.
// Unresolved variables are left as-is.
func (i *Interpolator) Interpolate(s string) string {
	return varRe.ReplaceAllStringFunc(s, func(match string) string {
		varName := match[2 : len(match)-1]
		if val, ok := i.vars[varName]; ok {
			return val
		}
		return match
	})
}

// InterpolateMap runs Interpolate on all values in a map.
func (i *Interpolator) InterpolateMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = i.Interpolate(v)
	}
	return result
}

// InterpolateSlice runs Interpolate on all elements in a slice.
func (i *Interpolator) InterpolateSlice(s []string) []string {
	if s == nil {
		return nil
	}
	result := make([]string, len(s))
	for idx, v := range s {
		result[idx] = i.Interpolate(v)
	}
	return result
}

// AvailableVars returns all resolved variable names and values.
func (i *Interpolator) AvailableVars() map[string]string {
	result := make(map[string]string, len(i.vars))
	for k, v := range i.vars {
		result[k] = v
	}
	return result
}

// UnresolvedVars returns variables in s that couldn't be resolved.
func (i *Interpolator) UnresolvedVars(s string) []string {
	matches := varRe.FindAllStringSubmatch(s, -1)
	seen := map[string]bool{}
	var unresolved []string
	for _, m := range matches {
		varName := m[1]
		if _, ok := i.vars[varName]; !ok && !seen[varName] {
			unresolved = append(unresolved, varName)
			seen[varName] = true
		}
	}
	return unresolved
}

// HasUnresolved returns true if s contains unresolved variables.
func (i *Interpolator) HasUnresolved(s string) bool {
	return len(i.UnresolvedVars(s)) > 0
}
