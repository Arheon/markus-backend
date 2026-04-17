package http

import (
	createnewmessage "github.com/Arheon/markus-backend/internal/message/infrastructure/server/http/handler/create_new_message"
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

	messageGroup := router.Group("message", authMiddleware.MiddlewareFunc())
	{
		createNewMessageHandler := createnewmessage.NewHandler(publisher)
		messageGroup.POST("/:roomID", createNewMessageHandler.Handle)
	}

	return nil
}
