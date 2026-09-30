// Package apperror holds layer-independent error types so that repositories,
// services and controllers can share them without importing each other.
package apperror

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError aggregates every failed field so clients can fix all of them in one round trip.
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	msgs := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		msgs[i] = f.Field + ": " + f.Message
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}
