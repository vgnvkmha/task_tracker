package team

import (
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h TeamHandler, tokens *infraauth.JWTService, legacyHeadersEnabled bool) {
	team := r.Group("/team", middleware.ActorMiddleware(tokens, legacyHeadersEnabled))

	team.POST("/create", h.Create)

	team.GET("/get_by_id/:id", h.GetByID)
	team.GET("/get_by_name/:team_name", h.GetByName)
	team.GET("/list_active", h.ListActive)
	team.GET("/list", h.List)

	team.PATCH("/update/:id", h.Update)

	team.DELETE("/delete_by_id/:id", h.DeleteByID)
}
