package materialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheKhiem7/GitCompass/internal/profile"
	"github.com/TheKhiem7/GitCompass/internal/routing"
)

type Manager struct{ Root string }

func (m Manager) Apply(globalConfig string, profiles []profile.Profile, rules []routing.Rule) error {
	if strings.TrimSpace(m.Root) == "" {
		return fmt.Errorf("managed configuration root is required")
	}
	if err := os.MkdirAll(filepath.Join(m.Root, "profiles"), 0o700); err != nil {
		return fmt.Errorf("create managed configuration directory: %w", err)
	}
	byID := map[string]profile.Profile{}
	for _, candidate := range profiles {
		byID[candidate.ID] = candidate
		if err := write(filepath.Join(m.Root, "profiles", candidate.ID+".gitconfig"), profileFragment(candidate)); err != nil {
			return err
		}
	}
	root, err := rootFragment(m.Root, byID, rules)
	if err != nil {
		return err
	}
	if err := write(filepath.Join(m.Root, "gitconfig"), root); err != nil {
		return err
	}
	return ensureInclude(globalConfig, filepath.Join(m.Root, "gitconfig"))
}

func profileFragment(candidate profile.Profile) string {
	return fmt.Sprintf("[user]\n\tname = %s\n\temail = %s\n", candidate.CommitName, candidate.CommitEmail)
}

func rootFragment(root string, profiles map[string]profile.Profile, rules []routing.Rule) (string, error) {
	var builder strings.Builder
	for _, kind := range []routing.RuleKind{routing.Default, routing.Folder, routing.Remote, routing.ExactRepository} {
		for _, rule := range rules {
			if rule.Kind != kind {
				continue
			}
			if _, ok := profiles[rule.ProfileID]; !ok {
				return "", fmt.Errorf("rule %q references an unknown Profile", rule.ID)
			}
			fragment := filepath.ToSlash(filepath.Join(root, "profiles", rule.ProfileID+".gitconfig"))
			switch kind {
			case routing.Default:
				fmt.Fprintf(&builder, "[include]\n\tpath = %s\n", fragment)
			case routing.Folder, routing.ExactRepository:
				fmt.Fprintf(&builder, "[includeIf \"gitdir/i:%s\"]\n\tpath = %s\n", filepath.ToSlash(rule.Pattern), fragment)
			case routing.Remote:
				fmt.Fprintf(&builder, "[includeIf \"hasconfig:remote.*.url:%s\"]\n\tpath = %s\n", rule.Pattern, fragment)
			}
		}
	}
	return builder.String(), nil
}

func ensureInclude(globalConfig, managed string) error {
	content, err := os.ReadFile(globalConfig)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read global Git config: %w", err)
	}
	line := "[include]\n\tpath = " + filepath.ToSlash(managed) + "\n"
	if strings.Contains(string(content), filepath.ToSlash(managed)) {
		return nil
	}
	return write(globalConfig, string(content)+line)
}

func write(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write managed configuration: %w", err)
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("replace managed configuration: %w", err)
	}
	return nil
}
