package user

import "errors"

var (
	// general
	ErrCreateUserFailed      = errors.New("failed to create user")
	ErrTransactionFailed     = errors.New("transaction failed")
	ErrOnlyManagersCanModify = errors.New("only managers can modify users")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserUpdateFailed      = errors.New("user update failed")
	ErrConflict              = errors.New("Conflict with user")

	// input / orchestration
	ErrRoleRequired     = errors.New("user role must be provided")
	ErrInvalidRole      = errors.New("input role is invalid")
	ErrInvalidInput     = errors.New("invalid user input")
	ErrInvalidUserID    = errors.New("invalid user_id")
	ErrInvalidTeamID    = errors.New("invalid team_id")
	ErrInvalidBirthDate = errors.New("invalid birth_date")

	// team-related
	ErrTeamNotFound    = errors.New("team not found")
	ErrTeamFetchFailed = errors.New("failed to fetch team")

	// personal data
	ErrPersonalDataCreateFailed = errors.New("failed to create personal data")
	ErrPersonalDataUpdateFailed = errors.New("failed to update personal data")
	ErrPersonalDataNotFound     = errors.New("personal data was not found")

	// user persistence
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPersistUser       = errors.New("failed to persist user")
)
