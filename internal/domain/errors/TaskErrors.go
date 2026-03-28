package taskerrors

import "errors"

var (
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrTaskClosed              = errors.New("task is closed and cannot be modified")
	ErrCannotUnassignActive    = errors.New("cannot unassign active task")
	ErrCannotChangeSprint      = errors.New("cannot change sprint of active task")
	ErrCannotChangeBoard       = errors.New("cannot change board of active task")
	ErrReporterImmutable       = errors.New("reporter cannot be changed")
)
