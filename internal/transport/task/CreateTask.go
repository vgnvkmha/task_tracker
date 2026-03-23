package dto

import "time"

type CreateTaskDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BoardID     int       `json:"board_id"`
	AssigneeID  int       `json:"assignee_id"`
	ReporetID   int       `json:"reporter_id"`
	DueTo       time.Time `json:"due_to"`
}
