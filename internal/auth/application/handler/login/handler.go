package login

import (
	"errors"

	"github.com/Arheon/markus-backend/internal/auth/domain/repository"
	domainHelpers "github.com/Arheon/markus-backend/internal/auth/infrastructure/helpers"
	"github.com/Arheon/markus-backend/internal/shared/domain/entity"
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

func NewHandler(repo repository.UserRepository) *Handler {
	return &Handler{
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

	userID := loginVals.Login
	password := loginVals.Password

	h.logger.Debug("Try to find user by userID")
	user, err := h.userRepo.GetUserByUsername(userID)
	if err != nil {
		h.logger.Errorf("Could't find user by userName %s", userID)
		return "", ErrInvalidLoginValues
	}

	if domainHelpers.CheckPasswordHash(password, user.Password) == false {
		return "", ErrInvalidLoginValues
	}

	return user, nil
}

// TODO: Переделать на нормальную авторизацию
func (h *Handler) HandleAuthorize(c *gin.Context, data any) bool {
	if v, ok := data.(*entity.User); ok && v.Username == "admin" {
		return true
	}
	return false
}
