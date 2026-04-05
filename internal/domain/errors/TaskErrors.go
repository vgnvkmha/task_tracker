package task_errors

import "errors"

var (
	ErrInvalidStatus           = errors.New("invalid status: should be to do, in progress, etc")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidRights           = errors.New("only administrators or team leaders can modify task at this stage")
	ErrInvalidRole             = errors.New("invalid user type")
	ErrImmutableTask           = errors.New("task is immutable at the current state")
	ErrTaskName                = errors.New("task name must be set")
	ErrTaskBoard               = errors.New("task board must exist and be running")
	ErrTaskUser                = errors.New("all users must be active")
	ErrInvalidTime             = errors.New("date must be in the future")
	ErrSameChange              = errors.New("nothing to change")
)
