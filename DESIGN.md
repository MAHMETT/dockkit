# DESIGN.md — dockkit Technical Design

## Architecture: Bottom-Up Layer Model

```
┌─────────────────────────────────────────────────────┐
│  Layer 6: CLI Commands                              │
│  dockkit up/down/restart/list/logs/setup/templates  │
├─────────────────────────────────────────────────────┤
│  Layer 5: TUI Screens                              │
│  dashboard, detail, picker, wizard, logs, templates │
├─────────────────────────────────────────────────────┤
│  Layer 4: TUI Framework                            │
│  app.go, styles, keys, messages, components         │
├─────────────────────────────────────────────────────┤
│  Layer 3: Conflict Detection                       │
│  detector, resolver, portcheck                      │
├─────────────────────────────────────────────────────┤
│  Layer 2: Docker Core                              │
│  client, containers, images, logs, compose, network │
├─────────────────────────────────────────────────────┤
│  Layer 1: Config & Templates                       │
│  config, templates, errors, registry                │
├─────────────────────────────────────────────────────┤
│  Layer 0: Foundation                               │
│  go.mod, main.go, cmd/root.go, Makefile            │
└─────────────────────────────────────────────────────┘
```

**Rule:** Each layer ONLY depends on layers below it. Never upward.

---

## Layer 1: Config & Templates

### Data Flow

```
User Input
    ↓
config.yaml (disk) ←→ Config struct (memory)
    ↓
Template YAML (embedded) ←→ Template struct (memory)
    ↓
Renderer: Config + Template → ComposeConfig (YAML bytes)
```

### Config Loading Priority

```
1. CLI flag (--config /path/to/file)
2. Environment variable (DOCKKIT_CONFIG)
3. Default path (~/.config/dockkit/config.yaml)
4. Built-in defaults (if no file exists)
```

### Template Rendering Flow

```
Template (YAML)
    ↓
Loader: parse YAML → Template struct
    ↓
Interpolator: replace ${VAR} with config values
    ↓
Renderer: generate docker-compose.yml content
    ↓
Writer: write to ~/.config/dockkit/services/<name>/<version>/
```

### Variable Interpolation Rules

| Pattern | Source | Example |
|---|---|---|
| `${CONFIG_<field>}` | User input from config wizard | `${CONFIG_PORT}` → `5432` |
| `${GENERAL_<key>}` | config.yaml general section | `${GENERAL_TIMEZONE}` → `Asia/Jakarta` |
| `${SERVICE_<key>}` | Service metadata | `${SERVICE_NAME}` → `postgresql` |
| `${VERSION_<key>}` | Version metadata | `${VERSION_NUMBER}` → `16` |

**Undefined variables:** Left as-is (not expanded), logged as warning.

---

## Layer 2: Docker Core

### Client Initialization

```
docker.NewClient()
    ↓
client.NewClientWithOpts(
    client.FromEnv,                          // read DOCKER_HOST etc
    client.WithAPIVersionNegotiation(),      // auto-negotiate version
)
    ↓
client.Ping(ctx) → verify Docker is reachable
    ↓
Client ready
```

### Batch Container Listing

```
// Instead of N individual Inspect calls:
containers, _ := client.ContainerList(ctx, container.ListOptions{
    All: true,  // include stopped
})

// Parse from list response (no extra API calls):
for _, c := range containers {
    state = parseContainerState(c)   // from c.State, c.Status
    health = parseHealth(c.Labels)   // from container inspect if needed
}
```

### Docker Compose Execution

```
compose.ComposeUp(serviceDir string) error
    ↓
exec.CommandContext(ctx, "docker", "compose",
    "-f", filepath.Join(serviceDir, "docker-compose.yml"),
    "--env-file", filepath.Join(serviceDir, ".env"),
    "up", "-d",
)
    ↓
stdout/stderr → capture for error messages
```

### Log Streaming

```
docker.StreamLogs(containerName string, opts LogOptions) (io.ReadCloser, error)
    ↓
client.ContainerLogs(ctx, containerID, container.LogsOptions{
    ShowStdout: true,
    ShowStderr: true,
    Follow:     opts.Follow,
    Tail:       opts.Tail,
    Since:      opts.Since,
})
    ↓
reader → demux with stdcopy.StdCopy → stdout + stderr
    ↓
pipe to TUI viewport
```

---

