package task_handler

import (
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h TaskHandler, tokens *infraauth.JWTService, legacyHeadersEnabled bool) {
	ui := r.Group("/ui/tasks", middleware.UIActorMiddleware(tokens))
	ui.GET("", h.ShowTasks)
	ui.POST("", h.CreateFromUI)
	ui.POST("/update", h.UpdateFromUI)
	ui.POST("/delete", h.DeleteFromUI)

	taskModal := r.Group("/tasks", middleware.UIActorMiddleware(tokens))
	taskModal.GET("/:task_id", h.ShowTaskModal)
	taskModal.GET("/:task_id/edit", h.ShowTaskEditModal)
	taskModal.PATCH("/:task_id", h.UpdateTaskModal)

	task := r.Group("/task", middleware.ActorMiddleware(tokens, legacyHeadersEnabled))

	task.POST("/create", h.Create)
	task.GET("/get_by_id/:id", h.GetByID)
	task.GET("/list", h.List)
	task.GET("/get_by_board/:board_id", h.GetByBoardID)
	task.GET("/get_by_assignee/:assignee_id", h.GetByAssigneeID)
	task.PATCH("/update/:id", h.Update)
	task.DELETE("/delete/:id", h.Delete)
}
