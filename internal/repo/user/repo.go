package user

import (
	"context"
	"database/sql"
	"errors"
	"task_tracker/internal/common_errors"
	personaldata "task_tracker/internal/domain/personal_data"
	"task_tracker/internal/domain/user"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type User = user.User

type ActiveUserProfile struct {
	User         *User
	PersonalData personaldata.PersonalData
	TeamName     string
}

type UserRepo interface {
	Create(ctx context.Context, user User) (*User, error)

	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	ListActive(ctx context.Context) ([]*User, error)
	ListActiveProfiles(ctx context.Context) ([]*ActiveUserProfile, error)
	ListProfiles(ctx context.Context) ([]*ActiveUserProfile, error)
	List(ctx context.Context) ([]*User, error)

	Update(ctx context.Context, id uuid.UUID, user User) (*User, error)
	Restore(ctx context.Context, id uuid.UUID) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	db *sql.DB
}

type userRowScanner interface {
	Scan(dest ...any) error
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func scanUser(scanner userRowScanner) (*User, error) {
	var (
		user      User
		teamID    uuid.NullUUID
		deletedAt sql.NullTime
	)

	err := scanner.Scan(
		&user.ID,
		&teamID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PersonalDataID,
		&user.IsActive,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if teamID.Valid {
		user.TeamID = &teamID.UUID
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func scanActiveUserProfile(scanner userRowScanner) (*ActiveUserProfile, error) {
	var (
		u            User
		teamID       uuid.NullUUID
		personalID   uuid.NullUUID
		firstName    sql.NullString
		lastName     sql.NullString
		age          sql.NullInt64
		birthDate    sql.NullTime
		deletedAt    sql.NullTime
		teamName     string
		personalData personaldata.PersonalData
	)

	err := scanner.Scan(
		&u.ID,
		&teamID,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.PersonalDataID,
		&u.IsActive,
		&deletedAt,
		&personalID,
		&firstName,
		&lastName,
		&age,
		&birthDate,
		&teamName,
	)
	if err != nil {
		return nil, err
	}

	if teamID.Valid {
		u.TeamID = &teamID.UUID
	}
	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}
	personalData.Id = u.PersonalDataID
	if personalID.Valid {
		personalData.Id = personalID.UUID
	}
	if firstName.Valid {
		personalData.FirstName = firstName.String
	}
	if lastName.Valid {
		personalData.LastName = lastName.String
	}
	if age.Valid {
		value := uint8(age.Int64)
		personalData.Age = &value
	}
	if birthDate.Valid {
		value := birthDate.Time
		personalData.BirthDate = &value
	}

	return &ActiveUserProfile{
		User:         &u,
		PersonalData: personalData,
		TeamName:     teamName,
	}, nil
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
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active, deleted_at
		FROM users
		WHERE email = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		user, err := scanUser(tx.QueryRowContext(
			ctx,
			query,
			email,
		))
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return user, nil
	}

	user, err := scanUser(r.db.QueryRowContext(
		ctx,
		query,
		email,
	))
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active, deleted_at
		FROM users
		WHERE id = $1
	`

	if tx, ok := db.GetTx(ctx); ok {
		user, err := scanUser(tx.QueryRowContext(
			ctx,
			query,
			id,
		))
		if err != nil {
			return nil, dberrors.Map(err)
		}
		return user, nil
	}

	user, err := scanUser(r.db.QueryRowContext(
		ctx,
		query,
		id,
	))
	if err != nil {
		return nil, dberrors.Map(err)
	}
	return user, nil
}

func (r *userRepo) ListActive(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active, deleted_at
		FROM users
		WHERE is_active = true
			AND deleted_at IS NULL
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
		u, err := scanUser(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return users, nil
}

func (r *userRepo) ListActiveProfiles(ctx context.Context) ([]*ActiveUserProfile, error) {

	const query = `
		SELECT
			u.id,
			u.team_id,
			u.email,
			u.password,
			u.role,
			u.personal_data_id,
			u.is_active,
			u.deleted_at,
			pd.id,
			pd.first_name,
			pd.last_name,
			pd.age,
			pd.birth_date,
			COALESCE(t.name, '')
		FROM users u
		LEFT JOIN personal_datas pd ON pd.id = u.personal_data_id
		LEFT JOIN teams t ON t.id = u.team_id
		WHERE u.is_active = true
			AND u.deleted_at IS NULL
		ORDER BY pd.first_name NULLS LAST, pd.last_name NULLS LAST, u.email
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
	profiles := make([]*ActiveUserProfile, 0)
	for rows.Next() {
		profile, err := scanActiveUserProfile(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return profiles, nil
}

func (r *userRepo) ListProfiles(ctx context.Context) ([]*ActiveUserProfile, error) {
	const query = `
		SELECT
			u.id,
			u.team_id,
			u.email,
			u.password,
			u.role,
			u.personal_data_id,
			u.is_active,
			u.deleted_at,
			pd.id,
			pd.first_name,
			pd.last_name,
			pd.age,
			pd.birth_date,
			COALESCE(t.name, '')
		FROM users u
		LEFT JOIN personal_datas pd ON pd.id = u.personal_data_id
		LEFT JOIN teams t ON t.id = u.team_id
		ORDER BY pd.first_name NULLS LAST, pd.last_name NULLS LAST, u.email
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

	profiles := make([]*ActiveUserProfile, 0)
	for rows.Next() {
		profile, err := scanActiveUserProfile(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return profiles, nil
}

func (r *userRepo) List(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, team_id, email, password, role, personal_data_id, is_active, deleted_at
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
		u, err := scanUser(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}

		users = append(users, u)
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
			AND deleted_at IS NULL
	`

	var (
		result sql.Result
		err    error
	)
	if tx, ok := db.GetTx(ctx); ok {
		result, err = tx.ExecContext(
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
		result, err = r.db.ExecContext(
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
	if rows, err := result.RowsAffected(); err == nil && rows == 0 {
		if deletedAt, err := r.deletedAt(ctx, id); err != nil {
			return nil, err
		} else if deletedAt.Valid {
			return nil, common_errors.ErrConflict
		}
		return nil, common_errors.ErrNotFound
	}
	return &user, nil
}

func (r *userRepo) deletedAt(ctx context.Context, id uuid.UUID) (sql.NullTime, error) {
	const query = `
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.NullTime{}, common_errors.ErrNotFound
		}
		return sql.NullTime{}, dberrors.Map(err)
	}
	return deletedAt, nil
}

func (r *userRepo) Restore(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `
		UPDATE users
		SET
			is_active = true,
			deleted_at = NULL
		WHERE id = $1
			AND (is_active = false OR deleted_at IS NOT NULL)
	`

	var (
		result sql.Result
		err    error
	)
	if tx, ok := db.GetTx(ctx); ok {
		result, err = tx.ExecContext(ctx, query, id)
	} else {
		result, err = r.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		return nil, dberrors.Map(err)
	}
	if rows, err := result.RowsAffected(); err == nil && rows == 0 {
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
		return nil, common_errors.ErrConflict
	}

	return r.GetByID(ctx, id)
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
