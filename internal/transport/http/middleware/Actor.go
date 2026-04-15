package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"
)

const actorKey = "actor" //TODO: remove hadrcode

func ActorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.GetHeader("X-User-ID") //TODO: remove after JWT integration
		role := ctx.GetHeader("X-User-Role")
		if !valueobjects.IsValidRole(role) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing role",
			})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			return
		}

		actor := auth.Actor{
			Id:   userID,
			Role: valueobjects.Role(role),
		}

		ctx.Set(actorKey, actor)
		ctx.Next()
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
			Id:   uuid.MustParse("3f1c2a6e-9b7d-4c8f-8a2e-1d5b6f7a9c10"),
			Role: "admin",
		})
		c.Next()
	}
}
