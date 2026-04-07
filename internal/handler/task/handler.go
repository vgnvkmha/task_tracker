package task_handler

import (
	task_service "task_tracker/internal/application/task"
	dto "task_tracker/internal/handler/task/dto"
	"task_tracker/internal/transport/http/middleware"

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
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), actor, params)
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
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	statusStr := ctx.Param("status")
	taskIdStr := ctx.Param("id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		ctx.JSON(400, gin.H{"invalid__task_id": taskId, "error": err})
		return
	}

	status, err := h.service.ChangeStatus(ctx.Request.Context(), actor, taskId, statusStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(201, gin.H{
		"task_id":    taskId,
		"user_id":    actor.Id,
		"user_role":  actor.Role,
		"new_status": status,
	})
}
