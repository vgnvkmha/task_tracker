package user

import (
	"context"
	"errors"
	"task_tracker/internal/application/common"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	personaldata "task_tracker/internal/domain/personal_data"
	"task_tracker/internal/domain/user"
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
	Login(ctx context.Context, email string, password string) (*User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (*Profile, error)

	ListActive(ctx context.Context) ([]*User, error)
	List(ctx context.Context) ([]*User, error)

	DeleteByID(ctx context.Context, actor auth.Actor, id uuid.UUID) error
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
		userRepo:    userRepo,
		dataRepo:    dataRepo,
		teamRepo:    teamRepo,
		logger:      logger,
		transaction: transaction,
	}
}

func (s *service) CreateRegister(ctx context.Context, userInput CreateUserInput) (*User, error) {
	var u *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		createdUser, err := s.createUser(ctx, "create_register", userInput)
		if err != nil {
			return err
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
		if !actor.Role.IsManagerRole() {
			logFailure(s.logger, "actor role cannot create users", nil,
				"operation", "create_by_actor",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
			)
			return ErrOnlyManagersCanModify
		}

		createdUser, err := s.createUser(ctx, "create_by_actor", userInput,
			"actor_id", actor.ID,
			"actor_role", actor.Role,
		)
		if err != nil {
			return err
		}
		u = createdUser
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "create_by_actor",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"email", userInput.Email,
			"team_name", userInput.TeamName,
		)
	}

	logSuccess(s.logger,
		"operation", "create_by_actor",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"user_id", u.ID,
		"email", u.Email,
		"role", u.Role,
		"team_id", u.TeamID,
	)

	return u, nil
}

