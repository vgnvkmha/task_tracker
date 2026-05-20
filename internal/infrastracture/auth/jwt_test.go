package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

const testSecret = "test-secret-with-at-least-thirty-two-bytes"

func TestJWTServiceGenerateAndParseAccessToken(t *testing.T) {
	service, err := NewJWTService(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("NewJWTService error = %v", err)
	}
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	userID := uuid.New()
	teamID := uuid.New()
	token, claims, err := service.GenerateAccessToken(userID, valueobjects.Admin, &teamID)
	if err != nil {
		t.Fatalf("GenerateAccessToken error = %v", err)
	}

	parsed, err := service.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken error = %v", err)
	}

	if parsed.Subject != userID {
		t.Fatalf("subject = %s, want %s", parsed.Subject, userID)
	}
	if parsed.Role != valueobjects.Admin {
		t.Fatalf("role = %s, want %s", parsed.Role, valueobjects.Admin)
	}
	if parsed.TeamID == nil || *parsed.TeamID != teamID {
		t.Fatalf("team_id = %v, want %s", parsed.TeamID, teamID)
	}
	if parsed.ExpiresAt != claims.ExpiresAt {
		t.Fatalf("expires_at = %d, want %d", parsed.ExpiresAt, claims.ExpiresAt)
	}
}

func TestJWTServiceRejectsTamperedToken(t *testing.T) {
	service, err := NewJWTService(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("NewJWTService error = %v", err)
	}

	token, _, err := service.GenerateAccessToken(uuid.New(), valueobjects.User, nil)
	if err != nil {
		t.Fatalf("GenerateAccessToken error = %v", err)
	}

	tampered := token[:len(token)-1] + "x"
	if _, err := service.ParseAccessToken(tampered); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("ParseAccessToken error = %v, want ErrInvalidToken", err)
	}
}

func TestJWTServiceRejectsExpiredToken(t *testing.T) {
	service, err := NewJWTService(testSecret, time.Minute)
	if err != nil {
		t.Fatalf("NewJWTService error = %v", err)
	}
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	token, _, err := service.GenerateAccessToken(uuid.New(), valueobjects.User, nil)
	if err != nil {
		t.Fatalf("GenerateAccessToken error = %v", err)
	}

	service.now = func() time.Time { return now.Add(time.Minute) }
	if _, err := service.ParseAccessToken(token); !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("ParseAccessToken error = %v, want ErrExpiredToken", err)
	}
}

func TestNewJWTServiceRequiresStrongSecret(t *testing.T) {
	if _, err := NewJWTService("short", 15*time.Minute); err == nil || !strings.Contains(err.Error(), "secret") {
		t.Fatalf("NewJWTService error = %v, want secret length error", err)
	}
}
