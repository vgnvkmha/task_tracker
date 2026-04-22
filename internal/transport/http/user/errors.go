package user

import "errors"

var (
	// --- 400 Bad Request
	ErrInvalidRequest   = errors.New("invalid request body")
	ErrInvalidUserID    = errors.New("invalid user_id")
	ErrInvalidTeamID    = errors.New("invalid team_id")
	ErrInvalidBirthDate = errors.New("invalid birth_date")
	ErrInvalidRole      = errors.New("invalid role")

	// --- 401 Unauthorized
	ErrUnauthorized = errors.New("unauthorized")

	// --- 403 Forbidden
	ErrForbidden = errors.New("forbidden")

	// --- 404 Not Found
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")

	// --- 409 Conflict
	ErrUserAlreadyExists = errors.New("user already exists")

	// --- 500 Internal
	ErrInternal = errors.New("internal server error")
)
