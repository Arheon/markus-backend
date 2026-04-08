package server

import (
	domainRepository "github.com/Arheon/markus-backend/internal/server/domain/repository"
	createnewserver "github.com/Arheon/markus-backend/internal/server/infrastructure/server/http/handler/create_server"
	getuserservers "github.com/Arheon/markus-backend/internal/server/infrastructure/server/http/handler/get_user_servers"
	"github.com/Arheon/markus-backend/pkg/outbox"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitController(injector do.Injector, router *gin.RouterGroup) error {
	authMiddleware, err := do.InvokeAs[*jwt.GinJWTMiddleware](injector)
	if err != nil {
		return err
	}

	publisher, err := do.InvokeAs[*outbox.Publisher](injector)
	if err != nil {
		return err
	}

	serverRepo, err := do.InvokeAs[domainRepository.ServerRepository](injector)
	if err != nil {
		return err
	}

	serverGroup := router.Group("server", authMiddleware.MiddlewareFunc())
	{
		createNewServerHandler := createnewserver.NewHandler(serverRepo, publisher)
		serverGroup.POST("/:userID", createNewServerHandler.Handle)
	}
	{
		getUserServersHandler := getuserservers.NewHandler(serverRepo)
		serverGroup.GET("/:userID", getUserServersHandler.Handle)
	}

	return nil
}
