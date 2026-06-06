package task_handler

import (
	"errors"
	"net/http"

	taskApplication "task_tracker/internal/application/task"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, taskApplication.ErrTaskNotFound):
		return http.StatusNotFound, "task not found"
	case errors.Is(err, taskApplication.ErrPermissionDenied):
		return http.StatusForbidden, "permission denied"
	case errors.Is(err, taskApplication.ErrInvalidInput):
		return http.StatusUnprocessableEntity, "invalid task input"
	case errors.Is(err, taskApplication.ErrInvalidStatus):
		return http.StatusUnprocessableEntity, "invalid task status"
	case errors.Is(err, taskApplication.ErrInvalidTransition):
		return http.StatusUnprocessableEntity, "invalid task status transition"
	case errors.Is(err, taskApplication.ErrInvalidTaskID):
		return http.StatusUnprocessableEntity, "invalid task id"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
