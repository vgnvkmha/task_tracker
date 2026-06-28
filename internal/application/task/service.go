package task

import (
	"context"
	"errors"

	"task_tracker/internal/application/common"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	domaintask "task_tracker/internal/domain/task"
	"task_tracker/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	module = "task"
	layer  = "service"
)

type Task = domaintask.Task

type TaskService interface {
	Create(ctx context.Context, actor auth.Actor, input CreateTaskInput) (*Task, error)
	GetByID(ctx context.Context, actor auth.Actor, id uuid.UUID) (*Task, error)
	FindMany(ctx context.Context, actor auth.Actor, filters TaskFilters) ([]*Task, error)
	FindByBoardID(ctx context.Context, actor auth.Actor, boardID uuid.UUID) ([]*Task, error)
	FindByAssigneeID(ctx context.Context, actor auth.Actor, assigneeID uuid.UUID) ([]*Task, error)
	Update(ctx context.Context, actor auth.Actor, id uuid.UUID, input UpdateTaskInput) (*Task, error)
	Delete(ctx context.Context, actor auth.Actor, id uuid.UUID) error
}

type service struct {
	taskRepo  repo.TaskRepo
	boardRepo repo.BoardRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(taskRepo repo.TaskRepo, boardRepo repo.BoardRepo, logger *zap.SugaredLogger, transaction common.TxManager) TaskService {
	return &service{
		taskRepo:    taskRepo,
		boardRepo:   boardRepo,
		logger:      logger,
		transaction: transaction,
	}
}

func (s *service) Create(ctx context.Context, actor auth.Actor, input CreateTaskInput) (*Task, error) {
	var result *Task

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if err := s.validateCreateBoardAccess(ctx, actor, input.BoardID); err != nil {
			return err
		}

		reporterID := actor.ID
		if input.ReporterID != nil {
			if !actor.Role.IsManagerRole() && *input.ReporterID != actor.ID {
				return ErrPermissionDenied
			}
			reporterID = *input.ReporterID
		}

		task, err := domaintask.New(
			uuid.New(),
			input.Name,
			input.Description,
			input.BoardID,
			input.DueTo,
			input.AssigneeID,
			reporterID,
			input.SprintID,
		)
		if err != nil {
			return mapDomainError(err)
		}

		created, err := s.taskRepo.Create(ctx, task)
		if err != nil {
			return mapRepoError(err, ErrCreateTaskFailed)
		}
		result = created
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "create",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"name", input.Name,
		)
	}

	logSuccess(s.logger,
		"operation", "create",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"task_id", result.Id,
	)
	return result, nil
}

func (s *service) validateCreateBoardAccess(ctx context.Context, actor auth.Actor, boardID *uuid.UUID) error {
	if boardID == nil || *boardID == uuid.Nil {
		return ErrBoardRequired
	}
	if actor.TeamID == nil || *actor.TeamID == uuid.Nil {
		return ErrActorTeamRequired
	}

	board, err := s.boardRepo.GetByID(ctx, *boardID)
	if err != nil {
		return mapRepoError(err, ErrBoardNotFound)
	}
	if board == nil {
		return ErrBoardNotFound
	}
	if board.TeamId != *actor.TeamID {
		return ErrBoardTeamMismatch
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, actor auth.Actor, id uuid.UUID) (*Task, error) {
	result, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		err = mapRepoError(err, ErrTaskNotFound)
		return nil, logError(err, s.logger,
			"operation", "get_by_id",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"task_id", id,
		)
	}
	if result == nil {
		return nil, ErrTaskNotFound
	}

	return result, nil
}

func (s *service) FindMany(ctx context.Context, actor auth.Actor, filters TaskFilters) ([]*Task, error) {

	if filters.Status != nil {
		if err := filters.Status.IsValid(); err != nil {
			return nil, logError(ErrInvalidStatus, s.logger,
				"operation", "find_many",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
			)
		}
	}
	result, err := s.taskRepo.FindMany(ctx, toRepoFilters(filters))
	if err != nil {
		mappedErr := mapRepoError(err, ErrTaskNotFound)
		return nil, logError(mappedErr, s.logger,
			"operation", "find_many",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
		)
	}

	return result, nil
}

