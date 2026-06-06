package board

import (
	"net/http"

	boardapp "task_tracker/internal/application/board"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BoardHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	ListByTeamID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type handler struct {
	service boardapp.BoardService
}

func New(service boardapp.BoardService) BoardHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	var req CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": common_errors.ErrBadRequest.Error()})
		return
	}

	input, err := req.ToApplicationInput()
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	board, err := h.service.Create(c.Request.Context(), actor, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": NewResponse(board)})
}

func (h *handler) GetByID(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": common_errors.ErrInvalidID.Error()})
		return
	}

	board, err := h.service.GetByID(c.Request.Context(), actor, id)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": NewResponse(board)})
}

func (h *handler) List(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	boards, err := h.service.List(c.Request.Context(), actor)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": NewResponses(boards)})
}

func (h *handler) ListByTeamID(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": common_errors.ErrInvalidID.Error()})
		return
	}

	boards, err := h.service.ListByTeamID(c.Request.Context(), actor, teamID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": NewResponses(boards)})
}

func (h *handler) Update(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": common_errors.ErrInvalidID.Error()})
		return
	}

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": common_errors.ErrBadRequest.Error()})
		return
	}

	input, err := req.ToApplicationInput()
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	board, err := h.service.Update(c.Request.Context(), actor, id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": NewResponse(board)})
}

func (h *handler) Delete(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": common_errors.ErrInvalidID.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), actor, id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board_deleted": id.String()})
}

func actorFromContext(c *gin.Context) (auth.Actor, bool) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor"})
		return auth.Actor{}, false
	}
	return actor, true
}
