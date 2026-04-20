package user

import (
	"time"

	"github.com/google/uuid"
)

type UpdateUserInput struct {
	UserID uuid.UUID

	Email    *string
	Password *string
	Role     *string
	TeamName *string

	FirstName *string
	LastName  *string
	Age       *uint8
	BirthDate *time.Time
}
