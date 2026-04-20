package user

import "github.com/google/uuid"

type Response struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	TeamID         *uuid.UUID `json:"team_id"`
	PersonalDataID uuid.UUID  `json:"personal_data_id"`
}
