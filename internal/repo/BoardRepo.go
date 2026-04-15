package repo

import (
	"context"
	"database/sql"
	"task_tracker/internal/domain/board"

	"github.com/google/uuid"
)

type BoardRepo interface {
	Create(ctx context.Context, b board.Board) (board.Board, error)
	Get(ctx context.Context, boardId uuid.UUID) (board.Board, error)
	Update(ctx context.Context, b board.Board) (board.Board, error)
}

type boardRepo struct {
	db *sql.DB
}

func NewBoardRepo(db *sql.DB) BoardRepo {
	return &boardRepo{
		db: db,
	}
}

func (r *boardRepo) Create(ctx context.Context, b board.Board) (board.Board, error) {
	const query = `
		INSERT INTO boards (id, team_id, is_public, name, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		b.Id,
		b.TeamId,
		b.IsPublic,
		b.Name,
		b.CreatedAt,
	)
	if err != nil {
		return board.Board{}, err
	}

	return b, nil
}

func (r *boardRepo) Get(ctx context.Context, boardId uuid.UUID) (board.Board, error) {
	var b board.Board

	query := `
		SELECT *
		FROM boards
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, boardId).Scan(
		&b.Id,
		&b.TeamId,
		&b.IsPublic,
		&b.Name,
		&b.CreatedAt,
	)

	if err != nil {
		return board.Board{}, err
	}
	return b, nil
}

func (r *boardRepo) Update(ctx context.Context, b board.Board) (board.Board, error) {
	//TODO: troubles are possible
	const query = `
		UPDATE boards
		SET
			team_id = $1,
			is_public = $2,
			name = $3,
			created_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		b.TeamId,
		b.IsPublic,
		b.Name,
		b.CreatedAt,
		b.Id,
	)

	if err != nil {
		return board.Board{}, err
	}

	return b, nil
}
