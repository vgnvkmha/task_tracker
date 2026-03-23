package models

import valueobjects "task_tracker/internal/domain/models/value_objects"

type PersonalData struct {
	email    string
	password string
	role     valueobjects.Role
}
