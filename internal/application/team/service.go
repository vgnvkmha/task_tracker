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

	Update(ctx context.Context, team UpdateTeamInput) (*Team, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
	return &service{}
}

func (s *service) Create(ctx context.Context, input CreateTeamInput) (*Team, error) {
	_, err := s.teamRepo.GetByName(ctx, input.Name)
	if !errors.Is(err, common_errors.ErrNotFound) {
		return nil, ErrTeamAlreadyExists
	}
	if input.LeaderID != nil {
		leader, err := s.userRepo.GetByID(ctx, *input.LeaderID)
		if errors.Is(err, common_errors.ErrNotFound) {
			return nil, ErrLeaderNotFound
		}
		if err == nil {
			return nil, ErrLeaderAlreadyHasTeam
		}
		if !leader.IsActive {
			return nil, ErrLeaderInactive
		}
	}

	team, err := domain_team.New(input.Name, input.Timezone, input.LeaderID)
	if err != nil {
		return nil, err
	}
	createdTeam, err := s.teamRepo.Create(ctx, *team)
	if errors.Is(err, common_errors.ErrPermissionDenied) {
		return nil, ErrPermissionDenied
	}
	if errors.Is(err, common_errors.ErrAlreadyExists) {
		return nil, ErrTeamAlreadyExists
	}
	return createdTeam, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Team, error) {

}

func (s *service) GetByName(ctx context.Context, name string) (*Team, error) {}

func (s *service) ListActive(ctx context.Context) ([]*Team, error) {}

func (s *service) List(ctx context.Context) ([]*Team, error) {}

func (s *service) Update(ctx context.Context, team UpdateTeamInput) (*Team, error) {}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {}
