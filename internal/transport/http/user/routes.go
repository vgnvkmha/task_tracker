package user

import (
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h UserHandler) {
	user := r.Group("/user", middleware.ActorMiddleware())

	user.POST("/create_register", h.CreateRegister)
	user.POST("/create_by_actor", h.CreateByActor)
	user.PATCH("/update", h.Update)
}
