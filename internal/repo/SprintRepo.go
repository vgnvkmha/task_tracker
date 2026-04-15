package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/sprint"

	"github.com/google/uuid"
)

type Sprint = sprint.Sprint

type SprintRepo interface {
	Create(ctx context.Context, sprint Sprint) (Sprint, error)
	Get(ctx context.Context, sprintId uuid.UUID) (Sprint, error)
	Update(ctx context.Context, sprint Sprint) (Sprint, error)
}

type sprintRepo struct {
	db *sql.DB
}

func NewSprintRepo(db *sql.DB) SprintRepo {
	return &sprintRepo{
		db: db,
	}
}

func (r *sprintRepo) Create(ctx context.Context, sprint Sprint) (Sprint, error) {
	const query = `
		INSERT INTO sprints (id, name, start_date, end_date, status, board_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		sprint.Id,
		sprint.Name,
		sprint.StartDate,
		sprint.EndDate,
		sprint.Status,
		sprint.BoardId,
	)
	if err != nil {
		return Sprint{}, err
	}

	return sprint, nil
}

func (r *sprintRepo) Get(ctx context.Context, sprintId uuid.UUID) (Sprint, error) {
	var sprint Sprint

	query := `
		SELECT *
		FROM sprints
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, sprintId).Scan(
		&sprint.Id,
		&sprint.Name,
		&sprint.StartDate,
		&sprint.EndDate,
		&sprint.Status,
		&sprint.BoardId,
	)

	if err != nil {
		return Sprint{}, err
	}

	return sprint, nil
}

func (r *sprintRepo) Update(ctx context.Context, sprint Sprint) (Sprint, error) {
	const query = `
		UPDATE teams
		SET
			name = $1,
			start_date = $2,
			end_date = $3,
			status = $4,
			board_id = $5
		WHERE id = $7
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		sprint.Name,
		sprint.StartDate,
		sprint.EndDate,
		sprint.Status,
		sprint.BoardId,
		sprint.Id,
	)
	if err != nil {
		return Sprint{}, err
	}
	return sprint, nil
}
