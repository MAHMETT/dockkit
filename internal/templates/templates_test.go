package templates

import (
	"strings"
	"testing"

	"github.com/MAHMETT/dockkit/internal/config"
)

func TestLoadBuiltin(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"postgresql", false},
		{"mysql", false},
		{"mariadb", false},
		{"redis", false},
		{"mongodb", false},
		{"minio", false},
		{"elasticsearch", false},
		{"memcached", false},
		{"nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := LoadBuiltin(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBuiltin(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tmpl == nil {
					t.Errorf("LoadBuiltin(%q) returned nil template", tt.name)
					return
				}
				if tmpl.Name == "" {
					t.Errorf("template Name is empty")
				}
				if len(tmpl.Versions) == 0 {
					t.Errorf("template has no versions")
				}
			}
		})
	}
}

func TestListBuiltin(t *testing.T) {
	names, err := ListBuiltin()
	if err != nil {
		t.Fatalf("ListBuiltin() error = %v", err)
	}
	if len(names) < 8 {
		t.Errorf("ListBuiltin() returned %d names, want >= 8", len(names))
	}
	// Check specific names exist
	nameSet := map[string]bool{}
	for _, n := range names {
		nameSet[n] = true
	}
	for _, want := range []string{"postgresql", "mysql", "redis", "mongodb"} {
		if !nameSet[want] {
			t.Errorf("ListBuiltin() missing %q", want)
		}
	}
}

