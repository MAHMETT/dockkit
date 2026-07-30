package conflict

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/MAHMETT/dockkit/internal/config"
)

type ConflictType int

const (
	ConflictPort ConflictType = iota
	ConflictContainerName
	ConflictNetwork
	ConflictVolume
)

type ConflictSeverity int

const (
	SeverityError ConflictSeverity = iota
	SeverityWarning
)

type Conflict struct {
	Type      ConflictType
	Severity  ConflictSeverity
	ServiceA  string
	ServiceB  string
	Resource  string
	Message   string
	Suggested string
}

type Detector struct {
	config *config.Config
}

func NewDetector(cfg *config.Config) *Detector {
	return &Detector{config: cfg}
}

func (d *Detector) Detect() []Conflict {
	var conflicts []Conflict
	conflicts = append(conflicts, d.detectPortConflicts()...)
	conflicts = append(conflicts, d.detectContainerNameConflicts()...)
	return conflicts
}

func (d *Detector) detectPortConflicts() []Conflict {
	var conflicts []Conflict
	portMap := map[int]string{}

	for name, svc := range d.config.Services {
		for ver, cfg := range svc.Versions {
			if !cfg.Enabled {
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

func (d *Detector) detectContainerNameConflicts() []Conflict {
	var conflicts []Conflict
	_ = strings.TrimSpace
	return conflicts
}

func (d *Detector) SuggestPort(port int) string {
	for offset := 1; offset <= 100; offset++ {
		candidate := port + offset
		if !isPortOccupied(candidate) && !d.isPortUsedByService(candidate) {
			return fmt.Sprintf("%d", candidate)
		}
	}
	return ""
}

func isPortOccupied(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

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
