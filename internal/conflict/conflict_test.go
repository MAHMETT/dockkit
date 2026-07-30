package conflict

import (
	"testing"

	"github.com/MAHMETT/dockkit/internal/config"
)

func TestDetect_NoConflicts(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"postgresql": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"16": {Enabled: true, Port: 5432, ContainerName: "pg-16", Image: "postgres:16"},
				},
			},
			"redis": {
				Prefix: "REDIS",
				Versions: map[string]config.ServiceVersion{
					"7": {Enabled: true, Port: 6379, ContainerName: "redis-7", Image: "redis:7"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	conflicts := detector.Detect()

	if conflicts.HasErrors() {
		t.Errorf("expected no errors, got %d", len(conflicts.Errors()))
	}
}

func TestDetect_PortConflict(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"pg16": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"16": {Enabled: true, Port: 5432, ContainerName: "pg-16", Image: "postgres:16"},
				},
			},
			"pg17": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"17": {Enabled: true, Port: 5432, ContainerName: "pg-17", Image: "postgres:17"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	conflicts := detector.Detect()

	if !conflicts.HasErrors() {
		t.Fatal("expected port conflict error")
	}

	portConflicts := conflicts.ByType()[ConflictPort]
	if len(portConflicts) == 0 {
		t.Fatal("expected port conflicts")
	}

	c := portConflicts[0]
	if c.Resource != "5432" {
		t.Errorf("expected resource '5432', got %q", c.Resource)
	}
	if c.Suggested == "" {
		t.Error("expected suggestion for alternative port")
	}
}

func TestDetect_ContainerNameConflict(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"pg16": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"16": {Enabled: true, Port: 5432, ContainerName: "shared-name", Image: "postgres:16"},
				},
			},
			"redis": {
				Prefix: "REDIS",
				Versions: map[string]config.ServiceVersion{
					"7": {Enabled: true, Port: 6379, ContainerName: "shared-name", Image: "redis:7"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	conflicts := detector.Detect()

	if !conflicts.HasErrors() {
		t.Fatal("expected container name conflict error")
	}

	nameConflicts := conflicts.ByType()[ConflictContainerName]
	if len(nameConflicts) == 0 {
		t.Fatal("expected container name conflicts")
	}

	c := nameConflicts[0]
	if c.Resource != "shared-name" {
		t.Errorf("expected resource 'shared-name', got %q", c.Resource)
	}
}

func TestDetect_DisabledServiceWarning(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"pg16": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"16": {Enabled: false, Port: 5432, ContainerName: "pg-16", Image: "postgres:16"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	conflicts := detector.Detect()

	if !conflicts.HasWarnings() {
		t.Fatal("expected disabled service warning")
	}
}

func TestDetect_NoPortZeroConflict(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"redis": {
				Prefix: "REDIS",
				Versions: map[string]config.ServiceVersion{
					"7": {Enabled: true, Port: 0, ContainerName: "redis-7", Image: "redis:7"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	conflicts := detector.Detect()

	// Port 0 should not cause a conflict (not configured)
	for _, c := range conflicts {
		if c.Type == ConflictPort && c.Severity == SeverityError {
			t.Error("port 0 should not cause a conflict")
		}
	}
}

func TestSuggestPort(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{},
	}

	detector := NewDetector(cfg)

	// SuggestPort should return a port string
	suggested := detector.SuggestPort(5432)
	if suggested == "" {
		t.Error("SuggestPort returned empty string")
	}

	// Suggested port should be different from input
	if suggested == "5432" {
		t.Error("suggested port should not be the same as input")
	}
}

func TestSuggestContainerName(t *testing.T) {
	existing := map[string]string{
		"my-container": "service",
	}

	name := SuggestContainerName("my-container", existing)
	if name == "my-container" {
		t.Error("should suggest alternative name")
	}
	if name != "my-container-2" {
		t.Errorf("expected 'my-container-2', got %q", name)
	}
}

func TestSuggestContainerName_NoConflict(t *testing.T) {
	existing := map[string]string{}

	name := SuggestContainerName("my-container", existing)
	if name != "my-container" {
		t.Errorf("expected 'my-container', got %q", name)
	}
}

func TestConflictList_HasErrors(t *testing.T) {
	cl := ConflictList{
		{Severity: SeverityWarning},
		{Severity: SeverityError},
	}

	if !cl.HasErrors() {
		t.Error("expected HasErrors to return true")
	}
}

func TestConflictList_HasWarnings(t *testing.T) {
	cl := ConflictList{
		{Severity: SeverityError},
		{Severity: SeverityWarning},
	}

	if !cl.HasWarnings() {
		t.Error("expected HasWarnings to return true")
	}
}

func TestConflictList_Errors(t *testing.T) {
	cl := ConflictList{
		{Severity: SeverityWarning, Message: "warning"},
		{Severity: SeverityError, Message: "error"},
		{Severity: SeverityError, Message: "error2"},
	}

	errors := cl.Errors()
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errors))
	}
}

func TestConflictList_Warnings(t *testing.T) {
	cl := ConflictList{
		{Severity: SeverityWarning, Message: "warning"},
		{Severity: SeverityError, Message: "error"},
	}

	warnings := cl.Warnings()
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(warnings))
	}
}

