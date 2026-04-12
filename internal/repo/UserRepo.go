package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Get(ctx context.Context, userId uuid.UUID) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	const query = `
		INSERT INTO users (id, team_id, personal_data_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Id,
		user.TeamId,
		user.DataId,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) Get(ctx context.Context, userId uuid.UUID) (models.User, error) {
	var user models.User

	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.TeamId,
		&user.DataId,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user models.User) (models.User, error) {
	//TODO: troubles are possible
	const query = `
		UPDATE users
		SET
			team_id = $1,
			personal_data_id = $2,
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.TeamId,
		user.DataId,
		user.Id,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
