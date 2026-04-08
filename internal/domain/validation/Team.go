package validation

import (
	task_errors "task_tracker/internal/domain/errors"
	valueobjects "task_tracker/internal/domain/models/value_objects"

	"github.com/google/uuid"
)

func IsAlloweToSeeTeamData(role valueobjects.Role, usersTeamid, teamId uuid.UUID) error {
	if usersTeamid == teamId || role.IsManagerRole() {
		return nil
	}
	return task_errors.ErrUnableToSeeTeamData
}
