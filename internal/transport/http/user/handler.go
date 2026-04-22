package user

import (
	"errors"
	"fmt"
	"net/http"
	user_application "task_tracker/internal/application/user"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateRegister(ctx *gin.Context)
	CreateByActor(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type handler struct {
	service user_application.UserService
}

func New(service user_application.UserService) UserHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateRegister(ctx *gin.Context) {
	var input CreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("TeamName: %#v\n", input.TeamName)
	inputModel := input.ToServiceInput()

	user, err := h.service.CreateRegister(ctx.Request.Context(), inputModel)
	if err != nil {
		switch {
		case errors.Is(err, user_application.ErrTeamNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "team not found",
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})
			return
		}
	}
	response := FromService(user)
	ctx.JSON(http.StatusOK, response)
}

func (h *handler) CreateByActor(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input CreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	inputModel := input.ToServiceInput()

	user, err := h.service.CreateByActor(
		ctx.Request.Context(),
		actor,
		inputModel,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := FromService(user)
	ctx.JSON(http.StatusOK, response)
}

func (h *handler) Update(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input UpdateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	inputModel, err := input.ToServiceInput()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	updatedUser, err := h.service.Update(ctx, actor, inputModel)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	response := FromService(updatedUser)
	ctx.JSON(http.StatusOK, response)

}
