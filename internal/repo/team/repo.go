package team

import (
	"context"
	"database/sql"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/team"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type Team = team.Team

type TeamRepo interface {
	Create(ctx context.Context, team Team) (*Team, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Team, error)
	GetByName(ctx context.Context, name string) (*Team, error)
	GetByLeaderID(ctx context.Context, id uuid.UUID) (*Team, error)

	ListActive(ctx context.Context) ([]*Team, error)
	List(ctx context.Context) ([]*Team, error)

	Update(ctx context.Context, team Team) (*Team, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) TeamRepo {
	return &teamRepo{
		db: db,
	}
}

func (r *teamRepo) Create(ctx context.Context, team Team) (*Team, error) {
	const query = `
		INSERT INTO teams (id, name, timezone, leader_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
	`

	if tx, ok := db.GetTx(ctx); ok {

		_, err := tx.ExecContext(
			ctx,
			query,
			team.ID,
			team.Name,
			team.Timezone,
			team.LeaderID,
			team.IsActive,
		)
		if err != nil {
			return nil, err
		}
		return &team, nil

	}
	_, err := r.db.ExecContext(
		ctx,
		query,
		team.ID,
		team.Name,
		team.Timezone,
		team.LeaderID,
		team.IsActive,
	)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &team, nil
}

func (r *teamRepo) GetByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	var team Team
	const query = `
		SELECT *
		FROM teams
		WHERE id = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		err := tx.QueryRowContext(
			ctx,
			query,
			id,
		).Scan(
			&team.ID,
			&team.Name,
			&team.Timezone,
			&team.LeaderID,
			&team.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &team, nil
	}
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&team.ID,
		&team.Name,
		&team.Timezone,
		&team.LeaderID,
		&team.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &team, nil

}

func (r *teamRepo) GetByName(ctx context.Context, teamName string) (*Team, error) {
	var team Team

	query := `
		SELECT *
		FROM teams
		WHERE name = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		err := tx.QueryRowContext(
			ctx,
			query,
			teamName,
		).Scan(
			&team.ID,
			&team.Name,
			&team.Timezone,
			&team.LeaderID,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &team, nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		teamName,
	).Scan(
		&team.ID,
		&team.Name,
		&team.Timezone,
		&team.LeaderID,
	)

	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &team, nil

}

func (r *teamRepo) GetByLeaderID(ctx context.Context, id uuid.UUID) (*Team, error) {
	var team Team

	query := `
		SELECT *
		FROM teams
		WHERE leader_id = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		err := tx.QueryRowContext(
			ctx,
			query,
			id,
		).Scan(
			&team.ID,
			&team.Name,
			&team.Timezone,
			&team.LeaderID,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &team, nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&team.ID,
		&team.Name,
		&team.Timezone,
		&team.LeaderID,
	)

	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &team, nil

}

func (r *teamRepo) ListActive(ctx context.Context) ([]*Team, error) {
	const query = `
        SELECT id, name, timezone, leader_id
        FROM teams
        WHERE is_active = true
    `

	var (
		rows *sql.Rows
		err  error
	)

	if tx, ok := db.GetTx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*Team

	for rows.Next() {
		var t Team

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Timezone,
			&t.LeaderID,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		teams = append(teams, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return teams, nil
}

func (r *teamRepo) List(ctx context.Context) ([]*Team, error) {
	const query = `
        SELECT id, name, timezone, leader_id
        FROM teams
    `

	var (
		rows *sql.Rows
		err  error
	)

	if tx, ok := db.GetTx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, dberrors.Map(err)
	}
	defer rows.Close()

	var teams []*Team

	for rows.Next() {
		var t Team

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Timezone,
			&t.LeaderID,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		teams = append(teams, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return teams, nil
}

func (r *teamRepo) Update(ctx context.Context, team Team) (*Team, error) {
	const query = `
		UPDATE teams
		SET
			name = $1,
			timezone = $2,
			leader_id = $3,
			is_active = $4,
		WHERE id = $5
	`
	var err error
	if tx, ok := db.GetTx(ctx); ok {
		_, err = tx.ExecContext(
			ctx,
			query,
			team.Name,
			team.Timezone,
			team.LeaderID,
			team.IsActive,
			team.ID,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			team.Name,
			team.Timezone,
			team.LeaderID,
			team.IsActive,
			team.ID,
		)
	}
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return &team, nil
}

func (r *teamRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
        UPDATE teams
        SET is_active = false
        WHERE id = $1;
    `

	var (
		res sql.Result
		err error
	)

	if tx, ok := db.GetTx(ctx); ok {
		res, err = tx.ExecContext(ctx, query, id)
	} else {
		res, err = r.db.ExecContext(ctx, query, id)
	}

	if err != nil {
		return dberrors.Map(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return dberrors.Map(err)
	}

	if rowsAffected == 0 {
		return common_errors.ErrNotFound
	}

	return nil
}
