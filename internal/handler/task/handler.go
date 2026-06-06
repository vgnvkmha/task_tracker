package task_handler

import (
	"fmt"
	"net/http"

	taskApplication "task_tracker/internal/application/task"
	domaintask "task_tracker/internal/domain/task"
	dto "task_tracker/internal/handler/task/dto"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	List(ctx *gin.Context)
	GetByBoardID(ctx *gin.Context)
	GetByAssigneeID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type handler struct {
	service taskApplication.TaskService
}

func New(service taskApplication.TaskService) TaskHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.TaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		fmt.Printf("bind error: %v\n", err)
		fmt.Printf("bind error type: %T\n", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), actor, request.ToServiceInput())
	if err != nil {
		status, msg := mapError(err)
		fmt.Print("-----------")
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) GetByID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	task, err := h.service.GetByID(ctx.Request.Context(), actor, id)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) List(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filters, ok := parseTaskFilters(ctx)
	if !ok {
		return
	}

	tasks, err := h.service.FindMany(ctx.Request.Context(), actor, filters)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) GetByBoardID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	boardID, ok := parseIDParam(ctx, "board_id")
	if !ok {
		return
	}

	tasks, err := h.service.FindByBoardID(ctx.Request.Context(), actor, boardID)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) GetByAssigneeID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	assigneeID, ok := parseIDParam(ctx, "assignee_id")
	if !ok {
		return
	}

	tasks, err := h.service.FindByAssigneeID(ctx.Request.Context(), actor, assigneeID)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) Update(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	var request dto.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := h.service.Update(ctx.Request.Context(), actor, id, request.ToServiceInput())
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) Delete(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), actor, id); err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task_deleted": id.String(),
	})
}

func parseTaskFilters(ctx *gin.Context) (taskApplication.TaskFilters, bool) {
	var filters taskApplication.TaskFilters

	if !parseOptionalUUIDQuery(ctx, "board_id", &filters.BoardID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "assignee_id", &filters.AssigneeID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "reporter_id", &filters.ReporterID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "sprint_id", &filters.SprintID) {
		return taskApplication.TaskFilters{}, false
	}
	if rawStatus := ctx.Query("status"); rawStatus != "" {
		status := domaintask.TaskStatus(rawStatus)
		filters.Status = &status
	}

	return filters, true
}

func parseOptionalUUIDQuery(ctx *gin.Context, name string, dest **uuid.UUID) bool {
	raw := ctx.Query(name)
	if raw == "" {
		return true
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid " + name})
		return false
	}

	*dest = &id
	return true
}

func parseIDParam(ctx *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return id, true
}
