package valueobjects

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email string

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return "", ErrInvalidEmail
	}

	_, err := mail.ParseAddress(raw)
	if err != nil {
		return "", ErrInvalidEmail
	}

	return Email(strings.ToLower(raw)), nil
}
