package user_service

import (
	"context"
	"task_tracker/internal/domain/auth"
	"task_tracker/internal/domain/user"
	"task_tracker/internal/handler/task/dto"
	"task_tracker/internal/repo"

	"go.uber.org/zap"
)

const (
	module = "user"
	layer  = "service"
)

type User = user.User

type UserService interface {
	CreateRegister(ctx context.Context, user User) (User, error)
	CreateByActor(ctx context.Context, actor auth.Actor, user User) (User, error)
	Update(ctx context.Context, actor auth.Actor, update dto.UpdateUser) (User, error)
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

func (s *userService) CreateRegister(ctx context.Context, user User) (User, error) {

	return user, nil
}

// TODO: make personal data and user creation transaction
func (s *userService) CreateByActor(ctx context.Context, actor auth.Actor, user User) (User, error) {
	const op = "create task"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
	}
	user, err := user
	return user, nil
}

func (s *userService) Update(ctx context.Context, actor auth.Actor, update dto.UpdateUser) (User, error) {
	return nil
}
