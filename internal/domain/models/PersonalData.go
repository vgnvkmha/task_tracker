package models

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"

	"github.com/google/uuid"
)

type PersonalData struct {
	Id       uuid.UUID
	Email    string
	Password string
	Role     valueobjects.Role
}
