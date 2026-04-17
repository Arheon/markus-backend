package identity

import (
	"fmt"

	"github.com/Arheon/markus-backend/internal/auth/domain/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo   repository.UserRepository
	logger *logrus.Logger
}

func NewHandler(repo repository.UserRepository, logger *logrus.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *Handler) Handle(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)

	if identity, ok := claims["id"]; ok {
		user, err := h.repo.GetUserByUsername(c, identity.(string))
		if err != nil {
			h.logger.Error(fmt.Sprintf("Undefined user with id %s", identity.(string)))
			return nil
		}

		return user
	}

	return nil
}
