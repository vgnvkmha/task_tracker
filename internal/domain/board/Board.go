package board

import (
	valueobjects "task_tracker/internal/domain/value_objects"
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id        uuid.UUID
	TeamId    uuid.UUID
	IsPublic  bool
	Status    valueobjects.Status
	Name      string
	CreatedAt time.Time
}

func New(teamId uuid.UUID, isPublic bool, name string) (*Board, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	return &Board{
		Id:        uuid.New(),
		TeamId:    teamId,
		IsPublic:  isPublic,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}
