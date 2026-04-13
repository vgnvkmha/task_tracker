package user

import (
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID
	TeamId         uuid.UUID
	Email          string
	Password       string
	Role           valueobjects.Role
	PersonalDataId uuid.UUID
}
