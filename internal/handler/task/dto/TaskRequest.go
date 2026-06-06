package dto

import (
	taskApplication "task_tracker/internal/application/task"
	domaintask "task_tracker/internal/domain/task"
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	BoardID     *uuid.UUID `json:"board_id"`
	DueTo       *time.Time `json:"due_to"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	ReporterID  *uuid.UUID `json:"reporter_id"`
	SprintID    *uuid.UUID `json:"sprint_id"`
}

func (r TaskRequest) ToServiceInput() taskApplication.CreateTaskInput {
	var dueTo time.Time
	if r.DueTo != nil {
		dueTo = *r.DueTo
	}

	return taskApplication.CreateTaskInput{
		Name:        r.Name,
		Description: r.Description,
		DueTo:       dueTo,
		ReporterID:  r.ReporterID,
		AssigneeID:  r.AssigneeID,
		BoardID:     r.BoardID,
		SprintID:    r.SprintID,
	}
}

type UpdateTaskRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	BoardID     *uuid.UUID `json:"board_id"`
	DueTo       *time.Time `json:"due_to"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	ReporterID  *uuid.UUID `json:"reporter_id"`
	SprintID    *uuid.UUID `json:"sprint_id"`
}

func (r UpdateTaskRequest) ToServiceInput() taskApplication.UpdateTaskInput {
	var status *domaintask.TaskStatus
	if r.Status != nil {
		value := domaintask.TaskStatus(*r.Status)
		status = &value
	}

	return taskApplication.UpdateTaskInput{
		Name:        r.Name,
		Description: r.Description,
		Status:      status,
		DueTo:       r.DueTo,
		ReporterID:  r.ReporterID,
		AssigneeID:  r.AssigneeID,
		BoardID:     r.BoardID,
		SprintID:    r.SprintID,
	}
}