func (s *service) createUser(ctx context.Context, operation string, userInput CreateUserInput, logFields ...any) (*User, error) {
	baseLogFields := appendLogFields([]any{"operation", operation}, logFields...)
	teamID := uuid.Nil

	if userInput.TeamName != nil {
		team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
		if err != nil {
			logFailure(s.logger, "get team by name failed", err,
				appendLogFields(baseLogFields, "team_name", *userInput.TeamName)...,
			)
			return nil, ErrTeamNotFound
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
		logFailure(s.logger, "build personal data failed", err,
			appendLogFields(baseLogFields, "email", userInput.Email)...,
		)
		return nil, err
	}

	mappedUser, err := user.New(
		teamID,
		personalData.Id,
		userInput.Email,
		userInput.Password,
		userInput.Role,
	)
	if err != nil {
		logFailure(s.logger, "build user failed", err,
			appendLogFields(baseLogFields,
				"email", userInput.Email,
				"personal_data_id", personalData.Id,
				"team_id", teamID,
			)...,
		)
		return nil, err
	}

	if _, err = s.dataRepo.Create(ctx, *personalData); err != nil {
		logFailure(s.logger, "create personal data failed", err,
			appendLogFields(baseLogFields, "personal_data_id", personalData.Id)...,
		)
		return nil, ErrPersonalDataCreateFailed
	}

	createdUser, err := s.userRepo.Create(ctx, *mappedUser)
	if err != nil {
		mappedErr := mapCreateError(err)
		logFailure(s.logger, "create user repo failed", err,
			appendLogFields(baseLogFields,
				"email", userInput.Email,
				"user_id", mappedUser.ID,
				"team_id", mappedUser.TeamID,
				"personal_data_id", mappedUser.PersonalDataID,
			)...,
		)
		if mappedErr != err {
			return nil, mappedErr
		}
		return nil, ErrCreateUserFailed
	}

	return createdUser, nil
}

func (s *service) Update(ctx context.Context, actor auth.Actor, userInput UpdateUserInput) (*User, error) {
	var updatedUser *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {

		if !actor.Role.IsManagerRole() {
			if actor.ID != userInput.UserID {
				return ErrOnlyManagersCanModify
			}
			if userInput.Role != nil || userInput.TeamId != nil || userInput.TeamName != nil {
				return ErrOnlyManagersCanModify
			}
		}

		existingUser, err := s.userRepo.GetByID(ctx, userInput.UserID)
		if err != nil {
			logFailure(s.logger, "get user by ID failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"email", userInput.UserID,
			)
			return ErrUserNotFound
		}
		pd, err := s.dataRepo.Get(ctx, existingUser.PersonalDataID)
		if err != nil {
			logFailure(s.logger, "get personal data failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"personal_data_id", existingUser.PersonalDataID,
			)
			return ErrPersonalDataNotFound
		}

		if err := pd.Update(userInput.FirstName, userInput.LastName, userInput.Age, userInput.BirthDate); err != nil {
			logFailure(s.logger, "validate personal data failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"personal_data_id", pd.Id,
			)
			return err
		}

		if _, err := s.dataRepo.Update(ctx, pd); err != nil {
			logFailure(s.logger, "update personal data failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"personal_data_id", pd.Id,
			)
			return ErrPersonalDataUpdateFailed
		}

		teamID := existingUser.TeamID
		if userInput.TeamId != nil {
			if _, err := s.teamRepo.GetByID(ctx, *userInput.TeamId); err != nil {
				logFailure(s.logger, "get team by ID failed", err,
					"operation", "update",
					"actor_id", actor.ID,
					"actor_role", actor.Role,
					"user_id", existingUser.ID,
					"team_id", *userInput.TeamId,
				)
				return ErrTeamNotFound
			}
			teamID = userInput.TeamId
		}

		if userInput.TeamName != nil {
			team, err := s.teamRepo.GetByName(ctx, *userInput.TeamName)
			if err != nil {
				logFailure(s.logger, "get team by name failed", err,
					"operation", "update",
					"actor_id", actor.ID,
					"actor_role", actor.Role,
					"user_id", existingUser.ID,
					"team_name", *userInput.TeamName,
				)
				return ErrTeamNotFound
			}
			teamID = &team.ID
		}

		if err := existingUser.Update(teamID, userInput.Email, userInput.Password, userInput.Role); err != nil {
			logFailure(s.logger, "update user domain failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"email", userInput.Email,
				"role", userInput.Role,
				"team_id", teamID,
			)
			return err
		}

		savedUser, err := s.userRepo.Update(ctx, existingUser.ID, *existingUser)
		if err != nil {
			mappedErr := mapUpdateError(err)
			logFailure(s.logger, "update user repo failed", err,
				"operation", "update",
				"actor_id", actor.ID,
				"actor_role", actor.Role,
				"user_id", existingUser.ID,
				"email", existingUser.Email,
				"team_id", existingUser.TeamID,
				"personal_data_id", existingUser.PersonalDataID,
			)
			if mappedErr != err {
				return mappedErr
			}
			return ErrUserUpdateFailed
		}

		updatedUser = savedUser
		return nil
	})

	if err != nil {
		return nil, logError(err, s.logger,
			"operation", "update",
			"actor_id", actor.ID,
			"actor_role", actor.Role,
			"user_id", userInput.UserID,
			"email", userInput.Email,
		)
	}

	logSuccess(s.logger,
		"operation", "update",
		"actor_id", actor.ID,
		"actor_role", actor.Role,
		"user_id", updatedUser.ID,
		"email", updatedUser.Email,
		"role", updatedUser.Role,
		"team_id", updatedUser.TeamID,
	)

	return updatedUser, nil
}

func (s *service) Login(ctx context.Context, email string, password string) (*User, error) {
	foundUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logFailure(s.logger, "get user by email failed", err,
			"operation", "login",
			"email", email,
		)
		return nil, ErrInvalidCredentials
	}

	if !foundUser.Password.Compare(password) {
		logFailure(s.logger, "password check failed", nil,
			"operation", "login",
			"email", email,
			"user_id", foundUser.ID,
		)
		return nil, ErrInvalidCredentials
	}

	logSuccess(s.logger,
		"operation", "login",
		"user_id", foundUser.ID,
		"email", foundUser.Email,
	)

	return foundUser, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var result *User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}
		result = user
		return nil
	})

	if err != nil {
		logFailure(s.logger, "get user by ID operation failed", err,
			"operation", "get_by_id",
			"id", id,
		)
		return nil, err
	}
	return result, nil
}

