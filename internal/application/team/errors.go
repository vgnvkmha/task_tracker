package team

import (
	"errors"
	"fmt"
	"task_tracker/internal/domain/team"
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
	ErrPermissionDenied = errors.New("only managers can create teams")

	// Generic service layer
	ErrInvalidInput = errors.New("invalid input")
)

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, team.ErrInvalidTZ):
		return fmt.Errorf("invalid timezone: %w", ErrInvalidInput)

	case errors.Is(err, team.ErrInvalidLeaderID):
		return fmt.Errorf("invalid leader ID: %w", ErrInvalidInput)

	case errors.Is(err, team.ErrEmptyName):
		return fmt.Errorf("empty name: %w", ErrInvalidInput)
	case errors.Is(err, team.ErrNameTooLong):
		return fmt.Errorf("name too long: %w", ErrInvalidInput)
	default:
		return err
	}
}
