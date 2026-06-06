package board

import "github.com/google/uuid"

type CreateBoardInput struct {
	TeamID   uuid.UUID
	Name     string
	IsPublic bool
}
