package personaldata

import (
	"time"
	"unicode"

	"github.com/google/uuid"
)

var minBirthDate = time.Date(1800, time.January, 1, 0, 0, 0, 0, time.UTC)

type PersonalData struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Age       *uint8
	BirthDate *time.Time
}

func New(firstName, lastName string, birthDate *time.Time, age *uint8) (*PersonalData, error) {
	if err := validate(firstName, lastName, birthDate); err != nil {
		return nil, err
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
	return validate(data.FirstName, data.LastName, data.BirthDate)
}

func validate(firstName, lastName string, birthDate *time.Time) error {
	if firstName == "" {
		return ErrFirstNameRequired
	}
	if lastName == "" {
		return ErrLastNameRequired
	}
	if hasDigit(firstName) {
		return ErrInvalidFirstName
	}
	if hasDigit(lastName) {
		return ErrInvalidLastName
	}

	if birthDate != nil && birthDate.After(time.Now()) {
		return ErrInvalidBirthDate
	}
	if birthDate != nil && birthDate.Before(minBirthDate) {
		return ErrBirthDateTooOld
	}
	return nil
}

func hasDigit(value string) bool {
	for _, r := range value {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
