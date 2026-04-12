package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/models"

	"github.com/google/uuid"
)

type PersonalDataRepo interface {
	Create(ctx context.Context, data models.PersonalData) (models.PersonalData, error)
	Get(ctx context.Context, dataId uuid.UUID) (models.PersonalData, error)
	Update(ctx context.Context, data models.PersonalData) (models.PersonalData, error)
}

type personalDataRepo struct {
	db *sql.DB
}

func NewPersonalDataRepo(db *sql.DB) PersonalDataRepo {
	return &personalDataRepo{
		db: db,
	}
}

func (r *personalDataRepo) Create(ctx context.Context, data models.PersonalData) (models.PersonalData, error) {
	const query = `
		INSERT INTO personal_datas (id, email. password, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		data.Id,
		data.Email,
		data.Password,
		data.Role,
	)
	if err != nil {
		return models.PersonalData{}, err
	}
	return data, nil
}

func (r *personalDataRepo) Get(ctx context.Context, dataId uuid.UUID) (models.PersonalData, error) {
	var data models.PersonalData

	const query = `
		SELECT *
		FROM personal_datas
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, dataId).Scan(
		&data.Id,
		&data.Email,
		&data.Password,
		&data.Role,
	)

	if err != nil {
		return models.PersonalData{}, err
	}
	return data, nil
}

func (r *personalDataRepo) Update(ctx context.Context, data models.PersonalData) (models.PersonalData, error) {
	const query = `
		UPDATE personal_datas
		SET
			email = $1,
			password = $2,
			role = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		data.Email,
		data.Password,
		data.Role,
		data.Id,
	)

	if err != nil {
		return models.PersonalData{}, err
	}

	return data, nil
}
