package user_service

import (
	"context"
	"task_tracker/internal/domain/models"
	"task_tracker/internal/handler/task/dto"
)

type UserService interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, update dto.UpdateUser) error
}

type service struct {
}
