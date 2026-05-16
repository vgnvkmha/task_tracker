package user

import (
	"time"

	userApplication "task_tracker/internal/application/user"
)

type CabinetView struct {
	ID        string
	Email     string
	Role      string
	TeamID    string
	FirstName string
	LastName  string
	Age       string
	BirthDate string
	IsActive  bool
}

func NewCabinetView(profile *userApplication.Profile) CabinetView {
	view := CabinetView{
		ID:        profile.User.ID.String(),
		Email:     string(profile.User.Email),
		Role:      string(profile.User.Role),
		FirstName: profile.PersonalData.FirstName,
		LastName:  profile.PersonalData.LastName,
		IsActive:  profile.User.IsActive,
	}

	if profile.User.TeamID != nil {
		view.TeamID = profile.User.TeamID.String()
	}
	if profile.PersonalData.Age != nil {
		view.Age = uint8ToString(*profile.PersonalData.Age)
	}
	if profile.PersonalData.BirthDate != nil {
		view.BirthDate = profile.PersonalData.BirthDate.Format(time.DateOnly)
	}

	return view
}

func uint8ToString(value uint8) string {
	if value == 0 {
		return "0"
	}

	var digits [3]byte
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = '0' + value%10
		value /= 10
	}
	return string(digits[i:])
}
