package team

import (
	"errors"
	"net/http"
	"task_tracker/internal/application/team"
	domain_team "task_tracker/internal/domain/team"
)

func mapError(err error) (int, string) {
	switch {

	// Team
	case errors.Is(err, team.ErrTeamAlreadyExists):
		return http.StatusConflict, "team already exists"

	case errors.Is(err, team.ErrTeamNotFound):
		return http.StatusNotFound, "team not found"

	case errors.Is(err, team.ErrTeamInactive):
		return http.StatusBadRequest, "team is inactive"

	// Leader / User
	case errors.Is(err, team.ErrLeaderNotFound):
		return http.StatusBadRequest, "leader not found"

	case errors.Is(err, team.ErrLeaderInactive):
		return http.StatusBadRequest, "leader inactive"

	case errors.Is(err, team.ErrLeaderAlreadyHasTeam):
		return http.StatusConflict, "leader already has a team"

	// Access
	case errors.Is(err, team.ErrPermissionDenied):
		return http.StatusForbidden, "permission denied"

	// Generic
	case errors.Is(err, team.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, domain_team.ErrInvalidTZ):
		return http.StatusBadRequest, "invalid timezone"
	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
