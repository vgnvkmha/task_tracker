package repo

import (
	"context"
	"database/sql"
	"fmt"
	"task_tracker/internal/domain/models"
)

type ITaskRepo interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Update(ctx context.Context, task models.Task) error
	Get(ctx context.Context, id uint32) (models.Task, error)
}

type TaskRepo struct {
	db *sql.DB
}

func (r *TaskRepo) Create(ctx context.Context, task models.Task) (models.Task, error) {
	const query = `
		INSERT INTO task (name, description, board_id, assignee_id, reporter_id, due_to, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	createdTask := task

	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.BoardId,
		task.AssigneeId,
		task.ReporterId,
		task.DueTo,
		task.Status,
	).Scan(
		&createdTask.Id,
		&createdTask.CreatedAt,
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("create task: %w", err)
	}

	return createdTask, nil
}

func (r *TaskRepo) Get(ctx context.Context, id uint32) (models.Task, error) {
	var task models.Task

	query := `
		SELECT *
		FROM task
		WHERE id = $1
		LIMIT = 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.Id,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.BoardId,
		&task.CreatedAt,
		&task.DueTo,
		&task.AssigneeId,
		&task.ReporterId,
	)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *TaskRepo) Update(ctx context.Context, task models.Task) error {
	query := `
		UPDATE task
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Status,
		task.Id,
	)

	return err
}
