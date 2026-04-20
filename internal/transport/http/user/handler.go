package user

import (
	"task_tracker/internal/application/user"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)

}

type handler struct {
	service user.UserService
}