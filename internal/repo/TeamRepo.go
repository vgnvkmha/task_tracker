package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/team"

	"github.com/google/uuid"
)

type Team = team.Team

type TeamRepo interface {
	Create(ctx context.Context, team Team) (Team, error)
	Get(ctx context.Context, teamId uuid.UUID) (Team, error)
	Update(ctx context.Context, team Team) (Team, error)
}

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) TeamRepo {
	return &teamRepo{
		db: db,
	}
}

func (r *teamRepo) Create(ctx context.Context, team Team) (Team, error) {
	const query = `
		INSERT INTO teams (id, name)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		team.ID,
		team.Name,
	)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (r *teamRepo) Get(ctx context.Context, teamId uuid.UUID) (Team, error) {
	var team Team

	query := `
		SELECT *
		FROM team
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, teamId).Scan(
		&team.ID,
		&team.Name,
	)

	if err != nil {
		return Team{}, err
	}

	return team, nil

}

func (r *teamRepo) Update(ctx context.Context, team Team) (Team, error) {
	const query = `
		UPDATE teams
		SET
			name = $1,
		WHERE id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		team.Name,
		team.ID,
	)
	if err != nil {
		return Team{}, err
	}
	return team, nil
}
