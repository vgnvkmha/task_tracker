package valueobjects

import err "task_tracker/internal/domain/errors"

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
	Closed     Status = "closed"
)

func (s Status) IsValid() error {
	switch s {
	case Todo, InProgress, Done, Closed:
		return err.InvalidStatus
	default:
		return nil
	}
}

func (s Status) IsImmutable() error {
	switch s {
	case InProgress, Done, Closed:
		return err.AdminCanModifyOnly
	default:
		return nil
	}
}
