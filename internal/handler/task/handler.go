package task_handler

import (
	"strconv"
	task_service "task_tracker/internal/application/task"
	dto "task_tracker/internal/transport/task"

	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	Create(ctx *gin.Context)

	ListActive(ctx *gin.Context)

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
	userId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"invalid_id": userId, "error": err})
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), uint32(userId), params)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"created_task": task,
	})
}

func (h *handler) ListActive(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	tasks, err := h.service.GetActiveTasksByTeam(ctx.Request.Context(), uint32(id))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{
		"active_tasks": tasks,
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

	status, err := h.service.ChangeStatus(ctx.Request.Context(), uint32(userId), uint32(taskId), input)
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
