package user_repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/user"
	"task_tracker/internal/infrastracture/db"
)

type User = user.User

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, email string) (User, error)
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
	const query = `
		INSERT INTO users (id, team_id, email, password, role, personal_data_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if tx, ok := db.GetTx(ctx); ok {
		_, err := tx.ExecContext(
			ctx,
			query,
			user.ID,
			user.TeamID,
			user.Email,
			user.Password,
			user.Role,
			user.PersonalDataID,
		)
		if err != nil {
			return User{}, err
		}
		return user, nil
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.TeamID,
		user.Email,
		user.Password,
		user.Role,
		user.PersonalDataID,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *userRepo) Get(ctx context.Context, email string) (User, error) {
	var user User

	query := `
		SELECT *
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.TeamID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PersonalDataID,
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
		user.TeamID,
		user.Email,
		user.Password,
		user.Role,
		user.PersonalDataID,
		user.ID,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
