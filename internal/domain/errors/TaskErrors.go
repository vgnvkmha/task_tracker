package errors_task

import "errors"

var (
	ErrInvalidStatus           = errors.New("invalid status: should be to do, in progress, etc")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidRights           = errors.New("only administrators or team leaders can modify task at this stage")
	ErrInvalidRole             = errors.New("invalid user type")
	ErrImmutableTask           = errors.New("task is immutable at the current state")
)
