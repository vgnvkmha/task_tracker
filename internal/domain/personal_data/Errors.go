package personaldata

import "errors"

var (
	// general
	ErrNotFound      = errors.New("personal data not found")
	ErrAlreadyExists = errors.New("personal data already exists")

	// validation
	ErrFirstNameRequired = errors.New("first name must be provided")
	ErrLastNameRequired  = errors.New("last name must be provided")
	ErrInvalidFirstName  = errors.New("first name must not contain digits")
	ErrInvalidLastName   = errors.New("last name must not contain digits")

	ErrInvalidBirthDate = errors.New("birth date cannot be in the future")
	ErrBirthDateTooOld  = errors.New("birth date cannot be before 1800")
	ErrBirthDateNotSet  = errors.New("birth date is not set")

	ErrNegativeAge = errors.New("user's age must be positive")

	// logical
	ErrTooYoung   = errors.New("user does not meet minimum age requirement")
	ErrInvalidAge = errors.New("invalid age calculated from birth date")
)
