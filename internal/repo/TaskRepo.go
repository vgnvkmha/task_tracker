package repo

import (
	"context"
	"database/sql"
	"fmt"
	"task_tracker/internal/domain/models"
)

type TaskRepo interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Update(ctx context.Context, task models.Task) error
	GetTask(ctx context.Context, taskId uint32) (models.Task, error)
	GetTeam(ctx context.Context, teamId uint32) (models.Team, error)
	GetUser(ctx context.Context, userId uint32) (models.User, error)    //TODO: make realization
	GetBoard(ctx context.Context, boardId uint32) (models.Board, error) //TODO: make realization
	GetActiveByTeamId(ctx context.Context, teamId uint32) ([]models.Task, error)
}

type repo struct {
	db *sql.DB
}

func (r *repo) Create(ctx context.Context, task models.Task) (models.Task, error) {
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

func (r *repo) GetTask(ctx context.Context, id uint32) (models.Task, error) {
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

func (r *repo) Update(ctx context.Context, task models.Task) error {
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

func (r *repo) GetTeam(ctx context.Context, id uint32) (models.Team, error) {
	var team models.Team

	query := `
		SELECT *
		FROM team
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.Name,
		&team.UsersId,
	)

	if err != nil {
		return models.Team{}, err
	}

	return team, nil

}

func (r *repo) GetActiveByTeamId(ctx context.Context, id uint32) ([]models.Task, error) {
	query := `
		SELECT id, name, description, status, board, due_to
		FROM task
		WHERE team_id = $1 AND status IN ($2, $3)
		ORDER BY due_to ASC
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

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
