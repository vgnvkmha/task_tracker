package user

import personaldata "task_tracker/internal/domain/personal_data"

type Profile struct {
	User         *User
	PersonalData personaldata.PersonalData
	TeamName     string
}
