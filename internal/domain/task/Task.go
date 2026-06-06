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
	BoardId     *uuid.UUID
	SprintId    *uuid.UUID
}

func New(
	id uuid.UUID,
	name string,
	description string,
	boardID *uuid.UUID,
	dueTo time.Time,
	assigneeID *uuid.UUID,
	reporterID uuid.UUID,
	sprintId *uuid.UUID,
) (Task, error) {

	if strings.TrimSpace(name) == "" {
		return Task{}, ErrTaskName
	}
	if reporterID == uuid.Nil {
		return Task{}, ErrTaskUser
	}

	if !dueTo.IsZero() && time.Now().After(dueTo) {
		return Task{}, ErrInvalidTime
	}

	now := time.Now()
	task := Task{
		Id:          id,
		Name:        strings.TrimSpace(name),
		Description: description,
		Status:      Todo,
		BoardId:     boardID,
		CreatedAt:   now,
		DueTo:       dueTo,
		UpdatedAt:   now,
		AssigneeId:  assigneeID,
		ReporterId:  reporterID,
		SprintId:    sprintId,
	}

	return task, nil
}

func (t *Task) Update(
	name *string,
	description *string,
	status *TaskStatus,
	dueTo *time.Time,
	reporterID *uuid.UUID,
	assigneeID *uuid.UUID,
	boardID *uuid.UUID,
	sprintID *uuid.UUID,
) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return ErrTaskName
		}
		t.Name = strings.TrimSpace(*name)
	}
	if description != nil {
		t.Description = *description
	}
	if status != nil {
		if err := status.IsValid(); err != nil {
			return err
		}
		t.Status = *status
	}
	if dueTo != nil {
		if !dueTo.IsZero() && time.Now().After(*dueTo) {
			return ErrInvalidTime
		}
		t.DueTo = *dueTo
	}
	if reporterID != nil {
		if *reporterID == uuid.Nil {
			return ErrTaskUser
		}
		t.ReporterId = *reporterID
	}
	if assigneeID != nil {
		t.AssigneeId = assigneeID
	}
	if boardID != nil {
		t.BoardId = boardID
	}
	if sprintID != nil {
		t.SprintId = sprintID
	}
	t.UpdatedAt = time.Now()
	return nil
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

func (t *Task) ChangeBoard(newBoardId *uuid.UUID) error {
	if t.BoardId != nil && newBoardId != nil && *t.BoardId == *newBoardId {
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
