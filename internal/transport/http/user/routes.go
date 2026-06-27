package user

import (
	"fmt"
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h UserHandler, tokens *infraauth.JWTService, legacyHeadersEnabled bool) {
	ui := r.Group("/ui/users")
	ui.GET("/create", h.ShowCreateForm)
	ui.GET("/success", h.ShowAuthSuccess)
	ui.POST("", h.SubmitCreateForm)
	ui.GET("/login", h.SubmitLoginForm)
	ui.POST("/restore", h.SubmitRestoreForm)
	ui.POST("/logout", h.Logout)

	uiAuth := ui.Group("", middleware.UIActorMiddleware(tokens))
	uiAuth.GET("/cabinet", h.ShowCabinet)
	uiAuth.GET("/admin", h.ShowAdmin)
	uiAuth.POST("/cabinet/update", h.UpdateCabinet)
	uiAuth.POST("/cabinet/team-change-request", h.SubmitTeamChangeRequest)
	uiAuth.POST("/cabinet/delete", h.DeleteCabinet)

	user := r.Group("/user", middleware.ActorMiddleware(tokens, legacyHeadersEnabled))
	fmt.Println("REGISTER ROUTES")

	user.POST("/create_register", h.CreateRegister)
	user.POST("/create_by_actor", h.CreateByActor)
	user.PATCH("/update", h.Update)

	user.GET("/get_by_id/:user_id", h.GetByID)
	user.GET("/list_active", h.ListActive)
	user.GET("/list", h.List)

	user.DELETE("/delete/:id", h.DeleteByID)

}
