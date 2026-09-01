package routing

import (
	"path/filepath"
	"strings"
)

type RuleKind string

const (
	ExactRepository RuleKind = "exact-repository"
	Remote          RuleKind = "remote"
	Folder          RuleKind = "folder"
	Default         RuleKind = "default"
)

type Rule struct {
	ID        string
	Kind      RuleKind
	ProfileID string
	Pattern   string
}

type Repository struct {
	Path       string
	RemoteURLs []string
	FetchURL   string
	PushURL    string
}

type Result struct {
	ProfileID   string
	Rule        Rule
	Explanation string
	Conflict    bool
	Materialize bool
}

func Resolve(repository Repository, rules []Rule) Result {
	path := normalisePath(repository.Path)
	for _, kind := range []RuleKind{ExactRepository, Remote, Folder, Default} {
		matches := matchingRules(path, repository.RemoteURLs, rules, kind)
		if len(matches) == 0 {
			continue
		}
		if kind == Remote && remoteProfilesConflict(matches) {
			return Result{Explanation: "remote rules resolve to different Profiles", Conflict: true}
		}
		winner := matches[0]
		return Result{ProfileID: winner.ProfileID, Rule: winner, Explanation: string(winner.Kind) + " rule selected the Profile", Materialize: true}
	}
	return Result{Explanation: "no routing rule matched"}
}

func matchingRules(path string, remoteURLs []string, rules []Rule, kind RuleKind) []Rule {
	matches := []Rule{}
	for _, rule := range rules {
		if rule.Kind != kind {
			continue
		}
		switch kind {
		case ExactRepository:
			if path == normalisePath(rule.Pattern) {
				matches = append(matches, rule)
			}
		case Folder:
			folder := normalisePath(rule.Pattern)
			if path == folder || strings.HasPrefix(path, folder+"/") {
				matches = append(matches, rule)
			}
		case Remote:
			for _, url := range remoteURLs {
				if strings.Contains(normaliseRemote(url), normaliseRemote(rule.Pattern)) {
					matches = append(matches, rule)
					break
				}
			}
		case Default:
			matches = append(matches, rule)
		}
	}
	return matches
}

func normaliseRemote(remote string) string {
	remote = strings.ToLower(strings.TrimSpace(remote))
	if at := strings.Index(remote, "@"); at >= 0 && !strings.Contains(remote[:at], "://") {
		remote = remote[at+1:]
	}
	return strings.ReplaceAll(remote, ":", "/")
}

func remoteProfilesConflict(rules []Rule) bool {
	profileID := rules[0].ProfileID
	for _, rule := range rules[1:] {
		if rule.ProfileID != profileID {
			return true
		}
	}
	return false
}

func normalisePath(path string) string {
	return strings.TrimSuffix(strings.ToLower(filepath.ToSlash(filepath.Clean(path))), "/")
}