func TestConflictList_ByType(t *testing.T) {
	cl := ConflictList{
		{Type: ConflictPort, Severity: SeverityError},
		{Type: ConflictContainerName, Severity: SeverityError},
		{Type: ConflictPort, Severity: SeverityWarning},
	}

	byType := cl.ByType()
	if len(byType[ConflictPort]) != 2 {
		t.Errorf("expected 2 port conflicts, got %d", len(byType[ConflictPort]))
	}
	if len(byType[ConflictContainerName]) != 1 {
		t.Errorf("expected 1 container name conflict, got %d", len(byType[ConflictContainerName]))
	}
}

func TestConflictList_Messages(t *testing.T) {
	cl := ConflictList{
		{Severity: SeverityError, Message: "error 1"},
		{Severity: SeverityWarning, Message: "warning 1"},
	}

	msgs := cl.Messages()
	if msgs == "" {
		t.Error("expected non-empty messages")
	}
}

func TestConflict_Error(t *testing.T) {
	c := Conflict{
		Type:     ConflictPort,
		Severity: SeverityError,
		Message:  "port conflict",
	}

	if c.Error() != "port conflict" {
		t.Errorf("expected 'port conflict', got %q", c.Error())
	}
}

func TestConflict_HasSuggestion(t *testing.T) {
	c1 := Conflict{Suggested: "5433"}
	c2 := Conflict{Suggested: ""}

	if !c1.HasSuggestion() {
		t.Error("expected HasSuggestion to return true")
	}
	if c2.HasSuggestion() {
		t.Error("expected HasSuggestion to return false")
	}
}

func TestResolve_PortConflict(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		General: config.GeneralConfig{
			Timezone:       "UTC",
			DefaultNetwork: "test-net",
		},
		Services: map[string]config.Service{
			"pg16": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"16": {Enabled: true, Port: 5432, ContainerName: "pg-16", Image: "postgres:16"},
				},
			},
			"pg17": {
				Prefix: "PG",
				Versions: map[string]config.ServiceVersion{
					"17": {Enabled: true, Port: 5432, ContainerName: "pg-17", Image: "postgres:17"},
				},
			},
		},
	}

	detector := NewDetector(cfg)
	resolver := NewResolver(detector)

	conflicts := detector.Detect()
	resolutions := resolver.ResolveAll(conflicts)

	hasAutoFix := false
	for _, res := range resolutions {
		if res.Action == ActionAutoFix {
			hasAutoFix = true
			if res.Fix == nil {
				t.Error("expected fix to be non-nil for auto-fix")
			}
		}
	}

	if !hasAutoFix {
		t.Error("expected at least one auto-fix resolution")
	}
}

func TestIsPortOccupied(t *testing.T) {
	// Port 1 (privileged) should not be occupied in test env
	// This is a basic smoke test
	occupied := isPortOccupied(1)
	_ = occupied // result depends on system
}

func TestConflictType_String(t *testing.T) {
	tests := []struct {
		ct   ConflictType
		want string
	}{
		{ConflictPort, "port"},
		{ConflictContainerName, "container_name"},
		{ConflictNetwork, "network"},
		{ConflictVolume, "volume"},
		{ConflictType(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.ct.String(); got != tt.want {
			t.Errorf("ConflictType(%d).String() = %q, want %q", tt.ct, got, tt.want)
		}
	}
}

func TestConflictSeverity_String(t *testing.T) {
	tests := []struct {
		cs   ConflictSeverity
		want string
	}{
		{SeverityError, "error"},
		{SeverityWarning, "warning"},
		{ConflictSeverity(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.cs.String(); got != tt.want {
			t.Errorf("ConflictSeverity(%d).String() = %q, want %q", tt.cs, got, tt.want)
		}
	}
}
