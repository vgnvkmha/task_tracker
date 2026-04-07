package task_handler

import (
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	h TaskHandler,
) {
	r.Use(middleware.ActorMiddleware())

}
