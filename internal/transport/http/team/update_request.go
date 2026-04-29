package team

import (
	"task_tracker/internal/application/team"
	"task_tracker/internal/common_errors"

	"github.com/google/uuid"
)

type UpdateTeamRequest struct {
	Name     *string `json:"name"`
	Timezone *string `json:"timezone"`
	LeaderID *string `json:"leader_id"`
}

type Team = team.UpdateTeamInput

func ApplyUpdateTeam(req UpdateTeamRequest) (*Team, error) {
	team := Team{} //TODO: delete questionable
	if req.Name != nil {
		team.Name = req.Name
	}
	if req.Timezone != nil {
		team.Timezone = req.Timezone
	}
	if req.LeaderID != nil {
		uuid, err := uuid.Parse(*req.LeaderID)
		if err != nil {
			return nil, common_errors.ErrInvalidID
		}
		team.LeaderID = &uuid
	}
	return &team, nil
}
