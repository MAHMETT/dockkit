# RULES.md — dockkit Coding Rules

## Go Code Style

### Naming

```go
// Exported: PascalCase
type ServiceConfig struct { ... }
func (c *Config) Validate() error { ... }

// Unexported: camelCase
func parseContainerState(c container.Summary) ServiceState { ... }
var portMap = map[int]string{}

// Interfaces: -er suffix
type DockerClient interface { ... }
type TemplateLoader interface { ... }

// Constants: PascalCase or camelCase
const MaxLogLines = 1000
const pollActiveInterval = 5 * time.Second
```

### File Organization

```
One concept per file.
Package-level concerns in package.go if small enough.

internal/config/
  config.go      → Config struct + Load/Save
  paths.go       → Path helpers
  validate.go    → Validation logic
  defaults.go    → Default values

NOT:
  config.go      → Everything in one 500-line file
```

### Error Handling

```go
// ALWAYS wrap errors with context
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config %s: %w", path, err)
    }
    // ...
}

// Use custom error types for user-facing errors
type DockkitError struct {
    Code    string
    Message string
    Detail  string
}

// Check errors with errors.Is/errors.As
if errors.Is(err, ErrPortConflict) {
    // show resolution dialog
}
```

### Comments

```go
// Good: explain WHY, not WHAT
// Batch API call: single ContainerList instead of N Inspect calls
containers, err := client.ContainerList(ctx, container.ListOptions{All: true})

// Bad: redundant comment
// List all containers
containers, err := client.ContainerList(ctx, container.ListOptions{All: true})

// Package doc: one line
// Package config manages dockkit configuration files.

// Exported functions: always have doc comment
// Load reads the config file from the given path.
func Load(path string) (*Config, error) { ... }
```

### Imports

```go
// Standard library
import (
    "context"
    "fmt"
    "os"
)

// External packages
import (
    tea "charm.land/bubbletea/v2"
    "github.com/spf13/cobra"
)

// Internal packages
import (
    "github.com/MAHMETT/dockkit/internal/config"
    "github.com/MAHMETT/dockkit/internal/docker"
)

// Grouped, sorted, no blank lines between groups in small files
// Blank lines between groups in large files for readability
```

---

## Project Rules

### Layer Dependency

```
Layer N can ONLY import Layer N-1 and below.

Layer 5 (Screens)  → Layer 4, 3, 2, 1
Layer 4 (TUI)      → Layer 3, 2, 1
Layer 3 (Conflict)  → Layer 2, 1
Layer 2 (Docker)    → Layer 1
Layer 1 (Config)    → (none, pure Go)
Layer 0 (Foundation)→ (none)

NEVER: Layer 1 importing Layer 2
NEVER: Layer 2 importing Layer 4
```

### File Naming

```
Go files:     lowercase, underscore (config.go, service_detail.go)
Test files:   *_test.go (config_test.go)
YAML files:   lowercase, hyphen (builtins/postgresql.yaml)
Dirs:         lowercase, no underscores (internal/tui/screens/)
```

### Package Size

```
Max files per package:  10
Max lines per file:     300
Max lines per function: 50
Max function params:    5

If package exceeds 10 files → split into sub-packages
If file exceeds 300 lines → split by concern
```

### Testing

```go
// Unit test: fast, no external deps
func TestConfigLoad(t *testing.T) {
    cfg, err := Load("testdata/config.yaml")
    if err != nil {
        t.Fatalf("Load() error = %v", err)
    }
    if cfg.General.Timezone != "Asia/Jakarta" {
        t.Errorf("Timezone = %s, want Asia/Jakarta", cfg.General.Timezone)
    }
}

// Table-driven tests for multiple cases
func TestSuggestPort(t *testing.T) {
    tests := []struct {
        name     string
        port     int
        occupied []int
        want     string
    }{
        {"available", 5432, nil, "5432"},
        {"occupied", 5432, []int{5432}, "5433"},
        {"both occupied", 5432, []int{5432, 5433}, "5434"},
    }
    // ...
}

// Test file location: same package as source
internal/config/config.go
internal/config/config_test.go
```

### Naming Test Functions

```go
func Test_FunctionName_Scenario(t *testing.T) { ... }
func TestConfigLoad_FileNotFound(t *testing.T) { ... }
func TestSuggestPort_AllOccupied(t *testing.T) { ... }
func TestConflictDetector_DuplicatePorts(t *testing.T) { ... }
```

---

## Commit Rules

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

```
feat:      New feature
fix:       Bug fix
refactor:  Code restructure (no behavior change)
test:      Adding tests
docs:      Documentation only
chore:     Build, deps, config
perf:      Performance improvement
style:     Formatting (no logic change)
```

### Scope (optional)

```
config, templates, docker, registry, conflict, tui, cmd, docs
```

### Examples

```
feat(docker): add batch container listing
fix(config): handle missing timezone field
refactor(templates): extract interpolator to separate file
test(conflict): add port conflict detection tests
docs: update DESIGN.md with error flow
chore: bump bubbletea to v2.0.8
```

### Rules

- Description: imperative mood, lowercase, no period
- Max 72 chars for subject line
- Body: explain what and why, not how
- Reference issues: "Closes #123"

---

## Performance Rules

### Budgets

```
TUI launch:      < 500ms cold, < 200ms warm
Docker status:   < 500ms (batch API)
Docker Hub:      < 2s (with cache)
Template render: < 50ms
Port probe:      < 100ms per port
Config save:     < 50ms
Memory:          < 20MB for 10 services
Goroutines:      < 10 total
```

### Rules

- Never block main goroutine (use tea.Cmd)
- Batch Docker API calls (1 call, not N)
- Cache Docker Hub responses (24h TTL)
- Debounce user input (100ms)
- Stream logs, don't buffer entire history
- Use context.WithTimeout on all external calls

---

## Security Rules

### Secrets

```
NEVER log passwords
NEVER print passwords in TUI (mask with ••••••••)
NEVER include passwords in error messages
NEVER commit .env files or config.yaml with real credentials
```

### File Permissions

```
config.yaml:      0600
.env files:       0600
templates/*.yaml: 0644
data directories: 0700
```

### Input Validation

```
Validate all user input before use
Port numbers: 1024-63553
Container names: Docker naming rules
Template variables: strict pattern matching
No shell execution of user input
```
