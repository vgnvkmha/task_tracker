package user

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	ActorId   uuid.UUID
	ActorRole string
	Email     string
	Password  string
	Role      *string
	TeamName  *string
	FirstName string
	LastName  string
	Age       *uint8
	BirthDate *time.Time
}