func TestInterpolator(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Timezone:       "Asia/Jakarta",
			DefaultNetwork: "my-net",
		},
	}

	interp := NewInterpolator(cfg, "postgresql", "16")
	interp.SetVar("CONFIG_PORT", "5432")
	interp.SetVar("CONFIG_USER", "admin")
	interp.SetVar("CONFIG_PASSWORD", "secret")
	interp.SetVar("CONFIG_DATABASE", "mydb")

	tests := []struct {
		input string
		want  string
	}{
		{"${GENERAL_TIMEZONE}", "Asia/Jakarta"},
		{"${GENERAL_DEFAULT_NETWORK}", "my-net"},
		{"${CONFIG_PORT}", "5432"},
		{"${CONFIG_USER}", "admin"},
		{"no vars here", "no vars here"},
		{"${CONFIG_PORT}:5432", "5432:5432"},
		{"${UNRESOLVED_VAR}", "${UNRESOLVED_VAR}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := interp.Interpolate(tt.input)
			if got != tt.want {
				t.Errorf("Interpolate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestInterpolator_InterpolateMap(t *testing.T) {
	interp := NewInterpolator(nil, "redis", "7")
	interp.SetVar("CONFIG_PORT", "6379")
	interp.SetVar("GENERAL_TIMEZONE", "UTC")

	input := map[string]string{
		"HOST": "${CONFIG_PORT}",
		"TZ":   "${GENERAL_TIMEZONE}",
		"RAW":  "no-interpolation",
	}

	got := interp.InterpolateMap(input)
	if got["HOST"] != "6379" {
		t.Errorf("HOST = %q, want %q", got["HOST"], "6379")
	}
	if got["TZ"] != "UTC" {
		t.Errorf("TZ = %q, want %q", got["TZ"], "UTC")
	}
	if got["RAW"] != "no-interpolation" {
		t.Errorf("RAW = %q, want %q", got["RAW"], "no-interpolation")
	}
}

func TestInterpolator_InterpolateSlice(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	interp.SetVar("CONFIG_PORT", "5432")

	input := []string{"${CONFIG_PORT}:5432", "data:/var/lib/test"}
	got := interp.InterpolateSlice(input)

	if got[0] != "5432:5432" {
		t.Errorf("got[0] = %q, want %q", got[0], "5432:5432")
	}
	if got[1] != "data:/var/lib/test" {
		t.Errorf("got[1] = %q, want %q", got[1], "data:/var/lib/test")
	}
}

func TestInterpolator_InterpolateSlice_Nil(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	got := interp.InterpolateSlice(nil)
	if got != nil {
		t.Errorf("InterpolateSlice(nil) = %v, want nil", got)
	}
}

func TestInterpolator_InterpolateMap_Nil(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	got := interp.InterpolateMap(nil)
	if got != nil {
		t.Errorf("InterpolateMap(nil) = %v, want nil", got)
	}
}

func TestInterpolator_UnresolvedVars(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	interp.SetVar("CONFIG_PORT", "5432")

	unresolved := interp.UnresolvedVars("port ${CONFIG_PORT} and ${MISSING} and ${ALSO_MISSING}")
	if len(unresolved) != 2 {
		t.Errorf("UnresolvedVars() returned %d vars, want 2", len(unresolved))
	}
}

func TestInterpolator_HasUnresolved(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	interp.SetVar("CONFIG_PORT", "5432")

	if interp.HasUnresolved("port ${CONFIG_PORT}") {
		t.Error("HasUnresolved() = true, want false")
	}
	if !interp.HasUnresolved("port ${MISSING}") {
		t.Error("HasUnresolved() = false, want true")
	}
}

func TestInterpolator_AvailableVars(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Timezone: "UTC",
		},
	}
	interp := NewInterpolator(cfg, "pg", "16")
	vars := interp.AvailableVars()

	if vars["GENERAL_TIMEZONE"] != "UTC" {
		t.Errorf("GENERAL_TIMEZONE = %q, want %q", vars["GENERAL_TIMEZONE"], "UTC")
	}
	if vars["SERVICE_NAME"] != "pg" {
		t.Errorf("SERVICE_NAME = %q, want %q", vars["SERVICE_NAME"], "pg")
	}
}

func TestInterpolator_SetServiceConfig(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	sv := &config.ServiceVersion{
		Port:          5432,
		User:          "admin",
		Password:      "secret",
		Database:      "mydb",
		Image:         "postgres:16",
		ContainerName: "my-pg",
	}
	interp.SetServiceConfig(sv)

	if interp.vars["CONFIG_PORT"] != "5432" {
		t.Errorf("CONFIG_PORT = %q, want %q", interp.vars["CONFIG_PORT"], "5432")
	}
	if interp.vars["CONFIG_USER"] != "admin" {
		t.Errorf("CONFIG_USER = %q, want %q", interp.vars["CONFIG_USER"], "admin")
	}
}

func TestInterpolator_SetServiceConfig_Nil(t *testing.T) {
	interp := NewInterpolator(nil, "test", "1")
	interp.SetServiceConfig(nil)
	// Should not panic
}

func TestRender(t *testing.T) {
	tmpl, err := LoadBuiltin("postgresql")
	if err != nil {
		t.Fatalf("LoadBuiltin() error = %v", err)
	}

	opts := RenderOptions{
		ServiceName:   "postgresql",
		Version:       "16",
		Port:          5433,
		User:          "admin",
		Password:      "secret",
		Database:      "mydb",
		ContainerName: "dockkit-postgresql-16",
		Timezone:      "Asia/Jakarta",
		Network:       "dockkit-network",
	}

	compose, err := Render(tmpl, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if compose.Version != "3.8" {
		t.Errorf("Version = %q, want %q", compose.Version, "3.8")
	}

	svc, ok := compose.Services["postgresql"]
	if !ok {
		t.Fatal("service 'postgresql' not found in compose")
	}

	if svc.Image != "postgres:16-alpine" {
		t.Errorf("Image = %q, want %q", svc.Image, "postgres:16-alpine")
	}
	if svc.ContainerName != "dockkit-postgresql-16" {
		t.Errorf("ContainerName = %q, want %q", svc.ContainerName, "dockkit-postgresql-16")
	}
	if svc.Restart != "unless-stopped" {
		t.Errorf("Restart = %q, want %q", svc.Restart, "unless-stopped")
	}
	if svc.Environment["POSTGRES_USER"] != "admin" {
		t.Errorf("POSTGRES_USER = %q, want %q", svc.Environment["POSTGRES_USER"], "admin")
	}
	if svc.Environment["POSTGRES_PASSWORD"] != "secret" {
		t.Errorf("POSTGRES_PASSWORD = %q, want %q", svc.Environment["POSTGRES_PASSWORD"], "secret")
	}
	if len(svc.Ports) == 0 {
		t.Error("Ports is empty")
	}
	if len(svc.Volumes) == 0 {
		t.Error("Volumes is empty")
	}
	if svc.Healthcheck == nil {
		t.Error("Healthcheck is nil")
	}
}

func TestRender_NilTemplate(t *testing.T) {
	_, err := Render(nil, RenderOptions{})
	if err == nil {
		t.Error("Render(nil) should return error")
	}
}

func TestRender_InvalidVersion(t *testing.T) {
	tmpl, _ := LoadBuiltin("postgresql")
	_, err := Render(tmpl, RenderOptions{Version: "999"})
	if err == nil {
		t.Error("Render() with invalid version should return error")
	}
}

func TestRenderYAML(t *testing.T) {
	tmpl, _ := LoadBuiltin("redis")
	opts := RenderOptions{
		ServiceName:   "redis",
		Version:       "7",
		Port:          6379,
		ContainerName: "dockkit-redis-7",
	}

	data, err := RenderYAML(tmpl, opts)
	if err != nil {
		t.Fatalf("RenderYAML() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "redis:7-alpine") {
		t.Error("YAML doesn't contain redis:7-alpine")
	}
	if !strings.Contains(yamlStr, "dockkit-redis-7") {
		t.Error("YAML doesn't contain dockkit-redis-7")
	}
}

func TestRenderToString(t *testing.T) {
	tmpl, _ := LoadBuiltin("mysql")
	opts := RenderOptions{
		ServiceName:   "mysql",
		Version:       "8.0",
		Port:          3306,
		ContainerName: "dockkit-mysql-8",
	}

	str, err := RenderToString(tmpl, opts)
	if err != nil {
		t.Fatalf("RenderToString() error = %v", err)
	}
	if str == "" {
		t.Error("RenderToString() returned empty string")
	}
}

func TestSanitizeServiceKey(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"PostgreSQL", "postgresql"},
		{"My SQL", "my-sql"},
		{"my_service", "my-service"},
		{"Redis", "redis"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeServiceKey(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeServiceKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tmpl, _ := LoadBuiltin("postgresql")
	errs := Validate(tmpl)
	if errs.HasErrors() {
		t.Errorf("Validate(postgresql) errors = %v", errs)
	}
}

func TestValidate_Nil(t *testing.T) {
	errs := Validate(nil)
	if !errs.HasErrors() {
		t.Error("Validate(nil) should have errors")
	}
}

func TestValidate_InvalidCategory(t *testing.T) {
	tmpl := &Template{
		Name:        "test",
		Description: "test template",
		Category:    "invalid",
		Versions: []VersionEntry{
			{Key: "1", Image: "test:1", DefaultPort: 5000},
		},
	}
	errs := Validate(tmpl)
	if !errs.HasErrors() {
		t.Error("Validate() should detect invalid category")
	}
}

func TestValidate_DuplicatePorts(t *testing.T) {
	tmpl := &Template{
		Name:        "test",
		Description: "test",
		Category:    "database",
		Versions: []VersionEntry{
			{Key: "1", Image: "test:1", DefaultPort: 5000},
			{Key: "2", Image: "test:2", DefaultPort: 5000},
		},
	}
	errs := Validate(tmpl)
	if !errs.HasErrors() {
		t.Error("Validate() should detect duplicate ports")
	}
}

func TestValidateConfigField(t *testing.T) {
	field := ConfigField{
		Key:      "port",
		Label:    "Port",
		Type:     "number",
		Required: true,
	}

	if err := ValidateConfigField(field, ""); err == nil {
		t.Error("ValidateConfigField() should fail on empty required field")
	}
	if err := ValidateConfigField(field, "abc"); err == nil {
		t.Error("ValidateConfigField() should fail on non-numeric input")
	}
	if err := ValidateConfigField(field, "5000"); err != nil {
		t.Errorf("ValidateConfigField() error = %v", err)
	}
}

func TestFindVersion(t *testing.T) {
	tmpl, _ := LoadBuiltin("postgresql")
	v := findVersion(tmpl, "16")
	if v == nil {
		t.Fatal("findVersion(16) returned nil")
	}
	if v.Image != "postgres:16-alpine" {
		t.Errorf("Image = %q, want %q", v.Image, "postgres:16-alpine")
	}

	v = findVersion(tmpl, "999")
	if v != nil {
		t.Error("findVersion(999) should return nil")
	}
}
