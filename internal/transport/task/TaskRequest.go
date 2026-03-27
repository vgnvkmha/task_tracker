package dto

import "time"

type TaskRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BoardID     uint32    `json:"board_id"`
	AssigneeID  uint32    `json:"assignee_id"`
	ReporetID   uint32    `json:"reporter_id"`
	DueTo       time.Time `json:"due_to"`
}
