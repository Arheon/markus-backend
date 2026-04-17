package http

import (
	"github.com/Arheon/markus-backend/internal/server/domain/repository"
	createroom "github.com/Arheon/markus-backend/internal/server/infrastructure/server/http/handler/create_room"
	createnewserver "github.com/Arheon/markus-backend/internal/server/infrastructure/server/http/handler/create_server"
	getuserservers "github.com/Arheon/markus-backend/internal/server/infrastructure/server/http/handler/get_user_servers"
	"github.com/Arheon/markus-backend/pkg/outbox"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

func InitController(injector do.Injector, router *gin.RouterGroup) error {
	logger, err := do.InvokeAs[*logrus.Logger](injector)
	if err != nil {
		return err
	}

	authMiddleware, err := do.InvokeAs[*jwt.GinJWTMiddleware](injector)
	if err != nil {
		return err
	}

	publisher, err := do.InvokeAs[*outbox.Publisher](injector)
	if err != nil {
		return err
	}

	serverRepo, err := do.InvokeAs[repository.ServerRepository](injector)
	if err != nil {
		return err
	}

	roomRepo, err := do.InvokeAs[repository.RoomRepository](injector)

	serverGroup := router.Group("server", authMiddleware.MiddlewareFunc())
	{
		createNewServerHandler := createnewserver.NewHandler(serverRepo, publisher)
		serverGroup.POST("/:userID", createNewServerHandler.Handle)
	}
	{
		getUserServersHandler := getuserservers.NewHandler(serverRepo)
		serverGroup.GET("/:userID", getUserServersHandler.Handle)
	}
	roomGroup := serverGroup.Group("room", authMiddleware.MiddlewareFunc())
	{
		createRoomHandler := createroom.NewHandler(roomRepo, publisher, logger)
		roomGroup.POST("/", createRoomHandler.Handle)
	}

	return nil
}
