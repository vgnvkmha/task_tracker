package dto

import "task_tracker/internal/domain/models"

func ToTaskResponse(t models.Task) TaskResponse {
	return TaskResponse{
		Id:          t.Id.String(),
		Name:        t.Name,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt.String(),
		DueTo:       t.DueTo.String(),
		UpdatedAt:   t.UpdatedAt.String(),
		ReporterId:  t.ReporterId.String(),
		AssigneeId:  t.AssigneeId.String(),
		BoardId:     t.BoardId.String(),
		SprintId:    t.SprintId.String(),
	}
}

func ToTaskResponses(tasks []models.Task) []TaskResponse {
	res := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		res = append(res, ToTaskResponse(t))
	}
	return res
}
