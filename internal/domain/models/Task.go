package models

import (
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
