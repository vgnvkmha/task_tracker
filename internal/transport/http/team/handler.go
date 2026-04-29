package team

import (
	"net/http"
	"task_tracker/internal/application/team"

	"github.com/gin-gonic/gin"
)

type TeamHandler interface {
	Create(ctx *gin.Context)

	GetByID(ctx *gin.Context)
	GetByName(ctx *gin.Context)

	ListActive(ctx *gin.Context)
	List(ctx *gin.Context)

	Update(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

type handler struct {
	service team.TeamService
}

func New(service team.TeamService) TeamHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(ctx *gin.Context) {
	var input CreateTeamRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	appInput, err := ToServiceInput(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	team, err := h.service.Create(ctx, *appInput)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	response := FromDomain(team)

	ctx.JSON(http.StatusCreated, response)
}

func (h *handler) GetByID(ctx *gin.Context) {}

func (h *handler) GetByName(ctx *gin.Context) {}

func (h *handler) ListActive(ctx *gin.Context) {}

func (h *handler) List(ctx *gin.Context) {}

func (h *handler) Update(ctx *gin.Context) {}

func (h *handler) DeleteByID(ctx *gin.Context) {}
