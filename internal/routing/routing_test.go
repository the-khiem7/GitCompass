package routing

import "testing"

func TestResolveUsesDeterministicPrecedence(t *testing.T) {
	rules := []Rule{{ID: "default", Kind: Default, ProfileID: "personal"}, {ID: "folder", Kind: Folder, ProfileID: "work", Pattern: "C:/Source"}, {ID: "remote", Kind: Remote, ProfileID: "team", Pattern: "github.com/acme"}, {ID: "exact", Kind: ExactRepository, ProfileID: "special", Pattern: "c:/source/app"}}
	result := Resolve(Repository{Path: "C:\\SOURCE\\App", RemoteURLs: []string{"https://github.com/acme/app.git"}}, rules)
	if result.ProfileID != "special" || result.Rule.ID != "exact" || !result.Materialize {
		t.Fatalf("Resolve() = %#v", result)
	}
}

func TestResolveBlocksConflictingRemoteProfiles(t *testing.T) {
	rules := []Rule{{Kind: Remote, ProfileID: "work", Pattern: "github.com/acme"}, {Kind: Remote, ProfileID: "personal", Pattern: "gitlab.com/home"}}
	result := Resolve(Repository{Path: "C:/repo", RemoteURLs: []string{"https://github.com/acme/repo.git", "git@gitlab.com:home/repo.git"}}, rules)
	if !result.Conflict || result.Materialize {
		t.Fatalf("Resolve() = %#v", result)
	}
}

func TestResolveFallsBackToDefault(t *testing.T) {
	result := Resolve(Repository{Path: "C:/other"}, []Rule{{Kind: Default, ProfileID: "personal"}})
	if result.ProfileID != "personal" || result.Rule.Kind != Default {
		t.Fatalf("Resolve() = %#v", result)
	}
}
