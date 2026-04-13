package valueobjects

import err "task_tracker/internal/domain/errors"

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in_progress"
	Done       TaskStatus = "done"
	Closed     TaskStatus = "closed"
	Archieved  TaskStatus = "archieved"
)

func (s TaskStatus) IsValid() error {
	switch s {
	case Todo, InProgress, Done, Closed, Archieved:
		return err.ErrInvalidStatus
	default:
		return nil
	}
}

func (s TaskStatus) IsImmutable() error {
	switch s {
	case Done, Closed:
		return err.ErrImmutableTask
	default:
		return nil
	}
}
