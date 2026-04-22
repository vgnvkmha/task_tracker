package team

import "github.com/google/uuid"

type Team struct {
	ID       uuid.UUID
	Name     string
	Timezone string
	LeaderID uuid.UUID
	IsActive bool
}
