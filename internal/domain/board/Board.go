package board

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id        uuid.UUID
	TeamId    uuid.UUID
	Name      string
	IsPublic  bool
	Status    BoardStatus
	CreatedAt time.Time
}

func New(teamId uuid.UUID, isPublic bool, name string) (*Board, error) {
	if name == "" {
		return nil, ErrEmptyName
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
