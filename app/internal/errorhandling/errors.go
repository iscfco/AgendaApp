package errorhandling

import (
	"errors"
)

var (
	ErrInternal        = errors.New("internal_error")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrDuplicatedError = errors.New("duplicated_error")
	ErrNotFoundError   = errors.New("not_found")
)
