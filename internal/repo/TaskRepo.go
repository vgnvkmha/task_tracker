package repo

import (
	"context"
	"database/sql"
	"fmt"
	"task_tracker/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type TaskRepo interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Update(ctx context.Context, task models.Task) error

	GetTask(ctx context.Context, taskId uuid.UUID) (models.Task, error)
	GetTeam(ctx context.Context, teamId uuid.UUID) (models.Team, error)
	GetUser(ctx context.Context, userId uuid.UUID) (models.User, error)
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	GetSprint(ctx context.Context, sprintId uuid.UUID) (models.Sprint, error)
	GetActiveTasksByTeam(ctx context.Context, teamId uuid.UUID) ([]models.Task, error)
}

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) TaskRepo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, task models.Task) (models.Task, error) {
	const query = `
		INSERT INTO task (id, name, description, status, created_at, due_to, board_id, reporter_id, assignee_id, sprint_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, &9, &10)
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Id,
		task.Name,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.DueTo,
		task.BoardId,
		task.ReporterId,
		task.AssigneeId,
		task.SprintId,
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("create task: %v", err)
	}

	return task, nil
}

func (r *repo) Update(ctx context.Context, task models.Task) error {
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

	return err
}

func (r *repo) GetTask(ctx context.Context, taskId uuid.UUID) (models.Task, error) {
	var task models.Task

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
		return models.Task{}, err
	}

	return task, nil
}

func (r *repo) GetTeam(ctx context.Context, teamId uuid.UUID) (models.Team, error) {
	var team models.Team

	query := `
		SELECT *
		FROM team
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, teamId).Scan(
		&team.ID,
		&team.Name,
		&team.UsersId,
	)

	if err != nil {
		return models.Team{}, err
	}

	return team, nil

}

func (r *repo) GetUser(ctx context.Context, userId uuid.UUID) (models.User, error) {
	var user models.User

	query := `
		SELECT *
		FROM user
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.TeamId,
		&user.Data,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *repo) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	var board models.Board

	query := `
		SELECT *
		FROM board
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, boardId).Scan(
		&board.Id,
		&board.TeamId,
		&board.IsPublic,
		&board.Name,
		&board.CreatedAt,
	)

	if err != nil {
		return models.Board{}, err
	}
	return board, nil
}

func (r *repo) GetSprint(ctx context.Context, sprintId uuid.UUID) (models.Sprint, error) {
	var sprint models.Sprint

	query := `
		SELECT *
		FROM sprint
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, sprintId).Scan(
		&sprint.Id,
		&sprint.Name,
		&sprint.StartDate,
		&sprint.EndDate,
		&sprint.Status,
		&sprint.BoardId,
		&sprint.TasksIds,
	)

	if err != nil {
		return models.Sprint{}, err
	}

	return sprint, nil
}

func (r *repo) GetActiveTasksByTeam(ctx context.Context, teamId uuid.UUID) ([]models.Task, error) {
	query := `
		SELECT id, name, description, status, board, due_to
		FROM task
		WHERE team_id = $1 AND status IN ($2, $3)
		ORDER BY due_to ASC
	`

	rows, err := r.db.QueryContext(ctx, query, teamId)
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
