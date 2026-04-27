package team

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID       uuid.UUID
	Name     string
	Timezone string
	LeaderID uuid.UUID
	IsActive bool
}

func New(name string, timezone *string, leaderID *uuid.UUID) (*Team, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	if timezone != nil && *timezone != "" {
		if _, err := time.LoadLocation(*timezone); err != nil {
			return nil, ErrInvalidTZ
		}
	}

	return &Team{
		ID:       uuid.New(),
		Name:     name,
		Timezone: timezone,
		LeaderID: leaderID,
		IsActive: true,
	}, nil
}
