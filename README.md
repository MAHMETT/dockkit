# dockkit

> Docker development infrastructure manager — interactive TUI tool for managing database, cache, storage, and search services.

## Features

- **Interactive TUI** — Modern terminal interface with keyboard & mouse support
- **Multi-service** — PostgreSQL, MySQL, MariaDB, Redis, MongoDB, MinIO, Elasticsearch, Memcached
- **Multi-version** — Run multiple versions of the same service simultaneously
- **One-click setup** — Configure and deploy services in seconds
- **Conflict detection** — Auto-detect port and container name conflicts
- **Config management** — Edit configs directly from the TUI
- **Log viewer** — Stream and filter container logs
- **Template system** — Use built-in or create custom service templates

## Installation

```bash
go install github.com/MAHMETT/dockkit@latest
```

## Quick Start

```bash
# Initialize config
dockkit init

# Launch TUI
dockkit

# Or use CLI commands
dockkit list
dockkit up postgresql-16
dockkit logs postgresql-16
```

## Requirements

- Go 1.25+
- Docker Engine / Docker Desktop
- `docker compose` v2 plugin

## Architecture

Bottom-up layer model:

```
Layer 6: CLI Commands          → dockkit up/down/list/logs
Layer 5: TUI Screens          → dashboard, detail, wizard, logs
Layer 4: TUI Framework        → Bubble Tea skeleton, styles, components
Layer 3: Conflict Detection   → port/name conflict engine
Layer 2: Docker Core          → SDK wrapper, container operations
Layer 1: Config & Templates   → data structures, config, templates
Layer 0: Foundation           → go.mod, main.go, Makefile
```

See [DESIGN.md](DESIGN.md) for detailed technical design.

## Development

```bash
git clone https://github.com/MAHMETT/dockkit.git
cd dockkit
go mod download
make dev
```

### Build

```bash
make build       # → bin/dockkit
make test        # run tests
make lint        # go vet
make clean       # remove binaries
```

### Project Structure

```
dockkit/
├── cmd/                  # Cobra commands
├── internal/
│   ├── config/           # Config load/save/validate
│   ├── templates/        # Template system
│   ├── docker/           # Docker SDK wrapper
│   ├── registry/         # Docker Hub API
│   ├── conflict/         # Conflict detection
│   ├── errors/           # Error types
│   └── tui/              # Bubble Tea TUI
│       ├── screens/      # Screen models
│       ├── components/   # Reusable UI components
│       └── messages/     # Message types
├── DESIGN.md             # Technical design
├── RULES.md              # Coding rules
├── PRD.md                # Product requirements
└── .opencode/AGENTS.md   # AI agent config
```

## Docs

- [PRD.md](PRD.md) — Product requirements (v2.1)
- [DESIGN.md](DESIGN.md) — Technical design & architecture
- [RULES.md](RULES.md) — Coding rules & conventions

## License

MIT
