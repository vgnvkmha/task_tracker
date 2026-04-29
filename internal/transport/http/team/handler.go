package team

import (
	"net/http"
	"task_tracker/internal/application/team"
	"task_tracker/internal/common_errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler interface {
	Create(c *gin.Context)

	GetByID(c *gin.Context)
	GetByName(c *gin.Context)

	ListActive(c *gin.Context)
	List(c *gin.Context)

	Update(c *gin.Context)
	DeleteByID(c *gin.Context)
}

type handler struct {
	service team.TeamService
}

func New(service team.TeamService) TeamHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(c *gin.Context) {
	var input CreateTeamRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest,
		})
		return
	}

	appInput, err := NewApplicationTeam(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest,
		})
		return
	}

	ctx := c.Request.Context()
	team, err := h.service.Create(ctx, *appInput)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	response := NewResponse(team)

	c.JSON(http.StatusCreated, response)
}

func (h *handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrInvalidID,
		})
		return
	}
	ctx := c.Request.Context()

	team, err := h.service.GetByID(ctx, uuid)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := NewResponse(team)

	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *handler) GetByName(c *gin.Context) {
	name := c.Param("team_name")
	ctx := c.Request.Context()
	team, err := h.service.GetByName(ctx, name)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := NewResponse(team)

	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *handler) ListActive(c *gin.Context) {
	var response []*Response
	ctx := c.Request.Context()
	activeTeams, err := h.service.ListActive(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = NewResponses(activeTeams)

	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *handler) List(c *gin.Context) {
	var response []*Response
	ctx := c.Request.Context()
	teams, err := h.service.List(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = NewResponses(teams)

	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})
}

func (h *handler) Update(c *gin.Context) {
	team_id := c.Param("id")
	id, err := uuid.Parse(team_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrInvalidID,
		})
		return
	}
	var req UpdateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest,
		})
		return
	}

	input, err := ApplyUpdateTeam(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx := c.Request.Context()
	team, err := h.service.Update(ctx, id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := NewResponse(team)
	c.JSON(http.StatusOK, gin.H{
		"team": response,
	})

}

func (h *handler) DeleteByID(c *gin.Context) {
	team_id := c.Param("id")
	uuid, err := uuid.Parse(team_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrInvalidID,
		})
		return

	}
	ctx := c.Request.Context()
	err = h.service.DeleteByID(ctx, uuid)
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
