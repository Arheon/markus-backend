package server

import (
	createnewserver "github.com/Arheon/markus-backend/internal/server/application/handler/create_new_server"
	getuserservers "github.com/Arheon/markus-backend/internal/server/application/handler/get_user_servers"
	domainRepository "github.com/Arheon/markus-backend/internal/server/domain/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitController(injector do.Injector, router *gin.RouterGroup) error {
	authMiddleware, err := do.InvokeAs[*jwt.GinJWTMiddleware](injector)
	if err != nil {
		return err
	}

	serverGroup := router.Group("server", authMiddleware.MiddlewareFunc())

	serverRepo, err := do.InvokeAs[domainRepository.ServerRepository](injector)
	if err != nil {
		return err
	}

	{
		createNewServerHandler := createnewserver.NewHandler(serverRepo)
		serverGroup.POST("/:userID", createNewServerHandler.Handle)
	}
	{
		getUserServersHandler := getuserservers.NewHandler(serverRepo)
		serverGroup.GET("/:userID", getUserServersHandler.Handle)
	}

	return nil
}
