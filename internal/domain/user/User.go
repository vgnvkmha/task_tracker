package user

import (
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID
	TeamId         *uuid.UUID
	Email          valueobjects.Email
	Password       valueobjects.Password
	Role           valueobjects.Role
	PersonalDataId uuid.UUID
}

func New(teamId, personalDataId uuid.UUID, emailRaw, passwordRaw, roleRaw string) (*User, error) {

	if personalDataId == uuid.Nil {
		return nil, ErrEmptyDataId
	}

	email, err := valueobjects.NewEmail(emailRaw)
	if err != nil {
		return nil, err
	}

	password, err := valueobjects.NewPassword(passwordRaw)
	if err != nil {
		return nil, err
	}

	if !valueobjects.IsValidRole(roleRaw) {
		return nil, ErrInvalidRole
	}

	role := valueobjects.Role(roleRaw)
	if teamId == uuid.Nil && role.IsManagerRole() {
		return nil, ErrManagerMustHaveTeam
	}

	user := &User{
		Id:             uuid.New(),
		TeamId:         &teamId,
		Email:          email,
		Password:       password,
		Role:           role,
		PersonalDataId: personalDataId,
	}

	return user, nil
}
