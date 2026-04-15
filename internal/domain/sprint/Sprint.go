package sprint

import (
	"time"

	uuid "github.com/google/uuid"
)

type Sprint struct {
	Id        uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	Status    SprintStatus
	BoardId   uuid.UUID
}

func New(name string, startDate, endDate time.Time, boardId uuid.UUID) (*Sprint, error) {
	if name == "" {

	}
	if startDate.After(endDate) {

	}
	return &Sprint{
		Id:        uuid.New(),
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    SprintDraft,
		BoardId:   boardId,
	}, nil
}
