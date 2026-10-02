package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"tripo-backend/internal/apperror"
)

const (
	maxNameLength  = 100
	maxEmailLength = 254
)

var (
	emailRegex = regexp.MustCompile(`^[\w.+-]+@[\w-]+(\.[\w-]+)+$`)
	upperRegex = regexp.MustCompile(`[A-Z]`)
	digitRegex = regexp.MustCompile(`[0-9]`)
)

type passwordRule struct {
	message string
	check   func(string) bool
}

var passwordRules = []passwordRule{
	{"must be at least 8 characters", func(v string) bool { return len([]rune(v)) >= 8 }},
	{"must contain an uppercase letter", upperRegex.MatchString},
	{"must contain a number", digitRegex.MatchString},
	{"must be at most 72 bytes", func(v string) bool { return len(v) <= 72 }},
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateID(id string) []apperror.FieldError {
	if _, err := uuid.Parse(id); err != nil {
		return []apperror.FieldError{{Field: "id", Message: "must be a valid UUID"}}
	}
	return nil
}

func validateName(name string) []apperror.FieldError {
	switch {
	case name == "":
		return []apperror.FieldError{{Field: "name", Message: "is required"}}
	case len([]rune(name)) > maxNameLength:
		return []apperror.FieldError{{Field: "name", Message: "is too long"}}
	}
	return nil
}

func validateEmail(email string) []apperror.FieldError {
	switch {
	case email == "":
		return []apperror.FieldError{{Field: "email", Message: "is required"}}
	case len(email) > maxEmailLength || !emailRegex.MatchString(email):
		return []apperror.FieldError{{Field: "email", Message: "must be a valid email address"}}
	}
	return nil
}

func validatePassword(password string) []apperror.FieldError {
	if password == "" {
		return []apperror.FieldError{{Field: "password", Message: "is required"}}
	}
	var failed []apperror.FieldError
	for _, r := range passwordRules {
		if !r.check(password) {
			failed = append(failed, apperror.FieldError{Field: "password", Message: r.message})
		}
	}
	return failed
}

func validateRequired(field, value string) []apperror.FieldError {
	if strings.TrimSpace(value) == "" {
		return []apperror.FieldError{{Field: field, Message: "is required"}}
	}
	return nil
}

func collect(groups ...[]apperror.FieldError) error {
	var all []apperror.FieldError
	for _, g := range groups {
		all = append(all, g...)
	}
	if len(all) == 0 {
		return nil
	}
	return &apperror.ValidationError{Fields: all}
}

func isNotFound(err error) bool { return errors.Is(err, apperror.ErrNotFound) }
