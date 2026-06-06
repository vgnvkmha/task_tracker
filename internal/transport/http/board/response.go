package board

import (
	boardapp "task_tracker/internal/application/board"
	"time"
)

type Response struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func NewResponse(board *boardapp.Board) Response {
	return Response{
		ID:        board.Id.String(),
		TeamID:    board.TeamId.String(),
		Name:      board.Name,
		IsPublic:  board.IsPublic,
		Status:    string(board.Status),
		CreatedAt: board.CreatedAt.Format(time.RFC3339),
	}
}

func NewResponses(boards []*boardapp.Board) []Response {
	response := make([]Response, 0, len(boards))
	for _, board := range boards {
		response = append(response, NewResponse(board))
	}
	return response
}
