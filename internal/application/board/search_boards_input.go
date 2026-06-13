package board

import "github.com/google/uuid"

type SearchBoardsInput struct {
	Query  *string
	TeamID *uuid.UUID
	UserID *uuid.UUID
}
