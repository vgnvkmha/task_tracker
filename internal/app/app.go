package app

import (
	"context"
	task_service "task_tracker/internal/application/task"
	"task_tracker/internal/configs"
	"task_tracker/internal/domain/logger"
	task_handler "task_tracker/internal/handler/task"
	"task_tracker/internal/repo"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func Run() error {
	postgresCfg, err := configs.LoadPostgres()
	if err != nil {
		return err
	}
	pDb, err := configs.New(*postgresCfg)
	if err != nil {
		return err
	}
	repo := repo.New(pDb)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, loggerErr := logger.New()
	if loggerErr != nil {
		return loggerErr
	}

	service := task_service.New(repo, logger)
	handler := task_handler.New(service)

	router := gin.Default()
	router.Use(middleware.MockActorMiddleware())
	task_handler.RegisterRoutes(router, handler)
	return router.RunTLS(":8080", "cert.pem", "key.pem")
}
