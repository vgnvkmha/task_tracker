package models

import (
	"strings"
	err "task_tracker/internal/domain/errors"
	task_errors "task_tracker/internal/domain/errors"
	valueobjects "task_tracker/internal/domain/models/value_objects"
	"task_tracker/internal/domain/validation"
	"time"

	uuid "github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID
	Name        string
	Description string
	Status      valueobjects.Status
	BoardId     uuid.UUID
	CreatedAt   time.Time
	DueTo       time.Time
	AssigneeId  *uuid.UUID
	ReporterId  uuid.UUID
	SprintId    *uuid.UUID
}

func NewTask(
	id uuid.UUID,
	name string,
	description string,
	boardID uuid.UUID,
	dueTo time.Time,
	assigneeID *uuid.UUID,
	reporterID uuid.UUID,
	sprintId *uuid.UUID,
) (Task, error) {

	if strings.TrimSpace(name) == "" {
		return Task{}, task_errors.ErrTaskName
	}
	if boardID == uuid.Nil {
		return Task{}, task_errors.ErrTaskBoard
	}
	if reporterID == uuid.Nil {
		return Task{}, task_errors.ErrTaskUser
	}

	if time.Now().After(dueTo) {
		return Task{}, task_errors.ErrInvalidTime
	}

	task := Task{
		Id:          id,
		Name:        name,
		Description: description,
		Status:      valueobjects.Todo,
		BoardId:     boardID,
		CreatedAt:   time.Now(),
		DueTo:       dueTo,
		AssigneeId:  assigneeID,
		ReporterId:  reporterID,
		SprintId:    sprintId,
	}

	return task, nil
}

func (t *Task) ChangeStatus(newStatus valueobjects.Status) error {
	err := newStatus.IsValid()
	if err != nil {
		return err
	}

	if err = validation.IsValidStatusTransition(t.Status, newStatus); err != nil {
		return err
	}

	t.Status = newStatus
	return nil
}

func (t *Task) ChangeBoard(newBoardId uuid.UUID) error {
	if t.BoardId == newBoardId {
		return err.ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrInvalidStatus
	}
	t.BoardId = newBoardId
	return nil
}

func (t *Task) ChangeReporter(newReporterId uuid.UUID) error {
	if t.ReporterId == newReporterId {
		return err.ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrInvalidRights
	}
	t.ReporterId = newReporterId
	return nil
}

func (t *Task) ChangeAssignee(newAssigneeId *uuid.UUID) error {
	if t.AssigneeId == newAssigneeId {
		return err.ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrInvalidRights
	}
	t.AssigneeId = newAssigneeId
	return nil
}

func (t *Task) ChangeSprint(newSprintId *uuid.UUID) error {
	if t.SprintId != nil && newSprintId != nil && *t.SprintId == *newSprintId {
		return err.ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return err.ErrImmutableTask
	}

	t.SprintId = newSprintId
	return nil
}
