package persistence

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/TheKhiem7/GitCompass/internal/profile"
)

func TestStoreProfileCRUD(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "gitcompass.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	service := profile.NewService(store)
	created, err := service.Create(t.Context(), profile.Profile{Name: "Work", CommitName: "Work User", CommitEmail: "work@example.com", HTTPSHelperRef: "manager"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	created.CommitName = "Updated User"
	updated, err := service.Update(t.Context(), created)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.CommitName != "Updated User" || updated.CreatedAt.IsZero() || updated.UpdatedAt.IsZero() {
		t.Errorf("Update() = %#v", updated)
	}
	profiles, err := service.List(t.Context())
	if err != nil || len(profiles) != 1 {
		t.Fatalf("List() = %#v, %v", profiles, err)
	}
	if err := service.Delete(t.Context(), created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = service.Get(t.Context(), created.ID)
	if !errors.Is(err, profile.ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}
