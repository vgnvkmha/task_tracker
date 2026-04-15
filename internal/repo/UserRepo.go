package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/user"

	"github.com/google/uuid"
)

type User = user.User

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, userId uuid.UUID) (User, error)
	Update(ctx context.Context, user User) (User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user User) (User, error) {
	//TODO: update entity
	const query = `
		INSERT INTO users (id, team_id, email, password, role, personal_data_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Id,
		user.TeamId,
		user.Email,
		user.Password,
		user.Role,
		user.PersonalDataId,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *userRepo) Get(ctx context.Context, userId uuid.UUID) (User, error) {
	var user User

	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.TeamId,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PersonalDataId,
	)

	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user User) (User, error) {
	//TODO: troubles are possible
	const query = `
		UPDATE users
		SET
			team_id = $1,
			email = $2,
			password = $3,
			role = $4,
			personal_data_id = $5,
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.TeamId,
		user.Email,
		user.Password,
		user.Role,
		user.PersonalDataId,
		user.Id,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
