package conflict

import (
	"fmt"
	"strconv"
	"strings"

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
	case ConflictNetwork, ConflictVolume, ConflictDisabled:
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
			Field:       fmt.Sprintf("services.%s.versions.port", c.ServiceB),
			OldValue:    c.Resource,
			NewValue:    c.Suggested,
			Description: fmt.Sprintf("Change port from %s to %s for %s", c.Resource, c.Suggested, c.ServiceB),
		},
	}
}

// resolveContainerNameConflict suggests a renamed container.
func (r *Resolver) resolveContainerNameConflict(c Conflict) Resolution {
	// Build list of existing container names
	existingNames := map[string]bool{}
	for _, svc := range r.detector.config.Services {
		for _, cfg := range svc.Versions {
			if cfg.ContainerName != "" {
				existingNames[cfg.ContainerName] = true
			}
		}
	}

	suggested := SuggestAvailableName(c.Resource, existingNames)
	field := fmt.Sprintf("services.%s.versions.container_name", c.ServiceB)

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

// SuggestContainerName suggests an alternative container name that doesn't conflict.
func SuggestContainerName(base string, existingNames map[string]string) string {
	boolMap := make(map[string]bool, len(existingNames))
	for k := range existingNames {
		boolMap[k] = true
	}
	return SuggestAvailableName(base, boolMap)
}

// SuggestAvailableName finds a name that doesn't exist in the given set.
func SuggestAvailableName(base string, existing map[string]bool) string {
	if !existing[base] {
		return base
	}
	for i := 2; i <= 100; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !existing[candidate] {
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

		// Parse field path: "services.<name>.versions.<ver>.<field>"
		parts := strings.Split(fix.Field, ".")
		if len(parts) < 5 || parts[0] != "services" {
			continue
		}

		serviceName := parts[1]
		versionKey := parts[3]
		fieldName := parts[4]

		svc, ok := cfg.Services[serviceName]
		if !ok {
			continue
		}

		ver, ok := svc.Versions[versionKey]
		if !ok {
			continue
		}

		switch fieldName {
		case "port":
			port, err := strconv.Atoi(fix.NewValue)
			if err == nil {
				ver.Port = port
				svc.Versions[versionKey] = ver
				cfg.Services[serviceName] = svc
				applied = append(applied, *fix)
			}
		case "container_name":
			ver.ContainerName = fix.NewValue
			svc.Versions[versionKey] = ver
			cfg.Services[serviceName] = svc
			applied = append(applied, *fix)
		}
	}

	return applied
}
