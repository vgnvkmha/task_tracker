package team

import (
	"task_tracker/internal/application/team"
)

type Response struct {
	Name     string  `json:"name"`
	Timezone *string `json:"timezone"`
	LeaderID *string `json:"leader_id"`
}

func FromDomain(team *team.Team) *Response {
	var (
		timezone *string
		leaderID *string
	)

	if team.Timezone != nil {
		tz := team.Timezone
		timezone = tz
	}

	if team.LeaderID != nil {
		id := team.LeaderID.String()
		leaderID = &id
	}

	return &Response{
		Name:     team.Name,
		Timezone: timezone,
		LeaderID: leaderID,
	}
}
