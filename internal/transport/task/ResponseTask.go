package dto

import "time"

type TaskResponse struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Assignee string    `json:"assignee"`
	Reporter string    `json:"reporter"`
	Board    string    `json:"board"`
	DueTo    time.Time `json:"due_to"`
}
