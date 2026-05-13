package dberrors

import (
	"errors"
	"testing"

	"task_tracker/internal/common_errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapUniqueViolation(t *testing.T) {
	err := Map(&pgconn.PgError{
		Code:           "23505",
		ConstraintName: "users_email_key",
	})

	if !errors.Is(err, common_errors.ErrAlreadyExists) {
		t.Fatalf("Map() error = %v, want ErrAlreadyExists", err)
	}
}
