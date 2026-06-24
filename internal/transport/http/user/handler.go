package user

import (
	"net/http"
	"net/url"
	user_application "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"
	"task_tracker/internal/perf"
	"task_tracker/internal/transport/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	ShowCreateForm(c *gin.Context)
	ShowAuthSuccess(c *gin.Context)
	ShowCabinet(c *gin.Context)
	SubmitCreateForm(c *gin.Context)
	SubmitLoginForm(c *gin.Context)
	Logout(c *gin.Context)
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
	tokens  TokenGenerator
}

type TokenGenerator interface {
	GenerateAccessToken(userID uuid.UUID, role valueobjects.Role, teamID *uuid.UUID) (string, auth.TokenClaims, error)
}

func New(service user_application.UserService, tokens ...TokenGenerator) UserHandler {
	var tokenGenerator TokenGenerator
	if len(tokens) > 0 {
		tokenGenerator = tokens[0]
	}

	return &handler{
		service: service,
		tokens:  tokenGenerator,
	}
}

func (h *handler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "user_create_page", gin.H{
		"title":        "Вход или регистрация",
		"auth_message": authMessage(c.Query("auth")),
	})
}

func authMessage(value string) string {
	switch value {
	case "expired":
		return "Сессия истекла. Войдите снова."
	case "required":
		return "Для доступа к кабинету пользователя нужно войти."
	case "logout":
		return "Вы вышли из аккаунта."
	default:
		return ""
	}
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

	if h.tokens != nil {
		token, claims, err := h.tokens.GenerateAccessToken(createdUser.ID, createdUser.Role, createdUser.TeamID)
		if err != nil {
			c.HTML(http.StatusOK, "user_create_result", gin.H{
				"error": "Не удалось создать сессию",
			})
			return
		}
		setAccessTokenCookie(c, token, claims.ExpiresAtTime())
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

	if h.tokens != nil {
		token, claims, err := h.tokens.GenerateAccessToken(loggedUser.ID, loggedUser.Role, loggedUser.TeamID)
		if err != nil {
			c.HTML(http.StatusOK, "user_create_result", gin.H{
				"error": "Не удалось создать сессию",
			})
			return
		}
		setAccessTokenCookie(c, token, claims.ExpiresAtTime())
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

func (h *handler) Logout(c *gin.Context) {
	clearAccessTokenCookie(c)
	noStore(c)
	replaceLocation(c, "/ui/users/create?auth=logout")
}

func setAccessTokenCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAccessTokenCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func noStore(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func replaceLocation(c *gin.Context, location string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta http-equiv="Cache-Control" content="no-store">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Logout</title>
</head>
<body>
  <script>window.location.replace("`+location+`");</script>
  <noscript><a href="`+location+`">Перейти на страницу входа</a></noscript>
</body>
</html>`))
}

func (h *handler) ShowCabinet(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.HTML(http.StatusUnauthorized, "user_auth_success_page", gin.H{
			"title": "Ошибка",
			"error": common_errors.ErrUnauthorized.Error(),
		})
		return
	}

	userID := actor.ID
	if rawUserID := c.Query("user_id"); rawUserID != "" {
		id, err := uuid.Parse(rawUserID)
		if err != nil {
			c.HTML(http.StatusUnprocessableEntity, "user_auth_success_page", gin.H{
				"title": "Ошибка",
				"error": mapUIError(user_application.ErrInvalidUserID),
			})
			return
		}
		if id != actor.ID {
			c.HTML(http.StatusForbidden, "user_auth_success_page", gin.H{
				"title": "Ошибка",
				"error": common_errors.ErrPermissionDenied.Error(),
			})
			return
		}
		userID = id
	}

	profile, err := h.service.GetProfileByID(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIError(err),
		})
		return
	}

	renderCabinet(c, profile, cabinetMessage(c.Query("message")), c.Query("error"))
}

func (h *handler) UpdateCabinet(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.HTML(http.StatusUnauthorized, "user_create_result", gin.H{
			"error": common_errors.ErrUnauthorized.Error(),
		})
		return
	}

	input, err := NewUpdateProfileFormRequest(c.Request)
	if err != nil {
		c.HTML(http.StatusOK, "user_create_result", gin.H{
			"error": mapUIFormError(err),
		})
		return
	}

	if actor.ID != input.UserID {
		c.HTML(http.StatusForbidden, "user_create_result", gin.H{
			"error": common_errors.ErrPermissionDenied.Error(),
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

	updateActor := auth.Actor{
		ID:     profile.User.ID,
		Role:   profile.User.Role,
		TeamID: profile.User.TeamID,
	}
	updatedUser, err := h.service.Update(c.Request.Context(), updateActor, input.ToServiceInput())
	if err != nil {
		redirectUI(c, cabinetURL(input.UserID.String(), "", mapUIError(err)), http.StatusSeeOther)
		return
	}
	if h.tokens != nil {
		token, claims, err := h.tokens.GenerateAccessToken(updatedUser.ID, updatedUser.Role, updatedUser.TeamID)
		if err == nil {
			setAccessTokenCookie(c, token, claims.ExpiresAtTime())
		}
	}

	redirectUI(c, cabinetURL(input.UserID.String(), "updated", ""), http.StatusSeeOther)
}

func (h *handler) DeleteCabinet(c *gin.Context) {
	actor, ok := middleware.GetActor(c)
	if !ok {
		c.HTML(http.StatusUnauthorized, "user_create_result", gin.H{
			"error": common_errors.ErrUnauthorized.Error(),
		})
		return
	}

	userID, err := uuid.Parse(c.PostForm("user_id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "user_create_result", gin.H{
			"error": mapUIError(user_application.ErrInvalidUserID),
		})
		return
	}

	if actor.ID != userID {
		c.HTML(http.StatusForbidden, "user_create_result", gin.H{
			"error": common_errors.ErrPermissionDenied.Error(),
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

	deleteActor := auth.Actor{
		ID:   profile.User.ID,
		Role: profile.User.Role,
	}
	if err := h.service.DeleteByID(c.Request.Context(), deleteActor, userID); err != nil {
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
	defer perf.Track(c.Request.Context(), "handler_total")()
	perf.LogStep(c.Request.Context(), "handler_start")

	ctx := c.Request.Context()
	serviceDone := perf.Track(ctx, "service.ListActiveProfiles")
	activeUsers, err := h.service.ListActiveProfiles(ctx)
	serviceDone()
	if err != nil {
		status, msg := mapError(err)
		jsonDone := perf.Track(ctx, "json_response")
		c.JSON(status, gin.H{
			"error": msg,
		})
		jsonDone()
		return
	}
	response := FromProfiles(activeUsers)

	jsonDone := perf.Track(ctx, "json_response")
	c.JSON(http.StatusOK, gin.H{
		"users":        response,
		"active users": response,
	})
	jsonDone()
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.service.ListActiveProfiles(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{
			"error": msg,
		})
		return
	}
	response := FromProfiles(users)

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
