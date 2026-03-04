package user

import (
	"github.com/Arheon/markus-backend/internal/user/application/handler/registration"
	userinfo "github.com/Arheon/markus-backend/internal/user/application/handler/user_info"
	"github.com/Arheon/markus-backend/internal/user/domain/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

func InitContrioller(injector do.Injector, router *gin.RouterGroup) error {
	authMiddleware, err := do.InvokeAs[*jwt.GinJWTMiddleware](injector)
	if err != nil {
		return err
	}

	logger, err := do.InvokeAs[*logrus.Logger](injector)
	if err != nil {
		return err
	}

	userRepository, err := do.InvokeAs[repository.UserRepository](injector)
	if err != nil {
		return err
	}

	registrationHandler := registration.NewHandler(logger, userRepository)
	router.POST("/registration", registrationHandler.Handle)
	{
		getUserHandler := userinfo.NewHandler(logger, userRepository)
		userGroup := router.Group("/user", authMiddleware.MiddlewareFunc())
		userGroup.GET("/:userID", getUserHandler.Handle)
	}

	return nil
}