## Layer 3: Conflict Detection

### Detection Flow

```
Config Save / Service Start
    ↓
ConflictDetector.Detect(config)
    ↓
┌──────────────────────────────────┐
│ 1. Port vs Port (config scan)    │
│    for each enabled service:     │
│      check portMap[port]         │
│      if exists → CONFLICT        │
├──────────────────────────────────┤
│ 2. Port vs Host (TCP probe)      │
│    for each enabled service:     │
│      net.Dial("tcp", port)       │
│      if success → CONFLICT       │
├──────────────────────────────────┤
│ 3. Container Name (Docker scan)  │
│    ContainerList(All: true)      │
│    check if name exists          │
│    if not managed by us → CONFLICT│
└──────────────────────────────────┘
    ↓
[]Conflict → resolution dialog or auto-fix
```

### Port Suggestion Algorithm

```
SuggestPort(originalPort int) string:
    for offset := 1; offset <= 100; offset++:
        candidate := originalPort + offset
        if !isPortOccupied(candidate) &&
           !isPortUsedByService(candidate) &&
           candidate > 1024:
            return candidate
    return ""  // no port found
```

---

## Layer 4: TUI Framework

### State Machine

```
┌──────────┐     navigate      ┌────────────────┐
│          │ ─────────────────► │                │
│ Dashboard│ ◄───────────────── │ Service Detail │
│          │     back           │                │
└────┬─────┘                    └────────┬───────┘
     │                                   │
     │ +                                 │ l
     │                                   │
     ▼                                   ▼
┌──────────────┐                ┌────────────────┐
│              │   Enter        │                │
│Service Picker│ ──────────────►│  Logs Viewer   │
│              │                │                │
└──────┬───────┘                └────────────────┘
       │ Enter
       ▼
┌──────────────┐
│              │
│Config Wizard │ ──Ctrl+S──► Install service
│              │
└──────────────┘
```

### Screen Interface

Every screen implements:

```go
type Screen interface {
    Init() tea.Cmd                    // initial commands
    Update(msg tea.Msg) (Screen, tea.Cmd)  // handle events
    View() string                     // render
    Title() string                    // screen title for header
}
```

### Message Flow

```
User Key Press
    ↓
app.Update(msg)
    ↓
CurrentScreen.Update(msg)
    ↓
Returns: (newModel, cmd)
    ↓
cmd → tea.Cmd → async operation
    ↓
Result message → app.Update → route to screen
```

### Component Composition

```
Screen = Header + Content + Footer

Header:
  Title + Version + [?] Help + [q] Back

Content:
  (varies per screen)

Footer:
  Status bar / Key hints / Toast messages
```

---

## Layer 5: TUI Screens

### Screen Dependencies

```
dashboard ──────► Docker (Layer 2)
    │
    ├──► service_detail ──► Docker
    │
    ├──► service_picker ──► Templates (Layer 1)
    │       │
    │       └──► config_wizard ──► Config + Conflict
    │               │
    │               └──► (install service)
    │
    ├──► logs_viewer ──► Docker
    │
    ├──► template_manager ──► Templates
    │       │
    │       └──► template_editor ──► Templates
    │
    └──► version_fetcher ──► Registry (Layer 2)
```

### Dashboard Data Flow

```
dashboard.Init()
    ↓
tea.Cmd: fetchServiceStates()
    ↓
Docker client.ContainerList(All: true)
    ↓
For each container:
  match against config.Services
  build ServiceState struct
    ↓
StatusUpdateMsg{states}
    ↓
dashboard.Update(msg)
    ↓
Render: table of services with status badges
```

### Config Wizard Flow

```
config_wizard.Init()
    ↓
Load template for selected service
    ↓
Render huh.Form with config_fields
    ↓
User fills form → Validate on blur
    ↓
Ctrl+S → ConflictDetector.Detect()
    ↓
NO conflict → Render compose → Write files → Install
YES conflict → Show resolution dialog
    ↓
Toast: "Service installed successfully"
Navigate: back to dashboard
```

---

## Layer 6: CLI Commands

### Command Pattern

```
cmd/up.go:

var upCmd = &cobra.Command{
    Use:   "up [service-version]",
    Short: "Start a service",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // 1. Parse arg (e.g., "postgresql-16")
        // 2. Load config
        // 3. Pre-flight checks (conflict detection)
        // 4. docker compose up
        // 5. Wait for healthy
        // 6. Print status
    },
}
```

