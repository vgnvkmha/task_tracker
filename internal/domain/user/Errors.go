package user

import "errors"

var (
	// general
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")

	// validation
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrEmptyEmail    = errors.New("email must be provided")
	ErrEmptyPassword = errors.New("password must be provided")
	ErrWeakPassword  = errors.New("password does not meet security requirements")
	ErrEmptyData     = errors.New("personal data must be set")

	// roles and rights
	ErrInvalidRole      = errors.New("invalid user role")
	ErrPermissionDenied = errors.New("permission denied")

	// status
	ErrInvalidStatus = errors.New("invalid user status")
	ErrInactiveUser  = errors.New("user is inactive")

	// logical
	ErrCannotDeleteSelf    = errors.New("user cannot delete themselves")
	ErrEmailAlreadyUsed    = errors.New("email is already in use")
	ErrManagerMustHaveTeam = errors.New("manager must have a team")
)
