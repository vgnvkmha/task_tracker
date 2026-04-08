package models

import "github.com/google/uuid"

type User struct {
	Id     uuid.UUID
	TeamId uuid.UUID
	Data   PersonalData
}
