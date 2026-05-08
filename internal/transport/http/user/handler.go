package user

import (
	"net/http"
	user_application "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
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
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel := input.ToServiceInput()

	user, err := h.service.CreateRegister(ctx.Request.Context(), inputModel)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(*user)
	ctx.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
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
			"error": common_errors.ErrBadRequest.Error(),
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
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(*user)
	ctx.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
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
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel, err := input.ToServiceInput()
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	updatedUser, err := h.service.Update(ctx.Request.Context(), actor, inputModel)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(*updatedUser)
	ctx.JSON(http.StatusOK, gin.H{
		"user": response,
	})

}
