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
	userRepo "task_tracker/internal/repo/user"

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
	userRepo userRepo.UserRepo
	dataRepo userRepo.PersonalDataRepo
	teamRepo team.TeamRepo

	logger      *zap.SugaredLogger
	transaction common.TxManager
}

func New(userRepo userRepo.UserRepo, dataRepo userRepo.PersonalDataRepo, teamRepo team.TeamRepo, logger *zap.SugaredLogger, transaction common.TxManager) UserService {
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
				s.logger.Errorw("get team by name failed",
					"error", err,
					"operation", "create_register",
					"team_name", *userInput.TeamName,
				)
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
			s.logger.Errorw("build personal data failed",
				"error", err,
				"operation", "create_register",
				"email", userInput.Email,
			)
			return err
		}

		if _, err = s.dataRepo.Create(ctx, *personalData); err != nil {
			s.logger.Errorw("create personal data failed",
				"error", err,
				"operation", "create_register",
				"personal_data_id", personalData.Id,
			)
			return ErrPersonalDataCreateFailed
		}

		if userInput.Role == nil {
			s.logger.Errorw("role is required",
				"operation", "create_register",
				"email", userInput.Email,
			)
			return ErrRoleRequired
		}

		mappedUser, err := user.New(
			teamID,
			personalData.Id,
			userInput.Email,
			userInput.Password,
			*userInput.Role,
		)
		if err != nil {
			s.logger.Errorw("build user failed",
				"error", err,
				"operation", "create_register",
				"email", userInput.Email,
				"personal_data_id", personalData.Id,
				"team_id", teamID,
			)
			return err
		}

		createdUser, err := s.userRepo.Create(ctx, *mappedUser)
		if err != nil {
			if errors.Is(err, user.ErrAlreadyExists) {
				return ErrUserAlreadyExists
			}
			s.logger.Errorw("create user repo failed",
				"error", err,
				"operation", "create_register",
				"email", userInput.Email,
				"user_id", mappedUser.ID,
				"team_id", mappedUser.TeamID,
				"personal_data_id", mappedUser.PersonalDataID,
			)
			return ErrCreateUserFailed
		}

		u = createdUser
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "create_register",
			"email", userInput.Email,
			"team_name", userInput.TeamName,
		)
	}

	logSuccess(s.logger,
		"operation", "create_register",
		"user_id", u.ID,
		"email", u.Email,
		"role", u.Role,
		"team_id", u.TeamID,
	)

	return u, nil
}

func (s *service) CreateByActor(ctx context.Context, actor auth.Actor, userInput CreateUserInput) (*User, error) {
	var u *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		var teamID uuid.UUID = uuid.Nil
		actorRole := valueobjects.Role(actor.Role)
		if !actorRole.IsManagerRole() {
			s.logger.Errorw("actor role cannot create users",
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
			)
			return ErrOnlyManagersCanCreate
		}
		if userInput.TeamName != nil {
			team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
			if err != nil {
				s.logger.Errorw("get team by name failed",
					"error", err,
					"operation", "create_by_actor",
					"actor_id", actor.Id,
					"actor_role", actor.Role,
					"team_name", *userInput.TeamName,
				)
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
			s.logger.Errorw("build personal data failed",
				"error", err,
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"email", userInput.Email,
			)
			return err
		}

		if _, err = s.dataRepo.Create(ctx, *personalData); err != nil {
			s.logger.Errorw("create personal data failed",
				"error", err,
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"personal_data_id", personalData.Id,
			)
			return ErrPersonalDataCreateFailed
		}

		if userInput.Role == nil {
			s.logger.Errorw("role is required",
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"email", userInput.Email,
			)
			return ErrRoleRequired
		}

		mappedUser, err := user.New(
			teamID,
			personalData.Id,
			userInput.Email,
			userInput.Password,
			*userInput.Role,
		)
		if err != nil {
			s.logger.Errorw("build user failed",
				"error", err,
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"email", userInput.Email,
				"personal_data_id", personalData.Id,
				"team_id", teamID,
			)
			return err
		}

		createdUser, err := s.userRepo.Create(ctx, *mappedUser)
		if err != nil {
			s.logger.Errorw("create user repo failed",
				"error", err,
				"operation", "create_by_actor",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"email", userInput.Email,
				"user_id", mappedUser.ID,
				"team_id", mappedUser.TeamID,
				"personal_data_id", mappedUser.PersonalDataID,
			)
			return ErrUserCreateFailed
		}

		u = createdUser
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "create_by_actor",
			"actor_id", actor.Id,
			"actor_role", actor.Role,
			"email", userInput.Email,
			"team_name", userInput.TeamName,
		)
	}

	logSuccess(s.logger,
		"operation", "create_by_actor",
		"actor_id", actor.Id,
		"actor_role", actor.Role,
		"user_id", u.ID,
		"email", u.Email,
		"role", u.Role,
		"team_id", u.TeamID,
	)

	return u, nil
}

