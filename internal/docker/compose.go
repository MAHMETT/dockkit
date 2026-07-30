package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ComposeUp runs docker compose up -d in the given directory.
func ComposeUp(ctx context.Context, serviceDir string) error {
	return composeExec(ctx, serviceDir, "up", "-d")
}

// ComposeDown runs docker compose down in the given directory.
func ComposeDown(ctx context.Context, serviceDir string) error {
	return composeExec(ctx, serviceDir, "down")
}

// ComposeRestart runs docker compose restart in the given directory.
func ComposeRestart(ctx context.Context, serviceDir string) error {
	return composeExec(ctx, serviceDir, "restart")
}

// ComposeLogs runs docker compose logs in the given directory.
func ComposeLogs(ctx context.Context, serviceDir string) (string, error) {
	return composeExecOutput(ctx, serviceDir, "logs", "--tail=100")
}

// ComposePS runs docker compose ps in the given directory.
func ComposePS(ctx context.Context, serviceDir string) (string, error) {
	return composeExecOutput(ctx, serviceDir, "ps")
}

// composeExec runs a docker compose command with timeout.
func composeExec(ctx context.Context, dir string, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	composeArgs := buildComposeArgs(dir, args)
	cmd := exec.CommandContext(ctx, "docker", composeArgs...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose %s: %w\n%s", args[0], err, stderr.String())
	}
	return nil
}

// composeExecOutput runs a docker compose command and returns stdout.
func composeExecOutput(ctx context.Context, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	composeArgs := buildComposeArgs(dir, args)
	cmd := exec.CommandContext(ctx, "docker", composeArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker compose %s: %w\n%s", args[0], err, stderr.String())
	}
	return stdout.String(), nil
}

// buildComposeArgs builds the docker compose command arguments.
func buildComposeArgs(dir string, args []string) []string {
	composeFile := filepath.Join(dir, "docker-compose.yml")
	envFile := filepath.Join(dir, ".env")

	result := []string{"compose", "-f", composeFile}

	if _, err := os.Stat(envFile); err == nil {
		result = append(result, "--env-file", envFile)
	}

	return append(result, args...)
}
