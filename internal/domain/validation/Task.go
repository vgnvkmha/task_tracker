package validation

import (
	task_errors "task_tracker/internal/domain/errors"
	vo "task_tracker/internal/domain/models/value_objects"
)

var allowedTransitions = map[vo.Status]vo.Status{
	vo.Todo:       vo.InProgress,
	vo.InProgress: vo.Done,
	vo.Done:       vo.Closed,
}

func IsValidStatusTransition(from, to vo.Status) error {
	if from.IsValid() != nil || to.IsValid() != nil {
		return task_errors.ErrInvalidStatusTransition
	}

	if from == vo.Closed {
		return task_errors.ErrImmutableTask
	}

	if allowedTransitions[from] == to {
		return nil
	}

	return task_errors.ErrInvalidStatusTransition
}
