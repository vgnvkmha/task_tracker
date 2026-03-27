package repo

import (
	"context"
	"task_tracker/internal/domain/models"
)

type ITaskRepo interface {
	createTask(task models.Task, ctx context.Context) (models.Task, error)
}

type TaskRepo struct {
}

func (r *TaskRepo) createTask(task models.Task, ctx context.Context) (models.Task, error) {
	return models.Task{}, nil
}
