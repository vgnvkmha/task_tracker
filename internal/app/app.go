package app

import (
	"context"
	"net/http"
	"strings"
	board_application "task_tracker/internal/application/board"
	task_service "task_tracker/internal/application/task"
	team_application "task_tracker/internal/application/team"
	"task_tracker/internal/application/user"
	"task_tracker/internal/configs"
	"task_tracker/internal/domain/logger"
	task_handler "task_tracker/internal/handler/task"
	infraauth "task_tracker/internal/infrastracture/auth"
	"task_tracker/internal/infrastracture/db"
	taskRepo "task_tracker/internal/repo"
	"task_tracker/internal/repo/team"
	userRepo "task_tracker/internal/repo/user"
	auth_handler "task_tracker/internal/transport/http/auth"
	board_handler "task_tracker/internal/transport/http/board"
	"task_tracker/internal/transport/http/middleware"
	team_handler "task_tracker/internal/transport/http/team"
	handler_user "task_tracker/internal/transport/http/user"

	"github.com/gin-gonic/gin"
)

func Run() error {
	postgresCfg, err := configs.LoadPostgres()
	if err != nil {
		return err
	}
	authCfg, err := configs.LoadAuth()
	if err != nil {
		return err
	}
	jwtService, err := infraauth.NewJWTService(authCfg.JWTSecret, authCfg.AccessTokenTTL)
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

	usersRepo := userRepo.NewUserRepo(pDb)
	personalDataRepo := userRepo.NewPersonalDataRepo(pDb)
	teamRepo := team.NewTeamRepo(pDb)
	tasksRepo := taskRepo.NewTaskRepo(pDb)
	boardsRepo := taskRepo.NewBoardRepo(pDb)

	txManager := db.NewTxManager(pDb)
	userService := user.New(usersRepo, personalDataRepo, teamRepo, logger, txManager)
	userHandler := handler_user.New(userService, jwtService)
	authHandler := auth_handler.New(userService, jwtService)

	teamService := team_application.New(teamRepo, usersRepo, logger, txManager)
	teamHandler := team_handler.New(teamService)
	taskService := task_service.New(tasksRepo, boardsRepo, logger, txManager)
	taskHandler := task_handler.New(taskService)
	boardService := board_application.New(boardsRepo, logger, txManager)
	boardHandler := board_handler.New(boardService)

	router := gin.Default()
	router.SetTrustedProxies(nil) //TODO: change later
	router.SetHTMLTemplate(handler_user.Templates())
	router.Use(staticCacheMiddleware())
	router.Static("/static", "./static")
	router.Use(middleware.MockActorMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/ui/login")
	})
	auth_handler.RegisterRoutes(router, authHandler)
	handler_user.RegisterRoutes(router, userHandler, jwtService, authCfg.LegacyHeadersEnabled)
	team_handler.RegisterRoutes(router, teamHandler, jwtService, authCfg.LegacyHeadersEnabled)
	task_handler.RegisterRoutes(router, taskHandler, jwtService, authCfg.LegacyHeadersEnabled)
	board_handler.RegisterRoutes(router, boardHandler, jwtService, authCfg.LegacyHeadersEnabled)
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "not_found_page", gin.H{
			"title": "Страница не найдена",
		})
	})
	return router.Run(":8080")
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	}
}
