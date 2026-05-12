package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/user"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type User = user.User

type UserRepo interface {
	Create(ctx context.Context, user User) (*User, error)

	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	ListActive(ctx context.Context) ([]*User, error)
	List(ctx context.Context) ([]*User, error)

	Update(ctx context.Context, id uuid.UUID, user User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user User) (*User, error) {
	const query = `
		INSERT INTO users (id, team_id, email, password, role, personal_data_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
			user.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &user, nil
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
		user.IsActive,
	)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active
		FROM users
		WHERE email = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		err := tx.QueryRowContext(
			ctx,
			query,
			email,
		).Scan(
			&user.ID,
			&user.TeamID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.PersonalDataID,
			&user.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &user, nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.TeamID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PersonalDataID,
		&user.IsActive,
	)
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return &user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User

	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active
		FROM users
		WHERE id = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		err := tx.QueryRowContext(
			ctx,
			query,
			id,
		).Scan(
			&user.ID,
			&user.TeamID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.PersonalDataID,
			&user.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return &user, nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.TeamID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PersonalDataID,
		&user.IsActive,
	)
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return &user, nil
}

func (r *userRepo) ListActive(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active
		FROM users
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
		return nil, dberrors.Map(err)
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.TeamID,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.PersonalDataID,
			&u.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return users, nil
}

func (r *userRepo) List(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active
		FROM users
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

	var users []*User

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.TeamID,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.PersonalDataID,
			&u.IsActive,
		)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return users, nil
}

func (r *userRepo) Update(ctx context.Context, id uuid.UUID, user User) (*User, error) {
	const query = `
		UPDATE users
		SET
			team_id = $1,
			email = $2,
			password = $3,
			role = $4,
			personal_data_id = $5,
			is_active = $6
		WHERE id = $7
	`

	var err error
	if tx, ok := db.GetTx(ctx); ok {
		_, err = tx.ExecContext(
			ctx,
			query,
			user.TeamID,
			user.Email,
			user.Password,
			user.Role,
			user.PersonalDataID,
			user.IsActive,
			id,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			user.TeamID,
			user.Email,
			user.Password,
			user.Role,
			user.PersonalDataID,
			user.IsActive,
			id,
		)
	}
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return &user, nil
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        SELECT deleted_at
        FROM users
        WHERE id = $1
    `

	var deletedAt sql.NullTime

	var err error

	if tx, ok := db.GetTx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, id).Scan(&deletedAt)
	} else {
		err = r.db.QueryRowContext(ctx, query, id).Scan(&deletedAt)
	}

	fmt.Printf("%+v\n", deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common_errors.ErrNotFound
		}

		return dberrors.Map(err)
	}

	if deletedAt.Valid {
		return common_errors.ErrConflict
	}

	updateQuery := `
        UPDATE users
        SET deleted_at = NOW(),
			is_active = false
        WHERE id = $1
    `

	if tx, ok := db.GetTx(ctx); ok {
		_, err = tx.ExecContext(ctx, updateQuery, id)
	} else {
		_, err = r.db.ExecContext(ctx, updateQuery, id)
	}

	if err != nil {
		return dberrors.Map(err)
	}

	return nil
}
