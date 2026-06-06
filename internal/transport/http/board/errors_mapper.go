package board

import (
	"errors"
	"net/http"

	boardapp "task_tracker/internal/application/board"
	"task_tracker/internal/common_errors"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, boardapp.ErrBoardNotFound):
		return http.StatusNotFound, "board not found"
	case errors.Is(err, boardapp.ErrBoardAlreadyExists):
		return http.StatusConflict, "board already exists"
	case errors.Is(err, boardapp.ErrPermissionDenied):
		return http.StatusForbidden, "permission denied"
	case errors.Is(err, boardapp.ErrInvalidBoardID),
		errors.Is(err, boardapp.ErrInvalidInput),
		errors.Is(err, boardapp.ErrInvalidStatus),
		errors.Is(err, common_errors.ErrInvalidID):
		return http.StatusUnprocessableEntity, err.Error()
	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
