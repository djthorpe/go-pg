package pg

import (
	"errors"
	"fmt"

	// Packages
	pgx "github.com/jackc/pgx/v5"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Err int

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ErrSuccess Err = iota
	ErrNotFound
	ErrNotImplemented
	ErrBadParameter
	ErrNotAvailable
)

// Error returns the string representation of the error.
func (e Err) Error() string {
	switch e {
	case ErrSuccess:
		return "success"
	case ErrNotFound:
		return "not found"
	case ErrNotImplemented:
		return "not implemented"
	case ErrBadParameter:
		return "bad parameter"
	case ErrNotAvailable:
		return "not available"
	default:
		return fmt.Sprint("Unknown error ", int(e))
	}
}

// With returns the error with additional context appended.
func (e Err) With(a ...any) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprint(a...))
}

// Withf returns the error with formatted context appended.
func (e Err) Withf(format string, a ...any) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprintf(format, a...))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func pgerror(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	} else {
		return err
	}
}
