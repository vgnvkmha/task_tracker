package user

import (
	"net/http"
	user_application "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	CreateRegister(c *gin.Context)
	CreateByActor(c *gin.Context)
	Update(c *gin.Context)
	GetByID(c *gin.Context)
	ListActive(c *gin.Context)
	List(c *gin.Context)
	DeleteByID(c *gin.Context)
}

type handler struct {
	service user_application.UserService
}

func New(service user_application.UserService) UserHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateRegister(c *gin.Context) {
	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel := input.ToServiceInput()

	ctx := c.Request.Context()
	user, err := h.service.CreateRegister(ctx, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(user)
	c.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
}

func (h *handler) CreateByActor(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel := input.ToServiceInput()

	ctx := c.Request.Context()
	user, err := h.service.CreateByActor(ctx, actor, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(user)
	c.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
}

func (h *handler) Update(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel, err := input.ToServiceInput()
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	ctx := c.Request.Context()
	updatedUser, err := h.service.Update(ctx, actor, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(updatedUser)
	c.JSON(http.StatusOK, gin.H{
		"user": response,
	})

}

func (h *handler) GetByID(c *gin.Context) {
	id := c.Param("user_id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid id",
		})
		return
	}
	ctx := c.Request.Context()
	user, err := h.service.GetByID(ctx, uuid)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(user)

	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *handler) ListActive(c *gin.Context) {
	var response []Response
	ctx := c.Request.Context()
	activeUsers, err := h.service.ListActive(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = FromDomainReponses(activeUsers)

	c.JSON(http.StatusOK, gin.H{
		"active_users": response,
	})
}

func (h *handler) List(c *gin.Context) {
	var response []Response
	ctx := c.Request.Context()
	activeUsers, err := h.service.List(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = FromDomainReponses(activeUsers)

	c.JSON(http.StatusOK, gin.H{
		"users": response,
	})
}

func (h *handler) DeleteByID(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	team_id := c.Param("id")
	uuid, err := uuid.Parse(team_id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": common_errors.ErrInvalidID,
		})
		return

	}
	ctx := c.Request.Context()
	err = h.service.DeleteByID(ctx, actor, uuid)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"team_deleted": team_id,
	})
}
