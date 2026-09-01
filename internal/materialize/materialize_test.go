package materialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheKhiem7/GitCompass/internal/profile"
	"github.com/TheKhiem7/GitCompass/internal/routing"
)

func TestApplyPreservesGlobalConfigAndIsIdempotent(t *testing.T) {
	directory := t.TempDir()
	global := filepath.Join(directory, "global.gitconfig")
	if err := os.WriteFile(global, []byte("[core]\n\teditor = code\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manager := Manager{Root: filepath.Join(directory, "managed")}
	profiles := []profile.Profile{{ID: "work", CommitName: "Work", CommitEmail: "work@example.com"}}
	rules := []routing.Rule{{ID: "default", Kind: routing.Default, ProfileID: "work"}, {ID: "folder", Kind: routing.Folder, ProfileID: "work", Pattern: "C:/Source"}}
	if err := manager.Apply(global, profiles, rules); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if err := manager.Apply(global, profiles, rules); err != nil {
		t.Fatalf("second Apply() error = %v", err)
	}
	content, err := os.ReadFile(global)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "editor = code") || strings.Count(string(content), "managed/gitconfig") != 1 {
		t.Fatalf("global config = %s", content)
	}
	root, err := os.ReadFile(filepath.Join(manager.Root, "gitconfig"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(root), "gitdir/i:C:/Source") {
		t.Fatalf("root fragment = %s", root)
	}
}
