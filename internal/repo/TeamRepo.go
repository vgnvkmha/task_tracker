package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/team"

	"github.com/google/uuid"
)

type TeamRepo interface {
	Create(ctx context.Context, t team.Team) (team.Team, error)
	Get(ctx context.Context, teamId uuid.UUID) (team.Team, error)
	Update(ctx context.Context, t team.Team) (team.Team, error)
}

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) TeamRepo {
	return &teamRepo{
		db: db,
	}
}

func (r *teamRepo) Create(ctx context.Context, t team.Team) (team.Team, error) {
	const query = `
		INSERT INTO teams (id, name)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		t.ID,
		t.Name,
	)
	if err != nil {
		return team.Team{}, err
	}

	return t, nil
}

func (r *teamRepo) Get(ctx context.Context, teamId uuid.UUID) (team.Team, error) {
	var t team.Team

	query := `
		SELECT *
		FROM team
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, teamId).Scan(
		&t.ID,
		&t.Name,
	)

	if err != nil {
		return team.Team{}, err
	}

	return t, nil

}

func (r *teamRepo) Update(ctx context.Context, t team.Team) (team.Team, error) {
	const query = `
		UPDATE teams
		SET
			name = $1,
		WHERE id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		t.Name,
		t.ID,
	)
	if err != nil {
		return team.Team{}, err
	}
	return t, nil
}
