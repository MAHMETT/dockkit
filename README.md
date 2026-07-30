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

## Development

```bash
git clone https://github.com/MAHMETT/dockkit.git
cd dockkit
go mod download
make dev
```

## License

MIT
