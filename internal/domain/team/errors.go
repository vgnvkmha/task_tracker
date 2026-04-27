package team

import "errors"

var (
	// validation
	ErrEmptyName       = errors.New("team name must be provided")
	ErrNameTooLong     = errors.New("team name is too long")
	ErrInvalidTZ       = errors.New("invalid timezone")
	ErrInvalidLeaderID = errors.New("invalid team leader ID")
)
