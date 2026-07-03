package models

import "errors"

// Common domain errors.
var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrBadRequest      = errors.New("bad request")
	ErrInternalServer  = errors.New("internal server error")
)
