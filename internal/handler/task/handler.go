package task

import (
	"task_tracker/internal/domain/service"

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
