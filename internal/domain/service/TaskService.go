package service

import (
	"context"
	"task_tracker/internal/domain/models"
	"task_tracker/internal/repo"
	dto "task_tracker/internal/transport/task"
)

type ITaskService interface {
	createTask(task dto.RequestTask, ctx context.Context) (models.Task, error)
	changeStatus(status string, ctx context.Context) error
	changeBoard(boardId uint32, ctx context.Context) error
}

type TaskService struct {
	repo *repo.ITaskRepo
}
