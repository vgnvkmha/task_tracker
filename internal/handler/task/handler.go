package task_handler

import (
	"strconv"
	"task_tracker/internal/domain/service"
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
	service service.TaskService
}

func New(service service.TaskService) handler { //TODO: must return interface
	return handler{
		service: service,
	}
}

func (h *handler) Create(ctx *gin.Context) {
	var params dto.TaskRequest

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), params)
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
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	status, err := h.service.ChangeStatus(ctx.Request.Context(), uint32(id), input)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(201, gin.H{
		"task_id":    id,
		"new_status": status,
	})
}