func (s *service) FindByBoardID(ctx context.Context, actor auth.Actor, boardID uuid.UUID) ([]*Task, error) {
	return s.FindMany(ctx, actor, TaskFilters{BoardID: &boardID})
}

func (s *service) FindByAssigneeID(ctx context.Context, actor auth.Actor, assigneeID uuid.UUID) ([]*Task, error) {
	return s.FindMany(ctx, actor, TaskFilters{AssigneeID: &assigneeID})
}

func (s *service) Update(ctx context.Context, actor auth.Actor, id uuid.UUID, input UpdateTaskInput) (*Task, error) {
	var result *Task

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		existing, err := s.taskRepo.GetByID(ctx, id)
		if err != nil {
			return mapRepoError(err, ErrTaskNotFound)
		}
		if existing == nil {
			return ErrTaskNotFound
		}
		if !canModify(actor, existing) {
			return ErrPermissionDenied
		}

		if input.ReporterID != nil && !actor.Role.IsManagerRole() && *input.ReporterID != actor.ID {
			return ErrPermissionDenied
		}

		if err := existing.Update(
			input.Name,
			input.Description,
			input.Status,
			input.DueTo,
			input.ReporterID,
			input.AssigneeID,
			input.BoardID,
			input.SprintID,
		); err != nil {
			return mapDomainError(err)
		}

		updated, err := s.taskRepo.Update(ctx, id, *existing)
		if err != nil {
			return mapRepoError(err, ErrUpdateTaskFailed)
		}
		if updated == nil {
			return ErrTaskNotFound
		}

		result = updated
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "update",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"task_id", id,
		)
	}

	logSuccess(s.logger,
		"operation", "update",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"task_id", result.Id,
	)
	return result, nil
}

func (s *service) Delete(ctx context.Context, actor auth.Actor, id uuid.UUID) error {
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		existing, err := s.taskRepo.GetByID(ctx, id)
		if err != nil {
			return mapRepoError(err, ErrTaskNotFound)
		}
		if existing == nil {
			return ErrTaskNotFound
		}
		if !canModify(actor, existing) {
			return ErrPermissionDenied
		}

		if err := s.taskRepo.Delete(ctx, id); err != nil {
			return mapRepoError(err, ErrDeleteTaskFailed)
		}
		return nil
	})

	if err != nil {
		return logError(err, s.logger,
			"operation", "delete",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"task_id", id,
		)
	}

	logSuccess(s.logger,
		"operation", "delete",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"task_id", id,
	)
	return nil
}

func canModify(actor auth.Actor, task *Task) bool {
	if actor.Role.IsManagerRole() {
		return true
	}
	if actor.ID == task.ReporterId {
		return true
	}
	return task.AssigneeId != nil && actor.ID == *task.AssigneeId
}

func toRepoFilters(filters TaskFilters) repo.TaskFilters {
	return repo.TaskFilters{
		BoardID:    filters.BoardID,
		AssigneeID: filters.AssigneeID,
		ReporterID: filters.ReporterID,
		SprintID:   filters.SprintID,
		Status:     filters.Status,
	}
}

func mapRepoError(err error, fallback error) error {
	switch {
	case errors.Is(err, common_errors.ErrNotFound):
		return ErrTaskNotFound
	case errors.Is(err, common_errors.ErrInvalidArgument):
		return ErrInvalidInput
	case errors.Is(err, common_errors.ErrConflict):
		return ErrInvalidInput
	default:
		return fallback
	}
}

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domaintask.ErrTaskName),
		errors.Is(err, domaintask.ErrTaskBoard),
		errors.Is(err, domaintask.ErrTaskUser):
		return ErrInvalidInput
	case errors.Is(err, domaintask.ErrInvalidTime):
		return ErrInvalidDueTo
	case errors.Is(err, domaintask.ErrInvalidStatus),
		errors.Is(err, domaintask.ErrInvalidRole):
		return ErrInvalidStatus
	case errors.Is(err, domaintask.ErrInvalidStatusTransition):
		return ErrInvalidTransition
	case errors.Is(err, domaintask.ErrInvalidRights),
		errors.Is(err, domaintask.ErrImmutableTask):
		return ErrPermissionDenied
	default:
		return err
	}
}
