package board

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const maxBoardNameLength = 255

type Board struct {
	Id        uuid.UUID
	TeamId    uuid.UUID
	Name      string
	IsPublic  bool
	Status    BoardStatus
	CreatedAt time.Time
}

func New(teamId uuid.UUID, isPublic bool, name string) (*Board, error) {
	name = strings.TrimSpace(name)
	if teamId == uuid.Nil {
		return nil, ErrOwnerRequired
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if len(name) > maxBoardNameLength {
		return nil, ErrNameTooLong
	}

	return &Board{
		Id:        uuid.New(),
		TeamId:    teamId,
		Name:      name,
		IsPublic:  isPublic,
		Status:    BoardActive,
		CreatedAt: time.Now(),
	}, nil
}

func (b *Board) ApplyChanges(teamId *uuid.UUID, name *string, isPublic *bool, status *BoardStatus) error {
	if teamId != nil {
		if *teamId == uuid.Nil {
			return ErrOwnerRequired
		}
		b.TeamId = *teamId
	}

	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			return ErrEmptyName
		}
		if len(trimmedName) > maxBoardNameLength {
			return ErrNameTooLong
		}
		b.Name = trimmedName
	}

	if isPublic != nil {
		b.IsPublic = *isPublic
	}

	if status != nil {
		if !status.IsValid() {
			return ErrInvalidStatus
		}
		b.Status = *status
	}

	return nil
}
