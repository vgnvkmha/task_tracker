package dto

import "github.com/google/uuid"

type UpdateUser struct {
	TeamId *uuid.UUID
	DataId *uuid.UUID
}
