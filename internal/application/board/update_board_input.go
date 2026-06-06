package board

import (
	domainboard "task_tracker/internal/domain/board"

	"github.com/google/uuid"
)

type UpdateBoardInput struct {
	TeamID   *uuid.UUID
	Name     *string
	IsPublic *bool
	Status   *domainboard.BoardStatus
}
