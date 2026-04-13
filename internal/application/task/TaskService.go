package task_service

import (
	"context"
	"task_tracker/internal/domain/auth"
	task_errors "task_tracker/internal/domain/errors"
	"task_tracker/internal/domain/models"
	vo "task_tracker/internal/domain/models/value_objects"
	"task_tracker/internal/domain/validation"
	dto "task_tracker/internal/handler/task/dto"
	"task_tracker/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	module = "task"
	layer  = "service"
)

type TaskService interface {
	Create(ctx context.Context, actor auth.Actor, task dto.TaskRequest) (models.Task, error)

	GetActiveTasksByTeam(ctx context.Context, actor auth.Actor, teamId uuid.UUID) ([]models.Task, error)
	GetTeamById(ctx context.Context, actor auth.Actor, teamId uuid.UUID) (models.Team, error)

	ChangeStatus(ctx context.Context, actor auth.Actor, taskId uuid.UUID, newStatus string) (vo.Status, error)
	ChangeBoard(ctx context.Context, actor auth.Actor, taskId, newBoardId uuid.UUID) (models.Board, error)
	ChangeAssign(ctx context.Context, actor auth.Actor, taskId, newAssignId uuid.UUID) (models.User, error)
	ChangeReporter(ctx context.Context, actor auth.Actor, taskId, newreporterId uuid.UUID) (models.User, error)
	ChangeSprint(ctx context.Context, actor auth.Actor, taskId, newSprintId uuid.UUID) (models.Sprint, error)
}

type service struct {
	repo   repo.TaskRepo
	logger *zap.SugaredLogger
}

func New(repo repo.TaskRepo, logger *zap.SugaredLogger) TaskService {
	return &service{
		repo: repo,
		logger: logger.With(
			"module", module,
			"layer", layer,
		),
	}
}

func (s *service) Create(ctx context.Context, actor auth.Actor, task dto.TaskRequest) (models.Task, error) {
	const op = "Create Task"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
	}

	user, err := s.repo.GetUser(ctx, actor.Id)
	if err != nil {
		return models.Task{}, logError(err, s.logger, loggingFields...)
	}

	usersData, err := s.repo.GetPersonalData(ctx, user.DataId)
	if err != nil {
		return models.Task{}, logError(err, s.logger, loggingFields...)
	}
	if !usersData.Role.IsManagerRole() {
		return models.Task{}, logError(task_errors.ErrInvalidRights, s.logger, loggingFields...)
	}

	id := uuid.New()
	model, err := models.NewTask(
		id,
		task.Name,
		task.Description,
		task.BoardID,
		task.DueTo,
		task.AssigneeID,
		task.ReporetID,
		task.SprintId,
	)
	if err != nil {
		return models.Task{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return s.repo.Create(ctx, model)
}

func (s *service) GetActiveTasksByTeam(ctx context.Context, actor auth.Actor, teamId uuid.UUID) ([]models.Task, error) {
	const op = "GetActiveTaskByTeam"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
		"team_id", teamId,
	}

	user, err := s.repo.GetUser(ctx, actor.Id)
	if err != nil {
		return nil, logError(err, s.logger, loggingFields...)
	}
	usersData, err := s.repo.GetPersonalData(ctx, user.DataId)
	if err != nil {
		return nil, logError(err, s.logger, loggingFields...)
	}
	if err = validation.IsAlloweToSeeTeamData(usersData.Role, user.TeamId, teamId); err != nil {
		return nil, err
	}

	tasks, err := s.repo.GetActiveTasksByTeam(ctx, teamId)
	if err != nil {
		return nil, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return tasks, nil
}

func (s *service) GetTeamById(ctx context.Context, actor auth.Actor, teamId uuid.UUID) (models.Team, error) {
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

func (s *service) ChangeStatus(ctx context.Context, actor auth.Actor, taskId uuid.UUID, newStatus string) (vo.Status, error) {
	const op = "Change Status"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
		"task_id", taskId,
		"new_status", newStatus,
	}
	newStatusVo, err := validation.ParseStatus(newStatus)
	if err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}

	user, err := s.repo.GetUser(ctx, actor.Id)
	if err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}
	usersData, err := s.repo.GetPersonalData(ctx, user.DataId)
	if err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}
	if !usersData.Role.IsManagerRole() || newStatusVo.IsValid() != nil || newStatusVo.IsImmutable() != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeStatus(newStatusVo)
	if err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}

	if err = s.repo.Update(ctx, task); err != nil {
		return newStatusVo, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return newStatusVo, nil
}

func (s *service) ChangeBoard(ctx context.Context, actor auth.Actor, taskId, newBoardId uuid.UUID) (models.Board, error) {
	const op = "Change Task Board"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
		"board_id", newBoardId,
		"task_id", taskId,
	}

	user, err := s.repo.GetUser(ctx, actor.Id)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}
	usersData, err := s.repo.GetPersonalData(ctx, user.DataId)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	if !usersData.Role.IsManagerRole() {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}
	board, err := s.repo.GetBoard(ctx, newBoardId)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	if err := task.ChangeBoard(newBoardId); err != nil {
		return models.Board{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return board, nil
}

func (s *service) ChangeAssign(ctx context.Context, actor auth.Actor, taskId, newAssignId uuid.UUID) (models.User, error) {
	const op = "Change Task Assignee"

	loggingFields := []any{
		"operation", op,
		"assign_id", newAssignId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	assignee, err := s.repo.GetUser(ctx, newAssignId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeAssignee(&newAssignId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return assignee, nil
}

func (s *service) ChangeReporter(ctx context.Context, actor auth.Actor, taskId, newreporterId uuid.UUID) (models.User, error) {
	const op = "Change Task Reporter"

	loggingFields := []any{
		"operation", op,
		"reporter_id", newreporterId,
		"task_id", taskId,
		"user_role", "undefined",
	}

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	reporter, err := s.repo.GetUser(ctx, newreporterId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	err = task.ChangeReporter(newreporterId)
	if err != nil {
		return models.User{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return reporter, nil
}

func (s *service) ChangeSprint(ctx context.Context, actor auth.Actor, taskId, newSprintId uuid.UUID) (models.Sprint, error) {
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
	newSprint, err := s.repo.GetSprint(ctx, newSprintId)
	if err != nil {
		return models.Sprint{}, logError(err, s.logger, loggingFields...)
	}
	err = task.ChangeSprint(&newSprintId)
	if err != nil {
		return models.Sprint{}, logError(err, s.logger, loggingFields...)
	}

	logSuccess(s.logger, loggingFields...)
	return newSprint, nil
}
