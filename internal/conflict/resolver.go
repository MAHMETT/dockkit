package conflict

import (
	"fmt"

	"github.com/MAHMETT/dockkit/internal/config"
)

// ResolutionAction represents the action to resolve a conflict.
type ResolutionAction int

const (
	ActionAutoFix ResolutionAction = iota // automatically apply suggestion
	ActionManual                          // require manual intervention
	ActionSkip                            // skip/ignore the conflict
)

// Resolution represents how to resolve a conflict.
type Resolution struct {
	Conflict Conflict
	Action   ResolutionAction
	Fix      *Fix // nil if Action != ActionAutoFix
}

// Fix describes a specific change to resolve a conflict.
type Fix struct {
	Field       string // config field to change (e.g., "services.postgresql.versions.16.port")
	OldValue    string
	NewValue    string
	Description string
}

// Resolver resolves conflicts automatically when possible.
type Resolver struct {
	detector *Detector
}

// NewResolver creates a new conflict resolver.
func NewResolver(detector *Detector) *Resolver {
	return &Resolver{detector: detector}
}

// ResolveAll resolves all detected conflicts.
func (r *Resolver) ResolveAll(conflicts ConflictList) []Resolution {
	var resolutions []Resolution
	for _, conflict := range conflicts {
		resolutions = append(resolutions, r.Resolve(conflict))
	}
	return resolutions
}

// Resolve resolves a single conflict.
func (r *Resolver) Resolve(conflict Conflict) Resolution {
	switch conflict.Type {
	case ConflictPort:
		return r.resolvePortConflict(conflict)
	case ConflictContainerName:
		return r.resolveContainerNameConflict(conflict)
	case ConflictNetwork:
		return Resolution{Conflict: conflict, Action: ActionSkip}
	case ConflictVolume:
		return Resolution{Conflict: conflict, Action: ActionSkip}
	default:
		return Resolution{Conflict: conflict, Action: ActionManual}
	}
}

// resolvePortConflict tries to find an alternative port.
func (r *Resolver) resolvePortConflict(c Conflict) Resolution {
	if c.Suggested == "" {
		return Resolution{Conflict: c, Action: ActionManual}
	}

	return Resolution{
		Conflict: c,
		Action:   ActionAutoFix,
		Fix: &Fix{
			Field:       fmt.Sprintf("services.%s.port", c.ServiceB),
			OldValue:    c.Resource,
			NewValue:    c.Suggested,
			Description: fmt.Sprintf("Change port from %s to %s for %s", c.Resource, c.Suggested, c.ServiceB),
		},
	}
}

// resolveContainerNameConflict suggests a renamed container.
func (r *Resolver) resolveContainerNameConflict(c Conflict) Resolution {
	suggested := c.Resource + "-2"
	field := fmt.Sprintf("services.%s.container_name", c.ServiceB)

	return Resolution{
		Conflict: c,
		Action:   ActionAutoFix,
		Fix: &Fix{
			Field:       field,
			OldValue:    c.Resource,
			NewValue:    suggested,
			Description: fmt.Sprintf("Rename container from %s to %s", c.Resource, suggested),
		},
	}
}

// SuggestContainerName suggests an alternative container name.
func SuggestContainerName(base string, existingNames map[string]string) string {
	if _, exists := existingNames[base]; !exists {
		return base
	}
	for i := 2; i <= 100; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if _, exists := existingNames[candidate]; !exists {
			return candidate
		}
	}
	return fmt.Sprintf("%s-%d", base, 101)
}

// AutoFix applies all auto-fixable resolutions to the config.
func AutoFix(cfg *config.Config, resolutions []Resolution) []Fix {
	var applied []Fix

	for _, res := range resolutions {
		if res.Action != ActionAutoFix || res.Fix == nil {
			continue
		}

		fix := res.Fix
		switch fix.Field {
		case "port":
			// Find and update the port in config
			for name, svc := range cfg.Services {
				for ver, v := range svc.Versions {
					key := fmt.Sprintf("%s %s", name, ver)
					if key == res.Conflict.ServiceB {
						v.Port = parsePort(fix.NewValue)
						svc.Versions[ver] = v
						cfg.Services[name] = svc
						applied = append(applied, *fix)
					}
				}
			}
		}
	}

	return applied
}

// parsePort parses a port string to int.
func parsePort(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
