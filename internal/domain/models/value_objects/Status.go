package valueobjects

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
	Closed     Status = "closed"
)

func (s Status) IsValid() bool {
	switch s {
	case Todo, InProgress, Done, Closed:
		return true
	default:
		return false
	}
}

func (s Status) IsImmutable() bool {
	switch s {
	case InProgress, Done, Closed:
		return true
	default:
		return false
	}
}
