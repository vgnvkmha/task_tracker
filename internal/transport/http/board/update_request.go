package board

import (
	boardapp "task_tracker/internal/application/board"
	domainboard "task_tracker/internal/domain/board"

	"task_tracker/internal/common_errors"

	"github.com/google/uuid"
)

type UpdateBoardRequest struct {
	TeamID   *string `json:"team_id"`
	Name     *string `json:"name"`
	IsPublic *bool   `json:"is_public"`
	Status   *string `json:"status"`
}

func (r UpdateBoardRequest) ToApplicationInput() (boardapp.UpdateBoardInput, error) {
	var teamID *uuid.UUID
	if r.TeamID != nil {
		parsedTeamID, err := uuid.Parse(*r.TeamID)
		if err != nil {
			return boardapp.UpdateBoardInput{}, common_errors.ErrInvalidID
		}
		teamID = &parsedTeamID
	}

	var status *domainboard.BoardStatus
	if r.Status != nil {
		parsedStatus := domainboard.BoardStatus(*r.Status)
		status = &parsedStatus
	}

	return boardapp.UpdateBoardInput{
		TeamID:   teamID,
		Name:     r.Name,
		IsPublic: r.IsPublic,
		Status:   status,
	}, nil
}
