package task

var allowedTransitions = map[TaskStatus]TaskStatus{
	Todo:       InProgress,
	InProgress: Done,
	Done:       Closed,
}

func IsValidStatusTransition(from, to TaskStatus) error {
	if from.IsValid() != nil || to.IsValid() != nil {
		return ErrInvalidStatusTransition
	}

	if from == Closed {
		return ErrImmutableTask
	}

	if allowedTransitions[from] == to {
		return nil
	}

	return ErrInvalidStatusTransition
}
