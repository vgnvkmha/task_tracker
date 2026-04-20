package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/team"
)

type Team = team.Team

type TeamRepo interface {
	Create(ctx context.Context, team Team) (Team, error)
	GetByName(ctx context.Context, teamName string) (Team, error)
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

func (r *teamRepo) GetByName(ctx context.Context, teamName string) (Team, error) {
	var team Team

	query := `
		SELECT *
		FROM team
		WHERE name = $1
	`

	err := r.db.QueryRowContext(ctx, query, teamName).Scan(
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
