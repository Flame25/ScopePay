package httpx

import "errors"

type StateError struct{ Message string }

func (e *StateError) Error() string { return e.Message }

func NewStateError(msg string) error { return &StateError{Message: msg} }

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
)
