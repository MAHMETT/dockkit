package conflict

import (
	"fmt"
	"net"
	"time"

	"github.com/MAHMETT/dockkit/internal/config"
)

// Detector scans config for conflicts.
type Detector struct {
	config *config.Config
}

// NewDetector creates a new conflict detector.
func NewDetector(cfg *config.Config) *Detector {
	return &Detector{config: cfg}
}

// Detect scans the config and returns all conflicts.
func (d *Detector) Detect() ConflictList {
	var conflicts ConflictList
	conflicts = append(conflicts, d.detectPortConflicts()...)
	conflicts = append(conflicts, d.detectContainerNameConflicts()...)
	conflicts = append(conflicts, d.detectDisabledServiceWarnings()...)
	return conflicts
}

// detectPortConflicts checks for port conflicts between enabled services.
func (d *Detector) detectPortConflicts() ConflictList {
	var conflicts ConflictList
	portMap := map[int]string{} // port -> "service version"

	for name, svc := range d.config.Services {
		for ver, cfg := range svc.Versions {
			if !cfg.Enabled {
				continue
			}
			if cfg.Port == 0 {
				continue
			}

			key := fmt.Sprintf("%s %s", name, ver)
			if existing, ok := portMap[cfg.Port]; ok {
				conflicts = append(conflicts, Conflict{
					Type:      ConflictPort,
					Severity:  SeverityError,
					ServiceA:  existing,
					ServiceB:  key,
					Resource:  fmt.Sprintf("%d", cfg.Port),
					Message:   fmt.Sprintf("Port %d already used by %s", cfg.Port, existing),
					Suggested: d.SuggestPort(cfg.Port),
				})
			} else {
				portMap[cfg.Port] = key
			}
		}
	}

	return conflicts
}

// detectContainerNameConflicts checks for container name conflicts.
func (d *Detector) detectContainerNameConflicts() ConflictList {
	var conflicts ConflictList
	nameMap := map[string]string{} // container_name -> "service version"

	for name, svc := range d.config.Services {
		for ver, cfg := range svc.Versions {
			if !cfg.Enabled {
				continue
			}
			if cfg.ContainerName == "" {
				continue
			}

			key := fmt.Sprintf("%s %s", name, ver)
			if existing, ok := nameMap[cfg.ContainerName]; ok {
				conflicts = append(conflicts, Conflict{
					Type:      ConflictContainerName,
					Severity:  SeverityError,
					ServiceA:  existing,
					ServiceB:  key,
					Resource:  cfg.ContainerName,
					Message:   fmt.Sprintf("Container name %q used by both %s and %s", cfg.ContainerName, existing, key),
					Suggested: SuggestContainerName(cfg.ContainerName, nameMap),
				})
			} else {
				nameMap[cfg.ContainerName] = key
			}
		}
	}

	return conflicts
}

// detectDisabledServiceWarnings warns about disabled services.
func (d *Detector) detectDisabledServiceWarnings() ConflictList {
	var conflicts ConflictList

	for name, svc := range d.config.Services {
		for ver, cfg := range svc.Versions {
			if !cfg.Enabled {
				key := fmt.Sprintf("%s %s", name, ver)
				conflicts = append(conflicts, Conflict{
					Type:     ConflictDisabled,
					Severity: SeverityWarning,
					ServiceA: key,
					Message:  fmt.Sprintf("Service %s is disabled", key),
				})
			}
		}
	}

	return conflicts
}

// SuggestPort finds the next available port starting from port+1.
func (d *Detector) SuggestPort(port int) string {
	for offset := 1; offset <= 100; offset++ {
		candidate := port + offset
		if candidate > 65535 {
			return ""
		}
		if candidate < 1024 {
			continue
		}
		if !isPortOccupied(candidate) && !d.isPortUsedByService(candidate) {
			return fmt.Sprintf("%d", candidate)
		}
	}
	return ""
}

// isPortOccupied checks if a port is in use on the host.
func isPortOccupied(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// isPortUsedByService checks if a port is configured for any enabled service.
func (d *Detector) isPortUsedByService(port int) bool {
	for _, svc := range d.config.Services {
		for _, cfg := range svc.Versions {
			if cfg.Enabled && cfg.Port == port {
				return true
			}
		}
	}
	return false
}
