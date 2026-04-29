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

func NewApplicationTeam(req CreateTeamRequest) (*applicationTeam, error) {
	var leaderID *uuid.UUID

	if req.LeaderID != nil {
		id, err := uuid.Parse(*req.LeaderID)
		if err != nil {
			return nil, common_errors.ErrInvalidID
		}
		leaderID = &id
	}

	return &applicationTeam{
		Name:     req.Name,
		Timezone: req.Timezone,
		LeaderID: leaderID,
	}, nil
}
