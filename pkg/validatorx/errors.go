package validatorx

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ExtractErrors extracts human-readable error messages from validator.ValidationErrors.
//
//nolint:gocyclo,cyclop,gocritic,nolintlint
func ExtractErrors(err error) []string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return []string{err.Error()}
	}

	out := make([]string, 0, len(verrs))

	for _, e := range verrs {
		field := e.Field()
		param := e.Param()

		switch e.Tag() {
		case "required":
			out = append(out, field+" is required")

		case "max":
			out = append(out, field+" must be at most "+param+" characters")

		case "min":
			out = append(out, field+" must be at least "+param+" characters")

		case "email":
			out = append(out, field+" must be a valid email address")

		case "username":
			out = append(out, field+" must contain only letters, and numbers")

		case "password":
			out = append(
				out,
				field+" must contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
			)

		case "tag":
			out = append(
				out,
				field+` allowed characters are letters, numbers, spaces, underscores, and hyphens`,
			)

		case "no_dups_str":
			out = append(out, field+" contains duplicate values")

		default:
			out = append(out, fmt.Sprintf("%s is invalid (%s)", field, e.Tag()))
		}
	}

	return out
}
