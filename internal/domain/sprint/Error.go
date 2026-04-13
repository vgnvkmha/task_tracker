package sprint

import "errors"

var (
	// general
	ErrNotFound      = errors.New("sprint not found")
	ErrAlreadyExists = errors.New("sprint already exists")

	// validation
	ErrEmptyName        = errors.New("sprint name must be provided")
	ErrInvalidDateRange = errors.New("sprint start date must be before end date")
	ErrInvalidDuration  = errors.New("invalid sprint duration")

	// status
	ErrInvalidStatus           = errors.New("invalid sprint status")
	ErrInvalidStatusTransition = errors.New("invalid sprint status transition")
	ErrSprintNotActive         = errors.New("sprint is not active")
	ErrSprintAlreadyCompleted  = errors.New("sprint is already completed")

	// logical
	ErrSprintNotStarted     = errors.New("sprint has not started yet")
	ErrSprintAlreadyStarted = errors.New("sprint is already started")
	ErrCannotModifySprint   = errors.New("cannot modify sprint in current state")

	// relations
	ErrTaskNotInSprint = errors.New("task does not belong to sprint")
	ErrSprintClosed    = errors.New("sprint is closed")
)
