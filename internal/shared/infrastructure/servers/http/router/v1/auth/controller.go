package auth

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitContrioller(injector do.Injector, router *gin.RouterGroup) error {
	authMiddleware, err := do.InvokeAs[*jwt.GinJWTMiddleware](injector)
	if err != nil {
		return err
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/refresh", authMiddleware.RefreshHandler)

	auth := router.Group("auth", authMiddleware.MiddlewareFunc())
	{
		auth.POST("/logout", authMiddleware.LogoutHandler)
	}

	return nil
}
