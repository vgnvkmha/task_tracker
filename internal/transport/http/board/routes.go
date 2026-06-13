package board

import (
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h BoardHandler, tokens *infraauth.JWTService, legacyHeadersEnabled bool) {
	boardsUI := r.Group("/boards", middleware.UIActorMiddleware(tokens))
	boardsUI.GET("/search", h.Search)

	board := r.Group("/board", middleware.ActorMiddleware(tokens, legacyHeadersEnabled))

	board.POST("/create", h.Create)
	board.GET("/get_by_id/:id", h.GetByID)
	board.GET("/list", h.List)
	board.GET("/list_by_team_id/:team_id", h.ListByTeamID)
	board.PATCH("/update/:id", h.Update)
	board.DELETE("/delete/:id", h.Delete)
}
