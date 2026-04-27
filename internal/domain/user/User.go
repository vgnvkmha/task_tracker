package user

import (
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	TeamID         *uuid.UUID
	Email          valueobjects.Email //TODO: make unique in DB
	Password       valueobjects.Password
	Role           valueobjects.Role
	PersonalDataID uuid.UUID
	IsActive       bool
}

func New(teamId, personalDataId uuid.UUID, emailRaw, passwordRaw, roleRaw string) (*User, error) {

	if personalDataId == uuid.Nil {
		return nil, ErrEmptyData
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
		ID:             uuid.New(),
		TeamID:         &teamId,
		Email:          email,
		Password:       password,
		Role:           role,
		PersonalDataID: personalDataId,
		IsActive:       true,
	}

	return user, nil
}
