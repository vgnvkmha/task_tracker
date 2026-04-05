package models

import (
	"errors"
	err "task_tracker/internal/domain/errors"
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

	err := newStatus.IsValid()
	if err != nil {
		return err
	}

	if !validation.IsValidStatusTransition(t.Status, newStatus) {
		return errors.New("invalid status transition")
	}

	t.Status = newStatus
	return nil
}

func (t *Task) ChangeBoard(newBoardId uint32) error {
	if t.BoardId == newBoardId {
		return errors.New("Same Board")
	}
	if t.Status.IsImmutable() != nil || t.Sprint.Status.IsImmutable() != nil {
		return err.ErrInvalidStatus
	}
	t.BoardId = newBoardId
	return nil
}

func (t *Task) ChangeReporter(newReporterId uint32) error {
	if t.ReporterId == newReporterId {
		return errors.New("Same Reporter")
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrInvalidRights
	}
	t.ReporterId = newReporterId
	return nil
}

func (t *Task) ChangeAssignee(newAssigneeId uint32) error {
	if t.AssigneeId == newAssigneeId {
		return errors.New("Same Assignee")
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrInvalidRights
	}
	t.AssigneeId = newAssigneeId
	return nil
}

func (t *Task) ChangeSprint(newSprintId uint32) error {
	if t.Sprint.ID == newSprintId {
		return errors.New("Same Sprint")
	}
	if t.Status.IsImmutable() != nil || t.Sprint.Status.IsImmutable() != nil {
		return err.ErrInvalidRights
	}
	t.Sprint.ID = newSprintId
	return nil
}
