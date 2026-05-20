package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	domainauth "task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

const (
	defaultIssuer   = "task_tracker"
	defaultAudience = "task_tracker_api"
)

type JWTService struct {
	secret   []byte
	ttl      time.Duration
	issuer   string
	audience string
	now      func() time.Time
}

func NewJWTService(secret string, ttl time.Duration) (*JWTService, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("jwt secret must be at least 32 bytes")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("jwt ttl must be positive")
	}

	return &JWTService{
		secret:   []byte(secret),
		ttl:      ttl,
		issuer:   defaultIssuer,
		audience: defaultAudience,
		now:      time.Now,
	}, nil
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID, role valueobjects.Role, teamID *uuid.UUID) (string, domainauth.TokenClaims, error) {
	now := s.now().UTC()
	claims := domainauth.TokenClaims{
		Subject:   userID,
		Role:      role,
		TeamID:    teamID,
		Issuer:    s.issuer,
		Audience:  s.audience,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(s.ttl).Unix(),
		TokenID:   newTokenID(),
	}

	token, err := s.sign(claims)
	return token, claims, err
}

func (s *JWTService) ParseAccessToken(token string) (domainauth.TokenClaims, error) {
	if strings.TrimSpace(token) == "" {
		return domainauth.TokenClaims{}, ErrMissingToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expected := s.signature(signingInput)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := decodeSegment(parts[0], &header); err != nil {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}

	var claims domainauth.TokenClaims
	if err := decodeSegment(parts[1], &claims); err != nil {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}
	if claims.Issuer != s.issuer || claims.Audience != s.audience {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}
	if claims.Subject == uuid.Nil || !valueobjects.IsValidRole(string(claims.Role)) {
		return domainauth.TokenClaims{}, ErrInvalidToken
	}
	if s.now().UTC().Unix() >= claims.ExpiresAt {
		return domainauth.TokenClaims{}, ErrExpiredToken
	}

	return claims, nil
}

func (s *JWTService) sign(claims domainauth.TokenClaims) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerSegment, err := encodeSegment(header)
	if err != nil {
		return "", err
	}
	claimsSegment, err := encodeSegment(claims)
	if err != nil {
		return "", err
	}

	signingInput := headerSegment + "." + claimsSegment
	return signingInput + "." + s.signature(signingInput), nil
}

func (s *JWTService) signature(signingInput string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeSegment(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func decodeSegment(segment string, dest any) error {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, dest)
}

func newTokenID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return uuid.NewString()
	}
	return base64.RawURLEncoding.EncodeToString(bytes[:])
}
