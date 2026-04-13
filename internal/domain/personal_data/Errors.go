package personaldata

import "errors"

var (
	// general
	ErrNotFound      = errors.New("personal data not found")
	ErrAlreadyExists = errors.New("personal data already exists")

	// validation
	ErrFirstNameRequired = errors.New("first name must be provided")
	ErrLastNameRequired  = errors.New("last name must be provided")

	ErrInvalidBirthDate = errors.New("birth date cannot be in the future")
	ErrBirthDateNotSet  = errors.New("birth date is not set")

	ErrNegativeAge = errors.New("user's age must be positive")

	// logical
	ErrTooYoung   = errors.New("user does not meet minimum age requirement")
	ErrInvalidAge = errors.New("invalid age calculated from birth date")
)
