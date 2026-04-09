package task_handler

import (
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	h TaskHandler,
) {
	r.POST("/create_task", h.Create)
	r.GET("/get_tasks_by_team/:id", h.ListActiveByTeam)

	r.POST("change/status", h.ChangeStatus)
	r.POST("change/board", h.ChangeBoard)
	r.POST("change/assign", h.ChangeAssign)
	r.POST("change/reporter", h.ChangeReporter)
	r.POST("change/sprint", h.ChangeSprint)

	r.Use(middleware.ActorMiddleware())

}
