package team

import (
	"context"
	"errors"
	"task_tracker/internal/application/common"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/repo/team"
	user_repo "task_tracker/internal/repo/user"

	domain_team "task_tracker/internal/domain/team"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Team = domain_team.Team

type TeamService interface {
	Create(ctx context.Context, input CreateTeamInput) (*Team, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Team, error)
	GetByName(ctx context.Context, name string) (*Team, error)

	ListActive(ctx context.Context) ([]*Team, error)
	List(ctx context.Context) ([]*Team, error)

	Update(ctx context.Context, id uuid.UUID, input *UpdateTeamInput) (*Team, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type service struct {
	teamRepo team.TeamRepo
	userRepo user_repo.UserRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(teamRepo team.TeamRepo, userRepo user_repo.UserRepo,
	logger *zap.SugaredLogger,
	transaction common.TxManager) TeamService {
	return &service{
		teamRepo:    teamRepo,
		userRepo:    userRepo,
		logger:      logger,
		transaction: transaction,
	}
}

func (s *service) Create(ctx context.Context, input CreateTeamInput) (*Team, error) {
	var result *Team

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if err := s.ensureTeamNameUnique(ctx, input.Name); err != nil {
			return err
		}

		if err := s.validateLeader(ctx, input.LeaderID, input.Name); err != nil {
			return err
		}

		team, err := domain_team.New(input.Name, input.Timezone, input.LeaderID)
		if err != nil {
			return err
		}

		created, err := s.teamRepo.Create(ctx, *team)
		if err != nil {
			return mapCreateError(err)
		}

		result = created
		return nil
	})

	if err != nil {
		s.logger.Errorw("create team failed",
			"error", err,
			"team_name", input.Name,
			"leader_id", input.LeaderID,
		)
		return nil, err
	}

	s.logger.Infow("team created",
		"team_id", result.ID,
		"name", result.Name,
		"leader_id", result.LeaderID,
	)

	return result, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	var result *Team
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		team, err := s.teamRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}
		result = team
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetByName(ctx context.Context, name string) (*Team, error) {
	var result *Team
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		team, err := s.teamRepo.GetByName(ctx, name)
		if err != nil {
			return mapGetError(err)
		}
		result = team
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) ListActive(ctx context.Context) ([]*Team, error) {
	var result []*Team
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		team, err := s.teamRepo.ListActive(ctx)
		if err != nil {
			return mapGetError(err)
		}
		result = team
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) List(ctx context.Context) ([]*Team, error) {
	var result []*Team
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		team, err := s.teamRepo.List(ctx)
		if err != nil {
			return mapGetError(err)
		}
		result = team
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *UpdateTeamInput) (*Team, error) {
	var result *Team
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		domainTeam, err := s.teamRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}
		if input.LeaderID != nil {
			_, err := s.userRepo.GetByID(ctx, *input.LeaderID)
			if err != nil {
				return mapGetError(err)
			}
		}
		domainTeam.ApplyChanges(input.Name, input.Timezone, input.LeaderID, input.IsActive)
		team, err := s.teamRepo.Update(ctx, id, *domainTeam)
		if err != nil {
			return mapGetError(err)
		}
		result = team
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		_, err := s.teamRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}

		err = s.teamRepo.Delete(ctx, id)
		if err != nil {
			return mapDeleteError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// helpers

func (s *service) ensureTeamNameUnique(ctx context.Context, name string) error {
	existing, err := s.teamRepo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, common_errors.ErrNotFound) {
		s.logger.Errorw("failed to get team by name",
			"error", err,
			"name", name,
		)
		return err
	}
	if existing != nil {
		return ErrTeamAlreadyExists
	}
	return nil
}

func (s *service) validateLeader(ctx context.Context, leaderID *uuid.UUID, teamName string) error {
	if leaderID == nil {
		return nil
	}

	leader, err := s.userRepo.GetByID(ctx, *leaderID)
	if err != nil {
		if errors.Is(err, common_errors.ErrNotFound) {
			return ErrLeaderNotFound
		}
		s.logger.Errorw("failed to get leader by id",
			"error", err,
			"leader_id", leaderID,
		)
		return err
	}

	if !leader.IsActive {
		return ErrLeaderInactive
	}

	team, err := s.teamRepo.GetByLeaderID(ctx, *leaderID)
	if err != nil && !errors.Is(err, common_errors.ErrNotFound) {
		s.logger.Errorw("failed to get team by leader id",
			"error", err,
			"leader_id", leaderID,
			"team_name", teamName,
		)
		return err
	}
	if team != nil {
		return ErrLeaderAlreadyHasTeam
	}

	return nil
}

func mapCreateError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return ErrPermissionDenied
	case errors.Is(err, common_errors.ErrAlreadyExists):
		return ErrTeamAlreadyExists
	case errors.Is(err, common_errors.ErrConflict):
		return team.ErrConflict
	case errors.Is(err, common_errors.ErrInvalidID):
		return team.ErrInvalidLeader
	default:
		return err
	}
}

func mapGetError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrNotFound):
		return common_errors.ErrNotFound
	case errors.Is(err, common_errors.ErrInvalidID):
		return common_errors.ErrInvalidID
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return common_errors.ErrPermissionDenied
	default:
		return err
	}
}

func mapUpdateError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return ErrPermissionDenied

	case errors.Is(err, common_errors.ErrNotFound):
		return ErrTeamNotFound

	case errors.Is(err, common_errors.ErrConflict):
		return team.ErrConflict

	case errors.Is(err, common_errors.ErrInvalidArgument):
		return ErrInvalidInput

	default:
		return err
	}
}

func mapDeleteError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return ErrPermissionDenied

	case errors.Is(err, common_errors.ErrNotFound):
		return ErrTeamNotFound

	default:
		return err
	}
}
