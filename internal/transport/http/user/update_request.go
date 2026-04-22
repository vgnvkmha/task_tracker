package user

import (
	"task_tracker/internal/application/user"
	"time"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	UserID string `json:"user_id"`

	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
	TeamID    *string `json:"team_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Age       *uint8  `json:"age"`
	BirthDate *string `json:"birth_date"`
}

func (r UpdateRequest) ToServiceInput() (user.UpdateUserInput, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return user.UpdateUserInput{}, ErrInvalidUserID
	}

	var teamID *uuid.UUID
	if r.TeamID != nil {
		tid, err := uuid.Parse(*r.TeamID)
		if err != nil {
			return user.UpdateUserInput{}, ErrInvalidTeamID
		}
		teamID = &tid
	}

	var birthDate *time.Time
	if r.BirthDate != nil {
		t, err := time.Parse(time.RFC3339, *r.BirthDate)
		if err != nil {
			return user.UpdateUserInput{}, ErrInvalidBirthDate
		}
		birthDate = &t
	}

	return user.UpdateUserInput{
		UserID:    userID,
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
		TeamId:    teamID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
		BirthDate: birthDate,
	}, nil
}
