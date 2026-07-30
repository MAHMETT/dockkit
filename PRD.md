# PRD: dockkit

> **Version:** 2.1
> **Last Updated:** 2026-07-30
> **Status:** Draft — Pending Review

---

## Table of Contents

1. [Overview](#overview)
2. [Glossary](#glossary)
3. [Tech Stack](#tech-stack)
4. [Installation](#installation)
5. [CLI Commands](#cli-commands)
6. [Config System](#config-system)
7. [Template System](#template-system)
8. [TUI Screens](#tui-screens)
9. [Keybindings](#keybindings)
10. [Data Model](#data-model)
11. [Architecture](#architecture)
12. [Conflict Detection & Resolution](#conflict-detection--resolution)
13. [Performance Optimization Strategy](#performance-optimization-strategy)
14. [Error Handling Matrix](#error-handling-matrix)
15. [Security Guidelines](#security-guidelines)
16. [Caching Strategy](#caching-strategy)
17. [Config Migration Strategy](#config-migration-strategy)
18. [Development Setup](#development-setup)
19. [Testing Strategy](#testing-strategy)
20. [Risk Mitigation](#risk-mitigation)
21. [MVP Scope](#mvp-scope)
22. [Implementation Roadmap](#implementation-roadmap)

---

## Overview

**dockkit** adalah CLI TUI tool untuk mengelola Docker-based development infrastructure services. Tool ini menyediakan interface interaktif untuk setup, konfigurasi, monitoring, dan orchestration berbagai service (database, cache, storage, dll) secara cepat dan terstruktur.

### Goals

- Single CLI tool untuk semua development infrastructure needs
- Zero manual docker-compose editing untuk service yang sudah ditemplate
- Interactive TUI yang modern, keyboard+mouse support
- Multi-version support untuk setiap service
- Portable, bisa dipindah antar mesin

### Non-Goals

- Tidak menggantikan Docker Desktop / Docker Engine
- Tidak managing application code, hanya infrastructure services
- Tidak production deployment tool (dev-only focus)

---

## Glossary

| Term | Definition |
|---|---|
| **Service** | A Docker-based infrastructure component (database, cache, storage, etc.) |
| **Template** | YAML definition describing how to configure and deploy a service |
| **Version** | A specific release of a service (e.g., PostgreSQL 16, Redis 7) |
| **Instance** | A configured and running version of a service |
| **Config** | User-level configuration stored in `~/.config/dockkit/` |
| **Container Name** | Unique identifier for a Docker container (e.g., `dockkit-postgresql-16`) |
| **Category** | Service grouping (database, cache, storage, search) |
| **Health** | Container health status (healthy, unhealthy, starting, none) |

---

## Tech Stack

| Component | Technology | Version | Module Path |
|---|---|---|---|
| Language | Go | **1.25+** | - |
| CLI Framework | Cobra | v1.8+ | `github.com/spf13/cobra` |
| TUI Framework | Bubble Tea | v2.0.8 | `charm.land/bubbletea/v2` |
| TUI Components | Bubbles | v2.1.1 | `charm.land/bubbles/v2` |
| TUI Forms | Huh | v2.0.3 | `charm.land/huh/v2` |
| TUI Styling | Lip Gloss | v2.0.5 | `charm.land/lipgloss/v2` |
| Config Management | Viper | v1.19+ | `github.com/spf13/viper` |
| Docker Client | moby/moby/client | v0.5.1 | `github.com/moby/moby/client` |
| Docker Types | moby/moby/api | v1.55.0 | `github.com/moby/moby/api` |
| Template Embed | embed (stdlib) | Go 1.16+ | - |
| YAML Parser | yaml.v3 | v3.0+ | `gopkg.in/yaml.v3` |
| Release | GoReleaser | v2 | - |
| Config Format | YAML | - | - |

> **Why moby/moby instead of docker/docker?**
> `github.com/docker/docker` is frozen with `+incompatible` tag. It pulls the entire monorepo. The new split modules (`github.com/moby/moby/client` + `github.com/moby/moby/api`) are the recommended approach for new projects.

---

## Installation

```bash
# Via go install
go install github.com/user/dockkit@latest

# Via Homebrew (after GoReleaser setup)
brew install user/tap/dockkit

# Via GitHub Releases
# Download binary from releases page

# Verify installation
dockkit version
```

### Requirements

- Go 1.25+ (for building from source)
- Docker Engine / Docker Desktop installed and running
- `docker compose` v2 plugin (for compose operations)

---

## CLI Commands

```
dockkit                          # Launch TUI (default)
dockkit init                     # Initialize config at ~/.config/dockkit/
dockkit list                     # List all services (non-TUI, for scripting)
dockkit up <service>             # Start a service
dockkit down <service>           # Stop a service
dockkit restart <service>        # Restart a service
dockkit logs <service>           # View service logs
dockkit setup <service>          # Setup new service (non-TUI wizard)
dockkit templates                # List available templates
dockkit template show <name>     # Show template content
dockkit template add <name>      # Add custom template
dockkit template remove <name>   # Remove custom template
dockkit version                  # Show version
```

### Command Flags

```bash
# Global flags
--config string      Config file path (default: ~/.config/dockkit/config.yaml)
--verbose            Enable debug logging
--no-color           Disable colored output

# Service flags
--version string     Service version (e.g., 16, 7, latest)
--port int           Override default port
--dry-run            Preview without executing

# Template flags
--from string        Copy from existing template
--category string    Filter by category
```

---

## Config System

### Directory Structure

```
~/.config/dockkit/
├── config.yaml                 # Main config
├── cache/                      # Docker Hub cache
│   └── tags/                   # Cached image tags
│       ├── postgresql.json
│       └── mysql.json
├── templates/                  # Custom templates (user-created)
│   ├── my-postgres.yaml
│   └── custom-redis.yaml
└── services/                   # Generated service configs
    ├── postgresql/
    │   └── 16/
    │       ├── docker-compose.yml
    │       ├── .env
    │       └── data/           # Persistent data (gitignored)
    └── redis/
        └── 7/
            ├── docker-compose.yml
            ├── .env
            └── data/
```

### config.yaml

```yaml
version: "1"

general:
  timezone: Asia/Jakarta
  default_network: dockkit-network
  auto_refresh: true
  refresh_interval: 30s

services:
  postgresql:
    prefix: PG
    versions:
      "16":
        enabled: true
        port: 5432
        container_name: dockkit-postgresql-16
        image: postgres:16-alpine
        user: postgres
        password: postgres
        database: postgres
      "17":
        enabled: false
        port: 5433
        container_name: dockkit-postgresql-17
        image: postgres:17-alpine
        user: postgres
        password: postgres
        database: postgres

  mysql:
    prefix: MYSQL
    versions:
      "8":
        enabled: true
        port: 3306
        container_name: dockkit-mysql-8
        image: mysql:8.0
        user: root
        password: mysql
        database: mysql

  redis:
    prefix: REDIS
    versions:
      "7":
        enabled: false
        port: 6379
        container_name: dockkit-redis-7
        image: redis:7-alpine
        password: ""
```

### Environment Naming Convention

```
<PREFIX>_<VERSION>_<CONFIG>
```

| Service | Prefix | Example |
|---|---|---|
| PostgreSQL | `PG` | `PG16_PORT`, `PG16_USER` |
| MySQL | `MYSQL` | `MYSQL8_PORT`, `MYSQL8_USER` |
| MariaDB | `MARIADB` | `MARIADB11_PORT` |
| Redis | `REDIS` | `REDIS7_PORT` |
| MongoDB | `MONGO` | `MONGO8_PORT` |
| MinIO | `MINIO` | `MINIO_PORT` |
| Elasticsearch | `ELASTIC` | `ELASTIC8_PORT` |
| Memcached | `MEMCACHED` | `MEMCACHED16_PORT` |

---

## Template System

### Template Format

```yaml
name: PostgreSQL
description: Advanced open source relational database
category: database
icon: "🐘"

# Available versions
versions:
  "15":
    image: postgres:15-alpine
    default_port: 5432
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${CONFIG_USER} -d ${CONFIG_DATABASE}"]
      interval: 10s
      timeout: 5s
      retries: 5
    env_vars:
      POSTGRES_USER: "${CONFIG_USER}"
      POSTGRES_PASSWORD: "${CONFIG_PASSWORD}"
      POSTGRES_DB: "${CONFIG_DATABASE}"
      TZ: "${GENERAL_TIMEZONE}"
      PGTZ: "${GENERAL_TIMEZONE}"
    volumes:
      - "./data:/var/lib/postgresql/data"
    networks:
      - "${GENERAL_DEFAULT_NETWORK}"
    ports:
      - "${CONFIG_PORT}:5432"
    command: null
    shm_size: null

  "16":
    image: postgres:16-alpine
    default_port: 5433
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${CONFIG_USER} -d ${CONFIG_DATABASE}"]
      interval: 10s
      timeout: 5s
      retries: 5
    env_vars:
      POSTGRES_USER: "${CONFIG_USER}"
      POSTGRES_PASSWORD: "${CONFIG_PASSWORD}"
      POSTGRES_DB: "${CONFIG_DATABASE}"
      TZ: "${GENERAL_TIMEZONE}"
      PGTZ: "${GENERAL_TIMEZONE}"
    volumes:
      - "./data:/var/lib/postgresql/data"
    networks:
      - "${GENERAL_DEFAULT_NETWORK}"
    ports:
      - "${CONFIG_PORT}:5432"

  "17":
    image: postgres:17-alpine
    default_port: 5434
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${CONFIG_USER} -d ${CONFIG_DATABASE}"]
      interval: 10s
      timeout: 5s
      retries: 5
    env_vars:
      POSTGRES_USER: "${CONFIG_USER}"
      POSTGRES_PASSWORD: "${CONFIG_PASSWORD}"
      POSTGRES_DB: "${CONFIG_DATABASE}"
      TZ: "${GENERAL_TIMEZONE}"
      PGTZ: "${GENERAL_TIMEZONE}"
    volumes:
      - "./data:/var/lib/postgresql/data"
    networks:
      - "${GENERAL_DEFAULT_NETWORK}"
    ports:
      - "${CONFIG_PORT}:5432"

# Configurable fields (for TUI form)
config_fields:
  - key: port
    label: "Port"
    type: number
    default: 5432
    required: true
    validation:
      min: 1024
      max: 65535
      unique_per_service: true

  - key: user
    label: "Username"
    type: text
    default: postgres
    required: true
    validation:
      min_length: 1
      max_length: 63
      pattern: "^[a-zA-Z_][a-zA-Z0-9_]*$"

  - key: password
    label: "Password"
    type: password
    default: postgres
    required: true
    validation:
      min_length: 4

  - key: database
    label: "Database Name"
    type: text
    default: postgres
    required: true
    validation:
      min_length: 1
      max_length: 63

# Dependencies (other services needed)
dependencies: []

# Post-setup commands
post_setup:
  - "echo 'PostgreSQL is ready on port ${CONFIG_PORT}'"
```

### Template Variable System

Variables are interpolated using `${VAR_NAME}` syntax.

| Variable | Source | Example |
|---|---|---|
| `CONFIG_<field>` | User config input | `CONFIG_PORT=5432` |
| `GENERAL_<key>` | `config.yaml` general section | `GENERAL_TIMEZONE=Asia/Jakarta` |
| `SERVICE_<key>` | Service metadata | `SERVICE_NAME=postgresql` |
| `VERSION_<key>` | Version metadata | `VERSION_NUMBER=16` |

### Built-in Templates

| Service | Versions | Category | Default Port |
|---|---|---|---|
| PostgreSQL | 15, 16, 17 | database | 5432 |
| MySQL | 8.0, 8.4, 9.0 | database | 3306 |
| MariaDB | 11 | database | 3306 |
| Redis | 7 | cache | 6379 |
| MongoDB | 7, 8 | database | 27017 |
| MinIO | latest | storage | 9000 |
| Elasticsearch | 8 | search | 9200 |
| Memcached | 1.6 | cache | 11211 |

### Custom Templates

```bash
# Create new template
dockkit template add myservice

# Copy from existing template
dockkit template add myservice --from postgresql

# Edit existing template
dockkit template edit myservice

# Delete custom template
dockkit template remove myservice
```

Custom templates are stored in `~/.config/dockkit/templates/`.

### Template Priority

1. Custom templates (user-created) — highest priority
2. Built-in templates (embedded in binary) — lowest priority

If a custom template has the same name as a built-in, the custom one takes precedence.

---

## TUI Screens

### 1. Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  dockkit v1.0.0                               [?] Help [q] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Services                                                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ ● PostgreSQL 16     :5432    healthy    up 2h 15m     │  │
│  │ ● Redis 7           :6379    healthy    up 2h 15m     │  │
│  │ ○ MySQL 8           :3306    stopped    -              │  │
│  │ ○ Elasticsearch 8   :9200    stopped    -              │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  [a] Start All   [A] Stop All   [+] Add   [R] Refresh       │
│                                                             │
│  Docker: running │ 4 services │ 2 active                     │
└─────────────────────────────────────────────────────────────┘
```

**Keybindings:**
| Key | Action |
|---|---|
| `↑/↓` | Navigate services |
| `Enter` | Service detail |
| `+` | Add service |
| `x` | Stop selected |
| `X` | Stop all (with confirm) |
| `s` | Start selected |
| `S` | Start all (with confirm) |
| `l` | View logs |
| `c` | Edit config |
| `r` | Restart selected |
| `R` | Refresh status |
| `?` | Help overlay |
| `q` | Quit |

### 2. Service Picker

```
┌─────────────────────────────────────────────────────────────┐
│  Add New Service                              [Esc] Back    │
├─────────────────────────────────────────────────────────────┤
│  Filter: [________________________________]                 │
│                                                             │
│  Database                                                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ PostgreSQL     v15  v16  v17     [Install]             │  │
│  │ MySQL          v8.0 v8.4 v9.0    [Install]             │  │
│  │ MariaDB        v11               [Install]             │  │
│  │ MongoDB        v7  v8            [Install]             │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Cache                                                      │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Redis          v7                 [Install]             │  │
│  │ Memcached      v1.6               [Install]             │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Storage                                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ MinIO          latest             [Install]             │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Search                                                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Elasticsearch  v8                 [Install]             │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

**Keybindings:**
| Key | Action |
|---|---|
| `↑/↓` | Navigate services |
| `←/→` | Navigate versions |
| `Enter` | Select version -> config wizard |
| `/` | Focus filter |
| `Esc` | Back to dashboard |

### 3. Config Wizard

```
┌─────────────────────────────────────────────────────────────┐
│  Configure PostgreSQL 16                      [Esc] Back    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Service Configuration                                      │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Port:            [5432      ]                          │  │
│  │ Username:        [postgres  ]                          │  │
│  │ Password:        [••••••••  ]                          │  │
│  │ Database:        [postgres  ]                          │  │
│  │ Container Name:  [auto-generated]                      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Advanced                                                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Timezone:        [Asia/Jakarta]                        │  │
│  │ Network:         [dockkit-network]                     │  │
│  │ Restart Policy:  [unless-stopped]                      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Preview                                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Image: postgres:16-alpine                             │  │
│  │ Port:  localhost:5432 -> container:5432               │  │
│  │ Health: pg_isready enabled                            │  │
│  │ Volume: ./data (persistent)                           │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  [Ctrl+S] Install    [Ctrl+T] Use Template    [Esc] Cancel  │
└─────────────────────────────────────────────────────────────┘
```

**Keybindings:**
| Key | Action |
|---|---|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Ctrl+S` | Install service |
| `Ctrl+T` | Choose template |
| `Esc` | Cancel (with confirm if changes made) |

### 4. Service Detail

```
┌─────────────────────────────────────────────────────────────┐
│  PostgreSQL 16                                [Esc] Back    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Status: ● Running (healthy)          Uptime: 2h 15m       │
│                                                             │
│  Details                                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Image:          postgres:16-alpine                    │  │
│  │ Container:      dockkit-postgresql-16                 │  │
│  │ Port:           localhost:5432                         │  │
│  │ Network:        dockkit-network                       │  │
│  │ Volume:         ./data                                │  │
│  │ Health:         healthy (5/5)                         │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Connection                                                 │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Host:     localhost                                   │  │
│  │ Port:     5432                                        │  │
│  │ User:     postgres                                    │  │
│  │ Password: postgres                                    │  │
│  │ Database: postgres                                    │  │
│  │ URI:      postgresql://postgres:postgres@localhost:5432/postgres │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Actions                                                    │
│  [s] Start  [x] Stop  [r] Restart  [l] Logs  [c] Config    │
│  [b] Backup  [d] Remove (with confirm)                      │
└─────────────────────────────────────────────────────────────┘
```

### 5. Logs Viewer

```
┌─────────────────────────────────────────────────────────────┐
│  Logs: PostgreSQL 16                         [Esc] Back     │
├─────────────────────────────────────────────────────────────┤
│  [F] Follow: ON    [T] Timestamps: ON    [L] Level: ALL    │
│  Search: [________________________________]                  │
├─────────────────────────────────────────────────────────────┤
│  10:15:23 [1] LOG:  starting PostgreSQL 16.4               │
│  10:15:23 [1] LOG:  listening on IPv4 0.0.0.0:5432        │
│  10:15:23 [1] LOG:  listening on Unix socket               │
│  10:15:24 [1] LOG:  database system is ready to accept     │
│  10:15:24 [1] LOG:  checkpoint start: immediate            │
│  10:16:00 [1] LOG:  checkpoint complete: wrote 3 buffers   │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  Lines: 156/156 │ Filter: 0 matches                         │
└─────────────────────────────────────────────────────────────┘
```

**Keybindings:**
| Key | Action |
|---|---|
| `↑/↓` | Scroll |
| `Page Up/Down` | Scroll page |
| `F` | Toggle follow |
| `T` | Toggle timestamps |
| `L` | Cycle log level (ALL, ERROR, WARN, INFO) |
| `/` | Focus search |
| `Esc` | Back |

### 6. Template Manager

```
┌─────────────────────────────────────────────────────────────┐
│  Templates                               [Esc] Back        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Built-in                                                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ PostgreSQL     3 versions  [View]                     │  │
│  │ MySQL          3 versions  [View]                     │  │
│  │ MariaDB        1 version   [View]                     │  │
│  │ Redis          1 version   [View]                     │  │
│  │ MongoDB        2 versions  [View]                     │  │
│  │ MinIO          1 version   [View]                     │  │
│  │ Elasticsearch  1 version   [View]                     │  │
│  │ Memcached      1 version   [View]                     │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Custom                                                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ my-postgres    1 version   [View] [Edit] [Delete]    │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  [n] New Template    [i] Import    [e] Export                │
└─────────────────────────────────────────────────────────────┘
```

### 7. Version Fetcher

```
┌─────────────────────────────────────────────────────────────┐
│  Fetching versions for PostgreSQL...              [Esc]     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Docker Hub Tags                                            │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 17.5    alpine, bookworm    [Select]                  │  │
│  │ 16.9    alpine, bookworm    [Select]                  │  │
│  │ 15.13   alpine, bookworm    [Select]                  │  │
│  │ 14.18   alpine, bookworm    [EOL]                     │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  Image Variants (for selected version)                      │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ postgres:17-alpine     ~80MB     [Recommended]        │  │
│  │ postgres:17-bookworm   ~400MB                         │  │
│  │ postgres:17             ~400MB                         │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  [Enter] Install    [R] Refresh    [Esc] Back               │
└─────────────────────────────────────────────────────────────┘
```

### 8. Template Editor

```
┌─────────────────────────────────────────────────────────────┐
│  Edit Template: my-postgres                 [Esc] Back      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  YAML Editor                                                │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ name: my-postgres                                     │  │
│  │ description: Custom PostgreSQL template               │  │
│  │ category: database                                    │  │
│  │ icon: "🐘"                                            │  │
│  │                                                       │  │
│  │ versions:                                             │  │
│  │   "16":                                               │  │
│  │     image: postgres:16-alpine                         │  │
│  │     default_port: 5432                                │  │
│  │     ...                                               │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  [Ctrl+S] Save    [Ctrl+V] Validate    [Esc] Cancel         │
└─────────────────────────────────────────────────────────────┘
```

### 9. Error Screen

```
┌─────────────────────────────────────────────────────────────┐
│  Error                                    [?] Help          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Docker is not running                                      │
│                                                             │
│  dockkit requires Docker Engine to be running.              │
│                                                             │
│  To fix this:                                               │
│  1. Start Docker Desktop, or                                │
│  2. Run: sudo systemctl start docker                       │
│                                                             │
│  [R] Retry    [q] Quit                                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Keybindings

### Global

| Key | Action |
|---|---|
| `?` | Toggle help overlay |
| `q` | Quit / Back |
| `Esc` | Back (with confirm if unsaved changes) |
| `Ctrl+C` | Force quit |
| `Tab` | Next focusable element |
| `Shift+Tab` | Previous focusable element |

### Mouse

| Action | Behavior |
|---|---|
| Left click | Select / Activate |
| Scroll wheel | Scroll lists and logs |
| Hover | Highlight (where supported) |

---

## Data Model

### ServiceTemplate

```go
type ServiceTemplate struct {
    Name         string            `yaml:"name"`
    Description  string            `yaml:"description"`
    Category     string            `yaml:"category"`     // database, cache, storage, search
    Icon         string            `yaml:"icon"`
    Source       TemplateSource    `yaml:"-"`             // built-in or custom
    Versions     []VersionTemplate `yaml:"versions"`      // ordered, not map
    ConfigFields []ConfigField     `yaml:"config_fields"`
    Dependencies []string          `yaml:"dependencies"`
    PostSetup    []string          `yaml:"post_setup"`
}

type TemplateSource int

const (
    TemplateBuiltIn TemplateSource = iota
    TemplateCustom
)

type VersionTemplate struct {
    Key         string            `yaml:"key"`           // "16", "7", "latest"
    Image       string            `yaml:"image"`
    DefaultPort int               `yaml:"default_port"`
    Healthcheck *HealthcheckConfig `yaml:"healthcheck,omitempty"`
    EnvVars     map[string]string `yaml:"env_vars"`
    Volumes     []string          `yaml:"volumes"`
    Networks    []string          `yaml:"networks"`
    Ports       []string          `yaml:"ports"`
    Command     *string           `yaml:"command,omitempty"`
    ShmSize     *string           `yaml:"shm_size,omitempty"`
}

type HealthcheckConfig struct {
    Test     []string `yaml:"test"`
    Interval string   `yaml:"interval"`
    Timeout  string   `yaml:"timeout"`
    Retries  int      `yaml:"retries"`
}
```

### ConfigField

```go
type ConfigField struct {
    Key            string           `yaml:"key"`
    Label          string           `yaml:"label"`
    Type           string           `yaml:"type"`        // text, number, password, select
    Default        string           `yaml:"default"`
    Required       bool             `yaml:"required"`
    Validation     *ValidationRules `yaml:"validation,omitempty"`
    Options        []string         `yaml:"options,omitempty"`    // for select type
    Placeholder    string           `yaml:"placeholder,omitempty"`
}

type ValidationRules struct {
    Min             *int    `yaml:"min,omitempty"`
    Max             *int    `yaml:"max,omitempty"`
    MinLength       *int    `yaml:"min_length,omitempty"`
    MaxLength       *int    `yaml:"max_length,omitempty"`
    Pattern         string  `yaml:"pattern,omitempty"`
    UniquePerService bool   `yaml:"unique_per_service,omitempty"`
}
```

### ServiceState

```go
type ServiceState struct {
    Name          string    `json:"name"`
    Version       string    `json:"version"`
    Status        string    `json:"status"`       // running, stopped, healthy, unhealthy
    ContainerID   string    `json:"container_id"`
    ContainerName string    `json:"container_name"`
    Image         string    `json:"image"`
    Port          int       `json:"port"`
    Uptime        string    `json:"uptime"`
    HealthStatus  string    `json:"health_status"` // healthy, unhealthy, starting, none
    CreatedAt     time.Time `json:"created_at"`
    StartedAt     *time.Time `json:"started_at,omitempty"`
}
```

### ServiceConfig

```go
type ServiceConfig struct {
    Version       string            `yaml:"version"`
    Enabled       bool              `yaml:"enabled"`
    Port          int               `yaml:"port"`
    ContainerName string            `yaml:"container_name"`
    Image         string            `yaml:"image"`
    User          string            `yaml:"user"`
    Password      string            `yaml:"password"`
    Database      string            `yaml:"database"`
    Extra         map[string]string `yaml:"extra,omitempty"`
}
```

### ComposeConfig (Generated)

```go
type ComposeConfig struct {
    Version  string                    `yaml:"version"`
    Services map[string]ComposeService `yaml:"services"`
    Networks map[string]ComposeNetwork `yaml:"networks,omitempty"`
    Volumes  map[string]ComposeVolume  `yaml:"volumes,omitempty"`
}

type ComposeService struct {
    Image         string            `yaml:"image"`
    ContainerName string            `yaml:"container_name"`
    Restart       string            `yaml:"restart"`
    Environment   map[string]string `yaml:"environment"`
    Ports         []string          `yaml:"ports"`
    Volumes       []string          `yaml:"volumes"`
    Healthcheck   *ComposeHealthcheck `yaml:"healthcheck,omitempty"`
    Networks      []string          `yaml:"networks"`
    Command       *string           `yaml:"command,omitempty"`
    ShmSize       *string           `yaml:"shm_size,omitempty"`
}

type ComposeHealthcheck struct {
    Test     []string `yaml:"test"`
    Interval string   `yaml:"interval"`
    Timeout  string   `yaml:"timeout"`
    Retries  int      `yaml:"retries"`
}

type ComposeNetwork struct {
    Name   string `yaml:"name"`
    Driver string `yaml:"driver"`
}

type ComposeVolume struct {
    Name   string `yaml:"name"`
    Driver string `yaml:"driver,omitempty"`
}
```

---

## Architecture

```
dockkit/
├── main.go                          # Entry point (minimal)
├── go.mod
├── go.sum
│
├── cmd/                             # Cobra commands
│   ├── root.go                      # Root command + persistent flags
│   ├── tui.go                       # Launch TUI (default)
│   ├── init.go                      # dockkit init
│   ├── list.go                      # dockkit list
│   ├── up.go                        # dockkit up
│   ├── down.go                      # dockkit down
│   ├── restart.go                   # dockkit restart
│   ├── logs.go                      # dockkit logs
│   ├── setup.go                     # dockkit setup
│   ├── templates.go                 # dockkit templates (parent)
│   ├── template_list.go             # dockkit templates list
│   ├── template_show.go             # dockkit template show
│   ├── template_add.go              # dockkit template add
│   ├── template_remove.go           # dockkit template remove
│   └── version.go                   # dockkit version
│
├── internal/
│   ├── config/                      # Config management
│   │   ├── config.go                # Load/save/validate config
│   │   ├── paths.go                 # XDG paths
│   │   ├── defaults.go              # Default config values
│   │   └── migrate.go               # Config migration
│   │
│   ├── docker/                      # Docker operations
│   │   ├── client.go                # Docker client wrapper
│   │   ├── containers.go            # Batch list/start/stop/restart
│   │   ├── images.go                # Pull/list/remove images
│   │   ├── logs.go                  # Stream logs
│   │   ├── compose.go               # Docker Compose exec wrapper
│   │   ├── health.go                # Health check operations
│   │   └── network.go               # Network create/remove
│   │
│   ├── conflict/                    # Conflict detection & resolution
│   │   ├── detector.go              # Scan config + host for conflicts
│   │   ├── resolver.go              # Auto-fix + suggestion engine
│   │   ├── portcheck.go             # TCP port probe
│   │   └── types.go                 # Conflict, ConflictType, Resolution
│   │
│   ├── perf/                        # Performance monitoring
│   │   ├── poller.go                # Adaptive polling controller
│   │   ├── stats.go                 # Runtime stats collection
│   │   └── cache.go                 # Generic in-memory cache with TTL
│   │
│   ├── registry/                    # Docker Hub API
│   │   ├── hub.go                   # Fetch tags from Docker Hub
│   │   ├── cache.go                 # Cache tag responses
│   │   └── types.go                 # Registry response types
│   │
│   ├── templates/                   # Template system
│   │   ├── loader.go                # Load built-in + custom templates
│   │   ├── renderer.go              # Render template to docker-compose
│   │   ├── validator.go             # Validate template
│   │   ├── interpolator.go          # Variable interpolation
│   │   └── builtins/                # Embedded templates (via go:embed)
│   │       ├── postgresql.yaml
│   │       ├── mysql.yaml
│   │       ├── mariadb.yaml
│   │       ├── redis.yaml
│   │       ├── mongodb.yaml
│   │       ├── minio.yaml
│   │       ├── elasticsearch.yaml
│   │       └── memcached.yaml
│   │
│   ├── errors/                      # Error handling
│   │   ├── errors.go                # Custom error types
│   │   ├── recovery.go              # Panic recovery
│   │   └── messages.go              # User-facing error messages
│   │
│   └── tui/                         # Bubble Tea TUI
│       ├── app.go                   # Main app model + router
│       ├── styles.go                # Lipgloss styles
│       ├── keys.go                  # Global keybindings
│       ├── theme.go                 # Color theme
│       │
│       ├── screens/                 # Screen models
│       │   ├── dashboard.go
│       │   ├── service_picker.go
│       │   ├── version_fetcher.go
│       │   ├── config_wizard.go
│       │   ├── service_detail.go
│       │   ├── logs_viewer.go
│       │   ├── template_manager.go
│       │   ├── template_editor.go
│       │   └── error_screen.go
│       │
│       ├── components/              # Reusable UI components
│       │   ├── service_card.go
│       │   ├── yaml_editor.go
│       │   ├── search_bar.go
│       │   ├── status_badge.go
│       │   ├── confirm_dialog.go
│       │   ├── toast.go
│       │   ├── help_overlay.go
│       │   ├── loading_spinner.go
│       │   ├── progress_bar.go
│       │   └── version_selector.go
│       │
│       └── messages/                # Custom message types
│           ├── navigation.go
│           ├── docker.go
│           └── config.go
│
├── templates/                       # Template YAML files (embedded)
│   ├── postgresql.yaml
│   ├── mysql.yaml
│   ├── mariadb.yaml
│   ├── redis.yaml
│   ├── mongodb.yaml
│   ├── minio.yaml
│   ├── elasticsearch.yaml
│   └── memcached.yaml
│
├── internal/tui/testdata/           # Test fixtures
│
├── Makefile                         # Dev commands
├── .goreleaser.yaml                 # Release config
├── .gitignore
└── README.md
```

---

## Conflict Detection & Resolution

### Overview

dockkit implements a 3-layer conflict detection system to prevent misconfigurations when multiple services run simultaneously.

### Conflict Types

| Type | Description | Example |
|---|---|---|
| `ConflictPort` | Two services use the same host port | PG16 and PG17 both on 5432 |
| `ConflictContainerName` | Two services use the same container name | Manual container named `dockkit-postgresql-16` |
| `ConflictNetwork` | Network name collision | External network named `dockkit-network` |
| `ConflictVolume` | Volume name collision | Two services sharing same volume path |

### Detection Layers

#### Layer 1: Config Validation (at save time)

```
Config Wizard → User fills port → Ctrl+S
    ↓
ConflictDetector.ScanAll()
    ↓
┌─────────────────────────────────────────────┐
│ Scan enabled services for port collisions   │
│   PG16 port=5432 ✓                          │
│   PG17 port=5432 ✗ CONFLICT with PG16       │
│                                             │
│ Scan host for occupied ports                │
│   port 5432 → process: postgres (host)      │
│   port 3306 → free ✓                        │
│                                             │
│ Scan container names                        │
│   dockkit-postgresql-16 → exists in Docker  │
└─────────────────────────────────────────────┘
    ↓
NO conflicts → Save config
YES conflicts → Block save, show resolution dialog
```

#### Layer 2: Pre-flight Check (before container start)

```
dockkit up postgresql-17
    ↓
Pre-flight checks:
  1. Port 5432 available? → NO (used by postgresql-16)
  2. Container name available? → YES
  3. Network accessible? → YES
  4. Volume writable? → YES
    ↓
CONFLICT DETECTED → Show resolution dialog
    ↓
User resolves → Retry start
```

#### Layer 3: Runtime Detection (after container failure)

```
docker compose up -d
    ↓
Container exits (exit code != 0)
    ↓
Parse logs for error patterns:
  - "bind: address already in use" → port conflict
  - "name is already in use" → container name conflict
  - "network already exists" → network conflict
    ↓
Show toast with specific error + suggestion
```

### Data Model

```go
type ConflictDetector struct {
    config       *Config
    dockerClient DockerClient
}

type Conflict struct {
    Type       ConflictType
    Severity   ConflictSeverity
    ServiceA   string  // existing service using resource
    ServiceB   string  // new service wanting resource
    Resource   string  // port number, name, etc.
    Message    string  // human-readable description
    Suggested  string  // suggested alternative
}

type ConflictType int

const (
    ConflictPort ConflictType = iota
    ConflictContainerName
    ConflictNetwork
    ConflictVolume
)

type ConflictSeverity int

const (
    SeverityError   ConflictSeverity = iota  // blocks operation
    SeverityWarning                           // warns but allows
)

type ConflictResolution int

const (
    ResolutionAutoFix  ConflictResolution = iota  // use suggested alternative
    ResolutionChoose                                // user picks alternative
    ResolutionCancel                                // abort operation
)
```

### Conflict Detector Implementation

```go
func (d *ConflictDetector) Detect() []Conflict {
    var conflicts []Conflict

    // 1. Port conflicts between configured services
    portMap := map[int]string{} // port -> "service version"
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
                    Suggested: d.suggestPort(cfg.Port),
                })
            } else {
                portMap[cfg.Port] = key
            }
        }
    }

    // 2. Port conflicts with host processes
    for name, svc := range d.config.Services {
        for ver, cfg := range svc.Versions {
            if !cfg.Enabled {
                continue
            }
            if d.isPortOccupiedOnHost(cfg.Port) {
                if !d.isOwnContainer(cfg.ContainerName) {
                    conflicts = append(conflicts, Conflict{
                        Type:     ConflictPort,
                        Severity: SeverityError,
                        ServiceB: fmt.Sprintf("%s %s", name, ver),
                        Resource: fmt.Sprintf("%d", cfg.Port),
                        Message:  fmt.Sprintf("Port %d is used by another process on host", cfg.Port),
                        Suggested: d.suggestPort(cfg.Port),
                    })
                }
            }
        }
    }

    // 3. Container name conflicts
    containers, _ := d.dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
    existingNames := map[string]bool{}
    for _, c := range containers {
        for _, n := range c.Names {
            existingNames[strings.TrimPrefix(n, "/")] = true
        }
    }
    for name, svc := range d.config.Services {
        for ver, cfg := range svc.Versions {
            if !cfg.Enabled {
                continue
            }
            if existingNames[cfg.ContainerName] {
                // Check if it's managed by dockkit or external
                if !d.isManagedByDockkit(cfg.ContainerName) {
                    conflicts = append(conflicts, Conflict{
                        Type:     ConflictContainerName,
                        Severity: SeverityError,
                        ServiceB: fmt.Sprintf("%s %s", name, ver),
                        Resource: cfg.ContainerName,
                        Message:  fmt.Sprintf("Container '%s' already exists (not managed by dockkit)", cfg.ContainerName),
                    })
                }
            }
        }
    }

    return conflicts
}

// suggestPort finds the next available port starting from port+1
func (d *ConflictDetector) suggestPort(port int) string {
    for offset := 1; offset <= 100; offset++ {
        candidate := port + offset
        if !d.isPortOccupiedOnHost(candidate) && !d.isPortUsedByService(candidate) {
            return fmt.Sprintf("%d", candidate)
        }
    }
    return ""
}

// isPortOccupiedOnHost checks if a port is in use on the host
func (d *ConflictDetector) isPortOccupiedOnHost(port int) bool {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

// isPortUsedByService checks if port is configured for any dockkit service
func (d *ConflictDetector) isPortUsedByService(port int) bool {
    for _, svc := range d.config.Services {
        for _, cfg := range svc.Versions {
            if cfg.Enabled && cfg.Port == port {
                return true
            }
        }
    }
    return false
}
```

### Resolution Dialog (TUI)

```
┌─────────────────────────────────────────────────────────────┐
│  Port Conflict Detected                       [Esc] Cancel  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  PostgreSQL 17 wants port 5432                              │
│  but PostgreSQL 16 already uses port 5432                   │
│                                                             │
│  ┌─ Options ──────────────────────────────────────────────┐ │
│  │                                                         │ │
│  │  > Use port 5434 (recommended)                          │ │
│  │    Choose different port                                │ │
│  │    Cancel install                                       │ │
│  │                                                         │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  [Enter] Select    [Esc] Cancel                              │
└─────────────────────────────────────────────────────────────┘
```

### Auto-port Allocation Rules

When user clicks "Use recommended" or when auto-fixing:

1. Start from `template.default_port + 1`
2. Skip ports already configured for any enabled service
3. Skip ports occupied on host (quick TCP probe)
4. Skip privileged ports (< 1024)
5. Maximum search range: +100 from original port
6. If no port found, require manual input

### Precedence

| Check | Priority | Blocks? |
|---|---|---|
| Port conflict (service vs service) | 1 | Yes |
| Port conflict (service vs host) | 2 | Yes |
| Container name conflict | 3 | Yes |
| Network exists | 4 | No (skip creation) |
| Volume exists | 5 | No (reuse) |

---

## Performance Optimization Strategy

### Resource Budget

| Resource | Budget | Measurement |
|---|---|---|
| CPU idle | < 5% | `top -p <pid>` |
| CPU active | < 20% | During status refresh |
| Memory RSS | < 20MB | `ps -o rss= -p <pid>` (10 services) |
| Goroutines | < 10 | `runtime.NumGoroutine()` |
| File handles | < 20 | `lsof -p <pid> | wc -l` |
| Network | < 1 req/sec avg | Docker API + Hub combined |
| Disk I/O | Minimal after startup | Embedded templates, cache |

### Docker API Optimization

#### Batch Operations (Critical)

```
BAD:  N API calls for N services
┌─────────────┐
│ for svc in services: │
│   client.ContainerInspect(name)  │  ← N calls, ~50ms each
└─────────────┘

GOOD: 1 API call for all services
┌─────────────┐
│ client.ContainerList(All: true)  │  ← 1 call, ~50ms total
│ for _, c := range containers:    │
│   parse from list response       │
└─────────────┘
```

#### Adaptive Polling

```
Polling frequency adapts to activity:

State unchanged 3x → Slow polling (30s)
State changed    → Fast polling (5s)
TUI unfocused    → Minimal polling (60s)
CLI mode         → One-shot (no polling)

┌─────────────────────────────────────────────┐
│                Poll Controller               │
│                                             │
│  State: idle                                │
│  Interval: 30s                              │
│  Last change: 2m ago                        │
│                                             │
│  ── container starts ──→                    │
│  State: active                              │
│  Interval: 5s                               │
│  Last change: now                           │
│                                             │
│  ── no changes for 3 polls ──→              │
│  State: idle                                │
│  Interval: 30s                              │
└─────────────────────────────────────────────┘
```

#### Background Workers

```
┌─────────────────────────────────────────────────────┐
│  Goroutine Map:                                     │
│                                                     │
│  Main (TUI event loop)                              │
│    ├── StatusWorker (every 5-30s)                   │
│    │     └── ContainerList → update states          │
│    │                                                │
│    ├── HubWorker (on-demand, cached 24h)            │
│    │     └── Fetch tags → save to cache             │
│    │                                                │
│    ├── LogWorker (streaming, per-service)           │
│    │     └── ContainerLogs → pipe to viewport       │
│    │                                                │
│    └── ComposeWorker (on user action)               │
│          └── exec docker compose up/down            │
│                                                     │
│  Rules:                                             │
│  - Never block main goroutine                       │
│  - Buffered channels (size 10)                      │
│  - Timeout on all API calls (5s)                    │
│  - Cancel context on screen change                  │
│  - One worker per concern                           │
└─────────────────────────────────────────────────────┘
```

### TUI Rendering Optimization

#### Minimal Re-renders

```go
// BAD: Re-render on every tick
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.tick++
    return m, tickCmd()  // re-renders every time
}

// GOOD: Re-render only on state change
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case statusUpdateMsg:
        if m.statesChanged(msg.states) {
            m.states = msg.states
            return m, nil  // triggers View() re-render
        }
        return m, nil  // no change, no re-render
    }
}
```

#### Virtual Scrolling (Logs)

```
BAD:  Render all 10,000 log lines
┌─────────────────────────────┐
│ line 1                      │
│ line 2                      │
│ ...                         │  ← 10,000 lines rendered
│ line 10000                  │
└─────────────────────────────┘

GOOD: Render only visible viewport
┌─────────────────────────────┐
│ viewport.go handles scroll  │
│ Only 30-50 lines rendered   │  ← 50 lines rendered
│ based on terminal height    │
└─────────────────────────────┘
```

#### Debounced Input

```
User types in search filter:

Keystroke 1: "p"     → debounce timer starts (100ms)
Keystroke 2: "po"    → timer resets
Keystroke 3: "pos"   → timer resets
Keystroke 4: "post"  → timer resets
Timer expires         → filter applied, list re-rendered

Result: 4 keystrokes → 1 filter operation
```

#### Static Content Caching

```go
// Cache rendered content that doesn't change
type RenderCache struct {
    helpOverlay string        // rendered once
    statusBar   string        // rendered once per tick
    templateList string       // invalidated on template change
}

// Only invalidate when underlying data changes
func (c *RenderCache) Invalidate(templateChanged bool) {
    if templateChanged {
        c.templateList = ""  // re-render next time
    }
    // helpOverlay never invalidates
}
```

### Startup Sequence (Optimized)

```
dockkit launch
    ↓ [0ms]
1. Parse CLI args              (< 1ms)
2. Load config.yaml            (< 1ms, cached in memory)
3. Load embedded templates     (< 1ms, go:embed)
4. Initialize TUI model        (< 10ms)
5. Launch Bubble Tea program   (< 50ms)
    ↓ [~60ms total]
6. TUI visible to user
    ↓ [background, non-blocking]
7. Ping Docker daemon          (< 100ms)
8. List all containers         (< 50ms)
9. Update service states       (< 10ms)
    ↓ [~200ms total perceived]
Ready.
```

### Memory Management

```go
// Limit log buffer in memory
const MaxLogLines = 1000

type LogBuffer struct {
    lines []string
    mu    sync.RWMutex
}

func (b *LogBuffer) Append(line string) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.lines = append(b.lines, line)
    if len(b.lines) > MaxLogLines {
        b.lines = b.lines[len(b.lines)-MaxLogLines:]  // trim oldest
    }
}

// Reuse ServiceState objects
var statePool = sync.Pool{
    New: func() interface{} {
        return &ServiceState{}
    },
}

func getServiceState() *ServiceState {
    return statePool.Get().(*ServiceState)
}

func putServiceState(s *ServiceState) {
    *s = ServiceState{}  // reset
    statePool.Put(s)
}
```

### Performance Targets

| Operation | Target | Strategy |
|---|---|---|
| TUI cold start | < 500ms | Lazy Docker init, embedded templates |
| TUI warm start | < 200ms | Config cache, skip Docker re-init |
| Docker status | < 500ms | Batch API, async polling |
| Docker Hub fetch | < 2s | 24h cache, fallback to embedded |
| Template render | < 50ms | In-memory, no disk I/O |
| Log streaming | 30fps | Event-driven, virtual scroll |
| Port conflict check | < 100ms | TCP probe with 100ms timeout |
| Config save | < 50ms | Atomic write, no re-read |

### Performance Monitoring (Debug Mode)

```bash
dockkit --verbose list
# Output:
# dockkit v1.0.0 (debug mode)
# Goroutines: 6
# Memory: 4.2 MB
# Docker API calls: 3 (cached: 2)
# Hub API calls: 0 (cached: 8)
# Last status poll: 2.3s ago
# Services: 4 (2 running, 2 stopped)
```

```go
type Stats struct {
    Uptime         time.Duration
    ServicesTotal  int
    ServicesActive int
    DockerAPICalls int64
    HubAPICalls    int64
    CacheHits      int64
    CacheMisses    int64
    Goroutines     int
    MemoryAlloc    uint64
    LastPollTime   time.Duration
}
```

---

## Error Handling Matrix

| Error | Detection | User Message | Action |
|---|---|---|---|
| Docker not running | `client.Ping()` fails | "Docker is not running. Start Docker Desktop or run `sudo systemctl start docker`" | Show error screen with retry |
| Container crashed | `State.Status == "exited"` | "Container stopped with exit code {N}" | Show status + suggest restart |
| Port conflict (service vs service) | ConflictDetector scan | "Port {port} already used by {service_a}. {service_b} cannot use the same port." | Show resolution dialog with suggested alternative |
| Port conflict (service vs host) | TCP probe + ConflictDetector | "Port {port} is already in use by another process on your system." | Show resolution dialog, suggest alternative port |
| Port conflict (runtime) | Container logs parse | "Container failed to start: port {port} already in use" | Show error toast + suggest fix |
| Container name conflict | ContainerList scan | "Container '{name}' already exists (not managed by dockkit)" | Offer to remove existing or rename |
| Container name taken by dockkit | ContainerList + label check | "Container '{name}' is already managed by dockkit" | Show which service owns it |
| Network exists | Create error | "Network '{name}' already exists" | Skip creation (idempotent) |
| Permission denied | Docker API 403 | "Permission denied. Try running with sudo or add user to docker group" | Show fix instructions |
| Disk space low | Docker API error | "Insufficient disk space" | Show cleanup suggestions |
| Config corrupted | YAML parse error | "Config file is corrupted. Backup created at {path}" | Auto-backup + reset |
| Partial failure | Write error | "Failed to save {file}. Changes rolled back." | Rollback all changes |
| Hub API rate limited | HTTP 429 | "Docker Hub rate limited. Using cached data." | Fallback to cache |
| Hub API offline | Network error | "Docker Hub unreachable. Using cached data." | Fallback to cache |
| Invalid template | Validation error | "Template error: {detail}" | Show validation message |
| No available port | suggestPort returns "" | "No available port found in range {start}-{end}. Please choose manually." | Show manual port input |

---

## Security Guidelines

### Password Handling

- Passwords stored in `config.yaml` and `.env` in plaintext (acceptable for dev infrastructure)
- Passwords NOT logged or printed in TUI (masked with `••••••••`)
- Passwords NOT included in error messages
- `.env` files created with `0600` permissions

### File Permissions

| File | Permissions | Owner |
|---|---|---|
| `~/.config/dockkit/config.yaml` | `0600` | User |
| `~/.config/dockkit/templates/*.yaml` | `0644` | User |
| `~/.config/dockkit/services/**/.env` | `0600` | User |
| `~/.config/dockkit/services/**/data/` | `0700` | User |

### Input Sanitization

- Template variables interpolated using strict pattern matching
- No shell execution of user input (templates only generate YAML)
- Container names validated against Docker naming rules
- Port numbers validated as integers 1024-65535

---

## Caching Strategy

### Docker Hub Cache

```
~/.config/dockkit/cache/tags/
├── postgresql.json    # Tags for postgres image
├── mysql.json         # Tags for mysql image
├── mariadb.json       # Tags for mariadb image
├── redis.json         # Tags for redis image
├── mongodb.json       # Tags for mongo image
├── minio.json         # Tags for minio image
├── elasticsearch.json # Tags for elasticsearch image
└── memcached.json     # Tags for memcached image
```

### Cache Format

```json
{
  "image": "postgres",
  "fetched_at": "2026-07-30T10:00:00Z",
  "ttl": "24h",
  "tags": [
    {"name": "17.5-alpine", "size": 80000000, "last_updated": "2026-07-15"},
    {"name": "16.9-alpine", "size": 79000000, "last_updated": "2026-07-15"}
  ]
}
```

### Cache Rules

- TTL: 24 hours
- Refresh on: manual refresh (R key), new service setup
- Fallback: use cached data if API unavailable
- Invalidation: manual (dockkit cache clear)

---

## Config Migration Strategy

### Version Field

```yaml
version: "1"
```

### Migration Rules

1. **Never break existing configs** — new fields get defaults
2. **Migration functions** per version bump
3. **Auto-backup** before migration
4. **Log migration** actions

### Migration Example

```go
func MigrateConfig(cfg *Config, fromVersion string) error {
    switch fromVersion {
    case "1":
        // v1 -> v2: add default_network if missing
        if cfg.General.DefaultNetwork == "" {
            cfg.General.DefaultNetwork = "dockkit-network"
        }
        return MigrateConfig(cfg, "2")
    case "2":
        // current version
        return nil
    default:
        return fmt.Errorf("unknown config version: %s", fromVersion)
    }
}
```

---

## Development Setup

### Prerequisites

```bash
# Go 1.25+
go version  # should be 1.25.0+

# Docker
docker version
docker compose version

# Development tools
go install github.com/air-verse/air@latest    # Hot reload
go install github.com/golangci/golangci-lint@latest  # Linter
```

### Setup

```bash
# Clone repo
git clone https://github.com/user/dockkit.git
cd dockkit

# Install dependencies
go mod download

# Run in dev mode
go run .

# Or with hot reload
air

# Run tests
make test

# Build
make build
```

### Makefile

```makefile
.PHONY: build test lint clean dev

dev:
	air

build:
	go build -o bin/dockkit .

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

install:
	go install .

version:
	@echo "dockkit $(shell git describe --tags --always)"
```

---

## Testing Strategy

### Unit Tests

| Package | Test Focus |
|---|---|
| `internal/config` | Config load/save/validate/migrate |
| `internal/templates` | Template loading, rendering, validation |
| `internal/docker` | Container operations (with mock client) |
| `internal/conflict` | Port conflict detection, auto-suggest, resolution |
| `internal/perf` | Adaptive polling, cache TTL, stats collection |
| `internal/registry` | Hub API parsing, cache logic |
| `internal/errors` | Error formatting, user messages |

### Integration Tests

| Test | What |
|---|---|
| Full setup flow | Template -> Config -> Compose -> Validate |
| Docker operations | Real Docker client (optional, CI only) |
| Config migration | v1 -> v2 -> v3 upgrade path |

### Mock Docker Client

```go
// Mock for testing without Docker daemon
type MockClient struct {
    Containers []container.Summary
    Images     []image.Summary
    PullError  error
}
```

### Test Commands

```bash
# Unit tests
go test ./...

# Integration tests (requires Docker)
go test -tags=integration ./...

# Coverage
go test -cover ./... -coverprofile=coverage.out

# Benchmark
go test -bench=. ./...
```

---

## Risk Mitigation

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Docker Hub rate limiting | High | Medium | 24h cache, fallback to embedded versions |
| Config corruption | Low | High | Atomic writes, auto-backup before write |
| Port conflicts (multi-service) | High | High | 3-layer detection, auto-suggest alternatives, resolution dialog |
| Port conflicts (host process) | Medium | Medium | TCP probe on config save, suggest alternative |
| Container name collisions | Medium | Low | `dockkit-` prefix, pre-flight scan, label-based ownership |
| Partial failure (compose + .env) | Low | High | Rollback mechanism, atomic operations |
| Terminal compatibility | Medium | Medium | Test on iTerm2, Alacritty, Windows Terminal, kitty |
| Go version requirement (1.25+) | Low | High | Clear error message, version check at startup |
| Docker Compose v1 (legacy) | Low | Low | Check `docker compose version`, error if v1 |
| Custom template corruption | Low | Low | Validation on load, graceful fallback to defaults |
| Network name conflicts | Medium | Low | Use unique names, create-if-not-exists pattern |
| Memory growth (log streaming) | Medium | Low | Max 1000 lines buffer, trim oldest |
| Polling overhead | Low | Low | Adaptive polling, batch API, background workers |
| Goroutine leak | Low | Medium | Context cancellation, timeout on all operations |

---

## MVP Scope (Phase 1)

### Must Have

- [ ] CLI entrypoint with Cobra
- [ ] TUI home dashboard
- [ ] Service list with status (running/stopped/healthy)
- [ ] Service picker (browse templates by category)
- [ ] Version fetcher from Docker Hub (with cache)
- [ ] Config wizard (form-based setup with validation)
- [ ] Docker compose generation from template
- [ ] Start/Stop/Restart service
- [ ] Basic template system (8 built-in templates)
- [ ] Config file persistence (~/.config/dockkit/)
- [ ] Logs viewer (with follow, timestamps, filter)
- [ ] Service detail view
- [ ] Keyboard navigation (full)
- [ ] Mouse support (click + scroll)
- [ ] Health check integration
- [ ] Error handling for all failure modes
- [ ] Confirmation dialogs for destructive actions
- [ ] Loading states for async operations
- [ ] Toast notifications for actions
- [ ] Port conflict detection (3-layer)
- [ ] Conflict resolution dialog with auto-suggest
- [ ] Adaptive Docker status polling
- [ ] Batch Docker API calls
- [ ] Docker Hub cache with 24h TTL

### Nice to Have (Phase 2)

- [ ] Custom template editor (in TUI)
- [ ] Import/Export templates
- [ ] Backup/Restore
- [ ] Multi-service orchestration (start all with profiles)
- [ ] Shell completions (bash, zsh, fish)
- [ ] Docker context switching
- [ ] Service dependency graph

### Future (Phase 3)

- [ ] Plugin system for custom service types
- [ ] Remote Docker host support
- [ ] Team sharing (sync configs via git)
- [ ] Service health monitoring dashboard
- [ ] Auto-cleanup unused images/containers
- [ ] Integration with mise/asdf for runtime management

---

## Implementation Roadmap

### Phase 0: Foundation (1-2 days)

1. Initialize Go module with correct versions
2. Create project skeleton
3. Implement config loading + XDG paths
4. Create all 8 template YAML files
5. Implement template loader + renderer
6. Implement error types + messages

### Phase 1: Docker Layer (2-3 days)

1. Docker client wrapper (moby/moby/client)
2. Container operations (batch list/start/stop/restart)
3. Image pull + Docker Hub API + caching
4. Docker Compose exec wrapper
5. Health check integration
6. Network management
7. Conflict detection engine (3-layer)
8. Port probe utility
9. Adaptive polling controller

### Phase 2: TUI Core (3-4 days)

1. Main app model + navigation system
2. Dashboard screen
3. Service picker screen
4. Config wizard screen (with huh forms)
5. Service detail screen
6. Error screen
7. Toast system + loading states

### Phase 3: TUI Features (2-3 days)

1. Version fetcher screen
2. Logs viewer screen
3. Template manager screen
4. Template editor screen
5. Help overlay
6. Confirmation dialogs

### Phase 4: CLI Commands (1 day)

1. Non-TUI commands (list, up, down, restart, logs)
2. Template management commands
3. Version + init commands

### Phase 5: Polish (1-2 days)

1. Error handling all screens
2. Mouse support all screens
3. Performance optimization
4. Terminal compatibility testing
5. Unit tests
6. Documentation

**Total estimated: 10-15 days**

---

## UI Design Principles

1. **Keyboard-first**: Every action available via keyboard
2. **Mouse-friendly**: Click and scroll work everywhere
3. **Consistent navigation**: `Esc` always goes back, `?` always shows help
4. **Visual feedback**: Toasts for actions, spinners for loading, progress for long ops
5. **Confirmation**: Destructive actions always ask first
6. **Responsive**: Adapts to terminal size, graceful degradation
7. **Dark theme default**: Modern, easy on the eyes
8. **Color coding**: Status colors (green=running, red=stopped, yellow=unhealthy)
9. **No dead ends**: Every screen has a way back
10. **Predictable**: Same action, same result, always

---

## License

TBD