func (s *service) GetProfileByID(ctx context.Context, id uuid.UUID) (*Profile, error) {
	var result *Profile
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		foundUser, err := s.userRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}

		pd, err := s.dataRepo.Get(ctx, foundUser.PersonalDataID)
		if err != nil {
			return ErrPersonalDataNotFound
		}

		teamName := ""
		if foundUser.TeamID != nil {
			team, err := s.teamRepo.GetByID(ctx, *foundUser.TeamID)
			if err != nil {
				logFailure(s.logger, "get team by ID failed", err,
					"operation", "get_profile_by_id",
					"id", id,
					"team_id", *foundUser.TeamID,
				)
			} else if team != nil {
				teamName = team.Name
			}
		}

		result = &Profile{
			User:         foundUser,
			PersonalData: pd,
			TeamName:     teamName,
		}
		return nil
	})

	if err != nil {
		logFailure(s.logger, "get profile by ID operation failed", err,
			"operation", "get_profile_by_id",
			"id", id,
		)
		return nil, err
	}
	return result, nil
}

func (s *service) ListActive(ctx context.Context) ([]*User, error) {
	var result []*User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.ListActive(ctx)
		if err != nil {
			return mapGetError(err)
		}
		result = user
		return nil
	})

	if err != nil {
		logFailure(s.logger, "list active get operation failed", err,
			"operation", "list_active",
		)
		return nil, err
	}
	return result, nil
}

func (s *service) List(ctx context.Context) ([]*User, error) {
	var result []*User
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.ListActive(ctx)
		if err != nil {
			return mapGetError(err)
		}
		result = user
		return nil
	})

	if err != nil {
		logFailure(s.logger, "list get operation failed", err,
			"operation", "list",
		)
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteByID(ctx context.Context, actor auth.Actor, id uuid.UUID) error {
	err := s.transaction.WithTx(ctx, func(ctx context.Context) error {

		if !actor.Role.IsManagerRole() && actor.ID != id {
			return ErrOnlyManagersCanModify
		}

		_, err := s.userRepo.GetByID(ctx, id)
		if err != nil {
			return mapGetError(err)
		}

		err = s.userRepo.Delete(ctx, id)
		if err != nil {
			return mapDeleteError(err)
		}
		return nil
	})

	if err != nil {
		logFailure(s.logger, "delete user by ID operation failed", err,
			"operation", "delete_by_id",
			"actor_role", actor.Role,
			"actor_id", actor.ID,
			"id", id,
		)
		return err
	}
	return nil
}

// helpers

func mapCreateError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return ErrOnlyManagersCanModify
	case errors.Is(err, user.ErrAlreadyExists):
		return ErrUserAlreadyExists
	case errors.Is(err, common_errors.ErrAlreadyExists):
		return ErrUserAlreadyExists
	case errors.Is(err, common_errors.ErrConflict):
		return ErrConflict
	case errors.Is(err, common_errors.ErrInvalidID):
		return ErrInvalidUserID
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
		return ErrOnlyManagersCanModify

	case errors.Is(err, common_errors.ErrAlreadyExists):
		return ErrUserAlreadyExists

	case errors.Is(err, common_errors.ErrNotFound): //TODO: bad case
		return common_errors.ErrNotFound

	case errors.Is(err, common_errors.ErrConflict):
		return ErrConflict

	case errors.Is(err, common_errors.ErrInvalidArgument):
		return ErrInvalidInput

	default:
		return err
	}
}

func mapDeleteError(err error) error {
	switch {
	case errors.Is(err, common_errors.ErrPermissionDenied):
		return ErrOnlyManagersCanModify

	case errors.Is(err, common_errors.ErrNotFound):
		return ErrUserNotFound
	case errors.Is(err, common_errors.ErrConflict):
		return ErrUserAlreadyDeleted
	default:
		return err
	}
}
