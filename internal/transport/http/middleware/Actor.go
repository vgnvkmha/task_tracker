package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"
	infraauth "task_tracker/internal/infrastracture/auth"
)

const actorKey = "actor" //TODO: remove hadrcode

type TokenParser interface {
	ParseAccessToken(token string) (auth.TokenClaims, error)
}

func ActorMiddleware(tokenParser TokenParser, legacyHeadersEnabled bool) gin.HandlerFunc {
	return actorMiddleware(tokenParser, legacyHeadersEnabled, abortUnauthorized)
}

func UIActorMiddleware(tokenParser TokenParser) gin.HandlerFunc {
	base := actorMiddleware(tokenParser, false, redirectToLogin)
	return func(ctx *gin.Context) {
		setNoStoreHeaders(ctx)
		base(ctx)
	}
}

func actorMiddleware(tokenParser TokenParser, legacyHeadersEnabled bool, unauthorized func(*gin.Context, string)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, hasAuthorization := bearerToken(ctx.GetHeader("Authorization"))
		if hasAuthorization {
			if token == "" || tokenParser == nil {
				unauthorized(ctx, "invalid token")
				return
			}
			claims, err := tokenParser.ParseAccessToken(token)
			if err != nil {
				unauthorized(ctx, authErrorMessage(err))
				return
			}

			ctx.Set(actorKey, claims.Actor())
			ctx.Next()
			return
		}

		if tokenParser != nil {
			if cookieToken, err := ctx.Cookie("access_token"); err == nil && cookieToken != "" {
				claims, err := tokenParser.ParseAccessToken(cookieToken)
				if err != nil {
					unauthorized(ctx, authErrorMessage(err))
					return
				}

				ctx.Set(actorKey, claims.Actor())
				ctx.Next()
				return
			}
		}

		if legacyHeadersEnabled {
			actor, ok := legacyActor(ctx, unauthorized)
			if !ok {
				return
			}
			log.Printf("WARN legacy header auth fallback used path=%s user_id=%s role=%s", ctx.FullPath(), actor.ID, actor.Role)
			ctx.Set(actorKey, actor)
			ctx.Next()
			return
		}

		unauthorized(ctx, "missing bearer token")
	}
}

func GetActor(ctx *gin.Context) (auth.Actor, bool) {
	val, exists := ctx.Get(actorKey)
	if !exists {
		return auth.Actor{}, false
	}

	actor, ok := val.(auth.Actor)
	return actor, ok
}

func MockActorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(actorKey, auth.Actor{
			ID:   uuid.MustParse("3f1c2a6e-9b7d-4c8f-8a2e-1d5b6f7a9c10"),
			Role: "admin",
		})
		c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", true
	}
	return parts[1], true
}

func legacyActor(ctx *gin.Context, unauthorized func(*gin.Context, string)) (auth.Actor, bool) {
	userIDStr := ctx.GetHeader("X-User-ID") //TODO: remove after JWT migration
	role := ctx.GetHeader("X-User-Role")
	if !valueobjects.IsValidRole(role) {
		unauthorized(ctx, "missing role")
		return auth.Actor{}, false
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		unauthorized(ctx, "invalid user")
		return auth.Actor{}, false
	}

	return auth.Actor{
		ID:   userID,
		Role: valueobjects.Role(role),
	}, true
}

func abortUnauthorized(ctx *gin.Context, message string) {
	clearAccessTokenCookie(ctx)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": message,
	})
}

func redirectToLogin(ctx *gin.Context, message string) {
	values := "?auth=required"
	if message == "expired token" {
		values = "?auth=expired"
	}
	clearAccessTokenCookie(ctx)
	setNoStoreHeaders(ctx)
	body := `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta http-equiv="Cache-Control" content="no-store">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Redirect</title>
</head>
<body>
  <script>window.location.replace("/ui/users/create` + values + `");</script>
  <noscript><a href="/ui/users/create` + values + `">Перейти на страницу входа</a></noscript>
</body>
</html>`
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(body))
	ctx.Abort()
}

func setNoStoreHeaders(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
}

func clearAccessTokenCookie(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func authErrorMessage(err error) string {
	switch {
	case errors.Is(err, infraauth.ErrExpiredToken):
		return "expired token"
	case errors.Is(err, infraauth.ErrMissingToken):
		return "missing bearer token"
	default:
		return "invalid token"
	}
}
