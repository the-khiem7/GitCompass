package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/TheKhiem7/GitCompass/internal/profile"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Store struct {
	database *sql.DB
}

func Open(path string) (*Store, error) {
	database, err := sql.Open("sqlite", filepath.ToSlash(path)+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}
	store := &Store{database: database}
	if err := store.migrate(context.Background()); err != nil {
		_ = database.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.database.Close()
}

func (s *Store) Create(ctx context.Context, candidate profile.Profile) (profile.Profile, error) {
	if candidate.ID == "" {
		candidate.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	candidate.CreatedAt = now
	candidate.UpdatedAt = now
	_, err := s.database.ExecContext(ctx, `INSERT INTO profiles (id, name, commit_name, commit_email, https_helper_ref, ssh_key_reference, signing_key_ref, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, candidate.ID, candidate.Name, candidate.CommitName, candidate.CommitEmail, candidate.HTTPSHelperRef, candidate.SSHKeyReference, candidate.SigningKeyRef, candidate.CreatedAt, candidate.UpdatedAt)
	if err != nil {
		return profile.Profile{}, fmt.Errorf("create profile: %w", err)
	}
	return candidate, nil
}

func (s *Store) Update(ctx context.Context, candidate profile.Profile) (profile.Profile, error) {
	candidate.UpdatedAt = time.Now().UTC()
	result, err := s.database.ExecContext(ctx, `UPDATE profiles SET name = ?, commit_name = ?, commit_email = ?, https_helper_ref = ?, ssh_key_reference = ?, signing_key_ref = ?, updated_at = ? WHERE id = ?`, candidate.Name, candidate.CommitName, candidate.CommitEmail, candidate.HTTPSHelperRef, candidate.SSHKeyReference, candidate.SigningKeyRef, candidate.UpdatedAt, candidate.ID)
	if err != nil {
		return profile.Profile{}, fmt.Errorf("update profile: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return profile.Profile{}, fmt.Errorf("read profile update result: %w", err)
	}
	if changed == 0 {
		return profile.Profile{}, profile.ErrNotFound
	}
	stored, err := s.Get(ctx, candidate.ID)
	if err != nil {
		return profile.Profile{}, err
	}
	return stored, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	result, err := s.database.ExecContext(ctx, `DELETE FROM profiles WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read profile deletion result: %w", err)
	}
	if changed == 0 {
		return profile.ErrNotFound
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (profile.Profile, error) {
	row := s.database.QueryRowContext(ctx, `SELECT id, name, commit_name, commit_email, https_helper_ref, ssh_key_reference, signing_key_ref, created_at, updated_at FROM profiles WHERE id = ?`, id)
	return scanProfile(row)
}

func (s *Store) List(ctx context.Context) ([]profile.Profile, error) {
	rows, err := s.database.QueryContext(ctx, `SELECT id, name, commit_name, commit_email, https_helper_ref, ssh_key_reference, signing_key_ref, created_at, updated_at FROM profiles ORDER BY name, id`)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()
	profiles := []profile.Profile{}
	for rows.Next() {
		candidate, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, candidate)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate profiles: %w", err)
	}
	return profiles, nil
}

type scanner interface {
	Scan(...any) error
}

func scanProfile(source scanner) (profile.Profile, error) {
	var candidate profile.Profile
	err := source.Scan(&candidate.ID, &candidate.Name, &candidate.CommitName, &candidate.CommitEmail, &candidate.HTTPSHelperRef, &candidate.SSHKeyReference, &candidate.SigningKeyRef, &candidate.CreatedAt, &candidate.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return profile.Profile{}, profile.ErrNotFound
	}
	if err != nil {
		return profile.Profile{}, fmt.Errorf("read profile: %w", err)
	}
	return candidate, nil
}

func (s *Store) migrate(ctx context.Context) error {
	transaction, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin SQLite migration: %w", err)
	}
	defer transaction.Rollback()
	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`,
		`CREATE TABLE IF NOT EXISTS profiles (id TEXT PRIMARY KEY, name TEXT NOT NULL, commit_name TEXT NOT NULL, commit_email TEXT NOT NULL, https_helper_ref TEXT NOT NULL, ssh_key_reference TEXT NOT NULL, signing_key_ref TEXT NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL)`,
		`INSERT OR IGNORE INTO schema_migrations (version) VALUES (1)`,
	}
	for _, statement := range statements {
		if _, err := transaction.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply SQLite migration: %w", err)
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit SQLite migration: %w", err)
	}
	return nil
}
