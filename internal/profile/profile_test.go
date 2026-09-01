package profile

import (
	"errors"
	"testing"
)

func TestValidateRejectsIncompleteAndInvalidProfiles(t *testing.T) {
	cases := []Profile{
		{},
		{Name: "Work", CommitName: "User", CommitEmail: "invalid", HTTPSHelperRef: "manager"},
		{Name: "Work", CommitName: "User", CommitEmail: "work@example.com"},
	}
	for _, candidate := range cases {
		if err := Validate(candidate); !errors.Is(err, ErrInvalid) {
			t.Errorf("Validate(%#v) error = %v, want ErrInvalid", candidate, err)
		}
	}
}
