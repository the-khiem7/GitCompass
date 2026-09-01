package git

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectReadsEffectiveContext(t *testing.T) {
	ctx := t.Context()
	directory := t.TempDir()
	runGit(t, ctx, directory, "init")
	runGit(t, ctx, directory, "config", "user.name", "Test User")
	runGit(t, ctx, directory, "config", "user.email", "test@example.com")
	runGit(t, ctx, directory, "config", "credential.helper", "manager")
	runGit(t, ctx, directory, "remote", "add", "work", "https://example.com/work/repo.git")

	runner, err := Discover(ctx)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	repository, err := runner.Inspect(ctx, directory)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if normalisePath(repository.Path) != normalisePath(directory) {
		t.Errorf("Path = %q, want %q", repository.Path, directory)
	}
	if repository.IdentityEmail != "test@example.com" {
		t.Errorf("IdentityEmail = %q", repository.IdentityEmail)
	}
	if repository.CredentialHelper != "manager" {
		t.Errorf("CredentialHelper = %q", repository.CredentialHelper)
	}
	if got := repository.Remotes["work"]; len(got) != 1 || got[0] != "https://example.com/work/repo.git" {
		t.Errorf("Remotes[work] = %#v", got)
	}
	if len(repository.Config) == 0 {
		t.Error("Config is empty")
	}
}

func normalisePath(path string) string {
	cleaned, err := filepath.EvalSymlinks(path)
	if err == nil {
		path = cleaned
	}
	return strings.ToLower(filepath.ToSlash(path))
}

func TestInspectRejectsNonRepository(t *testing.T) {
	runner, err := Discover(t.Context())
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	_, err = runner.Inspect(t.Context(), t.TempDir())
	if err != ErrNotRepository {
		t.Fatalf("Inspect() error = %v, want ErrNotRepository", err)
	}
}

func runGit(t *testing.T, ctx context.Context, directory string, args ...string) {
	t.Helper()
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = filepath.Clean(directory)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}
