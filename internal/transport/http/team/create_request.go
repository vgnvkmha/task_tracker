package team

import (
	"task_tracker/internal/application/team"
	"task_tracker/internal/common_errors"

	"github.com/google/uuid"
)

type CreateTeamRequest struct {
	Name     string  `json:"name" binding:"required"`
	Timezone *string `json:"timezone"`
	LeaderID *string `json:"leader_id"`
}
type applicationTeam = team.CreateTeamInput

func ToServiceInput(input CreateTeamRequest) (*applicationTeam, error) {
	var leaderID *uuid.UUID

	if input.LeaderID != nil {
		id, err := uuid.Parse(*input.LeaderID)
		if err != nil {
			return nil, common_errors.ErrInvalidID
		}
		leaderID = &id
	}

	return &applicationTeam{
		Name:     input.Name,
		Timezone: input.Timezone,
		LeaderID: leaderID,
	}, nil
}
