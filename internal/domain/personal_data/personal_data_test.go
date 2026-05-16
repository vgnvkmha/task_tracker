package personaldata

import (
	"errors"
	"testing"
	"time"
)

func TestNewRejectsFirstNameWithDigits(t *testing.T) {
	_, err := New("Ivan1", "Petrov", nil, nil)
	if !errors.Is(err, ErrInvalidFirstName) {
		t.Fatalf("New error = %v, want ErrInvalidFirstName", err)
	}
}

func TestNewRejectsLastNameWithDigits(t *testing.T) {
	_, err := New("Ivan", "Petrov2", nil, nil)
	if !errors.Is(err, ErrInvalidLastName) {
		t.Fatalf("New error = %v, want ErrInvalidLastName", err)
	}
}

func TestNewRejectsBirthDateBefore1800(t *testing.T) {
	birthDate := time.Date(1799, time.December, 31, 0, 0, 0, 0, time.UTC)

	_, err := New("Ivan", "Petrov", &birthDate, nil)
	if !errors.Is(err, ErrBirthDateTooOld) {
		t.Fatalf("New error = %v, want ErrBirthDateTooOld", err)
	}
}

func TestValidateRejectsUpdatedNameWithDigits(t *testing.T) {
	data := PersonalData{
		FirstName: "Ivan",
		LastName:  "Petrov3",
	}

	err := data.Validate()
	if !errors.Is(err, ErrInvalidLastName) {
		t.Fatalf("Validate error = %v, want ErrInvalidLastName", err)
	}
}
