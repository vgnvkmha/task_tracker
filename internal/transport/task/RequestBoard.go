package dto

type RequestBoard struct {
	TeamId   uint32 `json:"team_id"`
	IsPublic bool   `json:"is_public"`
	Name     string `json:"name"`
}
