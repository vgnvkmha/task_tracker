package dto

type TaskResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	DueTo       string `json:"due_to"`
	UpdatedAt   string `json:"updated_at"`
	ReporterId  string `json:"reporter_id"`
	AssigneeId  string `json:"assignee_id"`
	BoardId     string `json:"board_id"`
	SprintId    string `json:"sprint_id"`
}