### Non-TUI Output Format

```
$ dockkit list
SERVICE           VERSION  STATUS    PORT    HEALTH
postgresql        16       running   5432    healthy
redis             7        running   6379    healthy
mysql             8        stopped   -       -

$ dockkit up postgresql-17
Starting PostgreSQL 17...
✓ Container created (dockkit-postgresql-17)
✓ Image pulled (postgres:17-alpine)
✓ Service started on port 5434
✓ Health check passed

PostgreSQL 17 is ready.
  URI: postgresql://postgres:postgres@localhost:5434/postgres
```

---

## Error Propagation

```
Internal Error
    ↓
Wrap with context: fmt.Errorf("loading config: %w", err)
    ↓
Internal/errors classifies: errors.Classify(err)
    ↓
User-facing message: errors.UserMessage(classifiedErr)
    ↓
Display: Toast (TUI) or stderr (CLI)
```

### Error Types

| Code | Message | Action |
|---|---|---|
| `DOCKER_NOT_RUNNING` | "Docker is not running" | Show error screen / stderr |
| `PORT_CONFLICT` | "Port X already used by Y" | Resolution dialog |
| `CONFIG_CORRUPTED` | "Config corrupted, backup at X" | Auto-backup + reset |
| `HUB_OFFLINE` | "Docker Hub unreachable" | Fallback to cache |
| `PERMISSION_DENIED` | "Permission denied" | Show fix instructions |

---

## File Structure (Final)

```
dockkit/
├── main.go
├── go.mod / go.sum
├── Makefile
├── README.md
├── DESIGN.md
├── RULES.md
├── PRD.md
│
├── .opencode/
│   └── AGENTS.md
│
├── cmd/
│   ├── root.go
│   ├── tui.go
│   ├── init.go
│   ├── list.go
│   ├── up.go
│   ├── down.go
│   ├── restart.go
│   ├── logs.go
│   ├── setup.go
│   ├── templates.go
│   └── version.go
│
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   ├── paths.go
│   │   ├── defaults.go
│   │   ├── validate.go
│   │   └── config_test.go
│   │
│   ├── templates/
│   │   ├── loader.go
│   │   ├── renderer.go
│   │   ├── interpolator.go
│   │   ├── validator.go
│   │   ├── loader_test.go
│   │   └── builtins/
│   │       ├── postgresql.yaml
│   │       ├── mysql.yaml
│   │       ├── mariadb.yaml
│   │       ├── redis.yaml
│   │       ├── mongodb.yaml
│   │       ├── minio.yaml
│   │       ├── elasticsearch.yaml
│   │       └── memcached.yaml
│   │
│   ├── docker/
│   │   ├── client.go
│   │   ├── containers.go
│   │   ├── images.go
│   │   ├── logs.go
│   │   ├── compose.go
│   │   ├── network.go
│   │   ├── health.go
│   │   └── docker_test.go
│   │
│   ├── registry/
│   │   ├── hub.go
│   │   ├── cache.go
│   │   ├── types.go
│   │   └── hub_test.go
│   │
│   ├── conflict/
│   │   ├── detector.go
│   │   ├── resolver.go
│   │   ├── portcheck.go
│   │   ├── types.go
│   │   └── detector_test.go
│   │
│   ├── errors/
│   │   ├── errors.go
│   │   └── messages.go
│   │
│   └── tui/
│       ├── app.go
│       ├── styles.go
│       ├── keys.go
│       ├── theme.go
│       │
│       ├── messages/
│       │   ├── navigation.go
│       │   ├── docker.go
│       │   ├── config.go
│       │   └── system.go
│       │
│       ├── components/
│       │   ├── status_badge.go
│       │   ├── toast.go
│       │   ├── confirm_dialog.go
│       │   ├── loading_spinner.go
│       │   ├── help_overlay.go
│       │   ├── search_bar.go
│       │   └── progress_bar.go
│       │
│       └── screens/
│           ├── dashboard.go
│           ├── service_detail.go
│           ├── service_picker.go
│           ├── config_wizard.go
│           ├── version_fetcher.go
│           ├── logs_viewer.go
│           ├── template_manager.go
│           ├── template_editor.go
│           └── error_screen.go
│
└── templates/                    # (for reference only, embedded via go:embed)
```
