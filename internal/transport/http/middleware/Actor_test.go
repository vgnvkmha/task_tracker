package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainauth "task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"
	infraauth "task_tracker/internal/infrastracture/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeTokenParser struct {
	claims domainauth.TokenClaims
	err    error
}

func (p fakeTokenParser) ParseAccessToken(token string) (domainauth.TokenClaims, error) {
	if token == "valid-token" {
		return p.claims, p.err
	}
	return domainauth.TokenClaims{}, infraauth.ErrInvalidToken
}

func TestActorMiddlewareUsesJWTBeforeLegacyHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{
		claims: domainauth.TokenClaims{
			Subject: userID,
			Role:    valueobjects.User,
		},
	}, true), func(c *gin.Context) {
		actor, ok := GetActor(c)
		if !ok {
			t.Fatal("actor missing")
		}
		if actor.ID != userID || actor.Role != valueobjects.User {
			t.Fatalf("actor = %+v, want JWT actor", actor)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-ID", uuid.NewString())
	req.Header.Set("X-User-Role", string(valueobjects.Admin))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

func TestActorMiddlewareFallsBackToLegacyHeadersWhenEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{}, true), func(c *gin.Context) {
		actor, ok := GetActor(c)
		if !ok {
			t.Fatal("actor missing")
		}
		if actor.ID != userID || actor.Role != valueobjects.Admin {
			t.Fatalf("actor = %+v, want legacy actor", actor)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", userID.String())
	req.Header.Set("X-User-Role", string(valueobjects.Admin))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

func TestActorMiddlewareUsesAccessTokenCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{
		claims: domainauth.TokenClaims{
			Subject: userID,
			Role:    valueobjects.User,
		},
	}, true), func(c *gin.Context) {
		actor, ok := GetActor(c)
		if !ok {
			t.Fatal("actor missing")
		}
		if actor.ID != userID || actor.Role != valueobjects.User {
			t.Fatalf("actor = %+v, want cookie actor", actor)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "valid-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

func TestActorMiddlewareRejectsMissingTokenWhenLegacyDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{}, false), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestActorMiddlewareRejectsMalformedAuthorizationBeforeLegacyFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{}, true), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic abc")
	req.Header.Set("X-User-ID", uuid.NewString())
	req.Header.Set("X-User-Role", string(valueobjects.Admin))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "invalid token") {
		t.Fatalf("body = %s, want invalid token error", rec.Body.String())
	}
}

func TestActorMiddlewareRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", ActorMiddleware(fakeTokenParser{err: infraauth.ErrExpiredToken}, true), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "expired token") {
		t.Fatalf("body = %s, want expired token error", rec.Body.String())
	}
}

func TestUIActorMiddlewareRedirectsToLoginWhenTokenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ui/users/cabinet", UIActorMiddleware(fakeTokenParser{}), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ui/users/cabinet", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `window.location.replace("/ui/users/create?auth=required")`) {
		t.Fatalf("body = %s, want replace redirect to required login", rec.Body.String())
	}
	if cacheControl := rec.Header().Get("Cache-Control"); !strings.Contains(cacheControl, "no-store") {
		t.Fatalf("Cache-Control = %q, want no-store", cacheControl)
	}
}

func TestUIActorMiddlewareRedirectsToLoginWhenTokenExpired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ui/users/cabinet", UIActorMiddleware(fakeTokenParser{err: infraauth.ErrExpiredToken}), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ui/users/cabinet", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "valid-token"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `window.location.replace("/ui/users/create?auth=expired")`) {
		t.Fatalf("body = %s, want replace redirect to expired login", rec.Body.String())
	}
}
