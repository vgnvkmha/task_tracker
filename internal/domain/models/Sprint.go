package models

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"time"

	uuid "github.com/google/uuid"
)

type Sprint struct {
	Id        uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	Status    valueobjects.SprintStatus
	BoardId   uuid.UUID
	TasksIds  uuid.UUID
}
