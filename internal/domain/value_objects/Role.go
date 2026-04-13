package valueobjects

import errors_task "task_tracker/internal/domain/errors"

type Role string

const (
	Admin   Role = "admin"
	Captain Role = "captain"
	User    Role = "user"
	Guest   Role = "guest"
)

func (r Role) IsValid() (Role, error) {

	switch r {
	case Admin, Captain, User, Guest:
		return r, nil
	default:
		return "", errors_task.ErrInvalidRole
	}
}

func (r Role) IsManagerRole() bool {
	switch r {
	case Admin, Captain:
		return true
	default:
		return false
	}
}
