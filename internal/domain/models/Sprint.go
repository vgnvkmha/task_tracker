package models

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"time"
)

type Sprint struct {
	ID        int
	Name      string
	Goal      string
	StartDate time.Time
	EndDate   time.Time
	Status    valueobjects.SprintStatus
	BoardID   int

	CreatedAt time.Time
	UpdatedAt time.Time
}
