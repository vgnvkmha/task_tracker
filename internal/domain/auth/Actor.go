package auth

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"

	"github.com/google/uuid"
)

type Actor struct {
	Id   uuid.UUID
	Role valueobjects.Role
}
