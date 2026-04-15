package repo

import (
	"context"
	"database/sql"
	personaldata "task_tracker/internal/domain/personal_data"

	"github.com/google/uuid"
)

type PersonalDataRepo interface {
	Create(ctx context.Context, data personaldata.PersonalData) (personaldata.PersonalData, error)
	Get(ctx context.Context, dataId uuid.UUID) (personaldata.PersonalData, error)
	Update(ctx context.Context, data personaldata.PersonalData) (personaldata.PersonalData, error)
}

type personalDataRepo struct {
	db *sql.DB
}

func NewPersonalDataRepo(db *sql.DB) PersonalDataRepo {
	return &personalDataRepo{
		db: db,
	}
}

func (r *personalDataRepo) Create(ctx context.Context, data personaldata.PersonalData) (personaldata.PersonalData, error) {
	//TODO: change entity
	const query = `
		INSERT INTO personal_datas (id, first_name. last_name, age, birth_date)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		data.Id,
		data.FirstName,
		data.LastName,
		data.Age,
		data.BirthDate,
	)
	if err != nil {
		return personaldata.PersonalData{}, err
	}
	return data, nil
}

func (r *personalDataRepo) Get(ctx context.Context, dataId uuid.UUID) (personaldata.PersonalData, error) {
	var data personaldata.PersonalData

	const query = `
		SELECT *
		FROM personal_datas
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, dataId).Scan(
		&data.Id,
		&data.FirstName,
		&data.LastName,
		&data.Age,
		&data.BirthDate,
	)

	if err != nil {
		return personaldata.PersonalData{}, err
	}
	return data, nil
}

func (r *personalDataRepo) Update(ctx context.Context, data personaldata.PersonalData) (personaldata.PersonalData, error) {
	const query = `
		UPDATE personal_datas
		SET
			first_name = $1,
			last_name = $2,
			age = $3
			birth_date = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		data.FirstName,
		data.LastName,
		data.Age,
		data.BirthDate,
		data.Id,
	)

	if err != nil {
		return personaldata.PersonalData{}, err
	}

	return data, nil
}
