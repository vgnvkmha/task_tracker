package models

import (
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"time"
)

type Task struct {
	id          int
	name        string
	description string
	status      valueobjects.Status
	board       Board
	created_at  time.Time
	due_to      time.Time
	assignee    User
}
