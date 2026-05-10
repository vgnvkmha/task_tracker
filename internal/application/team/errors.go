package team

import (
	"errors"
)

var (
	// Team
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamInactive      = errors.New("team is inactive")

	// Leader / User
	ErrLeaderNotFound       = errors.New("leader not found")
	ErrLeaderInactive       = errors.New("leader is inactive")
	ErrLeaderAlreadyHasTeam = errors.New("leader has team already")

	// Access
	ErrPermissionDenied = errors.New("only managers can modify teams")

	// Generic service layer
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidID    = errors.New("invalid team ID")
)

func isUnexpectedError(err error) bool {
	switch {
	// Team
	case errors.Is(err, ErrTeamNotFound):
		return false

	case errors.Is(err, ErrTeamAlreadyExists):
		return false

	case errors.Is(err, ErrTeamInactive):
		return false

	// Leader / User
	case errors.Is(err, ErrLeaderNotFound):
		return false

	case errors.Is(err, ErrLeaderInactive):
		return false

	case errors.Is(err, ErrLeaderAlreadyHasTeam):
		return false

	// Access
	case errors.Is(err, ErrPermissionDenied):
		return false

	// Validation
	case errors.Is(err, ErrInvalidInput):
		return false

	case errors.Is(err, ErrInvalidID):
		return false

	default:
		return true
	}
}
