package validatorx_test

import (
	"strings"
	"testing"

	"github.com/minhhoccode111/go-clean-template-gin/pkg/validatorx"
)

// ---- no_dups_str ------------------------------------------------------------

func TestNoDupsStr(t *testing.T) {
	t.Parallel()

	v := validatorx.New()

	type payload struct {
		Tags []string `validate:"no_dups_str"`
	}

	tests := []struct {
		name  string
		input []string
		valid bool
	}{
		{"unique values", []string{"go", "rust", "python"}, true},
		{"duplicate values", []string{"go", "go"}, false},
		{"trimmed duplicates", []string{"go", " go"}, false}, // " go" trimmed == "go"
		{"single element", []string{"go"}, true},
		{"empty slice", []string{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(payload{Tags: tc.input})

			got := err == nil
			if got != tc.valid {
				t.Errorf("no_dups_str(%v): want valid=%v, got valid=%v (err: %v)", tc.input, tc.valid, got, err)
			}
		})
	}
}

// ---- tag --------------------------------------------------------------------

func TestTag(t *testing.T) {
	t.Parallel()

	v := validatorx.New()

	type payload struct {
		T string `validate:"required,tag"`
	}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"simple word", "golang", true},
		{"with hyphen", "sci-fi", true},
		{"with underscore", "my_tag", true},
		{"with internal space", "golang 101", true},
		{"leading space", " golang", false},
		{"trailing space", "golang ", false},
		{"only space", " ", false},
		{"special characters", "go@lang", false},
		{"single letter", "g", true},
		{"unicode letter", "日本語", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(payload{T: tc.input})

			got := err == nil
			if got != tc.valid {
				t.Errorf("tag(%q): want valid=%v, got valid=%v (err: %v)", tc.input, tc.valid, got, err)
			}
		})
	}
}

// ---- username ---------------------------------------------------------------

func TestUsername(t *testing.T) {
	t.Parallel()

	v := validatorx.New()

	type payload struct {
		U string `validate:"required,username"`
	}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"alphanumeric", "minhhoccode111", true},
		{"letters only", "john", true},
		{"digits only", "12345", true},
		{"unicode letters", "Ψuser42", true},
		{"with space", "john doe", false},
		{"with hyphen", "john-doe", false},
		{"with underscore", "john_doe", false},
		{"with special char", "john@doe", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(payload{U: tc.input})

			got := err == nil
			if got != tc.valid {
				t.Errorf("username(%q): want valid=%v, got valid=%v (err: %v)", tc.input, tc.valid, got, err)
			}
		})
	}
}

// ---- password ---------------------------------------------------------------

func TestPassword(t *testing.T) {
	t.Parallel()

	v := validatorx.New()

	type payload struct {
		P string `validate:"required,password"`
	}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"all requirements met", "P@ssw0rd", true},
		{"missing uppercase", "p@ssw0rd", false},
		{"missing lowercase", "P@SSW0RD", false},
		{"missing digit", "P@ssword", false},
		{"missing special char", "Passw0rd", false},
		{"only letters", "Password", false},
		{"empty string", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(payload{P: tc.input})

			got := err == nil
			if got != tc.valid {
				t.Errorf("password(%q): want valid=%v, got valid=%v (err: %v)", tc.input, tc.valid, got, err)
			}
		})
	}
}

// ---- ExtractErrors ----------------------------------------------------------

type extractErrorsPayload struct {
	Email    string   `validate:"required,email"`
	Username string   `validate:"required,min=2,max=50,username"`
	Password string   `validate:"required,min=8,max=50,password"` //nolint:gosec // test struct, not a real credential store
	Tags     []string `validate:"no_dups_str"`
}

func extractErrorsCases() []struct {
	name        string
	input       extractErrorsPayload
	wantMessage string
} {
	return []struct {
		name        string
		input       extractErrorsPayload
		wantMessage string
	}{
		{
			"required field missing",
			extractErrorsPayload{},
			"Email is required",
		},
		{
			"invalid email",
			extractErrorsPayload{Email: "not-an-email", Username: "validuser", Password: "P@ssw0rd"},
			"must be a valid email address",
		},
		{
			"username too short",
			extractErrorsPayload{Email: "a@b.com", Username: "x", Password: "P@ssw0rd"},
			"must be at least",
		},
		{
			"invalid username characters",
			extractErrorsPayload{Email: "a@b.com", Username: "bad user!", Password: "P@ssw0rd"},
			"must contain only letters",
		},
		{
			"weak password",
			extractErrorsPayload{Email: "a@b.com", Username: "validuser", Password: "weakpass"},
			"uppercase",
		},
		{
			"duplicate tags",
			extractErrorsPayload{Email: "a@b.com", Username: "validuser", Password: "P@ssw0rd", Tags: []string{"go", "go"}},
			"contains duplicate values",
		},
	}
}

func TestExtractErrors(t *testing.T) {
	t.Parallel()

	v := validatorx.New()

	for _, tc := range extractErrorsCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(tc.input)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}

			msgs := validatorx.ExtractErrors(err)
			found := false

			for _, m := range msgs {
				if contains(m, tc.wantMessage) {
					found = true

					break
				}
			}

			if !found {
				t.Errorf("ExtractErrors: want a message containing %q, got %v", tc.wantMessage, msgs)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
