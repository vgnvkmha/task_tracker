package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/board"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type Board = board.Board

type BoardRepo interface {
	Create(ctx context.Context, board Board) (*Board, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Board, error)
	List(ctx context.Context) ([]*Board, error)
	ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]*Board, error)
	Update(ctx context.Context, id uuid.UUID, board Board) (*Board, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type boardRepo struct {
	db *sql.DB
}

type boardRowScanner interface {
	Scan(dest ...any) error
}

func NewBoardRepo(db *sql.DB) BoardRepo {
	return &boardRepo{
		db: db,
	}
}

func (r *boardRepo) Create(ctx context.Context, board Board) (*Board, error) {
	const query = `
		INSERT INTO boards (id, team_id, is_public, name, status, created_at)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, NOW()))
		RETURNING id, team_id, is_public, name, status, created_at
	`

	result, err := scanBoard(queryBoardRow(ctx, r.db, query,
		board.Id,
		board.TeamId,
		board.IsPublic,
		board.Name,
		board.Status,
		nullableBoardTime(board.CreatedAt),
	))
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return result, nil
}

func (r *boardRepo) GetByID(ctx context.Context, id uuid.UUID) (*Board, error) {
	const query = `
		SELECT id, team_id, is_public, name, status, created_at
		FROM boards
		WHERE id = $1
	`

	board, err := scanBoard(queryBoardRow(ctx, r.db, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, dberrors.Map(err)
	}

	return board, nil
}

func (r *boardRepo) List(ctx context.Context) ([]*Board, error) {
	const query = `
		SELECT id, team_id, is_public, name, status, created_at
		FROM boards
		ORDER BY created_at DESC
	`

	rows, err := queryBoardRows(ctx, r.db, query)
	if err != nil {
		return nil, dberrors.Map(err)
	}
	defer rows.Close()

	return scanBoards(rows)
}

func (r *boardRepo) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]*Board, error) {
	const query = `
		SELECT id, team_id, is_public, name, status, created_at
		FROM boards
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	rows, err := queryBoardRows(ctx, r.db, query, teamID)
	if err != nil {
		return nil, dberrors.Map(err)
	}
	defer rows.Close()

	return scanBoards(rows)
}

func (r *boardRepo) Update(ctx context.Context, id uuid.UUID, board Board) (*Board, error) {
	const query = `
		UPDATE boards
		SET
			team_id = $1,
			is_public = $2,
			name = $3,
			status = $4
		WHERE id = $5
		RETURNING id, team_id, is_public, name, status, created_at
	`

	result, err := scanBoard(queryBoardRow(ctx, r.db, query,
		board.TeamId,
		board.IsPublic,
		board.Name,
		board.Status,
		id,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, dberrors.Map(err)
	}

	return result, nil
}

func (r *boardRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM boards
		WHERE id = $1
	`

	res, err := execBoard(ctx, r.db, query, id)
	if err != nil {
		return dberrors.Map(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return dberrors.Map(err)
	}
	if rowsAffected == 0 {
		return common_errors.ErrNotFound
	}

	return nil
}

func scanBoards(rows *sql.Rows) ([]*Board, error) {
	var boards []*Board

	for rows.Next() {
		board, err := scanBoard(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return boards, nil
}

func scanBoard(scanner boardRowScanner) (*Board, error) {
	var board Board
	err := scanner.Scan(
		&board.Id,
		&board.TeamId,
		&board.IsPublic,
		&board.Name,
		&board.Status,
		&board.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &board, nil
}

func queryBoardRow(ctx context.Context, sqlDB *sql.DB, query string, args ...any) *sql.Row {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return sqlDB.QueryRowContext(ctx, query, args...)
}

func queryBoardRows(ctx context.Context, sqlDB *sql.DB, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}
	return sqlDB.QueryContext(ctx, query, args...)
}

func execBoard(ctx context.Context, sqlDB *sql.DB, query string, args ...any) (sql.Result, error) {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return sqlDB.ExecContext(ctx, query, args...)
}

func nullableBoardTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}
