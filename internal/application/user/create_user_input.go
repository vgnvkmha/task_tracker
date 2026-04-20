package user

import (
	"time"
)

type CreateUserInput struct {
	Email     string
	Password  string
	Role      *string
	TeamName  *string
	FirstName string
	LastName  string
	Age       *uint8
	BirthDate *time.Time
}
