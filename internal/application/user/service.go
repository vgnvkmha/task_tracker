package user

import (
	"context"
	"errors"
	"task_tracker/internal/application/common"
	"task_tracker/internal/domain/auth"
	personaldata "task_tracker/internal/domain/personal_data"
	"task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"
	"task_tracker/internal/repo/team"
	user_repo "task_tracker/internal/repo/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	module = "user"
	layer  = "service"
)

type User = user.User

type UserService interface {
	CreateRegister(ctx context.Context, userInput CreateUserInput) (*User, error)
	CreateByActor(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (*User, error)
	Update(ctx context.Context, actor auth.Actor, userInput UpdateUserInput) (*User, error)
}

type service struct {
	userRepo user_repo.UserRepo
	dataRepo user_repo.PersonalDataRepo
	teamRepo team.TeamRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(userRepo user_repo.UserRepo, dataRepo user_repo.PersonalDataRepo, teamRepo team.TeamRepo, logger *zap.SugaredLogger, transaction common.TxManager) UserService {
	return &service{
		userRepo: userRepo,
		dataRepo: dataRepo,
		teamRepo: teamRepo,
		logger: logger.With(
			"module", module,
			"layer", layer,
		),
		transaction: transaction,
	}
}

func (s *service) CreateRegister(ctx context.Context, userInput CreateUserInput) (*User, error) {
	var u *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		var teamID uuid.UUID = uuid.Nil

		if userInput.TeamName != nil {
			team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
			if err != nil {
				return ErrTeamNotFound
			}
			teamID = team.ID
		}

		personalData, err := personaldata.New(
			userInput.FirstName,
			userInput.LastName,
			userInput.BirthDate,
			userInput.Age,
		)
		if err != nil {
			return err
		}

		if _, err = s.dataRepo.Create(ctx, *personalData); err != nil {
			return ErrPersonalDataCreateFailed
		}

		mappedUser, err := user.New(
			teamID,
			personalData.Id,
			userInput.Email,
			userInput.Password,
			*userInput.Role,
		)
		if err != nil {
			return err
		}

		createdUser, err := s.userRepo.Create(ctx, *mappedUser)
		if err != nil {
			if errors.Is(err, user.ErrAlreadyExists) {
				return ErrUserAlreadyExists
			}
			return ErrCreateUserFailed
		}

		u = createdUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) CreateByActor(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (*User, error) {
	var u *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		var teamID uuid.UUID = uuid.Nil
		actorRole := valueobjects.Role(actor.Role)
		if !actorRole.IsManagerRole() {
			return ErrOnlyManagersCanCreate
		}
		if userInput.TeamName != nil {
			team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
			if err != nil {
				return ErrTeamNotFound
			}
			teamID = team.ID
		}

		personalData, err := personaldata.New(
			userInput.FirstName,
			userInput.LastName,
			userInput.BirthDate,
			userInput.Age,
		)
		if err != nil {
			return err
		}

		if _, err = s.dataRepo.Create(ctx, *personalData); err != nil {
			return ErrPersonalDataCreateFailed
		}

		mappedUser, err := user.New(
			teamID,
			personalData.Id,
			userInput.Email,
			userInput.Password,
			*userInput.Role,
		)
		if err != nil {
			return err
		}

		createdUser, err := s.userRepo.Create(ctx, *mappedUser)
		if err != nil {
			return ErrUserCreateFailed
		}

		u = createdUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Update(ctx context.Context, actor auth.Actor, userInput UpdateUserInput) (*User, error) {
	var updatedUser *User

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {

		existingUser, err := s.userRepo.GetByEmail(ctx, *userInput.Email)
		if err != nil {
			return ErrUserNotFound
		}

		if userInput.TeamId != nil {
			team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
			if err != nil {
				return ErrTeamNotFound
			}
			existingUser.TeamID = &team.ID
		}

		if userInput.FirstName != nil ||
			userInput.LastName != nil ||
			userInput.BirthDate != nil ||
			userInput.Age != nil {

			pd, err := s.dataRepo.Get(ctx, existingUser.PersonalDataID)
			if err != nil {
				return ErrPersonalDataNotFound
			}

			if userInput.FirstName != nil {
				pd.FirstName = *userInput.FirstName
			}
			if userInput.LastName != nil {
				pd.LastName = *userInput.LastName
			}
			if userInput.BirthDate != nil {
				pd.BirthDate = userInput.BirthDate
			}
			if userInput.Age != nil {
				pd.Age = userInput.Age
			}

			if err := pd.Validate(); err != nil {
				return err
			}

			if _, err := s.dataRepo.Update(ctx, pd); err != nil {
				return ErrPersonalDataUpdateFailed
			}
		}

		email, err := valueobjects.NewEmail(*userInput.Email)
		if err != nil {
			return err
		}
		existingUser.Email = email

		password, err := valueobjects.NewPassword(*userInput.Password)
		if err != nil {
			return err
		}
		existingUser.Password = password

		if valueobjects.IsValidRole(*userInput.Role) {
			existingUser.Role = valueobjects.Role(*userInput.Role)
		} else {
			return ErrInvalidRole
		}

		savedUser, err := s.userRepo.Update(ctx, *existingUser)
		if err != nil {
			return ErrUserUpdateFailed
		}

		updatedUser = savedUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
