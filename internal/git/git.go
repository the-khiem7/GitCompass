package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrNotRepository = errors.New("path is not inside a Git work tree")

type Runner struct {
	Executable string
}

type Repository struct {
	Path             string
	GitDir           string
	IdentityName     string
	IdentityEmail    string
	CredentialHelper string
	Remotes          map[string][]string
	Config           []ConfigEntry
}

type ConfigEntry struct {
	Origin string
	Scope  string
	Key    string
	Value  string
}

func Discover(ctx context.Context) (Runner, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return Runner{}, fmt.Errorf("find Git executable: %w", err)
	}
	runner := Runner{Executable: path}
	if _, err := runner.run(ctx, "--version"); err != nil {
		return Runner{}, fmt.Errorf("verify Git executable: %w", err)
	}
	return runner, nil
}

func (r Runner) Inspect(ctx context.Context, path string) (Repository, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return Repository{}, fmt.Errorf("make repository path absolute: %w", err)
	}
	topLevel, err := r.runIn(ctx, path, "rev-parse", "--show-toplevel")
	if err != nil {
		return Repository{}, ErrNotRepository
	}
	gitDir, err := r.runIn(ctx, path, "rev-parse", "--git-dir")
	if err != nil {
		return Repository{}, fmt.Errorf("read Git directory: %w", err)
	}
	configOutput, err := r.runIn(ctx, path, "config", "--show-origin", "--show-scope", "--null", "--list")
	if err != nil {
		return Repository{}, fmt.Errorf("read effective Git configuration: %w", err)
	}
	config := parseConfig(configOutput)
	repo := Repository{
		Path:    strings.TrimSpace(topLevel),
		GitDir:  strings.TrimSpace(gitDir),
		Remotes: map[string][]string{},
		Config:  config,
	}
	for _, entry := range config {
		switch entry.Key {
		case "user.name":
			repo.IdentityName = entry.Value
		case "user.email":
			repo.IdentityEmail = entry.Value
		case "credential.helper":
			repo.CredentialHelper = entry.Value
		}
	}
	remoteNames, err := r.runIn(ctx, path, "remote")
	if err != nil {
		return Repository{}, fmt.Errorf("list remotes: %w", err)
	}
	for _, name := range strings.Fields(remoteNames) {
		urls, err := r.runIn(ctx, path, "remote", "get-url", "--all", name)
		if err != nil {
			return Repository{}, fmt.Errorf("read remote %q: %w", name, err)
		}
		repo.Remotes[name] = strings.Fields(urls)
	}
	return repo, nil
}

func (r Runner) run(ctx context.Context, args ...string) (string, error) {
	return r.runIn(ctx, "", args...)
}

func (r Runner) runIn(ctx context.Context, dir string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, r.Executable, args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func parseConfig(output string) []ConfigEntry {
	fields := strings.Split(output, "\x00")
	entries := make([]ConfigEntry, 0, len(fields)/3)
	for index := 0; index+2 < len(fields); index += 3 {
		keyValue := strings.SplitN(fields[index+2], "\n", 2)
		if len(keyValue) != 2 {
			continue
		}
		entries = append(entries, ConfigEntry{Scope: fields[index], Origin: fields[index+1], Key: keyValue[0], Value: keyValue[1]})
	}
	return entries
}
