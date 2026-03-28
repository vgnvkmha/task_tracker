package service

import (
	"context"
	"errors"
	"task_tracker/internal/domain/models"
	valueobjects "task_tracker/internal/domain/models/value_objects" //TODO: 2 vo imports
	vo "task_tracker/internal/domain/models/value_objects"
	"task_tracker/internal/repo"
	dto "task_tracker/internal/transport/task"

	"go.uber.org/zap"
)

type TaskService interface {
	Create(ctx context.Context, task dto.TaskRequest) (models.Task, error)

	GetActive(ctx context.Context) ([]models.Task, error)

	ChangeStatus(ctx context.Context, id uint32, status vo.Status) (valueobjects.Status, error)
	ChangeBoard(ctx context.Context, boardId uint32) (models.Board, error)
	ChangeAssign(ctx context.Context, assignId uint32) (models.User, error)
	ChangeReporter(ctx context.Context, reporterId uint32) (models.User, error)
	ChangeSprint(ctx context.Context, sprint models.Sprint) (models.Sprint, error)
}

type service struct {
	repo   repo.TaskRepo
	logger *zap.SugaredLogger
}

func New(repo repo.TaskRepo, logger *zap.SugaredLogger) *service {
	return &service{
		repo: repo,
		logger: logger.With(
			"module", "task",
			"layer", "service",
		),
	}
}

func (s *service) Create(ctx context.Context, task dto.TaskRequest) (models.Task, error) {
	model, err := models.NewTask(
		task.Name,
		task.Description,
		task.BoardID,
		task.AssigneeID,
		task.ReporetID,
		task.DueTo,
	)
	if err != nil {
		s.logger.Infow("Mapping error",
			"operation", "Create",
		)
		return models.Task{}, err
	}
	return s.repo.Create(ctx, model)
}

func (s *service) ChangeStatus(
	ctx context.Context,
	id uint32,
	status valueobjects.Status,
) error {

	const op = "task.ChangeStatus"

	if !status.IsValid() {
		err := errors.New("invalid status")

		s.logger.Infow("invalid status",
			"operation", op,
			"status", status,
		)

		return err
	}

	task, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Infow("task not found",
			"operation", op,
			"id", id,
		)
		return err
	}

	err = task.ChangeStatus(status)
	if err != nil {
		s.logger.Infow("invalid status transition",
			"operation", op,
			"from", task.Status,
			"to", status,
		)
		return err
	}

	return s.repo.Update(ctx, task)
}

func (s *service) ChangeBoard(ctx context.Context, boardId uint32) (models.Board, error) {
	return models.Board{}, nil
}

func (s *service) ChangeAssign(ctx context.Context, assignId uint32) (models.User, error) {
	return models.User{}, nil
}

func (s *service) ChangeReporter(ctx context.Context, reporterId uint32) (models.User, error) {
	return models.User{}, nil
}

func (s *service) ChangeSprint(ctx context.Context, sprint models.Sprint) (models.Sprint, error) {
	return models.Sprint{}, nil
}
