package task

import (
	domaintask "task_tracker/internal/domain/task"

	"github.com/google/uuid"
)

type TaskFilters struct {
	BoardID    *uuid.UUID
	AssigneeID *uuid.UUID
	ReporterID *uuid.UUID
	SprintID   *uuid.UUID
	Status     *domaintask.TaskStatus
}
