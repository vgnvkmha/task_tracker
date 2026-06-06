package dto

import "task_tracker/internal/domain/task"

func ToTaskResponse(t task.Task) TaskResponse {
	assigneeID := ""
	if t.AssigneeId != nil {
		assigneeID = t.AssigneeId.String()
	}
	boardID := ""
	if t.BoardId != nil {
		boardID = t.BoardId.String()
	}
	sprintID := ""
	if t.SprintId != nil {
		sprintID = t.SprintId.String()
	}

	dueTo := ""
	if !t.DueTo.IsZero() {
		dueTo = t.DueTo.String()
	}

	return TaskResponse{
		Id:          t.Id.String(),
		Name:        t.Name,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt.String(),
		DueTo:       dueTo,
		UpdatedAt:   t.UpdatedAt.String(),
		ReporterId:  t.ReporterId.String(),
		AssigneeId:  assigneeID,
		BoardId:     boardID,
		SprintId:    sprintID,
	}
}

func ToTaskResponses(tasks []*task.Task) []TaskResponse {
	res := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		if t != nil {
			res = append(res, ToTaskResponse(*t))
		}
	}
	return res
}
