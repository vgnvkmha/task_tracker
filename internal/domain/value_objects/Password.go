package valueobjects

import (
	"database/sql/driver"
	"errors"
	"fmt"

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
