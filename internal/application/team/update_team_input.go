package team

import "github.com/google/uuid"

type UpdateTeamInput struct {
	Name     *string
	Timezone *string
	LeaderID *uuid.UUID
	IsActive *bool
}
