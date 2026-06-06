package task_handler

import (
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h TaskHandler, tokens *infraauth.JWTService, legacyHeadersEnabled bool) {
	task := r.Group("/task", middleware.ActorMiddleware(tokens, legacyHeadersEnabled))

	task.POST("/create", h.Create)
	task.GET("/get_by_id/:id", h.GetByID)
	task.GET("/list", h.List)
	task.GET("/get_by_board/:board_id", h.GetByBoardID)
	task.GET("/get_by_assignee/:assignee_id", h.GetByAssigneeID)
	task.PATCH("/update/:id", h.Update)
	task.DELETE("/delete/:id", h.Delete)
}
