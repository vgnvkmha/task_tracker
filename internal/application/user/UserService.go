package user_service

import (
	"context"
	"task_tracker/internal/domain/auth"
	"task_tracker/internal/domain/models"
	"task_tracker/internal/handler/task/dto"
	"task_tracker/internal/repo"

	"go.uber.org/zap"
)

const (
	module = "user"
	layer  = "service"
)

type UserService interface {
	Register(ctx context.Context, user models.User) (models.User, error)
	CreateByActor(ctx context.Context, actor auth.Actor, user models.User) (models.User, error)
	Update(ctx context.Context, actor auth.Actor, update dto.UpdateUser) (models.User, error)
}

type userService struct {
	repo   repo.UserRepo
	logger *zap.SugaredLogger
}

func New(repo repo.UserRepo, logger *zap.SugaredLogger) UserService {
	return &userService{
		repo: repo,
		logger: logger.With(
			"module", module,
			"layer", layer,
		),
	}
}

func (s *userService) Register(ctx context.Context, user models.User) (models.User, error) {

	return user, nil
}
func (s *userService) CreateByActor(ctx context.Context, actor auth.Actor, user models.User) (models.User, error) {
	const op = "Create Task"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, actor auth.Actor, update dto.UpdateUser) (models.User, error) {
	return nil
}
