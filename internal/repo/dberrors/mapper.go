package dberrors

import (
	"database/sql"
	"errors"

	"task_tracker/internal/common_errors"

	"github.com/jackc/pgconn"
)

func Map(err error) error {
	if err == nil {
		return nil
	}

	// common SELECT case
	if errors.Is(err, sql.ErrNoRows) {
		return common_errors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		// uniqueness
		case "23505": // unique_violation
			return common_errors.ErrAlreadyExists

		// foreign keys
		case "23503": // foreign_key_violation
			return common_errors.ErrInvalidArgument

		// not null
		case "23502": // not_null_violation
			return common_errors.ErrInvalidArgument

		// check constraints
		case "23514": // check_violation
			return common_errors.ErrInvalidArgument

		// invalid input / data
		case "22P02": // invalid_text_representation (ex: UUID)
			return common_errors.ErrInvalidArgument

		// serialization / concurrency
		case "40001": // serialization_failure
			return common_errors.ErrConflict

		// deadlock
		case "40P01": // deadlock_detected
			return common_errors.ErrConflict
		}
	}

	// fallback
	return err
}
