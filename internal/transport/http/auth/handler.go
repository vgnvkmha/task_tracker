package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	userapplication "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	domainauth "task_tracker/internal/domain/auth"
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserLoginService interface {
	Login(ctx context.Context, email string, password string) (*userapplication.User, error)
}

type TokenGenerator interface {
	GenerateAccessToken(userID uuid.UUID, role valueobjects.Role, teamID *uuid.UUID) (string, domainauth.TokenClaims, error)
}

type Handler struct {
	users  UserLoginService
	tokens TokenGenerator
}

func New(users UserLoginService, tokens TokenGenerator) *Handler {
	return &Handler{
		users:  users,
		tokens: tokens,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": common_errors.ErrBadRequest.Error(),
		})
		return
	}

	user, err := h.users.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, userapplication.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": common_errors.ErrUnauthorized.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	token, claims, err := h.tokens.GenerateAccessToken(user.ID, user.Role, user.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	expiresIn := claims.ExpiresAtTime().Sub(time.Unix(claims.IssuedAt, 0)).Seconds()
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiresIn),
	})
}
