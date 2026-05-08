package user

import (
	"context"
	"database/sql"
	personaldata "task_tracker/internal/domain/personal_data"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type PersonalData = personaldata.PersonalData

type PersonalDataRepo interface {
	Create(ctx context.Context, data PersonalData) (PersonalData, error)
	Get(ctx context.Context, dataId uuid.UUID) (PersonalData, error)
	Update(ctx context.Context, data PersonalData) (PersonalData, error)
}

type personalDataRepo struct {
	db *sql.DB
}

func NewPersonalDataRepo(db *sql.DB) PersonalDataRepo {
	return &personalDataRepo{
		db: db,
	}
}

func (r *personalDataRepo) Create(ctx context.Context, data PersonalData) (PersonalData, error) {
	const query = `
		INSERT INTO personal_datas (id, first_name, last_name, age, birth_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	var err error
	if tx, ok := db.GetTx(ctx); ok {
		_, err = tx.ExecContext(
			ctx,
			query,
			data.Id,
			data.FirstName,
			data.LastName,
			data.Age,
			data.BirthDate,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			data.Id,
			data.FirstName,
			data.LastName,
			data.Age,
			data.BirthDate,
		)
	}
	if err != nil {
		return PersonalData{}, dberrors.Map(err)
	}
	return data, nil
}

func (r *personalDataRepo) Get(ctx context.Context, dataId uuid.UUID) (PersonalData, error) {
	var data PersonalData

	const query = `
		SELECT id, first_name, last_name, age, birth_date
		FROM personal_datas
		WHERE id = $1
	`

	var err error
	if tx, ok := db.GetTx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, dataId).Scan(
			&data.Id,
			&data.FirstName,
			&data.LastName,
			&data.Age,
			&data.BirthDate,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, dataId).Scan(
			&data.Id,
			&data.FirstName,
			&data.LastName,
			&data.Age,
			&data.BirthDate,
		)
	}

	if err != nil {
		return PersonalData{}, dberrors.Map(err)
	}
	return data, nil
}

func (r *personalDataRepo) Update(ctx context.Context, data PersonalData) (PersonalData, error) {
	const query = `
		UPDATE personal_datas
		SET
			first_name = $1,
			last_name = $2,
			age = $3,
			birth_date = $4
		WHERE id = $5
	`

	var err error
	if tx, ok := db.GetTx(ctx); ok {
		_, err = tx.ExecContext(
			ctx,
			query,
			data.FirstName,
			data.LastName,
			data.Age,
			data.BirthDate,
			data.Id,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			data.FirstName,
			data.LastName,
			data.Age,
			data.BirthDate,
			data.Id,
		)
	}

	if err != nil {
		return PersonalData{}, dberrors.Map(err)
	}

	return data, nil
}
