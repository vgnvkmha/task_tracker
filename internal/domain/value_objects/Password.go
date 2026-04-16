package valueobjects

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWeakPassword = errors.New("password is too weak")
)

type Password struct {
	hash []byte
}

func NewPassword(raw string) (Password, error) {
	if len(raw) < 8 {
		return Password{}, ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: hash}, nil
}

// For login
func (p Password) Compare(raw string) bool {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(raw)) == nil
}

// For DB save
func (p Password) Hash() string {
	return string(p.hash)
}
