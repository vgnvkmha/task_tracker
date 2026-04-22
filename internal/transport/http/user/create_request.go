package user

import (
	"task_tracker/internal/application/user"
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

func (r CreateRequest) ToServiceInput() user.CreateUserInput {
	return user.CreateUserInput{
		Email:     r.Email,
		Password:  r.Password,
		Role:      &r.Role,
		TeamName:  r.TeamName,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
		BirthDate: r.BirthDate,
	}
}
