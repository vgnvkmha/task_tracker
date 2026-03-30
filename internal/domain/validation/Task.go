package validation

import (
	vo "task_tracker/internal/domain/models/value_objects"
)

func IsValidStatusTransition(from, to vo.Status) bool {
	if err1, err2 := from.IsValid(), to.IsValid(); err1 != nil || err2 != nil {
		return false
	}

	switch from {
	case vo.Todo:
		return to == vo.InProgress
	case vo.InProgress:
		return to == vo.Done
	case vo.Done:
		return to == vo.Closed
	case vo.Closed:
		return false
	}

	return false
}
