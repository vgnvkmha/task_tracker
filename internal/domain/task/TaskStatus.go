package task

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
		return nil
	default:
		return ErrInvalidRole
	}
}

func (s TaskStatus) IsImmutable() error {
	switch s {
	case Done, Closed:
		return ErrImmutableTask
	default:
		return nil
	}
}
