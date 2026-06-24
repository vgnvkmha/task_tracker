package task

import "errors"

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrCreateTaskFailed = errors.New("failed to create task")
	ErrUpdateTaskFailed = errors.New("failed to update task")
	ErrDeleteTaskFailed = errors.New("failed to delete task")

	ErrInvalidInput      = errors.New("invalid task input")
	ErrInvalidTaskID     = errors.New("invalid task id")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrReporterRequired  = errors.New("task reporter is required")
	ErrBoardRequired     = errors.New("task board is required")
	ErrBoardNotFound     = errors.New("task board not found")
	ErrActorTeamRequired = errors.New("user team is required")
	ErrBoardTeamMismatch = errors.New("board belongs to another team")
	ErrInvalidStatus     = errors.New("invalid task status")
	ErrInvalidTransition = errors.New("invalid task status transition")
)
