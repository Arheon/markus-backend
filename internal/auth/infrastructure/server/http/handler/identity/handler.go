package identity

import (
	domainEntity "github.com/Arheon/markus-backend/internal/auth/domain/entity"
	sharedEntity "github.com/Arheon/markus-backend/internal/shared/domain/entity"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) Handle(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)

	if identity, ok := claims["id"]; ok {
		return &domainEntity.User{
			User: sharedEntity.User{
				ID: identity.(string),
			},
		}
	}

	return nil
}
