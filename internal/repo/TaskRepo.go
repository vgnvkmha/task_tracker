package repo

import (
	"context"
	"database/sql"
	"fmt"
	"task_tracker/internal/domain/task"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Task = task.Task

type TaskRepo interface {
	Create(ctx context.Context, task Task) (Task, error)
	Get(ctx context.Context, taskId uuid.UUID) (Task, error)
	Update(ctx context.Context, task Task) (Task, error)

	GetActiveByTeam(ctx context.Context, teamId uuid.UUID) ([]Task, error)
}

type taskRepo struct {
	db *sql.DB
}

func New(db *sql.DB) TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) Create(ctx context.Context, task Task) (Task, error) {
	const query = `
		INSERT INTO tasks (id, name, description, status, created_at, due_to, updated_at, reporter_id, assignee_id, board_id, sprint_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Id,
		task.Name,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.DueTo,
		task.UpdatedAt,
		task.ReporterId,
		task.AssigneeId,
		task.BoardId,
		task.SprintId,
	)

	if err != nil {
		return Task{}, fmt.Errorf("create task: %v", err)
	}

	return task, nil
}

func (r *taskRepo) Get(ctx context.Context, taskId uuid.UUID) (Task, error) {
	var task Task

	query := `
		SELECT *
		FROM task
		WHERE id = $1
		LIMIT = 1
	`
	err := r.db.QueryRowContext(ctx, query, taskId).Scan(
		&task.Id,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.BoardId,
		&task.CreatedAt,
		&task.DueTo,
		&task.UpdatedAt,
		&task.AssigneeId,
		&task.ReporterId,
	)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Update(ctx context.Context, task Task) (Task, error) {
	const query = `
		UPDATE task
		SET 
			name = $1,
			description = $2,
			status = $3,
			board_id = $4,
			created_at = $5,
			due_to = $6,
			updated_at = $7,
			assignee_id = $8,
			reporter_id = $9,
			sprint_id = $10
		WHERE id = $11
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Status,
		task.BoardId,
		task.CreatedAt,
		task.DueTo,
		time.Now(),
		task.AssigneeId,
		task.ReporterId,
		task.SprintId,
		task.Id,
	)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *taskRepo) GetActiveByTeam(ctx context.Context, teamId uuid.UUID) ([]Task, error) {
	query := `
		SELECT id, name, description, status, board, due_to
		FROM tasks
		WHERE team_id = $1 AND status IN ($2, $3)
		ORDER BY due_to ASC
	`

	rows, err := r.db.QueryContext(ctx, query, teamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.Status,
			&t.BoardId,
			&t.DueTo,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
