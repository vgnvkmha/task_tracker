package user

import "github.com/google/uuid"

type UpdateRequest struct {
	Email    *string    `json:"email"`
	Password *string    `json:"password"`
	Role     *string    `json:"role"`
	TeamID   *uuid.UUID `json:"team_id"`
}
