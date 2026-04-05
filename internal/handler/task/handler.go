package task_handler

import (
	"strconv"
	task_service "task_tracker/internal/application/task"
	dto "task_tracker/internal/handler/task/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler interface {
	Create(ctx *gin.Context)

	ListActiveByTeam(ctx *gin.Context)

	ChangeStatus(ctx *gin.Context)
	ChangeBoard(ctx *gin.Context)
	ChangeAssign(ctx *gin.Context)
	ChangeReporter(ctx *gin.Context)
	ChangeSprint(ctx *gin.Context)
}

type handler struct {
	service task_service.TaskService
}

func New(service task_service.TaskService) handler { //TODO: must return interface
	return handler{
		service: service,
	}
}

func (h *handler) Create(ctx *gin.Context) {
	var params dto.TaskRequest
	strId := ctx.Param("id")
	userId, err := uuid.Parse(strId)
	if err != nil {
		ctx.JSON(400, gin.H{
			"invalid_user_id": strId,
			"error":           err,
		})
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), uuid.UUID(userId), params)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToTaskResponse(task)
	ctx.JSON(201, gin.H{
		"created_task": response,
	})
}

func (h *handler) ListActiveByTeam(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	tasks, err := h.service.GetActiveTasksByTeam(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	results := dto.ToTaskResponses(tasks)
	ctx.JSON(201, gin.H{
		"active_tasks": results,
	})
}

func (h *handler) ChangeStatus(ctx *gin.Context) {
	input := ctx.Param("status")
	// roleStr := ctx.Param("role")
	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"invalid_user_id": userId, "error": err})
		return
	}
	taskId, err := strconv.ParseUint(ctx.Param("task_id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"invalid__task_id": taskId, "error": err})
		return
	}

	status, err := h.service.ChangeStatus(ctx.Request.Context(), userId, taskId, input)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(201, gin.H{
		"task_id":    taskId,
		"user_id":    userId,
		"new_status": status,
	})
}
