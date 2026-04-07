package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/models/value_objects"
)

const actorKey = "actor"

func ActorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.GetHeader("X-User-ID")
		role := ctx.GetHeader("X-User-Role")

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
