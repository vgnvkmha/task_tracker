package team

import "errors"

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
