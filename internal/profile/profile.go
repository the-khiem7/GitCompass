package profile

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("profile not found")
	ErrInvalid  = errors.New("invalid profile")
)

type Profile struct {
	ID              string
	Name            string
	CommitName      string
	CommitEmail     string
	HTTPSHelperRef  string
	SSHKeyReference string
	SigningKeyRef   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Repository interface {
	Create(context.Context, Profile) (Profile, error)
	Update(context.Context, Profile) (Profile, error)
	Delete(context.Context, string) error
	Get(context.Context, string) (Profile, error)
	List(context.Context) ([]Profile, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) Create(ctx context.Context, candidate Profile) (Profile, error) {
	if err := Validate(candidate); err != nil {
		return Profile{}, err
	}
	return s.repository.Create(ctx, candidate)
}

func (s Service) Update(ctx context.Context, candidate Profile) (Profile, error) {
	if strings.TrimSpace(candidate.ID) == "" {
		return Profile{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if err := Validate(candidate); err != nil {
		return Profile{}, err
	}
	return s.repository.Update(ctx, candidate)
}

func (s Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repository.Delete(ctx, id)
}

func (s Service) Get(ctx context.Context, id string) (Profile, error) {
	if strings.TrimSpace(id) == "" {
		return Profile{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Profile, error) {
	return s.repository.List(ctx)
}

func Validate(candidate Profile) error {
	if strings.TrimSpace(candidate.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if strings.TrimSpace(candidate.CommitName) == "" {
		return fmt.Errorf("%w: commit name is required", ErrInvalid)
	}
	address, err := mail.ParseAddress(candidate.CommitEmail)
	if err != nil || address.Address != candidate.CommitEmail {
		return fmt.Errorf("%w: commit email is invalid", ErrInvalid)
	}
	if strings.TrimSpace(candidate.HTTPSHelperRef) == "" && strings.TrimSpace(candidate.SSHKeyReference) == "" {
		return fmt.Errorf("%w: HTTPS helper reference or SSH key reference is required", ErrInvalid)
	}
	return nil
}
