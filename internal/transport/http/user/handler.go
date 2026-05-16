package user

import (
	"net/http"
	"net/url"
	user_application "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	ShowCreateForm(c *gin.Context)
	ShowAuthSuccess(c *gin.Context)
	ShowCabinet(c *gin.Context)
	SubmitCreateForm(c *gin.Context)
	SubmitLoginForm(c *gin.Context)
	UpdateCabinet(c *gin.Context)
	DeleteCabinet(c *gin.Context)
	CreateRegister(c *gin.Context)
	CreateByActor(c *gin.Context)
	Update(c *gin.Context)
	GetByID(c *gin.Context)
	ListActive(c *gin.Context)
	List(c *gin.Context)
	DeleteByID(c *gin.Context)
}

type handler struct {
	service user_application.UserService
}

func New(service user_application.UserService) UserHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "user_create_page", gin.H{
		"title": "Вход или регистрация",
	})
}

func (h *handler) ShowAuthSuccess(c *gin.Context) {
	message := "Успешная регистрация"
	switch c.Query("type") {
	case "login":
		message = "Успешный логин"
	case "delete":
		message = "Пользователь удален"
	}

	c.HTML(http.StatusOK, "user_auth_success_page", gin.H{
		"title":   message,
		"message": message,
	})
}

func (h *handler) SubmitCreateForm(c *gin.Context) {
	input, err := NewCreateFormRequest(c.Request)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIFormError(err),
		})
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.service.CreateRegister(ctx, input.ToServiceInput())
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	redirectUI(c, cabinetURL(createdUser.ID.String(), "registered", ""), http.StatusCreated)
}

func (h *handler) SubmitLoginForm(c *gin.Context) {
	input, err := NewLoginFormRequest(c.Request)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIFormError(err),
		})
		return
	}

	ctx := c.Request.Context()
	loggedUser, err := h.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	redirectUI(c, cabinetURL(loggedUser.ID.String(), "login", ""), http.StatusOK)
}

func redirectUI(c *gin.Context, location string, status int) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", location)
		c.Status(status)
		return
	}

	c.Redirect(http.StatusSeeOther, location)
}

func (h *handler) ShowCabinet(c *gin.Context) {
	id, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "user_auth_success_page", gin.H{
			"title": "Ошибка",
			"error": mapUIError(user_application.ErrInvalidUserID),
		})
		return
	}

	profile, err := h.service.GetProfileByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	renderCabinet(c, profile, cabinetMessage(c.Query("message")), c.Query("error"))
}

func (h *handler) UpdateCabinet(c *gin.Context) {
	input, err := NewUpdateProfileFormRequest(c.Request)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIFormError(err),
		})
		return
	}

	profile, err := h.service.GetProfileByID(c.Request.Context(), input.UserID)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	actor := auth.Actor{
		ID:   profile.User.ID,
		Role: profile.User.Role,
	}
	if _, err := h.service.Update(c.Request.Context(), actor, input.ToServiceInput()); err != nil {
		redirectUI(c, cabinetURL(input.UserID.String(), "", mapUIError(err)), http.StatusSeeOther)
		return
	}

	redirectUI(c, cabinetURL(input.UserID.String(), "updated", ""), http.StatusSeeOther)
}

func (h *handler) DeleteCabinet(c *gin.Context) {
	userID, err := uuid.Parse(c.PostForm("user_id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "user_create_result", gin.H{
			"error": mapUIError(user_application.ErrInvalidUserID),
		})
		return
	}

	profile, err := h.service.GetProfileByID(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	actor := auth.Actor{
		ID:   profile.User.ID,
		Role: profile.User.Role,
	}
	if err := h.service.DeleteByID(c.Request.Context(), actor, userID); err != nil {
		redirectUI(c, cabinetURL(userID.String(), "", mapUIError(err)), http.StatusSeeOther)
		return
	}

	redirectUI(c, "/ui/users/success?type=delete", http.StatusSeeOther)
}

func renderCabinet(c *gin.Context, profile *user_application.Profile, message string, errorMessage string) {
	c.HTML(http.StatusOK, "user_cabinet_page", gin.H{
		"title":   "Личный кабинет",
		"profile": NewCabinetView(profile),
		"message": message,
		"error":   errorMessage,
	})
}

func cabinetMessage(value string) string {
	switch value {
	case "registered":
		return "Успешная регистрация"
	case "login":
		return "Успешный логин"
	case "updated":
		return "Данные обновлены"
	default:
		return ""
	}
}

func cabinetURL(userID string, message string, errorMessage string) string {
	values := url.Values{}
	values.Set("user_id", userID)
	if message != "" {
		values.Set("message", message)
	}
	if errorMessage != "" {
		values.Set("error", errorMessage)
	}
	return "/ui/users/cabinet?" + values.Encode()
}

func (h *handler) CreateRegister(c *gin.Context) {
	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}
	inputModel := input.ToServiceInput()

	ctx := c.Request.Context()
	user, err := h.service.CreateRegister(ctx, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(user)
	c.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
}

func (h *handler) CreateByActor(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": common_errors.ErrUnauthorized.Error(),
		})
		return
	}

	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}

	inputModel := input.ToServiceInput()
	ctx := c.Request.Context()
	user, err := h.service.CreateByActor(ctx, actor, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	response := FromDomain(user)
	c.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
}

func (h *handler) Update(c *gin.Context) {
	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}

	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": common_errors.ErrPermissionDenied.Error(),
		})
		return
	}

	inputModel, err := input.ToServiceInput()
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	isManager := actor.Role.IsManagerRole()
	isSelf := actor.ID == inputModel.UserID
	if !isManager && !isSelf {
		c.JSON(http.StatusForbidden, gin.H{
			"error": common_errors.ErrPermissionDenied.Error(),
		})
		return
	}
	if !isManager && (input.Role != nil || input.TeamID != nil) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": common_errors.ErrPermissionDenied.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	updatedUser, err := h.service.Update(ctx, actor, inputModel)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	response := FromDomain(updatedUser)
	c.JSON(http.StatusOK, gin.H{
		"user": response,
	})

}

func (h *handler) GetByID(c *gin.Context) {
	id := c.Param("user_id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": common_errors.ErrInvalidID.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	user, err := h.service.GetByID(ctx, uuid)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromDomain(user)

	c.JSON(http.StatusOK, gin.H{
		"user": response,
	})
}

func (h *handler) ListActive(c *gin.Context) {
	var response []Response
	ctx := c.Request.Context()
	activeUsers, err := h.service.ListActive(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = FromDomainReponses(activeUsers)

	c.JSON(http.StatusOK, gin.H{
		"active users": response,
	})
}

func (h *handler) List(c *gin.Context) {
	var response []Response
	ctx := c.Request.Context()
	users, err := h.service.List(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response = FromDomainReponses(users)

	c.JSON(http.StatusOK, gin.H{
		"users": response,
	})
}

func (h *handler) DeleteByID(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": common_errors.ErrUnauthorized.Error(),
		})
		return
	}
	user_id := c.Param("id")
	uuid, err := uuid.Parse(user_id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": common_errors.ErrInvalidID.Error(),
		})
		return

	}
	ctx := c.Request.Context()
	err = h.service.DeleteByID(ctx, actor, uuid)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user deleted": user_id,
	})
}
