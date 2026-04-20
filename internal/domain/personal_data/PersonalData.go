package personaldata

import (
	"time"

	"github.com/google/uuid"
)

type PersonalData struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Age       *uint8
	BirthDate *time.Time
}

func New(firstName, lastName string, birthDate *time.Time, age *uint8) (*PersonalData, error) {
	if firstName == "" {
		return nil, ErrFirstNameRequired
	}
	if lastName == "" {
		return nil, ErrLastNameRequired
	}

	if birthDate != nil && birthDate.After(time.Now()) {
		return nil, ErrInvalidBirthDate
	}

	return &PersonalData{
		Id:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Age:       age,
	}, nil
}

func (data *PersonalData) Validate() error {
	if data.FirstName == "" {
		return ErrFirstNameRequired
	}
	if data.LastName == "" {
		return ErrLastNameRequired
	}
	if data.BirthDate != nil && data.BirthDate.After(time.Now()) {
		return ErrInvalidBirthDate
	}
	return nil
}
