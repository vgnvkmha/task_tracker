package service

import (
	"context"
	taskErr "task_tracker/internal/domain/errors"
	"task_tracker/internal/domain/models"
	valueobjects "task_tracker/internal/domain/models/value_objects" //TODO: 2 vo imports
	"task_tracker/internal/domain/validation"
	"task_tracker/internal/repo"
	dto "task_tracker/internal/transport/task"

	"go.uber.org/zap"
)

const msg = "Task Service Error"

type TaskService interface {
	Create(ctx context.Context, task dto.TaskRequest) (models.Task, error)

	GetActiveTasksByTeam(ctx context.Context, id uint32) ([]models.Task, error)
	GetTeamById(ctx context.Context, id uint32) (models.Team, error)

	ChangeStatus(ctx context.Context, taskId uint32, input string) (valueobjects.Status, error)
	ChangeBoard(ctx context.Context, taskId, boardId uint32) (models.Board, error)
	ChangeAssign(ctx context.Context, taskId, assignId uint32) (models.User, error)
	ChangeReporter(ctx context.Context, taskId, reporterId uint32) (models.User, error)
	ChangeSprint(ctx context.Context, taskId, sprintId uint32) (models.Sprint, error)
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

func (s *service) GetActiveTasksByTeam(ctx context.Context, id uint32) ([]models.Task, error) {
	tasks, err := s.repo.GetActiveByTeamId(ctx, id)
	if err != nil {
		s.logger.Infow("Getting Active Tasks Failure",
			"team_id", id,
			"error", err,
		)
		return []models.Task{}, err
	}
	return tasks, nil
}

func (s *service) GetTeamById(ctx context.Context, id uint32) (models.Team, error) {
	team, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		s.logger.Infow(msg,
			"id", id,
			"error", err,
		)
		return models.Team{}, err
	}
	return team, nil
}

// TODO: change method so it looks like others
func (s *service) ChangeStatus(
	ctx context.Context,
	id uint32,
	input string,
) error {
	const op = "Change Status"
	logError := func(err error) error {
		s.logger.Infow(msg,
			"operation", op,
			"error", err,
		)
		return err
	}
	newStatus, err := validation.ParseStatus(input)
	if err != nil {
		return logError(err)
	}

	if newStatus.IsValid() != nil || newStatus.IsImmutable() != nil {
		return logError(taskErr.InvalidStatus)
	}

	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		s.logger.Infow(msg,
			"operation", op,
			"id", id,
			"error", err,
		)
		return err
	}

	err = task.ChangeStatus(newStatus)
	if err != nil {
		s.logger.Infow(msg,
			"operation", op,
			"from", task.Status,
			"to", input,
			"error", err,
		)
		return err
	}

	return s.repo.Update(ctx, task)
}

func (s *service) ChangeBoard(ctx context.Context, taskId, boardId uint32) (models.Board, error) {
	const op = "ChangeBoard"
	logError := func(err error) error {
		s.logger.Infow(msg,
			"operation", op,
			"error", err,
		)
		return err
	}

	board, err := s.repo.GetBoard(ctx, boardId)
	if err != nil {
		return models.Board{}, logError(err)
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.Board{}, logError(err)
	}

	if err := task.ChangeBoard(boardId); err != nil {
		return models.Board{}, logError(err)
	}

	return board, nil
}

func (s *service) ChangeAssign(ctx context.Context, taskId, assignId uint32) (models.User, error) {
	const op = "Change Assignee"
	logError := func(err error) error {
		s.logger.Infow(msg,
			"operation", op,
			"error", err,
		)
		return err
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err)
	}

	assignee, err := s.repo.GetUser(ctx, assignId)
	if err != nil {
		return models.User{}, logError(err)
	}

	err = task.ChangeAssignee(assignId)
	if err != nil {
		return models.User{}, logError(err)
	}

	return assignee, nil
}

func (s *service) ChangeReporter(ctx context.Context, taskId, reporterId uint32) (models.User, error) {
	const op = "Change Reporter"
	logError := func(err error) error {
		s.logger.Infow(msg,
			"operation", op,
			"error", err,
		)
		return err
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err)
	}

	reporter, err := s.repo.GetUser(ctx, reporterId)
	if err != nil {
		return models.User{}, logError(err)
	}

	err = task.ChangeReporter(reporterId)
	if err != nil {
		return models.User{}, logError(err)
	}

	return reporter, nil
}

func (s *service) ChangeSprint(ctx context.Context, taskId, sprintId uint32) (models.Sprint, error) {
	const op = "Change Sprint"
	logError := func(err error) error {
		s.logger.Infow(msg,
			"operation", op,
			"error", err,
		)
		return err
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.Sprint{}, logError(err)
	}

	_, err = s.repo.GetUser(ctx, sprintId)
	if err != nil {
		return models.Sprint{}, logError(err)
	}

	err = task.ChangeSprint(sprintId)
	if err != nil {
		return models.Sprint{}, logError(err)
	}

	return task.Sprint, nil
}
