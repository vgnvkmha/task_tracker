package app

import (
	"context"
	// task_service "task_tracker/internal/application/task"
	"task_tracker/internal/application/user"
	"task_tracker/internal/configs"
	"task_tracker/internal/domain/logger"
	"task_tracker/internal/infrastracture/db"

	// task_handler "task_tracker/internal/handler/task"
	"task_tracker/internal/repo/team"
	user_repo "task_tracker/internal/repo/user"
	"task_tracker/internal/transport/http/middleware"
	handler_user "task_tracker/internal/transport/http/user"

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
	// repo := repo.New(pDb)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, loggerErr := logger.New()
	if loggerErr != nil {
		return loggerErr
	}

	userRepo := user_repo.NewUserRepo(pDb)
	personalDataRepo := user_repo.NewPersonalDataRepo(pDb)
	teamRepo := team.NewTeamRepo(pDb)

	txManager := db.NewTxManager(pDb)
	userService := user.New(userRepo, personalDataRepo, teamRepo, logger, txManager)
	userHandler := handler_user.New(userService)

	// taskService := task_service.New(repo, logger)
	// taskHandler := task_handler.New(service)

	router := gin.Default()
	router.Use(middleware.MockActorMiddleware())
	handler_user.RegisterRoutes(router, userHandler)
	// task_handler.RegisterRoutes(router, handler)
	return router.RunTLS(":8080", "cert.pem", "key.pem")
}
