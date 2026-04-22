package team

import "errors"

var (
	// general
	ErrNotFound      = errors.New("team not found")
	ErrAlreadyExists = errors.New("team already exists")

	// validation
	ErrEmptyName   = errors.New("team name must be provided")
	ErrNameTooLong = errors.New("team name is too long")
	ErrInvalidTZ   = errors.New("invalid timezone")

	// leader
	ErrLeaderNotFound  = errors.New("leader not found")
	ErrInvalidLeader   = errors.New("invalid leader")
	ErrLeaderNotInTeam = errors.New("leader must be a member of the team")

	// members
	ErrMemberAlreadyExists = errors.New("user already in team")
	ErrMemberNotFound      = errors.New("user is not a team member")

	// state
	ErrInactiveTeam = errors.New("team is inactive")

	// logical
	ErrCannotRemoveLeader = errors.New("cannot remove team leader")
)
