package user

import (
	"time"
)

type CreateRequest struct {
	Email     string     `json:"email" binding:"required"`
	Password  string     `json:"password" binding:"required"`
	Role      string     `json:"role" binding:"required"`
	TeamName  *string    `json:"team_name"`
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Age       *uint8     `json:"age"`
	BirthDate *time.Time `json:"birth_date"`
}
