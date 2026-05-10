package user

import (
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h UserHandler) {
	user := r.Group("/user/", middleware.ActorMiddleware())

	user.POST("create_register", h.CreateRegister)
	user.POST("create_by_actor", h.CreateByActor)
	user.PATCH("update", h.Update)

	user.GET("get_by_id/:user_id", h.GetByID)
	user.GET("list_active", h.ListActive)
	user.GET("list", h.List)

	user.POST("delete", h.DeleteByID)

}
