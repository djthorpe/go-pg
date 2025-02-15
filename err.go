package pg

import "fmt"

/////////////////////////////////////////////////////////////////////
// TYPES

type Err int

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ErrSuccess Err = iota
	ErrNotFound
	ErrNotImplemented
)

func (e Err) Error() string {
	switch e {
	case ErrSuccess:
		return "success"
	case ErrNotFound:
		return "not found"
	case ErrNotImplemented:
		return "not implemented"
	default:
		return fmt.Sprint("Unknown error ", int(e))
	}
}
