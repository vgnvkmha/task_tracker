package board

import (
	"context"
	"errors"

	"task_tracker/internal/application/common"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	domainboard "task_tracker/internal/domain/board"
	"task_tracker/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Board = domainboard.Board

type BoardService interface {
	Create(ctx context.Context, actor auth.Actor, input CreateBoardInput) (*Board, error)
	GetByID(ctx context.Context, actor auth.Actor, id uuid.UUID) (*Board, error)
	List(ctx context.Context, actor auth.Actor) ([]*Board, error)
	ListByTeamID(ctx context.Context, actor auth.Actor, teamID uuid.UUID) ([]*Board, error)
	Search(ctx context.Context, actor auth.Actor, input SearchBoardsInput) ([]*Board, error)
	Update(ctx context.Context, actor auth.Actor, id uuid.UUID, input UpdateBoardInput) (*Board, error)
	Delete(ctx context.Context, actor auth.Actor, id uuid.UUID) error
}

type service struct {
	boardRepo repo.BoardRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(boardRepo repo.BoardRepo, logger *zap.SugaredLogger, transaction common.TxManager) BoardService {
	return &service{
		boardRepo:   boardRepo,
		logger:      logger,
		transaction: transaction,
	}
}

func (s *service) Create(ctx context.Context, actor auth.Actor, input CreateBoardInput) (*Board, error) {
	var result *Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if !actor.Role.IsManagerRole() {
			return ErrPermissionDenied
		}

		board, err := domainboard.New(input.TeamID, input.IsPublic, input.Name)
		if err != nil {
			return mapDomainError(err)
		}

		created, err := s.boardRepo.Create(ctx, *board)
		if err != nil {
			return mapRepoError(err, ErrCreateBoardFailed)
		}

		result = created
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "create",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"name", input.Name,
		)
	}

	logSuccess(s.logger,
		"operation", "create",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"board_id", result.Id,
	)
	return result, nil
}

func (s *service) GetByID(ctx context.Context, actor auth.Actor, id uuid.UUID) (*Board, error) {
	var result *Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		board, err := s.boardRepo.GetByID(ctx, id)
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}
		if board == nil {
			return ErrBoardNotFound
		}

		result = board
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "get_by_id",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"board_id", id,
		)
	}

	return result, nil
}

func (s *service) List(ctx context.Context, actor auth.Actor) ([]*Board, error) {
	var result []*Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		boards, err := s.boardRepo.List(ctx)
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}

		result = boards
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "list",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
		)
	}

	return result, nil
}

func (s *service) ListByTeamID(ctx context.Context, actor auth.Actor, teamID uuid.UUID) ([]*Board, error) {
	var result []*Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		boards, err := s.boardRepo.ListByTeamID(ctx, teamID)
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}

		result = boards
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "list_by_team_id",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"team_id", teamID,
		)
	}

	return result, nil
}

func (s *service) Search(ctx context.Context, actor auth.Actor, input SearchBoardsInput) ([]*Board, error) {
	var result []*Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		boards, err := s.boardRepo.Search(ctx, repo.BoardSearchFilters{
			Query:  input.Query,
			TeamID: input.TeamID,
			UserID: input.UserID,
		})
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}

		result = boards
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "search",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
		)
	}

	return result, nil
}

func (s *service) Update(ctx context.Context, actor auth.Actor, id uuid.UUID, input UpdateBoardInput) (*Board, error) {
	var result *Board

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if !actor.Role.IsManagerRole() {
			return ErrPermissionDenied
		}

		existing, err := s.boardRepo.GetByID(ctx, id)
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}
		if existing == nil {
			return ErrBoardNotFound
		}

		if err := existing.ApplyChanges(input.TeamID, input.Name, input.IsPublic, input.Status); err != nil {
			return mapDomainError(err)
		}

		updated, err := s.boardRepo.Update(ctx, id, *existing)
		if err != nil {
			return mapRepoError(err, ErrUpdateBoardFailed)
		}
		if updated == nil {
			return ErrBoardNotFound
		}

		result = updated
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "update",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"board_id", id,
		)
	}

	logSuccess(s.logger,
		"operation", "update",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"board_id", result.Id,
	)
	return result, nil
}

func (s *service) Delete(ctx context.Context, actor auth.Actor, id uuid.UUID) error {
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if !actor.Role.IsManagerRole() {
			return ErrPermissionDenied
		}

		existing, err := s.boardRepo.GetByID(ctx, id)
		if err != nil {
			return mapRepoError(err, ErrBoardNotFound)
		}
		if existing == nil {
			return ErrBoardNotFound
		}

		if err := s.boardRepo.Delete(ctx, id); err != nil {
			return mapRepoError(err, ErrDeleteBoardFailed)
		}
		return nil
	})

	if err != nil {
		return logError(err, s.logger,
			"operation", "delete",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"board_id", id,
		)
	}

	logSuccess(s.logger,
		"operation", "delete",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"board_id", id,
	)
	return nil
}

func mapRepoError(err error, fallback error) error {
	switch {
	case errors.Is(err, common_errors.ErrNotFound):
		return ErrBoardNotFound
	case errors.Is(err, common_errors.ErrAlreadyExists):
		return ErrBoardAlreadyExists
	case errors.Is(err, common_errors.ErrInvalidArgument):
		return ErrInvalidInput
	case errors.Is(err, common_errors.ErrConflict):
		return ErrInvalidInput
	default:
		return fallback
	}
}

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domainboard.ErrEmptyName),
		errors.Is(err, domainboard.ErrNameTooLong),
		errors.Is(err, domainboard.ErrOwnerRequired),
		errors.Is(err, domainboard.ErrInvalidVisibility),
		errors.Is(err, domainboard.ErrInvalidCreationTime):
		return ErrInvalidInput
	case errors.Is(err, domainboard.ErrInvalidStatus):
		return ErrInvalidStatus
	case errors.Is(err, domainboard.ErrPermissionDenied),
		errors.Is(err, domainboard.ErrNotBoardMember):
		return ErrPermissionDenied
	default:
		return err
	}
}
