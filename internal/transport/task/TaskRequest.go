package dto

import (
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	BoardID     uuid.UUID  `json:"board_id"`
	DueTo       time.Time  `json:"due_to"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	ReporetID   uuid.UUID  `json:"reporter_id"`
	SprintId    *uuid.UUID `json:"sprint_id"`
}
