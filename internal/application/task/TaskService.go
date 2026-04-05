package task_service

import (
	"context"
	errors_task "task_tracker/internal/domain/errors"
	"task_tracker/internal/domain/models"
	vo "task_tracker/internal/domain/models/value_objects"
	"task_tracker/internal/domain/validation"
	"task_tracker/internal/repo"
	dto "task_tracker/internal/transport/task"

	"go.uber.org/zap"
)

const (
	module = "task"
	layer  = "service"
)

type TaskService interface {
	Create(ctx context.Context, userId uint32, task dto.TaskRequest) (models.Task, error)

	GetActiveTasksByTeam(ctx context.Context, id uint32) ([]models.Task, error)
	GetTeamById(ctx context.Context, id uint32) (models.Team, error)

	ChangeStatus(ctx context.Context, userId, taskId uint32, newStatus string) (vo.Status, error)
	ChangeBoard(ctx context.Context, userId, taskId, newBoardId uint32) (models.Board, error)
	ChangeAssign(ctx context.Context, userId, taskId, newAssignId uint32) (models.User, error)
	ChangeReporter(ctx context.Context, userId, taskId, newreporterId uint32) (models.User, error)
	ChangeSprint(ctx context.Context, userId, taskId, newSprintId uint32) (models.Sprint, error)
}

type service struct {
	repo   repo.TaskRepo
	logger *zap.SugaredLogger
}

func New(repo repo.TaskRepo, logger *zap.SugaredLogger) *service {
	return &service{
		repo: repo,
		logger: logger.With(
			"module", module,
			"layer", layer,
		),
	}
}

func (s *service) Create(ctx context.Context, userId uint32, task dto.TaskRequest) (models.Task, error) {
	const op = "Create Task"

	loggingFields := []any{
		"operation", op,
		"user_id", userId,
		"user_role", "undefined",
	}

	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return models.Task{}, logError(err, s.logger, loggingFields...)
	}
	if !user.Data.Role.CanModifyTaskInRestrictedState() {
		return models.Task{}, logError(errors_task.ErrInvalidRights, s.logger, loggingFields...)
	}

	//TODO: add validation in task package
	model, err := models.NewTask(
		task.Name,
		task.Description,
		task.BoardID,
		task.AssigneeID,
		task.ReporetID,
		task.DueTo,
	)
	if err != nil {
		return models.Task{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return s.repo.Create(ctx, model)
}

func (s *service) GetActiveTasksByTeam(ctx context.Context, teamId uint32) ([]models.Task, error) {
	const op = "Get Active Task By Team ID"

	loggingFields := []any{
		"operation", op,
		"team_id", teamId,
		"user_role", "undefined",
	}

	tasks, err := s.repo.GetActiveByTeamId(ctx, teamId)
	if err != nil {
		return []models.Task{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return tasks, nil
}

func (s *service) GetTeamById(ctx context.Context, teamId uint32) (models.Team, error) {
	const op = "Get Team by ID"

	loggingFields := []any{
		"operation", op,
		"team_id", teamId,
		"user_role", "undefined",
	}
	team, err := s.repo.GetTeam(ctx, teamId)
	if err != nil {
		return models.Team{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return team, nil
}

func (s *service) ChangeStatus(
	ctx context.Context,
	userId,
	taskId uint32,
	inputStatus string,
) error {
	const op = "Change Status"

	loggingFields := []any{
		"operation", op,
		"user_id", userId,
		"task_id", taskId,
		"user_role", "undefined",
		"new_status", inputStatus,
	}
	newStatus, err := validation.ParseStatus(inputStatus)
	if err != nil {
		return logError(err, s.logger, loggingFields...)
	}

	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return logError(err, s.logger, loggingFields...)
	}
	role, err := user.Data.Role.IsValid()

	//changing role for logging as soon as we get it
	loggingFields[7] = role
	if err != nil {
		return logError(err, s.logger, loggingFields...)
	}

	if !role.CanModifyTaskInRestrictedState() || newStatus.IsValid() != nil || newStatus.IsImmutable() != nil {
		return logError(err, s.logger, loggingFields...)
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeStatus(newStatus)
	if err != nil {
		return logError(err, s.logger, loggingFields...)
	}

	if err = s.repo.Update(ctx, task); err != nil {
		return logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return nil
}

func (s *service) ChangeBoard(ctx context.Context, taskId, boardId uint32) (models.Board, error) {
	const op = "Change Task Board"

	loggingFields := []any{
		"operation", op,
		"board_id", boardId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	board, err := s.repo.GetBoard(ctx, boardId)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	if err := task.ChangeBoard(boardId); err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return board, nil
}

func (s *service) ChangeAssign(ctx context.Context, taskId, assignId uint32) (models.User, error) {
	const op = "Change Task Assignee"

	loggingFields := []any{
		"operation", op,
		"assign_id", assignId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	assignee, err := s.repo.GetUser(ctx, assignId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeAssignee(assignId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return assignee, nil
}

func (s *service) ChangeReporter(ctx context.Context, taskId, reporterId uint32) (models.User, error) {
	const op = "Change Task Reporter"

	loggingFields := []any{
		"operation", op,
		"reporter_id", reporterId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	reporter, err := s.repo.GetUser(ctx, reporterId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeReporter(reporterId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return reporter, nil
}

func (s *service) ChangeSprint(ctx context.Context, userId, taskId, newSprintId uint32) (models.Sprint, error) {
	const op = "Change Task Sprint"

	loggingFields := []any{
		"operation", op,
		"sprint_id", newSprintId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.Sprint{}, logError(err, s.logger, loggingFields...)
	}

	_, err = s.repo.GetUser(ctx, newSprintId)
	if err != nil {
		return models.Sprint{}, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeSprint(newSprintId)
	if err != nil {
		return models.Sprint{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return task.Sprint, nil
}