func (s *service) Update(ctx context.Context, actor auth.Actor, userInput UpdateUserInput) (*User, error) {
	var updatedUser *User

	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		if userInput.Email == nil {
			s.logger.Errorw("email is required",
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", userInput.UserID,
			)
			return ErrInvalidInput
		}
		if userInput.Password == nil {
			s.logger.Errorw("password is required",
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", userInput.UserID,
				"email", *userInput.Email,
			)
			return ErrInvalidInput
		}
		if userInput.Role == nil {
			s.logger.Errorw("role is required",
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", userInput.UserID,
				"email", *userInput.Email,
			)
			return ErrRoleRequired
		}

		existingUser, err := s.userRepo.GetByEmail(ctx, *userInput.Email)
		if err != nil {
			s.logger.Errorw("get user by email failed",
				"error", err,
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"email", *userInput.Email,
			)
			return ErrUserNotFound
		}

		if userInput.TeamId != nil {
			team, err := s.teamRepo.GetByID(ctx, *userInput.TeamId)
			if err != nil {
				s.logger.Errorw("get team by id failed",
					"error", err,
					"operation", "update",
					"actor_id", actor.Id,
					"actor_role", actor.Role,
					"team_id", *userInput.TeamId,
					"user_id", existingUser.ID,
				)
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
				s.logger.Errorw("get personal data failed",
					"error", err,
					"operation", "update",
					"actor_id", actor.Id,
					"actor_role", actor.Role,
					"user_id", existingUser.ID,
					"personal_data_id", existingUser.PersonalDataID,
				)
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
				s.logger.Errorw("validate personal data failed",
					"error", err,
					"operation", "update",
					"actor_id", actor.Id,
					"actor_role", actor.Role,
					"user_id", existingUser.ID,
					"personal_data_id", pd.Id,
				)
				return err
			}

			if _, err := s.dataRepo.Update(ctx, pd); err != nil {
				s.logger.Errorw("update personal data failed",
					"error", err,
					"operation", "update",
					"actor_id", actor.Id,
					"actor_role", actor.Role,
					"user_id", existingUser.ID,
					"personal_data_id", pd.Id,
				)
				return ErrPersonalDataUpdateFailed
			}
		}

		email, err := valueobjects.NewEmail(*userInput.Email)
		if err != nil {
			s.logger.Errorw("build email failed",
				"error", err,
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"email", *userInput.Email,
			)
			return err
		}
		existingUser.Email = email

		password, err := valueobjects.NewPassword(*userInput.Password)
		if err != nil {
			s.logger.Errorw("build password failed",
				"error", err,
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
			)
			return err
		}
		existingUser.Password = password

		if valueobjects.IsValidRole(*userInput.Role) {
			existingUser.Role = valueobjects.Role(*userInput.Role)
		} else {
			s.logger.Errorw("invalid role",
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"role", *userInput.Role,
			)
			return ErrInvalidRole
		}

		savedUser, err := s.userRepo.Update(ctx, existingUser.ID, *existingUser)
		if err != nil {
			s.logger.Errorw("update user repo failed",
				"error", err,
				"operation", "update",
				"actor_id", actor.Id,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"email", existingUser.Email,
				"team_id", existingUser.TeamID,
				"personal_data_id", existingUser.PersonalDataID,
			)
			return ErrUserUpdateFailed
		}

		updatedUser = savedUser
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "update",
			"actor_id", actor.Id,
			"actor_role", actor.Role,
			"user_id", userInput.UserID,
			"email", userInput.Email,
		)
	}

	logSuccess(s.logger,
		"operation", "update",
		"actor_id", actor.Id,
		"actor_role", actor.Role,
		"user_id", updatedUser.ID,
		"email", updatedUser.Email,
		"role", updatedUser.Role,
		"team_id", updatedUser.TeamID,
	)

	return updatedUser, nil
}
