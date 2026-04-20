package user

import "errors"

var (
	// general
	ErrCreateUserFailed      = errors.New("failed to create user")
	ErrTransactionFailed     = errors.New("transaction failed")
	ErrOnlyManagersCanCreate = errors.New("only managers can create new users")
	ErrUserNotFound          = errors.New("user was not found")
	ErrUserUpdateFailed      = errors.New("user update failed")

	// input / orchestration
	ErrRoleRequired = errors.New("user role must be provided")
	ErrInvalidRole  = errors.New("input role is invalid")
	ErrInvalidInput = errors.New("invalid user input")

	// team-related
	ErrTeamNotFound    = errors.New("team not found")
	ErrTeamFetchFailed = errors.New("failed to fetch team")

	// personal data
	ErrPersonalDataCreateFailed = errors.New("failed to create personal data")
	ErrPersonalDataUpdateFailed = errors.New("failed to update personal data")
	ErrPersonalDataNotFound     = errors.New("personal data was not found")

	// user persistence
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserCreateFailed  = errors.New("failed to persist user")
)
