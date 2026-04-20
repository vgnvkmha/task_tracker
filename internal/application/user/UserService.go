package user

import (
	"context"
	"task_tracker/internal/application/common"
	"task_tracker/internal/domain/auth"
	personaldata "task_tracker/internal/domain/personal_data"
	"task_tracker/internal/domain/user"
	"task_tracker/internal/repo"
	user_repo "task_tracker/internal/repo/user"

	"go.uber.org/zap"
)

const (
	module = "user"
	layer  = "service"
)

type User = user.User

type UserService interface {
	CreateRegister(ctx context.Context, userInput CreateUserInput) (User, error)
	CreateByActor(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (User, error)
	Update(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (User, error)
}

type userService struct {
	userRepo user_repo.UserRepo
	dataRepo user_repo.PersonalDataRepo
	teamRepo repo.TeamRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(userRepo user_repo.UserRepo, dataRepo user_repo.PersonalDataRepo, logger *zap.SugaredLogger) UserService {
	return &userService{
		userRepo: userRepo,
		dataRepo: dataRepo,
		logger: logger.With(
			"module", module,
			"layer", layer,
		),
	}
}

func (s *userService) CreateRegister(ctx context.Context, userInput CreateUserInput) (User, error) {
	var mappedUser User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		var teamName *string
		if userInput.TeamName != nil {
			teamName = nil
		} else {
			teamName = userInput.TeamName
		}

		team, err := s.teamRepo.GetByName(ctx, *teamName)
		if err != nil {
			return err
		}
		personalData, err := personaldata.New(userInput.FirstName, userInput.LastName, userInput.BirthDate, userInput.Age)
		_, err = s.dataRepo.Create(ctx, *personalData)
		mappedUser, err := user.New(team.ID, personalData.Id, userInput.Email, userInput.Password, *userInput.Role)
		if err != nil {
			return err
		}
		createdUser, err := s.userRepo.Create(ctx, *mappedUser)
		if err != nil {
			return err
		}
		mappedUser.Id = createdUser.Id
		mappedUser.TeamId = createdUser.TeamId
		mappedUser.Email = createdUser.Email
		mappedUser.Password = createdUser.Password
		mappedUser.Role = createdUser.Role
		mappedUser.PersonalDataId = createdUser.PersonalDataId
		return nil
	})
	if err != nil {
		return User{}, err
	}
	return mappedUser, nil
}

// TODO: make personal data and user creation transaction
func (s *userService) CreateByActor(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (User, error) {
	const op = "create task"

	loggingFields := []any{
		"operation", op,
		"user_id", actor.Id,
		"user_role", actor.Role,
	}
	user, err := user
	return user, nil
}

func (s *userService) Update(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (User, error) {
	return nil
}
