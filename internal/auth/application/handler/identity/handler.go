package identity

import (
	"github.com/Arheon/markus-backend/internal/shared/domain/entity"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) Handle(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)
	return &entity.User{
		Username: claims["id"].(string),
	}
}
