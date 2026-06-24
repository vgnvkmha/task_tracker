package repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/task"
	"task_tracker/internal/infrastracture/db"
	"task_tracker/internal/perf"
	"task_tracker/internal/repo/dberrors"

	"github.com/google/uuid"
)

type Task = task.Task

type TaskFilters struct {
	BoardID    *uuid.UUID
	AssigneeID *uuid.UUID
	ReporterID *uuid.UUID
	SprintID   *uuid.UUID
	Status     *task.TaskStatus
}

type TaskRepo interface {
	Create(ctx context.Context, task Task) (*Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	FindMany(ctx context.Context, filters TaskFilters) ([]*Task, error)
	FindByBoardID(ctx context.Context, boardID uuid.UUID) ([]*Task, error)
	FindByAssigneeID(ctx context.Context, assigneeID uuid.UUID) ([]*Task, error)
	Update(ctx context.Context, id uuid.UUID, task Task) (*Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepo struct {
	db *sql.DB
}

type rowScanner interface {
	Scan(dest ...any) error
}

func NewTaskRepo(db *sql.DB) TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) Create(ctx context.Context, task Task) (*Task, error) {
	const query = `
		INSERT INTO tasks (
			id,
			name,
			description,
			status,
			created_at,
			due_to,
			updated_at,
			reporter_id,
			assignee_id,
			board_id,
			sprint_id
		)
		VALUES ($1, $2, $3, $4, COALESCE($5, NOW()), $6, COALESCE($7, NOW()), $8, $9, $10, $11)
		RETURNING id, name, description, status, created_at, due_to, updated_at, reporter_id, assignee_id, board_id, sprint_id
	`

	result, err := scanTask(queryRow(ctx, r.db, query,
		task.Id,
		task.Name,
		nullableString(task.Description),
		nullableStatus(task.Status),
		nullableTime(task.CreatedAt),
		nullableTime(task.DueTo),
		nullableTime(task.UpdatedAt),
		task.ReporterId,
		task.AssigneeId,
		task.BoardId,
		task.SprintId,
	))

	if err != nil {
		return nil, dberrors.Map(err)
	}

	return result, nil
}

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	const query = `
		SELECT id, name, description, status, created_at, due_to, updated_at, reporter_id, assignee_id, board_id, sprint_id
		FROM tasks
		WHERE id = $1
	`

	task, err := scanTask(queryRow(ctx, r.db, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, dberrors.Map(err)
	}

	return task, nil
}

func (r *taskRepo) FindMany(ctx context.Context, filters TaskFilters) ([]*Task, error) {
	defer perf.Track(ctx, "repo.FindManyTasks.inner")()

	query, args := buildFindManyQuery(filters)

	queryDone := perf.Track(ctx, "db.Query")
	rows, err := queryRows(ctx, r.db, query, args...)
	queryDone()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, dberrors.Map(err)
	}
	defer rows.Close()

	scanDone := perf.Track(ctx, "rows.ScanAll")
	defer scanDone()
	return scanTasks(rows)
}

func (r *taskRepo) FindByBoardID(ctx context.Context, boardID uuid.UUID) ([]*Task, error) {
	return r.FindMany(ctx, TaskFilters{
		BoardID: &boardID,
	})
}

func (r *taskRepo) FindByAssigneeID(ctx context.Context, assigneeID uuid.UUID) ([]*Task, error) {
	return r.FindMany(ctx, TaskFilters{
		AssigneeID: &assigneeID,
	})
}

func (r *taskRepo) Update(ctx context.Context, id uuid.UUID, task Task) (*Task, error) {
	const query = `
		UPDATE tasks
		SET
			name = $1,
			description = $2,
			status = $3,
			due_to = $4,
			updated_at = NOW(),
			reporter_id = $5,
			assignee_id = $6,
			board_id = $7,
			sprint_id = $8
		WHERE id = $9
		RETURNING id, name, description, status, created_at, due_to, updated_at, reporter_id, assignee_id, board_id, sprint_id
	`

	result, err := scanTask(queryRow(ctx, r.db, query,
		task.Name,
		nullableString(task.Description),
		nullableStatus(task.Status),
		nullableTime(task.DueTo),
		task.ReporterId,
		task.AssigneeId,
		task.BoardId,
		task.SprintId,
		id,
	))

	if err != nil {
		return nil, dberrors.Map(err)
	}

	return result, nil
}

func (r *taskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM tasks
		WHERE id = $1
	`

	res, err := exec(ctx, r.db, query, id)
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

func buildFindManyQuery(filters TaskFilters) (string, []any) {
	query := `
		SELECT id, name, description, status, created_at, due_to, updated_at, reporter_id, assignee_id, board_id, sprint_id
		FROM tasks
	`

	var (
		conditions []string
		args       []any
	)

	addFilter := func(column string, value any) {
		args = append(args, value)
		conditions = append(conditions, column+" = $"+strconv.Itoa(len(args)))
	}

	if filters.BoardID != nil {
		addFilter("board_id", *filters.BoardID)
	}
	if filters.AssigneeID != nil {
		addFilter("assignee_id", *filters.AssigneeID)
	}
	if filters.ReporterID != nil {
		addFilter("reporter_id", *filters.ReporterID)
	}
	if filters.SprintID != nil {
		addFilter("sprint_id", *filters.SprintID)
	}
	if filters.Status != nil {
		addFilter("status", *filters.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC"

	return query, args
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, dberrors.Map(err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.Map(err)
	}

	return tasks, nil
}

func scanTask(scanner rowScanner) (*Task, error) {
	var (
		task        Task
		description sql.NullString
		status      sql.NullString
		dueTo       sql.NullTime
		assigneeID  uuid.NullUUID
		boardID     uuid.NullUUID
		sprintID    uuid.NullUUID
	)

	err := scanner.Scan(
		&task.Id,
		&task.Name,
		&description,
		&status,
		&task.CreatedAt,
		&dueTo,
		&task.UpdatedAt,
		&task.ReporterId,
		&assigneeID,
		&boardID,
		&sprintID,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		task.Description = description.String
	}
	if status.Valid {
		task.Status = taskDomainStatus(status.String)
	}
	if dueTo.Valid {
		task.DueTo = dueTo.Time
	}
	if assigneeID.Valid {
		task.AssigneeId = &assigneeID.UUID
	}
	if boardID.Valid {
		task.BoardId = &boardID.UUID
	}
	if sprintID.Valid {
		task.SprintId = &sprintID.UUID
	}

	return &task, nil
}

func taskDomainStatus(status string) task.TaskStatus {
	return task.TaskStatus(status)
}

func queryRow(ctx context.Context, sqlDB *sql.DB, query string, args ...any) *sql.Row {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return sqlDB.QueryRowContext(ctx, query, args...)
}

func queryRows(ctx context.Context, sqlDB *sql.DB, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}
	return sqlDB.QueryContext(ctx, query, args...)
}

func exec(ctx context.Context, sqlDB *sql.DB, query string, args ...any) (sql.Result, error) {
	if tx, ok := db.GetTx(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return sqlDB.ExecContext(ctx, query, args...)
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullableStatus(value task.TaskStatus) any {
	if value == "" {
		return nil
	}
	return value
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}
