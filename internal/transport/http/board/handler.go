package board

import (
	"net/http"

	boardapp "task_tracker/internal/application/board"
	"task_tracker/internal/domain/auth"
	"task_tracker/internal/perf"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BoardHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	ListByTeamID(c *gin.Context)
	Search(c *gin.Context)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Некорректный идентификатор доски"})
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
	defer perf.Track(c.Request.Context(), "handler_total")()
	perf.LogStep(c.Request.Context(), "handler_start")

	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	serviceDone := perf.Track(c.Request.Context(), "service.ListBoards")
	boards, err := h.service.List(c.Request.Context(), actor)
	serviceDone()
	if err != nil {
		status, msg := mapError(err)
		jsonDone := perf.Track(c.Request.Context(), "json_response")
		c.JSON(status, gin.H{"error": msg})
		jsonDone()
		return
	}

	jsonDone := perf.Track(c.Request.Context(), "json_response")
	c.JSON(http.StatusOK, gin.H{"boards": NewResponses(boards)})
	jsonDone()
}

func (h *handler) ListByTeamID(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Некорректный идентификатор команды"})
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

func (h *handler) Search(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	input, ok := newSearchBoardsInput(c)
	if !ok {
		return
	}

	boards, err := h.service.Search(c.Request.Context(), actor, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.HTML(http.StatusOK, "boards_select", gin.H{
		"Boards": NewResponses(boards),
	})
}

func newSearchBoardsInput(c *gin.Context) (boardapp.SearchBoardsInput, bool) {
	input := boardapp.SearchBoardsInput{}

	if query := c.Query("q"); query != "" {
		input.Query = &query
	}
	if teamID, ok := parseOptionalUUIDQuery(c, "team_id"); !ok {
		return boardapp.SearchBoardsInput{}, false
	} else {
		input.TeamID = teamID
	}
	if userID, ok := parseOptionalUUIDQuery(c, "user_id"); !ok {
		return boardapp.SearchBoardsInput{}, false
	} else {
		input.UserID = userID
	}

	return input, true
}

func (h *handler) Update(c *gin.Context) {
	actor, ok := actorFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Некорректный идентификатор доски"})
		return
	}

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Некорректный идентификатор доски"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return auth.Actor{}, false
	}
	return actor, true
}

func parseOptionalUUIDQuery(c *gin.Context, name string) (*uuid.UUID, bool) {
	raw := c.Query(name)
	if raw == "" || raw == "all" {
		return nil, true
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": invalidBoardFilterMessage(name)})
		return nil, false
	}
	if id == uuid.Nil {
		return nil, true
	}

	return &id, true
}

func invalidBoardFilterMessage(name string) string {
	switch name {
	case "team_id":
		return "Некорректный идентификатор команды"
	case "user_id":
		return "Некорректный идентификатор пользователя"
	default:
		return "Некорректный фильтр досок"
	}
}
