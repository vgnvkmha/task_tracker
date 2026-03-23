package dto

import "time"

type TaskResponse struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeId  uint32    `json:"assignee_id"`
	ReporterId  uint32    `json:"reporter_id"`
	BoardId     uint32    `json:"board_id"`
	DueTo       time.Time `json:"due_to"`
}
