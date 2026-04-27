package team

import "github.com/google/uuid"

type CreateTeamInput struct {
	Name     string
	Timezone *string
	LeaderID *uuid.UUID
}
