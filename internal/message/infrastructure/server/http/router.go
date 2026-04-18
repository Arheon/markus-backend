package http

import (
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	createnewmessage "github.com/Arheon/markus-backend/internal/message/infrastructure/server/http/handler/create_new_message"
	getroommessages "github.com/Arheon/markus-backend/internal/message/infrastructure/server/http/handler/get_room_messages"
	"github.com/Arheon/markus-backend/pkg/outbox"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
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

	logger, err := do.InvokeAs[*logrus.Logger](injector)
	if err != nil {
		return err
	}

	messageRepo, err := do.InvokeAs[repository.MessageRepository](injector)
	if err != nil {
		return err
	}

	messageGroup := router.Group("message", authMiddleware.MiddlewareFunc())
	{
		createNewMessageHandler := createnewmessage.NewHandler(publisher, logger, messageRepo)
		messageGroup.POST("", createNewMessageHandler.Handle)
	}
	{
		getRoomMessagesHandler := getroommessages.NewHandler(logger, messageRepo)
		messageGroup.GET("/all", getRoomMessagesHandler.Handle)
	}

	return nil
}
