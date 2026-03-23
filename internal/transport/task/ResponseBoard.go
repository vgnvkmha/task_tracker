package dto

import "time"

type ResponseBoard struct {
	Id        uint32    `json:"id"`
	TeamId    uint32    `json:"team_id"`
	IsPublic  bool      `json:"is_public"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
