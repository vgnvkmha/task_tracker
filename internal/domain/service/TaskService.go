package service

import (
	"context"
	"task_tracker/internal/domain/models"
	"task_tracker/internal/repo"
	dto "task_tracker/internal/transport/task"
)

type ITaskService interface {
	createTask(task dto.TaskRequest, ctx context.Context) (models.Task, error)
	changeStatus(status string, ctx context.Context) error
	changeBoard(boardId uint32, ctx context.Context) error
}

type TaskService struct {
	repo *repo.ITaskRepo
}

func (s *TaskService) createTask(task dto.TaskRequest, ctx context.Context) (models.Task, error) {
	model, err := models.NewTask(
		task.Name,
		task.Description,
		task.BoardID,
		task.AssigneeID,
		task.ReporetID,
		task.DueTo,
	)
	if err != nil {
		return models.Task{}, err
	}
	return model, nil
}
