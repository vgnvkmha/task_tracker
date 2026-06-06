package task

import (
	"time"

	domaintask "task_tracker/internal/domain/task"

	"github.com/google/uuid"
)

type UpdateTaskInput struct {
	Name        *string
	Description *string
	Status      *domaintask.TaskStatus
	DueTo       *time.Time
	ReporterID  *uuid.UUID
	AssigneeID  *uuid.UUID
	BoardID     *uuid.UUID
	SprintID    *uuid.UUID
}
