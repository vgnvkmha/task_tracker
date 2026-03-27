package models

import (
	"errors"
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"task_tracker/internal/domain/validation"
	"time"
)

type Task struct {
	Id          uint32
	Name        string
	Description string
	Status      valueobjects.Status
	BoardId     uint32
	CreatedAt   time.Time
	DueTo       time.Time
	AssigneeId  uint32
	ReporterId  uint32
	Sprint      Sprint
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
		Name:        name,
		Description: description,
		Status:      valueobjects.InProgress,
		BoardId:     boardId,
		CreatedAt:   time.Now(),
		DueTo:       dueTo,
		AssigneeId:  assigneeId,
		ReporterId:  reporterId,
	}, nil
}

func (t *Task) ChangeStatus(newStatus valueobjects.Status) error {

	if !newStatus.IsValid() {
		return errors.New("invalid status")
	}

	if !validation.IsValidStatusTransition(t.Status, newStatus) {
		return errors.New("invalid status transition")
	}

	t.Status = newStatus
	return nil
}
