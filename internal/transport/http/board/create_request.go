package board

import (
	boardapp "task_tracker/internal/application/board"
	"task_tracker/internal/common_errors"

	"github.com/google/uuid"
)

type CreateBoardRequest struct {
	TeamID   string `json:"team_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsPublic bool   `json:"is_public"`
}

func (r CreateBoardRequest) ToApplicationInput() (boardapp.CreateBoardInput, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return boardapp.CreateBoardInput{}, common_errors.ErrInvalidID
	}

	return boardapp.CreateBoardInput{
		TeamID:   teamID,
		Name:     r.Name,
		IsPublic: r.IsPublic,
	}, nil
}
