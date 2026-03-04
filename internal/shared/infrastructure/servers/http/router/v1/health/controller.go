package health

import (
	"github.com/Arheon/markus-backend/internal/shared/application/handler/health"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func InitController(injector do.Injector, router *gin.RouterGroup) {
	healthGroup := router.Group("health")

	handler := health.NewHandler(injector)
	healthGroup.GET("/:type", handler.Handle)
}
