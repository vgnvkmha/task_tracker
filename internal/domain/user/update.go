package user

import (
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

func (u *User) Update(teamID *uuid.UUID, emailRaw, passwordRaw, roleRaw *string) error {
	if teamID != nil {
		u.TeamID = teamID
	}

	if emailRaw != nil {
		email, err := valueobjects.NewEmail(*emailRaw)
		if err != nil {
			return err
		}
		u.Email = email
	}

	if passwordRaw != nil {
		password, err := valueobjects.NewPassword(*passwordRaw)
		if err != nil {
			return err
		}
		u.Password = password
	}

	if roleRaw != nil {
		if !valueobjects.IsValidRole(*roleRaw) {
			return ErrInvalidRole
		}
		u.Role = valueobjects.Role(*roleRaw)
	}

	if u.TeamID == nil && u.Role.IsManagerRole() {
		return ErrManagerMustHaveTeam
	}

	return nil
}
