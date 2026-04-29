package team

import (
	"errors"
	"net/http"
	"task_tracker/internal/application/team"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, team.ErrTeamAlreadyExists):
		return http.StatusConflict, "team already exists"

	case errors.Is(err, team.ErrLeaderNotFound):
		return http.StatusBadRequest, "leader not found"

	case errors.Is(err, team.ErrLeaderInactive):
		return http.StatusBadRequest, "leader inactive"

	case errors.Is(err, team.ErrPermissionDenied):
		return http.StatusForbidden, "permission denied"

	case errors.Is(err, team.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
