package models

import (
	"errors"
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"time"
)

type Task struct {
	id          uint32
	name        string
	description string
	status      valueobjects.Status
	boardId     uint32
	created_at  time.Time
	due_to      time.Time
	assigneeId  uint32
	reporterId  uint32
}

func NewTask(
	name string,
	description string,
	boardId uint32,
	assigneeId uint32,
	reporterId uint32,
	dueTo time.Time,
) (Task, error) {

	if name == "" {
		return Task{}, errors.New("name is required")
	}

	return Task{
		name:        name,
		description: description,
		status:      valueobjects.InProgress,
		boardId:     boardId,
		created_at:  time.Now(),
		due_to:      dueTo,
		assigneeId:  assigneeId,
		reporterId:  reporterId,
	}, nil
}
