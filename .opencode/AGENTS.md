# AGENTS.md — dockkit AI Agent Config

## Project Context

dockkit is a Go CLI TUI tool for managing Docker development infrastructure services. Built with:
- **Go 1.25+** (currently 1.26.5)
- **Bubble Tea v2** (TUI framework)
- **Cobra** (CLI framework)
- **moby/moby/client** (Docker SDK)
- **Lip Gloss v2** (TUI styling)
- **Huh v2** (TUI forms)

Module: `github.com/MAHMETT/dockkit`

## Code Style

### Go Conventions
- Follow standard Go conventions (gofmt, go vet)
- Use `gofmt -s -w .` before committing
- Run `go vet ./...` to check for issues
- One concept per file, max 300 lines per file
- Always wrap errors with `fmt.Errorf("context: %w", err)`
- Use table-driven tests

### Bubble Tea Patterns
- Each screen is a separate Model in `internal/tui/screens/`
- Components are reusable in `internal/tui/components/`
- Messages go through `internal/tui/messages/`
- Never block main goroutine — use `tea.Cmd` for async
- Return `tea.Batch()` for multiple commands

### Docker SDK Patterns
- Always use `client.WithAPIVersionNegotiation()`
- Batch operations: `ContainerList` instead of N `ContainerInspect`
- Use `context.WithTimeout` for all API calls
- Handle `io.EOF` in log streaming

## File Structure

```
Layer 0: Foundation     → go.mod, main.go, cmd/root.go
Layer 1: Config         → internal/config/, internal/templates/, internal/errors/
Layer 2: Docker         → internal/docker/, internal/registry/
Layer 3: Conflict       → internal/conflict/
Layer 4: TUI Framework  → internal/tui/app.go, styles, keys, messages, components
Layer 5: TUI Screens    → internal/tui/screens/
Layer 6: CLI Commands   → cmd/
```

**Rule:** Layers only import downward. Never import upward.

## Build & Test

```bash
# Build
go build -o bin/dockkit .

# Test
go test ./...

# Test with coverage
go test -cover ./... -coverprofile=coverage.out

# Lint
go vet ./...

# Format
gofmt -s -w .
```

## Commit Messages

Format: `<type>(<scope>): <description>`

Types: feat, fix, refactor, test, docs, chore, perf, style

Examples:
```
feat(docker): add batch container listing
fix(config): handle missing timezone field
test(conflict): add port conflict detection tests
```

## Key Decisions

1. **Docker Compose via exec** — Not compose-spec/compose-go. Shell out to `docker compose` for simplicity.
2. **Templates embedded via go:embed** — No runtime file loading for built-in templates.
3. **Config in ~/.config/dockkit/** — XDG-compliant, separate from project data.
4. **Port conflict = 3-layer detection** — Config save, pre-flight, runtime.
5. **Adaptive polling** — 5s active, 30s idle, 60s paused.

## When Editing Files

1. Read the file first to understand context
2. Check neighboring files for patterns
3. Follow existing code style
4. Run `go vet ./...` after changes
5. Run `go build .` to verify compilation
6. Add tests for new functions

## Testing

- Unit tests: `*_test.go` in same package
- Test data: `testdata/` directories
- Table-driven tests preferred
- Test both success and error paths
- Mock Docker client for unit tests
