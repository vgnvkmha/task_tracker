package task

import (
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID
	Name        string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	DueTo       time.Time
	UpdatedAt   time.Time
	ReporterId  uuid.UUID
	AssigneeId  *uuid.UUID
	BoardId     uuid.UUID
	SprintId    *uuid.UUID
}

func New(
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
		return Task{}, ErrTaskName
	}
	if boardID == uuid.Nil {
		return Task{}, ErrTaskBoard
	}
	if reporterID == uuid.Nil {
		return Task{}, ErrTaskUser
	}

	if time.Now().After(dueTo) {
		return Task{}, ErrInvalidTime
	}

	task := Task{
		Id:          id,
		Name:        name,
		Description: description,
		Status:      Todo,
		BoardId:     boardID,
		CreatedAt:   time.Now(),
		DueTo:       dueTo,
		AssigneeId:  assigneeID,
		ReporterId:  reporterID,
		SprintId:    sprintId,
	}

	return task, nil
}

func (t *Task) ChangeStatus(newStatus TaskStatus) error {
	err := newStatus.IsValid()
	if err != nil {
		return err
	}

	if err = IsValidStatusTransition(t.Status, newStatus); err != nil {
		return err
	}

	t.Status = newStatus
	return nil
}

func (t *Task) ChangeBoard(newBoardId uuid.UUID) error {
	if t.BoardId == newBoardId {
		return ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return ErrInvalidRights
	}
	t.BoardId = newBoardId
	return nil
}

func (t *Task) ChangeReporter(newReporterId uuid.UUID) error {
	if t.ReporterId == newReporterId {
		return ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return ErrInvalidRights
	}
	t.ReporterId = newReporterId
	return nil
}

func (t *Task) ChangeAssignee(newAssigneeId *uuid.UUID) error {
	if t.AssigneeId == newAssigneeId {
		return ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return ErrInvalidRights
	}
	t.AssigneeId = newAssigneeId
	return nil
}

func (t *Task) ChangeSprint(newSprintId *uuid.UUID) error {
	if t.SprintId != nil && newSprintId != nil && *t.SprintId == *newSprintId {
		return ErrSameChange
	}
	if t.Status.IsImmutable() != nil {
		return ErrImmutableTask
	}

	t.SprintId = newSprintId
	return nil
}
