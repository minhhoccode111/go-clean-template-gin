package validatorx

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Precompiled regexes — compiled once at package init to avoid repeated allocation on every validation call.
//
//nolint:gochecknoglobals // package-level regex vars are intentional: compiled once, read-only after init.
var (
	// reTag: letters (any language), digits, underscores, hyphens, and internal spaces.
	// Must start and end with a letter, digit, underscore, or hyphen — no leading/trailing spaces.
	// Example matches: "sci-fi", "golang 101", "my_tag".
	reTag = regexp.MustCompile(`^[\p{L}0-9_-]([\p{L}0-9_ -]*[\p{L}0-9_-])?$`)

	// reUsername: letters (any language) and digits only. No spaces or special characters.
	// Example matches: "minhhoccode111", "Ψuser42".
	reUsername = regexp.MustCompile(`^[\p{L}0-9]+$`)

	// rePassword: Go's regexp uses RE2, which does not support lookaheads, so each
	// character-class rule is expressed as a separate sub-regex and checked individually.
	// All four must match for the password to be valid.
	rePassword = struct {
		upper, lower, digit, special *regexp.Regexp
	}{
		upper:   regexp.MustCompile(`[A-Z]`),                     // at least one uppercase letter
		lower:   regexp.MustCompile(`[a-z]`),                     // at least one lowercase letter
		digit:   regexp.MustCompile(`\d`),                        // at least one digit
		special: regexp.MustCompile(`[!@#~$%^&*()+|_{}<>?,./-]`), // at least one special character
	}
)

// New creates a new validator instance with custom validations.
func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	mustRegister(v, "no_dups_str", func(fl validator.FieldLevel) bool {
		slices, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}

		seen := make(map[string]struct{})

		for _, t := range slices {
			t = strings.TrimSpace(t)
			if _, exists := seen[t]; exists {
				return false
			}

			seen[t] = struct{}{}
		}

		return true
	})

	mustRegister(v, "tag", func(fl validator.FieldLevel) bool {
		return reTag.MatchString(fl.Field().String())
	})

	mustRegister(v, "username", func(fl validator.FieldLevel) bool {
		return reUsername.MatchString(fl.Field().String())
	})

	mustRegister(v, "password", func(fl validator.FieldLevel) bool {
		p := fl.Field().String()

		return rePassword.upper.MatchString(p) &&
			rePassword.lower.MatchString(p) &&
			rePassword.digit.MatchString(p) &&
			rePassword.special.MatchString(p)
	})

	return v
}

// mustRegister registers a custom validation and panics if it fails.
// Failures only occur when the tag name clashes with a built-in validator tag,
// so a panic here indicates a programming error that should be fixed at compile time.
func mustRegister(v *validator.Validate, tag string, fn validator.Func) {
	if err := v.RegisterValidation(tag, fn); err != nil {
		panic("validatorx: failed to register '" + tag + "': " + err.Error())
	}
}
