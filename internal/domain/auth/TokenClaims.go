package auth

import (
	"time"

	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type TokenClaims struct {
	Subject   uuid.UUID         `json:"sub"`
	Role      valueobjects.Role `json:"role"`
	TeamID    *uuid.UUID        `json:"team_id,omitempty"`
	Issuer    string            `json:"iss"`
	Audience  string            `json:"aud"`
	IssuedAt  int64             `json:"iat"`
	ExpiresAt int64             `json:"exp"`
	TokenID   string            `json:"jti"`
}

func (c TokenClaims) Actor() Actor {
	return Actor{
		ID:   c.Subject,
		Role: c.Role,
	}
}

func (c TokenClaims) ExpiresAtTime() time.Time {
	return time.Unix(c.ExpiresAt, 0)
}
