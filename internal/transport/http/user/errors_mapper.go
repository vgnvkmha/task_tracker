package user

import (
	"errors"
	"net/http"
	userApplication "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	personaldata "task_tracker/internal/domain/personal_data"
	domainUser "task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"
)

func mapError(err error) (int, string) {
	switch {

	case errors.Is(err, userApplication.ErrInvalidUserID):
		return http.StatusUnprocessableEntity, "invalid user_id"

	case errors.Is(err, userApplication.ErrInvalidTeamID):
		return http.StatusUnprocessableEntity, "invalid team_id"

	case errors.Is(err, userApplication.ErrInvalidBirthDate):
		return http.StatusUnprocessableEntity, "invalid birth_date"

	case errors.Is(err, personaldata.ErrBirthDateTooOld),
		errors.Is(err, personaldata.ErrInvalidBirthDate):
		return http.StatusUnprocessableEntity, "invalid birth_date"

	case errors.Is(err, personaldata.ErrInvalidFirstName):
		return http.StatusUnprocessableEntity, "invalid first_name"

	case errors.Is(err, personaldata.ErrInvalidLastName):
		return http.StatusUnprocessableEntity, "invalid last_name"

	case errors.Is(err, valueobjects.ErrInvalidEmail):
		return http.StatusUnprocessableEntity, "invalid email"

	case errors.Is(err, valueobjects.ErrWeakPassword):
		return http.StatusUnprocessableEntity,
			"password is too weak: password must be at least 8 characters long and include lowercase, uppercase, digit, and special character"

	case errors.Is(err, domainUser.ErrManagerMustHaveTeam):
		return http.StatusUnprocessableEntity, "manager must have a team"

	case errors.Is(err, userApplication.ErrUserNotFound),
		errors.Is(err, common_errors.ErrNotFound),
		errors.Is(err, domainUser.ErrNotFound):
		return http.StatusNotFound, "user not found"

	case errors.Is(err, userApplication.ErrPersonalDataNotFound):
		return http.StatusNotFound, "personal data not found"

	case errors.Is(err, userApplication.ErrTeamNotFound):
		return http.StatusNotFound, "team not found"

	case errors.Is(err, domainUser.ErrEmailAlreadyUsed):
		return http.StatusConflict, "email already used"

	case errors.Is(err, userApplication.ErrUserAlreadyExists),
		errors.Is(err, common_errors.ErrAlreadyExists),
		errors.Is(err, domainUser.ErrAlreadyExists):
		return http.StatusConflict, "user with this email already exists"

	case errors.Is(err, userApplication.ErrUserAlreadyDeleted):
		return http.StatusGone, "user already deleted"

	case errors.Is(err, userApplication.ErrOnlyManagersCanModify),
		errors.Is(err, common_errors.ErrPermissionDenied),
		errors.Is(err, domainUser.ErrPermissionDenied):
		return http.StatusForbidden, "forbidden"

	case errors.Is(err, userApplication.ErrInvalidRole),
		errors.Is(err, userApplication.ErrInvalidInput),
		errors.Is(err, userApplication.ErrRoleRequired),
		errors.Is(err, domainUser.ErrInvalidRole):
		return http.StatusBadRequest, "invalid user input"

	case errors.Is(err, domainUser.ErrEmptyData):
		return http.StatusBadRequest, "personal data must be set"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
