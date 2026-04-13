package personaldata

import (
	"time"

	"github.com/google/uuid"
)

type PersonalData struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Age       uint8
	BirthDate *time.Time
}

func New(firstName, lastName string, birthDate *time.Time, age int8) (*PersonalData, error) {
	if firstName == "" {
		return nil, ErrFirstNameRequired
	}
	if lastName == "" {
		return nil, ErrLastNameRequired
	}

	if birthDate != nil && birthDate.After(time.Now()) {
		return nil, ErrInvalidBirthDate
	}

	if age < 0 {
		return nil, ErrInvalidAge
	}

	return &PersonalData{
		Id:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Age:       uint8(age),
	}, nil
}
