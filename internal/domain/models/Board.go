package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id        uuid.UUID
	TeamId    uuid.UUID
	IsPublic  bool
	Name      string
	CreatedAt time.Time
}
