package models

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"time"
)

type Sprint struct {
	ID        uint32
	Name      string
	StartDate time.Time
	EndDate   time.Time
	Status    valueobjects.SprintStatus
	BoardID   uint32
	TasksId   uint32
}
