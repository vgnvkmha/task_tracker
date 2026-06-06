package task

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskInput struct {
	Name        string
	Description string
	DueTo       time.Time
	ReporterID  *uuid.UUID
	AssigneeID  *uuid.UUID
	BoardID     *uuid.UUID
	SprintID    *uuid.UUID
}
