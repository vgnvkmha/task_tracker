package common_errors

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrConflict      = errors.New("conflict")
)
