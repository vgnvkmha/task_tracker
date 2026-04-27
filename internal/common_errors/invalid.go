package common_errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidState    = errors.New("invalid state")
)
