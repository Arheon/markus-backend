package login

import (
	"errors"

	"github.com/Arheon/markus-backend/internal/auth/domain/repository"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidLoginValues = errors.New("invalid username or password")
)

type Handler struct {
	logger   *logrus.Logger
	userRepo repository.UserRepository
}

func NewHandler(logger *logrus.Logger, repo repository.UserRepository) *Handler {
	return &Handler{
		logger:   logger,
		userRepo: repo,
	}
}

func (h *Handler) HandleAuthenticate(ctx *gin.Context) (any, error) {
	var loginVals LoginForm
	if err := ctx.ShouldBind(&loginVals); err != nil {
		h.logger.Error("Bad authenticate request", map[string]any{
			"Error":   err,
			"Context": ctx,
		})

		return "", jwt.ErrMissingLoginValues
	}

	userID := loginVals.UserName
	password := loginVals.Password

	h.logger.Debug("Try to find user by userID")
	user, err := h.userRepo.GetUserByUsername(ctx, userID)
	if err != nil {
		h.logger.Errorf("Could't find user by userName %s", userID)
		return "", ErrInvalidLoginValues
	}

	if helpers.CheckPasswordHash(password, user.Password) == false {
		return "", ErrInvalidLoginValues
	}

	h.logger.Debug("ABOBA2", user.ID, user.Password)

	return user, nil
}

// TODO: Переделать на нормальную авторизацию
func (h *Handler) HandleAuthorize(c *gin.Context, data any) bool {
	return true
}
