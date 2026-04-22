package user

import (
	"task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type Response struct {
	ID             uuid.UUID          `json:"id"`
	Email          valueobjects.Email `json:"email"`
	Role           valueobjects.Role  `json:"role"`
	TeamID         *uuid.UUID         `json:"team_id"`
	PersonalDataID uuid.UUID          `json:"personal_data_id"`
}

func FromService(user user.User) Response {
	return Response{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		TeamID:         user.TeamID,
		PersonalDataID: user.PersonalDataID,
	}
}
