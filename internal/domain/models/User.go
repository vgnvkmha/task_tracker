package models

import valueobjects "task_tracker/internal/domain/models/value_objects"

type User struct {
	id   int
	team Team
	data valueobjects.PersonalData
}
