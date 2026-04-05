package models

import valueobjects "task_tracker/internal/domain/models/value_objects"

type PersonalData struct {
	Email    string
	Password string
	Role     valueobjects.Role
}
