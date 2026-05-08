package valueobjects

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWeakPassword = errors.New("password is too weak")
)

type Password struct {
	hash []byte
}

func NewPassword(raw string) (Password, error) {
	if !isStrongPassword(raw) {
		return Password{}, ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: hash}, nil
}

func isStrongPassword(raw string) bool {
	if len(raw) < 8 {
		return false
	}

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, r := range raw {
		switch {
		case unicode.IsLower(r):
			hasLower = true

		case unicode.IsUpper(r):
			hasUpper = true

		case unicode.IsDigit(r):
			hasDigit = true

		case unicode.IsPunct(r), unicode.IsSymbol(r):
			hasSpecial = true
		}

		if hasLower && hasUpper && hasDigit && hasSpecial {
			return true
		}
	}

	return false
}

// For login
func (p Password) Compare(raw string) bool {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(raw)) == nil
}

// For DB save
func (p Password) Hash() string {
	return string(p.hash)
}

func (p Password) Value() (driver.Value, error) {
	return p.Hash(), nil
}

func (p *Password) Scan(value any) error {
	switch v := value.(type) {
	case string:
		p.hash = []byte(v)
		return nil
	case []byte:
		p.hash = append(p.hash[:0], v...)
		return nil
	case nil:
		p.hash = nil
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Password", value)
	}
}
